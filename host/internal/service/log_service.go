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
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/config"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/db"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/models"
)

// LogService manages device log requests and storage.
type LogService struct {
	store *db.Store
	mqtt  MQTTPublisher
	cfg   *config.Config
}

// NewLogService creates a new LogService.
func NewLogService(store *db.Store, mqtt MQTTPublisher, cfg *config.Config) *LogService {
	return &LogService{store: store, mqtt: mqtt, cfg: cfg}
}

// SetMQTT sets the MQTT publisher.
func (s *LogService) SetMQTT(mqtt MQTTPublisher) {
	s.mqtt = mqtt
}

// RequestLog publishes a log request to a device.
func (s *LogService) RequestLog(deviceID string) error {
	scheme := "http"
	if s.cfg.HTTPS {
		scheme = "https"
	}
	uploadURL := fmt.Sprintf("%s://%s:%d/upload/log/%s", scheme, s.cfg.HostAddress, s.cfg.HTTPPort, deviceID)
	expires := float32(time.Now().Add(5 * time.Minute).Unix())
	req := models.RequestLogRequestContent{
		Expires:    &expires,
		RemotePath: &uploadURL,
	}
	payload, _ := json.Marshal(req)
	topic := fmt.Sprintf("cdd/%s/log/subscription", deviceID)
	return s.mqtt.Publish(topic, payload, true) // retained — device picks up active subscription on reconnect
}

// StoreLog saves a device log to the database (overwrites previous).
func (s *LogService) StoreLog(deviceID string, data []byte) error {
	return s.store.UpsertLog(&db.DeviceLog{
		DeviceID:   deviceID,
		LogData:    data,
		UploadedAt: time.Now().UTC().Format(time.RFC3339),
		LogSizeKB:  len(data) / 1024,
	})
}
