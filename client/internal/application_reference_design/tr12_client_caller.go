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
// SDK REST client — makes HTTP calls to the running CDD SDK daemon.
package application_reference_design

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	cddsdkgo "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"
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

// Connect calls PUT /connect.
func (c *SDKClient) Connect(hostID string, registration *cddsdkgo.DeviceRegistration) (*cddsdkgo.ConnectResponseContent, error) {
	body := cddsdkgo.ConnectRequestContent{
		HostId:       hostID,
		Registration: *registration,
	}
	return doTypedPUT[cddsdkgo.ConnectResponseContent](c, "/connect", body)
}

// Disconnect calls PUT /disconnect.
func (c *SDKClient) Disconnect() (*cddsdkgo.DisconnectResponseContent, error) {
	return doTypedPUT[cddsdkgo.DisconnectResponseContent](c, "/disconnect", nil)
}

// GetState calls GET /get_state.
func (c *SDKClient) GetState() (*cddsdkgo.GetConnectionStatusResponseContent, error) {
	return doTypedGET[cddsdkgo.GetConnectionStatusResponseContent](c, "/get_state")
}

// ReportStatus calls PUT /report_status.
func (c *SDKClient) ReportStatus(status *cddsdkgo.DeviceStatus) (*cddsdkgo.ReportStatusResponseContent, error) {
	body := map[string]interface{}{"status": status}
	return doTypedPUT[cddsdkgo.ReportStatusResponseContent](c, "/report_status", body)
}

// ReportActualConfiguration calls PUT /report_actual_configuration.
func (c *SDKClient) ReportActualConfiguration(config *cddsdkgo.DeviceConfiguration) (*cddsdkgo.ReportActualConfigurationResponseContent, error) {
	body := map[string]interface{}{"configuration": config}
	return doTypedPUT[cddsdkgo.ReportActualConfigurationResponseContent](c, "/report_actual_configuration", body)
}

// GetConfiguration calls GET /get_configuration.
func (c *SDKClient) GetConfiguration() (*cddsdkgo.GetConfigurationResponseContent, error) {
	return doTypedGET[cddsdkgo.GetConfigurationResponseContent](c, "/get_configuration")
}

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

func doTypedGET[T any](c *SDKClient, path string) (*T, error) {
	req, err := http.NewRequest(http.MethodGet, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}
	return doRequestTyped[T](c, req)
}

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
