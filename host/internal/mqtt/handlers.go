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
package mqtt

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/broker"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/db"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/models"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/service"
)

// Handlers manages internal MQTT subscriptions.
type Handlers struct {
	store     *db.Store
	broker    *broker.Broker
	deviceSvc *service.DeviceService
}

// NewHandlers creates MQTT message handlers.
func NewHandlers(store *db.Store, b *broker.Broker, deviceSvc *service.DeviceService) *Handlers {
	return &Handlers{store: store, broker: b, deviceSvc: deviceSvc}
}

// Subscribe registers all wildcard subscriptions.
func (h *Handlers) Subscribe() {
	h.broker.Subscribe("cdd/+/registration/report", h.handleRegistration)
	h.broker.Subscribe("cdd/+/status/report", h.handleStatus)
	h.broker.Subscribe("cdd/+/config/actual/report", h.handleActualConfig)
	h.broker.Subscribe("cdd/+/deprovision/ack", h.handleDeprovision)
}

func (h *Handlers) handleRegistration(topic string, payload []byte) {
	deviceID := extractDeviceID(topic)
	if deviceID == "" {
		return
	}
	// Store the registration JSON as-is (no lossy transformation)
	var raw json.RawMessage
	if err := json.Unmarshal(payload, &raw); err != nil {
		log.Printf("[mqtt] invalid registration from %s: %v", deviceID, err)
		return
	}
	// The payload may be wrapped in {"message": ...} per ReportMessage
	var report struct {
		Message json.RawMessage `json:"message"`
	}
	if err := json.Unmarshal(payload, &report); err == nil && len(report.Message) > 0 {
		raw = report.Message
	}
	if err := h.store.UpdateDeviceRegistration(deviceID, raw); err != nil {
		log.Printf("[mqtt] error storing registration for %s: %v", deviceID, err)
	} else {
		log.Printf("[mqtt] registration updated for %s", deviceID)
	}
}

func (h *Handlers) handleStatus(topic string, payload []byte) {
	deviceID := extractDeviceID(topic)
	if deviceID == "" {
		return
	}
	var raw json.RawMessage
	if err := json.Unmarshal(payload, &raw); err != nil {
		log.Printf("[mqtt] invalid status from %s: %v", deviceID, err)
		return
	}
	var report struct {
		Message json.RawMessage `json:"message"`
	}
	if err := json.Unmarshal(payload, &report); err == nil && len(report.Message) > 0 {
		raw = report.Message
	}
	if err := h.store.UpdateDeviceStatus(deviceID, raw); err != nil {
		log.Printf("[mqtt] error storing status for %s: %v", deviceID, err)
	}
}

func (h *Handlers) handleActualConfig(topic string, payload []byte) {
	deviceID := extractDeviceID(topic)
	if deviceID == "" {
		return
	}
	var raw json.RawMessage
	if err := json.Unmarshal(payload, &raw); err != nil {
		log.Printf("[mqtt] invalid actual config from %s: %v", deviceID, err)
		return
	}
	var report struct {
		Message json.RawMessage `json:"message"`
	}
	if err := json.Unmarshal(payload, &report); err == nil && len(report.Message) > 0 {
		raw = report.Message
	}
	if err := h.store.UpdateDeviceActualConfig(deviceID, raw); err != nil {
		log.Printf("[mqtt] error storing actual config for %s: %v", deviceID, err)
	}
}

func (h *Handlers) handleDeprovision(topic string, payload []byte) {
	deviceID := extractDeviceID(topic)
	if deviceID == "" {
		return
	}

	device, err := h.store.GetDevice(deviceID)
	if err != nil || device == nil {
		return
	}

	var msg models.DeprovisionRequest
	json.Unmarshal(payload, &msg)

	if device.State == "DEPROVISIONED" {
		// Phase 2: device acknowledged — full cleanup
		log.Printf("[mqtt] device %s acknowledged deprovision, cleaning up", deviceID)
		if err := h.deviceSvc.FullCleanup(deviceID); err != nil {
			log.Printf("[mqtt] cleanup error for %s: %v", deviceID, err)
		}
	} else {
		// Device-initiated deprovision — immediate full cleanup
		log.Printf("[mqtt] device %s initiated deprovision, cleaning up", deviceID)
		if err := h.deviceSvc.FullCleanup(deviceID); err != nil {
			log.Printf("[mqtt] cleanup error for %s: %v", deviceID, err)
		}
	}
}

func extractDeviceID(topic string) string {
	parts := strings.SplitN(topic, "/", 3)
	if len(parts) < 2 {
		return ""
	}
	return parts[1]
}
