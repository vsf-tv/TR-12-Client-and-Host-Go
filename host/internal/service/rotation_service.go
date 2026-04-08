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
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/ca"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/config"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/db"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/models"
)

// RotationService handles background certificate rotation.
type RotationService struct {
	store    *db.Store
	ca       *ca.CA
	mqtt     MQTTPublisher
	cfg      *config.Config
}

// NewRotationService creates a new RotationService.
func NewRotationService(store *db.Store, ca *ca.CA, mqtt MQTTPublisher, cfg *config.Config) *RotationService {
	return &RotationService{store: store, ca: ca, mqtt: mqtt, cfg: cfg}
}

// SetMQTT sets the MQTT publisher.
func (s *RotationService) SetMQTT(mqtt MQTTPublisher) {
	s.mqtt = mqtt
}

// Start runs the background rotation, pairing cleanup, and registration expiry goroutines.
func (s *RotationService) Start(ctx context.Context, deviceSvc *DeviceService) {
	go s.rotationLoop(ctx)
	go s.pairingCleanupLoop(ctx)
	go s.registrationExpiryLoop(ctx, deviceSvc)
}

func (s *RotationService) rotationLoop(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.rotateExpiredCerts()
		}
	}
}

func (s *RotationService) rotateExpiredCerts() {
	threshold := time.Now().UTC().Add(-time.Duration(s.cfg.RotationIntervalDays) * 24 * time.Hour).Format(time.RFC3339)
	devices, err := s.store.GetDevicesNeedingRotation(threshold)
	if err != nil {
		log.Printf("[rotation] error querying devices: %v", err)
		return
	}
	for _, d := range devices {
		if err := s.rotateDevice(d); err != nil {
			log.Printf("[rotation] error rotating %s: %v", d.DeviceID, err)
		} else {
			log.Printf("[rotation] rotated certs for %s", d.DeviceID)
		}
	}
}

func (s *RotationService) rotateDevice(d *models.Device) error {
	newCert, err := s.ca.SignCSR([]byte(d.CSRPEM), d.DeviceID, s.cfg.CertExpiryDays)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	newExpires := now.Add(time.Duration(s.cfg.CertExpiryDays) * 24 * time.Hour).Format(time.RFC3339)

	if err := s.store.UpdateDeviceCerts(
		d.DeviceID,
		string(newCert),
		d.CurrentCertPEM,
		newExpires,
		d.CertExpiresAt,
		now.Format(time.RFC3339),
	); err != nil {
		return err
	}

	mqttURI := fmt.Sprintf("tls://%s:%d", s.cfg.HostAddress, s.cfg.MQTTPort)
	rotate := models.RotateCertificatesRequestContent{
		MqttUri:    mqttURI,
		DeviceCert: string(newCert),
		Region:     "local",
	}
	payload, _ := json.Marshal(rotate)
	topic := fmt.Sprintf("cdd/%s/certs/update", d.DeviceID)
	return s.mqtt.Publish(topic, payload, true)
}

func (s *RotationService) pairingCleanupLoop(ctx context.Context) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now().UTC().Format(time.RFC3339)
			ids, err := s.store.GetExpiredPairingDevices(now)
			if err != nil {
				log.Printf("[pairing-cleanup] error: %v", err)
				continue
			}
			for _, id := range ids {
				if err := s.store.DeleteDevice(id); err != nil {
					log.Printf("[pairing-cleanup] error deleting %s: %v", id, err)
				} else {
					log.Printf("[pairing-cleanup] removed expired pairing %s", id)
				}
			}
		}
	}
}

func (s *RotationService) registrationExpiryLoop(ctx context.Context, deviceSvc *DeviceService) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			now := time.Now().UTC().Format(time.RFC3339)
			ids, err := s.store.GetExpiredRegistrationDevices(now)
			if err != nil {
				log.Printf("[reg-expiry] error: %v", err)
				continue
			}
			for _, id := range ids {
				if err := s.store.UpdateDeviceState(id, "DEPROVISIONED", true); err != nil {
					log.Printf("[reg-expiry] error deprovisioning %s: %v", id, err)
				} else {
					log.Printf("[reg-expiry] deprovisioned expired device %s", id)
				}
			}
		}
	}
}
