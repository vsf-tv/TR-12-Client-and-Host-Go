// Copyright 2025 Amazon.com Inc
// Licensed under the Apache License, Version 2.0
//
// ClientApplication — the main ARD run loop. Mirrors application.py.
// Connects to the SDK, polls for configuration, reports status, and
// simulates thumbnail emission.
package ard

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/pkg/cddmodels"
)

// ClientApplication simulates a 1-channel encoder device.
type ClientApplication struct {
	shim           *Tr12Shim
	sdk            *SDKClient
	basePath       string
	registration   *cddmodels.DeviceRegistration
	currentConfig  *cddmodels.DeviceConfiguration
	latestConfigID string
	running        bool
	thumbnailSDI   *ThumbnailSimulator
	thumbnailHDMI  *ThumbnailSimulator
}

// NewClientApplication creates a new ARD instance.
func NewClientApplication(sdkURL, basePath string) (*ClientApplication, error) {
	// Load registration into typed struct
	regFile := filepath.Join(basePath, "payloads", "1_channel_encoder", "registration.json")
	data, err := os.ReadFile(regFile)
	if err != nil {
		return nil, fmt.Errorf("cannot read registration file %s: %w", regFile, err)
	}
	var registration cddmodels.DeviceRegistration
	if err := json.Unmarshal(data, &registration); err != nil {
		return nil, fmt.Errorf("invalid registration JSON: %w", err)
	}

	app := &ClientApplication{
		shim:         NewTr12Shim(),
		sdk:          NewSDKClient(sdkURL),
		basePath:     basePath,
		registration: &registration,
		running:      true,
	}

	// Set up thumbnail simulators
	sdiDir := filepath.Join(basePath, "application_reference", "thumbnail_images_sdi")
	hdmiDir := filepath.Join(basePath, "application_reference", "thumbnail_images_hdmi")

	app.thumbnailSDI, err = NewThumbnailSimulator(sdiDir, "/tmp/image_sdi.jpg", 2, "sdi")
	if err != nil {
		fmt.Printf("Warning: SDI thumbnail simulator not available: %v\n", err)
	}
	app.thumbnailHDMI, err = NewThumbnailSimulator(hdmiDir, "/tmp/image_hdmi.jpg", 2, "hdmi")
	if err != nil {
		fmt.Printf("Warning: HDMI thumbnail simulator not available: %v\n", err)
	}

	return app, nil
}

// Stop signals the application to shut down.
func (a *ClientApplication) Stop() {
	a.running = false
	if _, err := a.sdk.Disconnect(); err != nil {
		fmt.Printf("Error during disconnect: %v\n", err)
	}
	if a.thumbnailSDI != nil {
		a.thumbnailSDI.Stop()
	}
	if a.thumbnailHDMI != nil {
		a.thumbnailHDMI.Stop()
	}
	a.shim.CB.Encoder.Stop()
}

// RunLoop is the main application loop — mirrors application.py run_loop().
func (a *ClientApplication) RunLoop(hostID string) {
	if a.thumbnailSDI != nil {
		a.thumbnailSDI.Start()
	}
	if a.thumbnailHDMI != nil {
		a.thumbnailHDMI.Start()
	}

	for a.running {
		fmt.Println("........................")

		resp, err := a.sdk.Connect(hostID, a.registration)
		if err != nil {
			fmt.Printf("An error occurred: %v\n", err)
			time.Sleep(3 * time.Second)
			continue
		}

		fmt.Printf("Success: %v State: %s  error: %v DeviceID: %s  message: %s\n",
			resp.Success, resp.State, resp.Error, resp.GetDeviceId(), resp.Message)

		if resp.Success && resp.State == "PAIRING" {
			fmt.Printf("Device is not paired. Pairing Code: %s Expires in: %ds.\n",
				resp.GetPairingCode(), int(resp.GetExpires()))
		}

		if resp.Success && resp.State == "CONNECTED" {
			a.getConfiguration()
			a.reportStatus()
			a.reportActualConfiguration()
		}

		time.Sleep(3 * time.Second)
	}

	if a.thumbnailSDI != nil {
		a.thumbnailSDI.Stop()
	}
	if a.thumbnailHDMI != nil {
		a.thumbnailHDMI.Stop()
	}
	fmt.Println("Exiting")
}

func (a *ClientApplication) reportStatus() {
	status := a.shim.GetDeviceStatus()
	resp, err := a.sdk.ReportStatus(status)
	if err != nil {
		fmt.Printf("report_status error: %v\n", err)
		return
	}
	fmt.Printf("report_status Success: %v  State: %v  error: %v  message: %v\n",
		resp.Success, resp.State, resp.Error, resp.Message)
}

func (a *ClientApplication) reportActualConfiguration() {
	if a.currentConfig == nil {
		fmt.Println("No configuration to report")
		return
	}
	resp, err := a.sdk.ReportActualConfiguration(a.currentConfig)
	if err != nil {
		fmt.Printf("report_actual_configuration error: %v\n", err)
		return
	}
	fmt.Printf("report_actual_configuration Success: %v  State: %v  error: %v  message: %v\n",
		resp.Success, resp.State, resp.Error, resp.Message)
}

func (a *ClientApplication) getConfiguration() {
	resp, err := a.sdk.GetConfiguration()
	if err != nil {
		fmt.Printf("get_configuration error: %v\n", err)
		return
	}

	fmt.Printf("get_configuration Success: %v State: %v  error: %v  message: %v\n",
		resp.Success, resp.State, resp.Error, resp.Message)

	if resp.Configuration == nil {
		return
	}

	updateID := resp.Configuration.GetUpdateId()
	payload := resp.Configuration.Payload

	if updateID != "" && updateID != a.latestConfigID {
		fmt.Printf("New update. update_id: %s\n", updateID)
		a.latestConfigID = updateID

		if payload != nil {
			// Convert map[string]interface{} payload to *cddmodels.DeviceConfiguration
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				fmt.Printf("get_configuration: failed to marshal payload: %v\n", err)
				return
			}
			var deviceConfig cddmodels.DeviceConfiguration
			if err := json.Unmarshal(payloadBytes, &deviceConfig); err != nil {
				fmt.Printf("get_configuration: failed to parse payload as DeviceConfiguration: %v\n", err)
				return
			}

			fmt.Println("[SHIM TEST] apply_desired_configuration.")
			success := a.shim.ApplyDesiredConfiguration(&deviceConfig)
			if success {
				a.currentConfig = &deviceConfig
				a.shim.PrintActualConfig(a.registration)
			}
		}
	}
}
