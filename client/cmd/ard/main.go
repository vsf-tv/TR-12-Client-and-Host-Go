// Copyright 2025 Amazon.com Inc
// Licensed under the Apache License, Version 2.0
//
// Application Reference Design (ARD) — simulates a 1-channel encoder device
// that makes REST calls on the CDD SDK daemon.
//
// Usage:
//
//	ard --host_id <host_id>
//	ard --host_id vsf_test_host --sdk_url http://127.0.0.1:8603
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/ard"
)

func main() {
	hostID := flag.String("host_id", "", "Host ID to connect to (required)")
	sdkURL := flag.String("sdk_url", "http://127.0.0.1:8603", "Base URL of the running CDD SDK")
	flag.Parse()

	if *hostID == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Resolve base path for payloads and thumbnail images
	execPath, err := os.Executable()
	if err != nil {
		execPath, _ = os.Getwd()
	}
	basePath := filepath.Dir(execPath)
	// If running via `go run`, fall back to working directory
	if _, err := os.Stat(filepath.Join(basePath, "payloads")); os.IsNotExist(err) {
		basePath, _ = os.Getwd()
	}

	app, err := ard.NewClientApplication(*sdkURL, basePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize ARD: %v\n", err)
		os.Exit(1)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		fmt.Println("\nReceived shutdown signal, cleaning up...")
		app.Stop()
		time.Sleep(500 * time.Millisecond)
		os.Exit(0)
	}()

	fmt.Printf("ARD connecting to host: %s via SDK at %s\n", *hostID, *sdkURL)
	app.RunLoop(*hostID)
}
