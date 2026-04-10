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
// TR-12 Callback implementations — mirrors tr12_callbacks.py.
// Callbacks are stateless delegates — all configuration state lives in Encoder.
package application_reference_design

import (
	"fmt"

	cddsdkgo "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"
)

// ArdCallbacks implements DeviceCallbacks for the Application Reference Design.
// It delegates all state to the Encoder, which is the source of truth for
// current device configuration — mirroring how a real device integration works.
type ArdCallbacks struct {
	Encoder *Encoder
}

// NewArdCallbacks creates ARD callbacks with a new encoder instance.
func NewArdCallbacks() *ArdCallbacks {
	return &ArdCallbacks{Encoder: NewEncoder()}
}

// Ensure ArdCallbacks implements DeviceCallbacks at compile time.
var _ DeviceCallbacks = (*ArdCallbacks)(nil)

// --- Update (set) callbacks ---

func (cb *ArdCallbacks) UpdateDeviceKeyValue(key, value string) {
	fmt.Printf("[UPDATE] Device setting: %s = %s\n", key, value)
	cb.Encoder.SetDeviceSetting(key, value)
}

func (cb *ArdCallbacks) UpdateChannelSettings(channelID, key, value string) {
	fmt.Printf("[UPDATE] Channel %s setting: %s = %s\n", channelID, key, value)
	cb.Encoder.SetChannelSetting(channelID, key, value)
}

func (cb *ArdCallbacks) UpdateChannelProfile(channelID, profileID string) {
	fmt.Printf("[UPDATE] Channel %s profile: %s\n", channelID, profileID)
}

func (cb *ArdCallbacks) UpdateChannelConnection(channelID string, connection *cddsdkgo.Connection) {
	fmt.Printf("[UPDATE] Channel %s connection: %+v\n", channelID, connection)
	cb.Encoder.HandleTransportConfigChange(channelID, connection)
}

func (cb *ArdCallbacks) UpdateChannelState(channelID string, state cddsdkgo.ChannelState) {
	fmt.Printf("[UPDATE] Channel %s state: %s\n", channelID, state)
	cb.Encoder.HandleUpdateState(channelID, state)
}

// --- Get (read-back) callbacks ---

func (cb *ArdCallbacks) GetDeviceUpdatedValue(key string) (string, bool) {
	return cb.Encoder.GetDeviceSetting(key)
}

func (cb *ArdCallbacks) GetChannelUpdatedValue(channelID, key string) (string, bool) {
	return cb.Encoder.GetChannelSetting(channelID, key)
}

func (cb *ArdCallbacks) GetChannelProfileValue(channelID string) (string, bool) {
	return "", false // Use simple settings by default
}

func (cb *ArdCallbacks) GetChannelConnection(channelID string) *cddsdkgo.Connection {
	return cb.Encoder.GetChannelConnection(channelID)
}

func (cb *ArdCallbacks) GetChannelState(channelID string) cddsdkgo.ChannelState {
	return cb.Encoder.GetChannelState(channelID)
}

func (cb *ArdCallbacks) GetDeviceStatus() []cddsdkgo.StatusValue {
	// Device status doesn't depend on any specific channel
	return []cddsdkgo.StatusValue{
		{Name: "cpu", Value: "31", Info: "Current CPU % utilization."},
		{Name: "temp", Value: "76", Info: "CPU in degrees C."},
		{Name: "model", Value: "Talon", Info: "Hardware device model identifier."},
		{Name: "serial", Value: "123456789", Info: "Device serial number."},
	}
}

func (cb *ArdCallbacks) GetChannelStatus(channelID string) []cddsdkgo.StatusValue {
	if cb.Encoder.RunningChannel(channelID) {
		return []cddsdkgo.StatusValue{
			{Name: "bitrate", Value: GetSimulatedBitrate(), Info: "Bitrate Mbps configured on the video encoder."},
		}
	}
	return []cddsdkgo.StatusValue{
		{Name: "bitrate", Value: "0", Info: "Bitrate Mbps configured on the video encoder."},
	}
}

func (cb *ArdCallbacks) GetChannelHealth(channelID string) *cddsdkgo.Health {
	return cb.Encoder.GetChannelHealth(channelID)
}

func (cb *ArdCallbacks) GetDeviceHealth() *cddsdkgo.Health {
	h := cddsdkgo.HealthyAsHealth(cddsdkgo.NewHealthy(map[string]interface{}{}))
	return &h
}
