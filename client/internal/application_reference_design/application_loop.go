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
// Mirrors the original ARD RunLoop logic exactly.
package application_reference_design

import (
	"context"
	"encoding/json"
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
	latestConfigID string
	log            *cddlogger.CDDLogger
	// StateCallback is called after each connect response with the current state details.
	StateCallback func(state, pairingCode, deviceID string)
	// ConfigAppliedCallback is called after a configuration update is successfully applied.
	ConfigAppliedCallback func(configID string)
}

// NewApplicationLoop creates a loop with the given callbacks and registration.
func NewApplicationLoop(sdkURL string, callbacks DeviceCallbacks, registration *cddsdkgo.DeviceRegistration) *ApplicationLoop {
	return &ApplicationLoop{
		callbacks:    callbacks,
		shim:         NewTr12ShimWithCallbacks(callbacks),
		sdk:          NewSDKClient(sdkURL),
		registration: registration,
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
// Each iteration:
//  1. Call connect — get connection state
//  2. If CONNECTED: get configuration
//     - If updateId changed: apply config to device, read back actual, report actual
//  3. If CONNECTED: report status
func (l *ApplicationLoop) Run(ctx context.Context, hostID string) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// 1. Connect
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
			l.logf("[LOOP] pairing code: %s (expires in %ds)", resp.GetPairingCode(), int(resp.GetExpires()))
		}

		if l.StateCallback != nil {
			l.StateCallback(resp.State, resp.GetPairingCode(), resp.GetDeviceId())
		}

		if resp.State == "CONNECTED" {
			// 2. Get configuration — apply only if updateId changed
			l.processConfiguration()

			// 3. Report status
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

// processConfiguration gets the current configuration from the SDK.
// If the updateId has changed since last processed, applies it to the device
// via the shim callbacks, reads back actual state, and reports it.
// Does NOT report actual configuration if nothing changed.
func (l *ApplicationLoop) processConfiguration() {
	resp, err := l.sdk.GetConfiguration()
	if err != nil {
		l.logf("[LOOP] get_configuration error: %v", err)
		return
	}
	if resp.Configuration == nil {
		l.logf("[LOOP] get_configuration: no configuration yet")
		return
	}

	updateID := resp.Configuration.GetUpdateId()
	l.logf("[LOOP] get_configuration updateId=%s latestId=%s", updateID, l.latestConfigID)

	// Skip if no update or already processed
	if updateID == "" || updateID == l.latestConfigID {
		return
	}

	payload := resp.Configuration.Payload
	if payload == nil {
		l.logf("[LOOP] configuration payload is nil")
		return
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		l.logf("[LOOP] marshal payload error: %v", err)
		return
	}
	l.logf("[LOOP] new configuration update_id=%s payload=%s", updateID, string(payloadBytes))

	var deviceConfig cddsdkgo.DeviceConfiguration
	if err := json.Unmarshal(payloadBytes, &deviceConfig); err != nil {
		l.logf("[LOOP] parse DeviceConfiguration error: %v", err)
		return
	}

	l.logf("[LOOP] applying configuration: channels=%d", len(deviceConfig.Channels))

	// Apply desired configuration to the device via callbacks
	if !l.shim.ApplyDesiredConfiguration(&deviceConfig) {
		l.logf("[LOOP] ApplyDesiredConfiguration returned false")
		return
	}

	// Mark this updateId as processed
	l.latestConfigID = updateID

	if l.ConfigAppliedCallback != nil {
		l.ConfigAppliedCallback(updateID)
	}

	// Read back actual state from the device and report it
	actual := l.shim.GetActualConfiguration(l.registration)
	actualBytes, _ := json.Marshal(actual)
	l.logf("[LOOP] reporting actual configuration: %s", string(actualBytes))

	reportResp, err := l.sdk.ReportActualConfiguration(actual)
	if err != nil {
		l.logf("[LOOP] report_actual_configuration error: %v", err)
		return
	}
	l.logf("[LOOP] report_actual_configuration state=%s message=%s", reportResp.State, reportResp.Message)
}

func (l *ApplicationLoop) reportStatus() {
	status := l.shim.GetDeviceStatus()
	resp, err := l.sdk.ReportStatus(status)
	if err != nil {
		l.logf("[LOOP] report_status error: %v", err)
		return
	}
	l.logf("[LOOP] report_status state=%s", resp.State)
}
