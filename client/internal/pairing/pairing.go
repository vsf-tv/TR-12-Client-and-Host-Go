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
// Pairing manages the TR-12 pairing and authentication flow with the host service.
package pairing

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/credentials"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/models"
	tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
)

const maxTimeoutSec = 5

// Pairing manages the pairing process with the host service.
type Pairing struct {
	Certs        *credentials.Store
	DeviceType   string
	HostID       string
	PairingURL   string
	AuthURL      string
	StartTime    int64
	PairResponse *models.CreatePairingCodeResponseContent
	AuthResponse *models.AuthenticatePairingCodeResponseContent
	httpClient   *http.Client
}

// New creates a new Pairing instance.
func New(certs *credentials.Store, deviceType, hostID, pairingURL, authURL string) *Pairing {
	// Skip TLS verification for pairing/auth — the CA cert is received as part of
	// the auth response, so we can't verify it beforehand. The device cert received
	// is then used for all subsequent MQTT connections which DO verify against the CA.
	httpClient := &http.Client{
		Timeout: maxTimeoutSec * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	return &Pairing{
		Certs:      certs,
		DeviceType: deviceType,
		HostID:     hostID,
		PairingURL: pairingURL,
		AuthURL:    authURL,
		StartTime:  time.Now().Unix(),
		httpClient: httpClient,
	}
}

// getSuccessData returns the CreatePairingCodeSuccessData if available.
func (p *Pairing) getSuccessData() *tr12models.CreatePairingCodeSuccessData {
	if p.PairResponse == nil || p.PairResponse.Result.Success == nil {
		return nil
	}
	return &p.PairResponse.Result.Success.Success
}

// IsExpired returns true if the pairing code has expired.
func (p *Pairing) IsExpired() bool {
	sd := p.getSuccessData()
	if sd == nil {
		return false
	}
	elapsed := time.Now().Unix() - p.StartTime
	return float64(elapsed) > float64(sd.PairingTimeoutSeconds)
}

// ExpiresIn returns seconds until the pairing code expires.
func (p *Pairing) ExpiresIn() int {
	sd := p.getSuccessData()
	if sd == nil {
		return 0
	}
	remaining := float64(sd.PairingTimeoutSeconds) - float64(time.Now().Unix()-p.StartTime)
	if remaining < 0 {
		return 0
	}
	return int(remaining)
}

// GetPairingCode returns the current pairing code.
func (p *Pairing) GetPairingCode() string {
	sd := p.getSuccessData()
	if sd != nil {
		return sd.PairingCode
	}
	return ""
}

// GetNewPairingCode requests a new pairing code from the host service.
func (p *Pairing) GetNewPairingCode() error {
	if err := p.Certs.GenerateKeysAndCSR(); err != nil {
		return fmt.Errorf("failed to generate keys: %w", err)
	}
	reqBody := models.CreatePairingCodeRequestContent{
		DeviceType:                tr12models.DeviceType(p.DeviceType),
		HostId:                    p.HostID,
		CertificateSigningRequest: p.Certs.CSR,
		Version:                   tr12models.ProtocolVersion{Version: tr12models.PtrString(models.ProtocolVersionString)},
	}
	body, _ := json.Marshal(reqBody)
	log.Printf("[PAIR] POST %s/pair  body=%s", p.PairingURL, string(body))
	resp, err := p.httpClient.Post(p.PairingURL+"/pair", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("pairing unable to connect to the service: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[PAIR] Response status=%d body=%s", resp.StatusCode, string(respBody))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("pairing API error - StatusCode: %d - Response: %s", resp.StatusCode, string(respBody))
	}
	var pairResp models.CreatePairingCodeResponseContent
	if err := json.Unmarshal(respBody, &pairResp); err != nil {
		return fmt.Errorf("pairing service response was not valid: %w", err)
	}
	p.PairResponse = &pairResp
	p.StartTime = time.Now().Unix()

	// Check for failure
	if pairResp.Result.Failure != nil {
		reason := pairResp.Result.Failure.Failure.Reason
		switch reason {
		case tr12models.VERSION_NOT_SUPPORTED:
			return fmt.Errorf("TR-12 version not supported: %s", reason)
		case tr12models.DEVICE_TYPE_NOT_SUPPORTED:
			return fmt.Errorf("device type not supported: %s", reason)
		case tr12models.HOST_ID_MISMATCH:
			return fmt.Errorf("host ID does not match the host endpoint: %s", reason)
		default:
			return fmt.Errorf("unknown pairing failure: %s", reason)
		}
	}
	if pairResp.Result.Success == nil {
		return fmt.Errorf("unexpected pairing response: no success or failure data")
	}
	return nil
}

// AuthenticatePairingCode polls the auth endpoint. Returns true if claimed and certs written.
func (p *Pairing) AuthenticatePairingCode() (bool, error) {
	sd := p.getSuccessData()
	if sd == nil {
		return false, fmt.Errorf("no pairing code to authenticate")
	}
	reqBody := models.AuthenticatePairingCodeRequestContent{
		DeviceId:    sd.DeviceId,
		PairingCode: sd.PairingCode,
		AccessCode:  sd.AccessCode,
	}
	body, _ := json.Marshal(reqBody)
	log.Printf("[AUTH] POST %s/authenticate  body=%s", p.AuthURL, string(body))
	resp, err := p.httpClient.Post(p.AuthURL+"/authenticate", "application/json", bytes.NewReader(body))
	if err != nil {
		return false, fmt.Errorf("auth unable to connect to the service: %w", err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	log.Printf("[AUTH] Response status=%d body=%s", resp.StatusCode, string(respBody))
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("auth API error - StatusCode: %d - Response: %s", resp.StatusCode, string(respBody))
	}
	var authResp models.AuthenticatePairingCodeResponseContent
	if err := json.Unmarshal(respBody, &authResp); err != nil {
		return false, fmt.Errorf("auth service response was not valid: %w", err)
	}
	p.AuthResponse = &authResp
	log.Printf("[AUTH] Parsed: status=%s mqttUri=%q regionName=%q hasHostSettings=%v hasCaCertificate=%v hasDeviceCertificate=%v",
		authResp.Status, authResp.GetMqttUri(), authResp.GetRegionName(),
		authResp.HasHostSettings(), authResp.HasCaCertificate(), authResp.HasDeviceCertificate())

	switch authResp.Status {
	case tr12models.STANDBY:
		return false, nil
	case tr12models.CLAIMED:
		log.Printf("[AUTH] Device CLAIMED — writing certs to filesystem for deviceId=%s", sd.DeviceId)
		if err := p.Certs.WriteToFilesystem(sd.DeviceId, &authResp); err != nil {
			return false, fmt.Errorf("unable to write certs to disk: %w", err)
		}
		return true, nil
	default:
		return false, fmt.Errorf("unexpected auth status: %s", authResp.Status)
	}
}
