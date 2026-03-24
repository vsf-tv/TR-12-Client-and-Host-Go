// Copyright 2025 Amazon.com Inc
// Licensed under the Apache License, Version 2.0
//
// TR-12 Callback implementations — mirrors tr12_callbacks.py.
// Provides get/set callbacks for device and channel settings.
package ard

import (
	"fmt"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/pkg/cddmodels"
)

// Callbacks implements device-specific get/set logic for the TR-12 shim.
type Callbacks struct {
	Encoder *Encoder
}

// NewCallbacks creates callbacks with a new encoder instance.
func NewCallbacks() *Callbacks {
	return &Callbacks{Encoder: NewEncoder()}
}

// --- Update (set) callbacks ---

func (cb *Callbacks) UpdateDeviceKeyValue(key, value string) {
	fmt.Printf("[UPDATE] Device setting: %s = %s\n", key, value)
}

func (cb *Callbacks) UpdateChannelSettings(channelID, key, value string) {
	fmt.Printf("[UPDATE] Channel %s setting: %s = %s\n", channelID, key, value)
}

func (cb *Callbacks) UpdateChannelProfile(channelID, profileID string) {
	fmt.Printf("[UPDATE] Channel %s profile: %s\n", channelID, profileID)
}

func (cb *Callbacks) UpdateChannelConnection(channelID string, connection *cddmodels.Connection) {
	fmt.Printf("[UPDATE] Channel %s connection: %+v\n", channelID, connection)
	cb.Encoder.HandleTransportConfigChange(channelID, connection)
}

func (cb *Callbacks) UpdateChannelState(channelID string, state cddmodels.ChannelState) {
	fmt.Printf("[UPDATE] Channel %s state: %s\n", channelID, state)
	cb.Encoder.HandleUpdateState(channelID, state)
}

// --- Get (read-back) callbacks ---

func (cb *Callbacks) GetDeviceUpdatedValue(key string) (string, bool) {
	defaults := map[string]string{
		"sync_clock_source": "NTP",
	}
	v, ok := defaults[key]
	return v, ok
}

func (cb *Callbacks) GetChannelUpdatedValue(channelID, key string) (string, bool) {
	// BUG WORKAROUND: Host service sends Setting.name instead of Setting.id
	defaults := map[string]string{
		"resolution":     "1920x1080",
		"framerate":      "30",
		"max_bitrate":    "10000",
		"rate_control":   "CBR",
		"codec":          "H.264",
		"gop_size":       "60",
		"selected_input": "SDI1",
	}
	v, ok := defaults[key]
	return v, ok
}

func (cb *Callbacks) GetChannelProfileValue(channelID string) (string, bool) {
	return "", false // Use simple settings by default
}

func (cb *Callbacks) GetChannelConnection(channelID string) *cddmodels.Connection {
	srtProto := cddmodels.SrtCallerTransportProtocol{
		Ip:                         "127.0.0.1",
		Port:                       5000,
		MinimumLatencyMilliseconds: 200,
	}
	streamID := "test_stream"
	srtProto.StreamId = &streamID

	tp := cddmodels.SrtCallerAsTransportProtocol(
		cddmodels.NewSrtCaller(srtProto),
	)
	conn := cddmodels.NewConnection()
	conn.SetTransportProtocol(tp)
	return conn
}

func (cb *Callbacks) GetChannelState(channelID string) cddmodels.ChannelState {
	return cb.Encoder.GetChannelState(channelID)
}

func (cb *Callbacks) GetDeviceStatus() []cddmodels.StatusValue {
	if cb.Encoder.Running() {
		return []cddmodels.StatusValue{
			{Name: "cpu", Value: "61", Info: "Current CPU % utilization."},
			{Name: "temp", Value: "84", Info: "CPU in degrees C."},
			{Name: "model", Value: "Talon", Info: "Hardware device model identifier."},
			{Name: "serial", Value: "123456789", Info: "Device serial number."},
		}
	}
	return []cddmodels.StatusValue{
		{Name: "cpu", Value: "31", Info: "Current CPU % utilization."},
		{Name: "temp", Value: "76", Info: "CPU in degrees C."},
		{Name: "model", Value: "Talon", Info: "Hardware device model identifier."},
		{Name: "serial", Value: "123456789", Info: "Device serial number."},
	}
}

func (cb *Callbacks) GetChannelStatus(channelID string) []cddmodels.StatusValue {
	if cb.Encoder.Running() {
		return []cddmodels.StatusValue{
			{Name: "bitrate", Value: GetSimulatedBitrate(), Info: "Bitrate Mbps configured on the video encoder."},
		}
	}
	return []cddmodels.StatusValue{
		{Name: "bitrate", Value: "0", Info: "Bitrate Mbps configured on the video encoder."},
	}
}
