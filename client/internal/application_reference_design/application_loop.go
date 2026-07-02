// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// ApplicationLoop — the reusable TR-12 connect/status/config run loop.
//
// Design: per-channel work is strictly version-gated.
//
//   The main loop only dispatches a goroutine for a channel when its config
//   version changes. There is no persistent state-watching loop between version
//   changes — the host and operator control the channel lifecycle.
//
//   For desired=ACTIVE: goroutine stops the channel if running (waits for fully
//   stopped), applies settings and protocol (safe because channel is stopped),
//   then issues start and walks away. The device is responsible for reaching
//   "started"; we do not wait or retry.
//
//   For desired=IDLE: goroutine issues stop and walks away.
//
//   A newer version dispatched while a goroutine is running cancels the in-flight
//   goroutine so it never operates against stale settings.
//
//   Device-level standard settings are applied synchronously in the main loop
//   (no state dependency) and reported immediately.
package application_reference_design

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/cddlogger"
	cddsdkgo "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"
)

// pollInterval is how often stopAndWait checks whether the channel has stopped.
const pollInterval = 500 * time.Millisecond

// stopTimeout is the maximum time stopAndWait will wait for a channel to reach
// "stopped" before giving up. If exceeded, the channel is marked DEGRADED and
// the config apply is aborted for this cycle.
const stopTimeout = 5 * time.Second

// RawStatusGetter is an optional interface that callbacks can implement to
// return the device's raw channel status string ("stopped", "stopping",
// "starting", "started"). Used by stopAndWait to correctly distinguish the
// fully-stopped state from a mid-transition "stopping" state.
// If not implemented, GetChannelState returning IDLE is used as the fallback.
type RawStatusGetter interface {
	GetChannelRawStatus(channelID string) string
}

// ApplicationLoop drives the TR-12 lifecycle for a device.
type ApplicationLoop struct {
	callbacks    DeviceCallbacks
	shim         *Tr12Shim
	sdk          *SDKClient
	registration *cddsdkgo.DeviceRegistration

	mu                     sync.Mutex
	latestDeviceConfigId   string
	latestChannelConfigIds map[string]string
	channelWorkers         map[string]context.CancelFunc // at most one in-flight goroutine per channel
	lastSeenConfig         *cddsdkgo.DesiredDeviceConfiguration // for actual-config reports from goroutines

	reportedInitialActualConfig bool

	log *cddlogger.CDDLogger

	StateCallback         func(state, pairingCode, deviceID string)
	ConfigAppliedCallback func(deviceConfigId string)
}

// NewApplicationLoop creates a loop with the given callbacks and registration.
func NewApplicationLoop(sdkURL string, callbacks DeviceCallbacks, registration *cddsdkgo.DeviceRegistration) *ApplicationLoop {
	return &ApplicationLoop{
		callbacks:              callbacks,
		shim:                   NewTr12ShimWithCallbacks(callbacks),
		sdk:                    NewSDKClient(sdkURL),
		registration:           registration,
		latestChannelConfigIds: make(map[string]string),
		channelWorkers:         make(map[string]context.CancelFunc),
	}
}

// SetLogger attaches a logger to the loop. Call before Run().
func (l *ApplicationLoop) SetLogger(log *cddlogger.CDDLogger) {
	l.log = log
}

func (l *ApplicationLoop) logf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	if l.log != nil {
		l.log.Info(msg)
	} else {
		fmt.Println(msg)
	}
}

// Run executes the loop until ctx is cancelled.
func (l *ApplicationLoop) Run(ctx context.Context, hostID string) {
	wasConnected := false
	for {
		select {
		case <-ctx.Done():
			l.cancelAllWorkers()
			return
		default:
		}

		resp, err := l.sdk.Connect(hostID, l.registration)
		if err != nil {
			l.logf("[LOOP] connect error: %v", err)
			if wasConnected {
				l.logf("[LOOP] lost connection — resetting config state")
				l.resetConfigState()
				wasConnected = false
			}
			select {
			case <-ctx.Done():
				l.cancelAllWorkers()
				return
			case <-time.After(3 * time.Second):
			}
			continue
		}

		l.logf("[LOOP] state=%s deviceId=%s", resp.State, resp.GetDeviceId())

		if resp.State == "PAIRING" {
			l.logf("[LOOP] pairing code: %s (expires in %ds)", resp.GetPairingCode(), int(resp.GetExpiresSeconds()))
		}

		if l.StateCallback != nil {
			l.StateCallback(resp.State, resp.GetPairingCode(), resp.GetDeviceId())
		}

		if resp.State == "CONNECTED" {
			wasConnected = true
			l.reportInitialActualConfig()
			l.processConfiguration(ctx)
			l.reportStatus()
		} else if wasConnected {
			l.logf("[LOOP] transitioned away from CONNECTED (now %s) — resetting config state", resp.State)
			l.resetConfigState()
			wasConnected = false
		}

		select {
		case <-ctx.Done():
			l.cancelAllWorkers()
			return
		case <-time.After(3 * time.Second):
		}
	}
}

// processConfiguration fetches desired config and, for each channel whose
// version has changed, dispatches a per-channel goroutine to apply the change.
// Channels whose version is unchanged are left completely untouched.
// Device-level settings are applied synchronously and reported immediately.
func (l *ApplicationLoop) processConfiguration(ctx context.Context) {
	resp, err := l.sdk.GetConfiguration()
	if err != nil {
		l.logf("[LOOP] get_configuration error: %v", err)
		return
	}
	if resp.Configuration == nil || resp.Configuration.Payload == nil {
		l.logf("[LOOP] get_configuration: no configuration yet")
		return
	}

	cfg := resp.Configuration.Payload

	l.mu.Lock()
	l.lastSeenConfig = cfg
	latestChannelConfigIds := make(map[string]string, len(l.latestChannelConfigIds))
	for k, v := range l.latestChannelConfigIds {
		latestChannelConfigIds[k] = v
	}
	latestDeviceConfigId := l.latestDeviceConfigId
	l.mu.Unlock()

	l.logf("[LOOP] get_configuration deviceVersion=%s channels=%d", cfg.Version, len(cfg.Channels))

	channelDispatched := false
	for _, ch := range cfg.Channels {
		if latestChannelConfigIds[ch.Id] == ch.Version {
			continue // version unchanged — nothing to do
		}
		l.logf("[LOOP] channel %s version %s → %s desired=%s",
			ch.Id, latestChannelConfigIds[ch.Id], ch.Version, ch.State)
		l.dispatchChannelWork(ctx, ch)
		channelDispatched = true
	}

	// Device-level standard settings have no state dependency — apply
	// synchronously in the main loop and report immediately.
	deviceSettingsApplied := false
	if cfg.Version != latestDeviceConfigId {
		if len(cfg.StandardSettings) > 0 {
			l.logf("[LOOP] applying device standardSettings (version %s → %s)",
				latestDeviceConfigId, cfg.Version)
			for _, kv := range cfg.StandardSettings {
				l.callbacks.UpdateDeviceKeyValue(kv.Id, kv.Value)
			}
			deviceSettingsApplied = true
		}
		l.mu.Lock()
		l.latestDeviceConfigId = cfg.Version
		l.mu.Unlock()
	}

	if !channelDispatched && !deviceSettingsApplied {
		return
	}

	// If only device settings changed (no channel goroutines in flight), report
	// now. Channel goroutines report actual config when they complete.
	if deviceSettingsApplied && !channelDispatched {
		l.reportActual(cfg)
	}
}

// dispatchChannelWork cancels any in-flight goroutine for a channel and starts
// a fresh one to apply the new desired configuration.
func (l *ApplicationLoop) dispatchChannelWork(ctx context.Context, ch cddsdkgo.DesiredChannelConfiguration) {
	l.mu.Lock()
	if cancel, ok := l.channelWorkers[ch.Id]; ok {
		cancel()
	}
	workerCtx, cancel := context.WithCancel(ctx)
	l.channelWorkers[ch.Id] = cancel
	l.mu.Unlock()

	go l.channelWorkerFn(workerCtx, ch)
}

// channelWorkerFn applies a new desired configuration for one channel.
//
//   desired=ACTIVE: stops the channel if not already stopped (waits for stopped),
//   applies settings and protocol, then issues start (fire and forget).
//
//   desired=IDLE: issues stop and exits immediately (fire and forget).
//
// In both cases: updates latestChannelConfigIds and reports actual config when done.
// Context cancellation (new version arriving) exits early without updating the version —
// the next cycle will re-dispatch for the newer version.
func (l *ApplicationLoop) channelWorkerFn(ctx context.Context, ch cddsdkgo.DesiredChannelConfiguration) {
	chID := ch.Id

	if ch.State == cddsdkgo.CHANNELSTATE_ACTIVE {
		// Device requires channel to be fully stopped before settings can be applied.
		// Use raw status when available — GetChannelState returns IDLE for both
		// "stopped" and "stopping", so we must check the raw value to avoid
		// applying settings or issuing start while the channel is mid-stop.
		if !l.isFullyStopped(chID) {
			l.logf("[LOOP] channel %s stopping before settings apply", chID)
			if !l.stopAndWait(ctx, chID) {
			if ctx.Err() != nil {
				l.logf("[LOOP] channel %s stop aborted (new version or shutdown)", chID)
			} else {
				// Timed out — channel is stuck stopping. Mark DEGRADED and report
				// so the host is informed; do not apply settings to a stuck channel.
				l.logf("[LOOP] channel %s failed to stop within %s — aborting config apply", chID, stopTimeout)
				msg := fmt.Sprintf("channel failed to stop within %s — config not applied", stopTimeout)
				if cb, ok := l.callbacks.(interface{ SetChannelDegradedHealth(string, string) }); ok {
					cb.SetChannelDegradedHealth(chID, msg)
				}
				l.mu.Lock()
				cfg := l.lastSeenConfig
				appliedVersions := make(map[string]string, len(l.latestChannelConfigIds))
				for k, v := range l.latestChannelConfigIds {
					appliedVersions[k] = v
				}
				l.mu.Unlock()
				if cfg != nil {
					actual := l.shim.GetActualConfiguration(l.registration, cfg, appliedVersions)
					l.sdk.ReportActualConfiguration(actual)
				}
			}
			return
		}
		}

		if ctx.Err() != nil {
			return
		}

		// Channel is stopped — safe to apply settings and protocol.
		if cb, ok := l.callbacks.(interface{ BeginChannelUpdate(string) }); ok {
			cb.BeginChannelUpdate(chID)
		}
		l.applyChannelConfigSync(ch)

		// Evaluate settings-apply health (e.g. Talon API errors during PATCH calls).
		if cb, ok := l.callbacks.(interface{ EvalChannelHealth(string) }); ok {
			cb.EvalChannelHealth(chID)
		}

		if ctx.Err() != nil {
			return
		}

		// Issue start — fire and forget. The device transitions on its own.
		l.logf("[LOOP] channel %s issuing start", chID)
		l.issueChannelState(ctx, chID, cddsdkgo.CHANNELSTATE_ACTIVE)

	} else {
		// Desired = IDLE: ensure stopped, apply settings (device requires stopped
		// state to accept settings), then walk away. Settings are applied now so
		// they take effect on the next start without needing another push.
		if !l.isFullyStopped(chID) {
			l.logf("[LOOP] channel %s stopping before settings apply (desired=IDLE)", chID)
			if !l.stopAndWait(ctx, chID) {
				l.logf("[LOOP] channel %s stop aborted", chID)
				return
			}
		}
		if cb, ok := l.callbacks.(interface{ BeginChannelUpdate(string) }); ok {
			cb.BeginChannelUpdate(chID)
		}
		l.applyChannelConfigSync(ch)
		if cb, ok := l.callbacks.(interface{ EvalChannelHealth(string) }); ok {
			cb.EvalChannelHealth(chID)
		}
	}

	if ctx.Err() != nil {
		return
	}

	// Mark this version as handled and report actual config.
	l.mu.Lock()
	l.latestChannelConfigIds[chID] = ch.Version
	cfg := l.lastSeenConfig
	appliedVersions := make(map[string]string, len(l.latestChannelConfigIds))
	for k, v := range l.latestChannelConfigIds {
		appliedVersions[k] = v
	}
	l.mu.Unlock()

	if l.ConfigAppliedCallback != nil && cfg != nil {
		composite := cfg.Version
		for _, c := range cfg.Channels {
			if v, ok := appliedVersions[c.Id]; ok {
				composite += ":" + v
			}
		}
		l.ConfigAppliedCallback(composite)
	}

	actual := l.shim.GetActualConfiguration(l.registration, cfg, appliedVersions)
	l.logf("[LOOP] channel %s reporting actual configuration", chID)
	if reportResp, err := l.sdk.ReportActualConfiguration(actual); err != nil {
		l.logf("[LOOP] channel %s report_actual error: %v", chID, err)
	} else {
		l.logf("[LOOP] channel %s report_actual state=%s message=%s",
			chID, reportResp.State, reportResp.Message)
	}
}

// isFullyStopped returns true when the channel is confirmed stopped (not stopping).
func (l *ApplicationLoop) isFullyStopped(chID string) bool {
	if rawGetter, hasRaw := l.callbacks.(RawStatusGetter); hasRaw {
		return rawGetter.GetChannelRawStatus(chID) == "stopped"
	}
	return l.callbacks.GetChannelState(chID) == cddsdkgo.CHANNELSTATE_IDLE
}

// stopAndWait issues a stop command and polls until the channel is fully stopped,
// the context is cancelled, or stopTimeout is exceeded.
// Returns true when stopped. Returns false on context cancellation (caller should
// check ctx.Err() == nil to detect timeout vs. cancellation).
func (l *ApplicationLoop) stopAndWait(ctx context.Context, chID string) bool {
	l.issueChannelState(ctx, chID, cddsdkgo.CHANNELSTATE_IDLE)

	rawGetter, hasRaw := l.callbacks.(RawStatusGetter)
	deadline := time.Now().Add(stopTimeout)
	for {
		var stopped bool
		if hasRaw {
			stopped = rawGetter.GetChannelRawStatus(chID) == "stopped"
		} else {
			stopped = l.callbacks.GetChannelState(chID) == cddsdkgo.CHANNELSTATE_IDLE
		}
		if stopped {
			return true
		}
		if time.Now().After(deadline) {
			return false // timeout — ctx.Err() will be nil, distinguishing from cancellation
		}
		select {
		case <-ctx.Done():
			return false
		case <-time.After(pollInterval):
		}
	}
}

// issueChannelState calls the context-aware state callback when available,
// falling back to the plain UpdateChannelState.
func (l *ApplicationLoop) issueChannelState(ctx context.Context, chID string, state cddsdkgo.ChannelState) {
	if cb, ok := l.callbacks.(interface {
		UpdateChannelStateWithContext(context.Context, string, cddsdkgo.ChannelState)
	}); ok {
		cb.UpdateChannelStateWithContext(ctx, chID, state)
	} else {
		l.callbacks.UpdateChannelState(chID, state)
	}
}

// applyChannelConfigSync applies settings and protocol for one channel.
// Must be called when the channel is stopped (device requirement for Osprey).
func (l *ApplicationLoop) applyChannelConfigSync(ch cddsdkgo.DesiredChannelConfiguration) {
	chID := ch.Id
	if ch.ChannelSettings != nil {
		if ch.ChannelSettings.StandardSettings != nil {
			for _, kv := range ch.ChannelSettings.StandardSettings.StandardSettings {
				l.callbacks.UpdateChannelSettings(chID, kv.Id, kv.Value)
			}
		} else if ch.ChannelSettings.Profile != nil {
			l.callbacks.UpdateChannelProfile(chID, ch.ChannelSettings.Profile.Profile.Id)
		}
	}
	if ch.Protocol != nil {
		l.callbacks.UpdateChannelConnection(chID, ch.Protocol)
	}
}

// reportActual builds and sends an actual config report to the SDK.
func (l *ApplicationLoop) reportActual(cfg *cddsdkgo.DesiredDeviceConfiguration) {
	l.mu.Lock()
	appliedVersions := make(map[string]string, len(l.latestChannelConfigIds))
	for k, v := range l.latestChannelConfigIds {
		appliedVersions[k] = v
	}
	l.mu.Unlock()

	actual := l.shim.GetActualConfiguration(l.registration, cfg, appliedVersions)
	l.logf("[LOOP] reporting actual configuration")
	if reportResp, err := l.sdk.ReportActualConfiguration(actual); err != nil {
		l.logf("[LOOP] report_actual_configuration error: %v", err)
	} else {
		l.logf("[LOOP] report_actual_configuration state=%s message=%s",
			reportResp.State, reportResp.Message)
	}
}

// cancelAllWorkers cancels all in-flight channel goroutines.
func (l *ApplicationLoop) cancelAllWorkers() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, cancel := range l.channelWorkers {
		cancel()
	}
	l.channelWorkers = make(map[string]context.CancelFunc)
}

func (l *ApplicationLoop) resetConfigState() {
	l.cancelAllWorkers()
	l.mu.Lock()
	l.latestDeviceConfigId = ""
	l.latestChannelConfigIds = make(map[string]string)
	l.lastSeenConfig = nil
	l.mu.Unlock()
	l.reportedInitialActualConfig = false
}

// Disconnect calls the SDK disconnect endpoint.
func (l *ApplicationLoop) Disconnect() {
	if _, err := l.sdk.Disconnect(); err != nil {
		l.logf("[LOOP] disconnect error: %v", err)
	}
}

func (l *ApplicationLoop) reportInitialActualConfig() {
	if l.reportedInitialActualConfig {
		return
	}
	l.mu.Lock()
	appliedVersions := make(map[string]string, len(l.latestChannelConfigIds))
	for k, v := range l.latestChannelConfigIds {
		appliedVersions[k] = v
	}
	l.mu.Unlock()

	actual := l.shim.GetActualConfiguration(l.registration, nil, appliedVersions)
	l.logf("[LOOP] reporting initial actual configuration")
	if _, err := l.sdk.ReportActualConfiguration(actual); err != nil {
		l.logf("[LOOP] initial report_actual_configuration error: %v", err)
		return
	}
	l.reportedInitialActualConfig = true
}

func (l *ApplicationLoop) reportStatus() {
	status := l.shim.GetDeviceStatus(l.registration)
	resp, err := l.sdk.ReportStatus(status)
	if err != nil {
		l.logf("[LOOP] report_status error: %v", err)
		return
	}
	l.logf("[LOOP] report_status state=%s", resp.State)
}
