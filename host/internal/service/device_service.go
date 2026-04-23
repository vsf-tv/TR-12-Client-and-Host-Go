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
	"encoding/json"
	"strconv"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/ca"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/config"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/db"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/models"
	tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
)

// MQTTPublisher is the interface for publishing MQTT messages.
type MQTTPublisher interface {
	Publish(topic string, payload []byte, retain bool) error
}

// DeviceService handles device pairing, authentication, and management.
type DeviceService struct {
	store *db.Store
	ca    *ca.CA
	mqtt  MQTTPublisher
	cfg   *config.Config
}

// NewDeviceService creates a new DeviceService.
func NewDeviceService(store *db.Store, ca *ca.CA, mqtt MQTTPublisher, cfg *config.Config) *DeviceService {
	return &DeviceService{store: store, ca: ca, mqtt: mqtt, cfg: cfg}
}

// SetMQTT sets the MQTT publisher (used for deferred wiring).
func (s *DeviceService) SetMQTT(mqtt MQTTPublisher) {
	s.mqtt = mqtt
}

// Pair handles a device pairing request.
func (s *DeviceService) Pair(req models.CreatePairingCodeRequestContent) (*models.CreatePairingCodeResponseContent, error) {
	// Validate host ID
	if req.HostId != s.cfg.ServiceID {
		return failPair(models.PairFailureHostIDMismatch), nil
	}
	// Validate version
	if req.Version.GetVersion() == "" {
		return failPair(models.PairFailureVersionNotSupported), nil
	}
	// Validate device type
	validType := false
	for _, dt := range []string{"SOURCE", "DESTINATION", "BOTH"} {
		if strings.EqualFold(string(req.DeviceType), dt) {
			validType = true
			break
		}
	}
	if !validType {
		return failPair(models.PairFailureDeviceTypeNotSupported), nil
	}

	deviceID := generateDeviceID()
	pairingCode := generatePairingCode()
	accessCode := generateAccessCode()

	// Sign CSR
	certPEM, err := s.ca.SignCSR([]byte(req.CertificateSigningRequest), deviceID, s.cfg.CertExpiryDays)
	if err != nil {
		return nil, fmt.Errorf("sign CSR: %w", err)
	}

	now := time.Now().UTC()
	device := &models.Device{
		DeviceID:         deviceID,
		AccountID:        "",
		DeviceType:       strings.ToUpper(string(req.DeviceType)),
		State:            "PAIRING",
		PairedAt:         now.Format(time.RFC3339),
		CurrentCertPEM:   string(certPEM),
		CertExpiresAt:    now.Add(time.Duration(s.cfg.CertExpiryDays) * 24 * time.Hour).Format(time.RFC3339),
		CSRPEM:           req.CertificateSigningRequest,
		PairingCode:      pairingCode,
		AccessCode:       accessCode,
		PairingExpiresAt: now.Add(time.Duration(s.cfg.PairingTimeout) * time.Second).Format(time.RFC3339),
	}
	if err := s.store.InsertDevice(device); err != nil {
		return nil, fmt.Errorf("insert device: %w", err)
	}

	successData := models.CreatePairingCodeSuccessData{
		DeviceId:              deviceID,
		PairingCode:           pairingCode,
		AccessCode:            accessCode,
		PairingTimeoutSeconds: float32(s.cfg.PairingTimeout),
	}
	return tr12models.NewCreatePairingCodeResponseContent(
		tr12models.SuccessAsCreatePairingCodeResult(tr12models.NewSuccess(successData)),
	), nil
}

// Authenticate handles device authentication polling.
func (s *DeviceService) Authenticate(req models.AuthenticatePairingCodeRequestContent) (*models.AuthenticatePairingCodeResponseContent, error) {
	device, err := s.store.GetDevice(req.DeviceId)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, ErrNotFound
	}
	if device.PairingCode != req.PairingCode || device.AccessCode != req.AccessCode {
		return nil, ErrUnauthorized
	}

	// Check if pairing expired
	if device.PairingExpiresAt != "" {
		expires, _ := time.Parse(time.RFC3339, device.PairingExpiresAt)
		if time.Now().UTC().After(expires) {
			return nil, fmt.Errorf("pairing expired")
		}
	}

	if device.State == "PAIRING" {
		resp := tr12models.NewAuthenticatePairingCodeResponseContent(models.AuthStatusSTANDBY)
		return resp, nil
	}

	// Device is claimed — return full auth response
	mqttURI := fmt.Sprintf("tls://%s:%d", s.cfg.HostAddress, s.cfg.MQTTPort)
	hs := buildHostSettings(device.DeviceID, s.cfg.PairingTimeout)
	resp := tr12models.NewAuthenticatePairingCodeResponseContent(models.AuthStatusCLAIMED)
	resp.SetCaCertificate(string(s.ca.CACertPEM))
	resp.SetDeviceCertificate(device.CurrentCertPEM)
	resp.SetMqttUri(mqttURI)
	resp.SetRegionName("local")
	resp.SetHostSettings(*hs)
	return resp, nil
}

// Claim associates a device with an account.
func (s *DeviceService) Claim(pairingCode, accountID string, expirationDays int, locationName, deviceName string, rotationIntervalDays int) error {
	device, err := s.store.GetDeviceByPairingCode(pairingCode)
	if err != nil {
		return err
	}
	if device == nil {
		return ErrNotFound
	}
	if device.State != "PAIRING" {
		return ErrConflict
	}
	// Check expiry
	if device.PairingExpiresAt != "" {
		expires, _ := time.Parse(time.RFC3339, device.PairingExpiresAt)
		if time.Now().UTC().After(expires) {
			return fmt.Errorf("pairing expired")
		}
	}
	if expirationDays <= 0 {
		expirationDays = 730
	}
	if rotationIntervalDays < 30 {
		rotationIntervalDays = 365
	}
	if rotationIntervalDays > 5*365 {
		rotationIntervalDays = 5 * 365
	}
	// Truncate strings to max 40 chars
	if len(locationName) > 40 {
		locationName = locationName[:40]
	}
	if len(deviceName) > 40 {
		deviceName = deviceName[:40]
	}
	regExpires := time.Now().UTC().Add(time.Duration(expirationDays) * 24 * time.Hour).Format(time.RFC3339)
	return s.store.ClaimDevice(device.DeviceID, accountID, regExpires, locationName, deviceName, rotationIntervalDays)
}

// ListDevices returns all non-deprovisioned devices for an account.
func (s *DeviceService) ListDevices(accountID string) ([]models.DeviceSummary, error) {
	devices, err := s.store.ListDevicesByAccount(accountID)
	if err != nil {
		return nil, err
	}
	summaries := make([]models.DeviceSummary, 0, len(devices))
	for _, d := range devices {
		if d.State == "DEPROVISIONED" {
			continue // hide deprovisioned devices from the list
		}
		summaries = append(summaries, models.DeviceSummary{
			DeviceID:      d.DeviceID,
			Message:       "",
			Errors:        []string{},
			OnlineDetails: formatOnlineDetails(d),
			Online:        d.Online,
			LocationName:  d.LocationName,
			DeviceName:    d.DeviceName,
		})
	}
	return summaries, nil
}

// DescribeDevice returns full device details.
func (s *DeviceService) DescribeDevice(deviceID, accountID string) (*models.DeviceDetail, error) {
	device, err := s.store.GetDevice(deviceID)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, ErrNotFound
	}
	if device.AccountID != accountID {
		return nil, ErrForbidden
	}
	onlineDetails := formatOnlineDetails(device)
	certExp := formatCertExpiration(device.CertExpiresAt)
	return &models.DeviceDetail{
		DeviceID:            device.DeviceID,
		Message:             "",
		Errors:              []string{},
		Registration:        device.Registration,
		Configuration:       device.DesiredConfig,
		ActualConfiguration: device.ActualConfig,
		Status:              device.Status,
		Online:              device.Online,
		OnlineDetails:       onlineDetails,
		CertExpiration:      certExp,
		DeviceMetadata: models.DeviceMetadata{
			Online:         device.Online,
			OnlineDetails:  onlineDetails,
			CertExpiration: certExp,
			SourceIP:       device.SourceIP,
			DeviceType:     device.DeviceType,
			AccountID:      device.AccountID,
			PairedAt:       device.PairedAt,
			LocationName:   device.LocationName,
			DeviceName:     device.DeviceName,
		},
	}, nil
}

// UpdateDeviceMetadata updates editable device metadata fields.
func (s *DeviceService) UpdateDeviceMetadata(deviceID, accountID string, meta *models.UpdateDeviceMetadata) error {
	device, err := s.store.GetDevice(deviceID)
	if err != nil {
		return err
	}
	if device == nil {
		return ErrNotFound
	}
	if device.AccountID != accountID {
		return ErrForbidden
	}
	name := meta.Name
	location := meta.Location
	rotInterval := meta.RotationIntervalDays
	if len(name) > 40 {
		name = name[:40]
	}
	if len(location) > 40 {
		location = location[:40]
	}
	if rotInterval < 30 {
		rotInterval = 30
	}
	if rotInterval > 5*365 {
		rotInterval = 5 * 365
	}
	return s.store.UpdateDeviceMetadata(deviceID, name, location, rotInterval)
}

// UpdateConfiguration validates and pushes desired config to a device.
func (s *DeviceService) UpdateConfiguration(deviceID, accountID string, cfgJSON json.RawMessage) error {
	device, err := s.store.GetDevice(deviceID)
	if err != nil {
		return err
	}
	if device == nil {
		return ErrNotFound
	}
	if device.AccountID != accountID {
		return ErrForbidden
	}
	if device.State == "DEPROVISIONED" {
		return fmt.Errorf("%w: device is deprovisioned", ErrConflict)
	}

	// Validate config against registration (if registration exists)
	if len(device.Registration) > 0 {
		if err := validateConfiguration(cfgJSON, device.Registration); err != nil {
			return fmt.Errorf("%w: %s", ErrBadRequest, err.Error())
		}
	}

	updateID, err := s.store.UpdateDeviceDesiredConfig(deviceID, cfgJSON)
	if err != nil {
		return err
	}

	// Build the outgoing payload with smart per-entity configurationId stamping.
	// configurationId is now a STRING (epoch seconds) — no float32 precision issues.
	// Each entity that changes gets its own independent time.Now().Unix() call.

	type prevChannel struct {
		ConfigurationId string
		State           string
		Settings        string
		Connection      string
	}
	prevDeviceConfigId := ""
	prevDeviceSettings := ""
	prevChannels := map[string]prevChannel{}

	if len(device.DesiredConfig) > 0 {
		var prev struct {
			ConfigurationId string          `json:"configurationId"`
			SimpleSettings  json.RawMessage `json:"standardSettings"`
			Channels        []struct {
				Id              string          `json:"id"`
				ConfigurationId string          `json:"configurationId"`
				State           string          `json:"state"`
				Settings        json.RawMessage `json:"settings"`
				Connection      json.RawMessage `json:"connection"`
			} `json:"channels"`
		}
		if json.Unmarshal(device.DesiredConfig, &prev) == nil {
			prevDeviceConfigId = prev.ConfigurationId
			prevDeviceSettings = canonicalJSON(prev.SimpleSettings)
			for _, ch := range prev.Channels {
				if ch.Id != "" {
					prevChannels[ch.Id] = prevChannel{
						ConfigurationId: ch.ConfigurationId,
						State:           ch.State,
						Settings:        canonicalJSON(ch.Settings),
						Connection:      canonicalJSON(ch.Connection),
					}
				}
			}
		}
	}

	var newCfg struct {
		SimpleSettings json.RawMessage `json:"standardSettings"`
		Channels       []struct {
			Id         string          `json:"id"`
			State      string          `json:"state"`
			Settings   json.RawMessage `json:"settings"`
			Connection json.RawMessage `json:"connection"`
		} `json:"channels"`
	}
	json.Unmarshal(cfgJSON, &newCfg)

	newDeviceSettings := canonicalJSON(newCfg.SimpleSettings)
	deviceConfigId := prevDeviceConfigId
	if prevDeviceConfigId == "" || newDeviceSettings != prevDeviceSettings {
		deviceConfigId = strconv.FormatInt(time.Now().UnixNano(), 10)
		log.Printf("[HOST UpdateConfig] device simpleSettings changed → configurationId=%s", deviceConfigId)
	} else {
		log.Printf("[HOST UpdateConfig] device simpleSettings unchanged → configurationId=%s", deviceConfigId)
	}

	wrapped := map[string]interface{}{}
	json.Unmarshal(cfgJSON, &wrapped)
	wrapped["updateId"] = updateID
	wrapped["configurationId"] = deviceConfigId

	if channels, ok := wrapped["channels"].([]interface{}); ok {
		for _, ch := range channels {
			chMap, ok := ch.(map[string]interface{})
			if !ok {
				continue
			}
			chID, _ := chMap["id"].(string)
			newChConfigId := strconv.FormatInt(time.Now().UnixNano(), 10) // default: first push
			for _, newCh := range newCfg.Channels {
				if newCh.Id != chID {
					continue
				}
				prev, hasPrev := prevChannels[chID]
				if hasPrev &&
					newCh.State == prev.State &&
					canonicalJSON(newCh.Settings) == prev.Settings &&
					canonicalJSON(newCh.Connection) == prev.Connection {
					newChConfigId = prev.ConfigurationId
					log.Printf("[HOST UpdateConfig] channel %s unchanged → configurationId=%s", chID, newChConfigId)
				} else {
					newChConfigId = strconv.FormatInt(time.Now().UnixNano(), 10)
					log.Printf("[HOST UpdateConfig] channel %s changed (state:%s→%s settings_changed=%v connection_changed=%v) → configurationId=%s",
						chID, prev.State, newCh.State,
						canonicalJSON(newCh.Settings) != prev.Settings,
						canonicalJSON(newCh.Connection) != prev.Connection,
						newChConfigId)
				}
				break
			}
			chMap["configurationId"] = newChConfigId
		}
	}

	payload, _ := json.Marshal(wrapped)
	topic := fmt.Sprintf("cdd/%s/config/update", deviceID)

	// Store the stamped config so future pushes can compare correctly (no counter bump).
	if err := s.store.StoreDeviceDesiredConfig(deviceID, payload); err != nil {
		log.Printf("[HOST UpdateConfig] failed to store stamped config: %v", err)
	}

	log.Printf("[HOST UpdateConfig] deviceID=%s state=%s online=%v topic=%s updateID=%d payloadLen=%d",
		deviceID, device.State, device.Online, topic, updateID, len(payload))

	if err := s.mqtt.Publish(topic, payload, true); err != nil { // retained — device picks up latest config on reconnect
		log.Printf("[HOST UpdateConfig] MQTT publish FAILED: %v", err)
		return err
	}
	log.Printf("[HOST UpdateConfig] MQTT publish succeeded for device=%s", deviceID)
	return nil
}

// Deprovision marks a device as deprovisioned (Phase 1).
func (s *DeviceService) Deprovision(deviceID, accountID string) error {
	device, err := s.store.GetDevice(deviceID)
	if err != nil {
		return err
	}
	if device == nil {
		return ErrNotFound
	}
	if device.AccountID != accountID {
		return ErrForbidden
	}
	if device.State == "DEPROVISIONED" {
		return nil // idempotent
	}

	if err := s.store.UpdateDeviceState(deviceID, "DEPROVISIONED", true); err != nil {
		return err
	}

	reason := tr12models.DEPROVISIONED
	t := time.Now().UTC()
	msg := models.DeprovisionRequest{Reason: &reason, Timestamp: t}
	payload, _ := json.Marshal(msg)
	topic := fmt.Sprintf("cdd/%s/deprovision", deviceID)
	return s.mqtt.Publish(topic, payload, true) // retained — offline device deprovisions itself on reconnect
}

// FullCleanup removes a device and all associated data (Phase 2 or device-initiated).
// Also clears any retained MQTT messages for the device's topics.
func (s *DeviceService) FullCleanup(deviceID string) error {
	if err := s.store.DeleteThumbnailsByDevice(deviceID); err != nil {
		return err
	}
	if err := s.store.DeleteLogsByDevice(deviceID); err != nil {
		return err
	}
	// Clear retained messages by publishing empty payloads (MQTT standard)
	for _, topic := range []string{
		fmt.Sprintf("cdd/%s/deprovision", deviceID),
		fmt.Sprintf("cdd/%s/config/update", deviceID),
		fmt.Sprintf("cdd/%s/certs/update", deviceID),
		fmt.Sprintf("cdd/%s/thumbnail/subscription", deviceID),
		fmt.Sprintf("cdd/%s/log/subscription", deviceID),
	} {
		_ = s.mqtt.Publish(topic, []byte{}, true)
	}
	return s.store.DeleteDevice(deviceID)
}

// RotateCredentials generates a new cert for a device and publishes via MQTT.
func (s *DeviceService) RotateCredentials(deviceID, accountID string) error {
	device, err := s.store.GetDevice(deviceID)
	if err != nil {
		return err
	}
	if device == nil {
		return ErrNotFound
	}
	if device.AccountID != accountID {
		return ErrForbidden
	}
	if device.State != "ACTIVE" {
		return fmt.Errorf("%w: device not active", ErrConflict)
	}

	log.Printf("[ROTATE] device=%s csrPEM length=%d currentCert length=%d", deviceID, len(device.CSRPEM), len(device.CurrentCertPEM))

	newCert, err := s.ca.SignCSR([]byte(device.CSRPEM), deviceID, s.cfg.CertExpiryDays)
	if err != nil {
		return err
	}

	sameAsCurrent := string(newCert) == device.CurrentCertPEM
	log.Printf("[ROTATE] newCert length=%d sameAsCurrent=%v", len(newCert), sameAsCurrent)

	now := time.Now().UTC()
	newExpires := now.Add(time.Duration(s.cfg.CertExpiryDays) * 24 * time.Hour).Format(time.RFC3339)

	if err := s.store.UpdateDeviceCerts(
		deviceID,
		string(newCert),
		device.CurrentCertPEM,
		newExpires,
		device.CertExpiresAt,
		now.Format(time.RFC3339),
	); err != nil {
		return err
	}

	mqttURI := fmt.Sprintf("tls://%s:%d", s.cfg.HostAddress, s.cfg.MQTTPort)
	rotate := tr12models.RotateCertificatesRequestContent{
		MqttUri:           mqttURI,
		DeviceCertificate: string(newCert),
		RegionName:        tr12models.PtrString("local"),
	}
	payload, _ := json.Marshal(rotate)
	topic := fmt.Sprintf("cdd/%s/certs/update", deviceID)
	log.Printf("[ROTATE] publishing to %s (retained=true) payload length=%d", topic, len(payload))
	return s.mqtt.Publish(topic, payload, true) // retained
}

// --- Helpers ---

func failPair(reason models.CreatePairingCodeFailureReason) *models.CreatePairingCodeResponseContent {
	return tr12models.NewCreatePairingCodeResponseContent(
		tr12models.FailureAsCreatePairingCodeResult(tr12models.NewFailure(
			models.CreatePairingCodeFailureData{Reason: reason},
		)),
	)
}

func buildHostSettings(deviceID string, pairingTimeout int) *models.HostSettings {
	hs := tr12models.NewHostSettings(
		"x-amzn-mqtt-ca",
		float32(pairingTimeout),
		1,
		30,
		fmt.Sprintf("cdd/%s/config/update", deviceID),
		fmt.Sprintf("cdd/%s/thumbnail/subscription", deviceID),
		fmt.Sprintf("cdd/%s/schema/report", deviceID),
		fmt.Sprintf("cdd/%s/registration/report", deviceID),
		fmt.Sprintf("cdd/%s/status/report", deviceID),
		fmt.Sprintf("cdd/%s/config/actual/report", deviceID),
		fmt.Sprintf("cdd/%s/certs/update", deviceID),
		fmt.Sprintf("cdd/%s/deprovision/ack", deviceID),
		fmt.Sprintf("cdd/%s/deprovision", deviceID),
		fmt.Sprintf("cdd/%s/log/subscription", deviceID),
	)
	log.Printf("[HOST buildHostSettings] deviceID=%s subUpdateTopic=%q", deviceID, hs.SubUpdateTopic)
	return hs
}

func generateDeviceID() string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	b := make([]byte, 21)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[n.Int64()]
	}
	return string(b)
}

func generatePairingCode() string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[n.Int64()]
	}
	return string(b)
}

func generateAccessCode() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 32)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		b[i] = chars[n.Int64()]
	}
	return string(b)
}

func formatOnlineDetails(d *models.Device) string {
	if !d.Online {
		if d.LastSeen == "" {
			return "offline"
		}
		return "offline since " + d.LastSeen
	}
	if d.LastSeen == "" {
		return "online"
	}
	t, err := time.Parse(time.RFC3339, d.LastSeen)
	if err != nil {
		return "online"
	}
	dur := time.Since(t)
	days := int(dur.Hours()) / 24
	hours := int(dur.Hours()) % 24
	mins := int(dur.Minutes()) % 60
	return fmt.Sprintf("online: %dd %dh %dm", days, hours, mins)
}

func formatCertExpiration(expiresAt string) string {
	if expiresAt == "" {
		return "unknown"
	}
	t, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return "unknown"
	}
	dur := time.Until(t)
	if dur < 0 {
		return "expired"
	}
	days := int(dur.Hours()) / 24
	hours := int(dur.Hours()) % 24
	mins := int(dur.Minutes()) % 60
	return fmt.Sprintf("%dd %dh %dm", days, hours, mins)
}

// validateConfiguration checks desired config against device registration.
func validateConfiguration(cfgJSON, regJSON json.RawMessage) error {
	var cfg struct {
		Channels       []struct {
			ID         string          `json:"id"`
			State      string          `json:"state,omitempty"`
			Settings   json.RawMessage `json:"settings,omitempty"`
			Connection json.RawMessage `json:"connection,omitempty"`
		} `json:"channels"`
		SimpleSettings []struct{ Key string `json:"key"` } `json:"standardSettings,omitempty"`
	}
	if err := json.Unmarshal(cfgJSON, &cfg); err != nil {
		return fmt.Errorf("invalid configuration JSON: %w", err)
	}

	var reg struct {
		Channels []struct {
			ID                  string                            `json:"id"`
			Name                string                            `json:"name"`
			SimpleSettings      []struct{ ID string `json:"id"` } `json:"standardSettings,omitempty"`
			Profiles            []struct{ ID string `json:"id"` } `json:"profiles,omitempty"`
			ConnectionProtocols []string                          `json:"connectionProtocols,omitempty"`
		} `json:"channels"`
		SimpleSettings []struct{ ID string `json:"id"` } `json:"standardSettings,omitempty"`
	}
	if err := json.Unmarshal(regJSON, &reg); err != nil {
		return fmt.Errorf("invalid registration JSON: %w", err)
	}

	// Validate device-level simpleSettings
	if len(cfg.SimpleSettings) > 0 {
		regDevSettings := map[string]bool{}
		for _, s := range reg.SimpleSettings {
			regDevSettings[s.ID] = true
		}
		for _, s := range cfg.SimpleSettings {
			if !regDevSettings[s.Key] {
				validKeys := make([]string, len(reg.SimpleSettings))
				for i, rs := range reg.SimpleSettings {
					validKeys[i] = rs.ID
				}
				return fmt.Errorf("unknown device-level setting key %q, valid: %v", s.Key, validKeys)
			}
		}
	}

	regChannels := map[string]int{}
	for i, ch := range reg.Channels {
		regChannels[ch.ID] = i
	}

	for _, cfgCh := range cfg.Channels {
		idx, ok := regChannels[cfgCh.ID]
		if !ok {
			validIDs := make([]string, len(reg.Channels))
			for i, ch := range reg.Channels {
				validIDs[i] = ch.ID
			}
			return fmt.Errorf("unknown channel ID %q, valid: %v", cfgCh.ID, validIDs)
		}

		if cfgCh.State != "" && cfgCh.State != "ACTIVE" && cfgCh.State != "IDLE" {
			return fmt.Errorf("invalid channel state %q for channel %q, must be ACTIVE or IDLE", cfgCh.State, cfgCh.ID)
		}

		regCh := reg.Channels[idx]

		if len(cfgCh.Settings) > 0 {
			var settings struct {
				SimpleSettings []struct{ Key string `json:"key"` } `json:"standardSettings,omitempty"`
				Profile        *struct{ ID string `json:"id"` }   `json:"profile,omitempty"`
			}
			json.Unmarshal(cfgCh.Settings, &settings)

			if settings.Profile != nil {
				validProfile := false
				for _, p := range regCh.Profiles {
					if p.ID == settings.Profile.ID {
						validProfile = true
						break
					}
				}
				if !validProfile {
					validIDs := make([]string, len(regCh.Profiles))
					for i, p := range regCh.Profiles {
						validIDs[i] = p.ID
					}
					return fmt.Errorf("unknown profile ID %q for channel %q, valid: %v", settings.Profile.ID, cfgCh.ID, validIDs)
				}
			}
			if len(settings.SimpleSettings) > 0 {
				regSettings := map[string]bool{}
				for _, s := range regCh.SimpleSettings {
					regSettings[s.ID] = true
				}
				for _, s := range settings.SimpleSettings {
					if !regSettings[s.Key] {
						return fmt.Errorf("unknown setting key %q for channel %q", s.Key, cfgCh.ID)
					}
				}
			}
		}

		// Validate connection transport protocol
		if len(cfgCh.Connection) > 0 {
			if err := validateConnection(cfgCh.Connection, regCh.ConnectionProtocols, cfgCh.ID); err != nil {
				return err
			}
		}
	}
	return nil
}

// protocolFieldMap defines the valid JSON field names for each transport protocol type.
// Required fields are marked true, optional fields are marked false.
var protocolFieldMap = map[string]map[string]bool{
	"srtCaller": {
		"address":                     true,
		"port":                        true,
		"minimumLatencyMilliseconds":  false,
		"streamId":                    false,
		"encryption":                  false,
	},
	"srtListener": {
		"port":                        true,
		"minimumLatencyMilliseconds":  false,
		"streamId":                    false,
		"encryption":                  false,
		"interface":                   false,
	},
	"ristCaller": {
		"address":                     true,
		"port":                        true,
		"minimumLatencyMilliseconds":  false,
		"streamId":                    false,
		"encryption":                  false,
	},
	"ristListener": {
		"port":                        true,
		"minimumLatencyMilliseconds":  false,
		"streamId":                    false,
		"encryption":                  false,
		"interface":                   false,
	},
	"zixiPush": {
		"streamId":                    false,
		"address":                     true,
		"port":                        true,
		"minimumLatencyMilliseconds":  false,
		"encryption":                  false,
	},
	"zixiPull": {
		"streamId":                    false,
		"address":                     true,
		"port":                        true,
		"minimumLatencyMilliseconds":  false,
		"encryption":                  false,
	},
	"rtp": {
		"address":                     true,
		"port":                        true,
		"sourceAddressFilter":         false,
		"rtpPayloadType":              false,
		"fecConfig":                   false,
	},
}

// protocolKeyToRegistration maps the JSON union key (e.g. "srtCaller") to the
// TransportProtocolName enum value in the registration (e.g. "SRT_CALLER").
var protocolKeyToRegistration = map[string]string{
	"srtCaller":    "SRT_CALLER",
	"srtListener":  "SRT_LISTENER",
	"ristCaller":   "RIST_CALLER",
	"ristListener": "RIST_LISTENER",
	"zixiPush":     "ZIXI_PUSH",
	"zixiPull":     "ZIXI_PULL",
	"rtp":          "RTP",
}

func validateConnection(connJSON json.RawMessage, registeredProtocols []string, channelID string) error {
	var conn struct {
		TransportProtocol json.RawMessage `json:"transportProtocol,omitempty"`
	}
	if err := json.Unmarshal(connJSON, &conn); err != nil {
		return fmt.Errorf("invalid connection JSON for channel %q: %w", channelID, err)
	}
	if len(conn.TransportProtocol) == 0 {
		return nil
	}

	// Parse the transport protocol union — it's a JSON object with a single key
	// like {"srtCaller": {...}} or {"ristListener": {...}}
	var protoMap map[string]json.RawMessage
	if err := json.Unmarshal(conn.TransportProtocol, &protoMap); err != nil {
		return fmt.Errorf("invalid transportProtocol JSON for channel %q: %w", channelID, err)
	}

	for protoKey, protoBody := range protoMap {
		// Check that the protocol type is known
		regValue, knownProto := protocolKeyToRegistration[protoKey]
		if !knownProto {
			validKeys := make([]string, 0, len(protocolKeyToRegistration))
			for k := range protocolKeyToRegistration {
				validKeys = append(validKeys, k)
			}
			return fmt.Errorf("unknown transport protocol type %q for channel %q, valid: %v", protoKey, channelID, validKeys)
		}

		// Check that the protocol type is registered for this channel
		if len(registeredProtocols) > 0 {
			found := false
			for _, rp := range registeredProtocols {
				if rp == regValue {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("transport protocol %q is not supported by channel %q, registered protocols: %v", protoKey, channelID, registeredProtocols)
			}
		}

		// Validate field names against the Smithy schema
		validFields, ok := protocolFieldMap[protoKey]
		if !ok {
			continue // shouldn't happen since we checked knownProto above
		}

		var fieldMap map[string]interface{}
		if err := json.Unmarshal(protoBody, &fieldMap); err != nil {
			return fmt.Errorf("invalid %s body for channel %q: %w", protoKey, channelID, err)
		}

		for fieldName := range fieldMap {
			if _, valid := validFields[fieldName]; !valid {
				validNames := make([]string, 0, len(validFields))
				for k := range validFields {
					validNames = append(validNames, k)
				}
				return fmt.Errorf("unknown field %q in %s for channel %q, valid fields: %v", fieldName, protoKey, channelID, validNames)
			}
		}

		// Check required fields are present
		for fieldName, required := range validFields {
			if required {
				if _, present := fieldMap[fieldName]; !present {
					return fmt.Errorf("missing required field %q in %s for channel %q", fieldName, protoKey, channelID)
				}
			}
		}
	}

	return nil
}

// canonicalJSON normalizes a JSON value by round-tripping through interface{}.
// This ensures consistent key ordering and number representation for comparison.
// Returns "" for nil or invalid input.
func canonicalJSON(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var v interface{}
	if err := json.Unmarshal(raw, &v); err != nil {
		return string(raw)
	}
	b, err := json.Marshal(v)
	if err != nil {
		return string(raw)
	}
	return string(b)
}
