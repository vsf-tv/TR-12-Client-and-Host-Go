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
package sdk

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/models"
)

// newTestSDK creates a minimal CddSdk in a temp directory.
func newTestSDK(t *testing.T) *CddSdk {
	t.Helper()
	tmp := t.TempDir()
	sdk, err := New(tmp, "test-device", "SOURCE", tmp, tmp)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	t.Cleanup(func() { sdk.Shutdown() })
	return sdk
}

// writeFakeCerts writes minimal credential files so deleteCredentials has something to remove.
func writeFakeCerts(t *testing.T, certsPath, deviceLocalID, hostID string) string {
	t.Helper()
	dir := filepath.Join(certsPath, deviceLocalID, hostID)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	for _, name := range []string{"ca_cert", "device_cert", "priv_key", "connection_settings", "host_settings"} {
		os.WriteFile(filepath.Join(dir, name), []byte("fake"), 0600)
	}
	return dir
}

// --- Deprovision ---

// TestDeprovision_NotConnected_NoForce: returns "must use --force" without deprovisioning.
func TestDeprovision_NotConnected_NoForce(t *testing.T) {
	sdk := newTestSDK(t)
	// SDK starts DISCONNECTED — not connected to any host
	resp := sdk.Deprovision("some-host", false)
	if !resp.Success {
		t.Fatalf("expected success=true, got false: %s", resp.Message)
	}
	if !strings.Contains(resp.Message, "force") {
		t.Fatalf("expected 'force' in message, got: %s", resp.Message)
	}
}

// TestDeprovision_NotConnected_Force: deletes certs even when not connected.
func TestDeprovision_NotConnected_Force(t *testing.T) {
	sdk := newTestSDK(t)
	certsDir := writeFakeCerts(t, sdk.certsPath, sdk.deviceLocalID, "some-host")

	// Verify certs exist before
	if _, err := os.Stat(certsDir); os.IsNotExist(err) {
		t.Fatal("expected certs dir to exist before deprovision")
	}

	resp := sdk.Deprovision("some-host", true)
	if !resp.Success {
		t.Fatalf("expected success=true, got false: %s", resp.Message)
	}

	// Certs dir should be gone
	if _, err := os.Stat(certsDir); !os.IsNotExist(err) {
		t.Fatal("expected certs dir to be deleted after force deprovision")
	}
}

// TestDeprovision_Force_StateReset: state is DISCONNECTED after deprovision.
func TestDeprovision_Force_StateReset(t *testing.T) {
	sdk := newTestSDK(t)
	writeFakeCerts(t, sdk.certsPath, sdk.deviceLocalID, "some-host")

	sdk.Deprovision("some-host", true)

	if sdk.state != models.StateDisconnected {
		t.Fatalf("expected DISCONNECTED after deprovision, got %s", sdk.state)
	}
}

// TestDeprovision_WrongHost_NoOp: deprovisioning a different host than connected is a no-op.
func TestDeprovision_WrongHost_NoOp(t *testing.T) {
	sdk := newTestSDK(t)
	// Manually set hostID to simulate being connected to "host-A"
	sdk.hostID = "host-A"
	sdk.state = models.StateConnected

	certsDir := writeFakeCerts(t, sdk.certsPath, sdk.deviceLocalID, "host-A")

	// Deprovision "host-B" with force — should not delete host-A's certs
	// because deleteCredentials only deletes when hostID matches or is empty
	sdk.Deprovision("host-B", true)

	// host-A certs should still exist (deleteCredentials is a no-op for mismatched host)
	if _, err := os.Stat(certsDir); os.IsNotExist(err) {
		t.Fatal("expected host-A certs to survive deprovision of host-B")
	}
}

// --- connectedTo ---

func TestConnectedTo_True(t *testing.T) {
	sdk := newTestSDK(t)
	sdk.state = models.StateConnected
	sdk.hostID = "my-host"
	if !sdk.connectedTo("my-host") {
		t.Fatal("expected connectedTo=true")
	}
}

func TestConnectedTo_WrongHost(t *testing.T) {
	sdk := newTestSDK(t)
	sdk.state = models.StateConnected
	sdk.hostID = "my-host"
	if sdk.connectedTo("other-host") {
		t.Fatal("expected connectedTo=false for wrong host")
	}
}

func TestConnectedTo_NotConnected(t *testing.T) {
	sdk := newTestSDK(t)
	sdk.state = models.StateDisconnected
	sdk.hostID = "my-host"
	if sdk.connectedTo("my-host") {
		t.Fatal("expected connectedTo=false when DISCONNECTED")
	}
}

// --- State machine helpers ---

func TestIs_SingleState(t *testing.T) {
	sdk := newTestSDK(t)
	sdk.state = models.StateConnected
	if !sdk.is(models.StateConnected) {
		t.Fatal("expected is(CONNECTED)=true")
	}
	if sdk.is(models.StateDisconnected) {
		t.Fatal("expected is(DISCONNECTED)=false")
	}
}

func TestIs_MultipleStates(t *testing.T) {
	sdk := newTestSDK(t)
	sdk.state = models.StateReconnecting
	if !sdk.is(models.StateConnected, models.StateReconnecting) {
		t.Fatal("expected is(CONNECTED, RECONNECTING)=true when RECONNECTING")
	}
}
