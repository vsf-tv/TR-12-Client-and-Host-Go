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
	subscriptions sync.Map // key: "deviceId:channelId" → expiry time.Time
}

// NewThumbnailService creates a new ThumbnailService.
func NewThumbnailService(store *db.Store, mqtt MQTTPublisher, cfg *config.Config) *ThumbnailService {
	return &ThumbnailService{store: store, mqtt: mqtt, cfg: cfg}
}

// SetMQTT sets the MQTT publisher.
func (s *ThumbnailService) SetMQTT(mqtt MQTTPublisher) {
	s.mqtt = mqtt
}

// RequestThumbnail creates a subscription and publishes to the device.
func (s *ThumbnailService) RequestThumbnail(deviceID, channelID string) error {
	scheme := "http"
	if s.cfg.HTTPS {
		scheme = "https"
	}
	uploadURL := fmt.Sprintf("%s://%s:%d/upload/thumbnail/%s/%s", scheme, s.cfg.HostAddress, s.cfg.HTTPPort, deviceID, channelID)
	expiry := time.Now().Add(120 * time.Second)

	period := float32(5)
	maxSize := float32(500)
	headers := map[string]string{"Content-Type": "image/jpeg"}
	sub := models.RequestThumbnailRequestContent{
		Requests: map[string]models.ThumbnailRequest{
			channelID: {
				PeriodSeconds: &period,
				ExpiresAt:     &expiry,
				MaxSizeKB:     &maxSize,
				RemotePath:    &uploadURL,
				Headers:       &headers,
			},
		},
	}
	payload, _ := json.Marshal(sub)
	topic := fmt.Sprintf("cdd/%s/thumbnail/subscription", deviceID)
	log.Printf("[THUMB] requesting thumbnail for device=%s channel=%s remotePath=%s", deviceID, channelID, uploadURL)
	s.subscriptions.Store(deviceID+":"+channelID, expiry)
	return s.mqtt.Publish(topic, payload, true)
}

// StoreThumbnail saves a thumbnail to the database.
func (s *ThumbnailService) StoreThumbnail(deviceID, channelID string, data []byte, imageType string) error {
	return s.store.UpsertThumbnail(&db.Thumbnail{
		DeviceID:    deviceID,
		ChannelID:   channelID,
		ImageData:   data,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		ImageType:   imageType,
		ImageSizeKB: len(data) / 1024,
	})
}

// GetThumbnail retrieves a thumbnail, requesting one if no subscription exists.
func (s *ThumbnailService) GetThumbnail(deviceID, channelID string) (*db.Thumbnail, bool, error) {
	key := deviceID + ":" + channelID
	if val, ok := s.subscriptions.Load(key); ok {
		if time.Now().Before(val.(time.Time)) {
			t, err := s.store.GetThumbnail(deviceID, channelID)
			return t, true, err
		}
		s.subscriptions.Delete(key)
	}
	if err := s.RequestThumbnail(deviceID, channelID); err != nil {
		return nil, false, err
	}
	t, err := s.store.GetThumbnail(deviceID, channelID)
	return t, false, err
}