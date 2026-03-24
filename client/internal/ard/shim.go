// Copyright 2025 Amazon.com Inc
// Licensed under the Apache License, Version 2.0
//
// TR-12 Model Shim — walks registration/configuration structures and dispatches
// to device-specific callbacks. Mirrors tr12_shim.py.
package ard

import (
	"fmt"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/pkg/cddmodels"
)

// Tr12Shim bridges TR-12 model structures to device callbacks.
type Tr12Shim struct {
	CB *Callbacks
}

// NewTr12Shim creates a new shim with callbacks.
func NewTr12Shim() *Tr12Shim {
	return &Tr12Shim{CB: NewCallbacks()}
}

// ApplyDesiredConfiguration walks a desired DeviceConfiguration and pushes values to the device.
func (s *Tr12Shim) ApplyDesiredConfiguration(desired *cddmodels.DeviceConfiguration) bool {
	if desired == nil {
		return false
	}

	// Device-level simple settings
	for _, kv := range desired.SimpleSettings {
		s.CB.UpdateDeviceKeyValue(kv.Key, kv.Value)
	}

	// Per-channel configuration
	for _, ch := range desired.Channels {
		s.applyChannel(ch)
	}
	return true
}

func (s *Tr12Shim) applyChannel(chCfg cddmodels.ChannelConfiguration) {
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
func (s *Tr12Shim) GetActualConfiguration(registration *cddmodels.DeviceRegistration) *cddmodels.DeviceConfiguration {
	result := &cddmodels.DeviceConfiguration{}

	// Device-level simple settings
	var deviceSettings []cddmodels.IdAndValue
	for _, setting := range registration.SimpleSettings {
		if val, found := s.CB.GetDeviceUpdatedValue(setting.Id); found {
			deviceSettings = append(deviceSettings, cddmodels.IdAndValue{
				Key: setting.Id, Value: val,
			})
		}
	}
	if len(deviceSettings) > 0 {
		result.SimpleSettings = deviceSettings
	}

	// Channels
	var channels []cddmodels.ChannelConfiguration
	for _, regCh := range registration.Channels {
		channels = append(channels, s.buildChannelConfig(regCh))
	}
	result.Channels = channels
	return result
}

func (s *Tr12Shim) buildChannelConfig(regCh cddmodels.Channel) cddmodels.ChannelConfiguration {
	chID := regCh.Id
	chCfg := cddmodels.ChannelConfiguration{
		Id:    chID,
		State: s.CB.GetChannelState(chID),
	}

	// Check profiles first
	hasProfile := false
	if len(regCh.Profiles) > 0 {
		if profileID, found := s.CB.GetChannelProfileValue(chID); found {
			chCfg.Settings = &cddmodels.SettingsChoice{
				Profile: &cddmodels.Profile{
					Profile: cddmodels.SettingProfile{Id: profileID},
				},
			}
			hasProfile = true
		}
	}

	// Simple settings if no profile
	if !hasProfile && len(regCh.SimpleSettings) > 0 {
		var kvList []cddmodels.IdAndValue
		for _, setting := range regCh.SimpleSettings {
			// BUG WORKAROUND: use Setting.Name instead of Setting.Id
			if val, found := s.CB.GetChannelUpdatedValue(chID, setting.Name); found {
				kvList = append(kvList, cddmodels.IdAndValue{
					Key: setting.Name, Value: val,
				})
			}
		}
		if len(kvList) > 0 {
			chCfg.Settings = &cddmodels.SettingsChoice{
				SimpleSettings: &cddmodels.SimpleSettings{
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

// GetDeviceStatus builds the typed device status payload.
func (s *Tr12Shim) GetDeviceStatus() *cddmodels.DeviceStatus {
	channelID := "CH01"
	return &cddmodels.DeviceStatus{
		Status: s.CB.GetDeviceStatus(),
		Channels: []cddmodels.ChannelStatus{
			{
				Id:     channelID,
				State:  s.CB.GetChannelState(channelID),
				Status: s.CB.GetChannelStatus(channelID),
			},
		},
	}
}

// PrintActualConfig is a debug helper.
func (s *Tr12Shim) PrintActualConfig(registration *cddmodels.DeviceRegistration) {
	actual := s.GetActualConfiguration(registration)
	fmt.Printf("[SHIM TEST] get_actual_configuration: %+v\n", actual)
}
