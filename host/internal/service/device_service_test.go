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
package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"strings"
	"testing"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/ca"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/config"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/db"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/models"
)

// mockMQTT captures published messages for assertions.
type mockMQTT struct {
	messages []mqttMsg
}

type mqttMsg struct {
	Topic   string
	Payload []byte
	Retain  bool
}

func (m *mockMQTT) Publish(topic string, payload []byte, retain bool) error {
	m.messages = append(m.messages, mqttMsg{Topic: topic, Payload: payload, Retain: retain})
	return nil
}

func newTestDeviceService(t *testing.T) (*DeviceService, *mockMQTT, *db.Store) {
	t.Helper()
	store, err := db.New(":memory:")
	if err != nil {
		t.Fatalf("db.New: %v", err)
	}
	t.Cleanup(func() { store.Close() })

	caInst, err := ca.New(store, "127.0.0.1")
	if err != nil {
		t.Fatalf("ca.New: %v", err)
	}

	mqtt := &mockMQTT{}
	cfg := &config.Config{
		ServiceID:      "test-host",
		ServiceName:    "Test Host",
		HostAddress:    "127.0.0.1",
		MQTTPort:       8883,
		CertExpiryDays: 30,
		PairingTimeout: 1800,
	}
	svc := NewDeviceService(store, caInst, mqtt, cfg)
	return svc, mqtt, store
}

func generateTestCSR(t *testing.T) string {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}
	template := &x509.CertificateRequest{
		Subject: pkix.Name{Organization: []string{"Test"}},
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, template, key)
	if err != nil {
		t.Fatalf("create CSR: %v", err)
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER}))
}

// doPair is a helper that pairs a device and returns the device ID, pairing code, and access code.
func doPair(t *testing.T, svc *DeviceService) (deviceID, pairingCode, accessCode string) {
	t.Helper()
	csr := generateTestCSR(t)
	resp, err := svc.Pair(models.PairRequestContent{
		HostId:     "test-host",
		Version:    "1.0",
		DeviceType: "SOURCE",
		Csr:        csr,
	})
	if err != nil {
		t.Fatalf("Pair: %v", err)
	}
	result := resp.GetResult()
	if result.Success == nil {
		t.Fatal("expected Success result from Pair")
	}
	data := result.Success.Success
	return data.DeviceId, data.PairingCode, data.AccessCode
}

// --- Pair ---

func TestPair_Success(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	deviceID, pairingCode, accessCode := doPair(t, svc)

	if len(deviceID) != 21 {
		t.Fatalf("expected 21-char device ID, got %d: %q", len(deviceID), deviceID)
	}
	if len(pairingCode) != 6 {
		t.Fatalf("expected 6-char pairing code, got %d: %q", len(pairingCode), pairingCode)
	}
	if len(accessCode) != 32 {
		t.Fatalf("expected 32-char access code, got %d: %q", len(accessCode), accessCode)
	}
}

func TestPair_HostIDMismatch(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	resp, err := svc.Pair(models.PairRequestContent{
		HostId:     "wrong-host",
		Version:    "1.0",
		DeviceType: "SOURCE",
		Csr:        generateTestCSR(t),
	})
	if err != nil {
		t.Fatalf("Pair: %v", err)
	}
	result := resp.GetResult()
	if result.Failure == nil {
		t.Fatal("expected Failure result for host ID mismatch")
	}
}

func TestPair_BadDeviceType(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	resp, _ := svc.Pair(models.PairRequestContent{
		HostId:     "test-host",
		Version:    "1.0",
		DeviceType: "INVALID",
		Csr:        generateTestCSR(t),
	})
	result := resp.GetResult()
	if result.Failure == nil {
		t.Fatal("expected Failure for bad device type")
	}
}

func TestPair_EmptyVersion(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	resp, _ := svc.Pair(models.PairRequestContent{
		HostId:     "test-host",
		Version:    "",
		DeviceType: "SOURCE",
		Csr:        generateTestCSR(t),
	})
	result := resp.GetResult()
	if result.Failure == nil {
		t.Fatal("expected Failure for empty version")
	}
}

// --- Authenticate ---

func TestAuthenticate_Standby(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	deviceID, pairingCode, accessCode := doPair(t, svc)

	resp, err := svc.Authenticate(models.AuthenticateRequestContent{
		DeviceId:    deviceID,
		PairingCode: pairingCode,
		AccessCode:  accessCode,
	})
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if resp.GetStatus() != models.AuthStatusSTANDBY {
		t.Fatalf("expected STANDBY, got %v", resp.GetStatus())
	}
}

func TestAuthenticate_Claimed(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	deviceID, pairingCode, accessCode := doPair(t, svc)

	// Claim the device
	if err := svc.Claim(pairingCode, "acc-1", 730); err != nil {
		t.Fatalf("Claim: %v", err)
	}

	resp, err := svc.Authenticate(models.AuthenticateRequestContent{
		DeviceId:    deviceID,
		PairingCode: pairingCode,
		AccessCode:  accessCode,
	})
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if resp.GetStatus() != models.AuthStatusCLAIMED {
		t.Fatalf("expected CLAIMED, got %v", resp.GetStatus())
	}
	if resp.GetCaCert() == "" {
		t.Fatal("expected CA cert in claimed response")
	}
	if resp.GetMqttUri() == "" {
		t.Fatal("expected MQTT URI in claimed response")
	}
}

func TestAuthenticate_WrongCredentials(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	deviceID, _, _ := doPair(t, svc)

	_, err := svc.Authenticate(models.AuthenticateRequestContent{
		DeviceId:    deviceID,
		PairingCode: "WRONG1",
		AccessCode:  "wrong",
	})
	if err != ErrUnauthorized {
		t.Fatalf("expected ErrUnauthorized, got %v", err)
	}
}

func TestAuthenticate_NotFound(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	_, err := svc.Authenticate(models.AuthenticateRequestContent{
		DeviceId:    "nonexistent",
		PairingCode: "ABC123",
		AccessCode:  "secret",
	})
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

// --- Claim ---

func TestClaim_Success(t *testing.T) {
	svc, _, store := newTestDeviceService(t)
	_, pairingCode, _ := doPair(t, svc)

	if err := svc.Claim(pairingCode, "acc-1", 730); err != nil {
		t.Fatalf("Claim: %v", err)
	}

	// Verify device is now ACTIVE
	d, _ := store.GetDeviceByPairingCode(pairingCode)
	if d.State != "ACTIVE" {
		t.Fatalf("expected ACTIVE, got %q", d.State)
	}
	if d.AccountID != "acc-1" {
		t.Fatalf("expected acc-1, got %q", d.AccountID)
	}
}

func TestClaim_NotFound(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	err := svc.Claim("NOPE00", "acc-1", 730)
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestClaim_AlreadyClaimed(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	_, pairingCode, _ := doPair(t, svc)
	svc.Claim(pairingCode, "acc-1", 730)

	err := svc.Claim(pairingCode, "acc-2", 730)
	if err != ErrConflict {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

// --- ListDevices ---

func TestListDevices(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	_, pc1, _ := doPair(t, svc)
	_, pc2, _ := doPair(t, svc)
	svc.Claim(pc1, "acc-1", 730)
	svc.Claim(pc2, "acc-1", 730)

	summaries, err := svc.ListDevices("acc-1")
	if err != nil {
		t.Fatalf("ListDevices: %v", err)
	}
	if len(summaries) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(summaries))
	}
}

// --- DescribeDevice ---

func TestDescribeDevice_Success(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	deviceID, pc, _ := doPair(t, svc)
	svc.Claim(pc, "acc-1", 730)

	detail, err := svc.DescribeDevice(deviceID, "acc-1")
	if err != nil {
		t.Fatalf("DescribeDevice: %v", err)
	}
	if detail.DeviceID != deviceID {
		t.Fatalf("expected %q, got %q", deviceID, detail.DeviceID)
	}
}

func TestDescribeDevice_NotFound(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	_, err := svc.DescribeDevice("nonexistent", "acc-1")
	if err != ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDescribeDevice_Forbidden(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	deviceID, pc, _ := doPair(t, svc)
	svc.Claim(pc, "acc-1", 730)

	_, err := svc.DescribeDevice(deviceID, "acc-other")
	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

// --- UpdateConfiguration ---

func TestUpdateConfiguration_Success(t *testing.T) {
	svc, mqtt, store := newTestDeviceService(t)
	deviceID, pc, _ := doPair(t, svc)
	svc.Claim(pc, "acc-1", 730)

	// Set registration so validation passes
	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Channel 1","simpleSettings":[{"id":"brightness"}],"profiles":[{"id":"prof1"}]}]}`)
	store.UpdateDeviceRegistration(deviceID, reg)

	cfg := json.RawMessage(`{"channels":[{"id":"ch1","state":"ACTIVE"}]}`)
	if err := svc.UpdateConfiguration(deviceID, "acc-1", cfg); err != nil {
		t.Fatalf("UpdateConfiguration: %v", err)
	}

	// Verify MQTT message was published
	if len(mqtt.messages) == 0 {
		t.Fatal("expected MQTT message")
	}
	expectedTopic := "cdd/" + deviceID + "/config/update"
	if mqtt.messages[0].Topic != expectedTopic {
		t.Fatalf("expected topic %q, got %q", expectedTopic, mqtt.messages[0].Topic)
	}
}

func TestUpdateConfiguration_ValidationError(t *testing.T) {
	svc, _, store := newTestDeviceService(t)
	deviceID, pc, _ := doPair(t, svc)
	svc.Claim(pc, "acc-1", 730)

	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Channel 1"}]}`)
	store.UpdateDeviceRegistration(deviceID, reg)

	// Unknown channel ID
	cfg := json.RawMessage(`{"channels":[{"id":"ch-unknown","state":"ACTIVE"}]}`)
	err := svc.UpdateConfiguration(deviceID, "acc-1", cfg)
	if err == nil || !strings.Contains(err.Error(), "unknown channel ID") {
		t.Fatalf("expected validation error about unknown channel, got %v", err)
	}
}

// --- Deprovision ---

func TestDeprovision(t *testing.T) {
	svc, mqtt, store := newTestDeviceService(t)
	deviceID, pc, _ := doPair(t, svc)
	svc.Claim(pc, "acc-1", 730)

	if err := svc.Deprovision(deviceID, "acc-1"); err != nil {
		t.Fatalf("Deprovision: %v", err)
	}

	d, _ := store.GetDevice(deviceID)
	if d.State != "DEPROVISIONED" {
		t.Fatalf("expected DEPROVISIONED, got %q", d.State)
	}

	// Verify MQTT deprovision message
	found := false
	for _, msg := range mqtt.messages {
		if strings.Contains(msg.Topic, "deprovision") {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected deprovision MQTT message")
	}

	// Idempotent
	if err := svc.Deprovision(deviceID, "acc-1"); err != nil {
		t.Fatalf("second Deprovision should be idempotent: %v", err)
	}
}

func TestDeprovision_Forbidden(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	deviceID, pc, _ := doPair(t, svc)
	svc.Claim(pc, "acc-1", 730)

	err := svc.Deprovision(deviceID, "acc-other")
	if err != ErrForbidden {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

// --- FullCleanup ---

func TestFullCleanup(t *testing.T) {
	svc, _, store := newTestDeviceService(t)
	deviceID, pc, _ := doPair(t, svc)
	svc.Claim(pc, "acc-1", 730)

	// Add thumbnail and log
	store.UpsertThumbnail(&db.Thumbnail{DeviceID: deviceID, SourceID: "ch1", ImageData: []byte{1}, Timestamp: "now", ImageType: "jpeg", ImageSizeKB: 1})
	store.UpsertLog(&db.DeviceLog{DeviceID: deviceID, LogData: []byte("log"), UploadedAt: "now", LogSizeKB: 1})

	if err := svc.FullCleanup(deviceID); err != nil {
		t.Fatalf("FullCleanup: %v", err)
	}

	d, _ := store.GetDevice(deviceID)
	if d != nil {
		t.Fatal("expected device to be deleted")
	}
	thumb, _ := store.GetThumbnail(deviceID, "ch1")
	if thumb != nil {
		t.Fatal("expected thumbnail to be deleted")
	}
	log, _ := store.GetLog(deviceID)
	if log != nil {
		t.Fatal("expected log to be deleted")
	}
}

// --- RotateCredentials ---

func TestRotateCredentials(t *testing.T) {
	svc, mqtt, store := newTestDeviceService(t)
	deviceID, pc, _ := doPair(t, svc)
	svc.Claim(pc, "acc-1", 730)

	if err := svc.RotateCredentials(deviceID, "acc-1"); err != nil {
		t.Fatalf("RotateCredentials: %v", err)
	}

	d, _ := store.GetDevice(deviceID)
	if d.CurrentCertPEM == "" {
		t.Fatal("expected current cert to be set")
	}
	if d.PreviousCertPEM == "" {
		t.Fatal("expected previous cert to be set")
	}
	if d.LastRotationAt == "" {
		t.Fatal("expected last rotation to be set")
	}

	// Verify MQTT message
	found := false
	for _, msg := range mqtt.messages {
		if strings.Contains(msg.Topic, "certs/update") {
			found = true
			if !msg.Retain {
				t.Fatal("expected cert rotation message to be retained")
			}
			break
		}
	}
	if !found {
		t.Fatal("expected cert rotation MQTT message")
	}
}

func TestRotateCredentials_NotActive(t *testing.T) {
	svc, _, _ := newTestDeviceService(t)
	deviceID, _, _ := doPair(t, svc) // still in PAIRING state

	err := svc.RotateCredentials(deviceID, "")
	if err == nil || !strings.Contains(err.Error(), "not active") {
		t.Fatalf("expected 'not active' error, got %v", err)
	}
}

// --- validateConfiguration ---

func TestValidateConfiguration_ValidConfig(t *testing.T) {
	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Ch1","simpleSettings":[{"id":"brightness"}],"profiles":[{"id":"prof1"}]}]}`)
	cfg := json.RawMessage(`{"channels":[{"id":"ch1","state":"ACTIVE","settings":{"simpleSettings":[{"key":"brightness"}],"profile":{"id":"prof1"}}}]}`)
	if err := validateConfiguration(cfg, reg); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}
}

func TestValidateConfiguration_UnknownChannel(t *testing.T) {
	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Ch1"}]}`)
	cfg := json.RawMessage(`{"channels":[{"id":"ch-bad"}]}`)
	err := validateConfiguration(cfg, reg)
	if err == nil || !strings.Contains(err.Error(), "unknown channel ID") {
		t.Fatalf("expected unknown channel error, got %v", err)
	}
}

func TestValidateConfiguration_UnknownSettingKey(t *testing.T) {
	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Ch1","simpleSettings":[{"id":"brightness"}]}]}`)
	cfg := json.RawMessage(`{"channels":[{"id":"ch1","settings":{"simpleSettings":[{"key":"contrast"}]}}]}`)
	err := validateConfiguration(cfg, reg)
	if err == nil || !strings.Contains(err.Error(), "unknown setting key") {
		t.Fatalf("expected unknown setting key error, got %v", err)
	}
}

func TestValidateConfiguration_UnknownProfile(t *testing.T) {
	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Ch1","profiles":[{"id":"prof1"}]}]}`)
	cfg := json.RawMessage(`{"channels":[{"id":"ch1","settings":{"profile":{"id":"prof-bad"}}}]}`)
	err := validateConfiguration(cfg, reg)
	if err == nil || !strings.Contains(err.Error(), "unknown profile ID") {
		t.Fatalf("expected unknown profile error, got %v", err)
	}
}

func TestValidateConfiguration_InvalidChannelState(t *testing.T) {
	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Ch1"}]}`)
	cfg := json.RawMessage(`{"channels":[{"id":"ch1","state":"BROKEN"}]}`)
	err := validateConfiguration(cfg, reg)
	if err == nil || !strings.Contains(err.Error(), "invalid channel state") {
		t.Fatalf("expected invalid channel state error, got %v", err)
	}
}

func TestValidateConfiguration_DeviceLevelSettingValid(t *testing.T) {
	reg := json.RawMessage(`{"simpleSettings":[{"id":"clock_source"}],"channels":[{"id":"ch1","name":"Ch1"}]}`)
	cfg := json.RawMessage(`{"simpleSettings":[{"key":"clock_source","value":"PTP"}],"channels":[{"id":"ch1","state":"ACTIVE"}]}`)
	if err := validateConfiguration(cfg, reg); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}
}

func TestValidateConfiguration_DeviceLevelSettingUnknown(t *testing.T) {
	reg := json.RawMessage(`{"simpleSettings":[{"id":"clock_source"}],"channels":[{"id":"ch1","name":"Ch1"}]}`)
	cfg := json.RawMessage(`{"simpleSettings":[{"key":"bogus_setting","value":"foo"}],"channels":[{"id":"ch1","state":"ACTIVE"}]}`)
	err := validateConfiguration(cfg, reg)
	if err == nil || !strings.Contains(err.Error(), "unknown device-level setting key") {
		t.Fatalf("expected unknown device-level setting key error, got %v", err)
	}
	if !strings.Contains(err.Error(), "bogus_setting") {
		t.Fatalf("expected error to mention bogus_setting, got %v", err)
	}
}

func TestValidateConfiguration_ConnectionValid(t *testing.T) {
	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Ch1","connectionProtocols":["SRT_CALLER"]}]}`)
	cfg := json.RawMessage(`{"channels":[{"id":"ch1","state":"ACTIVE","connection":{"transportProtocol":{"srtCaller":{"ip":"1.2.3.4","port":9000,"minimumLatencyMilliseconds":200}}}}]}`)
	if err := validateConfiguration(cfg, reg); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}
}

func TestValidateConfiguration_ConnectionUnknownProtocolType(t *testing.T) {
	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Ch1","connectionProtocols":["SRT_CALLER"]}]}`)
	cfg := json.RawMessage(`{"channels":[{"id":"ch1","state":"ACTIVE","connection":{"transportProtocol":{"rtmpCaller":{"ip":"1.2.3.4"}}}}]}`)
	err := validateConfiguration(cfg, reg)
	if err == nil || !strings.Contains(err.Error(), "unknown transport protocol type") {
		t.Fatalf("expected unknown transport protocol type error, got %v", err)
	}
}

func TestValidateConfiguration_ConnectionProtocolNotRegistered(t *testing.T) {
	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Ch1","connectionProtocols":["SRT_LISTENER"]}]}`)
	cfg := json.RawMessage(`{"channels":[{"id":"ch1","state":"ACTIVE","connection":{"transportProtocol":{"srtCaller":{"ip":"1.2.3.4","port":9000,"minimumLatencyMilliseconds":200}}}}]}`)
	err := validateConfiguration(cfg, reg)
	if err == nil || !strings.Contains(err.Error(), "not supported by channel") {
		t.Fatalf("expected protocol not supported error, got %v", err)
	}
}

func TestValidateConfiguration_ConnectionUnknownField(t *testing.T) {
	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Ch1","connectionProtocols":["SRT_CALLER"]}]}`)
	// "address" and "ports" are wrong — should be "ip" and "port"
	cfg := json.RawMessage(`{"channels":[{"id":"ch1","state":"ACTIVE","connection":{"transportProtocol":{"srtCaller":{"address":"1.2.3.4","ports":9000,"minimumLatencyMilliseconds":200}}}}]}`)
	err := validateConfiguration(cfg, reg)
	if err == nil || !strings.Contains(err.Error(), "unknown field") {
		t.Fatalf("expected unknown field error, got %v", err)
	}
	if !strings.Contains(err.Error(), "address") {
		t.Fatalf("expected error to mention 'address', got %v", err)
	}
}

func TestValidateConfiguration_ConnectionMissingRequiredField(t *testing.T) {
	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Ch1","connectionProtocols":["SRT_CALLER"]}]}`)
	// Missing "ip" which is required for srtCaller
	cfg := json.RawMessage(`{"channels":[{"id":"ch1","state":"ACTIVE","connection":{"transportProtocol":{"srtCaller":{"port":9000,"minimumLatencyMilliseconds":200}}}}]}`)
	err := validateConfiguration(cfg, reg)
	if err == nil || !strings.Contains(err.Error(), "missing required field") {
		t.Fatalf("expected missing required field error, got %v", err)
	}
	if !strings.Contains(err.Error(), `"ip"`) {
		t.Fatalf("expected error to mention 'ip', got %v", err)
	}
}

func TestValidateConfiguration_ConnectionSrtListenerValid(t *testing.T) {
	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Ch1","connectionProtocols":["SRT_LISTENER"]}]}`)
	cfg := json.RawMessage(`{"channels":[{"id":"ch1","state":"ACTIVE","connection":{"transportProtocol":{"srtListener":{"port":9000,"minimumLatencyMilliseconds":200}}}}]}`)
	if err := validateConfiguration(cfg, reg); err != nil {
		t.Fatalf("expected valid, got %v", err)
	}
}

func TestValidateConfiguration_ConnectionWithOptionalFields(t *testing.T) {
	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Ch1","connectionProtocols":["SRT_CALLER"]}]}`)
	cfg := json.RawMessage(`{"channels":[{"id":"ch1","state":"ACTIVE","connection":{"transportProtocol":{"srtCaller":{"ip":"1.2.3.4","port":9000,"minimumLatencyMilliseconds":200,"streamId":"test123"}}}}]}`)
	if err := validateConfiguration(cfg, reg); err != nil {
		t.Fatalf("expected valid with optional streamId, got %v", err)
	}
}

// --- Helper function tests ---

func TestGenerateDeviceID_Format(t *testing.T) {
	id := generateDeviceID()
	if len(id) != 21 {
		t.Fatalf("expected 21 chars, got %d: %q", len(id), id)
	}
}

func TestGeneratePairingCode_Format(t *testing.T) {
	code := generatePairingCode()
	if len(code) != 6 {
		t.Fatalf("expected 6 chars, got %d: %q", len(code), code)
	}
	if strings.ToUpper(code) != code {
		t.Fatalf("expected uppercase, got %q", code)
	}
}

func TestFormatOnlineDetails(t *testing.T) {
	tests := []struct {
		name     string
		device   *models.Device
		contains string
	}{
		{"offline no last seen", &models.Device{Online: false}, "offline"},
		{"offline with last seen", &models.Device{Online: false, LastSeen: "2025-01-01T00:00:00Z"}, "offline since"},
		{"online no last seen", &models.Device{Online: true}, "online"},
		{"online with last seen", &models.Device{Online: true, LastSeen: time.Now().UTC().Format(time.RFC3339)}, "online:"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatOnlineDetails(tt.device)
			if !strings.Contains(result, tt.contains) {
				t.Fatalf("expected %q to contain %q", result, tt.contains)
			}
		})
	}
}

func TestFormatCertExpiration(t *testing.T) {
	if formatCertExpiration("") != "unknown" {
		t.Fatal("expected unknown for empty")
	}
	if formatCertExpiration("not-a-date") != "unknown" {
		t.Fatal("expected unknown for invalid date")
	}
	past := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	if formatCertExpiration(past) != "expired" {
		t.Fatalf("expected expired, got %q", formatCertExpiration(past))
	}
	future := time.Now().Add(48 * time.Hour).UTC().Format(time.RFC3339)
	result := formatCertExpiration(future)
	if strings.Contains(result, "unknown") || strings.Contains(result, "expired") {
		t.Fatalf("expected valid duration, got %q", result)
	}
}
