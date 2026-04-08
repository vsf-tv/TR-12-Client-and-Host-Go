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
package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/api"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/broker"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/ca"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/config"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/db"
	mqtthandlers "github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/mqtt"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/service"
)

func main() {
	cfg := config.Parse()
	if cfg.HostAddress == "" {
		fmt.Fprintln(os.Stderr, "error: --host-address is required")
		os.Exit(1)
	}

	log.Println("TR-12 Host Service starting...")

	// 1. Open database
	store, err := db.New(cfg.DBPath)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer store.Close()
	log.Printf("  Database: %s", cfg.DBPath)

	// 2. Load or generate CA, server cert, JWT secret
	authority, err := ca.New(store, cfg.HostAddress)
	if err != nil {
		log.Fatalf("CA: %v", err)
	}
	jwtSecret, err := ca.GetJWTSecret(store)
	if err != nil {
		log.Fatalf("JWT secret: %v", err)
	}

	// 3. Build TLS config for MQTT broker
	tlsConfig := authority.TLSConfig()

	// 4. Start embedded MQTT broker
	mqttBroker := broker.New(cfg.MQTTPort, tlsConfig, store)
	go func() {
		if err := mqttBroker.Start(); err != nil {
			log.Fatalf("MQTT broker: %v", err)
		}
	}()
	log.Printf("  MQTT Broker: tls://%s:%d", cfg.HostAddress, cfg.MQTTPort)

	// 5. Create services
	accountSvc := service.NewAccountService(store, jwtSecret, cfg.JWTExpiryHours)
	deviceSvc := service.NewDeviceService(store, authority, mqttBroker, cfg)
	thumbnailSvc := service.NewThumbnailService(store, mqttBroker, cfg)
	logSvc := service.NewLogService(store, mqttBroker, cfg)
	rotationSvc := service.NewRotationService(store, authority, mqttBroker, cfg)

	// 6. Subscribe to wildcard device topics
	handlers := mqtthandlers.NewHandlers(store, mqttBroker, deviceSvc)
	handlers.Subscribe()

	// 7. Start background goroutines
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rotationSvc.Start(ctx, deviceSvc)

	// 8. Start HTTP/HTTPS server
	router := api.NewRouter(deviceSvc, accountSvc, thumbnailSvc, logSvc, cfg.ConsoleDir)
	httpAddr := fmt.Sprintf(":%d", cfg.HTTPPort)
	srv := &http.Server{Addr: httpAddr, Handler: router}

	go func() {
		scheme := "http"
		if cfg.HTTPS {
			scheme = "https"
		}
		log.Printf("  HTTP API: %s://%s:%d", scheme, cfg.HostAddress, cfg.HTTPPort)
		if cfg.ConsoleDir != "" {
			log.Printf("  Console: %s://%s:%d/console/", scheme, cfg.HostAddress, cfg.HTTPPort)
		}
		log.Printf("  Service ID: %s", cfg.ServiceID)
		log.Println("  Ready.")

		if cfg.HTTPS {
			var certFile, keyFile string
			if cfg.TLSCert != "" && cfg.TLSKey != "" {
				// Use provided cert/key (e.g. Let's Encrypt)
				certFile = cfg.TLSCert
				keyFile = cfg.TLSKey
			} else {
				// Fall back to self-signed service CA cert
				certPEM, _ := store.GetConfig("server_cert_pem")
				keyPEM, _ := store.GetConfig("server_key_pem")
				cf, _ := os.CreateTemp("", "tr12-cert-*.pem")
				kf, _ := os.CreateTemp("", "tr12-key-*.pem")
				cf.Write(certPEM)
				kf.Write(keyPEM)
				cf.Close()
				kf.Close()
				certFile = cf.Name()
				keyFile = kf.Name()
				defer os.Remove(certFile)
				defer os.Remove(keyFile)
			}
			httpsSrv := &http.Server{
				Addr:      httpAddr,
				Handler:   router,
				TLSConfig: &tls.Config{MinVersion: tls.VersionTLS12},
			}
			if err := httpsSrv.ListenAndServeTLS(certFile, keyFile); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTPS server: %v", err)
			}
		} else {
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatalf("HTTP server: %v", err)
			}
		}
	}()

	// 9. Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	cancel()
	srv.Shutdown(context.Background())
	mqttBroker.Stop()
	store.Close()
	log.Println("Stopped.")
}
