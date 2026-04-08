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
	ConsoleDir          string
	HTTPS               bool
	TLSCert             string
	TLSKey              string
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
	flag.StringVar(&cfg.ConsoleDir, "console-dir", "", "Path to console dist/ directory (serves at /console/)")
	flag.BoolVar(&cfg.HTTPS, "https", false, "Enable HTTPS for the HTTP API (uses the service CA cert)")
	flag.StringVar(&cfg.TLSCert, "tls-cert", "", "Path to TLS certificate file (e.g. Let's Encrypt fullchain.pem)")
	flag.StringVar(&cfg.TLSKey, "tls-key", "", "Path to TLS private key file (e.g. Let's Encrypt privkey.pem)")
	flag.Parse()
	return cfg
}
