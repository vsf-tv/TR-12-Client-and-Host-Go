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
// CDD SDK main entry point — starts the TR-12 Client Device Discovery SDK process.
// Exposes the same CLI arguments and REST API as the Python SDK.
//
// Usage:
//
//	cdd-sdk --internal_device_id <id> --certs_path <path> --log_path <path> \
//	        --ip <ip> --port <port> --device_type <SOURCE|DESTINATION|BOTH>
package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/api"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/sdk"
)

func main() {
	deviceID := flag.String("internal_device_id", "", "Enter a device name (required)")
	certsPath := flag.String("certs_path", "", "Enter a path for persistent cert storage (required)")
	logPath := flag.String("log_path", "", "Enter a writable path for log storage (required)")
	ip := flag.String("ip", "", "IP on which the SDK will host REST APIs (required)")
	port := flag.String("port", "", "Port on which the SDK will host REST APIs (required)")
	deviceType := flag.String("device_type", "", "Device type: SOURCE|DESTINATION|BOTH (required)")
	flag.Parse()

	if *deviceID == "" || *certsPath == "" || *logPath == "" || *ip == "" || *port == "" || *deviceType == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Determine base path for host_configuration lookup — relative to the executable
	execPath, err := os.Executable()
	if err != nil {
		execPath, _ = os.Getwd()
	}
	basePath := filepath.Dir(execPath)
	// If running via `go run`, use the working directory instead
	if _, err := os.Stat(filepath.Join(basePath, "host_configuration")); os.IsNotExist(err) {
		basePath, _ = os.Getwd()
	}

	sdkClient, err := sdk.New(*certsPath, *deviceID, *deviceType, *logPath, basePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize SDK: %v\n", err)
		os.Exit(1)
	}

	server := api.NewServer(sdkClient)

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		sdkClient.Shutdown()
		os.Exit(0)
	}()

	fmt.Printf("CDD SDK (Go) starting on %s:%s\n", *ip, *port)
	if err := server.Run(*ip, *port); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		sdkClient.Shutdown()
		os.Exit(1)
	}
}
