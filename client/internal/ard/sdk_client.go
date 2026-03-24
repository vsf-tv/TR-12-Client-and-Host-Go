// Copyright 2025 Amazon.com Inc
// Licensed under the Apache License, Version 2.0
//
// SDK REST client — makes HTTP calls to the running CDD SDK daemon.
// Deserializes JSON responses into typed generated model structs.
package ard

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/pkg/cddmodels"

	localmodels "github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/models"
)

// SDKClient wraps HTTP calls to the CDD SDK REST API.
type SDKClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewSDKClient creates a new SDK REST client.
func NewSDKClient(baseURL string) *SDKClient {
	return &SDKClient{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// Connect calls PUT /connect with a typed ConnectRequestContent body.
func (c *SDKClient) Connect(hostID string, registration *cddmodels.DeviceRegistration) (*cddmodels.ConnectResponseContent, error) {
	body := cddmodels.ConnectRequestContent{
		HostId:       hostID,
		Registration: *registration,
	}
	return doTypedPUT[cddmodels.ConnectResponseContent](c, "/connect", body)
}

// Disconnect calls PUT /disconnect.
func (c *SDKClient) Disconnect() (*cddmodels.DisconnectResponseContent, error) {
	return doTypedPUT[cddmodels.DisconnectResponseContent](c, "/disconnect", nil)
}

// GetState calls GET /get_state.
func (c *SDKClient) GetState() (*cddmodels.GetConnectionStatusResponseContent, error) {
	return doTypedGET[cddmodels.GetConnectionStatusResponseContent](c, "/get_state")
}

// ReportStatus calls PUT /report_status with a typed DeviceStatus body.
func (c *SDKClient) ReportStatus(status *cddmodels.DeviceStatus) (*cddmodels.ReportStatusResponseContent, error) {
	body := map[string]interface{}{
		"status": status,
	}
	return doTypedPUT[cddmodels.ReportStatusResponseContent](c, "/report_status", body)
}

// ReportActualConfiguration calls PUT /report_actual_configuration with a typed DeviceConfiguration body.
func (c *SDKClient) ReportActualConfiguration(config *cddmodels.DeviceConfiguration) (*cddmodels.ReportActualConfigurationResponseContent, error) {
	body := map[string]interface{}{
		"configuration": config,
	}
	return doTypedPUT[cddmodels.ReportActualConfigurationResponseContent](c, "/report_actual_configuration", body)
}

// GetConfiguration calls GET /get_configuration.
// Uses the local model (not the generated one) because the SDK returns raw MQTT payload
// which may contain fields unknown to the strict generated model.
func (c *SDKClient) GetConfiguration() (*localmodels.GetConfigurationResponseContent, error) {
	return doTypedGET[localmodels.GetConfigurationResponseContent](c, "/get_configuration")
}

// doTypedPUT sends a PUT request and deserializes the response into T.
func doTypedPUT[T any](c *SDKClient, path string, body interface{}) (*T, error) {
	var reqBody io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal error: %w", err)
		}
		reqBody = bytes.NewReader(data)
	}
	req, err := http.NewRequest(http.MethodPut, c.baseURL+path, reqBody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return doRequestTyped[T](c, req)
}

// doTypedGET sends a GET request and deserializes the response into T.
func doTypedGET[T any](c *SDKClient, path string) (*T, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	return doRequestTyped[T](c, req)
}

// doRequestTyped executes the HTTP request and unmarshals the JSON response into T.
func doRequestTyped[T any](c *SDKClient, req *http.Request) (*T, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("SDK request failed: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	return &result, nil
}
