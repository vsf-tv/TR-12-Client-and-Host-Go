package service

import (
	"crypto/rand"
	"encoding/json"
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
func (s *DeviceService) Pair(req models.PairRequestContent) (*models.PairResponseContent, error) {
	// Validate host ID
	if req.HostId != s.cfg.ServiceID {
		return failPair(models.PairFailureHostIDMismatch), nil
	}
	// Validate version
	if req.Version == "" {
		return failPair(models.PairFailureVersionNotSupported), nil
	}
	// Validate device type
	validType := false
	for _, dt := range []string{"SOURCE", "DESTINATION", "BOTH"} {
		if strings.EqualFold(req.DeviceType, dt) {
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
	certPEM, err := s.ca.SignCSR([]byte(req.Csr), deviceID, s.cfg.CertExpiryDays)
	if err != nil {
		return nil, fmt.Errorf("sign CSR: %w", err)
	}

	now := time.Now().UTC()
	device := &models.Device{
		DeviceID:         deviceID,
		AccountID:        "",
		DeviceType:       strings.ToUpper(req.DeviceType),
		State:            "PAIRING",
		PairedAt:         now.Format(time.RFC3339),
		CurrentCertPEM:   string(certPEM),
		CertExpiresAt:    now.Add(time.Duration(s.cfg.CertExpiryDays) * 24 * time.Hour).Format(time.RFC3339),
		CSRPEM:           req.Csr,
		PairingCode:      pairingCode,
		AccessCode:       accessCode,
		PairingExpiresAt: now.Add(time.Duration(s.cfg.PairingTimeout) * time.Second).Format(time.RFC3339),
	}
	if err := s.store.InsertDevice(device); err != nil {
		return nil, fmt.Errorf("insert device: %w", err)
	}

	successData := models.PairSuccessData{
		DeviceId:              deviceID,
		PairingCode:           pairingCode,
		AccessCode:            accessCode,
		PairingTimeoutSeconds: float32(s.cfg.PairingTimeout),
	}
	return tr12models.NewPairResponseContent(
		tr12models.SuccessAsPairResult(tr12models.NewSuccess(successData)),
	), nil
}

// Authenticate handles device authentication polling.
func (s *DeviceService) Authenticate(req models.AuthenticateRequestContent) (*models.AuthenticateResponseContent, error) {
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
		resp := tr12models.NewAuthenticateResponseContent(models.AuthStatusSTANDBY)
		return resp, nil
	}

	// Device is claimed — return full auth response
	mqttURI := fmt.Sprintf("tls://%s:%d", s.cfg.HostAddress, s.cfg.MQTTPort)
	hs := buildHostSettings(device.DeviceID, s.cfg.PairingTimeout)
	resp := tr12models.NewAuthenticateResponseContent(models.AuthStatusCLAIMED)
	resp.SetCaCert(string(s.ca.CACertPEM))
	resp.SetDeviceCert(device.CurrentCertPEM)
	resp.SetMqttUri(mqttURI)
	resp.SetRegion("local")
	resp.SetHostSettings(*hs)
	return resp, nil
}

// Claim associates a device with an account.
func (s *DeviceService) Claim(pairingCode, accountID string, expirationDays int) error {
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
	regExpires := time.Now().UTC().Add(time.Duration(expirationDays) * 24 * time.Hour).Format(time.RFC3339)
	return s.store.ClaimDevice(device.DeviceID, accountID, regExpires)
}

// ListDevices returns all devices for an account.
func (s *DeviceService) ListDevices(accountID string) ([]models.DeviceSummary, error) {
	devices, err := s.store.ListDevicesByAccount(accountID)
	if err != nil {
		return nil, err
	}
	summaries := make([]models.DeviceSummary, 0, len(devices))
	for _, d := range devices {
		summaries = append(summaries, models.DeviceSummary{
			DeviceID:      d.DeviceID,
			Message:       "",
			Errors:        []string{},
			OnlineDetails: formatOnlineDetails(d),
			Online:        d.Online,
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
		},
	}, nil
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

	// Wrap config with updateId and publish via MQTT
	wrapped := map[string]interface{}{}
	json.Unmarshal(cfgJSON, &wrapped)
	wrapped["updateId"] = updateID
	payload, _ := json.Marshal(wrapped)
	topic := fmt.Sprintf("cdd/%s/config/update", deviceID)

	log.Printf("[HOST UpdateConfig] deviceID=%s state=%s online=%v topic=%s updateID=%d payloadLen=%d",
		deviceID, device.State, device.Online, topic, updateID, len(payload))

	if err := s.mqtt.Publish(topic, payload, false); err != nil {
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
	t := float32(time.Now().Unix())
	msg := models.DeprovisionDeviceRequestContent{Reason: &reason, Time: &t}
	payload, _ := json.Marshal(msg)
	topic := fmt.Sprintf("cdd/%s/deprovision", deviceID)
	return s.mqtt.Publish(topic, payload, false)
}

// FullCleanup removes a device and all associated data (Phase 2 or device-initiated).
func (s *DeviceService) FullCleanup(deviceID string) error {
	if err := s.store.DeleteThumbnailsByDevice(deviceID); err != nil {
		return err
	}
	if err := s.store.DeleteLogsByDevice(deviceID); err != nil {
		return err
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
		MqttUri:    mqttURI,
		DeviceCert: string(newCert),
		Region:     "local",
	}
	payload, _ := json.Marshal(rotate)
	topic := fmt.Sprintf("cdd/%s/certs/update", deviceID)
	log.Printf("[ROTATE] publishing to %s (retained=true) payload length=%d", topic, len(payload))
	return s.mqtt.Publish(topic, payload, true) // retained
}

// --- Helpers ---

func failPair(reason models.PairFailureReason) *models.PairResponseContent {
	return tr12models.NewPairResponseContent(
		tr12models.FailureAsPairResult(tr12models.NewFailure(
			models.PairFailureData{Reason: reason},
		)),
	)
}

func buildHostSettings(deviceID string, pairingTimeout int) *models.HostSettings {
	hs := tr12models.NewHostSettings(
		"mqtt",
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
		fmt.Sprintf("cdd/%s/deprovision", deviceID),
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
		SimpleSettings []struct{ Key string `json:"key"` } `json:"simpleSettings,omitempty"`
	}
	if err := json.Unmarshal(cfgJSON, &cfg); err != nil {
		return fmt.Errorf("invalid configuration JSON: %w", err)
	}

	var reg struct {
		Channels []struct {
			ID                  string                            `json:"id"`
			Name                string                            `json:"name"`
			SimpleSettings      []struct{ ID string `json:"id"` } `json:"simpleSettings,omitempty"`
			Profiles            []struct{ ID string `json:"id"` } `json:"profiles,omitempty"`
			ConnectionProtocols []string                          `json:"connectionProtocols,omitempty"`
		} `json:"channels"`
		SimpleSettings []struct{ ID string `json:"id"` } `json:"simpleSettings,omitempty"`
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
				SimpleSettings []struct{ Key string `json:"key"` } `json:"simpleSettings,omitempty"`
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
		"ip":                          true,
		"port":                        true,
		"minimumLatencyMilliseconds":  true,
		"streamId":                    false,
		"encryption":                  false,
	},
	"srtListener": {
		"port":                        true,
		"minimumLatencyMilliseconds":  true,
		"streamId":                    false,
		"encryption":                  false,
		"interface":                   false,
	},
	"ristCaller": {
		"ip":                          true,
		"port":                        true,
		"minimumLatencyMilliseconds":  true,
		"streamId":                    false,
		"encryption":                  false,
	},
	"ristListener": {
		"port":                        true,
		"minimumLatencyMilliseconds":  true,
		"streamId":                    false,
		"encryption":                  false,
		"interface":                   false,
	},
	"zixiCaller": {
		"streamId":                    true,
		"ip":                          true,
		"port":                        true,
		"minimumLatencyMilliseconds":  true,
		"encryption":                  false,
	},
	"zixiListener": {
		"streamId":                    true,
		"port":                        true,
		"minimumLatencyMilliseconds":  true,
		"encryption":                  false,
		"interface":                   false,
	},
}

// protocolKeyToRegistration maps the JSON union key (e.g. "srtCaller") to the
// SupportedProtocol enum value in the registration (e.g. "SRT_CALLER").
var protocolKeyToRegistration = map[string]string{
	"srtCaller":     "SRT_CALLER",
	"srtListener":   "SRT_LISTENER",
	"ristCaller":    "RIST_CALLER",
	"ristListener":  "RIST_LISTENER",
	"zixiCaller":    "ZIXI_CALLER",
	"zixiListener":  "ZIXI_LISTENER",
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
