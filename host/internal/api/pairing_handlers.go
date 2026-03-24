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
	var req models.PairRequestContent
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "code": 400})
		return
	}
	log.Printf("[HOST /pair] Request: hostId=%s deviceType=%s version=%s", req.HostId, req.DeviceType, req.Version)
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
	var req models.AuthenticateRequestContent
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
