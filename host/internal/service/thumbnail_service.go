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
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/config"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/db"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/models"
)

// ThumbnailService manages thumbnail subscriptions and storage.
type ThumbnailService struct {
	store         *db.Store
	mqtt          MQTTPublisher
	cfg           *config.Config
	subscriptions sync.Map // key: "deviceId:sourceId" → expiry time.Time
}

// NewThumbnailService creates a new ThumbnailService.
func NewThumbnailService(store *db.Store, mqtt MQTTPublisher, cfg *config.Config) *ThumbnailService {
	return &ThumbnailService{store: store, mqtt: mqtt, cfg: cfg}
}

// SetMQTT sets the MQTT publisher.
func (s *ThumbnailService) SetMQTT(mqtt MQTTPublisher) {
	s.mqtt = mqtt
}

// lookupThumbnailLocalPath reads the device's registration JSON from the DB
// and returns the localPath for the given sourceID from the thumbnails array.
func (s *ThumbnailService) lookupThumbnailLocalPath(deviceID, sourceID string) string {
	device, err := s.store.GetDevice(deviceID)
	if err != nil || device == nil || device.Registration == nil {
		log.Printf("[THUMB] no registration found for device %s", deviceID)
		return ""
	}
	// Registration is stored as raw JSON; parse just the thumbnails array.
	var reg struct {
		Thumbnails []struct {
			ID        string `json:"id"`
			LocalPath string `json:"localPath"`
		} `json:"thumbnails"`
	}
	if err := json.Unmarshal(device.Registration, &reg); err != nil {
		log.Printf("[THUMB] failed to parse registration for %s: %v", deviceID, err)
		return ""
	}
	for _, t := range reg.Thumbnails {
		if t.ID == sourceID {
			return t.LocalPath
		}
	}
	log.Printf("[THUMB] no thumbnail entry for source %q in device %s registration", sourceID, deviceID)
	return ""
}

// RequestThumbnail creates a subscription and publishes to the device.
func (s *ThumbnailService) RequestThumbnail(deviceID, sourceID string) error {
	scheme := "http"
	if s.cfg.HTTPS {
		scheme = "https"
	}
	uploadURL := fmt.Sprintf("%s://%s:%d/upload/thumbnail/%s/%s", scheme, s.cfg.HostAddress, s.cfg.HTTPPort, deviceID, sourceID)
	localPath := s.lookupThumbnailLocalPath(deviceID, sourceID)
	expiry := time.Now().Add(120 * time.Second)

	period := float32(5)
	exp := float32(expiry.Unix())
	maxSize := float32(500)
	sub := models.RequestThumbnailRequestContent{
		Requests: map[string]models.ThumbnailRequest{
			sourceID: {
				Period:          &period,
				Expires:         &exp,
				MaxSizeKilobyte: &maxSize,
				LocalPath:       &localPath,
				RemotePath:      &uploadURL,
			},
		},
	}
	payload, _ := json.Marshal(sub)
	topic := fmt.Sprintf("cdd/%s/thumbnail/subscription", deviceID)
	log.Printf("[THUMB] requesting thumbnail for device=%s source=%s localPath=%s remotePath=%s", deviceID, sourceID, localPath, uploadURL)
	s.subscriptions.Store(deviceID+":"+sourceID, expiry)
	return s.mqtt.Publish(topic, payload, true) // retained — device picks up active subscription on reconnect
}

// StoreThumbnail saves a thumbnail to the database.
func (s *ThumbnailService) StoreThumbnail(deviceID, sourceID string, data []byte, imageType string) error {
	return s.store.UpsertThumbnail(&db.Thumbnail{
		DeviceID:    deviceID,
		SourceID:    sourceID,
		ImageData:   data,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		ImageType:   imageType,
		ImageSizeKB: len(data) / 1024,
	})
}

// GetThumbnail retrieves a thumbnail, requesting one if no subscription exists.
func (s *ThumbnailService) GetThumbnail(deviceID, sourceID string) (*db.Thumbnail, bool, error) {
	// Check if subscription is active
	key := deviceID + ":" + sourceID
	if val, ok := s.subscriptions.Load(key); ok {
		if time.Now().Before(val.(time.Time)) {
			// Subscription active, try to get thumbnail
			t, err := s.store.GetThumbnail(deviceID, sourceID)
			return t, true, err
		}
		s.subscriptions.Delete(key)
	}
	// No active subscription — request one
	if err := s.RequestThumbnail(deviceID, sourceID); err != nil {
		return nil, false, err
	}
	// Try to return existing thumbnail if available
	t, err := s.store.GetThumbnail(deviceID, sourceID)
	return t, false, err
}
