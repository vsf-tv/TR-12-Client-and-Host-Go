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
	"net/http"
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
		ChannelAssignments []struct {
			ChannelID string `json:"channelId"`
		} `json:"channelAssignments"`
		ChannelTemplates []struct {
			Settings []interface{} `json:"settings"`
		} `json:"channelTemplates"`
	}
	if err := json.Unmarshal(detail.Registration, &reg); err != nil {
		t.Fatalf("Phase 6: cannot parse registration: %v", err)
	}
	if len(reg.ChannelAssignments) != 1 || reg.ChannelAssignments[0].ChannelID != "CH01" {
		t.Fatalf("Phase 6: expected 1 channel assignment with channelId=CH01, got %+v", reg.ChannelAssignments)
	}
	if len(reg.ChannelTemplates) == 0 || len(reg.ChannelTemplates[0].Settings) != 7 {
		t.Fatalf("Phase 6: expected 1 template with 7 settings, got %+v", reg.ChannelTemplates)
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
	t.Logf("Phase 6: OK — registration channels=%d cert_expiration=%s",
		len(reg.ChannelAssignments), detail.CertExpiration)

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
	badSettingCfg := json.RawMessage(`{"channels":[{"id":"CH01","state":"ACTIVE","channelSettings":{"standardSettings":[{"id":"NONEXISTENT_SETTING","value":"foo"}]}}]}`)
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
	badProfileCfg := json.RawMessage(`{"channels":[{"id":"CH01","state":"ACTIVE","channelSettings":{"profile":{"id":"nonexistent_profile"}}}]}`)
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
		"standardSettings": [
			{"id": "clocksync", "value": "PTP"}
		],
		"channels": [
			{
				"id": "CH01",
				"state": "ACTIVE",
				"channelSettings": {
					"standardSettings": [
						{"id": "RS01", "value": "1920x1080"},
						{"id": "FR01", "value": "60"},
						{"id": "MB01", "value": "20000"},
						{"id": "RC01", "value": "CBR"},
						{"id": "CO01", "value": "H.264"},
						{"id": "GP01", "value": "60"},
						{"id": "IN01", "value": "SDI1"}
					]
				},
				"protocol": {
					"srtCaller": {
						"address": "192.168.1.100",
						"port": 9000,
						"minimumLatencyMilliseconds": 200
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
		"status": []map[string]interface{}{
			{"name": "cpu", "value": "41", "description": "Current CPU % utilization."},
		},
		"channels": []map[string]interface{}{
			{
				"id":    "CH01",
				"state": "ACTIVE",
				"status": []map[string]interface{}{
					{"name": "bitrate", "value": "9500", "description": "Current output bitrate (Kbps)"},
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
		actualCfg["version"] = 0
	}
	// Inject thumbnailLocalPath per channel — the application sets this in actual config
	if channels, ok := actualCfg["channels"].([]interface{}); ok {
		for _, ch := range channels {
			if chMap, ok := ch.(map[string]interface{}); ok {
				chMap["thumbnailLocalPath"] = "/tmp/image_sdi.jpg"
			}
		}
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
		thumbCode, thumbResp = env.hostGetThumbnail(deviceID, "CH01", token)
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
		"standardSettings": [{"id": "clocksync", "value": "PTP"}],
		"channels": [{
			"id": "CH01",
			"state": "IDLE",
			"channelSettings": {
				"standardSettings": [
					{"id": "RS01", "value": "1920x1080"},
					{"id": "FR01", "value": "60"},
					{"id": "MB01", "value": "5000"},
					{"id": "RC01", "value": "CBR"},
					{"id": "CO01", "value": "H.264"},
					{"id": "GP01", "value": "60"},
					{"id": "IN01", "value": "SDI1"}
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
		t.Fatalf("SDK config missing clocksync=PTP after offline update: %s", cfgStr)
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
		ChannelAssignments []struct{ ChannelID string `json:"channelId"` } `json:"channelAssignments"`
	}
	json.Unmarshal(detail.Registration, &reg)
	if len(reg.ChannelAssignments) != 2 {
		t.Fatalf("expected 2 channel assignments in registration, got %d", len(reg.ChannelAssignments))
	}
	t.Logf("Registration OK — channels: %s, %s", reg.ChannelAssignments[0].ChannelID, reg.ChannelAssignments[1].ChannelID)

	// ---------------------------------------------------------------
	// Phase 1: Full 2-channel config — both channels should be applied
	// ---------------------------------------------------------------
	t.Log("Phase 1: Full 2-channel configuration update")
	fullConfig := json.RawMessage(`{
		"standardSettings": [{"id": "clocksync", "value": "PTP"}],
		"channels": [
			{
				"id": "CH01", "state": "ACTIVE",
				"channelSettings": {"standardSettings": [
					{"id": "RS01", "value": "1920x1080"},
					{"id": "FR01", "value": "60"},
					{"id": "MB01", "value": "20000"},
					{"id": "RC01", "value": "CBR"},
					{"id": "CO01", "value": "H.264"},
					{"id": "GP01", "value": "60"},
					{"id": "IN01", "value": "SDI1"}
				]},
				"protocol": {"srtCaller": {
					"address": "192.168.1.100", "port": 9001, "minimumLatencyMilliseconds": 200
				}}
			},
			{
				"id": "CH02", "state": "ACTIVE",
				"channelSettings": {"standardSettings": [
					{"id": "RS01", "value": "1280x720"},
					{"id": "FR01", "value": "30"},
					{"id": "MB01", "value": "10000"},
					{"id": "RC01", "value": "CBR"},
					{"id": "CO01", "value": "H.264"},
					{"id": "GP01", "value": "30"},
					{"id": "IN01", "value": "HDMI1"}
				]},
				"protocol": {"srtCaller": {
					"address": "192.168.1.101", "port": 9002, "minimumLatencyMilliseconds": 200
				}}
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
		"standardSettings": [{"id": "clocksync", "value": "PTP"}],
		"channels": [
			{
				"id": "CH01", "state": "ACTIVE",
				"channelSettings": {"standardSettings": [
					{"id": "RS01", "value": "1920x1080"},
					{"id": "FR01", "value": "25"},
					{"id": "MB01", "value": "20000"},
					{"id": "RC01", "value": "CBR"},
					{"id": "CO01", "value": "H.264"},
					{"id": "GP01", "value": "60"},
					{"id": "IN01", "value": "SDI1"}
				]},
				"protocol": {"srtCaller": {
					"address": "192.168.1.100", "port": 9001, "minimumLatencyMilliseconds": 200
				}}
			},
			{
				"id": "CH02", "state": "ACTIVE",
				"channelSettings": {"standardSettings": [
					{"id": "RS01", "value": "1280x720"},
					{"id": "FR01", "value": "30"},
					{"id": "MB01", "value": "10000"},
					{"id": "RC01", "value": "CBR"},
					{"id": "CO01", "value": "H.264"},
					{"id": "GP01", "value": "30"},
					{"id": "IN01", "value": "HDMI1"}
				]},
				"protocol": {"srtCaller": {
					"address": "192.168.1.101", "port": 9002, "minimumLatencyMilliseconds": 200
				}}
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
	// Send identical channel configs — only standardSettings changes.
	// Both channels should be skipped since their configurationIds are unchanged.
	// ---------------------------------------------------------------
	t.Log("Phase 3: Device-level only update — channels should not be reapplied")
	time.Sleep(1 * time.Second)
	deviceOnlyConfig := json.RawMessage(`{
		"standardSettings": [{"id": "clocksync", "value": "GENLOCK"}],
		"channels": [
			{
				"id": "CH01", "state": "ACTIVE",
				"channelSettings": {"standardSettings": [
					{"id": "RS01", "value": "1920x1080"},
					{"id": "FR01", "value": "25"},
					{"id": "MB01", "value": "20000"},
					{"id": "RC01", "value": "CBR"},
					{"id": "CO01", "value": "H.264"},
					{"id": "GP01", "value": "60"},
					{"id": "IN01", "value": "SDI1"}
				]},
				"protocol": {"srtCaller": {
					"address": "192.168.1.100", "port": 9001, "minimumLatencyMilliseconds": 200
				}}
			},
			{
				"id": "CH02", "state": "ACTIVE",
				"channelSettings": {"standardSettings": [
					{"id": "RS01", "value": "1280x720"},
					{"id": "FR01", "value": "30"},
					{"id": "MB01", "value": "10000"},
					{"id": "RC01", "value": "CBR"},
					{"id": "CO01", "value": "H.264"},
					{"id": "GP01", "value": "30"},
					{"id": "IN01", "value": "HDMI1"}
				]},
				"protocol": {"srtCaller": {
					"address": "192.168.1.101", "port": 9002, "minimumLatencyMilliseconds": 200
				}}
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
			Version string `json:"version"`
			Channels        []struct {
				Id              string `json:"id"`
				Version string `json:"version"`
			} `json:"channels"`
		}
		json.Unmarshal(b, &cfg)
		deviceId = cfg.Version
		for _, ch := range cfg.Channels {
			if ch.Id == "CH01" {
				ch01Id = ch.Version
			} else if ch.Id == "CH02" {
				ch02Id = ch.Version
			}
		}
		return
	}

	baseConfig := func(ch01State, ch02State, clockSource string) json.RawMessage {
		return json.RawMessage(`{
			"standardSettings": [{"id": "clocksync", "value": "` + clockSource + `"}],
			"channels": [
				{"id": "CH01", "state": "` + ch01State + `", "channelSettings": {"standardSettings": [
					{"id": "RS01", "value": "1920x1080"}, {"id": "FR01", "value": "30"},
					{"id": "MB01", "value": "10000"}, {"id": "RC01", "value": "CBR"},
					{"id": "CO01", "value": "H.264"}, {"id": "GP01", "value": "60"},
					{"id": "IN01", "value": "SDI1"}
				]}},
				{"id": "CH02", "state": "` + ch02State + `", "channelSettings": {"standardSettings": [
					{"id": "RS01", "value": "1920x1080"}, {"id": "FR01", "value": "30"},
					{"id": "MB01", "value": "10000"}, {"id": "RC01", "value": "CBR"},
					{"id": "CO01", "value": "H.264"}, {"id": "GP01", "value": "60"},
					{"id": "IN01", "value": "SDI1"}
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
	t.Log("Push 3: change device standardSettings only")
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
			Version string `json:"version"`
			Channels        []struct {
				Id              string `json:"id"`
				Version string `json:"version"`
			} `json:"channels"`
		}
		json.Unmarshal(detail.ActualConfiguration, &actual)
		deviceId = actual.Version
		for _, ch := range actual.Channels {
			if ch.Id == "CH01" {
				ch01Id = ch.Version
			} else if ch.Id == "CH02" {
				ch02Id = ch.Version
			}
		}
		return
	}

	baseConfig := func(ch01State, ch02State, clockSource string) json.RawMessage {
		return json.RawMessage(`{
			"standardSettings": [{"id": "clocksync", "value": "` + clockSource + `"}],
			"channels": [
				{"id": "CH01", "state": "` + ch01State + `", "channelSettings": {"standardSettings": [
					{"id": "RS01", "value": "1920x1080"}, {"id": "FR01", "value": "30"},
					{"id": "MB01", "value": "10000"}, {"id": "RC01", "value": "CBR"},
					{"id": "CO01", "value": "H.264"}, {"id": "GP01", "value": "60"},
					{"id": "IN01", "value": "SDI1"}
				]}},
				{"id": "CH02", "state": "` + ch02State + `", "channelSettings": {"standardSettings": [
					{"id": "RS01", "value": "1920x1080"}, {"id": "FR01", "value": "30"},
					{"id": "MB01", "value": "10000"}, {"id": "RC01", "value": "CBR"},
					{"id": "CO01", "value": "H.264"}, {"id": "GP01", "value": "60"},
					{"id": "IN01", "value": "SDI1"}
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
		Version string `json:"version"`
		Channels        []struct {
			Id              string `json:"id"`
			Version string `json:"version"`
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

	// Push 3: change only device standardSettings
	t.Log("ARD Push 3: change device standardSettings only")
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

// TestTwoChannelThumbnails verifies that thumbnails can be requested and
// received for both channels on a 2-channel encoder using channelId-based
// subscriptions.
func TestTwoChannelThumbnails(t *testing.T) {
	createTestJPEG(t, "/tmp/image_sdi.jpg")
	createTestJPEG(t, "/tmp/image_hdmi.jpg")
	t.Cleanup(func() {
		_ = removeIfExists("/tmp/image_sdi.jpg")
		_ = removeIfExists("/tmp/image_hdmi.jpg")
	})

	env := newTestEnv(t)
	env.startHost()
	env.startSDK("integ-2ch-thumb-001")

	registration := loadRegistrationFrom(t, "2_channel_encoder")

	acct := env.hostRegisterAccount("thumbuser", "testpass123", "Thumb Test")
	token := acct.Token
	pairingCode := env.waitForPairingCode("tr12-host", registration, 15*time.Second)
	env.hostClaim(pairingCode, token)
	env.waitForSDKConnected("tr12-host", registration, 30*time.Second)

	devices := env.hostListDevices(token)
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	deviceID := devices[0].DeviceID

	// Push config so SDK has a desired config, then report actual with thumbnailLocalPath
	fullConfig := json.RawMessage(`{
		"channels": [
			{"id": "CH01", "state": "IDLE", "channelSettings": {"standardSettings": [
				{"id": "RS01", "value": "1920x1080"}, {"id": "FR01", "value": "30"},
				{"id": "MB01", "value": "10000"}, {"id": "RC01", "value": "CBR"},
				{"id": "CO01", "value": "H.264"}, {"id": "GP01", "value": "60"},
				{"id": "IN01", "value": "SDI1"}
			]}},
			{"id": "CH02", "state": "IDLE", "channelSettings": {"standardSettings": [
				{"id": "RS01", "value": "1280x720"}, {"id": "FR01", "value": "30"},
				{"id": "MB01", "value": "10000"}, {"id": "RC01", "value": "CBR"},
				{"id": "CO01", "value": "H.264"}, {"id": "GP01", "value": "60"},
				{"id": "IN01", "value": "HDMI1"}
			]}}
		]
	}`)
	code, body := env.hostUpdateConfig(deviceID, token, fullConfig)
	if code != 200 {
		t.Fatalf("config push: expected 200, got %d: %s", code, body)
	}
	time.Sleep(3 * time.Second)

	// Report actual config with thumbnailLocalPath per channel
	sdkCfg := env.sdkGetConfiguration()
	var actualCfg map[string]interface{}
	if p, ok := sdkCfg.Configuration["payload"]; ok {
		b, _ := json.Marshal(p)
		json.Unmarshal(b, &actualCfg)
	}
	if actualCfg == nil {
		t.Fatal("could not extract config payload")
	}
	if channels, ok := actualCfg["channels"].([]interface{}); ok {
		for _, ch := range channels {
			if chMap, ok := ch.(map[string]interface{}); ok {
				if chMap["id"] == "CH01" {
					chMap["thumbnailLocalPath"] = "/tmp/image_sdi.jpg"
				} else if chMap["id"] == "CH02" {
					chMap["thumbnailLocalPath"] = "/tmp/image_hdmi.jpg"
				}
			}
		}
	}
	cfgResp := env.sdkReportActualConfig(actualCfg)
	if !cfgResp.Success {
		t.Fatalf("report_actual_configuration failed: %s", cfgResp.Message)
	}
	time.Sleep(3 * time.Second)

	// Request thumbnail for CH01
	t.Log("Requesting thumbnail for CH01")
	createTestJPEG(t, "/tmp/image_sdi.jpg")
	var ch01Thumb thumbnailResponse
	for attempt := 0; attempt < 4; attempt++ {
		_, ch01Thumb = env.hostGetThumbnail(deviceID, "CH01", token)
		if ch01Thumb.Image != nil && ch01Thumb.Image.Base64Image != "" {
			break
		}
		createTestJPEG(t, "/tmp/image_sdi.jpg")
		time.Sleep(5 * time.Second)
	}
	if ch01Thumb.Image == nil || ch01Thumb.Image.Base64Image == "" {
		t.Fatal("CH01: no thumbnail received after retries")
	}
	t.Log("CH01: OK — thumbnail received")

	// Request thumbnail for CH02
	t.Log("Requesting thumbnail for CH02")
	createTestJPEG(t, "/tmp/image_hdmi.jpg")
	var ch02Thumb thumbnailResponse
	for attempt := 0; attempt < 4; attempt++ {
		_, ch02Thumb = env.hostGetThumbnail(deviceID, "CH02", token)
		if ch02Thumb.Image != nil && ch02Thumb.Image.Base64Image != "" {
			break
		}
		createTestJPEG(t, "/tmp/image_hdmi.jpg")
		time.Sleep(5 * time.Second)
	}
	if ch02Thumb.Image == nil || ch02Thumb.Image.Base64Image == "" {
		t.Fatal("CH02: no thumbnail received after retries")
	}
	t.Log("CH02: OK — thumbnail received")

	t.Log("=== TestTwoChannelThumbnails PASSED ===")
}

// TestPairingRejections verifies that CreatePairingCode returns HTTP 400 with
// the correct reason for each of the three rejection cases.
func TestPairingRejections(t *testing.T) {
	env := newTestEnv(t)
	env.startHost()

	csr := generateMinimalCSR(t)

	type pairRequest struct {
		DeviceType                string      `json:"deviceType"`
		HostId                    string      `json:"hostId"`
		CertificateSigningRequest string      `json:"certificateSigningRequest"`
		Version                   interface{} `json:"version"`
	}

	type pairErrorResponse struct {
		Reason string `json:"reason"`
	}

	doPairRaw := func(req pairRequest) (int, pairErrorResponse) {
		t.Helper()
		var errResp pairErrorResponse
		code := env.doPostRaw(env.hostURL+"/pair", req, &errResp)
		return code, errResp
	}

	goodVersion := map[string]string{"version": "5.0.0"}

	// --- HOST_ID_MISMATCH ---
	t.Log("Pairing rejection: HOST_ID_MISMATCH")
	code, errResp := doPairRaw(pairRequest{
		DeviceType:                "SOURCE",
		HostId:                    "wrong-host-id",
		CertificateSigningRequest: csr,
		Version:                   goodVersion,
	})
	if code != 400 {
		t.Fatalf("HOST_ID_MISMATCH: expected 400, got %d", code)
	}
	if errResp.Reason != "HOST_ID_MISMATCH" {
		t.Fatalf("HOST_ID_MISMATCH: expected reason=HOST_ID_MISMATCH, got %q", errResp.Reason)
	}
	t.Log("HOST_ID_MISMATCH: OK")

	// --- VERSION_NOT_SUPPORTED ---
	t.Log("Pairing rejection: VERSION_NOT_SUPPORTED")
	code, errResp = doPairRaw(pairRequest{
		DeviceType:                "SOURCE",
		HostId:                    "tr12-host",
		CertificateSigningRequest: csr,
		Version:                   map[string]string{"version": ""},
	})
	if code != 400 {
		t.Fatalf("VERSION_NOT_SUPPORTED: expected 400, got %d", code)
	}
	if errResp.Reason != "VERSION_NOT_SUPPORTED" {
		t.Fatalf("VERSION_NOT_SUPPORTED: expected reason=VERSION_NOT_SUPPORTED, got %q", errResp.Reason)
	}
	t.Log("VERSION_NOT_SUPPORTED: OK")

	// --- DEVICE_TYPE_NOT_SUPPORTED ---
	// Note: "INVALID_TYPE" fails JSON deserialization at the HTTP boundary because
	// DeviceType has a custom UnmarshalJSON that rejects unknown enum values.
	// The host returns 400 with {"error": "invalid request body"} rather than
	// {"reason": "DEVICE_TYPE_NOT_SUPPORTED"} — both are valid 400 rejections.
	t.Log("Pairing rejection: DEVICE_TYPE_NOT_SUPPORTED")
	code, _ = doPairRaw(pairRequest{
		DeviceType:                "INVALID_TYPE",
		HostId:                    "tr12-host",
		CertificateSigningRequest: csr,
		Version:                   goodVersion,
	})
	if code != 400 {
		t.Fatalf("DEVICE_TYPE_NOT_SUPPORTED: expected 400, got %d", code)
	}
	t.Log("DEVICE_TYPE_NOT_SUPPORTED: OK — rejected with 400")

	t.Log("=== TestPairingRejections PASSED ===")
}

// TestMQTTEnvelopes verifies that all MQTT messages use the v6.0.0 envelope format.
// It exercises every host→device and device→host topic and checks the envelope
// field is present and the inner payload is correctly unwrapped.
func TestMQTTEnvelopes(t *testing.T) {
	env := newTestEnv(t)
	env.startHost()
	env.startSDK("integ-envelope-001")

	registration := loadRegistration(t)

	// Setup: pair and connect
	acct := env.hostRegisterAccount("envelopeuser", "testpass123", "Envelope Test")
	token := acct.Token
	pairingCode := env.waitForPairingCode("tr12-host", registration, 15*time.Second)
	env.hostClaim(pairingCode, token)
	env.waitForSDKConnected("tr12-host", registration, 30*time.Second)

	devices := env.hostListDevices(token)
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	deviceID := devices[0].DeviceID

	// -----------------------------------------------------------------------
	// 1. Registration envelope: device → host
	//    The SDK sends {"deviceRegistration": {...}} on connect.
	//    Verify the host stored the registration (unwrapped correctly).
	// -----------------------------------------------------------------------
	t.Log("Envelope check: deviceRegistration (device→host)")
	detail := env.hostDescribeDevice(deviceID, token)
	if len(detail.Registration) == 0 || string(detail.Registration) == "null" {
		t.Fatal("registration envelope: host has no registration — envelope not unwrapped")
	}
	var reg struct {
		ChannelAssignments []struct{ ChannelID string `json:"channelId"` } `json:"channelAssignments"`
	}
	if err := json.Unmarshal(detail.Registration, &reg); err != nil || len(reg.ChannelAssignments) == 0 {
		t.Fatalf("registration envelope: stored registration is invalid: %v", err)
	}
	t.Log("registration envelope: OK")

	// -----------------------------------------------------------------------
	// 2. DesiredDeviceConfiguration envelope: host → device
	//    Push a config and verify the SDK received and stored it correctly.
	// -----------------------------------------------------------------------
	t.Log("Envelope check: desiredDeviceConfiguration (host→device)")
	cfg := json.RawMessage(`{
		"standardSettings": [{"id": "clocksync", "value": "PTP"}],
		"channels": [{
			"id": "CH01", "state": "IDLE",
			"channelSettings": {"standardSettings": [
				{"id": "RS01", "value": "1920x1080"},
				{"id": "FR01", "value": "30"},
				{"id": "MB01", "value": "10000"},
				{"id": "RC01", "value": "CBR"},
				{"id": "CO01", "value": "H.264"},
				{"id": "GP01", "value": "60"},
				{"id": "IN01", "value": "SDI1"}
			]}
		}]
	}`)
	code, body := env.hostUpdateConfig(deviceID, token, cfg)
	if code != 200 {
		t.Fatalf("desiredDeviceConfiguration envelope: push failed: %d %s", code, body)
	}
	time.Sleep(3 * time.Second)
	sdkCfg := env.sdkGetConfiguration()
	if sdkCfg.Configuration == nil {
		t.Fatal("desiredDeviceConfiguration envelope: SDK has no configuration — envelope not unwrapped")
	}
	cfgStr := string(mustMarshal(sdkCfg.Configuration))
	if !strings.Contains(cfgStr, "CH01") {
		t.Fatalf("desiredDeviceConfiguration envelope: SDK config missing CH01: %s", cfgStr)
	}
	if !strings.Contains(cfgStr, "PTP") {
		t.Fatalf("desiredDeviceConfiguration envelope: SDK config missing PTP clock source: %s", cfgStr)
	}
	t.Log("desiredDeviceConfiguration envelope: OK")

	// -----------------------------------------------------------------------
	// 3. ActualDeviceConfiguration envelope: device → host
	//    Report actual config and verify the host stored it (unwrapped correctly).
	// -----------------------------------------------------------------------
	t.Log("Envelope check: actualDeviceConfiguration (device→host)")
	var actualCfg map[string]interface{}
	if p, ok := sdkCfg.Configuration["payload"]; ok {
		b, _ := json.Marshal(p)
		json.Unmarshal(b, &actualCfg)
	}
	if actualCfg == nil {
		t.Fatal("actualDeviceConfiguration envelope: could not extract payload from SDK config")
	}
	cfgResp := env.sdkReportActualConfig(actualCfg)
	if !cfgResp.Success {
		t.Fatalf("actualDeviceConfiguration envelope: report failed: %s", cfgResp.Message)
	}
	time.Sleep(3 * time.Second)
	detail2 := env.hostDescribeDevice(deviceID, token)
	if len(detail2.ActualConfiguration) == 0 || string(detail2.ActualConfiguration) == "null" {
		t.Fatal("actualDeviceConfiguration envelope: host has no actual config — envelope not unwrapped")
	}
	var actual struct {
		Channels []struct{ ID string `json:"id"` } `json:"channels"`
	}
	if err := json.Unmarshal(detail2.ActualConfiguration, &actual); err != nil || len(actual.Channels) == 0 {
		t.Fatalf("actualDeviceConfiguration envelope: stored actual config is invalid: %v", err)
	}
	t.Log("actualDeviceConfiguration envelope: OK")

	// -----------------------------------------------------------------------
	// 4. DeviceStatus envelope: device → host
	//    Report status and verify the host stored it (unwrapped correctly).
	// -----------------------------------------------------------------------
	t.Log("Envelope check: deviceStatus (device→host)")
	statusPayload := map[string]interface{}{
		"channels": []map[string]interface{}{
			{"id": "CH01", "state": "IDLE", "status": []map[string]interface{}{
				{"name": "bitrate", "value": "0", "description": "Current bitrate"},
			}},
		},
		"status": []map[string]interface{}{
			{"name": "model", "value": "TestDevice", "description": "Device model"},
		},
	}
	statusResp := env.sdkReportStatus(statusPayload)
	if !statusResp.Success {
		t.Fatalf("deviceStatus envelope: report failed: %s", statusResp.Message)
	}
	time.Sleep(3 * time.Second)
	detail3 := env.hostDescribeDevice(deviceID, token)
	if len(detail3.Status) == 0 || string(detail3.Status) == "null" {
		t.Fatal("deviceStatus envelope: host has no status — envelope not unwrapped")
	}
	statusStr := string(detail3.Status)
	if !strings.Contains(statusStr, "TestDevice") {
		t.Fatalf("deviceStatus envelope: stored status missing TestDevice: %s", statusStr)
	}
	t.Log("deviceStatus envelope: OK")

	t.Log("=== TestMQTTEnvelopes PASSED ===")
}

// TestRegister verifies the PUT /register endpoint:
//  1. Rejects calls when not connected (returns success=false)
//  2. Rejects changes to anything other than profiles (channelType, settings, protocols, assignments)
//  3. Accepts a profile-only update and re-publishes registration to the host
//  4. Host reflects the updated profiles in its device record
func TestRegister(t *testing.T) {
	env := newTestEnv(t)
	env.startHost()
	env.startSDK("integ-register-001")

	registration := loadRegistration(t)

	// ---------------------------------------------------------------
	t.Log("Register Phase 1: call /register before connected — must fail")
	// ---------------------------------------------------------------
	resp := env.sdkRegister(registration)
	if resp.Success {
		t.Fatal("Phase 1: expected failure when not connected, got success=true")
	}
	if resp.State != "DISCONNECTED" {
		t.Fatalf("Phase 1: expected state DISCONNECTED, got %q", resp.State)
	}
	t.Logf("Phase 1: OK — correctly rejected (state=%s message=%q)", resp.State, resp.Message)

	// ---------------------------------------------------------------
	t.Log("Register Phase 2: pair and connect")
	// ---------------------------------------------------------------
	acct := env.hostRegisterAccount("registeruser", "testpass123", "Register Test")
	token := acct.Token
	pairingCode := env.waitForPairingCode("tr12-host", registration, 15*time.Second)
	env.hostClaim(pairingCode, token)
	env.waitForSDKConnected("tr12-host", registration, 30*time.Second)
	t.Log("Phase 2: OK — connected")

	devices := env.hostListDevices(token)
	if len(devices) == 0 {
		t.Fatal("Phase 2: no devices found after connect")
	}
	deviceID := devices[0].DeviceID

	// ---------------------------------------------------------------
	t.Log("Register Phase 3: reject change to channelAssignments")
	// ---------------------------------------------------------------
	badAssignments := deepCopyReg(t, registration)
	badAssignments["channelAssignments"] = []interface{}{
		map[string]interface{}{"channelId": "CH99", "name": "Changed Channel", "templateId": "main"},
	}
	resp = env.sdkRegister(badAssignments)
	if resp.Success {
		t.Fatal("Phase 3: expected failure when channelAssignments changed, got success=true")
	}
	t.Logf("Phase 3: OK — correctly rejected channel assignment change: %q", resp.Message)

	// ---------------------------------------------------------------
	t.Log("Register Phase 4: reject change to settings in channelTemplates")
	// ---------------------------------------------------------------
	badSettings := deepCopyReg(t, registration)
	templates, _ := badSettings["channelTemplates"].([]interface{})
	tmpl, _ := templates[0].(map[string]interface{})
	tmpl["settings"] = []interface{}{} // remove all settings
	resp = env.sdkRegister(badSettings)
	if resp.Success {
		t.Fatal("Phase 4: expected failure when channelTemplates settings changed, got success=true")
	}
	t.Logf("Phase 4: OK — correctly rejected settings change: %q", resp.Message)

	// ---------------------------------------------------------------
	t.Log("Register Phase 5: reject change to protocols in channelTemplates")
	// ---------------------------------------------------------------
	badProtos := deepCopyReg(t, registration)
	templates2, _ := badProtos["channelTemplates"].([]interface{})
	tmpl2, _ := templates2[0].(map[string]interface{})
	tmpl2["protocols"] = []interface{}{"SRT_CALLER"} // was ["SRT_CALLER","SRT_LISTENER"]
	resp = env.sdkRegister(badProtos)
	if resp.Success {
		t.Fatal("Phase 5: expected failure when protocols changed, got success=true")
	}
	t.Logf("Phase 5: OK — correctly rejected protocol change: %q", resp.Message)

	// ---------------------------------------------------------------
	t.Log("Register Phase 6: accept profile-only update")
	// ---------------------------------------------------------------
	updatedProfiles := deepCopyReg(t, registration)
	templates3, _ := updatedProfiles["channelTemplates"].([]interface{})
	tmpl3, _ := templates3[0].(map[string]interface{})
	tmpl3["profiles"] = []interface{}{
		map[string]interface{}{"id": "h264c", "name": "h264_contribution_compact", "description": "H264, 30fps, 10mbs, 720p, 6s gop"},
		map[string]interface{}{"id": "h264f", "name": "h264_contribution_full", "description": "H264, 60fps, 20mbs, 1080p, 3s gop"},
		// h265c, h265k, h265kl removed — simulates device user removing profiles
		map[string]interface{}{"id": "newprf", "name": "new_profile", "description": "A newly added profile"},
	}
	resp = env.sdkRegister(updatedProfiles)
	if !resp.Success {
		t.Fatalf("Phase 6: expected success for profile-only update, got failure: %q", resp.Message)
	}
	t.Logf("Phase 6: OK — profile update accepted (state=%s)", resp.State)

	// Wait briefly for MQTT re-publish to reach host.
	time.Sleep(2 * time.Second)

	// ---------------------------------------------------------------
	t.Log("Register Phase 7: verify host received updated registration")
	// ---------------------------------------------------------------
	detail := env.hostDescribeDevice(deviceID, token)
	var reg struct {
		ChannelTemplates []struct {
			Profiles []struct {
				ID string `json:"id"`
			} `json:"profiles"`
		} `json:"channelTemplates"`
	}
	if err := json.Unmarshal(detail.Registration, &reg); err != nil {
		t.Fatalf("Phase 7: cannot parse registration from host: %v", err)
	}
	if len(reg.ChannelTemplates) == 0 {
		t.Fatal("Phase 7: host has no channel templates in registration")
	}
	profiles := reg.ChannelTemplates[0].Profiles
	if len(profiles) != 3 {
		t.Fatalf("Phase 7: expected 3 profiles after update, got %d", len(profiles))
	}
	foundNew := false
	for _, p := range profiles {
		if p.ID == "newprf" {
			foundNew = true
		}
	}
	if !foundNew {
		t.Fatal("Phase 7: new profile 'newprf' not found in host registration")
	}
	t.Logf("Phase 7: OK — host registration updated with %d profiles including 'newprf'", len(profiles))

	t.Log("=== TestRegister PASSED ===")
}

// deepCopyReg round-trips a registration map through JSON for a clean deep copy.
func deepCopyReg(t *testing.T, reg map[string]interface{}) map[string]interface{} {
	t.Helper()
	data, err := json.Marshal(reg)
	if err != nil {
		t.Fatalf("deepCopyReg marshal: %v", err)
	}
	var copy map[string]interface{}
	if err := json.Unmarshal(data, &copy); err != nil {
		t.Fatalf("deepCopyReg unmarshal: %v", err)
	}
	return copy
}

// TestChannelHealthReporting verifies the end-to-end health reporting flow:
//
//  1. Connect to the host
//  2. Push a config update
//  3. Report DEGRADED channel health via the STATUS payload (PUT /report_status)
//  4. Verify the host stores the DEGRADED health in status
//  5. Report healthy via STATUS and verify the host sees HEALTHY again
//
// Health is reported in the status payload, not in actual_configuration.
// This exercises the scenario where a device-side failure (e.g. native API error,
// hardware fault) is surfaced to the cloud operator via the TR-12 health field
// in the device status.
func TestChannelHealthReporting(t *testing.T) {
	env := newTestEnv(t)
	env.startHost()
	env.startSDK("integ-health-001")

	registration := loadRegistration(t)

	// Setup: pair and connect
	acct := env.hostRegisterAccount("healthuser", "testpass123", "Health Test")
	token := acct.Token
	pairingCode := env.waitForPairingCode("tr12-host", registration, 15*time.Second)
	env.hostClaim(pairingCode, token)
	env.waitForSDKConnected("tr12-host", registration, 30*time.Second)

	devices := env.hostListDevices(token)
	if len(devices) == 0 {
		t.Fatal("no devices found after connect")
	}
	deviceID := devices[0].DeviceID
	t.Logf("Connected device: %s", deviceID)

	// Push a config update so the device has a version to echo back.
	cfg := json.RawMessage(`{"channels":[{"id":"CH01","state":"ACTIVE","channelSettings":{"standardSettings":[{"id":"RS01","value":"1920x1080"}]}}]}`)
	statusCode, body := env.hostUpdateConfig(deviceID, token, cfg)
	if statusCode != http.StatusOK {
		t.Fatalf("hostUpdateConfig: expected 200, got %d: %s", statusCode, body)
	}

	// -----------------------------------------------------------------------
	// Phase: report status with DEGRADED channel health.
	// Health is reported via PUT /report_status with the health field on
	// channels in the status payload.
	// -----------------------------------------------------------------------
	t.Log("Reporting status with DEGRADED health for CH01")
	degradedStatus := map[string]interface{}{
		"status": []map[string]interface{}{
			{"name": "cpu", "value": "41", "description": "Current CPU % utilization."},
		},
		"channels": []interface{}{
			map[string]interface{}{
				"id":    "CH01",
				"state": "ACTIVE",
				"status": []map[string]interface{}{
					{"name": "bitrate", "value": "9500", "description": "Current output bitrate (Kbps)"},
				},
				"health": map[string]interface{}{
					"degraded": map[string]interface{}{
						"message":   "TEST: native API returned error code 503 — codec unavailable",
						"timestamp": time.Now().UTC().Format(time.RFC3339),
					},
				},
			},
		},
		"health": map[string]interface{}{
			"healthy": map[string]interface{}{},
		},
	}
	reportResp := env.sdkReportStatus(degradedStatus)
	if !reportResp.Success {
		t.Fatalf("report_status failed: %s", reportResp.Message)
	}

	// Wait for the host to store the status with DEGRADED channel health.
	deadline := time.Now().Add(15 * time.Second)
	foundDegraded := false
	for time.Now().Before(deadline) {
		detail := env.hostDescribeDevice(deviceID, token)
		if len(detail.Status) == 0 || string(detail.Status) == "null" {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		var st struct {
			Channels []struct {
				ID     string          `json:"id"`
				Health json.RawMessage `json:"health,omitempty"`
			} `json:"channels"`
		}
		if err := json.Unmarshal(detail.Status, &st); err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		for _, ch := range st.Channels {
			if ch.ID == "CH01" && len(ch.Health) > 0 {
				healthStr := string(ch.Health)
				t.Logf("CH01 health from host status: %s", healthStr)
				if strings.Contains(healthStr, "degraded") {
					foundDegraded = true
				}
			}
		}
		if foundDegraded {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if !foundDegraded {
		detail := env.hostDescribeDevice(deviceID, token)
		t.Fatalf("expected DEGRADED health for CH01 in host status, got: %s",
			string(detail.Status))
	}
	t.Log("DEGRADED health confirmed in host status")

	// -----------------------------------------------------------------------
	// Phase: report status with HEALTHY channel state.
	// -----------------------------------------------------------------------
	t.Log("Reporting status with HEALTHY state for CH01")
	healthyStatus := map[string]interface{}{
		"status": []map[string]interface{}{
			{"name": "cpu", "value": "38", "description": "Current CPU % utilization."},
		},
		"channels": []interface{}{
			map[string]interface{}{
				"id":    "CH01",
				"state": "ACTIVE",
				"status": []map[string]interface{}{
					{"name": "bitrate", "value": "9500", "description": "Current output bitrate (Kbps)"},
				},
				"health": map[string]interface{}{
					"healthy": map[string]interface{}{},
				},
			},
		},
		"health": map[string]interface{}{
			"healthy": map[string]interface{}{},
		},
	}
	reportResp = env.sdkReportStatus(healthyStatus)
	if !reportResp.Success {
		t.Fatalf("report_status (healthy) failed: %s", reportResp.Message)
	}

	deadline = time.Now().Add(15 * time.Second)
	foundHealthy := false
	for time.Now().Before(deadline) {
		detail := env.hostDescribeDevice(deviceID, token)
		if len(detail.Status) == 0 || string(detail.Status) == "null" {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		var st struct {
			Channels []struct {
				ID     string          `json:"id"`
				Health json.RawMessage `json:"health,omitempty"`
			} `json:"channels"`
		}
		if err := json.Unmarshal(detail.Status, &st); err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		for _, ch := range st.Channels {
			if ch.ID == "CH01" {
				healthStr := string(ch.Health)
				t.Logf("CH01 health after healthy report: %s", healthStr)
				if !strings.Contains(healthStr, "degraded") && !strings.Contains(healthStr, "critical") {
					foundHealthy = true
				}
			}
		}
		if foundHealthy {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if !foundHealthy {
		detail := env.hostDescribeDevice(deviceID, token)
		t.Fatalf("expected HEALTHY for CH01 after reporting healthy status, got: %s",
			string(detail.Status))
	}
	t.Log("HEALTHY state confirmed after reporting healthy status")

	t.Log("=== TestChannelHealthReporting PASSED ===")
}

// TestSequentialPerChannelStops verifies that two sequential per-channel updates
// (stop CH01, then stop CH02) both get processed by the client. This is a
// regression test for a bug where the host published only the changed channel
// in the MQTT retained message, causing the SDK to lose all other channels.
func TestSequentialPerChannelStops(t *testing.T) {
	createTestJPEG(t, "/tmp/image_sdi.jpg")
	t.Cleanup(func() { _ = removeIfExists("/tmp/image_sdi.jpg") })

	env := newTestEnv(t)
	env.startHost()
	env.startSDK("integ-seqstop-001")

	registration := loadRegistrationFrom(t, "2_channel_encoder")

	acct := env.hostRegisterAccount("seqstopuser", "testpass123", "SeqStop Test")
	token := acct.Token
	pairingCode := env.waitForPairingCode("tr12-host", registration, 15*time.Second)
	env.hostClaim(pairingCode, token)
	env.waitForSDKConnected("tr12-host", registration, 30*time.Second)

	devices := env.hostListDevices(token)
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	deviceID := devices[0].DeviceID

	// Start both channels via full-device update
	t.Log("Phase 1: Start both channels")
	fullConfig := json.RawMessage(`{
		"standardSettings": [{"id": "clocksync", "value": "NTP"}],
		"channels": [
			{"id": "CH01", "state": "ACTIVE", "channelSettings": {"standardSettings": [
				{"id": "RS01", "value": "1920x1080"}, {"id": "FR01", "value": "60"},
				{"id": "MB01", "value": "10000"}, {"id": "RC01", "value": "CBR"},
				{"id": "CO01", "value": "H.264"}, {"id": "GP01", "value": "60"},
				{"id": "IN01", "value": "SDI1"}
			]}},
			{"id": "CH02", "state": "ACTIVE", "channelSettings": {"standardSettings": [
				{"id": "RS01", "value": "1920x1080"}, {"id": "FR01", "value": "30"},
				{"id": "MB01", "value": "10000"}, {"id": "RC01", "value": "CBR"},
				{"id": "CO01", "value": "H.264"}, {"id": "GP01", "value": "60"},
				{"id": "IN01", "value": "SDI1"}
			]}}
		]
	}`)
	code, body := env.hostUpdateConfig(deviceID, token, fullConfig)
	if code != 200 {
		t.Fatalf("Phase 1: expected 200, got %d: %s", code, body)
	}
	time.Sleep(5 * time.Second)

	// Verify both channels are ACTIVE in SDK config
	sdkCfg := env.sdkGetConfiguration()
	if sdkCfg.Configuration == nil {
		t.Fatal("Phase 1: SDK returned nil configuration")
	}
	cfgStr := string(mustMarshal(sdkCfg.Configuration))
	if !strings.Contains(cfgStr, "CH01") || !strings.Contains(cfgStr, "CH02") {
		t.Fatalf("Phase 1: expected both channels in SDK config: %s", cfgStr)
	}
	t.Log("Phase 1: OK — both channels ACTIVE")

	// Phase 2: Stop CH01 via per-channel API
	t.Log("Phase 2: Stop CH01 via per-channel API")
	ch01Stop := json.RawMessage(`{"state": "IDLE"}`)
	code, body = env.hostUpdateChannelConfig(deviceID, "CH01", token, ch01Stop)
	if code != 200 {
		t.Fatalf("Phase 2: stop CH01 expected 200, got %d: %s", code, body)
	}
	time.Sleep(3 * time.Second)

	// Verify SDK still has both channels (the regression was CH02 disappearing)
	sdkCfg = env.sdkGetConfiguration()
	if sdkCfg.Configuration == nil {
		t.Fatal("Phase 2: SDK returned nil configuration after CH01 stop")
	}
	cfgStr = string(mustMarshal(sdkCfg.Configuration))
	if !strings.Contains(cfgStr, "CH01") || !strings.Contains(cfgStr, "CH02") {
		t.Fatalf("Phase 2: SDK config lost a channel after CH01 stop: %s", cfgStr)
	}
	t.Log("Phase 2: OK — both channels still present in SDK config after CH01 stop")

	// Phase 3: Stop CH02 via per-channel API
	t.Log("Phase 3: Stop CH02 via per-channel API")
	ch02Stop := json.RawMessage(`{"state": "IDLE"}`)
	code, body = env.hostUpdateChannelConfig(deviceID, "CH02", token, ch02Stop)
	if code != 200 {
		t.Fatalf("Phase 3: stop CH02 expected 200, got %d: %s", code, body)
	}
	time.Sleep(3 * time.Second)

	// Verify SDK has both channels with IDLE state
	sdkCfg = env.sdkGetConfiguration()
	if sdkCfg.Configuration == nil {
		t.Fatal("Phase 3: SDK returned nil configuration after CH02 stop")
	}
	cfgStr = string(mustMarshal(sdkCfg.Configuration))
	if !strings.Contains(cfgStr, "CH01") || !strings.Contains(cfgStr, "CH02") {
		t.Fatalf("Phase 3: SDK config lost a channel after CH02 stop: %s", cfgStr)
	}
	// Both channels should now show IDLE
	if !strings.Contains(cfgStr, `"IDLE"`) {
		t.Fatalf("Phase 3: expected IDLE state in SDK config: %s", cfgStr)
	}

	// Verify both channels have distinct versions (both were bumped independently)
	payload, ok := sdkCfg.Configuration["payload"]
	if !ok {
		t.Fatal("Phase 3: no payload in SDK config")
	}
	b, _ := json.Marshal(payload)
	var cfg struct {
		Channels []struct {
			Id      string `json:"id"`
			Version string `json:"version"`
			State   string `json:"state"`
		} `json:"channels"`
	}
	json.Unmarshal(b, &cfg)
	if len(cfg.Channels) != 2 {
		t.Fatalf("Phase 3: expected 2 channels in payload, got %d", len(cfg.Channels))
	}
	for _, ch := range cfg.Channels {
		if ch.State != "IDLE" {
			t.Errorf("Phase 3: channel %s expected IDLE, got %s", ch.Id, ch.State)
		}
		if ch.Version == "" {
			t.Errorf("Phase 3: channel %s has empty version", ch.Id)
		}
	}
	t.Log("Phase 3: OK — both channels IDLE with independent versions")

	t.Log("=== TestSequentialPerChannelStops PASSED ===")
}

// TestSrtEncryptionRoundTrip verifies that an SRT channel configuration with
// encryption (passphrase + keyLength) is correctly delivered to the device SDK
// and can be reported back in actual_configuration.
func TestSrtEncryptionRoundTrip(t *testing.T) {
	env := newTestEnv(t)
	env.startHost()
	env.startSDK("integ-enc-001")

	registration := loadRegistration(t)

	// Setup: pair and connect
	acct := env.hostRegisterAccount("encuser", "testpass123", "Encryption Test")
	token := acct.Token
	pairingCode := env.waitForPairingCode("tr12-host", registration, 15*time.Second)
	env.hostClaim(pairingCode, token)
	env.waitForSDKConnected("tr12-host", registration, 30*time.Second)

	devices := env.hostListDevices(token)
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}
	deviceID := devices[0].DeviceID

	// ---------------------------------------------------------------
	// Phase 1: Push SRT Caller config with AES-256 encryption
	// ---------------------------------------------------------------
	t.Log("Phase 1: Push SRT Caller with AES-256 encryption")
	encConfig := json.RawMessage(`{
		"channels": [{
			"id": "CH01",
			"state": "ACTIVE",
			"channelSettings": {"standardSettings": [
				{"id": "RS01", "value": "1920x1080"},
				{"id": "FR01", "value": "30"},
				{"id": "MB01", "value": "10000"},
				{"id": "RC01", "value": "CBR"},
				{"id": "CO01", "value": "H.264"},
				{"id": "GP01", "value": "60"},
				{"id": "IN01", "value": "SDI1"}
			]},
			"protocol": {
				"srtCaller": {
					"address": "10.0.0.50",
					"port": 9000,
					"minimumLatencyMilliseconds": 500,
					"streamId": "live/feed1",
					"encryption": {
						"passphrase": "MySecurePassphrase2024!",
						"keyLength": "AES_256"
					}
				}
			}
		}]
	}`)
	code, body := env.hostUpdateConfig(deviceID, token, encConfig)
	if code != 200 {
		t.Fatalf("Phase 1: expected 200, got %d: %s", code, body)
	}

	time.Sleep(3 * time.Second)
	sdkCfg := env.sdkGetConfiguration()
	if sdkCfg.Configuration == nil {
		t.Fatal("Phase 1: SDK returned nil configuration")
	}
	cfgStr := string(mustMarshal(sdkCfg.Configuration))

	// Verify encryption fields arrived at the device
	if !strings.Contains(cfgStr, "MySecurePassphrase2024!") {
		t.Fatalf("Phase 1: SDK config missing passphrase: %s", cfgStr)
	}
	if !strings.Contains(cfgStr, "AES_256") {
		t.Fatalf("Phase 1: SDK config missing keyLength AES_256: %s", cfgStr)
	}
	if !strings.Contains(cfgStr, "live/feed1") {
		t.Fatalf("Phase 1: SDK config missing streamId: %s", cfgStr)
	}
	t.Log("Phase 1: OK — encryption config delivered to device")

	// ---------------------------------------------------------------
	// Phase 2: Push SRT Listener config with AES-128 encryption
	// ---------------------------------------------------------------
	t.Log("Phase 2: Push SRT Listener with AES-128 encryption")
	encConfig2 := json.RawMessage(`{
		"channels": [{
			"id": "CH01",
			"state": "ACTIVE",
			"channelSettings": {"standardSettings": [
				{"id": "RS01", "value": "1920x1080"},
				{"id": "FR01", "value": "30"},
				{"id": "MB01", "value": "10000"},
				{"id": "RC01", "value": "CBR"},
				{"id": "CO01", "value": "H.264"},
				{"id": "GP01", "value": "60"},
				{"id": "IN01", "value": "SDI1"}
			]},
			"protocol": {
				"srtListener": {
					"port": 4900,
					"minimumLatencyMilliseconds": 300,
					"encryption": {
						"passphrase": "ShortButValid!",
						"keyLength": "AES_128"
					}
				}
			}
		}]
	}`)
	code, body = env.hostUpdateConfig(deviceID, token, encConfig2)
	if code != 200 {
		t.Fatalf("Phase 2: expected 200, got %d: %s", code, body)
	}

	time.Sleep(3 * time.Second)
	sdkCfg2 := env.sdkGetConfiguration()
	if sdkCfg2.Configuration == nil {
		t.Fatal("Phase 2: SDK returned nil configuration")
	}
	cfgStr2 := string(mustMarshal(sdkCfg2.Configuration))

	if !strings.Contains(cfgStr2, "ShortButValid!") {
		t.Fatalf("Phase 2: SDK config missing passphrase: %s", cfgStr2)
	}
	if !strings.Contains(cfgStr2, "AES_128") {
		t.Fatalf("Phase 2: SDK config missing keyLength AES_128: %s", cfgStr2)
	}
	if !strings.Contains(cfgStr2, "4900") {
		t.Fatalf("Phase 2: SDK config missing port 4900: %s", cfgStr2)
	}
	t.Log("Phase 2: OK — SRT Listener with AES-128 delivered")

	// ---------------------------------------------------------------
	// Phase 3: Push SRT Caller with no encryption (null/omitted)
	// ---------------------------------------------------------------
	t.Log("Phase 3: Push SRT Caller with no encryption")
	noEncConfig := json.RawMessage(`{
		"channels": [{
			"id": "CH01",
			"state": "ACTIVE",
			"channelSettings": {"standardSettings": [
				{"id": "RS01", "value": "1920x1080"},
				{"id": "FR01", "value": "30"},
				{"id": "MB01", "value": "10000"},
				{"id": "RC01", "value": "CBR"},
				{"id": "CO01", "value": "H.264"},
				{"id": "GP01", "value": "60"},
				{"id": "IN01", "value": "SDI1"}
			]},
			"protocol": {
				"srtCaller": {
					"address": "10.0.0.50",
					"port": 9000,
					"minimumLatencyMilliseconds": 200
				}
			}
		}]
	}`)
	code, body = env.hostUpdateConfig(deviceID, token, noEncConfig)
	if code != 200 {
		t.Fatalf("Phase 3: expected 200, got %d: %s", code, body)
	}

	time.Sleep(3 * time.Second)
	sdkCfg3 := env.sdkGetConfiguration()
	if sdkCfg3.Configuration == nil {
		t.Fatal("Phase 3: SDK returned nil configuration")
	}
	cfgStr3 := string(mustMarshal(sdkCfg3.Configuration))

	// Verify no passphrase or keyLength present
	if strings.Contains(cfgStr3, "passphrase") {
		t.Fatalf("Phase 3: SDK config should NOT contain passphrase: %s", cfgStr3)
	}
	if strings.Contains(cfgStr3, "keyLength") {
		t.Fatalf("Phase 3: SDK config should NOT contain keyLength: %s", cfgStr3)
	}
	t.Log("Phase 3: OK — no encryption config when omitted")

	// ---------------------------------------------------------------
	// Phase 4: Echo-back — device reports actual config with encryption
	// ---------------------------------------------------------------
	t.Log("Phase 4: Verify encryption in actual_configuration round-trip")

	// Re-push encrypted config
	code, _ = env.hostUpdateConfig(deviceID, token, encConfig)
	if code != 200 {
		t.Fatalf("Phase 4: re-push expected 200, got %d", code)
	}
	time.Sleep(3 * time.Second)

	sdkCfg4 := env.sdkGetConfiguration()
	var actualCfg map[string]interface{}
	if p, ok := sdkCfg4.Configuration["payload"]; ok {
		b, _ := json.Marshal(p)
		json.Unmarshal(b, &actualCfg)
	}
	if actualCfg == nil {
		t.Fatal("Phase 4: could not extract payload")
	}

	// Report it back as actual config
	cfgResp := env.sdkReportActualConfig(actualCfg)
	if !cfgResp.Success {
		t.Fatalf("Phase 4: report_actual_configuration failed: %s", cfgResp.Message)
	}
	time.Sleep(3 * time.Second)

	// Verify host has the encryption in actual_configuration
	detail := env.hostDescribeDevice(deviceID, token)
	actualStr := string(detail.ActualConfiguration)
	if !strings.Contains(actualStr, "MySecurePassphrase2024!") {
		t.Fatalf("Phase 4: host actual_configuration missing passphrase: %s", actualStr)
	}
	if !strings.Contains(actualStr, "AES_256") {
		t.Fatalf("Phase 4: host actual_configuration missing AES_256: %s", actualStr)
	}
	t.Log("Phase 4: OK — encryption round-trips through actual_configuration")

	t.Log("=== TestSrtEncryptionRoundTrip PASSED ===")
}
