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
// Returns (response, pairingErr, internalErr) — pairingErr is a 400, internalErr is a 500.
func (s *DeviceService) Pair(req models.CreatePairingCodeRequestContent) (*models.CreatePairingCodeResponseContent, *models.CreatePairingCodeFailureReason, error) {
	// Validate host ID
	if req.HostId != s.cfg.ServiceID {
		r := models.PairFailureHostIDMismatch
		return nil, &r, nil
	}
	// Validate version
	if req.Version.GetVersion() == "" {
		r := models.PairFailureVersionNotSupported
		return nil, &r, nil
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
		r := models.PairFailureDeviceTypeNotSupported
		return nil, &r, nil
	}

	deviceID := generateDeviceID()
	pairingCode := generatePairingCode()
	accessCode := generateAccessCode()

	// Sign CSR
	certPEM, err := s.ca.SignCSR([]byte(req.CertificateSigningRequest), deviceID, s.cfg.CertExpiryDays)
	if err != nil {
		return nil, nil, fmt.Errorf("sign CSR: %w", err)
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
		return nil, nil, fmt.Errorf("insert device: %w", err)
	}

	resp := tr12models.NewCreatePairingCodeResponseContent(
		deviceID, pairingCode, accessCode, float32(s.cfg.PairingTimeout),
	)
	return resp, nil, nil
}

// Authenticate handles device authentication polling.
func (s *DeviceService) Authenticate(req models.AuthenticatePairingCodeRequestContent) (*models.AuthenticatePairingCodeResponseContent, error) {
	device, err := s.store.GetDevice(req.DeviceId)
	if err != nil {
		return nil, err
	}
	if device == nil {
		// Return STANDBY rather than 404 — the spec does not define an error case for
		// AuthenticatePairingCode. Returning not-found would let callers enumerate valid device IDs.
		resp := tr12models.NewAuthenticatePairingCodeResponseContent(models.AuthStatusSTANDBY)
		return resp, nil
	}
	if device.PairingCode != req.PairingCode || device.AccessCode != req.AccessCode {
		// Return STANDBY rather than an error — the spec does not define an error case for
		// AuthenticatePairingCode. Returning STANDBY for bad credentials prevents callers
		// from distinguishing a valid-but-unclaimed device from garbage credentials.
		resp := tr12models.NewAuthenticatePairingCodeResponseContent(models.AuthStatusSTANDBY)
		return resp, nil
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
// credentialExpirationDays, if > 0, re-signs the device certificate with that
// lifetime so the operator can set a per-device cert duration at claim time.
// If 0 or omitted the cert signed at pair time (service default, typically 30 days) is kept.
func (s *DeviceService) Claim(pairingCode, accountID string, expirationDays, credentialExpirationDays int, locationName, deviceName string, rotationIntervalDays int) error {
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

	// Re-sign the device certificate if a custom credential lifetime was requested.
	if credentialExpirationDays > 0 && device.CSRPEM != "" {
		newCert, err := s.ca.SignCSR([]byte(device.CSRPEM), device.DeviceID, credentialExpirationDays)
		if err != nil {
			return fmt.Errorf("re-sign cert for claim: %w", err)
		}
		newExpires := time.Now().UTC().Add(time.Duration(credentialExpirationDays) * 24 * time.Hour).Format(time.RFC3339)
		if err := s.store.SetInitialCert(device.DeviceID, string(newCert), newExpires); err != nil {
			return fmt.Errorf("update cert at claim: %w", err)
		}
	}

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
	if device.State == "DEPROVISIONED" {
		return nil, ErrNotFound // hide deprovisioned devices from get as well as list
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

// UpdateChannelConfig updates a single channel's content, unconditionally bumps
// that channel's version, and preserves all other channels' existing versions.
// The console sends only the fields for the one channel being updated.
func (s *DeviceService) UpdateChannelConfig(deviceID, accountID, channelID string, channelCfg json.RawMessage) error {
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

	// Validate that channelID exists in the registration.
	if len(device.Registration) > 0 {
		if err := validateChannelExistsInRegistration(channelID, device.Registration); err != nil {
			return fmt.Errorf("%w: %s", ErrNotFound, err.Error())
		}
	}

	// Load the stored full config and merge the updated channel into it.
	merged, _, err := s.mergeChannelUpdate(device.DesiredConfig, channelID, channelCfg)
	if err != nil {
		return fmt.Errorf("merge channel update: %w", err)
	}

	// Only publish the single updated channel in MQTT — not the full config.
	// This prevents the device from re-processing channels that weren't part of this update.
	// UPDATE: Publish the full config. The SDK replaces its cached config with whatever
	// arrives on MQTT (retained message), so partial updates cause other channels to
	// disappear. The client-side ApplicationLoop is version-gated and will skip
	// channels whose version hasn't changed — so publishing all channels is safe.
	return s.persistAndPublish(device, merged, nil)
}

// UpdateDeviceSettings updates only the device-level standardSettings, unconditionally
// bumps the device version, and preserves all channel versions unchanged.
func (s *DeviceService) UpdateDeviceSettings(deviceID, accountID string, settingsJSON json.RawMessage) error {
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

	// Merge the new standardSettings into the stored config, preserving all channel versions.
	merged, err := s.mergeDeviceSettingsUpdate(device.DesiredConfig, settingsJSON)
	if err != nil {
		return fmt.Errorf("merge device settings: %w", err)
	}

	// Device settings update publishes all channels — the device needs the full config
	// to apply the new standardSettings alongside the current channel state.
	return s.persistAndPublish(device, merged, nil)
}

// mergeChannelUpdate loads the stored full config and replaces the target channel's content,
// assigning a fresh version to that channel only. All other channel versions are preserved.
// Returns the merged full config (for DB storage) and the single updated channel entry
// (for MQTT — so the device only sees the one channel that changed).
func (s *DeviceService) mergeChannelUpdate(storedConfig json.RawMessage, channelID string, channelCfg json.RawMessage) (merged map[string]interface{}, mqttChannel interface{}, err error) {
	merged = s.loadOrInitFullConfig(storedConfig)

	// Parse incoming channel update.
	var incoming map[string]interface{}
	if err = json.Unmarshal(channelCfg, &incoming); err != nil {
		return nil, nil, fmt.Errorf("invalid channel config JSON: %w", err)
	}

	newVersion := strconv.FormatInt(time.Now().UnixNano(), 10)

	channels, _ := merged["channels"].([]interface{})
	found := false
	for i, ch := range channels {
		chMap, ok := ch.(map[string]interface{})
		if !ok {
			continue
		}
		if chMap["id"] == channelID {
			// Merge incoming fields — preserve version key explicitly.
			for k, v := range incoming {
				chMap[k] = v
			}
			chMap["id"] = channelID
			chMap["version"] = newVersion
			channels[i] = chMap
			mqttChannel = chMap
			found = true
			break
		}
	}

	if !found {
		// Channel not in stored config yet — add it with the incoming data.
		newCh := map[string]interface{}{"id": channelID, "version": newVersion}
		for k, v := range incoming {
			newCh[k] = v
		}
		channels = append(channels, newCh)
		mqttChannel = newCh
	}

	merged["channels"] = channels
	// The device model requires a top-level version on DesiredDeviceConfiguration.
	// Bump it alongside the channel version so the MQTT payload is always valid.
	merged["version"] = newVersion
	log.Printf("[HOST UpdateChannelConfig] channel %s → version=%s", channelID, newVersion)
	return merged, mqttChannel, nil
}

// mergeDeviceSettingsUpdate replaces standardSettings in the stored config and bumps
// the device-level version. Channel versions are left untouched.
func (s *DeviceService) mergeDeviceSettingsUpdate(storedConfig json.RawMessage, settingsJSON json.RawMessage) (map[string]interface{}, error) {
	merged := s.loadOrInitFullConfig(storedConfig)

	var incoming map[string]interface{}
	if err := json.Unmarshal(settingsJSON, &incoming); err != nil {
		return nil, fmt.Errorf("invalid settings JSON: %w", err)
	}

	if ss, ok := incoming["standardSettings"]; ok {
		merged["standardSettings"] = ss
	}

	newVersion := strconv.FormatInt(time.Now().UnixNano(), 10)
	merged["version"] = newVersion
	log.Printf("[HOST UpdateDeviceSettings] device → version=%s", newVersion)
	return merged, nil
}

// loadOrInitFullConfig parses the stored desired config into a mutable map.
// If no stored config exists, returns an empty skeleton with a channels slice.
func (s *DeviceService) loadOrInitFullConfig(storedConfig json.RawMessage) map[string]interface{} {
	result := map[string]interface{}{"channels": []interface{}{}}
	if len(storedConfig) > 0 {
		json.Unmarshal(storedConfig, &result)
		// Ensure channels key is always a slice.
		if result["channels"] == nil {
			result["channels"] = []interface{}{}
		}
	}
	return result
}

// persistAndPublish saves the merged config and publishes it via MQTT.
// mqttChannels: if non-nil, only these channels are included in the MQTT payload.
//               The DB always stores the full merged config regardless.
// This lets per-channel updates avoid re-triggering unchanged channels on the device.
func (s *DeviceService) persistAndPublish(device *models.Device, merged map[string]interface{}, mqttChannels []interface{}) error {
	// Remove envelope fields before storing.
	delete(merged, "updateId")

	cfgPayload, err := json.Marshal(merged)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}

	updateID, err := s.store.UpdateDeviceDesiredConfig(device.DeviceID, cfgPayload)
	if err != nil {
		return err
	}

	// Build the MQTT payload. If mqttChannels is set, use only those channels so the
	// device does not re-process channels that weren't part of this update.
	mqttConfig := merged
	if mqttChannels != nil {
		mqttConfig = make(map[string]interface{}, len(merged))
		for k, v := range merged {
			mqttConfig[k] = v
		}
		mqttConfig["channels"] = mqttChannels
	}

	topic := fmt.Sprintf("cdd/%s/config/update", device.DeviceID)
	envelopePayload, _ := json.Marshal(map[string]interface{}{
		"updateId":                   updateID,
		"desiredDeviceConfiguration": mqttConfig,
	})

	log.Printf("[HOST publish] deviceID=%s topic=%s updateID=%d payloadLen=%d mqttChannelCount=%d",
		device.DeviceID, topic, updateID, len(envelopePayload), channelCount(mqttConfig))

	if err := s.mqtt.Publish(topic, envelopePayload, true); err != nil {
		log.Printf("[HOST publish] MQTT publish FAILED: %v", err)
		return err
	}
	return nil
}

func channelCount(cfg map[string]interface{}) int {
	if chs, ok := cfg["channels"].([]interface{}); ok {
		return len(chs)
	}
	return 0
}

// UpdateConfiguration is the legacy full-device update path (kept for backward compat).
// It unconditionally bumps every channel's version — use UpdateChannelConfig for
// per-channel operations.
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

	if len(device.Registration) > 0 {
		if err := validateConfiguration(cfgJSON, device.Registration); err != nil {
			return fmt.Errorf("%w: %s", ErrBadRequest, err.Error())
		}
	}

	var full map[string]interface{}
	json.Unmarshal(cfgJSON, &full)

	// Bump device version.
	full["version"] = strconv.FormatInt(time.Now().UnixNano(), 10)

	// Unconditionally bump every channel version.
	if channels, ok := full["channels"].([]interface{}); ok {
		for _, ch := range channels {
			if chMap, ok := ch.(map[string]interface{}); ok {
				chMap["version"] = strconv.FormatInt(time.Now().UnixNano(), 10)
			}
		}
	}

	return s.persistAndPublish(device, full, nil)
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

	reason := tr12models.DEPROVISIONREASON_DEPROVISIONED
	t := time.Now().UTC()
	deprov := tr12models.DeviceSubscribesToDeprovisionResponseContent{
		Reason:    &reason,
		Timestamp: t,
	}
	payload, _ := json.Marshal(deprov)
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
	rotate := tr12models.DeviceSubscribesToCertificateRotationResponseContent{
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

func buildHostSettings(deviceID string, pairingTimeout int) *models.HostSettings {
	hs := tr12models.NewHostSettings(
		"x-amzn-mqtt-ca",
		float32(pairingTimeout),
		1,
		30,
		fmt.Sprintf("cdd/%s/config/update", deviceID),
		fmt.Sprintf("cdd/%s/thumbnail/subscription", deviceID),
		fmt.Sprintf("cdd/%s/registration/report", deviceID),
		fmt.Sprintf("cdd/%s/status/report", deviceID),
		fmt.Sprintf("cdd/%s/config/actual/report", deviceID),
		fmt.Sprintf("cdd/%s/certs/update", deviceID),
		fmt.Sprintf("cdd/%s/deprovision/ack", deviceID),
		fmt.Sprintf("cdd/%s/deprovision", deviceID),
		fmt.Sprintf("cdd/%s/log/subscription", deviceID),
	)
	log.Printf("[HOST buildHostSettings] deviceID=%s desiredConfigTopic=%q", deviceID, hs.DeviceSubscribesToDesiredConfigurationTopic)
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
			ID              string          `json:"id"`
			State           string          `json:"state,omitempty"`
			ChannelSettings json.RawMessage `json:"channelSettings,omitempty"`
			Protocol        json.RawMessage `json:"protocol,omitempty"`
		} `json:"channels"`
		StandardSettings []struct{ Id string `json:"id"` } `json:"standardSettings,omitempty"`
	}
	if err := json.Unmarshal(cfgJSON, &cfg); err != nil {
		return fmt.Errorf("invalid configuration JSON: %w", err)
	}

	var reg struct {
		ChannelTemplates []struct {
			ID              string                            `json:"id"`
			Name            string                            `json:"name"`
			ChannelSettings []struct{ ID string `json:"id"` } `json:"settings,omitempty"`
			Profiles        []struct{ ID string `json:"id"` } `json:"profiles,omitempty"`
			Protocols       []string                          `json:"protocols,omitempty"`
		} `json:"channelTemplates"`
		ChannelAssignments []struct {
			ChannelID  string `json:"channelId"`
			TemplateID string `json:"templateId"`
		} `json:"channelAssignments"`
		DeviceRegistrationSettings []struct{ ID string `json:"id"` } `json:"settings,omitempty"`
	}
	if err := json.Unmarshal(regJSON, &reg); err != nil {
		return fmt.Errorf("invalid registration JSON: %w", err)
	}

	// Build template lookup and resolved channel map from assignments
	type resolvedRegChannel struct {
		ID       string
		Settings []struct{ ID string `json:"id"` }
		Profiles []struct{ ID string `json:"id"` }
		Protocols []string
	}
	templateByID := make(map[string]int, len(reg.ChannelTemplates))
	for i, tmpl := range reg.ChannelTemplates {
		templateByID[tmpl.ID] = i
	}
	regChannels := make(map[string]resolvedRegChannel, len(reg.ChannelAssignments))
	for _, assignment := range reg.ChannelAssignments {
		if idx, ok := templateByID[assignment.TemplateID]; ok {
			tmpl := reg.ChannelTemplates[idx]
			regChannels[assignment.ChannelID] = resolvedRegChannel{
				ID:        assignment.ChannelID,
				Settings:  tmpl.ChannelSettings,
				Profiles:  tmpl.Profiles,
				Protocols: tmpl.Protocols,
			}
		}
	}

	// Validate device-level standardSettings
	if len(cfg.StandardSettings) > 0 {
		regDevSettings := map[string]bool{}
		for _, s := range reg.DeviceRegistrationSettings {
			regDevSettings[s.ID] = true
		}
		for _, s := range cfg.StandardSettings {
			if !regDevSettings[s.Id] {
				validKeys := make([]string, len(reg.DeviceRegistrationSettings))
				for i, rs := range reg.DeviceRegistrationSettings {
					validKeys[i] = rs.ID
				}
				return fmt.Errorf("unknown device-level setting key %q, valid: %v", s.Id, validKeys)
			}
		}
	}

	for _, cfgCh := range cfg.Channels {
		regCh, ok := regChannels[cfgCh.ID]
		if !ok {
			validIDs := make([]string, 0, len(regChannels))
			for id := range regChannels {
				validIDs = append(validIDs, id)
			}
			return fmt.Errorf("unknown channel ID %q, valid: %v", cfgCh.ID, validIDs)
		}

		if cfgCh.State != "" && cfgCh.State != "ACTIVE" && cfgCh.State != "IDLE" {
			return fmt.Errorf("invalid channel state %q for channel %q, must be ACTIVE or IDLE", cfgCh.State, cfgCh.ID)
		}

		if len(cfgCh.ChannelSettings) > 0 {
			// ChannelSettings is a union: either {"standardSettings": [...]} or {"profile": {"id": "..."}}
			var settings struct {
				StandardSettings []struct{ Id string `json:"id"` } `json:"standardSettings,omitempty"`
				Profile          *struct{ ID string `json:"id"` }  `json:"profile,omitempty"`
			}
			json.Unmarshal(cfgCh.ChannelSettings, &settings)

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
			if len(settings.StandardSettings) > 0 {
				regSettings := map[string]bool{}
				for _, s := range regCh.Settings {
					regSettings[s.ID] = true
				}
				for _, s := range settings.StandardSettings {
					if !regSettings[s.Id] {
						return fmt.Errorf("unknown setting key %q for channel %q", s.Id, cfgCh.ID)
					}
				}
			}
		}

		// Validate transport protocol
		if len(cfgCh.Protocol) > 0 {
			if err := validateProtocol(cfgCh.Protocol, regCh.Protocols, cfgCh.ID); err != nil {
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
		"address":                    true,
		"port":                       true,
		"minimumLatencyMilliseconds": false,
		"streamId":                   false,
		"encryption":                 false,
	},
	"srtListener": {
		"port":                       true,
		"minimumLatencyMilliseconds": false,
		"streamId":                   false,
		"encryption":                 false,
		"interface":                  false,
	},
	"ristSimpleCaller": {
		"address":                    true,
		"port":                       true,
		"minimumLatencyMilliseconds": false,
		"encryption":                 false,
	},
	"ristSimpleListener": {
		"port":                       true,
		"minimumLatencyMilliseconds": false,
		"encryption":                 false,
		"interface":                  false,
	},
	"zixiPushSender": {
		"address":                    true,
		"streamId":                   false,
		"port":                       false,
		"maximumLatencyMilliseconds": false,
		"encryption":                 false,
	},
	"zixiPushReceiver": {
		"address":                    true,
		"streamId":                   false,
		"port":                       false,
		"maximumLatencyMilliseconds": false,
		"encryption":                 false,
	},
	"zixiPullSender": {
		"streamId":                   true,
		"port":                       true,
		"maximumLatencyMilliseconds": false,
		"encryption":                 false,
		"interface":                  false,
	},
	"zixiPullReceiver": {
		"streamId":                   true,
		"address":                    true,
		"port":                       false,
		"maximumLatencyMilliseconds": false,
		"encryption":                 false,
	},
	"rtp": {
		"address":             true,
		"port":                true,
		"sourceAddressFilter": false,
	},
}

// protocolKeyToRegistration maps the JSON union key (e.g. "srtCaller") to the
// TransportProtocolName enum value in the registration (e.g. "SRT_CALLER").
var protocolKeyToRegistration = map[string]string{
	"srtCaller":          "SRT_CALLER",
	"srtListener":        "SRT_LISTENER",
	"ristSimpleCaller":   "RIST_SIMPLE_SENDER",
	"ristSimpleListener": "RIST_SIMPLE_RECEIVER",
	"zixiPushSender":     "ZIXI_PUSH_SENDER",
	"zixiPushReceiver":   "ZIXI_PUSH_RECEIVER",
	"zixiPullSender":     "ZIXI_PULL_SENDER",
	"zixiPullReceiver":   "ZIXI_PULL_RECEIVER",
	"rtp":                "RTP",
}

func validateProtocol(protoJSON json.RawMessage, registeredProtocols []string, channelID string) error {
	// Parse the transport protocol union — it's a JSON object with a single key
	// like {"srtCaller": {...}} or {"ristListener": {...}}
	var protoMap map[string]json.RawMessage
	if err := json.Unmarshal(protoJSON, &protoMap); err != nil {
		return fmt.Errorf("invalid protocol JSON for channel %q: %w", channelID, err)
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

// validateChannelExistsInRegistration checks that channelID is known in the device registration.
func validateChannelExistsInRegistration(channelID string, regJSON json.RawMessage) error {
	var reg struct {
		ChannelAssignments []struct {
			ChannelID string `json:"channelId"`
		} `json:"channelAssignments"`
		// Legacy flat format
		Channels []struct {
			ID string `json:"id"`
		} `json:"channels"`
	}
	if err := json.Unmarshal(regJSON, &reg); err != nil {
		return nil // can't parse registration, skip validation
	}
	for _, a := range reg.ChannelAssignments {
		if a.ChannelID == channelID {
			return nil
		}
	}
	for _, ch := range reg.Channels {
		if ch.ID == channelID {
			return nil
		}
	}
	return fmt.Errorf("channel %q not found in device registration", channelID)
}
