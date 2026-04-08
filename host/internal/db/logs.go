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
