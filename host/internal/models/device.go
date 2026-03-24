package models

import "encoding/json"

// Device represents a device record in the registry.
type Device struct {
	DeviceID             string           `json:"device_id"`
	AccountID            string           `json:"account_id"`
	DeviceType           string           `json:"device_type"`
	State                string           `json:"state"` // PAIRING, ACTIVE, DEPROVISIONED
	Registration         json.RawMessage  `json:"registration,omitempty"`
	DesiredConfig        json.RawMessage  `json:"desired_config,omitempty"`
	ActualConfig         json.RawMessage  `json:"actual_config,omitempty"`
	Status               json.RawMessage  `json:"status,omitempty"`
	Online               bool             `json:"online"`
	LastSeen             string           `json:"last_seen,omitempty"`
	SourceIP             string           `json:"source_ip,omitempty"`
	PairedAt             string           `json:"paired_at"`
	RegistrationExpiresAt string          `json:"registration_expires_at,omitempty"`
	CurrentCertPEM       string           `json:"-"`
	PreviousCertPEM      string           `json:"-"`
	CertExpiresAt        string           `json:"cert_expires_at,omitempty"`
	PrevCertExpiresAt    string           `json:"-"`
	LastRotationAt       string           `json:"last_rotation_at,omitempty"`
	CSRPEM               string           `json:"-"`
	PairingCode          string           `json:"-"`
	AccessCode           string           `json:"-"`
	PairingExpiresAt     string           `json:"-"`
	ConfigUpdateID       int              `json:"config_update_id"`
}

// DeviceSummary is the shape returned by ListDevices.
type DeviceSummary struct {
	DeviceID      string `json:"device_id"`
	Message       string `json:"message"`
	Errors        []string `json:"errors"`
	OnlineDetails string `json:"online_details"`
	Online        bool   `json:"online"`
}

// DeviceMetadata nested in DescribeDevice response.
type DeviceMetadata struct {
	Online        bool   `json:"online"`
	OnlineDetails string `json:"online_details"`
	CertExpiration string `json:"cert_expiration"`
	SourceIP      string `json:"source_ip"`
	DeviceType    string `json:"device_type"`
	AccountID     string `json:"account_id"`
	PairedAt      string `json:"paired_at"`
}

// DeviceDetail is the shape returned by DescribeDevice.
type DeviceDetail struct {
	DeviceID            string          `json:"device_id"`
	Message             string          `json:"message"`
	Errors              []string        `json:"errors"`
	Registration        json.RawMessage `json:"registration,omitempty"`
	Configuration       json.RawMessage `json:"configuration,omitempty"`
	ActualConfiguration json.RawMessage `json:"actual_configuration,omitempty"`
	Status              json.RawMessage `json:"status,omitempty"`
	Online              bool            `json:"online"`
	OnlineDetails       string          `json:"online_details"`
	CertExpiration      string          `json:"cert_expiration"`
	DeviceMetadata      DeviceMetadata  `json:"device_metadata"`
}

// ThumbnailResponse returned by GetThumbnail.
type ThumbnailResponse struct {
	Message string         `json:"message"`
	Image   *ThumbnailImage `json:"image,omitempty"`
}

// ThumbnailImage within a ThumbnailResponse.
type ThumbnailImage struct {
	Base64Image string `json:"base64_image"`
	Timestamp   string `json:"timestamp"`
	ImageType   string `json:"image_type"`
	MaxSizeKB   int    `json:"max_size_KB"`
	ImageSizeKB int    `json:"image_size_KB"`
}
