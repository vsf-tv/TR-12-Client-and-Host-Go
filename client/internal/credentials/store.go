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
	"strings"
	"sync"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/utils"
	tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
)

// Two-phase commit constants for atomic multi-file updates.
//
//   Phase 1: write every changed file as <path>.new (fsync each; fsync dir).
//   Phase 2: write new-creds-saved.done (the commit marker) then rename each
//            .new to its final name, then remove the marker.
//
// Recovery on startup / before any read or write: if the marker exists, finish
// any remaining renames and delete it. If it does not exist, delete any orphan
// .new files (they are from an aborted Phase 1). This makes rotation atomic across
// process crash and power loss.
const (
	newSuffix          = ".new"
	rotationDoneMarker = "new-creds-saved.done"
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
	// Complete or discard any pending rotation from a prior unclean shutdown.
	if err := recoverPendingRotation(dir); err != nil {
		return nil, fmt.Errorf("credential recovery failed: %w", err)
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

// writeNewFile writes data to path+".new" and fsyncs the file.
// It does NOT rename — commitPendingRotation does that after every .new is on disk.
func writeNewFile(path string, data []byte, perm os.FileMode) error {
	tmp := path + newSuffix
	f, err := os.OpenFile(tmp, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		_ = os.Remove(tmp)
		return err
	}
	return f.Close()
}

// fsyncDir persists directory metadata (rename / create / unlink) to disk.
func fsyncDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	return d.Sync()
}

// commitPendingRotation is Phase 2: writes the done marker, then renames every
// .new file in the directory to its final name, then removes the marker.
// Idempotent if crashed mid-way — the next call to recoverPendingRotation finishes it.
func commitPendingRotation(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir for commit: %w", err)
	}
	var toRename []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), newSuffix) {
			toRename = append(toRename, e.Name())
		}
	}
	if len(toRename) == 0 {
		return nil // nothing prepared, nothing to commit
	}
	if err := fsyncDir(dir); err != nil {
		return fmt.Errorf("fsync dir before marker: %w", err)
	}
	marker := filepath.Join(dir, rotationDoneMarker)
	f, err := os.OpenFile(marker, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("create done marker: %w", err)
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return fmt.Errorf("fsync done marker: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("close done marker: %w", err)
	}
	// Marker is now durably on disk — this is the commit point.
	if err := fsyncDir(dir); err != nil {
		return fmt.Errorf("fsync dir after marker: %w", err)
	}
	for _, name := range toRename {
		tmp := filepath.Join(dir, name)
		real := filepath.Join(dir, strings.TrimSuffix(name, newSuffix))
		if err := os.Rename(tmp, real); err != nil {
			return fmt.Errorf("rename %s -> %s: %w", tmp, real, err)
		}
	}
	if err := fsyncDir(dir); err != nil {
		return fmt.Errorf("fsync dir after renames: %w", err)
	}
	if err := os.Remove(marker); err != nil {
		return fmt.Errorf("remove done marker: %w", err)
	}
	return nil
}

// recoverPendingRotation brings the on-disk state to consistency after a crash:
//   - If the done marker exists: finish any remaining renames, remove the marker.
//     Any .new files were valid Phase 2 pending work.
//   - Otherwise: any .new files are orphaned Phase 1 work — delete them.
//
// Safe to call when the directory is empty or does not exist.
func recoverPendingRotation(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read dir for recovery: %w", err)
	}
	marker := filepath.Join(dir, rotationDoneMarker)
	_, markerErr := os.Stat(marker)
	doneExists := markerErr == nil

	if doneExists {
		for _, e := range entries {
			if !strings.HasSuffix(e.Name(), newSuffix) {
				continue
			}
			tmp := filepath.Join(dir, e.Name())
			real := filepath.Join(dir, strings.TrimSuffix(e.Name(), newSuffix))
			if err := os.Rename(tmp, real); err != nil {
				return fmt.Errorf("recovery rename %s: %w", tmp, err)
			}
			log.Printf("[CREDS] recovered pending rename: %s -> %s", tmp, real)
		}
		_ = fsyncDir(dir)
		if err := os.Remove(marker); err != nil {
			return fmt.Errorf("recovery remove marker: %w", err)
		}
		_ = fsyncDir(dir)
		return nil
	}
	// No marker — any .new files are orphaned from an aborted Phase 1.
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), newSuffix) {
			p := filepath.Join(dir, e.Name())
			_ = os.Remove(p)
			log.Printf("[CREDS] discarded orphaned .new file: %s", p)
		}
	}
	return nil
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
	// Complete or discard any pending rotation before observing state.
	if err := recoverPendingRotation(s.Dir); err != nil {
		return false, fmt.Errorf("recover before read: %w", err)
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
// Uses the two-phase commit pattern so a crash mid-write leaves either the old
// state or the new state on disk — never a partial mix.
func (s *Store) WriteToFilesystem(deviceID string, auth *tr12models.AuthenticatePairingCodeResponseContent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(s.Dir, 0755); err != nil {
		return fmt.Errorf("unable to create certs directory: %w", err)
	}
	if err := recoverPendingRotation(s.Dir); err != nil {
		return fmt.Errorf("recover before write: %w", err)
	}

	cs := &ConnectionSettings{
		DeviceID: deviceID,
		URI:      auth.GetMqttUri(),
		Region:   auth.GetRegionName(),
	}
	csData, _ := json.Marshal(cs)
	hsData, _ := json.Marshal(auth.HostSettings)

	writes := []struct {
		path string
		data []byte
	}{
		{s.CACertFile, []byte(auth.GetCaCertificate())},
		{s.DeviceCertFile, []byte(auth.GetDeviceCertificate())},
		{s.PrivKeyFile, []byte(s.PrivKey)},
		{s.ConnSettingsFile, csData},
		{s.HostSettingsFile, hsData},
	}
	for _, w := range writes {
		if err := writeNewFile(w.path, w.data, 0600); err != nil {
			_ = recoverPendingRotation(s.Dir) // clean up whichever .new files did land
			return fmt.Errorf("unable to write %s.new: %w", filepath.Base(w.path), err)
		}
	}
	log.Printf("[CREDS] Writing connection_settings to %s: deviceId=%s uri=%q regionName=%s",
		s.ConnSettingsFile, deviceID, auth.GetMqttUri(), auth.GetRegionName())
	if err := commitPendingRotation(s.Dir); err != nil {
		// Recovery on the next call to any store operation will finish or discard.
		return fmt.Errorf("commit initial credentials: %w", err)
	}
	s.ConnSettings = cs
	s.HostSettings = auth.HostSettings
	return nil
}

// RotateCerts updates the device cert, optional CA cert, and connection settings
// if any have changed. Returns true if any were updated.
//
// Uses the two-phase commit pattern: every changed file is written as <name>.new
// with fsync; then a marker file is created and each .new is renamed to its
// final name. A crash between renames is recovered on next startup — the on-disk
// state is always either fully-old or fully-new.
//
// CA cert is optional: when the payload omits it, the existing ca_cert on disk
// is reused (backward-compatible with hosts that do not send the CA).
func (s *Store) RotateCerts(rotate *tr12models.DeviceSubscribesToCertificateRotationResponseContent) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := recoverPendingRotation(s.Dir); err != nil {
		return false, fmt.Errorf("recover before rotate: %w", err)
	}

	needUpdate := false
	var pendingConn *ConnectionSettings

	currentCert, _ := os.ReadFile(s.DeviceCertFile)
	if rotate.DeviceCertificate != string(currentCert) {
		if err := writeNewFile(s.DeviceCertFile, []byte(rotate.DeviceCertificate), 0600); err != nil {
			_ = recoverPendingRotation(s.Dir)
			return false, fmt.Errorf("write device_cert.new: %w", err)
		}
		needUpdate = true
	}
	if newCA := rotate.GetCaCertificate(); newCA != "" {
		currentCA, _ := os.ReadFile(s.CACertFile)
		if newCA != string(currentCA) {
			if err := writeNewFile(s.CACertFile, []byte(newCA), 0600); err != nil {
				_ = recoverPendingRotation(s.Dir)
				return false, fmt.Errorf("write ca_cert.new: %w", err)
			}
			needUpdate = true
		}
	}
	if s.ConnSettings != nil && (rotate.MqttUri != s.ConnSettings.URI || rotate.GetRegionName() != s.ConnSettings.Region) {
		cs := *s.ConnSettings
		cs.URI = rotate.MqttUri
		cs.Region = rotate.GetRegionName()
		csData, _ := json.Marshal(cs)
		if err := writeNewFile(s.ConnSettingsFile, csData, 0600); err != nil {
			_ = recoverPendingRotation(s.Dir)
			return false, fmt.Errorf("write connection_settings.new: %w", err)
		}
		pendingConn = &cs
		needUpdate = true
	}

	if !needUpdate {
		return false, nil
	}
	if err := commitPendingRotation(s.Dir); err != nil {
		// A crash inside commit is safe — the next recoverPendingRotation completes it.
		return false, fmt.Errorf("commit rotation: %w", err)
	}
	if pendingConn != nil {
		s.ConnSettings = pendingConn
	}
	return true, nil
}

// Deprovision removes all credentials from the filesystem.
func (s *Store) Deprovision() error {
	return os.RemoveAll(s.Dir)
}
