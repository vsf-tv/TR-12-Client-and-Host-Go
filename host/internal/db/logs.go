package db

import "database/sql"

// DeviceLog represents a stored device log.
type DeviceLog struct {
	DeviceID  string
	LogData   []byte
	UploadedAt string
	LogSizeKB int
}

// UpsertLog stores or replaces a log for a device.
func (s *Store) UpsertLog(l *DeviceLog) error {
	_, err := s.DB.Exec(`INSERT INTO device_logs (device_id, log_data, uploaded_at, log_size_kb)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(device_id) DO UPDATE SET
			log_data = excluded.log_data, uploaded_at = excluded.uploaded_at, log_size_kb = excluded.log_size_kb`,
		l.DeviceID, l.LogData, l.UploadedAt, l.LogSizeKB,
	)
	return err
}

// GetLog retrieves a log by device ID.
func (s *Store) GetLog(deviceID string) (*DeviceLog, error) {
	l := &DeviceLog{}
	err := s.DB.QueryRow(
		"SELECT device_id, log_data, uploaded_at, log_size_kb FROM device_logs WHERE device_id = ?", deviceID,
	).Scan(&l.DeviceID, &l.LogData, &l.UploadedAt, &l.LogSizeKB)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return l, nil
}

// DeleteLogsByDevice removes the log for a device.
func (s *Store) DeleteLogsByDevice(deviceID string) error {
	_, err := s.DB.Exec("DELETE FROM device_logs WHERE device_id = ?", deviceID)
	return err
}
