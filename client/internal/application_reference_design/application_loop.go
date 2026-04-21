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
package application_reference_design

import (
	"context"
	"fmt"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/cddlogger"
	cddsdkgo "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"
)

// ApplicationLoop drives the TR-12 lifecycle for a device.
type ApplicationLoop struct {
	callbacks      DeviceCallbacks
	shim           *Tr12Shim
	sdk            *SDKClient
	registration   *cddsdkgo.DeviceRegistration

	// latestDeviceConfigId tracks the last applied DeviceConfiguration.configurationId.
	// This ID only bumps when device-level simpleSettings change.
	latestDeviceConfigId string

	// latestChannelConfigIds tracks the last applied ChannelConfiguration.configurationId per channel.
	// Each channel's ID bumps independently when that channel's state/settings/connection change.
	latestChannelConfigIds map[string]string

	log *cddlogger.CDDLogger

	// StateCallback is called after each connect response with the current state details.
	StateCallback func(state, pairingCode, deviceID string)

	// ConfigAppliedCallback is called after any configuration update is applied.
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
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		resp, err := l.sdk.Connect(hostID, l.registration)
		if err != nil {
			l.logf("[LOOP] connect error: %v", err)
			select {
			case <-ctx.Done():
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
			l.processConfiguration()
			l.reportStatus()
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(3 * time.Second):
		}
	}
}

// Disconnect calls the SDK disconnect endpoint.
func (l *ApplicationLoop) Disconnect() {
	if _, err := l.sdk.Disconnect(); err != nil {
		l.logf("[LOOP] disconnect error: %v", err)
	}
}

// processConfiguration checks DeviceConfiguration and each ChannelConfiguration independently
// using their configurationId fields to determine what needs to be applied.
//
// DeviceConfiguration.configurationId — bumped by host when device-level simpleSettings change.
// ChannelConfiguration.configurationId — bumped by host when that channel's state/settings/connection change.
// Each is tracked and applied independently.
func (l *ApplicationLoop) processConfiguration() {
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

	l.logf("[LOOP] get_configuration deviceConfigId=%s latestDeviceConfigId=%s channels=%d",
		cfg.ConfigurationId, l.latestDeviceConfigId, len(cfg.Channels))

	anyApplied := false

	// --- Per-channel: apply if ChannelConfiguration.configurationId changed ---
	for _, ch := range cfg.Channels {
		lastId, seen := l.latestChannelConfigIds[ch.Id]
		if seen && ch.ConfigurationId == lastId {
			continue
		}
		l.logf("[LOOP] applying channel %s configurationId=%s (was %s)", ch.Id, ch.ConfigurationId, lastId)
		l.shim.applyChannel(ch)
		l.latestChannelConfigIds[ch.Id] = ch.ConfigurationId
		anyApplied = true
	}

	// --- Device-level: apply standardSettings if DeviceConfiguration.configurationId changed ---
	if cfg.ConfigurationId != l.latestDeviceConfigId {
		if len(cfg.StandardSettings) > 0 {
			l.logf("[LOOP] applying device standardSettings (configurationId %s → %s)",
				l.latestDeviceConfigId, cfg.ConfigurationId)
			for _, kv := range cfg.StandardSettings {
				l.callbacks.UpdateDeviceKeyValue(kv.Key, kv.Value)
			}
			anyApplied = true
		}
		// Always track the device configurationId so we don't re-process it
		l.latestDeviceConfigId = cfg.ConfigurationId
	}

	if !anyApplied {
		return
	}

	if l.ConfigAppliedCallback != nil {
		// Build a composite ID from all applied configurationIds so the callback
		// fires whenever any entity (device or channel) was updated.
		composite := cfg.ConfigurationId
		for _, ch := range cfg.Channels {
			composite += ":" + ch.ConfigurationId
		}
		l.ConfigAppliedCallback(composite)
	}

	// Report actual configuration, echoing the configurationIds the ARD actually applied
	actual := l.shim.GetActualConfiguration(l.registration, cfg, l.latestChannelConfigIds)
	l.logf("[LOOP] reporting actual configuration")

	reportResp, err := l.sdk.ReportActualConfiguration(actual)
	if err != nil {
		l.logf("[LOOP] report_actual_configuration error: %v", err)
		return
	}
	l.logf("[LOOP] report_actual_configuration state=%s message=%s", reportResp.State, reportResp.Message)
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
