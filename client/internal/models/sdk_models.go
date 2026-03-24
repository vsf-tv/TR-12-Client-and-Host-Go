// Copyright 2025 Amazon.com Inc
// Licensed under the Apache License, Version 2.0
//
// Client-side SDK API response models — re-exports from the Smithy-generated cddmodels package.
// ConfigurationData and GetConfigurationResponseContent are kept local because the generated
// ConfigurationData.Payload is *DeviceConfiguration, but the SDK stores raw MQTT JSON as
// map[string]interface{}.
package models

import (
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/pkg/cddmodels"
)

// ---- Type aliases for generated CDD SDK models ----

type ErrorDetails = cddmodels.ErrorDetails
type ConnectResponseContent = cddmodels.ConnectResponseContent
type DisconnectResponseContent = cddmodels.DisconnectResponseContent
type GetConnectionStatusResponseContent = cddmodels.GetConnectionStatusResponseContent
type DeprovisionResponseContent = cddmodels.DeprovisionResponseContent
type ReportStatusResponseContent = cddmodels.ReportStatusResponseContent
type ReportActualConfigurationResponseContent = cddmodels.ReportActualConfigurationResponseContent

// ConfigurationData is kept local — generated ConfigurationData.Payload is *DeviceConfiguration,
// but SDK stores raw MQTT JSON as map[string]interface{}.
type ConfigurationData struct {
	Payload  map[string]interface{} `json:"payload,omitempty"`
	UpdateID string                 `json:"updateId"`
}

// GetUpdateId returns the UpdateID field (mirrors the generated model's getter).
func (c *ConfigurationData) GetUpdateId() string {
	if c == nil {
		return ""
	}
	return c.UpdateID
}

// GetConfigurationResponseContent is kept local to use our local ConfigurationData.
type GetConfigurationResponseContent struct {
	Success       bool               `json:"success"`
	State         string             `json:"state"`
	Message       string             `json:"message"`
	Configuration *ConfigurationData `json:"configuration,omitempty"`
	Error         *cddmodels.ErrorDetails `json:"error,omitempty"`
}

// States enumerates the SDK state machine values.
const (
	StateDisconnected = "DISCONNECTED"
	StatePairing      = "PAIRING"
	StateConnecting   = "CONNECTING"
	StateConnected    = "CONNECTED"
	StateReconnecting = "RECONNECTING"
)
