# TR-12 Client SDK and Application Reference Design — Go

A Go implementation of the TR-12 Client Device Discovery SDK, providing discovery, monitoring, and connection management of streaming video devices using an internet-secure, cloud and NAT friendly, scalable pairing and communication protocol.

This is a full port of the [Python CDD SDK](https://github.com/vsf-tv/gccg-cdd) with identical CLI arguments and REST API surface.

## TR-12 Working Group

> Draft design documents related to this project are currently being discussed and revised in the VSF Bi-Weekly Forum.
> For access, please reach out to Brad Gilmer <brad@gilmer.tv> or Brian Rundle <brundle@amazon.com>.

## Quick Start

### 1. Build

```bash
cd client

# Build the SDK daemon
go build -o bin/cdd-sdk ./cmd/cdd-sdk/

# Build the Application Reference Design
go build -o bin/ard ./cmd/application_reference_design/
```

### 2. Start the SDK

```bash
mkdir -p /tmp/cdd_certs /tmp/cdd_logs

./bin/cdd-sdk \
  --internal_device_id my_device_123 \
  --certs_path /tmp/cdd_certs \
  --log_path /tmp/cdd_logs \
  --ip 127.0.0.1 \
  --port 8603 \
  --device_type SOURCE
```

### 3. Start the ARD

In a second terminal:

```bash
# 1-channel encoder against the local TR-12 Host Service:
./bin/ard --host_id tr12-host --registration_file cmd/application_reference_design/payloads/1_channel_encoder/registration.json

# 2-channel encoder:
./bin/ard --host_id tr12-host --registration_file cmd/application_reference_design/payloads/2_channel_encoder/registration.json
```

The ARD will display a pairing code. Claim the device on the host service (see the [Host README](../host/README.md)), then the ARD transitions to `CONNECTED` and begins reporting status and configuration.

### 4. Cross-Compilation

Go's built-in cross-compilation requires no additional toolchain:

```bash
# Linux x86_64 (EC2, server)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/cdd-sdk-linux-amd64 ./cmd/cdd-sdk/

# Linux ARM64 (embedded devices, Raspberry Pi, PetaLinux)
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/cdd-sdk-arm64 ./cmd/cdd-sdk/
```

The resulting binary is fully self-contained — no runtime dependencies, no Python interpreter, no pip packages.

---

## Architecture

The SDK runs as a standalone process hosting a REST API on localhost. A device application (like the ARD) makes API requests to the SDK. The SDK handles pairing/auth via HTTPS and real-time communication with the host service via MQTT/TLS.

```
┌──────────────────┐       REST API        ┌──────────────┐     MQTT/TLS      ┌──────────────────┐
│  Device App /    │ ───────────────────── │   CDD SDK    │ ────────────────── │  TR-12 Host      │
│  ARD             │   localhost:port      │   (Go)       │   Port 443        │  Service         │
└──────────────────┘                       └──────────────┘                    └──────────────────┘
```

### State Machine

```
DISCONNECTED → PAIRING → CONNECTING → CONNECTED
                                          ↓
                                     RECONNECTING
```

---

## SDK CLI Arguments

| Argument | Description |
|---|---|
| `--internal_device_id` | Unique device identifier (required) |
| `--certs_path` | Persistent directory for X.509 credential storage (required) |
| `--log_path` | Writable directory for JSON log files (required) |
| `--ip` | IP address for the local REST API (required) |
| `--port` | Port for the local REST API (required) |
| `--device_type` | Device type: `SOURCE`, `DESTINATION`, or `BOTH` (required) |

## ARD CLI Arguments

| Argument | Description |
|---|---|
| `--host_id` | Host ID to connect to (required) |
| `--sdk_url` | Base URL of the running SDK (default: `http://127.0.0.1:8603`) |
| `--registration_file` | Path to registration JSON file (default: `payloads/1_channel_encoder/registration.json`) |

---

## Application Reference Design (ARD)

The ARD is a companion program that simulates an encoder device making REST calls on the SDK daemon. It demonstrates the full TR-12 lifecycle and serves as the template for real device integrations.

### What the ARD Does

- Calls `PUT /connect` in a loop until the SDK reaches `CONNECTED` state
- Displays the pairing code when in `PAIRING` state
- Once connected, polls `GET /get_configuration` for host-service updates
- Applies desired configuration via the TR-12 shim (settings, connection, state)
- Reports device status via `PUT /report_status`
- Reports actual configuration via `PUT /report_actual_configuration` (including `thumbnailLocalPath` per channel)
- Manages an ffmpeg subprocess for SRT streaming when the host sets channel state to `ACTIVE`

### DeviceCallbacks Interface

The core integration pattern is the `DeviceCallbacks` interface (`internal/application_reference_design/callbacks_interface.go`). Any device integration implements this interface — the `ApplicationLoop` and `Tr12Shim` only know the interface, never the concrete implementation.

**Apply (set) side** — called when the host sends a new desired configuration:
- `UpdateDeviceKeyValue` — device-level setting (e.g. clock source)
- `UpdateChannelSettings` — channel simple setting (e.g. framerate, codec, bitrate)
- `UpdateChannelProfile` — profile selection
- `UpdateChannelConnection` — transport protocol config (SRT caller/listener, RIST, etc.)
- `UpdateChannelState` — ACTIVE or IDLE

**Read-back (get) side** — called when building the actual configuration to report back:
- `GetDeviceUpdatedValue`, `GetChannelUpdatedValue`, `GetChannelProfileValue`
- `GetChannelConnection`, `GetChannelState`
- `GetDeviceHealth`, `GetChannelHealth`
- `GetDeviceStatus`, `GetChannelStatus`
- `GetChannelThumbnailLocalPath` — returns the local filesystem path to the thumbnail image for a channel based on the device's current input configuration

**One reference implementation ships with this repo:**

`ArdCallbacks` (`internal/application_reference_design/mock_encoder_callbacks.go`) — the reference implementation using a simulated ffmpeg encoder. Study this to understand the pattern.

Any real device integration implements the same interface and wires into `ApplicationLoop`:

```go
// ARD binary — mock encoder:
callbacks := ard.NewArdCallbacks()
loop := ard.NewApplicationLoop(sdkURL, callbacks, &registration)

// Real device — same pattern, different callbacks:
callbacks := device.NewMyDeviceCallbacks(deviceURL)
loop := ard.NewApplicationLoop(sdkURL, callbacks, &registration)
```

### Thumbnail Model (v2.0.2)

Thumbnails are channel-centric. The host subscribes by `channelId`, and the SDK resolves the local file path from the latest `ActualConfiguration` reported by the application.

- The application sets `thumbnailLocalPath` per channel in ActualConfiguration via `GetChannelThumbnailLocalPath`
- The ARD defaults to SDI1 (`/tmp/image_sdi.jpg`) when no configuration has been pushed yet
- When the host changes a channel's input (e.g. SDI1 → HDMI1), the application updates the path in the next ActualConfiguration report
- The SDK picks up the new path on its next upload cycle — no restart needed

### ARD Data Files

The ARD expects these files relative to its binary location:

- `payloads/1_channel_encoder/registration.json` — 1-channel device registration
- `payloads/2_channel_encoder/registration.json` — 2-channel device registration
- `thumbnails/thumbnail_images_sdi/` — SDI thumbnail source images
- `thumbnails/thumbnail_images_hdmi/` — HDMI thumbnail source images

---

## SDK REST API Reference

All endpoints are served on `http://<ip>:<port>`.

### PUT /connect

Initiates or continues the connection/pairing flow with a host service.

```json
{
  "hostId": "tr12-host",
  "registration": { ... }
}
```

Response includes `state`, `pairingCode` (if pairing), `deviceId`, and `regionName` (if connected).

### PUT /disconnect

Disconnects from the host service and resets state to `DISCONNECTED`.

### GET /get_state

Returns the current connection state.

```json
{
  "success": true,
  "state": "CONNECTED",
  "message": ""
}
```

### PUT /report_status

Publishes a device status payload to the host service.

```json
{
  "status": { ... }
}
```

### PUT /report_actual_configuration

Publishes the device's actual configuration to the host service. The SDK stores this internally for thumbnail path resolution.

```json
{
  "configuration": { ... }
}
```

### GET /get_configuration

Returns the latest desired configuration received from the host service.

### PUT /deprovision

Removes the device from the host service and deletes local credentials.

```json
{
  "hostId": "tr12-host"
}
```

Use `?force=true` to deprovision while not connected.

---

## Host Configuration

Host configuration files are JSON files in the `host_configuration/` directory, named `<host_id>.json`. The SDK looks for this directory relative to the binary location, falling back to the working directory.

Example for the local TR-12 Host Service:

```json
{
  "serviceId": "tr12-host",
  "serviceName": "My TR-12 Host",
  "deviceTypes": ["SOURCE", "DESTINATION", "BOTH"],
  "createPairingCodeUrl": "http://127.0.0.1:8080",
  "authenticatePairingCodeUrl": "http://127.0.0.1:8080",
  "thumbnailMaximumSizeKB": 100,
  "logFileMaximumSizeKB": 500
}
```

## Credential Storage

Credentials are stored at `<certs_path>/<internal_device_id>/<host_id>/`:

- Use the same `--internal_device_id` to reconnect with existing credentials
- Use a different `--internal_device_id` to start fresh
- Each host ID gets its own subfolder, so one device identity can connect to multiple hosts

## Logging

The SDK writes JSON-formatted rotating log files to the `--log_path` directory. Log files are capped at 500 KB with up to 3 rotated backups. Logs are also printed to stdout.

When the host service requests log uploads, the SDK automatically uploads rotated log files to the provided pre-signed URL.

---

## Requirements

- Go 1.22 or newer
- Outbound HTTPS access on port 443
- Persistent read/write filesystem for credential storage

## Dependencies

Managed via Go modules (`go.mod`):

- `github.com/gin-gonic/gin` — HTTP server
- `github.com/gin-contrib/cors` — CORS middleware
- `github.com/eclipse/paho.mqtt.golang` — MQTT client

## Project Structure

```
client/
├── cmd/
│   ├── cdd-sdk/main.go                          # SDK entry point
│   ├── application_reference_design/             # ARD entry point + payloads
│   │   ├── main.go
│   │   ├── payloads/
│   │   │   ├── 1_channel_encoder/registration.json
│   │   │   └── 2_channel_encoder/registration.json
│   │   └── thumbnails/
│   │       ├── thumbnail_images_sdi/
│   │       └── thumbnail_images_hdmi/
├── host_configuration/                           # Host config JSON files
├── internal/
│   ├── api/server.go                             # Gin REST API server
│   ├── application_reference_design/
│   │   ├── application.go                        # ARD initialization
│   │   ├── application_loop.go                   # Reusable TR-12 lifecycle loop
│   │   ├── callbacks_interface.go                # DeviceCallbacks interface
│   │   ├── mock_encoder_callbacks.go             # ArdCallbacks implementation
│   │   ├── mock_encoder_device.go                # ffmpeg encoder simulation
│   │   ├── tr12_client_caller.go                 # SDK REST client
│   │   ├── tr12_model_traversal.go               # TR-12 model shim
│   │   └── thumbnail_simulator.go                # Thumbnail image cycling
│   ├── cddlogger/logger.go                       # JSON rotating file logger
│   ├── credentials/store.go                      # X.509 cert persistence
│   ├── models/tr12_models.go                     # State constants + model aliases
│   ├── pairing/pairing.go                        # Pairing/auth flow
│   ├── sdk/
│   │   ├── sdk.go                                # Core SDK struct and state machine
│   │   ├── connect.go                            # Public API methods
│   │   └── mqtt.go                               # MQTT connection and callbacks
│   ├── thumbnails/manager.go                     # Thumbnail upload manager
│   └── utils/utils.go                            # TLS, upload, throttle, key gen
├── go.mod
└── go.sum
```

## TR-12 Protocol Reference

- Smithy Models: https://github.com/vsf-tv/TR-12-Models
- Draft Protocol: https://github.com/vsf-tv/TR-12-Models/blob/main/VSF_TR-12-ClientDeviceDiscoveryDraft.pdf

## License

Apache License, Version 2.0
