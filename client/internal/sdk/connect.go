// Copyright 2025 Amazon.com Inc
// Licensed under the Apache License, Version 2.0
package sdk

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/credentials"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/models"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/utils"
	tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
)

// Connect implements the TR-12 connect state machine.
func (s *CddSdk) Connect(registration map[string]interface{}, hostID string) models.ConnectResponseContent {
	s.apiLock.Lock()
	defer s.apiLock.Unlock()
	s.logger.Info("Connect")

	// Initialize if needed
	if s.hostID == "" || s.hostID != hostID {
		s.reset()
		if err := s.initializeHost(registration, hostID); err != nil {
			s.reset()
			return models.ConnectResponseContent{
				Success: false, State: s.state,
				Message: fmt.Sprintf("Error in connect() %v", err),
				Error:   utils.ExceptionToErrorDetails(err),
			}
		}
	}

	var resp models.ConnectResponseContent
	var err error

	switch s.state {
	case models.StateConnecting:
		resp = models.ConnectResponseContent{
			Success: true, State: s.state, Message: "Connecting to the service",
		}
	case models.StateConnected:
		resp = models.ConnectResponseContent{
			Success:  true, State: s.state, Message: "Connected",
			DeviceId: tr12models.PtrString(s.certs.GetDeviceID()),
			Region:   tr12models.PtrString(s.certs.GetRegion()),
		}
	case models.StateReconnecting:
		resp = models.ConnectResponseContent{
			Success:  true, State: s.state, Message: "Reconnecting...",
			DeviceId: tr12models.PtrString(s.certs.GetDeviceID()),
			Region:   tr12models.PtrString(s.certs.GetRegion()),
		}
	case models.StatePairing:
		resp, err = s.handlePairingState()
	case models.StateDisconnected:
		resp, err = s.handleDisconnectedState()
	}

	if err != nil {
		s.logger.Errorf("Exception: %v", err)
		s.reset()
		return models.ConnectResponseContent{
			Success: false, State: s.state,
			Message: fmt.Sprintf("Error in connect() %v", err),
			Error:   utils.ExceptionToErrorDetails(err),
		}
	}
	return resp
}

func (s *CddSdk) handlePairingState() (models.ConnectResponseContent, error) {
	loaded, err := s.loadCerts()
	if err != nil {
		return models.ConnectResponseContent{}, err
	}
	if loaded {
		return s.startConnect()
	}
	// Expired?
	if s.pairer.IsExpired() {
		s.reset()
		return models.ConnectResponseContent{
			Success: false, State: s.state,
			Message: "Pairing code expired. Reconnect to get a new one.",
		}, nil
	}
	// Poll for credentials
	claimed, err := s.pairer.AuthenticatePairingCode()
	if err != nil {
		return models.ConnectResponseContent{}, err
	}
	if claimed {
		loaded, err := s.loadCerts()
		if err != nil {
			return models.ConnectResponseContent{}, err
		}
		if loaded {
			hs, _ := s.certs.GetHostSettings()
			if hs != nil {
				s.initThrottles(int(hs.MinIntervalPubSeconds))
			}
			return s.startConnect()
		}
		s.reset()
		return models.ConnectResponseContent{}, fmt.Errorf("device was authenticated, but couldn't load certs")
	}
	// Still waiting
	expires := tr12models.PtrFloat32(float32(s.pairer.ExpiresIn()))
	return models.ConnectResponseContent{
		Success:     true, State: models.StatePairing,
		Message:     "Waiting for device to be claimed",
		PairingCode: tr12models.PtrString(s.pairer.GetPairingCode()),
		Expires:     expires,
	}, nil
}

func (s *CddSdk) handleDisconnectedState() (models.ConnectResponseContent, error) {
	loaded, err := s.loadCerts()
	if err != nil {
		return models.ConnectResponseContent{}, err
	}
	if loaded {
		hs, _ := s.certs.GetHostSettings()
		if hs != nil {
			s.statusThrottle = utils.NewThrottle(int(hs.MinIntervalPubSeconds))
		}
		return s.startConnect()
	}
	// Need to pair
	if err := s.pairer.GetNewPairingCode(); err != nil {
		return models.ConnectResponseContent{}, err
	}
	s.transition(models.StatePairing)
	expires := tr12models.PtrFloat32(float32(s.pairer.ExpiresIn()))
	return models.ConnectResponseContent{
		Success:     true, State: models.StatePairing,
		Message:     "Connecting pending. Waiting for device to be claimed",
		PairingCode: tr12models.PtrString(s.pairer.GetPairingCode()),
		Expires:     expires,
	}, nil
}

// Disconnect gracefully disconnects from the host service.
func (s *CddSdk) Disconnect() models.DisconnectResponseContent {
	s.apiLock.Lock()
	defer s.apiLock.Unlock()
	s.logger.Info("Disconnect")
	s.reset()
	return models.DisconnectResponseContent{
		Success: true, State: models.StateDisconnected, Message: "Disconnected",
	}
}

// GetConnectionStatus returns the current connection state.
func (s *CddSdk) GetConnectionStatus() models.GetConnectionStatusResponseContent {
	s.logger.Info("Get Connection Status")
	return models.GetConnectionStatusResponseContent{
		Success: true, State: s.state, Message: "",
	}
}

// Deprovision removes the device from the host service and deletes credentials.
func (s *CddSdk) Deprovision(hostID string, force bool) models.DeprovisionResponseContent {
	s.apiLock.Lock()
	defer s.apiLock.Unlock()
	s.logger.Info("Deprovision")

	if !s.connectedTo(hostID) && !force {
		return models.DeprovisionResponseContent{
			Success: true, State: s.state,
			Message: fmt.Sprintf("Must use --force to deprovision client while not CONNECTED to: %s", hostID),
		}
	}
	if err := s.handleDeprovision(hostID); err != nil {
		return models.DeprovisionResponseContent{
			Success: false, State: s.state,
			Message: fmt.Sprintf("Error in Deprovision: %v", err),
			Error:   utils.ExceptionToErrorDetails(err),
		}
	}
	return models.DeprovisionResponseContent{
		Success: true, State: s.state,
		Message: fmt.Sprintf("Deprovisioned credentials for host: %s", hostID),
	}
}

// GetConfiguration returns the latest cached configuration.
func (s *CddSdk) GetConfiguration() models.GetConfigurationResponseContent {
	s.apiLock.Lock()
	defer s.apiLock.Unlock()
	s.logger.Info("Get Configuration")
	config := &models.ConfigurationData{
		Payload:  s.configPayload,
		UpdateID: s.configUpdateID,
	}
	return models.GetConfigurationResponseContent{
		Success: true, State: s.state,
		Message:       "Latest configuration provided",
		Configuration: config,
	}
}

// ReportStatus publishes a status payload to the host service.
func (s *CddSdk) ReportStatus(payload map[string]interface{}) models.ReportStatusResponseContent {
	s.apiLock.Lock()
	defer s.apiLock.Unlock()
	s.logger.Info("Report Status")

	if err := s.canPublishNow(s.statusThrottle); err != nil {
		return models.ReportStatusResponseContent{
			Success: false, State: s.state, Message: err.Error(),
			Error: utils.ExceptionToErrorDetails(err),
		}
	}
	hs, err := s.certs.GetHostSettings()
	if err != nil {
		return models.ReportStatusResponseContent{
			Success: false, State: s.state, Message: err.Error(),
			Error: utils.ExceptionToErrorDetails(err),
		}
	}
	if err := s.doPublishMessage(payload, hs.PubReportStatusTopic); err != nil {
		return models.ReportStatusResponseContent{
			Success: false, State: s.state,
			Message: fmt.Sprintf("Status update not sent: %v", err),
			Error:   utils.ExceptionToErrorDetails(err),
		}
	}
	return models.ReportStatusResponseContent{
		Success: true, State: s.state, Message: "Status update sent",
	}
}

// ReportConfiguration publishes an actual configuration payload to the host service.
func (s *CddSdk) ReportConfiguration(payload map[string]interface{}) models.ReportActualConfigurationResponseContent {
	s.apiLock.Lock()
	defer s.apiLock.Unlock()
	s.logger.Info("Report Configuration")

	if err := s.canPublishNow(s.configThrottle); err != nil {
		return models.ReportActualConfigurationResponseContent{
			Success: false, State: s.state, Message: err.Error(),
			Error: utils.ExceptionToErrorDetails(err),
		}
	}
	hs, err := s.certs.GetHostSettings()
	if err != nil {
		return models.ReportActualConfigurationResponseContent{
			Success: false, State: s.state, Message: err.Error(),
			Error: utils.ExceptionToErrorDetails(err),
		}
	}
	if err := s.doPublishMessage(payload, hs.PubReportActualConfigurationTopic); err != nil {
		return models.ReportActualConfigurationResponseContent{
			Success: false, State: s.state,
			Message: fmt.Sprintf("Configuration update not sent: %v", err),
			Error:   utils.ExceptionToErrorDetails(err),
		}
	}
	return models.ReportActualConfigurationResponseContent{
		Success: true, State: s.state, Message: "Configuration update sent",
	}
}

func (s *CddSdk) handleDeprovision(hostID string) error {
	defer s.reset()
	if s.connectedTo(hostID) {
		s.informHostServiceDeprovision(hostID)
	}
	return s.deleteCredentials(hostID)
}

func (s *CddSdk) informHostServiceDeprovision(hostID string) {
	if !s.connectedTo(hostID) {
		return
	}
	s.logger.Infof("Publishing deprovision message to host service: %s", hostID)
	hs, err := s.certs.GetHostSettings()
	if err != nil {
		return
	}
	reason := tr12models.DEPROVISIONED
	t := tr12models.PtrFloat32(float32(time.Now().Unix()))
	msg := models.DeprovisionDeviceRequestContent{
		Reason: &reason,
		Time:   t,
	}
	data, _ := json.Marshal(msg)
	if s.mqttClient != nil {
		token := s.mqttClient.Publish(hs.PubDeprovisionTopic, 1, false, data)
		token.WaitTimeout(3 * time.Second)
	}
	time.Sleep(3 * time.Second)
}

func (s *CddSdk) deleteCredentials(hostID string) error {
	if s.hostID == "" || s.hostID == hostID {
		s.logger.Infof("Deleting credentials for %s", hostID)
		store, err := credentials.NewStore(s.certsPath, s.deviceLocalID, hostID)
		if err != nil {
			return err
		}
		return store.Deprovision()
	}
	return nil
}
