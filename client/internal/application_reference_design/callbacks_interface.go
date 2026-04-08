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
// DeviceCallbacks — the interface a device integrator implements to bridge
// TR-12 model operations to the device's native control API.
package application_reference_design

import cddsdkgo "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"

// DeviceCallbacks is the integration interface between the TR-12 shim and a
// device's native control plane. Implement this interface to make any device
// TR-12 compliant.
//
// The interface has two sides:
//   - Update (set) methods — called when applying a desired configuration from the host
//   - Get (read-back) methods — called when building the actual configuration to report back
type DeviceCallbacks interface {
	// --- Apply (set) side ---

	// UpdateDeviceKeyValue applies a device-level setting (e.g. clock source).
	UpdateDeviceKeyValue(key, value string)

	// UpdateChannelSettings applies a channel-level simple setting.
	UpdateChannelSettings(channelID, key, value string)

	// UpdateChannelProfile applies a profile selection to a channel.
	UpdateChannelProfile(channelID, profileID string)

	// UpdateChannelConnection applies transport protocol configuration to a channel.
	// IMPORTANT: TR-12 communicates desired configuration once. The device is
	// responsible for retrying until the desired state is achieved.
	UpdateChannelConnection(channelID string, connection *cddsdkgo.Connection)

	// UpdateChannelState applies the desired channel state (ACTIVE or IDLE).
	// The device must retry until the state is achieved — the host will not re-send.
	UpdateChannelState(channelID string, state cddsdkgo.ChannelState)

	// --- Read-back (get) side ---

	// GetDeviceUpdatedValue returns the current value of a device-level setting.
	GetDeviceUpdatedValue(key string) (string, bool)

	// GetChannelUpdatedValue returns the current value of a channel-level setting.
	GetChannelUpdatedValue(channelID, key string) (string, bool)

	// GetChannelProfileValue returns the currently active profile ID for a channel,
	// or ("", false) if simple settings are in use.
	GetChannelProfileValue(channelID string) (string, bool)

	// GetChannelConnection returns the current transport protocol configuration.
	GetChannelConnection(channelID string) *cddsdkgo.Connection

	// GetChannelState returns the current channel state.
	GetChannelState(channelID string) cddsdkgo.ChannelState

	// GetDeviceStatus returns device-level status values (CPU, temp, model, serial, etc.)
	GetDeviceStatus() []cddsdkgo.StatusValue

	// GetChannelStatus returns per-channel status values (bitrate, output state, etc.)
	GetChannelStatus(channelID string) []cddsdkgo.StatusValue
}
