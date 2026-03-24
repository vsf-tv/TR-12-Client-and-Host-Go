// Copyright 2025 Amazon.com Inc
// Licensed under the Apache License, Version 2.0
//
// CddSdk is the core TR-12 Client Device Discovery SDK engine.
// It manages the state machine, MQTT connection, pairing, authentication,
// and all pub/sub operations with the host service.
package sdk

import (
	"fmt"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/cddlogger"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/credentials"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/models"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/pairing"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/thumbnails"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/utils"
	tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
)

// CddSdk is the main SDK struct.
type CddSdk struct {
	certsPath     string
	deviceLocalID string
	deviceType    string
	basePath      string // path to the executable's directory for host_configuration lookup

	certs              *credentials.Store
	pairer             *pairing.Pairing
	logger             *cddlogger.CDDLogger
	mqttClient         mqtt.Client
	thumbnailManager   *thumbnails.Manager
	state              string
	hostID             string
	registration       map[string]interface{}
	configPayload      map[string]interface{} // raw MQTT config payload
	configUpdateID     string                 // update ID for config
	updateID           *utils.UpdateID
	statusThrottle     *utils.Throttle
	configThrottle     *utils.Throttle
	regDelivered       bool
	logRequest         tr12models.RequestLogRequestContent
	processingLogPut   bool
	logSpewDetected    int64

	apiLock sync.Mutex
}

// New creates a new CddSdk instance.
func New(certsPath, deviceLocalID, deviceType, logPath, basePath string) (*CddSdk, error) {
	if err := utils.ValidatePathExistsAndWriteable(certsPath); err != nil {
		return nil, err
	}
	sdk := &CddSdk{
		certsPath:     certsPath,
		deviceLocalID: deviceLocalID,
		deviceType:    deviceType,
		basePath:      basePath,
		state:         models.StateDisconnected,
		updateID:      utils.NewUpdateID(),
	}
	var err error
	sdk.logger, err = cddlogger.New(logPath, "", sdk.reportLogs)
	if err != nil {
		return nil, err
	}
	sdk.thumbnailManager = thumbnails.NewManager(sdk.logger)
	sdk.certs, _ = credentials.NewStore(certsPath, deviceLocalID, "undefined")
	sdk.initThrottles(1)
	return sdk, nil
}

// Shutdown gracefully stops all threads and disconnects.
func (s *CddSdk) Shutdown() {
	s.thumbnailManager.StopAll()
	if s.mqttClient != nil && s.mqttClient.IsConnected() {
		s.mqttClient.Disconnect(250)
	}
	s.logger.Close()
}

func (s *CddSdk) reset() {
	s.registration = nil
	s.updateID = utils.NewUpdateID()
	s.initThrottles(1)
	s.hostID = ""
	s.logger.UpdateDeviceID("")
	s.configPayload = nil
	s.configUpdateID = ""
	s.regDelivered = false
	s.thumbnailManager.StopAll()
	if s.mqttClient != nil && s.mqttClient.IsConnected() {
		s.mqttClient.Disconnect(250)
	}
	s.mqttClient = nil
	s.transition(models.StateDisconnected)
}

func (s *CddSdk) initThrottles(intervalSeconds int) {
	s.statusThrottle = utils.NewThrottle(intervalSeconds)
	s.configThrottle = utils.NewThrottle(intervalSeconds)
}

func (s *CddSdk) initializeHost(registration map[string]interface{}, hostID string) error {
	s.registration = registration
	hostConfig, err := utils.GetHostConfiguration(hostID, s.basePath)
	if err != nil {
		return err
	}
	s.certs, err = credentials.NewStore(s.certsPath, s.deviceLocalID, hostID)
	if err != nil {
		return err
	}
	s.pairer = pairing.New(s.certs, s.deviceType, hostConfig.ServiceId, hostConfig.PairingUrl, hostConfig.AuthUrl)
	s.hostID = hostID
	return nil
}

func (s *CddSdk) transition(state string) {
	if state == models.StateConnected || state == models.StateConnecting {
		s.logger.UpdateDeviceID(s.certs.GetDeviceID())
	}
	if state == models.StateDisconnected {
		s.logger.Dump()
	}
	s.logger.Infof("Setting state to %s", state)
	s.state = state
}

func (s *CddSdk) is(states ...string) bool {
	for _, st := range states {
		if s.state == st {
			return true
		}
	}
	return false
}

func (s *CddSdk) connectedTo(hostID string) bool {
	return s.is(models.StateConnected) && hostID == s.hostID
}

func (s *CddSdk) loadCerts() (bool, error) {
	return s.certs.ReadFromFilesystem()
}

func (s *CddSdk) canPublishNow(throttle *utils.Throttle) error {
	if !s.is(models.StateConnected) {
		return fmt.Errorf("must be CONNECTED to publish message to host")
	}
	if !throttle.CanPublish() {
		return fmt.Errorf("request throttled: too many requests")
	}
	return nil
}
