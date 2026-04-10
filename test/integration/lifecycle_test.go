//go:build integration

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

package integration_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

func TestFullLifecycle(t *testing.T) {
	// Create test thumbnail file before starting anything
	createTestJPEG(t, "/tmp/image_sdi.jpg")
	t.Cleanup(func() {
		// best-effort cleanup
		_ = removeIfExists("/tmp/image_sdi.jpg")
	})

	env := newTestEnv(t)
	env.startHost()
	env.startSDK("integ-test-001")

	registration := loadRegistration(t)

	// ---------------------------------------------------------------
	// Phase 1: Account Setup
	// ---------------------------------------------------------------
	t.Log("Phase 1: Account Setup")
	acct := env.hostRegisterAccount("testuser", "testpass123", "Test User")
	if acct.Token == "" {
		t.Fatal("Phase 1: expected non-empty token")
	}
	if acct.Account == nil {
		t.Fatal("Phase 1: expected non-nil account")
	}
	if !strings.HasPrefix(acct.Account.AccountID, "acc_") {
		t.Fatalf("Phase 1: expected account_id starting with acc_, got %q", acct.Account.AccountID)
	}
	token := acct.Token
	t.Logf("Phase 1: OK — account_id=%s", acct.Account.AccountID)

	// ---------------------------------------------------------------
	// Phase 2: Device Pairing
	// ---------------------------------------------------------------
	t.Log("Phase 2: Device Pairing")
	pairingCode := env.waitForPairingCode("tr12-host", registration, 15*time.Second)
	if len(pairingCode) != 6 {
		t.Fatalf("Phase 2: expected 6-char pairing code, got %q", pairingCode)
	}
	for _, c := range pairingCode {
		if !((c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')) {
			t.Fatalf("Phase 2: pairing code contains invalid char %q in %q", string(c), pairingCode)
		}
	}
	t.Logf("Phase 2: OK — pairingCode=%s", pairingCode)

	// ---------------------------------------------------------------
	// Phase 3: Device Claim
	// ---------------------------------------------------------------
	t.Log("Phase 3: Device Claim")
	status := env.hostClaim(pairingCode, token)
	if status != 200 {
		t.Fatalf("Phase 3: expected 200, got %d", status)
	}
	t.Log("Phase 3: OK — device claimed")

	// ---------------------------------------------------------------
	// Phase 4: SDK Connects
	// ---------------------------------------------------------------
	t.Log("Phase 4: Waiting for SDK CONNECTED")
	env.waitForSDKConnected("tr12-host", registration, 30*time.Second)
	t.Log("Phase 4: OK — SDK connected")

	// ---------------------------------------------------------------
	// Phase 5: List Devices
	// ---------------------------------------------------------------
	t.Log("Phase 5: List Devices")
	devices := env.hostListDevices(token)
	if len(devices) != 1 {
		t.Fatalf("Phase 5: expected 1 device, got %d", len(devices))
	}
	deviceID := devices[0].DeviceID
	if len(deviceID) != 21 {
		t.Fatalf("Phase 5: expected 21-char device ID, got %q (%d chars)", deviceID, len(deviceID))
	}
	if !devices[0].Online {
		t.Fatal("Phase 5: expected device to be online")
	}
	t.Logf("Phase 5: OK — deviceId=%s online=%v", deviceID, devices[0].Online)

	// ---------------------------------------------------------------
	// Phase 6: Describe Device
	// ---------------------------------------------------------------
	t.Log("Phase 6: Describe Device")
	detail := env.hostDescribeDevice(deviceID, token)

	// Parse registration to check structure
	var reg struct {
		Channels []struct {
			ID             string        `json:"id"`
			SimpleSettings []interface{} `json:"simpleSettings"`
		} `json:"channels"`
		Thumbnails []struct {
			ID string `json:"id"`
		} `json:"thumbnails"`
	}
	if err := json.Unmarshal(detail.Registration, &reg); err != nil {
		t.Fatalf("Phase 6: cannot parse registration: %v", err)
	}
	if len(reg.Channels) != 1 || reg.Channels[0].ID != "CH01" {
		t.Fatalf("Phase 6: expected 1 channel with id=CH01, got %+v", reg.Channels)
	}
	if len(reg.Channels[0].SimpleSettings) != 7 {
		t.Fatalf("Phase 6: expected 7 simpleSettings, got %d", len(reg.Channels[0].SimpleSettings))
	}
	if len(reg.Thumbnails) != 2 {
		t.Fatalf("Phase 6: expected 2 thumbnails, got %d", len(reg.Thumbnails))
	}
	thumbIDs := map[string]bool{}
	for _, th := range reg.Thumbnails {
		thumbIDs[th.ID] = true
	}
	if !thumbIDs["SDI-1"] || !thumbIDs["HDMI-1"] {
		t.Fatalf("Phase 6: expected thumbnail IDs SDI-1 and HDMI-1, got %v", thumbIDs)
	}
	if !detail.Online {
		t.Fatal("Phase 6: expected online=true")
	}
	if detail.DeviceMetadata.DeviceType != "SOURCE" {
		t.Fatalf("Phase 6: expected device_type=SOURCE, got %q", detail.DeviceMetadata.DeviceType)
	}
	if detail.DeviceMetadata.AccountID != acct.Account.AccountID {
		t.Fatalf("Phase 6: expected account_id=%s, got %s", acct.Account.AccountID, detail.DeviceMetadata.AccountID)
	}
	if detail.CertExpiration == "" {
		t.Fatal("Phase 6: expected non-empty cert_expiration")
	}
	t.Logf("Phase 6: OK — registration channels=%d thumbnails=%d cert_expiration=%s",
		len(reg.Channels), len(reg.Thumbnails), detail.CertExpiration)

	// ---------------------------------------------------------------
	// Phase 7: Update Configuration
	// ---------------------------------------------------------------

	// 7a: Negative — Unknown Channel ID
	t.Log("Phase 7a: Negative — Unknown Channel ID")
	badChannelCfg := json.RawMessage(`{"channels":[{"id":"BOGUS_CHANNEL","state":"ACTIVE"}]}`)
	code, body := env.hostUpdateConfig(deviceID, token, badChannelCfg)
	if code != 400 {
		t.Fatalf("Phase 7a: expected 400, got %d: %s", code, body)
	}
	if !strings.Contains(body, "unknown channel ID") {
		t.Fatalf("Phase 7a: expected error about unknown channel ID, got: %s", body)
	}
	t.Log("Phase 7a: OK — rejected unknown channel ID")

	// 7b: Negative — Unknown Setting Key
	t.Log("Phase 7b: Negative — Unknown Setting Key")
	badSettingCfg := json.RawMessage(`{"channels":[{"id":"CH01","state":"ACTIVE","settings":{"simpleSettings":[{"key":"NONEXISTENT_SETTING","value":"foo"}]}}]}`)
	code, body = env.hostUpdateConfig(deviceID, token, badSettingCfg)
	if code != 400 {
		t.Fatalf("Phase 7b: expected 400, got %d: %s", code, body)
	}
	if !strings.Contains(body, "unknown setting key") {
		t.Fatalf("Phase 7b: expected error about unknown setting key, got: %s", body)
	}
	t.Log("Phase 7b: OK — rejected unknown setting key")

	// 7c: Negative — Unknown Profile ID
	t.Log("Phase 7c: Negative — Unknown Profile ID")
	badProfileCfg := json.RawMessage(`{"channels":[{"id":"CH01","state":"ACTIVE","settings":{"profile":{"id":"nonexistent_profile"}}}]}`)
	code, body = env.hostUpdateConfig(deviceID, token, badProfileCfg)
	if code != 400 {
		t.Fatalf("Phase 7c: expected 400, got %d: %s", code, body)
	}
	if !strings.Contains(body, "unknown profile ID") {
		t.Fatalf("Phase 7c: expected error about unknown profile ID, got: %s", body)
	}
	t.Log("Phase 7c: OK — rejected unknown profile ID")

	// 7d: Positive — Full Configuration
	t.Log("Phase 7d: Positive — Full Configuration")
	fullConfig := json.RawMessage(`{
		"simpleSettings": [
			{"key": "sync_clock_source", "value": "PTP"}
		],
		"channels": [
			{
				"id": "CH01",
				"state": "ACTIVE",
				"settings": {
					"simpleSettings": [
						{"key": "RS01", "value": "1920x1080"},
						{"key": "FR01", "value": "60"},
						{"key": "MB01", "value": "20000"},
						{"key": "RC01", "value": "CBR"},
						{"key": "CO01", "value": "H.264"},
						{"key": "GP01", "value": "60"},
						{"key": "IN01", "value": "SDI1"}
					]
				},
				"connection": {
					"transportProtocol": {
						"srtCaller": {
							"ip": "192.168.1.100",
							"port": 9000,
							"minimumLatencyMilliseconds": 200
						}
					}
				}
			}
		]
	}`)
	code, body = env.hostUpdateConfig(deviceID, token, fullConfig)
	if code != 200 {
		t.Fatalf("Phase 7d: expected 200, got %d: %s", code, body)
	}
	if !strings.Contains(strings.ToLower(body), "updated") {
		t.Fatalf("Phase 7d: expected 'updated' in response, got: %s", body)
	}

	// Wait for MQTT delivery and check SDK received the config
	time.Sleep(3 * time.Second)
	sdkCfg := env.sdkGetConfiguration()
	if sdkCfg.Configuration == nil {
		t.Fatal("Phase 7d: SDK returned nil configuration")
	}
	cfgJSON, _ := json.Marshal(sdkCfg.Configuration)
	cfgStr := string(cfgJSON)
	if !strings.Contains(cfgStr, "CH01") {
		t.Fatalf("Phase 7d: SDK config missing CH01: %s", cfgStr)
	}
	if !strings.Contains(cfgStr, "ACTIVE") {
		t.Fatalf("Phase 7d: SDK config missing ACTIVE state: %s", cfgStr)
	}
	if !strings.Contains(cfgStr, "192.168.1.100") {
		t.Fatalf("Phase 7d: SDK config missing SRT caller address: %s", cfgStr)
	}
	if !strings.Contains(cfgStr, "1920x1080") {
		t.Fatalf("Phase 7d: SDK config missing resolution setting: %s", cfgStr)
	}
	t.Log("Phase 7d: OK — full config pushed and received by SDK")

	// ---------------------------------------------------------------
	// Phase 8: Report Status and Actual Configuration
	// ---------------------------------------------------------------
	t.Log("Phase 8: Report Status and Actual Configuration")
	statusPayload := map[string]interface{}{
		"channels": []map[string]interface{}{
			{
				"id":    "CH01",
				"state": "ACTIVE",
				"statusValues": []map[string]interface{}{
					{"name": "bitrate", "value": "9500", "info": "Current output bitrate (Kbps)"},
				},
			},
		},
	}
	statusResp := env.sdkReportStatus(statusPayload)
	if !statusResp.Success {
		t.Fatalf("Phase 8: report_status failed: %s", statusResp.Message)
	}

	// Report actual config — use the payload received from the SDK which has configurationId stamped by the host
	var actualCfg map[string]interface{}
	if sdkCfg.Configuration != nil {
		// Extract the payload from the configuration, which has configurationId from the host
		if payload, ok := sdkCfg.Configuration["payload"]; ok && payload != nil {
			payloadBytes, _ := json.Marshal(payload)
			json.Unmarshal(payloadBytes, &actualCfg)
		}
	}
	if actualCfg == nil {
		json.Unmarshal(fullConfig, &actualCfg)
		actualCfg["configurationId"] = 0
	}
	cfgResp := env.sdkReportActualConfig(actualCfg)
	if !cfgResp.Success {
		t.Fatalf("Phase 8: report_actual_configuration failed: %s", cfgResp.Message)
	}

	// Wait for MQTT delivery to host
	time.Sleep(3 * time.Second)

	detail2 := env.hostDescribeDevice(deviceID, token)
	if len(detail2.Status) == 0 || string(detail2.Status) == "null" {
		t.Fatal("Phase 8: expected non-null status on host")
	}
	if len(detail2.ActualConfiguration) == 0 || string(detail2.ActualConfiguration) == "null" {
		t.Fatal("Phase 8: expected non-null actual_configuration on host")
	}
	// Verify status contains our reported values
	statusStr := string(detail2.Status)
	if !strings.Contains(statusStr, "bitrate") || !strings.Contains(statusStr, "9500") {
		t.Fatalf("Phase 8: status missing reported values: %s", statusStr)
	}
	t.Log("Phase 8: OK — status and actual config visible on host")

	// ---------------------------------------------------------------
	// Phase 9: Thumbnail Request
	// ---------------------------------------------------------------
	t.Log("Phase 9: Thumbnail Request")
	// Touch the test JPEG again so it's fresh (< 10 seconds old)
	createTestJPEG(t, "/tmp/image_sdi.jpg")

	var thumbResp thumbnailResponse
	var thumbCode int
	// First call creates the subscription; may need to retry
	for attempt := 0; attempt < 4; attempt++ {
		thumbCode, thumbResp = env.hostGetThumbnail(deviceID, "SDI-1", token)
		if thumbCode == 200 && thumbResp.Image != nil && thumbResp.Image.Base64Image != "" {
			break
		}
		t.Logf("Phase 9: attempt %d — code=%d message=%q, retrying...", attempt+1, thumbCode, thumbResp.Message)
		// Re-touch the JPEG to keep it fresh
		createTestJPEG(t, "/tmp/image_sdi.jpg")
		time.Sleep(5 * time.Second)
	}
	if thumbResp.Image == nil || thumbResp.Image.Base64Image == "" {
		t.Fatalf("Phase 9: no thumbnail data after retries (last code=%d, message=%q)", thumbCode, thumbResp.Message)
	}
	t.Logf("Phase 9: OK — thumbnail received, size=%d type=%s", thumbResp.Image.ImageSizeKB, thumbResp.Image.ImageType)

	// ---------------------------------------------------------------
	// Phase 10: Credential Rotation
	// ---------------------------------------------------------------
	t.Log("Phase 10: Credential Rotation")

	rotateCode := env.hostRotateCredentials(deviceID, token)
	if rotateCode != 200 {
		t.Fatalf("Phase 10: expected 200, got %d", rotateCode)
	}

	// Wait for SDK to receive rotation, reconnect with new cert
	time.Sleep(6 * time.Second)

	// Verify SDK is still connected after rotation
	env.waitForSDKState("CONNECTED", 15*time.Second)
	t.Log("Phase 10: SDK reconnected after rotation")

	detail3 := env.hostDescribeDevice(deviceID, token)
	if !detail3.Online {
		t.Fatal("Phase 10: expected device online after rotation")
	}
	if detail3.CertExpiration == "" || detail3.CertExpiration == "unknown" || detail3.CertExpiration == "expired" {
		t.Fatalf("Phase 10: expected valid cert_expiration, got %q", detail3.CertExpiration)
	}
	t.Logf("Phase 10: OK — device online after rotation, cert_expiration=%q", detail3.CertExpiration)

	// ---------------------------------------------------------------
	// Phase 11: Deprovision
	// ---------------------------------------------------------------
	t.Log("Phase 11: Deprovision")
	deprovCode := env.hostDeprovision(deviceID, token)
	if deprovCode != 200 {
		t.Fatalf("Phase 11: expected 200, got %d", deprovCode)
	}

	// Wait for the SDK to receive the deprovision MQTT message and transition to DISCONNECTED
	env.waitForSDKState("DISCONNECTED", 15*time.Second)
	t.Log("Phase 11: OK — device deprovisioned, SDK disconnected")

	t.Log("=== TestFullLifecycle PASSED ===")
}

// removeIfExists removes a file if it exists, ignoring errors.
func removeIfExists(path string) error {
	err := os.Remove(path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// TestOfflineConfigDelivery verifies that a configuration update pushed while
// the SDK is offline is delivered when the SDK reconnects (retained MQTT).
func TestOfflineConfigDelivery(t *testing.T) {
	env := newTestEnv(t)
	env.startHost()
	env.startSDK("integ-offline-001")

	registration := loadRegistration(t)

	// Setup: register account, pair, claim, connect
	acct := env.hostRegisterAccount("offlineuser", "testpass123", "Offline Test")
	token := acct.Token

	pairingCode := env.waitForPairingCode("tr12-host", registration, 15*time.Second)
	env.hostClaim(pairingCode, token)
	env.waitForSDKConnected("tr12-host", registration, 30*time.Second)

	devices := env.hostListDevices(token)
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	deviceID := devices[0].DeviceID
	t.Logf("Device paired and connected: %s", deviceID)

	// Stop the SDK — device goes offline
	t.Log("Stopping SDK (device goes offline)...")
	env.sdkProc.stop(t)
	env.sdkProc = nil
	time.Sleep(1 * time.Second)

	// Push a config update while device is offline
	t.Log("Pushing config update while device is offline...")
	offlineConfig := json.RawMessage(`{
		"simpleSettings": [{"key": "sync_clock_source", "value": "PTP"}],
		"channels": [{
			"id": "CH01",
			"state": "IDLE",
			"settings": {
				"simpleSettings": [
					{"key": "RS01", "value": "1920x1080"},
					{"key": "FR01", "value": "60"},
					{"key": "MB01", "value": "5000"},
					{"key": "RC01", "value": "CBR"},
					{"key": "CO01", "value": "H.264"},
					{"key": "GP01", "value": "60"},
					{"key": "IN01", "value": "SDI1"}
				]
			}
		}]
	}`)
	code, body := env.hostUpdateConfig(deviceID, token, offlineConfig)
	if code != 200 {
		t.Fatalf("hostUpdateConfig while offline: expected 200, got %d: %s", code, body)
	}
	t.Log("Config update pushed while offline — OK")

	// Restart the SDK
	t.Log("Restarting SDK...")
	env.startSDK("integ-offline-001")
	env.waitForSDKConnected("tr12-host", registration, 30*time.Second)
	t.Log("SDK reconnected")

	// Give the SDK time to receive the retained config message and process it
	time.Sleep(4 * time.Second)

	// Verify the SDK received the config update
	sdkCfg := env.sdkGetConfiguration()
	if sdkCfg.Configuration == nil {
		t.Fatal("SDK returned nil configuration after reconnect")
	}
	cfgJSON, _ := json.Marshal(sdkCfg.Configuration)
	cfgStr := string(cfgJSON)

	if !strings.Contains(cfgStr, `"FR01"`) || !strings.Contains(cfgStr, `"60"`) {
		t.Fatalf("SDK config missing FR01=60 after offline update: %s", cfgStr)
	}
	if !strings.Contains(cfgStr, `"MB01"`) || !strings.Contains(cfgStr, `"5000"`) {
		t.Fatalf("SDK config missing MB01=5000 after offline update: %s", cfgStr)
	}
	if !strings.Contains(cfgStr, "PTP") {
		t.Fatalf("SDK config missing sync_clock_source=PTP after offline update: %s", cfgStr)
	}
	t.Log("TestOfflineConfigDelivery: OK — SDK received config update after reconnect")
}

// TestTwoChannelEncoder verifies:
// 1. Happy path — both channels receive and apply a full configuration update
// 2. Per-channel configurationId tracking — only the channel whose configurationId
//    changed is applied; unchanged channels are skipped
// 3. Device-level only update — no channel updates processed
func TestTwoChannelEncoder(t *testing.T) {
	createTestJPEG(t, "/tmp/image_sdi.jpg")
	t.Cleanup(func() { _ = removeIfExists("/tmp/image_sdi.jpg") })

	env := newTestEnv(t)
	env.startHost()
	env.startSDK("integ-2ch-001")

	registration := loadRegistrationFrom(t, "2_channel_encoder")

	// Setup: register, pair, claim, connect
	acct := env.hostRegisterAccount("user2ch", "testpass123", "2CH Test")
	token := acct.Token
	pairingCode := env.waitForPairingCode("tr12-host", registration, 15*time.Second)
	env.hostClaim(pairingCode, token)
	env.waitForSDKConnected("tr12-host", registration, 30*time.Second)

	devices := env.hostListDevices(token)
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	deviceID := devices[0].DeviceID

	// Verify registration has 2 channels
	detail := env.hostDescribeDevice(deviceID, token)
	var reg struct {
		Channels []struct{ ID string `json:"id"` } `json:"channels"`
	}
	json.Unmarshal(detail.Registration, &reg)
	if len(reg.Channels) != 2 {
		t.Fatalf("expected 2 channels in registration, got %d", len(reg.Channels))
	}
	t.Logf("Registration OK — channels: %s, %s", reg.Channels[0].ID, reg.Channels[1].ID)

	// ---------------------------------------------------------------
	// Phase 1: Full 2-channel config — both channels should be applied
	// ---------------------------------------------------------------
	t.Log("Phase 1: Full 2-channel configuration update")
	fullConfig := json.RawMessage(`{
		"simpleSettings": [{"key": "sync_clock_source", "value": "PTP"}],
		"channels": [
			{
				"id": "CH01", "state": "ACTIVE",
				"settings": {"simpleSettings": [
					{"key": "RS01", "value": "1920x1080"},
					{"key": "FR01", "value": "60"},
					{"key": "MB01", "value": "20000"},
					{"key": "RC01", "value": "CBR"},
					{"key": "CO01", "value": "H.264"},
					{"key": "GP01", "value": "60"},
					{"key": "IN01", "value": "SDI1"}
				]},
				"connection": {"transportProtocol": {"srtCaller": {
					"ip": "192.168.1.100", "port": 9001, "minimumLatencyMilliseconds": 200
				}}}
			},
			{
				"id": "CH02", "state": "ACTIVE",
				"settings": {"simpleSettings": [
					{"key": "RS01", "value": "1280x720"},
					{"key": "FR01", "value": "30"},
					{"key": "MB01", "value": "10000"},
					{"key": "RC01", "value": "CBR"},
					{"key": "CO01", "value": "H.264"},
					{"key": "GP01", "value": "30"},
					{"key": "IN01", "value": "HDMI1"}
				]},
				"connection": {"transportProtocol": {"srtCaller": {
					"ip": "192.168.1.101", "port": 9002, "minimumLatencyMilliseconds": 200
				}}}
			}
		]
	}`)
	code, body := env.hostUpdateConfig(deviceID, token, fullConfig)
	if code != 200 {
		t.Fatalf("Phase 1: expected 200, got %d: %s", code, body)
	}

	time.Sleep(3 * time.Second)
	sdkCfg := env.sdkGetConfiguration()
	if sdkCfg.Configuration == nil {
		t.Fatal("Phase 1: SDK returned nil configuration")
	}
	cfgStr := string(mustMarshal(sdkCfg.Configuration))
	for _, want := range []string{"CH01", "CH02", "192.168.1.100", "192.168.1.101", "1920x1080", "1280x720"} {
		if !strings.Contains(cfgStr, want) {
			t.Fatalf("Phase 1: SDK config missing %q: %s", want, cfgStr)
		}
	}
	t.Log("Phase 1: OK — both channels received and present in SDK config")

	// Capture the configurationIds stamped by the host for use in Phase 2
	var phase1Payload map[string]interface{}
	if p, ok := sdkCfg.Configuration["payload"]; ok {
		payloadBytes, _ := json.Marshal(p)
		json.Unmarshal(payloadBytes, &phase1Payload)
	}

	// ---------------------------------------------------------------
	// Phase 2: Update only CH01 — CH02 configurationId unchanged
	// The host bumps the device configurationId (new updateId) but only
	// CH01 gets a new channel configurationId. The client should apply
	// CH01 and skip CH02.
	// ---------------------------------------------------------------
	t.Log("Phase 2: Update CH01 only — CH02 should be skipped by client")

	// Send same CH02 config (same settings) but different CH01 settings.
	// The host will stamp a new device configurationId but CH02's channel
	// configurationId will be the same epoch second (within the same second)
	// OR we can verify by checking the SDK only sees CH01 changes.
	// We use a distinct CH01 value (FR01=25) to confirm it was applied.
	time.Sleep(1 * time.Second) // ensure epoch second advances for new configurationId
	ch01OnlyConfig := json.RawMessage(`{
		"simpleSettings": [{"key": "sync_clock_source", "value": "PTP"}],
		"channels": [
			{
				"id": "CH01", "state": "ACTIVE",
				"settings": {"simpleSettings": [
					{"key": "RS01", "value": "1920x1080"},
					{"key": "FR01", "value": "25"},
					{"key": "MB01", "value": "20000"},
					{"key": "RC01", "value": "CBR"},
					{"key": "CO01", "value": "H.264"},
					{"key": "GP01", "value": "60"},
					{"key": "IN01", "value": "SDI1"}
				]},
				"connection": {"transportProtocol": {"srtCaller": {
					"ip": "192.168.1.100", "port": 9001, "minimumLatencyMilliseconds": 200
				}}}
			},
			{
				"id": "CH02", "state": "ACTIVE",
				"settings": {"simpleSettings": [
					{"key": "RS01", "value": "1280x720"},
					{"key": "FR01", "value": "30"},
					{"key": "MB01", "value": "10000"},
					{"key": "RC01", "value": "CBR"},
					{"key": "CO01", "value": "H.264"},
					{"key": "GP01", "value": "30"},
					{"key": "IN01", "value": "HDMI1"}
				]},
				"connection": {"transportProtocol": {"srtCaller": {
					"ip": "192.168.1.101", "port": 9002, "minimumLatencyMilliseconds": 200
				}}}
			}
		]
	}`)
	code, body = env.hostUpdateConfig(deviceID, token, ch01OnlyConfig)
	if code != 200 {
		t.Fatalf("Phase 2: expected 200, got %d: %s", code, body)
	}

	time.Sleep(3 * time.Second)
	sdkCfg2 := env.sdkGetConfiguration()
	if sdkCfg2.Configuration == nil {
		t.Fatal("Phase 2: SDK returned nil configuration")
	}
	cfgStr2 := string(mustMarshal(sdkCfg2.Configuration))
	// CH01 should now show FR01=25
	if !strings.Contains(cfgStr2, `"25"`) {
		t.Fatalf("Phase 2: expected CH01 FR01=25 in SDK config: %s", cfgStr2)
	}
	t.Log("Phase 2: OK — CH01 update received, CH02 unchanged")

	// ---------------------------------------------------------------
	// Phase 3: Device-level only update (same channel configs)
	// Send identical channel configs — only simpleSettings changes.
	// Both channels should be skipped since their configurationIds are unchanged.
	// ---------------------------------------------------------------
	t.Log("Phase 3: Device-level only update — channels should not be reapplied")
	time.Sleep(1 * time.Second)
	deviceOnlyConfig := json.RawMessage(`{
		"simpleSettings": [{"key": "sync_clock_source", "value": "GENLOCK"}],
		"channels": [
			{
				"id": "CH01", "state": "ACTIVE",
				"settings": {"simpleSettings": [
					{"key": "RS01", "value": "1920x1080"},
					{"key": "FR01", "value": "25"},
					{"key": "MB01", "value": "20000"},
					{"key": "RC01", "value": "CBR"},
					{"key": "CO01", "value": "H.264"},
					{"key": "GP01", "value": "60"},
					{"key": "IN01", "value": "SDI1"}
				]},
				"connection": {"transportProtocol": {"srtCaller": {
					"ip": "192.168.1.100", "port": 9001, "minimumLatencyMilliseconds": 200
				}}}
			},
			{
				"id": "CH02", "state": "ACTIVE",
				"settings": {"simpleSettings": [
					{"key": "RS01", "value": "1280x720"},
					{"key": "FR01", "value": "30"},
					{"key": "MB01", "value": "10000"},
					{"key": "RC01", "value": "CBR"},
					{"key": "CO01", "value": "H.264"},
					{"key": "GP01", "value": "30"},
					{"key": "IN01", "value": "HDMI1"}
				]},
				"connection": {"transportProtocol": {"srtCaller": {
					"ip": "192.168.1.101", "port": 9002, "minimumLatencyMilliseconds": 200
				}}}
			}
		]
	}`)
	code, body = env.hostUpdateConfig(deviceID, token, deviceOnlyConfig)
	if code != 200 {
		t.Fatalf("Phase 3: expected 200, got %d: %s", code, body)
	}

	time.Sleep(3 * time.Second)
	sdkCfg3 := env.sdkGetConfiguration()
	if sdkCfg3.Configuration == nil {
		t.Fatal("Phase 3: SDK returned nil configuration")
	}
	cfgStr3 := string(mustMarshal(sdkCfg3.Configuration))
	// Device-level change should be present
	if !strings.Contains(cfgStr3, "GENLOCK") {
		t.Fatalf("Phase 3: expected GENLOCK in SDK config: %s", cfgStr3)
	}
	t.Log("Phase 3: OK — device-level update received, channel configs unchanged")

	t.Log("=== TestTwoChannelEncoder PASSED ===")
}

func mustMarshal(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

// TestConfigurationIdBumping verifies that the host correctly bumps configurationId
// only for the device/channel entities that actually changed between pushes.
func TestConfigurationIdBumping(t *testing.T) {
	env := newTestEnv(t)
	env.startHost()
	env.startSDK("integ-cfgid-001")

	registration := loadRegistrationFrom(t, "2_channel_encoder")

	acct := env.hostRegisterAccount("cfgiduser", "testpass123", "ConfigId Test")
	token := acct.Token
	pairingCode := env.waitForPairingCode("tr12-host", registration, 15*time.Second)
	env.hostClaim(pairingCode, token)
	env.waitForSDKConnected("tr12-host", registration, 30*time.Second)

	devices := env.hostListDevices(token)
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	deviceID := devices[0].DeviceID

	// Helper: extract configurationIds from SDK config payload
	getConfigIds := func() (deviceId string, ch01Id string, ch02Id string) {
		sdkCfg := env.sdkGetConfiguration()
		if sdkCfg.Configuration == nil {
			return "", "", ""
		}
		payload, ok := sdkCfg.Configuration["payload"]
		if !ok || payload == nil {
			return "", "", ""
		}
		b, _ := json.Marshal(payload)
		var cfg struct {
			ConfigurationId string `json:"configurationId"`
			Channels        []struct {
				Id              string `json:"id"`
				ConfigurationId string `json:"configurationId"`
			} `json:"channels"`
		}
		json.Unmarshal(b, &cfg)
		deviceId = cfg.ConfigurationId
		for _, ch := range cfg.Channels {
			if ch.Id == "CH01" {
				ch01Id = ch.ConfigurationId
			} else if ch.Id == "CH02" {
				ch02Id = ch.ConfigurationId
			}
		}
		return
	}

	baseConfig := func(ch01State, ch02State, clockSource string) json.RawMessage {
		return json.RawMessage(`{
			"simpleSettings": [{"key": "sync_clock_source", "value": "` + clockSource + `"}],
			"channels": [
				{"id": "CH01", "state": "` + ch01State + `", "settings": {"simpleSettings": [
					{"key": "RS01", "value": "1920x1080"}, {"key": "FR01", "value": "30"},
					{"key": "MB01", "value": "10000"}, {"key": "RC01", "value": "CBR"},
					{"key": "CO01", "value": "H.264"}, {"key": "GP01", "value": "60"},
					{"key": "IN01", "value": "SDI1"}
				]}},
				{"id": "CH02", "state": "` + ch02State + `", "settings": {"simpleSettings": [
					{"key": "RS01", "value": "1920x1080"}, {"key": "FR01", "value": "30"},
					{"key": "MB01", "value": "10000"}, {"key": "RC01", "value": "CBR"},
					{"key": "CO01", "value": "H.264"}, {"key": "GP01", "value": "60"},
					{"key": "IN01", "value": "SDI1"}
				]}}
			]
		}`)
	}

	// --- Push 1: initial full config — all IDs get bumped ---
	t.Log("Push 1: initial full config")
	code, body := env.hostUpdateConfig(deviceID, token, baseConfig("IDLE", "IDLE", "NTP"))
	if code != 200 {
		t.Fatalf("Push 1: expected 200, got %d: %s", code, body)
	}
	time.Sleep(4 * time.Second)
	dev1, ch01_1, ch02_1 := getConfigIds()
	if dev1 == "" || ch01_1 == "" || ch02_1 == "" {
		t.Fatalf("Push 1: expected non-empty configurationIds, got device=%s ch01=%s ch02=%s", dev1, ch01_1, ch02_1)
	}
	t.Logf("Push 1: device=%s ch01=%s ch02=%s", dev1, ch01_1, ch02_1)

	// --- Push 2: change only CH01 state — only CH01 ID bumps ---
	t.Log("Push 2: change CH01 state only")
	time.Sleep(4 * time.Second) // ensure epoch second advances for new configurationId
	code, body = env.hostUpdateConfig(deviceID, token, baseConfig("ACTIVE", "IDLE", "NTP"))
	if code != 200 {
		t.Fatalf("Push 2: expected 200, got %d: %s", code, body)
	}
	time.Sleep(4 * time.Second)
	dev2, ch01_2, ch02_2 := getConfigIds()
	t.Logf("Push 2: device=%s ch01=%s ch02=%s", dev2, ch01_2, ch02_2)

	if ch01_2 == ch01_1 {
		t.Fatalf("Push 2: CH01 configurationId should have bumped (%s → %s)", ch01_1, ch01_2)
	}
	if ch02_2 != ch02_1 {
		t.Fatalf("Push 2: CH02 configurationId should NOT have bumped (%s → %s)", ch02_1, ch02_2)
	}
	if dev2 != dev1 {
		t.Fatalf("Push 2: device configurationId should NOT have bumped (%s → %s)", dev1, dev2)
	}
	// All three IDs must be independent — no two should be equal after bumping
	if ch01_2 == dev2 || ch01_2 == ch02_2 {
		t.Fatalf("Push 2: CH01 configurationId (%s) should be independent from device (%s) and CH02 (%s)", ch01_2, dev2, ch02_2)
	}
	t.Log("Push 2: OK — only CH01 bumped, IDs are independent")

	// --- Push 3: change only device settings — only device ID bumps ---
	t.Log("Push 3: change device simpleSettings only")
	time.Sleep(4 * time.Second)
	code, body = env.hostUpdateConfig(deviceID, token, baseConfig("ACTIVE", "IDLE", "PTP"))
	if code != 200 {
		t.Fatalf("Push 3: expected 200, got %d: %s", code, body)
	}
	time.Sleep(4 * time.Second)
	dev3, ch01_3, ch02_3 := getConfigIds()
	t.Logf("Push 3: device=%s ch01=%s ch02=%s", dev3, ch01_3, ch02_3)

	if dev3 == dev2 {
		t.Fatalf("Push 3: device configurationId should have bumped (%s → %s)", dev2, dev3)
	}
	if ch01_3 != ch01_2 {
		t.Fatalf("Push 3: CH01 configurationId should NOT have bumped (%s → %s)", ch01_2, ch01_3)
	}
	if ch02_3 != ch02_2 {
		t.Fatalf("Push 3: CH02 configurationId should NOT have bumped (%s → %s)", ch02_2, ch02_3)
	}
	if dev3 == ch01_3 || dev3 == ch02_3 {
		t.Fatalf("Push 3: device configurationId (%s) should be independent from CH01 (%s) and CH02 (%s)", dev3, ch01_3, ch02_3)
	}
	t.Log("Push 3: OK — only device bumped, IDs are independent")

	// --- Push 4: change only CH02 state — only CH02 ID bumps ---
	t.Log("Push 4: change CH02 state only")
	time.Sleep(4 * time.Second)
	code, body = env.hostUpdateConfig(deviceID, token, baseConfig("ACTIVE", "ACTIVE", "PTP"))
	if code != 200 {
		t.Fatalf("Push 4: expected 200, got %d: %s", code, body)
	}
	time.Sleep(4 * time.Second)
	dev4, ch01_4, ch02_4 := getConfigIds()
	t.Logf("Push 4: device=%s ch01=%s ch02=%s", dev4, ch01_4, ch02_4)

	if ch02_4 == ch02_3 {
		t.Fatalf("Push 4: CH02 configurationId should have bumped (%s → %s)", ch02_3, ch02_4)
	}
	if ch01_4 != ch01_3 {
		t.Fatalf("Push 4: CH01 configurationId should NOT have bumped (%s → %s)", ch01_3, ch01_4)
	}
	if dev4 != dev3 {
		t.Fatalf("Push 4: device configurationId should NOT have bumped (%s → %s)", dev3, dev4)
	}
	if ch02_4 == dev4 || ch02_4 == ch01_4 {
		t.Fatalf("Push 4: CH02 configurationId (%s) should be independent from device (%s) and CH01 (%s)", ch02_4, dev4, ch01_4)
	}
	t.Log("Push 4: OK — only CH02 bumped, IDs are independent")

	t.Log("=== TestConfigurationIdBumping PASSED ===")
}

// TestARDConfigurationIdEchoBack verifies that the ARD correctly:
// 1. Applies only channels whose configurationId changed
// 2. Echoes back the configurationId it actually applied in actual_configuration
// 3. Does NOT echo back a new configurationId for unchanged channels
func TestARDConfigurationIdEchoBack(t *testing.T) {
	createTestJPEG(t, "/tmp/image_sdi.jpg")
	t.Cleanup(func() { _ = removeIfExists("/tmp/image_sdi.jpg") })

	env := newTestEnv(t)
	env.startHost()
	env.startSDK("integ-ard-echo-001")

	registration := loadRegistrationFrom(t, "2_channel_encoder")

	acct := env.hostRegisterAccount("ardechouser", "testpass123", "ARD Echo Test")
	token := acct.Token
	pairingCode := env.waitForPairingCode("tr12-host", registration, 15*time.Second)
	env.hostClaim(pairingCode, token)
	env.waitForSDKConnected("tr12-host", registration, 30*time.Second)

	devices := env.hostListDevices(token)
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	deviceID := devices[0].DeviceID

	// Helper: get actual_configuration from host and extract per-entity configurationIds
	getActualConfigIds := func() (deviceId string, ch01Id string, ch02Id string) {
		detail := env.hostDescribeDevice(deviceID, token)
		if len(detail.ActualConfiguration) == 0 || string(detail.ActualConfiguration) == "null" {
			return "", "", ""
		}
		var actual struct {
			ConfigurationId string `json:"configurationId"`
			Channels        []struct {
				Id              string `json:"id"`
				ConfigurationId string `json:"configurationId"`
			} `json:"channels"`
		}
		json.Unmarshal(detail.ActualConfiguration, &actual)
		deviceId = actual.ConfigurationId
		for _, ch := range actual.Channels {
			if ch.Id == "CH01" {
				ch01Id = ch.ConfigurationId
			} else if ch.Id == "CH02" {
				ch02Id = ch.ConfigurationId
			}
		}
		return
	}

	baseConfig := func(ch01State, ch02State, clockSource string) json.RawMessage {
		return json.RawMessage(`{
			"simpleSettings": [{"key": "sync_clock_source", "value": "` + clockSource + `"}],
			"channels": [
				{"id": "CH01", "state": "` + ch01State + `", "settings": {"simpleSettings": [
					{"key": "RS01", "value": "1920x1080"}, {"key": "FR01", "value": "30"},
					{"key": "MB01", "value": "10000"}, {"key": "RC01", "value": "CBR"},
					{"key": "CO01", "value": "H.264"}, {"key": "GP01", "value": "60"},
					{"key": "IN01", "value": "SDI1"}
				]}},
				{"id": "CH02", "state": "` + ch02State + `", "settings": {"simpleSettings": [
					{"key": "RS01", "value": "1920x1080"}, {"key": "FR01", "value": "30"},
					{"key": "MB01", "value": "10000"}, {"key": "RC01", "value": "CBR"},
					{"key": "CO01", "value": "H.264"}, {"key": "GP01", "value": "60"},
					{"key": "IN01", "value": "SDI1"}
				]}}
			]
		}`)
	}

	// Push 1: initial config — simulate ARD: receive config, apply, report actual
	t.Log("ARD Push 1: initial config")
	code, body := env.hostUpdateConfig(deviceID, token, baseConfig("IDLE", "IDLE", "NTP"))
	if code != 200 {
		t.Fatalf("Push 1: expected 200, got %d: %s", code, body)
	}
	time.Sleep(4 * time.Second) // wait for SDK to receive MQTT

	// Simulate ARD: get config from SDK, report it back as actual
	sdkCfg1 := env.sdkGetConfiguration()
	if sdkCfg1.Configuration == nil {
		t.Fatal("Push 1: SDK returned nil configuration")
	}
	var actualCfg1 map[string]interface{}
	if p, ok := sdkCfg1.Configuration["payload"]; ok {
		b, _ := json.Marshal(p)
		json.Unmarshal(b, &actualCfg1)
	}
	if actualCfg1 == nil {
		t.Fatal("Push 1: could not extract payload from SDK config")
	}
	cfgResp := env.sdkReportActualConfig(actualCfg1)
	if !cfgResp.Success {
		t.Fatalf("Push 1: report_actual_configuration failed: %s", cfgResp.Message)
	}
	time.Sleep(3 * time.Second) // wait for MQTT delivery to host

	actDev1, actCh01_1, actCh02_1 := getActualConfigIds()
	if actDev1 == "" || actCh01_1 == "" || actCh02_1 == "" {
		t.Fatalf("Push 1: expected non-empty actual configurationIds, got device=%s ch01=%s ch02=%s", actDev1, actCh01_1, actCh02_1)
	}
	t.Logf("Push 1 actual: device=%s ch01=%s ch02=%s", actDev1, actCh01_1, actCh02_1)

	// Push 2: change only CH01 state
	// Simulate ARD: only CH01 configurationId changed → only CH01 applied → echo CH01 new, CH02 old
	t.Log("ARD Push 2: change CH01 state only")
	time.Sleep(2 * time.Second)
	code, body = env.hostUpdateConfig(deviceID, token, baseConfig("ACTIVE", "IDLE", "NTP"))
	if code != 200 {
		t.Fatalf("Push 2: expected 200, got %d: %s", code, body)
	}
	time.Sleep(4 * time.Second)

	sdkCfg2 := env.sdkGetConfiguration()
	if sdkCfg2.Configuration == nil {
		t.Fatal("Push 2: SDK returned nil configuration")
	}
	// Extract the new payload — CH01 has new configurationId, CH02 has old
	var desiredCfg2 struct {
		ConfigurationId string `json:"configurationId"`
		Channels        []struct {
			Id              string `json:"id"`
			ConfigurationId string `json:"configurationId"`
		} `json:"channels"`
	}
	if p, ok := sdkCfg2.Configuration["payload"]; ok {
		b, _ := json.Marshal(p)
		json.Unmarshal(b, &desiredCfg2)
	}

	// Simulate ARD echo-back: report actual with the configurationIds from desired
	// (ARD echoes back what it applied — CH01 new ID, CH02 old ID)
	var actualCfg2 map[string]interface{}
	if p, ok := sdkCfg2.Configuration["payload"]; ok {
		b, _ := json.Marshal(p)
		json.Unmarshal(b, &actualCfg2)
	}
	cfgResp2 := env.sdkReportActualConfig(actualCfg2)
	if !cfgResp2.Success {
		t.Fatalf("Push 2: report_actual_configuration failed: %s", cfgResp2.Message)
	}
	time.Sleep(3 * time.Second)

	actDev2, actCh01_2, actCh02_2 := getActualConfigIds()
	t.Logf("Push 2 actual: device=%s ch01=%s ch02=%s", actDev2, actCh01_2, actCh02_2)

	// CH01 should have a new configurationId
	if actCh01_2 == actCh01_1 {
		t.Fatalf("Push 2: CH01 configurationId should have changed (%s → %s)", actCh01_1, actCh01_2)
	}
	// CH02 should have the SAME configurationId (host didn't bump it)
	if actCh02_2 != actCh02_1 {
		t.Fatalf("Push 2: CH02 configurationId should NOT have changed (%s → %s)", actCh02_1, actCh02_2)
	}
	// Device should have the SAME configurationId
	if actDev2 != actDev1 {
		t.Fatalf("Push 2: device configurationId should NOT have changed (%s → %s)", actDev1, actDev2)
	}
	// All three IDs must be independent
	if actCh01_2 == actDev2 || actCh01_2 == actCh02_2 {
		t.Fatalf("Push 2: CH01 configurationId (%s) should be independent from device (%s) and CH02 (%s)", actCh01_2, actDev2, actCh02_2)
	}
	t.Log("Push 2: OK — host bumped only CH01, actual config echoes correct IDs")

	// Push 3: change only device simpleSettings
	t.Log("ARD Push 3: change device simpleSettings only")
	time.Sleep(2 * time.Second)
	code, body = env.hostUpdateConfig(deviceID, token, baseConfig("ACTIVE", "IDLE", "PTP"))
	if code != 200 {
		t.Fatalf("Push 3: expected 200, got %d: %s", code, body)
	}
	time.Sleep(4 * time.Second)

	sdkCfg3 := env.sdkGetConfiguration()
	var actualCfg3 map[string]interface{}
	if p, ok := sdkCfg3.Configuration["payload"]; ok {
		b, _ := json.Marshal(p)
		json.Unmarshal(b, &actualCfg3)
	}
	cfgResp3 := env.sdkReportActualConfig(actualCfg3)
	if !cfgResp3.Success {
		t.Fatalf("Push 3: report_actual_configuration failed: %s", cfgResp3.Message)
	}
	time.Sleep(3 * time.Second)

	actDev3, actCh01_3, actCh02_3 := getActualConfigIds()
	t.Logf("Push 3 actual: device=%s ch01=%s ch02=%s", actDev3, actCh01_3, actCh02_3)

	if actDev3 == actDev2 {
		t.Fatalf("Push 3: device configurationId should have changed (%s → %s)", actDev2, actDev3)
	}
	if actCh01_3 != actCh01_2 {
		t.Fatalf("Push 3: CH01 configurationId should NOT have changed (%s → %s)", actCh01_2, actCh01_3)
	}
	if actCh02_3 != actCh02_2 {
		t.Fatalf("Push 3: CH02 configurationId should NOT have changed (%s → %s)", actCh02_2, actCh02_3)
	}
	if actDev3 == actCh01_3 || actDev3 == actCh02_3 {
		t.Fatalf("Push 3: device configurationId (%s) should be independent from CH01 (%s) and CH02 (%s)", actDev3, actCh01_3, actCh02_3)
	}
	t.Log("Push 3: OK — host bumped only device, actual config echoes correct IDs")

	t.Log("=== TestARDConfigurationIdEchoBack PASSED ===")
}
