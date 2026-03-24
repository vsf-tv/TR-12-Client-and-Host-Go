package db

import (
	"database/sql"
	"encoding/json"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/models"
)

// InsertDevice creates a new device record.
func (s *Store) InsertDevice(d *models.Device) error {
	_, err := s.DB.Exec(`INSERT INTO devices (
		device_id, account_id, device_type, state, registration, desired_config, actual_config, status,
		online, last_seen, source_ip, paired_at, registration_expires_at,
		current_cert_pem, previous_cert_pem, cert_expires_at, prev_cert_expires_at, last_rotation_at,
		csr_pem, pairing_code, access_code, pairing_expires_at, config_update_id
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		d.DeviceID, d.AccountID, d.DeviceType, d.State,
		nullableJSON(d.Registration), nullableJSON(d.DesiredConfig), nullableJSON(d.ActualConfig), nullableJSON(d.Status),
		boolToInt(d.Online), d.LastSeen, d.SourceIP, d.PairedAt, d.RegistrationExpiresAt,
		d.CurrentCertPEM, d.PreviousCertPEM, d.CertExpiresAt, d.PrevCertExpiresAt, d.LastRotationAt,
		d.CSRPEM, d.PairingCode, d.AccessCode, d.PairingExpiresAt, d.ConfigUpdateID,
	)
	return err
}

// GetDevice retrieves a device by ID.
func (s *Store) GetDevice(deviceID string) (*models.Device, error) {
	d := &models.Device{}
	var online int
	var reg, desCfg, actCfg, status sql.NullString
	var lastSeen, sourceIP, regExpires sql.NullString
	var currentCert, previousCert, certExpires, prevCertExpires, lastRotation sql.NullString
	var csrPEM, pairingCode, accessCode, pairingExpires sql.NullString
	err := s.DB.QueryRow(`SELECT
		device_id, account_id, device_type, state, registration, desired_config, actual_config, status,
		online, last_seen, source_ip, paired_at, registration_expires_at,
		current_cert_pem, previous_cert_pem, cert_expires_at, prev_cert_expires_at, last_rotation_at,
		csr_pem, pairing_code, access_code, pairing_expires_at, config_update_id
		FROM devices WHERE device_id = ?`, deviceID).Scan(
		&d.DeviceID, &d.AccountID, &d.DeviceType, &d.State,
		&reg, &desCfg, &actCfg, &status,
		&online, &lastSeen, &sourceIP, &d.PairedAt, &regExpires,
		&currentCert, &previousCert, &certExpires, &prevCertExpires, &lastRotation,
		&csrPEM, &pairingCode, &accessCode, &pairingExpires, &d.ConfigUpdateID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	d.Online = online != 0
	d.Registration = jsonOrNil(reg)
	d.DesiredConfig = jsonOrNil(desCfg)
	d.ActualConfig = jsonOrNil(actCfg)
	d.Status = jsonOrNil(status)
	d.LastSeen = nullStr(lastSeen)
	d.SourceIP = nullStr(sourceIP)
	d.RegistrationExpiresAt = nullStr(regExpires)
	d.CurrentCertPEM = nullStr(currentCert)
	d.PreviousCertPEM = nullStr(previousCert)
	d.CertExpiresAt = nullStr(certExpires)
	d.PrevCertExpiresAt = nullStr(prevCertExpires)
	d.LastRotationAt = nullStr(lastRotation)
	d.CSRPEM = nullStr(csrPEM)
	d.PairingCode = nullStr(pairingCode)
	d.AccessCode = nullStr(accessCode)
	d.PairingExpiresAt = nullStr(pairingExpires)
	return d, nil
}

// GetDeviceByPairingCode looks up a device by pairing code.
func (s *Store) GetDeviceByPairingCode(code string) (*models.Device, error) {
	d := &models.Device{}
	var online int
	var reg, desCfg, actCfg, status sql.NullString
	var lastSeen, sourceIP, regExpires sql.NullString
	var currentCert, previousCert, certExpires, prevCertExpires, lastRotation sql.NullString
	var csrPEM, pairingCode, accessCode, pairingExpires sql.NullString
	err := s.DB.QueryRow(`SELECT
		device_id, account_id, device_type, state, registration, desired_config, actual_config, status,
		online, last_seen, source_ip, paired_at, registration_expires_at,
		current_cert_pem, previous_cert_pem, cert_expires_at, prev_cert_expires_at, last_rotation_at,
		csr_pem, pairing_code, access_code, pairing_expires_at, config_update_id
		FROM devices WHERE pairing_code = ?`, code).Scan(
		&d.DeviceID, &d.AccountID, &d.DeviceType, &d.State,
		&reg, &desCfg, &actCfg, &status,
		&online, &lastSeen, &sourceIP, &d.PairedAt, &regExpires,
		&currentCert, &previousCert, &certExpires, &prevCertExpires, &lastRotation,
		&csrPEM, &pairingCode, &accessCode, &pairingExpires, &d.ConfigUpdateID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	d.Online = online != 0
	d.Registration = jsonOrNil(reg)
	d.DesiredConfig = jsonOrNil(desCfg)
	d.ActualConfig = jsonOrNil(actCfg)
	d.Status = jsonOrNil(status)
	d.LastSeen = nullStr(lastSeen)
	d.SourceIP = nullStr(sourceIP)
	d.RegistrationExpiresAt = nullStr(regExpires)
	d.CurrentCertPEM = nullStr(currentCert)
	d.PreviousCertPEM = nullStr(previousCert)
	d.CertExpiresAt = nullStr(certExpires)
	d.PrevCertExpiresAt = nullStr(prevCertExpires)
	d.LastRotationAt = nullStr(lastRotation)
	d.CSRPEM = nullStr(csrPEM)
	d.PairingCode = nullStr(pairingCode)
	d.AccessCode = nullStr(accessCode)
	d.PairingExpiresAt = nullStr(pairingExpires)
	return d, nil
}

// ListDevicesByAccount returns all devices for an account.
func (s *Store) ListDevicesByAccount(accountID string) ([]*models.Device, error) {
	rows, err := s.DB.Query(`SELECT
		device_id, account_id, device_type, state, online, last_seen, paired_at, cert_expires_at
		FROM devices WHERE account_id = ? ORDER BY paired_at DESC`, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var devices []*models.Device
	for rows.Next() {
		d := &models.Device{}
		var online int
		var lastSeen, certExpires sql.NullString
		if err := rows.Scan(&d.DeviceID, &d.AccountID, &d.DeviceType, &d.State, &online, &lastSeen, &d.PairedAt, &certExpires); err != nil {
			return nil, err
		}
		d.Online = online != 0
		d.LastSeen = nullStr(lastSeen)
		d.CertExpiresAt = nullStr(certExpires)
		devices = append(devices, d)
	}
	return devices, rows.Err()
}

// UpdateDeviceState sets the device state and optionally clears config/status fields.
func (s *Store) UpdateDeviceState(deviceID, state string, clearData bool) error {
	if clearData {
		_, err := s.DB.Exec(
			"UPDATE devices SET state = ?, desired_config = NULL, actual_config = NULL, status = NULL WHERE device_id = ?",
			state, deviceID,
		)
		return err
	}
	_, err := s.DB.Exec("UPDATE devices SET state = ? WHERE device_id = ?", state, deviceID)
	return err
}

// ClaimDevice marks a device as claimed by an account.
// Pairing fields are preserved so the SDK can complete its final /authenticate poll.
func (s *Store) ClaimDevice(deviceID, accountID, registrationExpiresAt string) error {
	_, err := s.DB.Exec(`UPDATE devices SET
		account_id = ?, state = 'ACTIVE', registration_expires_at = ?
		WHERE device_id = ?`,
		accountID, registrationExpiresAt, deviceID,
	)
	return err
}

// UpdateDeviceRegistration stores the registration JSON.
func (s *Store) UpdateDeviceRegistration(deviceID string, registration json.RawMessage) error {
	_, err := s.DB.Exec("UPDATE devices SET registration = ? WHERE device_id = ?", string(registration), deviceID)
	return err
}

// UpdateDeviceStatus stores the status JSON.
func (s *Store) UpdateDeviceStatus(deviceID string, status json.RawMessage) error {
	_, err := s.DB.Exec("UPDATE devices SET status = ? WHERE device_id = ?", string(status), deviceID)
	return err
}

// UpdateDeviceActualConfig stores the actual configuration JSON.
func (s *Store) UpdateDeviceActualConfig(deviceID string, config json.RawMessage) error {
	_, err := s.DB.Exec("UPDATE devices SET actual_config = ? WHERE device_id = ?", string(config), deviceID)
	return err
}

// UpdateDeviceDesiredConfig stores the desired configuration and increments update ID.
func (s *Store) UpdateDeviceDesiredConfig(deviceID string, config json.RawMessage) (int, error) {
	result := s.DB.QueryRow(
		"UPDATE devices SET desired_config = ?, config_update_id = config_update_id + 1 WHERE device_id = ? RETURNING config_update_id",
		string(config), deviceID,
	)
	var updateID int
	err := result.Scan(&updateID)
	return updateID, err
}

// UpdateDeviceOnline sets the online state and last_seen.
func (s *Store) UpdateDeviceOnline(deviceID string, online bool, sourceIP, lastSeen string) error {
	_, err := s.DB.Exec(
		"UPDATE devices SET online = ?, source_ip = ?, last_seen = ? WHERE device_id = ?",
		boolToInt(online), sourceIP, lastSeen, deviceID,
	)
	return err
}

// UpdateDeviceCerts updates certificate fields after rotation.
func (s *Store) UpdateDeviceCerts(deviceID, currentCert, previousCert, certExpires, prevCertExpires, lastRotation string) error {
	_, err := s.DB.Exec(`UPDATE devices SET
		current_cert_pem = ?, previous_cert_pem = ?, cert_expires_at = ?, prev_cert_expires_at = ?, last_rotation_at = ?
		WHERE device_id = ?`,
		currentCert, previousCert, certExpires, prevCertExpires, lastRotation, deviceID,
	)
	return err
}

// RevokePreviousCert clears the previous cert fields for a device.
func (s *Store) RevokePreviousCert(deviceID string) error {
	_, err := s.DB.Exec("UPDATE devices SET previous_cert_pem = NULL, prev_cert_expires_at = NULL WHERE device_id = ?", deviceID)
	return err
}

// DeleteDevice removes a device record.
func (s *Store) DeleteDevice(deviceID string) error {
	_, err := s.DB.Exec("DELETE FROM devices WHERE device_id = ?", deviceID)
	return err
}

// GetExpiredPairingDevices returns devices in PAIRING state past their expiry.
func (s *Store) GetExpiredPairingDevices(now string) ([]string, error) {
	rows, err := s.DB.Query(
		"SELECT device_id FROM devices WHERE state = 'PAIRING' AND pairing_expires_at IS NOT NULL AND pairing_expires_at < ?", now,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// GetExpiredRegistrationDevices returns active devices past their registration expiry.
func (s *Store) GetExpiredRegistrationDevices(now string) ([]string, error) {
	rows, err := s.DB.Query(
		"SELECT device_id FROM devices WHERE state = 'ACTIVE' AND registration_expires_at IS NOT NULL AND registration_expires_at < ?", now,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// GetDevicesNeedingRotation returns active devices whose cert is older than the given threshold.
func (s *Store) GetDevicesNeedingRotation(threshold string) ([]*models.Device, error) {
	rows, err := s.DB.Query(`SELECT device_id, csr_pem, current_cert_pem, cert_expires_at, last_rotation_at
		FROM devices WHERE state = 'ACTIVE' AND csr_pem IS NOT NULL AND
		(last_rotation_at IS NULL OR last_rotation_at < ?)`, threshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var devices []*models.Device
	for rows.Next() {
		d := &models.Device{}
		if err := rows.Scan(&d.DeviceID, &d.CSRPEM, &d.CurrentCertPEM, &d.CertExpiresAt, &d.LastRotationAt); err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}
	return devices, rows.Err()
}

// Helper functions

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func nullableJSON(data json.RawMessage) interface{} {
	if len(data) == 0 {
		return nil
	}
	return string(data)
}

func jsonOrNil(ns sql.NullString) json.RawMessage {
	if !ns.Valid || ns.String == "" {
		return nil
	}
	return json.RawMessage(ns.String)
}

func nullStr(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}
