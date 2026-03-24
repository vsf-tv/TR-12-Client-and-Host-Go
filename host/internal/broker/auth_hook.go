package broker

import (
	"bytes"
	"crypto/tls"
	"log"
	"time"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/packets"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/db"
)

// AuthHook authenticates MQTT clients via client certificates.
type AuthHook struct {
	mqtt.HookBase
	store *db.Store
}

// ID returns the hook identifier.
func (h *AuthHook) ID() string {
	return "auth-hook"
}

// Provides indicates which hook methods this hook implements.
func (h *AuthHook) Provides(b byte) bool {
	return bytes.Contains([]byte{
		mqtt.OnConnectAuthenticate,
		mqtt.OnDisconnect,
	}, []byte{b})
}

// OnConnectAuthenticate validates the client certificate against known devices.
func (h *AuthHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {
	if cl.Net.Conn == nil {
		return false
	}

	// Check if this is a TLS connection with a client cert
	tlsConn, ok := cl.Net.Conn.(*tls.Conn)
	if !ok {
		// Non-TLS connection (internal inline client) — allow
		return true
	}

	state := tlsConn.ConnectionState()
	if len(state.PeerCertificates) == 0 {
		// No client cert — allow for internal clients
		return true
	}

	deviceID := state.PeerCertificates[0].Subject.CommonName
	device, err := h.store.GetDevice(deviceID)
	if err != nil || device == nil {
		log.Printf("[auth] rejected unknown device %s", deviceID)
		return false
	}
	if device.State == "DEPROVISIONED" {
		log.Printf("[auth] rejected deprovisioned device %s", deviceID)
		return false
	}

	// Update online state
	now := time.Now().UTC().Format(time.RFC3339)
	sourceIP := cl.Net.Conn.RemoteAddr().String()
	h.store.UpdateDeviceOnline(deviceID, true, sourceIP, now)

	// If device has both current and previous certs, revoke previous on connect
	if device.PreviousCertPEM != "" && device.CurrentCertPEM != "" {
		h.store.RevokePreviousCert(deviceID)
	}

	log.Printf("[auth] device %s connected from %s", deviceID, sourceIP)
	return true
}

// OnDisconnect updates device online state.
func (h *AuthHook) OnDisconnect(cl *mqtt.Client, err error, expire bool) {
	if cl.Net.Conn == nil {
		return
	}

	// Use client ID as device ID
	deviceID := string(cl.ID)
	if deviceID == "" {
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	h.store.UpdateDeviceOnline(deviceID, false, "", now)
	log.Printf("[auth] device %s disconnected", deviceID)
}
