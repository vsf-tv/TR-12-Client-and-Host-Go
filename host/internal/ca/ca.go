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
package ca

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"database/sql"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/db"
)

// CA manages the service certificate authority.
type CA struct {
	CACert    *x509.Certificate
	CAKey     crypto.PrivateKey
	CACertPEM []byte
	ServerTLS tls.Certificate
	store     *db.Store
}

// New loads or generates the CA, server cert, and JWT secret.
func New(store *db.Store, hostAddress string) (*CA, error) {
	c := &CA{store: store}

	caCertPEM, err := store.GetConfig("ca_cert_pem")
	if err == sql.ErrNoRows {
		// First run — generate everything
		if err := c.generate(store, hostAddress); err != nil {
			return nil, fmt.Errorf("generate CA: %w", err)
		}
		return c, nil
	}
	if err != nil {
		return nil, fmt.Errorf("load CA cert: %w", err)
	}

	caKeyPEM, err := store.GetConfig("ca_key_pem")
	if err != nil {
		return nil, fmt.Errorf("load CA key: %w", err)
	}
	serverCertPEM, err := store.GetConfig("server_cert_pem")
	if err != nil {
		return nil, fmt.Errorf("load server cert: %w", err)
	}
	serverKeyPEM, err := store.GetConfig("server_key_pem")
	if err != nil {
		return nil, fmt.Errorf("load server key: %w", err)
	}

	c.CACertPEM = caCertPEM
	block, _ := pem.Decode(caCertPEM)
	c.CACert, err = x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse CA cert: %w", err)
	}
	block, _ = pem.Decode(caKeyPEM)
	c.CAKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse CA key: %w", err)
	}
	c.ServerTLS, err = tls.X509KeyPair(serverCertPEM, serverKeyPEM)
	if err != nil {
		return nil, fmt.Errorf("parse server TLS: %w", err)
	}

	// Regenerate server cert if the host address has changed
	serverCert, parseErr := x509.ParseCertificate(c.ServerTLS.Certificate[0])
	needsRegen := parseErr != nil
	if !needsRegen {
		if ip := net.ParseIP(hostAddress); ip != nil {
			found := false
			for _, certIP := range serverCert.IPAddresses {
				if certIP.Equal(ip) { found = true; break }
			}
			needsRegen = !found
		} else {
			found := false
			for _, dns := range serverCert.DNSNames {
				if dns == hostAddress { found = true; break }
			}
			needsRegen = !found
		}
	}
	if needsRegen {
		fmt.Printf("Host address changed to %s — regenerating server certificate\n", hostAddress)
		newServerCertPEM, newServerKeyPEM, err := c.generateServerCert(hostAddress)
		if err != nil {
			return nil, fmt.Errorf("regenerate server cert: %w", err)
		}
		if err := store.SetConfig("server_cert_pem", newServerCertPEM); err != nil {
			return nil, fmt.Errorf("save server cert: %w", err)
		}
		if err := store.SetConfig("server_key_pem", newServerKeyPEM); err != nil {
			return nil, fmt.Errorf("save server key: %w", err)
		}
		c.ServerTLS, err = tls.X509KeyPair(newServerCertPEM, newServerKeyPEM)
		if err != nil {
			return nil, fmt.Errorf("parse new server TLS: %w", err)
		}
	}
	return c, nil
}

func (c *CA) generate(store *db.Store, hostAddress string) error {
	// Generate CA key (RSA 4096)
	caKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}
	caTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "TR-12 Host Service CA", Organization: []string{"TR-12"}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caCertDER, err := x509.CreateCertificate(rand.Reader, caTemplate, caTemplate, &caKey.PublicKey, caKey)
	if err != nil {
		return err
	}
	c.CACert, _ = x509.ParseCertificate(caCertDER)
	c.CAKey = caKey
	c.CACertPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caCertDER})
	caKeyDER, _ := x509.MarshalPKCS8PrivateKey(caKey)
	caKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: caKeyDER})

	// Generate server cert signed by CA
	serverCertPEM, serverKeyPEM, err := c.generateServerCert(hostAddress)
	if err != nil {
		return fmt.Errorf("generate server cert: %w", err)
	}
	c.ServerTLS, err = tls.X509KeyPair(serverCertPEM, serverKeyPEM)
	if err != nil {
		return err
	}

	// Generate JWT secret
	jwtSecret := make([]byte, 32)
	if _, err := rand.Read(jwtSecret); err != nil {
		return err
	}

	// Store everything
	if err := store.SetConfig("ca_cert_pem", c.CACertPEM); err != nil {
		return err
	}
	if err := store.SetConfig("ca_key_pem", caKeyPEM); err != nil {
		return err
	}
	if err := store.SetConfig("server_cert_pem", serverCertPEM); err != nil {
		return err
	}
	if err := store.SetConfig("server_key_pem", serverKeyPEM); err != nil {
		return err
	}
	return store.SetConfig("jwt_secret", jwtSecret)
}

func (c *CA) generateServerCert(hostAddress string) (certPEM, keyPEM []byte, err error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	// Use a fixed service name in the SAN — not the transient IP.
	// The device validates against the CA cert (mutual TLS), so hostname
	// verification is redundant and breaks on IP changes.
	template := &x509.Certificate{
		SerialNumber: big.NewInt(2),
		Subject:      pkix.Name{CommonName: "tr12-host", Organization: []string{"TR-12"}},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{"tr12-host", "localhost"},
	}
	// Also add the current IP so standard TLS clients work without InsecureSkipVerify
	if ip := net.ParseIP(hostAddress); ip != nil {
		template.IPAddresses = []net.IP{ip}
	} else {
		template.DNSNames = append(template.DNSNames, hostAddress)
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, c.CACert, &key.PublicKey, c.CAKey)
	if err != nil {
		return nil, nil, err
	}
	certPEM = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	keyDER, _ := x509.MarshalPKCS8PrivateKey(key)
	keyPEM = pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyDER})
	return certPEM, keyPEM, nil
}

// SignCSR signs a device CSR and returns the device certificate PEM.
func (c *CA) SignCSR(csrPEM []byte, deviceID string, days int) ([]byte, error) {
	block, _ := pem.Decode(csrPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode CSR PEM")
	}
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse CSR: %w", err)
	}
	if err := csr.CheckSignature(); err != nil {
		return nil, fmt.Errorf("invalid CSR signature: %w", err)
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, err
	}
	template := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: deviceID, Organization: []string{"TR-12 Device"}},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(time.Duration(days) * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	certDER, err := x509.CreateCertificate(rand.Reader, template, c.CACert, csr.PublicKey, c.CAKey)
	if err != nil {
		return nil, fmt.Errorf("sign certificate: %w", err)
	}
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER}), nil
}

// TLSConfig returns a TLS config for the MQTT broker with mutual TLS.
func (c *CA) TLSConfig() *tls.Config {
	caPool := x509.NewCertPool()
	caPool.AddCert(c.CACert)
	return &tls.Config{
		Certificates: []tls.Certificate{c.ServerTLS},
		ClientAuth:   tls.RequireAnyClientCert,
		ClientCAs:    caPool,
	}
}

// HTTPTLSConfig returns a TLS config for the HTTP server — no client cert required.
func (c *CA) HTTPTLSConfig() *tls.Config {
	return &tls.Config{
		Certificates: []tls.Certificate{c.ServerTLS},
		MinVersion:   tls.VersionTLS12,
	}
}

// GetJWTSecret loads the JWT signing secret from the database.
func GetJWTSecret(store *db.Store) ([]byte, error) {
	return store.GetConfig("jwt_secret")
}
