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
// CredentialStore persists identity, X.509 certs, and host settings on the filesystem.
// Mirrors the Python SDK's CredentialStore class.
package credentials

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/utils"
	tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
)

// ConnectionSettings persists the device identity and MQTT endpoint.
type ConnectionSettings struct {
	DeviceID string `json:"device_id"`
	URI      string `json:"uri"`
	Region   string `json:"region"`
}

// Store manages certificate and settings persistence.
type Store struct {
	mu               sync.Mutex
	DeviceLocalID    string
	Base             string
	Dir              string
	CACertFile       string
	DeviceCertFile   string
	PrivKeyFile      string
	HostSettingsFile string
	ConnSettingsFile string
	HostSettings     *tr12models.HostSettings
	ConnSettings     *ConnectionSettings
	PrivKey          string
	PubKey           string
	CSR              string
}

// NewStore creates a new credential store for the given host.
func NewStore(base, deviceLocalID, hostID string) (*Store, error) {
	dir := filepath.Join(base, deviceLocalID, hostID)
	log.Printf("[CREDS] Store path: %s (base=%s deviceLocalID=%s hostID=%s)", dir, base, deviceLocalID, hostID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("certs directory is not writable: %s: %w", base, err)
	}
	return &Store{
		DeviceLocalID:    deviceLocalID,
		Base:             base,
		Dir:              dir,
		CACertFile:       filepath.Join(dir, "ca_cert"),
		DeviceCertFile:   filepath.Join(dir, "device_cert"),
		PrivKeyFile:      filepath.Join(dir, "priv_key"),
		HostSettingsFile: filepath.Join(dir, "host_settings"),
		ConnSettingsFile: filepath.Join(dir, "connection_settings"),
	}, nil
}

// GetDeviceID returns the connected device ID.
func (s *Store) GetDeviceID() string {
	if s.ConnSettings == nil {
		return ""
	}
	return s.ConnSettings.DeviceID
}

// GetURI returns the MQTT broker URI.
func (s *Store) GetURI() string {
	if s.ConnSettings == nil {
		return ""
	}
	return s.ConnSettings.URI
}

// GetRegion returns the connected region.
func (s *Store) GetRegion() string {
	if s.ConnSettings == nil {
		return ""
	}
	return s.ConnSettings.Region
}

// GetHostSettings returns the host settings or an error if not initialized.
func (s *Store) GetHostSettings() (*tr12models.HostSettings, error) {
	if s.HostSettings == nil {
		return nil, fmt.Errorf("host settings not initialized, likely not connected")
	}
	return s.HostSettings, nil
}

// GenerateKeysAndCSR generates RSA keys and a CSR for pairing.
func (s *Store) GenerateKeysAndCSR() error {
	if s.CSR != "" {
		return nil // Already generated
	}
	pub, priv, err := utils.GenerateClientKeys()
	if err != nil {
		return err
	}
	s.PubKey = pub
	s.PrivKey = priv
	csr, err := utils.GenerateCSR(priv)
	if err != nil {
		return err
	}
	s.CSR = csr
	return nil
}

// ReadFromFilesystem loads persisted certs. Returns true if certs exist and are valid.
func (s *Store) ReadFromFilesystem() (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := os.Stat(s.Dir); os.IsNotExist(err) {
		return false, nil
	}
	if _, err := os.Stat(s.CACertFile); os.IsNotExist(err) {
		return false, nil
	}
	for _, f := range []string{s.CACertFile, s.DeviceCertFile, s.PrivKeyFile, s.ConnSettingsFile, s.HostSettingsFile} {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			return false, fmt.Errorf("missing: %s. Should deregister and re-pair the device", f)
		}
	}
	data, err := os.ReadFile(s.ConnSettingsFile)
	if err != nil {
		return false, fmt.Errorf("invalid connection_settings file: %w", err)
	}
	var cs ConnectionSettings
	if err := json.Unmarshal(data, &cs); err != nil {
		return false, fmt.Errorf("invalid connection_settings: %w", err)
	}
	s.ConnSettings = &cs
	log.Printf("[CREDS] Loaded connection_settings from %s: deviceId=%s uri=%q region=%s",
		s.ConnSettingsFile, cs.DeviceID, cs.URI, cs.Region)

	data, err = os.ReadFile(s.HostSettingsFile)
	if err != nil {
		return false, fmt.Errorf("invalid host_settings: %w", err)
	}
	var hs tr12models.HostSettings
	if err := json.Unmarshal(data, &hs); err != nil {
		return false, fmt.Errorf("invalid host_settings: %w", err)
	}
	s.HostSettings = &hs
	return true, nil
}

// WriteToFilesystem saves certs and settings after successful authentication.
func (s *Store) WriteToFilesystem(deviceID string, auth *tr12models.AuthenticatePairingCodeResponseContent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(s.Dir, 0755); err != nil {
		return fmt.Errorf("unable to create certs directory: %w", err)
	}
	if err := os.WriteFile(s.CACertFile, []byte(auth.GetCaCertificate()), 0600); err != nil {
		return fmt.Errorf("unable to write ca_cert: %w", err)
	}
	if err := os.WriteFile(s.DeviceCertFile, []byte(auth.GetDeviceCertificate()), 0600); err != nil {
		return fmt.Errorf("unable to write device_cert: %w", err)
	}
	if err := os.WriteFile(s.PrivKeyFile, []byte(s.PrivKey), 0600); err != nil {
		return fmt.Errorf("unable to write priv_key: %w", err)
	}
	s.ConnSettings = &ConnectionSettings{
		DeviceID: deviceID,
		URI:      auth.GetMqttUri(),
		Region:   auth.GetRegionName(),
	}
	log.Printf("[CREDS] Writing connection_settings to %s: deviceId=%s uri=%q regionName=%s",
		s.ConnSettingsFile, deviceID, auth.GetMqttUri(), auth.GetRegionName())
	csData, _ := json.Marshal(s.ConnSettings)
	if err := os.WriteFile(s.ConnSettingsFile, csData, 0600); err != nil {
		return fmt.Errorf("unable to write connection_settings: %w", err)
	}
	hsData, _ := json.Marshal(auth.HostSettings)
	if err := os.WriteFile(s.HostSettingsFile, hsData, 0600); err != nil {
		return fmt.Errorf("unable to write host_settings: %w", err)
	}
	s.HostSettings = auth.HostSettings
	return nil
}

// RotateCerts updates the device cert and connection settings if changed. Returns true if updated.
func (s *Store) RotateCerts(rotate *tr12models.RotateCertificatesRequestContent) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	needUpdate := false
	currentCert, _ := os.ReadFile(s.DeviceCertFile)
	if rotate.DeviceCertificate != string(currentCert) {
		if err := os.WriteFile(s.DeviceCertFile, []byte(rotate.DeviceCertificate), 0600); err != nil {
			return false, fmt.Errorf("unable to write rotated device_cert: %w", err)
		}
		needUpdate = true
	}
	if s.ConnSettings != nil && (rotate.MqttUri != s.ConnSettings.URI || rotate.GetRegionName() != s.ConnSettings.Region) {
		s.ConnSettings.URI = rotate.MqttUri
		s.ConnSettings.Region = rotate.GetRegionName()
		csData, _ := json.Marshal(s.ConnSettings)
		if err := os.WriteFile(s.ConnSettingsFile, csData, 0600); err != nil {
			return false, fmt.Errorf("unable to write rotated connection_settings: %w", err)
		}
		needUpdate = true
	}
	return needUpdate, nil
}

// Deprovision removes all credentials from the filesystem.
func (s *Store) Deprovision() error {
	return os.RemoveAll(s.Dir)
}
