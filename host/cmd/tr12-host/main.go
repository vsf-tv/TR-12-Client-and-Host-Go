package main

import (
	"context"
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

	// 8. Start HTTP server
	router := api.NewRouter(deviceSvc, accountSvc, thumbnailSvc, logSvc)
	httpAddr := fmt.Sprintf(":%d", cfg.HTTPPort)
	srv := &http.Server{Addr: httpAddr, Handler: router}

	go func() {
		log.Printf("  HTTP API: http://%s:%d", cfg.HostAddress, cfg.HTTPPort)
		log.Printf("  Service ID: %s", cfg.ServiceID)
		log.Println("  Ready.")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server: %v", err)
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
