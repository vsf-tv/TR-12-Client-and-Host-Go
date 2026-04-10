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
// TR-12 Model Shim — walks registration/configuration structures and dispatches
// to device-specific callbacks. Mirrors tr12_shim.py.
package application_reference_design

import (
	"fmt"

	cddsdkgo "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"
)

// Tr12Shim bridges TR-12 model structures to device callbacks.
type Tr12Shim struct {
	CB DeviceCallbacks
}

// NewTr12Shim creates a new shim with ARD callbacks (for the ARD binary).
func NewTr12Shim() *Tr12Shim {
	return &Tr12Shim{CB: NewArdCallbacks()}
}

// NewTr12ShimWithCallbacks creates a shim with a custom DeviceCallbacks implementation.
func NewTr12ShimWithCallbacks(cb DeviceCallbacks) *Tr12Shim {
	return &Tr12Shim{CB: cb}
}

// ApplyDesiredConfiguration walks a desired DeviceConfiguration and pushes all values to the device.
// For selective per-channel application, use applyChannel directly from ApplicationLoop.
func (s *Tr12Shim) ApplyDesiredConfiguration(desired *cddsdkgo.DeviceConfiguration) bool {
	if desired == nil {
		return false
	}
	for _, kv := range desired.SimpleSettings {
		s.CB.UpdateDeviceKeyValue(kv.Key, kv.Value)
	}
	for _, ch := range desired.Channels {
		s.applyChannel(ch)
	}
	return true
}

// applyChannel applies a single channel's desired configuration to the device.
func (s *Tr12Shim) applyChannel(chCfg cddsdkgo.ChannelConfiguration) {
	chID := chCfg.Id

	// Settings (SimpleSettings or Profile via oneOf union)
	if chCfg.Settings != nil {
		settings := chCfg.Settings
		if settings.SimpleSettings != nil {
			for _, kv := range settings.SimpleSettings.SimpleSettings {
				s.CB.UpdateChannelSettings(chID, kv.Key, kv.Value)
			}
		} else if settings.Profile != nil {
			s.CB.UpdateChannelProfile(chID, settings.Profile.Profile.Id)
		}
	}

	// Connection
	if chCfg.Connection != nil {
		s.CB.UpdateChannelConnection(chID, chCfg.Connection)
	}

	// State (apply last so settings/connection are in place first)
	s.CB.UpdateChannelState(chID, chCfg.State)
}

// GetActualConfiguration reads back current values using the registration as a template.
// appliedChannelIds contains the configurationId last applied per channel by the ARD.
// The device-level configurationId is echoed from desired.
func (s *Tr12Shim) GetActualConfiguration(registration *cddsdkgo.DeviceRegistration, desired *cddsdkgo.DeviceConfiguration, appliedChannelIds map[string]string) *cddsdkgo.DeviceConfiguration {
	result := &cddsdkgo.DeviceConfiguration{}

	// Echo device-level configurationId from desired
	if desired != nil {
		result.ConfigurationId = desired.ConfigurationId
	}

	// Device-level health
	if health := s.CB.GetDeviceHealth(); health != nil {
		result.Health = health
	}

	// Device-level simple settings
	var deviceSettings []cddsdkgo.IdAndValue
	for _, setting := range registration.SimpleSettings {
		if val, found := s.CB.GetDeviceUpdatedValue(setting.Id); found {
			deviceSettings = append(deviceSettings, cddsdkgo.IdAndValue{
				Key: setting.Id, Value: val,
			})
		}
	}
	if len(deviceSettings) > 0 {
		result.SimpleSettings = deviceSettings
	}

	// Channels — echo back the configurationId the ARD actually applied per channel
	var channels []cddsdkgo.ChannelConfiguration
	for _, regCh := range registration.Channels {
		chCfg := s.buildChannelConfig(regCh)
		chCfg.ConfigurationId = appliedChannelIds[regCh.Id]
		channels = append(channels, chCfg)
	}
	result.Channels = channels
	return result
}

func (s *Tr12Shim) buildChannelConfig(regCh cddsdkgo.Channel) cddsdkgo.ChannelConfiguration {
	chID := regCh.Id
	chCfg := cddsdkgo.ChannelConfiguration{
		Id:    chID,
		State: s.CB.GetChannelState(chID),
	}

	// Health — report device health for this channel
	if health := s.CB.GetChannelHealth(chID); health != nil {
		chCfg.Health = health
	}

	// Check profiles first
	hasProfile := false
	if len(regCh.Profiles) > 0 {
		if profileID, found := s.CB.GetChannelProfileValue(chID); found {
			chCfg.Settings = &cddsdkgo.SettingsChoice{
				Profile: &cddsdkgo.Profile{
					Profile: cddsdkgo.SettingProfile{Id: profileID},
				},
			}
			hasProfile = true
		}
	}

	// Simple settings if no profile
	if !hasProfile && len(regCh.SimpleSettings) > 0 {
		var kvList []cddsdkgo.IdAndValue
		for _, setting := range regCh.SimpleSettings {
			if val, found := s.CB.GetChannelUpdatedValue(chID, setting.Id); found {
				kvList = append(kvList, cddsdkgo.IdAndValue{
					Key: setting.Id, Value: val,
				})
			}
		}
		if len(kvList) > 0 {
			chCfg.Settings = &cddsdkgo.SettingsChoice{
				SimpleSettings: &cddsdkgo.SimpleSettings{
					SimpleSettings: kvList,
				},
			}
		}
	}

	// Connection
	conn := s.CB.GetChannelConnection(chID)
	if conn != nil {
		chCfg.Connection = conn
	}

	return chCfg
}

// GetDeviceStatus builds the typed device status payload using registration channels.
func (s *Tr12Shim) GetDeviceStatus(registration *cddsdkgo.DeviceRegistration) *cddsdkgo.DeviceStatus {
	var channelStatuses []cddsdkgo.ChannelStatus
	for _, regCh := range registration.Channels {
		channelStatuses = append(channelStatuses, cddsdkgo.ChannelStatus{
			Id:     regCh.Id,
			State:  s.CB.GetChannelState(regCh.Id),
			Status: s.CB.GetChannelStatus(regCh.Id),
		})
	}
	return &cddsdkgo.DeviceStatus{
		Status:   s.CB.GetDeviceStatus(),
		Channels: channelStatuses,
	}
}

// PrintActualConfig is a debug helper.
func (s *Tr12Shim) PrintActualConfig(registration *cddsdkgo.DeviceRegistration) {
	actual := s.GetActualConfiguration(registration, nil, map[string]string{})
	fmt.Printf("[SHIM TEST] get_actual_configuration: %+v\n", actual)
}
