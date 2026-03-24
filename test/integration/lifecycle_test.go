//go:build integration

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

	// Report actual config matching the desired config from 7d
	var actualCfg map[string]interface{}
	json.Unmarshal(fullConfig, &actualCfg)
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
