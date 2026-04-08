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
package api

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/models"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/service"
)

// DeviceHandlers handles device management endpoints.
type DeviceHandlers struct {
	deviceSvc    *service.DeviceService
	thumbnailSvc *service.ThumbnailService
}

// NewDeviceHandlers creates device handlers.
func NewDeviceHandlers(deviceSvc *service.DeviceService, thumbnailSvc *service.ThumbnailService) *DeviceHandlers {
	return &DeviceHandlers{deviceSvc: deviceSvc, thumbnailSvc: thumbnailSvc}
}

// ListDevices handles GET /devices.
func (h *DeviceHandlers) ListDevices(c *gin.Context) {
	accountID := c.GetString("account_id")
	devices, err := h.deviceSvc.ListDevices(accountID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": 500})
		return
	}
	c.JSON(http.StatusOK, devices)
}

// DescribeDevice handles GET /device/:deviceId.
func (h *DeviceHandlers) DescribeDevice(c *gin.Context) {
	accountID := c.GetString("account_id")
	deviceID := c.Param("deviceId")
	detail, err := h.deviceSvc.DescribeDevice(deviceID, accountID)
	if err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, detail)
}

// UpdateConfiguration handles PUT /device/:deviceId.
func (h *DeviceHandlers) UpdateConfiguration(c *gin.Context) {
	accountID := c.GetString("account_id")
	deviceID := c.Param("deviceId")
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "code": 400})
		return
	}

	// Try to parse as the new structured request first
	var req models.UpdateDeviceRequest
	if jsonErr := json.Unmarshal(body, &req); jsonErr == nil && (req.Metadata != nil || req.DeviceConfiguration != nil) {
		// New format: { metadata: {...}, deviceConfiguration: {...} }
		if req.Metadata != nil {
			if err := h.deviceSvc.UpdateDeviceMetadata(deviceID, accountID, req.Metadata); err != nil {
				writeServiceError(c, err)
				return
			}
		}
		if req.DeviceConfiguration != nil {
			cfgBytes, _ := json.Marshal(req.DeviceConfiguration)
			if err := h.deviceSvc.UpdateConfiguration(deviceID, accountID, cfgBytes); err != nil {
				writeServiceError(c, err)
				return
			}
		}
	} else {
		// Legacy format: raw DeviceConfiguration JSON
		if err := h.deviceSvc.UpdateConfiguration(deviceID, accountID, body); err != nil {
			writeServiceError(c, err)
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"device_id": deviceID, "message": "Device updated", "error": ""})
}

// Claim handles PUT /authorize/:pairingCode.
func (h *DeviceHandlers) Claim(c *gin.Context) {
	accountID := c.GetString("account_id")
	pairingCode := c.Param("pairingCode")
	var req models.ClaimRequest
	c.ShouldBindJSON(&req) // optional body
	if err := h.deviceSvc.Claim(pairingCode, accountID, req.ExpirationDays, req.LocationName, req.DeviceName, req.RotationIntervalDays); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Device claimed", "pairing_code": pairingCode})
}

// Deprovision handles PUT /deprovision/:deviceId.
func (h *DeviceHandlers) Deprovision(c *gin.Context) {
	accountID := c.GetString("account_id")
	deviceID := c.Param("deviceId")
	if err := h.deviceSvc.Deprovision(deviceID, accountID); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"device_id": deviceID, "message": "Device deprovisioned"})
}

// GetThumbnail handles GET /thumbnail/:deviceId.
func (h *DeviceHandlers) GetThumbnail(c *gin.Context) {
	accountID := c.GetString("account_id")
	deviceID := c.Param("deviceId")
	sourceID := c.Query("source")
	if sourceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "source query parameter required", "code": 400})
		return
	}

	// Verify device ownership
	_, err := h.deviceSvc.DescribeDevice(deviceID, accountID)
	if err != nil {
		writeServiceError(c, err)
		return
	}

	thumb, subscribed, err := h.thumbnailSvc.GetThumbnail(deviceID, sourceID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": 500})
		return
	}
	if thumb == nil {
		msg := "thumbnail requested, waiting for device"
		if subscribed {
			msg = "subscription active, no thumbnail yet"
		}
		c.JSON(http.StatusOK, models.ThumbnailResponse{Message: msg})
		return
	}
	c.JSON(http.StatusOK, models.ThumbnailResponse{
		Message: "ok",
		Image: &models.ThumbnailImage{
			Base64Image: base64.StdEncoding.EncodeToString(thumb.ImageData),
			Timestamp:   thumb.Timestamp,
			ImageType:   thumb.ImageType,
			MaxSizeKB:   100,
			ImageSizeKB: thumb.ImageSizeKB,
		},
	})
}

// RotateCredentials handles PUT /credentials/:deviceId.
func (h *DeviceHandlers) RotateCredentials(c *gin.Context) {
	accountID := c.GetString("account_id")
	deviceID := c.Param("deviceId")
	if err := h.deviceSvc.RotateCredentials(deviceID, accountID); err != nil {
		writeServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"device_id": deviceID, "message": "Credentials rotated"})
}

func writeServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "not found", "code": 404})
	case errors.Is(err, service.ErrForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden", "code": 403})
	case errors.Is(err, service.ErrConflict):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error(), "code": 409})
	case errors.Is(err, service.ErrBadRequest):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "code": 400})
	case errors.Is(err, service.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized", "code": 401})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": 500})
	}
}
