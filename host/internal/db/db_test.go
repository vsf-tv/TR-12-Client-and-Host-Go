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
package db

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/models"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := New(":memory:")
	if err != nil {
		t.Fatalf("New(:memory:): %v", err)
	}
	t.Cleanup(func() { s.Close() })
	return s
}

func makeDevice(id, accountID, deviceType, state string) *models.Device {
	return &models.Device{
		DeviceID:   id,
		AccountID:  accountID,
		DeviceType: deviceType,
		State:      state,
		PairedAt:   time.Now().UTC().Format(time.RFC3339),
	}
}

// --- Device CRUD ---

func TestInsertAndGetDevice(t *testing.T) {
	s := newTestStore(t)
	d := makeDevice("dev-1", "acc-1", "SOURCE", "PAIRING")
	d.PairingCode = "ABC123"
	d.AccessCode = "secret"
	d.PairingExpiresAt = time.Now().Add(30 * time.Minute).UTC().Format(time.RFC3339)

	if err := s.InsertDevice(d); err != nil {
		t.Fatalf("InsertDevice: %v", err)
	}

	got, err := s.GetDevice("dev-1")
	if err != nil {
		t.Fatalf("GetDevice: %v", err)
	}
	if got == nil {
		t.Fatal("expected device, got nil")
	}
	if got.DeviceID != "dev-1" || got.DeviceType != "SOURCE" || got.State != "PAIRING" {
		t.Fatalf("unexpected device: %+v", got)
	}
	if got.PairingCode != "ABC123" {
		t.Fatalf("expected pairing code ABC123, got %q", got.PairingCode)
	}
}

func TestGetDevice_NotFound(t *testing.T) {
	s := newTestStore(t)
	d, err := s.GetDevice("nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d != nil {
		t.Fatal("expected nil for nonexistent device")
	}
}

func TestGetDeviceByPairingCode(t *testing.T) {
	s := newTestStore(t)
	d := makeDevice("dev-2", "", "SOURCE", "PAIRING")
	d.PairingCode = "XYZ789"
	s.InsertDevice(d)

	got, err := s.GetDeviceByPairingCode("XYZ789")
	if err != nil {
		t.Fatalf("GetDeviceByPairingCode: %v", err)
	}
	if got == nil || got.DeviceID != "dev-2" {
		t.Fatalf("expected dev-2, got %+v", got)
	}

	// Not found
	got, _ = s.GetDeviceByPairingCode("NOPE")
	if got != nil {
		t.Fatal("expected nil for unknown pairing code")
	}
}

func TestListDevicesByAccount(t *testing.T) {
	s := newTestStore(t)
	s.InsertDevice(makeDevice("d1", "acc-1", "SOURCE", "ACTIVE"))
	s.InsertDevice(makeDevice("d2", "acc-1", "DESTINATION", "ACTIVE"))
	s.InsertDevice(makeDevice("d3", "acc-2", "SOURCE", "ACTIVE"))

	devices, err := s.ListDevicesByAccount("acc-1")
	if err != nil {
		t.Fatalf("ListDevicesByAccount: %v", err)
	}
	if len(devices) != 2 {
		t.Fatalf("expected 2 devices, got %d", len(devices))
	}

	devices, _ = s.ListDevicesByAccount("acc-2")
	if len(devices) != 1 {
		t.Fatalf("expected 1 device, got %d", len(devices))
	}

	devices, _ = s.ListDevicesByAccount("acc-none")
	if len(devices) != 0 {
		t.Fatalf("expected 0 devices, got %d", len(devices))
	}
}


// --- Device State Updates ---

func TestUpdateDeviceState(t *testing.T) {
	s := newTestStore(t)
	d := makeDevice("d1", "acc-1", "SOURCE", "PAIRING")
	d.DesiredConfig = json.RawMessage(`{"channels":[]}`)
	s.InsertDevice(d)

	// Without clearing data
	s.UpdateDeviceState("d1", "ACTIVE", false)
	got, _ := s.GetDevice("d1")
	if got.State != "ACTIVE" {
		t.Fatalf("expected ACTIVE, got %q", got.State)
	}
	if got.DesiredConfig == nil {
		t.Fatal("expected DesiredConfig to be preserved")
	}

	// With clearing data
	s.UpdateDeviceState("d1", "DEPROVISIONED", true)
	got, _ = s.GetDevice("d1")
	if got.State != "DEPROVISIONED" {
		t.Fatalf("expected DEPROVISIONED, got %q", got.State)
	}
	if got.DesiredConfig != nil {
		t.Fatal("expected DesiredConfig to be cleared")
	}
}

func TestClaimDevice(t *testing.T) {
	s := newTestStore(t)
	d := makeDevice("d1", "", "SOURCE", "PAIRING")
	s.InsertDevice(d)

	regExpires := time.Now().Add(730 * 24 * time.Hour).UTC().Format(time.RFC3339)
	if err := s.ClaimDevice("d1", "acc-1", regExpires, "", "", 365); err != nil {
		t.Fatalf("ClaimDevice: %v", err)
	}

	got, _ := s.GetDevice("d1")
	if got.State != "ACTIVE" || got.AccountID != "acc-1" {
		t.Fatalf("expected ACTIVE/acc-1, got %q/%q", got.State, got.AccountID)
	}
	if got.RegistrationExpiresAt == "" {
		t.Fatal("expected RegistrationExpiresAt to be set")
	}
}

func TestUpdateDeviceRegistration(t *testing.T) {
	s := newTestStore(t)
	s.InsertDevice(makeDevice("d1", "acc-1", "SOURCE", "ACTIVE"))

	reg := json.RawMessage(`{"channels":[{"id":"ch1","name":"Channel 1"}]}`)
	if err := s.UpdateDeviceRegistration("d1", reg); err != nil {
		t.Fatalf("UpdateDeviceRegistration: %v", err)
	}
	got, _ := s.GetDevice("d1")
	if string(got.Registration) != string(reg) {
		t.Fatalf("expected %s, got %s", reg, got.Registration)
	}
}

func TestUpdateDeviceStatus(t *testing.T) {
	s := newTestStore(t)
	s.InsertDevice(makeDevice("d1", "acc-1", "SOURCE", "ACTIVE"))

	status := json.RawMessage(`{"channels":[{"id":"ch1","state":"ACTIVE"}]}`)
	s.UpdateDeviceStatus("d1", status)
	got, _ := s.GetDevice("d1")
	if string(got.Status) != string(status) {
		t.Fatalf("expected %s, got %s", status, got.Status)
	}
}

func TestUpdateDeviceActualConfig(t *testing.T) {
	s := newTestStore(t)
	s.InsertDevice(makeDevice("d1", "acc-1", "SOURCE", "ACTIVE"))

	cfg := json.RawMessage(`{"channels":[{"id":"ch1","state":"IDLE"}]}`)
	s.UpdateDeviceActualConfig("d1", cfg)
	got, _ := s.GetDevice("d1")
	if string(got.ActualConfig) != string(cfg) {
		t.Fatalf("expected %s, got %s", cfg, got.ActualConfig)
	}
}

func TestUpdateDeviceDesiredConfig(t *testing.T) {
	s := newTestStore(t)
	s.InsertDevice(makeDevice("d1", "acc-1", "SOURCE", "ACTIVE"))

	cfg := json.RawMessage(`{"channels":[]}`)
	id1, err := s.UpdateDeviceDesiredConfig("d1", cfg)
	if err != nil {
		t.Fatalf("UpdateDeviceDesiredConfig: %v", err)
	}
	if id1 != 1 {
		t.Fatalf("expected updateID=1, got %d", id1)
	}

	id2, _ := s.UpdateDeviceDesiredConfig("d1", cfg)
	if id2 != 2 {
		t.Fatalf("expected updateID=2, got %d", id2)
	}
}

func TestUpdateDeviceOnline(t *testing.T) {
	s := newTestStore(t)
	s.InsertDevice(makeDevice("d1", "acc-1", "SOURCE", "ACTIVE"))

	now := time.Now().UTC().Format(time.RFC3339)
	s.UpdateDeviceOnline("d1", true, "192.168.1.1", now)
	got, _ := s.GetDevice("d1")
	if !got.Online {
		t.Fatal("expected online=true")
	}
	if got.SourceIP != "192.168.1.1" {
		t.Fatalf("expected 192.168.1.1, got %q", got.SourceIP)
	}
}

func TestUpdateDeviceCerts(t *testing.T) {
	s := newTestStore(t)
	s.InsertDevice(makeDevice("d1", "acc-1", "SOURCE", "ACTIVE"))

	now := time.Now().UTC().Format(time.RFC3339)
	expires := time.Now().Add(30 * 24 * time.Hour).UTC().Format(time.RFC3339)
	s.UpdateDeviceCerts("d1", "new-cert", "old-cert", expires, now, now)

	got, _ := s.GetDevice("d1")
	if got.CurrentCertPEM != "new-cert" {
		t.Fatalf("expected new-cert, got %q", got.CurrentCertPEM)
	}
	if got.PreviousCertPEM != "old-cert" {
		t.Fatalf("expected old-cert, got %q", got.PreviousCertPEM)
	}
}

func TestRevokePreviousCert(t *testing.T) {
	s := newTestStore(t)
	d := makeDevice("d1", "acc-1", "SOURCE", "ACTIVE")
	d.PreviousCertPEM = "old"
	d.PrevCertExpiresAt = "2025-01-01T00:00:00Z"
	s.InsertDevice(d)

	s.RevokePreviousCert("d1")
	got, _ := s.GetDevice("d1")
	if got.PreviousCertPEM != "" || got.PrevCertExpiresAt != "" {
		t.Fatal("expected previous cert fields to be cleared")
	}
}

func TestDeleteDevice(t *testing.T) {
	s := newTestStore(t)
	s.InsertDevice(makeDevice("d1", "acc-1", "SOURCE", "ACTIVE"))

	s.DeleteDevice("d1")
	got, _ := s.GetDevice("d1")
	if got != nil {
		t.Fatal("expected nil after delete")
	}
}

// --- Expiry Queries ---

func TestGetExpiredPairingDevices(t *testing.T) {
	s := newTestStore(t)
	d := makeDevice("d1", "", "SOURCE", "PAIRING")
	d.PairingExpiresAt = time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	s.InsertDevice(d)

	d2 := makeDevice("d2", "", "SOURCE", "PAIRING")
	d2.PairingExpiresAt = time.Now().Add(1 * time.Hour).UTC().Format(time.RFC3339)
	s.InsertDevice(d2)

	ids, err := s.GetExpiredPairingDevices(time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("GetExpiredPairingDevices: %v", err)
	}
	if len(ids) != 1 || ids[0] != "d1" {
		t.Fatalf("expected [d1], got %v", ids)
	}
}

func TestGetExpiredRegistrationDevices(t *testing.T) {
	s := newTestStore(t)
	d := makeDevice("d1", "acc-1", "SOURCE", "ACTIVE")
	d.RegistrationExpiresAt = time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	s.InsertDevice(d)

	d2 := makeDevice("d2", "acc-1", "SOURCE", "ACTIVE")
	d2.RegistrationExpiresAt = time.Now().Add(1 * time.Hour).UTC().Format(time.RFC3339)
	s.InsertDevice(d2)

	ids, err := s.GetExpiredRegistrationDevices(time.Now().UTC().Format(time.RFC3339))
	if err != nil {
		t.Fatalf("GetExpiredRegistrationDevices: %v", err)
	}
	if len(ids) != 1 || ids[0] != "d1" {
		t.Fatalf("expected [d1], got %v", ids)
	}
}

func TestGetDevicesNeedingRotation(t *testing.T) {
	s := newTestStore(t)
	d := makeDevice("d1", "acc-1", "SOURCE", "ACTIVE")
	d.CSRPEM = "csr-data"
	d.CurrentCertPEM = "cert"
	d.CertExpiresAt = "2025-12-01T00:00:00Z"
	// No LastRotationAt — should be included
	s.InsertDevice(d)

	d2 := makeDevice("d2", "acc-1", "SOURCE", "ACTIVE")
	d2.CSRPEM = "csr-data"
	d2.LastRotationAt = time.Now().UTC().Format(time.RFC3339) // recent
	s.InsertDevice(d2)

	threshold := time.Now().Add(-1 * time.Hour).UTC().Format(time.RFC3339)
	devices, err := s.GetDevicesNeedingRotation(threshold)
	if err != nil {
		t.Fatalf("GetDevicesNeedingRotation: %v", err)
	}
	if len(devices) != 1 || devices[0].DeviceID != "d1" {
		t.Fatalf("expected [d1], got %v", devices)
	}
}

// --- Accounts ---

func TestCreateAndGetAccount(t *testing.T) {
	s := newTestStore(t)
	a := &models.Account{
		AccountID:    "acc_12345678",
		Username:     "testuser",
		PasswordHash: "hash",
		DisplayName:  "Test User",
		CreatedAt:    time.Now().UTC().Format(time.RFC3339),
	}
	if err := s.CreateAccount(a); err != nil {
		t.Fatalf("CreateAccount: %v", err)
	}

	got, err := s.GetAccountByUsername("testuser")
	if err != nil {
		t.Fatalf("GetAccountByUsername: %v", err)
	}
	if got == nil || got.AccountID != "acc_12345678" {
		t.Fatalf("expected acc_12345678, got %+v", got)
	}

	got, err = s.GetAccountByID("acc_12345678")
	if err != nil {
		t.Fatalf("GetAccountByID: %v", err)
	}
	if got == nil || got.Username != "testuser" {
		t.Fatalf("expected testuser, got %+v", got)
	}
}

func TestGetAccount_NotFound(t *testing.T) {
	s := newTestStore(t)
	got, _ := s.GetAccountByUsername("nobody")
	if got != nil {
		t.Fatal("expected nil")
	}
	got, _ = s.GetAccountByID("acc_00000000")
	if got != nil {
		t.Fatal("expected nil")
	}
}

func TestDuplicateAccount(t *testing.T) {
	s := newTestStore(t)
	a := &models.Account{AccountID: "acc_1", Username: "user1", PasswordHash: "h", DisplayName: "U", CreatedAt: "now"}
	s.CreateAccount(a)
	a2 := &models.Account{AccountID: "acc_2", Username: "user1", PasswordHash: "h", DisplayName: "U", CreatedAt: "now"}
	if err := s.CreateAccount(a2); err == nil {
		t.Fatal("expected error for duplicate username")
	}
}

// --- Thumbnails ---

func TestThumbnailCRUD(t *testing.T) {
	s := newTestStore(t)
	s.InsertDevice(makeDevice("d1", "acc-1", "SOURCE", "ACTIVE"))

	thumb := &Thumbnail{
		DeviceID:    "d1",
		SourceID:    "ch1",
		ImageData:   []byte{0xFF, 0xD8, 0xFF},
		Timestamp:   "2025-01-01T00:00:00Z",
		ImageType:   "image/jpeg",
		ImageSizeKB: 50,
	}
	if err := s.UpsertThumbnail(thumb); err != nil {
		t.Fatalf("UpsertThumbnail: %v", err)
	}

	got, err := s.GetThumbnail("d1", "ch1")
	if err != nil {
		t.Fatalf("GetThumbnail: %v", err)
	}
	if got == nil {
		t.Fatal("expected thumbnail")
	}
	if len(got.ImageData) != 3 || got.ImageType != "image/jpeg" {
		t.Fatalf("unexpected thumbnail: %+v", got)
	}

	// Upsert replaces
	thumb.ImageData = []byte{0x01, 0x02}
	thumb.ImageSizeKB = 1
	s.UpsertThumbnail(thumb)
	got, _ = s.GetThumbnail("d1", "ch1")
	if len(got.ImageData) != 2 {
		t.Fatalf("expected updated image data, got %d bytes", len(got.ImageData))
	}

	// Not found
	got, _ = s.GetThumbnail("d1", "ch-none")
	if got != nil {
		t.Fatal("expected nil for unknown source")
	}

	// Delete
	s.DeleteThumbnailsByDevice("d1")
	got, _ = s.GetThumbnail("d1", "ch1")
	if got != nil {
		t.Fatal("expected nil after delete")
	}
}

// --- Logs ---

func TestLogCRUD(t *testing.T) {
	s := newTestStore(t)
	s.InsertDevice(makeDevice("d1", "acc-1", "SOURCE", "ACTIVE"))

	l := &DeviceLog{
		DeviceID:   "d1",
		LogData:    []byte("log content here"),
		UploadedAt: "2025-01-01T00:00:00Z",
		LogSizeKB:  1,
	}
	if err := s.UpsertLog(l); err != nil {
		t.Fatalf("UpsertLog: %v", err)
	}

	got, err := s.GetLog("d1")
	if err != nil {
		t.Fatalf("GetLog: %v", err)
	}
	if got == nil || string(got.LogData) != "log content here" {
		t.Fatalf("unexpected log: %+v", got)
	}

	// Upsert replaces
	l.LogData = []byte("updated")
	s.UpsertLog(l)
	got, _ = s.GetLog("d1")
	if string(got.LogData) != "updated" {
		t.Fatalf("expected updated, got %q", got.LogData)
	}

	// Not found
	got, _ = s.GetLog("d-none")
	if got != nil {
		t.Fatal("expected nil")
	}

	// Delete
	s.DeleteLogsByDevice("d1")
	got, _ = s.GetLog("d1")
	if got != nil {
		t.Fatal("expected nil after delete")
	}
}

// --- Config ---

func TestConfigCRUD(t *testing.T) {
	s := newTestStore(t)

	if err := s.SetConfig("test_key", []byte("test_value")); err != nil {
		t.Fatalf("SetConfig: %v", err)
	}
	val, err := s.GetConfig("test_key")
	if err != nil {
		t.Fatalf("GetConfig: %v", err)
	}
	if string(val) != "test_value" {
		t.Fatalf("expected test_value, got %q", val)
	}

	// Upsert
	s.SetConfig("test_key", []byte("updated"))
	val, _ = s.GetConfig("test_key")
	if string(val) != "updated" {
		t.Fatalf("expected updated, got %q", val)
	}

	// Not found
	_, err = s.GetConfig("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent key")
	}
}
