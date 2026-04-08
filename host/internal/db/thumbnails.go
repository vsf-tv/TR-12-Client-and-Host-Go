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

// Thumbnail represents a stored thumbnail.
type Thumbnail struct {
	DeviceID    string
	SourceID    string
	ImageData   []byte
	Timestamp   string
	ImageType   string
	ImageSizeKB int
}

// UpsertThumbnail stores or replaces a thumbnail for a device+source.
func (s *Store) UpsertThumbnail(t *Thumbnail) error {
	_, err := s.DB.Exec(`INSERT INTO thumbnails (device_id, source_id, image_data, timestamp, image_type, image_size_kb)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(device_id, source_id) DO UPDATE SET
			image_data = excluded.image_data, timestamp = excluded.timestamp,
			image_type = excluded.image_type, image_size_kb = excluded.image_size_kb`,
		t.DeviceID, t.SourceID, t.ImageData, t.Timestamp, t.ImageType, t.ImageSizeKB,
	)
	return err
}

// GetThumbnail retrieves a thumbnail by device and source.
func (s *Store) GetThumbnail(deviceID, sourceID string) (*Thumbnail, error) {
	t := &Thumbnail{}
	err := s.DB.QueryRow(
		"SELECT device_id, source_id, image_data, timestamp, image_type, image_size_kb FROM thumbnails WHERE device_id = ? AND source_id = ?",
		deviceID, sourceID,
	).Scan(&t.DeviceID, &t.SourceID, &t.ImageData, &t.Timestamp, &t.ImageType, &t.ImageSizeKB)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return t, nil
}

// DeleteThumbnailsByDevice removes all thumbnails for a device.
func (s *Store) DeleteThumbnailsByDevice(deviceID string) error {
	_, err := s.DB.Exec("DELETE FROM thumbnails WHERE device_id = ?", deviceID)
	return err
}
