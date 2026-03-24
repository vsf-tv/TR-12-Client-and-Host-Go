package config

import "flag"

// Config holds all CLI configuration for the host service.
type Config struct {
	HTTPPort            int
	MQTTPort            int
	DBPath              string
	ServiceID           string
	ServiceName         string
	HostAddress         string
	CertExpiryDays      int
	RotationIntervalDays int
	PairingTimeout      int
	JWTExpiryHours      int
	LogLevel            string
}

// Parse reads CLI flags and returns a Config.
func Parse() *Config {
	cfg := &Config{}
	flag.IntVar(&cfg.HTTPPort, "http-port", 8080, "Management API HTTP port")
	flag.IntVar(&cfg.MQTTPort, "mqtt-port", 8883, "MQTT broker TLS port")
	flag.StringVar(&cfg.DBPath, "db-path", "./tr12-host.db", "Path to SQLite database file")
	flag.StringVar(&cfg.ServiceID, "service-id", "tr12-host", "Service identifier")
	flag.StringVar(&cfg.ServiceName, "service-name", "TR-12 Host Service", "Human-readable service name")
	flag.StringVar(&cfg.HostAddress, "host-address", "", "Externally reachable address (required)")
	flag.IntVar(&cfg.CertExpiryDays, "cert-expiry-days", 30, "Device certificate validity period in days")
	flag.IntVar(&cfg.RotationIntervalDays, "rotation-interval-days", 30, "Auto-rotation interval in days")
	flag.IntVar(&cfg.PairingTimeout, "pairing-timeout", 1800, "Pairing code timeout in seconds")
	flag.IntVar(&cfg.JWTExpiryHours, "jwt-expiry-hours", 24, "JWT token expiration in hours")
	flag.StringVar(&cfg.LogLevel, "log-level", "info", "Log verbosity: debug, info, warn, error")
	flag.Parse()
	return cfg
}
