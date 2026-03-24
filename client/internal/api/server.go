// Copyright 2025 Amazon.com Inc
// Licensed under the Apache License, Version 2.0
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
}

// --- Route handlers ---

type connectRequest struct {
	HostID       string                 `json:"hostId"`
	HostIDSnake  string                 `json:"host_id"`
	Registration map[string]interface{} `json:"registration"`
}

func (s *Server) connect(c *gin.Context) {
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
	Status map[string]interface{} `json:"status"`
}

func (s *Server) reportStatus(c *gin.Context) {
	var req reportStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.Status == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status is required"})
		return
	}
	resp := s.sdk.ReportStatus(req.Status)
	c.JSON(http.StatusOK, resp)
}

type reportConfigRequest struct {
	Configuration map[string]interface{} `json:"configuration"`
}

func (s *Server) reportConfiguration(c *gin.Context) {
	var req reportConfigRequest
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
