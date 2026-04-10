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
package pairing

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/credentials"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/models"
	tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
)

// newTestStore creates a credential store in a temp directory.
func newTestStore(t *testing.T) *credentials.Store {
	t.Helper()
	store, err := credentials.NewStore(t.TempDir(), "test-device", "test-host")
	if err != nil {
		t.Fatalf("NewStore: %v", err)
	}
	// Pre-generate keys so pairing doesn't need to do it
	if err := store.GenerateKeysAndCSR(); err != nil {
		t.Fatalf("GenerateKeysAndCSR: %v", err)
	}
	return store
}

// newTestPairing creates a Pairing pointed at the given server URL.
func newTestPairing(t *testing.T, serverURL string) *Pairing {
	t.Helper()
	return New(newTestStore(t), "SOURCE", "test-host", serverURL, serverURL)
}

// pairSuccessResponse builds a valid pair success JSON response.
func pairSuccessResponse(deviceID, pairingCode, accessCode string, timeoutSec int) []byte {
	timeout := float32(timeoutSec)
	resp := models.CreatePairingCodeResponseContent{
		Result: models.CreatePairingCodeResult{
			Success: &tr12models.Success{
				Success: tr12models.CreatePairingCodeSuccessData{
					DeviceId:              deviceID,
					PairingCode:           pairingCode,
					AccessCode:            accessCode,
					PairingTimeoutSeconds: timeout,
				},
			},
		},
	}
	b, _ := json.Marshal(resp)
	return b
}

// pairFailureResponse builds a pair failure JSON response.
func pairFailureResponse(reason tr12models.CreatePairingCodeFailureReason) []byte {
	resp := models.CreatePairingCodeResponseContent{
		Result: models.CreatePairingCodeResult{
			Failure: &tr12models.Failure{
				Failure: tr12models.CreatePairingCodeFailureData{Reason: reason},
			},
		},
	}
	b, _ := json.Marshal(resp)
	return b
}

// authResponse builds an authenticate response JSON.
func authResponse(status tr12models.AuthStatus, mqttURI, region string) []byte {
	caCert := "-----BEGIN CERTIFICATE-----\nfake-ca\n-----END CERTIFICATE-----\n"
	deviceCert := "-----BEGIN CERTIFICATE-----\nfake-device\n-----END CERTIFICATE-----\n"
	subTopic := "cdd/dev1/config/update"
	pubTopic := "cdd/dev1/registration/report"
	statusTopic := "cdd/dev1/status/report"
	actualTopic := "cdd/dev1/config/actual/report"
	schemaTopic := "cdd/dev1/schema/report"
	certsTopic := "cdd/dev1/certs/update"
	deprovTopic := "cdd/dev1/deprovision"
	pubDeprovTopic := "cdd/dev1/deprovision/ack"
	thumbTopic := "cdd/dev1/thumbnail/subscription"
	logTopic := "cdd/dev1/log/subscription"
	proto := "mqtt"
	timeout := float32(300)
	interval := float32(1)
	keepalive := float32(30)

	hs := &tr12models.HostSettings{
		IotProtocolName:                    proto,
		PairingTimeoutSeconds:              timeout,
		MinIntervalPubSeconds:              interval,
		MqttKeepaliveSeconds:               keepalive,
		SubUpdateTopic:                     subTopic,
		PubReportRegistrationTopic:         pubTopic,
		PubReportStatusTopic:               statusTopic,
		PubReportActualConfigurationTopic:  actualTopic,
		PubReportSchemaTopic:               schemaTopic,
		SubUpdateCertsTopic:                certsTopic,
		SubDeprovisionTopic:                deprovTopic,
		PubDeprovisionTopic:                pubDeprovTopic,
		SubUpdateThumbnailSubscriptionTopic: thumbTopic,
		SubUpdateLogSubscriptionTopic:      logTopic,
	}

	resp := tr12models.AuthenticatePairingCodeResponseContent{
		Status:        status,
		CaCertificate: &caCert,
		DeviceCertificate: &deviceCert,
		MqttUri:       &mqttURI,
		RegionName:    &region,
		HostSettings:  hs,
	}
	b, _ := json.Marshal(resp)
	return b
}

// --- IsExpired / ExpiresIn / GetPairingCode ---

func TestIsExpired_NoPairResponse(t *testing.T) {
	p := &Pairing{}
	if p.IsExpired() {
		t.Fatal("expected not expired when no pair response")
	}
}

func TestIsExpired_NotExpired(t *testing.T) {
	timeout := float32(300)
	p := &Pairing{
		StartTime: time.Now().Unix(),
		PairResponse: &models.CreatePairingCodeResponseContent{
			Result: models.CreatePairingCodeResult{
				Success: &tr12models.Success{
					Success: tr12models.CreatePairingCodeSuccessData{
						PairingTimeoutSeconds: timeout,
					},
				},
			},
		},
	}
	if p.IsExpired() {
		t.Fatal("expected not expired within timeout")
	}
}

func TestIsExpired_Expired(t *testing.T) {
	timeout := float32(1)
	p := &Pairing{
		StartTime: time.Now().Unix() - 10, // 10 seconds ago
		PairResponse: &models.CreatePairingCodeResponseContent{
			Result: models.CreatePairingCodeResult{
				Success: &tr12models.Success{
					Success: tr12models.CreatePairingCodeSuccessData{
						PairingTimeoutSeconds: timeout,
					},
				},
			},
		},
	}
	if !p.IsExpired() {
		t.Fatal("expected expired")
	}
}

func TestExpiresIn(t *testing.T) {
	timeout := float32(300)
	p := &Pairing{
		StartTime: time.Now().Unix() - 10,
		PairResponse: &models.CreatePairingCodeResponseContent{
			Result: models.CreatePairingCodeResult{
				Success: &tr12models.Success{
					Success: tr12models.CreatePairingCodeSuccessData{
						PairingTimeoutSeconds: timeout,
					},
				},
			},
		},
	}
	remaining := p.ExpiresIn()
	if remaining < 285 || remaining > 295 {
		t.Fatalf("expected ~290s remaining, got %d", remaining)
	}
}

func TestGetPairingCode_NoResponse(t *testing.T) {
	p := &Pairing{}
	if p.GetPairingCode() != "" {
		t.Fatal("expected empty pairing code when no response")
	}
}

func TestGetPairingCode(t *testing.T) {
	timeout := float32(300)
	p := &Pairing{
		PairResponse: &models.CreatePairingCodeResponseContent{
			Result: models.CreatePairingCodeResult{
				Success: &tr12models.Success{
					Success: tr12models.CreatePairingCodeSuccessData{
						PairingCode:           "ABC123",
						PairingTimeoutSeconds: timeout,
					},
				},
			},
		},
	}
	if got := p.GetPairingCode(); got != "ABC123" {
		t.Fatalf("expected ABC123, got %q", got)
	}
}

// --- GetNewPairingCode ---

func TestGetNewPairingCode_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/pair" {
			w.Header().Set("Content-Type", "application/json")
			w.Write(pairSuccessResponse("dev-001", "XYZ789", "secret", 300))
		}
	}))
	defer srv.Close()

	p := newTestPairing(t, srv.URL)
	if err := p.GetNewPairingCode(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.GetPairingCode() != "XYZ789" {
		t.Fatalf("expected pairing code XYZ789, got %q", p.GetPairingCode())
	}
	if p.ExpiresIn() <= 0 {
		t.Fatal("expected positive ExpiresIn after successful pair")
	}
}

func TestGetNewPairingCode_HostIDMismatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(pairFailureResponse(tr12models.HOST_ID_MISMATCH))
	}))
	defer srv.Close()

	p := newTestPairing(t, srv.URL)
	err := p.GetNewPairingCode()
	if err == nil {
		t.Fatal("expected error for HOST_ID_MISMATCH")
	}
	if !strings.Contains(err.Error(), "host ID") {
		t.Fatalf("expected 'host ID' in error, got: %v", err)
	}
}

func TestGetNewPairingCode_VersionNotSupported(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(pairFailureResponse(tr12models.VERSION_NOT_SUPPORTED))
	}))
	defer srv.Close()

	p := newTestPairing(t, srv.URL)
	err := p.GetNewPairingCode()
	if err == nil || !strings.Contains(err.Error(), "version") {
		t.Fatalf("expected version error, got: %v", err)
	}
}

func TestGetNewPairingCode_DeviceTypeNotSupported(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(pairFailureResponse(tr12models.DEVICE_TYPE_NOT_SUPPORTED))
	}))
	defer srv.Close()

	p := newTestPairing(t, srv.URL)
	err := p.GetNewPairingCode()
	if err == nil || !strings.Contains(err.Error(), "device type") {
		t.Fatalf("expected device type error, got: %v", err)
	}
}

func TestGetNewPairingCode_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	p := newTestPairing(t, srv.URL)
	err := p.GetNewPairingCode()
	if err == nil || !strings.Contains(err.Error(), "500") {
		t.Fatalf("expected HTTP 500 error, got: %v", err)
	}
}

func TestGetNewPairingCode_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`not valid json`))
	}))
	defer srv.Close()

	p := newTestPairing(t, srv.URL)
	err := p.GetNewPairingCode()
	if err == nil || !strings.Contains(err.Error(), "not valid") {
		t.Fatalf("expected parse error, got: %v", err)
	}
}

func TestGetNewPairingCode_ConnectionRefused(t *testing.T) {
	p := newTestPairing(t, "http://127.0.0.1:1") // nothing listening
	err := p.GetNewPairingCode()
	if err == nil {
		t.Fatal("expected connection error")
	}
}

// --- AuthenticatePairingCode ---

func TestAuthenticatePairingCode_NoPairingCode(t *testing.T) {
	p := &Pairing{Certs: newTestStore(t)}
	_, err := p.AuthenticatePairingCode()
	if err == nil || !strings.Contains(err.Error(), "no pairing code") {
		t.Fatalf("expected 'no pairing code' error, got: %v", err)
	}
}

func TestAuthenticatePairingCode_Standby(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(authResponse(tr12models.STANDBY, "", ""))
	}))
	defer srv.Close()

	p := newTestPairing(t, srv.URL)
	// Inject a pair response so AuthenticatePairingCode has a pairing code to use
	timeout := float32(300)
	p.PairResponse = &models.CreatePairingCodeResponseContent{
		Result: models.CreatePairingCodeResult{
			Success: &tr12models.Success{
				Success: tr12models.CreatePairingCodeSuccessData{
					DeviceId:              "dev-001",
					PairingCode:           "ABC123",
					AccessCode:            "secret",
					PairingTimeoutSeconds: timeout,
				},
			},
		},
	}

	claimed, err := p.AuthenticatePairingCode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if claimed {
		t.Fatal("expected not claimed for STANDBY")
	}
}

func TestAuthenticatePairingCode_Claimed(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(authResponse(tr12models.CLAIMED, "tls://127.0.0.1:8883", "local"))
	}))
	defer srv.Close()

	p := newTestPairing(t, srv.URL)
	timeout := float32(300)
	p.PairResponse = &models.CreatePairingCodeResponseContent{
		Result: models.CreatePairingCodeResult{
			Success: &tr12models.Success{
				Success: tr12models.CreatePairingCodeSuccessData{
					DeviceId:              "dev-001",
					PairingCode:           "ABC123",
					AccessCode:            "secret",
					PairingTimeoutSeconds: timeout,
				},
			},
		},
	}

	claimed, err := p.AuthenticatePairingCode()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !claimed {
		t.Fatal("expected claimed=true for CLAIMED status")
	}
	// Verify certs were written to filesystem
	if p.Certs.GetDeviceID() != "dev-001" {
		t.Fatalf("expected device ID dev-001, got %q", p.Certs.GetDeviceID())
	}
	if p.Certs.GetRegion() != "local" {
		t.Fatalf("expected region local, got %q", p.Certs.GetRegion())
	}
}

func TestAuthenticatePairingCode_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "server error", http.StatusInternalServerError)
	}))
	defer srv.Close()

	p := newTestPairing(t, srv.URL)
	timeout := float32(300)
	p.PairResponse = &models.CreatePairingCodeResponseContent{
		Result: models.CreatePairingCodeResult{
			Success: &tr12models.Success{
				Success: tr12models.CreatePairingCodeSuccessData{
					DeviceId:              "dev-001",
					PairingCode:           "ABC123",
					AccessCode:            "secret",
					PairingTimeoutSeconds: timeout,
				},
			},
		},
	}

	_, err := p.AuthenticatePairingCode()
	if err == nil || !strings.Contains(err.Error(), "500") {
		t.Fatalf("expected HTTP 500 error, got: %v", err)
	}
}
