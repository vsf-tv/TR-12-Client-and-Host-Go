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
// REST API server hosting the TR-12 client-side endpoints.
// Mirrors the Python SDK's Flask server_flask.py.
package api

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/sdk"
	cddsdkgo "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"
)

// Server wraps the Gin engine and SDK client.
type Server struct {
	engine *gin.Engine
	sdk    *sdk.CddSdk
}

// NewServer creates a new API server.
func NewServer(sdkClient *sdk.CddSdk) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: false,
	}))

	s := &Server{engine: engine, sdk: sdkClient}
	s.setupRoutes()
	return s
}

// Run starts the HTTP server.
func (s *Server) Run(host string, port string) error {
	return s.engine.Run(host + ":" + port)
}

func (s *Server) setupRoutes() {
	s.engine.PUT("/connect", s.connect)
	s.engine.PUT("/disconnect", s.disconnect)
	s.engine.GET("/get_state", s.getState)
	s.engine.PUT("/report_status", s.reportStatus)
	s.engine.PUT("/report_actual_configuration", s.reportConfiguration)
	s.engine.GET("/get_configuration", s.getConfiguration)
	s.engine.PUT("/deprovision", s.deprovision)
	s.engine.PUT("/register", s.register)
}

// --- Route handlers ---

type connectRequest struct {
	HostID       string                      `json:"hostId"`
	HostIDSnake  string                      `json:"host_id"`
	Registration *cddsdkgo.DeviceRegistration `json:"registration"`
}

func (s *Server) connect(c *gin.Context) {
	// Deserialise registration into the typed DeviceRegistration struct.
	// This validates structure at the HTTP boundary — unknown fields are
	// ignored by encoding/json, but required fields and type mismatches
	// are caught here before reaching the SDK or MQTT.
	var req connectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is required"})
		return
	}
	hostID := req.HostID
	if hostID == "" {
		hostID = req.HostIDSnake
	}
	if hostID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "host_id is required"})
		return
	}
	if req.Registration == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "registration is required"})
		return
	}
	resp := s.sdk.Connect(req.Registration, hostID)
	c.JSON(http.StatusOK, resp)
}

func (s *Server) disconnect(c *gin.Context) {
	resp := s.sdk.Disconnect()
	c.JSON(http.StatusOK, resp)
}

func (s *Server) getState(c *gin.Context) {
	resp := s.sdk.GetConnectionStatus()
	c.JSON(http.StatusOK, resp)
}

type reportStatusRequest struct {
	Status *cddsdkgo.DeviceStatus `json:"status"`
}

func (s *Server) reportStatus(c *gin.Context) {
	// Deserialise directly into the typed DeviceStatus struct.
	// This validates structure at the HTTP boundary — malformed payloads
	// are rejected with 400 before they can reach the MQTT publish path.
	var req reportStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required and must be a valid DeviceStatus"})
		return
	}
	resp := s.sdk.ReportStatus(req.Status)
	c.JSON(http.StatusOK, resp)
}

func (s *Server) reportConfiguration(c *gin.Context) {
	var req struct {
		Configuration *cddsdkgo.ActualDeviceConfiguration `json:"configuration"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Configuration == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "configuration is required"})
		return
	}
	resp := s.sdk.ReportConfiguration(req.Configuration)
	c.JSON(http.StatusOK, resp)
}

func (s *Server) getConfiguration(c *gin.Context) {
	resp := s.sdk.GetConfiguration()
	c.JSON(http.StatusOK, resp)
}

type deprovisionRequest struct {
	HostID      string `json:"hostId"`
	HostIDSnake string `json:"host_id"`
}

func (s *Server) deprovision(c *gin.Context) {
	var req deprovisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body is required"})
		return
	}
	hostID := req.HostID
	if hostID == "" {
		hostID = req.HostIDSnake
	}
	if hostID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hostId is required"})
		return
	}
	forceStr := c.Query("force")
	force := strings.EqualFold(forceStr, "true") || forceStr == "1"
	resp := s.sdk.Deprovision(hostID, force)
	c.JSON(http.StatusOK, resp)
}

type registerRequest struct {
	Registration *cddsdkgo.DeviceRegistration `json:"registration"`
}

func (s *Server) register(c *gin.Context) {
	// Deserialise into typed DeviceRegistration — validates structure at the HTTP boundary.
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Registration == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "registration is required and must be a valid DeviceRegistration"})
		return
	}
	resp := s.sdk.Register(req.Registration)
	c.JSON(http.StatusOK, resp)
}