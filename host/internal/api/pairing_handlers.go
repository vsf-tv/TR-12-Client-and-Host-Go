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
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/models"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/service"
)

// PairingHandlers handles /pair and /authenticate endpoints.
type PairingHandlers struct {
	deviceSvc *service.DeviceService
}

// NewPairingHandlers creates pairing handlers.
func NewPairingHandlers(deviceSvc *service.DeviceService) *PairingHandlers {
	return &PairingHandlers{deviceSvc: deviceSvc}
}

// Pair handles POST /pair.
func (h *PairingHandlers) Pair(c *gin.Context) {
	var req models.CreatePairingCodeRequestContent
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "code": 400})
		return
	}
	log.Printf("[HOST /pair] Request: hostId=%s deviceType=%s version=%s", req.HostId, req.DeviceType, req.Version.GetVersion())
	resp, err := h.deviceSvc.Pair(req)
	if err != nil {
		log.Printf("[HOST /pair] Error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": 500})
		return
	}
	respJSON, _ := json.Marshal(resp)
	log.Printf("[HOST /pair] Response: %s", string(respJSON))
	c.JSON(http.StatusOK, resp)
}

// Authenticate handles POST /authenticate.
func (h *PairingHandlers) Authenticate(c *gin.Context) {
	var req models.AuthenticatePairingCodeRequestContent
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "code": 400})
		return
	}
	log.Printf("[HOST /authenticate] Request: deviceId=%s pairingCode=%s", req.DeviceId, req.PairingCode)
	resp, err := h.deviceSvc.Authenticate(req)
	if err != nil {
		log.Printf("[HOST /authenticate] Error: %v", err)
		if err == service.ErrNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "device not found", "code": 404})
			return
		}
		if err == service.ErrUnauthorized {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials", "code": 401})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": 500})
		return
	}
	respJSON, _ := json.Marshal(resp)
	log.Printf("[HOST /authenticate] Response: status=%s (full: %s)", resp.Status, string(respJSON))
	c.JSON(http.StatusOK, resp)
}
