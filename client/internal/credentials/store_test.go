package credentials

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
)

func TestNewStore(t *testing.T) {
	dir := t.TempDir()
	s, err := NewStore(dir, "device-local-1", "host-1")
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	if s.Dir != filepath.Join(dir, "device-local-1", "host-1") {
		t.Fatalf("unexpected Dir: %s", s.Dir)
	}
	info, err := os.Stat(s.Dir)
	if err != nil || !info.IsDir() {
		t.Fatalf("expected directory to be created at %s", s.Dir)
	}
}

func TestGetDeviceID_NilConnSettings(t *testing.T) {
	s := &Store{}
	if s.GetDeviceID() != "" {
		t.Fatal("expected empty device ID when ConnSettings is nil")
	}
}

func TestGetURI_NilConnSettings(t *testing.T) {
	s := &Store{}
	if s.GetURI() != "" {
		t.Fatal("expected empty URI when ConnSettings is nil")
	}
}

func TestGetRegion_NilConnSettings(t *testing.T) {
	s := &Store{}
	if s.GetRegion() != "" {
		t.Fatal("expected empty region when ConnSettings is nil")
	}
}

func TestGetters_WithConnSettings(t *testing.T) {
	s := &Store{ConnSettings: &ConnectionSettings{
		DeviceID: "dev-123",
		URI:      "tls://broker:8883",
		Region:   "us-west-2",
	}}
	if s.GetDeviceID() != "dev-123" {
		t.Fatalf("expected dev-123, got %q", s.GetDeviceID())
	}
	if s.GetURI() != "tls://broker:8883" {
		t.Fatalf("expected tls://broker:8883, got %q", s.GetURI())
	}
	if s.GetRegion() != "us-west-2" {
		t.Fatalf("expected us-west-2, got %q", s.GetRegion())
	}
}

func TestGetHostSettings_Nil(t *testing.T) {
	s := &Store{}
	_, err := s.GetHostSettings()
	if err == nil {
		t.Fatal("expected error when HostSettings is nil")
	}
}

func TestGenerateKeysAndCSR(t *testing.T) {
	s := &Store{}
	if err := s.GenerateKeysAndCSR(); err != nil {
		t.Fatalf("GenerateKeysAndCSR: %v", err)
	}
	if s.PubKey == "" || s.PrivKey == "" || s.CSR == "" {
		t.Fatal("expected PubKey, PrivKey, and CSR to be populated")
	}
	// Second call should be a no-op
	origCSR := s.CSR
	if err := s.GenerateKeysAndCSR(); err != nil {
		t.Fatalf("second GenerateKeysAndCSR: %v", err)
	}
	if s.CSR != origCSR {
		t.Fatal("expected CSR to remain unchanged on second call")
	}
}

func TestReadFromFilesystem_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	s := &Store{Dir: filepath.Join(dir, "nonexistent")}
	ok, err := s.ReadFromFilesystem()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected false for nonexistent dir")
	}
}

func TestReadFromFilesystem_MissingFiles(t *testing.T) {
	dir := t.TempDir()
	storeDir := filepath.Join(dir, "dev", "host")
	os.MkdirAll(storeDir, 0755)
	// Write only ca_cert — the rest are missing
	os.WriteFile(filepath.Join(storeDir, "ca_cert"), []byte("ca"), 0600)

	s := &Store{
		Dir:              storeDir,
		CACertFile:       filepath.Join(storeDir, "ca_cert"),
		DeviceCertFile:   filepath.Join(storeDir, "device_cert"),
		PrivKeyFile:      filepath.Join(storeDir, "priv_key"),
		ConnSettingsFile: filepath.Join(storeDir, "connection_settings"),
		HostSettingsFile: filepath.Join(storeDir, "host_settings"),
	}
	ok, err := s.ReadFromFilesystem()
	if err == nil {
		t.Fatal("expected error for missing files")
	}
	if ok {
		t.Fatal("expected false")
	}
}

func TestWriteAndReadFromFilesystem(t *testing.T) {
	dir := t.TempDir()
	s, err := NewStore(dir, "dev1", "host1")
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	s.PrivKey = "fake-private-key"

	hs := tr12models.NewHostSettings(
		"mqtt", 1800, 1, 30,
		"sub/update", "sub/thumb", "pub/schema", "pub/reg",
		"pub/status", "pub/config", "sub/certs", "sub/deprov",
		"pub/deprov", "sub/log",
	)
	auth := tr12models.NewAuthenticateResponseContent(tr12models.CLAIMED)
	auth.SetCaCert("fake-ca-cert")
	auth.SetDeviceCert("fake-device-cert")
	auth.SetMqttUri("tls://localhost:8883")
	auth.SetRegion("local")
	auth.SetHostSettings(*hs)

	if err := s.WriteToFilesystem("device-abc", auth); err != nil {
		t.Fatalf("WriteToFilesystem: %v", err)
	}

	// Verify files exist
	for _, f := range []string{s.CACertFile, s.DeviceCertFile, s.PrivKeyFile, s.ConnSettingsFile, s.HostSettingsFile} {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected file %s to exist", f)
		}
	}

	// Read back
	s2, _ := NewStore(dir, "dev1", "host1")
	ok, err := s2.ReadFromFilesystem()
	if err != nil {
		t.Fatalf("ReadFromFilesystem: %v", err)
	}
	if !ok {
		t.Fatal("expected true from ReadFromFilesystem")
	}
	if s2.GetDeviceID() != "device-abc" {
		t.Fatalf("expected device-abc, got %q", s2.GetDeviceID())
	}
	if s2.GetURI() != "tls://localhost:8883" {
		t.Fatalf("expected tls://localhost:8883, got %q", s2.GetURI())
	}
	if s2.GetRegion() != "local" {
		t.Fatalf("expected local, got %q", s2.GetRegion())
	}
}

func TestRotateCerts_Changed(t *testing.T) {
	dir := t.TempDir()
	s, _ := NewStore(dir, "dev1", "host1")
	s.PrivKey = "key"

	hs := tr12models.NewHostSettings("mqtt", 1800, 1, 30, "a", "b", "c", "d", "e", "f", "g", "h", "i", "j")
	auth := tr12models.NewAuthenticateResponseContent(tr12models.CLAIMED)
	auth.SetCaCert("ca")
	auth.SetDeviceCert("old-cert")
	auth.SetMqttUri("tls://old:8883")
	auth.SetRegion("old-region")
	auth.SetHostSettings(*hs)
	s.WriteToFilesystem("dev1", auth)

	rotate := &tr12models.RotateCertificatesRequestContent{
		MqttUri:    "tls://new:8883",
		DeviceCert: "new-cert",
		Region:     "new-region",
	}
	updated, err := s.RotateCerts(rotate)
	if err != nil {
		t.Fatalf("RotateCerts: %v", err)
	}
	if !updated {
		t.Fatal("expected updated=true")
	}

	// Verify cert file changed
	data, _ := os.ReadFile(s.DeviceCertFile)
	if string(data) != "new-cert" {
		t.Fatalf("expected new-cert, got %q", string(data))
	}

	// Verify connection settings updated
	csData, _ := os.ReadFile(s.ConnSettingsFile)
	var cs ConnectionSettings
	json.Unmarshal(csData, &cs)
	if cs.URI != "tls://new:8883" || cs.Region != "new-region" {
		t.Fatalf("expected updated conn settings, got URI=%q Region=%q", cs.URI, cs.Region)
	}
}

func TestRotateCerts_NothingChanged(t *testing.T) {
	dir := t.TempDir()
	s, _ := NewStore(dir, "dev1", "host1")
	s.PrivKey = "key"

	hs := tr12models.NewHostSettings("mqtt", 1800, 1, 30, "a", "b", "c", "d", "e", "f", "g", "h", "i", "j")
	auth := tr12models.NewAuthenticateResponseContent(tr12models.CLAIMED)
	auth.SetCaCert("ca")
	auth.SetDeviceCert("same-cert")
	auth.SetMqttUri("tls://same:8883")
	auth.SetRegion("same-region")
	auth.SetHostSettings(*hs)
	s.WriteToFilesystem("dev1", auth)

	rotate := &tr12models.RotateCertificatesRequestContent{
		MqttUri:    "tls://same:8883",
		DeviceCert: "same-cert",
		Region:     "same-region",
	}
	updated, err := s.RotateCerts(rotate)
	if err != nil {
		t.Fatalf("RotateCerts: %v", err)
	}
	if updated {
		t.Fatal("expected updated=false when nothing changed")
	}
}

func TestDeprovision(t *testing.T) {
	dir := t.TempDir()
	s, _ := NewStore(dir, "dev1", "host1")
	s.PrivKey = "key"

	hs := tr12models.NewHostSettings("mqtt", 1800, 1, 30, "a", "b", "c", "d", "e", "f", "g", "h", "i", "j")
	auth := tr12models.NewAuthenticateResponseContent(tr12models.CLAIMED)
	auth.SetCaCert("ca")
	auth.SetDeviceCert("cert")
	auth.SetMqttUri("tls://x:8883")
	auth.SetRegion("local")
	auth.SetHostSettings(*hs)
	s.WriteToFilesystem("dev1", auth)

	if err := s.Deprovision(); err != nil {
		t.Fatalf("Deprovision: %v", err)
	}
	if _, err := os.Stat(s.Dir); !os.IsNotExist(err) {
		t.Fatal("expected directory to be removed after deprovision")
	}
}
