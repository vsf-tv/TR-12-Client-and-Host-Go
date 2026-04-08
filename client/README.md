# Client Device Discovery (CDD) Client SDK: TR-12 — Go Implementation

A Go implementation of the TR-12 Client Device Discovery SDK, providing discovery, monitoring, and connection management of streaming video devices using an internet-secure, cloud and NAT friendly, scalable pairing and communication protocol.

This is a full port of the [Python CDD SDK](https://github.com/vsf-tv/gccg-cdd) with identical CLI arguments and REST API surface.

## TR-12 Working Group

> Draft design documents related to this project are currently being discussed and revised in the VSF Bi-Weekly Forum.
> For access, please reach out to Brad Gilmer <brad@gilmer.tv> or Brian Rundle <brundle@amazon.com>.

## Architecture

The Go SDK runs as a standalone process hosting a REST API on localhost. A device application uses the TR-12 models to make API requests to the SDK process. The SDK handles connecting to and communicating with the host service via HTTPS (pairing/auth) and MQTT (pub/sub).

```
┌──────────────────┐       REST API        ┌──────────────┐     MQTT/TLS      ┌──────────────────┐
│  Device App /    │ ───────────────────── │   CDD SDK    │ ────────────────── │  TR-12 Host      │
│  Application     │   localhost:port      │   (Go)       │   Port 443        │  Service         │
└──────────────────┘                       └──────────────┘                    └──────────────────┘
```

## State Machine

```
DISCONNECTED → PAIRING → CONNECTING → CONNECTED
                                          ↓
                                     RECONNECTING
```

## Requirements

- Go 1.22 or newer
- Outbound HTTPS access on port 443
- Persistent read/write filesystem for credential storage

## Dependencies

Managed via Go modules (`go.mod`):

- `github.com/gin-gonic/gin` — HTTP server
- `github.com/gin-contrib/cors` — CORS middleware
- `github.com/eclipse/paho.mqtt.golang` — MQTT client

## Build

```bash
cd client
go build -o bin/cdd-sdk ./cmd/cdd-sdk/
```

This produces a single static binary `bin/cdd-sdk` (~13 MB).

### Cross-Compilation

Go's built-in cross-compilation requires no additional toolchain — just set `GOOS`, `GOARCH`, and `CGO_ENABLED=0`:

```bash
# Linux x86_64 (EC2, server)
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/cdd-sdk-linux-amd64 ./cmd/cdd-sdk/

# Linux ARM64 (embedded devices, Raspberry Pi, PetaLinux)
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o bin/cdd-sdk-arm64 ./cmd/cdd-sdk/
```

The resulting binary is fully self-contained — no runtime dependencies, no Python interpreter, no pip packages. Copy it to the target device and run it directly.

## Usage

```bash
./bin/cdd-sdk \
  --internal_device_id <device_name> \
  --certs_path <persistent_cert_directory> \
  --log_path <writable_log_directory> \
  --ip <listen_ip> \
  --port <listen_port> \
  --device_type <SOURCE|DESTINATION|BOTH>
```

### CLI Arguments

| Argument | Description |
|---|---|
| `--internal_device_id` | Unique device identifier (required) |
| `--certs_path` | Persistent directory for X.509 credential storage (required) |
| `--log_path` | Writable directory for JSON log files (required) |
| `--ip` | IP address for the local REST API (required) |
| `--port` | Port for the local REST API (required) |
| `--device_type` | Device type: `SOURCE`, `DESTINATION`, or `BOTH` (required) |

### Example

```bash
export CERTS_PATH="$HOME/cdd_certs"
mkdir -p $CERTS_PATH
export ID="my_device_123"

mkdir -p $CERTS_PATH /tmp/cdd_logs

./bin/cdd-sdk --internal_device_id $ID -certs_path $CERTS_PATH --log_path /tmp/cdd_logs --ip 127.0.0.1 --port 8603 --device_type SOURCE
```

## REST API Endpoints

All endpoints are served on `http://<ip>:<port>`.

### PUT /connect

Initiates or continues the connection/pairing flow with a host service.

```json
{
  "hostId": "vsf_test_host",
  "registration": { ... }
}
```

Response includes `state`, `pairingCode` (if pairing), `deviceId`, and `region` (if connected).

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

Publishes the device's actual configuration to the host service.

```json
{
  "configuration": { ... }
}
```

### GET /get_configuration

Returns the latest configuration received from the host service.

### PUT /deprovision

Removes the device from the host service and deletes local credentials.

```json
{
  "hostId": "vsf_test_host"
}
```

Use `?force=true` to deprovision while not connected.

## Host Configuration

Host configuration files are JSON files in the `host_configuration/` directory, named `<host_id>.json`. The SDK looks for this directory relative to the binary location, falling back to the working directory.

The included `vsf_test_host.json` points to the VSF test endpoint for development and testing.

### Using the Local TR-12 Host Service

To use the self-hosted [TR-12 Host Service](../host/) instead of the VSF test endpoint, create a host config file:

```bash
cat > host_configuration/local_go_host.json << 'EOF'
{
  "serviceId": "tr12-host",
  "serviceName": "My TR-12 Host",
  "deviceTypes": ["SOURCE", "DESTINATION", "BOTH"],
  "pairingUrl": "http://127.0.0.1:8080",
  "authUrl": "http://127.0.0.1:8080",
  "thumbnailMaxSizeKB": 100,
  "logFileMaxSizeKB": 500
}
EOF
```

Then use `--host_id local_go_host` when starting the ARD. See the [TR-12 Host Service README](../host/README.md) for how to start the host, register an account, and claim devices.

## Credential Storage

Credentials are stored at `<certs_path>/<internal_device_id>/<host_id>/`. The `--internal_device_id` flag controls which subfolder is used:

- Use the same `--internal_device_id` to reconnect with existing credentials
- Use a different `--internal_device_id` to start fresh with no cached credentials
- Each host ID gets its own subfolder, so one device identity can connect to multiple hosts

## Quick Start with VSF Test Endpoint

1. Build the SDK:
   ```bash
   cd client
   go build -o bin/cdd-sdk ./cmd/cdd-sdk/
   ```

2. Start the SDK:
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

3. From your application (or curl), initiate a connection:
   ```bash
   curl -X PUT http://127.0.0.1:8603/connect \
     -H "Content-Type: application/json" \
     -d '{"hostId": "vsf_test_host", "registration": {"deviceName": "My Encoder"}}'
   ```

4. The response will include a `pairingCode`. Claim the device on the VSF test endpoint:
   ```
   PUT <base_endpoint>/authorize/<pairing_code>
   ```

5. Call `/connect` again — the SDK will complete pairing and transition to `CONNECTED`.

6. Once connected, use `/report_status`, `/report_actual_configuration`, and `/get_configuration` to interact with the host service.

## Logging

The SDK writes JSON-formatted rotating log files to the `--log_path` directory. Log files are capped at 500 KB with up to 3 rotated backups. Logs are also printed to stdout.

When the host service requests log uploads, the SDK automatically uploads rotated log files to the provided pre-signed URL.

## Security

The SDK persists X.509 credentials on disk at `<certs_path>/<internal_device_id>/<host_id>/`. Using a different `--internal_device_id` creates a separate credential set with no overlap. While the protocol implements credential rotation to limit certificate lifespan, securing credentials on disk is the host system's responsibility.

See the [Python SDK README](https://github.com/vsf-tv/gccg-cdd) for detailed security best practices.

## Application Reference Design (ARD)

The ARD is a companion program that simulates a 1-channel encoder device making REST calls on the SDK daemon. It mirrors the Python ARD at `gccg-cdd/src/thumbnails/`.

### DeviceCallbacks Interface

The core integration pattern is the `DeviceCallbacks` interface (`internal/application_reference_design/device_callbacks.go`). Any device integration implements this interface — the `ApplicationLoop` and `Tr12Shim` only know the interface, never the concrete implementation.

**Apply (set) side** — called when the host sends a new desired configuration:
- `UpdateDeviceKeyValue` — device-level setting (e.g. clock source)
- `UpdateChannelSettings` — channel simple setting (e.g. framerate, codec, bitrate)
- `UpdateChannelProfile` — profile selection
- `UpdateChannelConnection` — transport protocol config (SRT caller/listener, RIST, etc.)
- `UpdateChannelState` — ACTIVE or IDLE

**Read-back (get) side** — called when building the actual configuration to report back:
- `GetDeviceUpdatedValue`, `GetChannelUpdatedValue`, `GetChannelProfileValue`
- `GetChannelConnection`, `GetChannelState`
- `GetDeviceStatus`, `GetChannelStatus`

**Two implementations ship with this repo:**

`ArdCallbacks` (`internal/application_reference_design/callbacks.go`) — the reference implementation using a simulated ffmpeg encoder. Study this to understand the pattern.

`OspreyCallbacks` (`internal/device/osprey_callbacks.go`, gitignored) — a real device implementation that calls the Osprey encoder's HTTP API. This is the template for any real device integration.

Both are wired into `ApplicationLoop` identically:

```go
// ARD binary — mock encoder:
callbacks := ard.NewArdCallbacks()
loop := ard.NewApplicationLoop(sdkURL, callbacks, &registration)

// proprietary-osprey-encoder-bridge — real device:
callbacks := device.NewOspreyCallbacks(deviceURL, deviceType)
loop := ard.NewApplicationLoop(sdkURL, callbacks, &registration)
```

The `ApplicationLoop` drives the TR-12 lifecycle: connect → get configuration → apply via callbacks → read back actual state → report. The `Tr12Shim` walks TR-12 model structures and dispatches to the callbacks.

### What the ARD Does

- Calls `PUT /connect` in a loop until the SDK reaches `CONNECTED` state
- Displays the pairing code when in `PAIRING` state
- Once connected, polls `GET /get_configuration` for host-service updates
- Applies desired configuration via the TR-12 shim (settings, connection, state)
- Reports device status via `PUT /report_status`
- Reports actual configuration via `PUT /report_actual_configuration`
- Simulates thumbnail emission by cycling through sample images to `/tmp/image_sdi.jpg` and `/tmp/image_hdmi.jpg`
- Manages an ffmpeg subprocess for SRT streaming when the host sets channel state to `ACTIVE`

### Build the ARD

```bash
cd client
go build -o bin/ard ./cmd/application_reference_design/
```

### Run the ARD

First, start the SDK daemon in Terminal 1:
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

Then start the ARD in Terminal 2:
```bash
# Against the VSF test endpoint:
./bin/ard --host_id vsf_test_host

# Or against the local TR-12 Host Service:
./bin/ard --host_id local_go_host
```

The ARD accepts these flags:

| Argument | Description |
|---|---|
| `--host_id` | Host ID to connect to (required) |
| `--sdk_url` | Base URL of the running SDK (default: `http://127.0.0.1:8603`) |

### ARD Workflow

1. The ARD calls `/connect` on the SDK with the registration payload and host ID
2. If no credentials exist, the SDK enters `PAIRING` state and returns a pairing code
3. Claim the device on the host service:
   - VSF test endpoint: `PUT <base_endpoint>/authorize/<pairing_code>`
   - Local host service: `curl -X PUT http://127.0.0.1:8080/authorize/<pairing_code> -H "Authorization: Bearer $TOKEN"`
4. The ARD's next `/connect` call completes pairing → `CONNECTED`
5. The ARD then loops: get configuration → report status → report actual configuration
6. When the host sends a configuration with `"state": "ACTIVE"` and SRT connection details, the ARD starts ffmpeg
7. When the host sends `"state": "IDLE"`, the ARD stops ffmpeg

### ARD Data Files

The ARD expects these files relative to its working directory:

- `payloads/1_channel_encoder/registration.json` — device registration payload
- `cmd/application_reference_design/thumbnails/thumbnail_images_sdi/` — SDI thumbnail source images
- `cmd/application_reference_design/thumbnails/thumbnail_images_hdmi/` — HDMI thumbnail source images

## Project Structure

```
client/
├── cmd/
│   ├── cdd-sdk/main.go              # SDK entry point, CLI flag parsing
│   ├── ard/main.go                  # ARD entry point
│   └── proprietary-osprey-encoder-bridge/                 # Device integration shim (gitignored — device-specific)
│       ├── main.go
│       ├── device_deploy.sh
│       ├── osprey_registration.json
│       └── web/                     # Device web console overrides
├── host_configuration/               # Host config JSON files
│   └── vsf_test_host.json
├── payloads/
│   └── 1_channel_encoder/
│       ├── registration.json         # Device registration payload
│       └── configuration.json        # Example desired configuration
│   └── thumbnails/
│   ├── thumbnail_images_sdi/         # Sample SDI thumbnails
│   └── thumbnail_images_hdmi/        # Sample HDMI thumbnails
├── internal/
│   ├── api/server.go                 # Gin REST API server
│   ├── ard/
│   │   ├── application.go           # ARD main run loop
│   │   ├── application_loop.go      # Reusable ApplicationLoop (ARD + proprietary-osprey-encoder-bridge)
│   │   ├── callbacks.go             # ArdCallbacks — reference DeviceCallbacks implementation
│   │   ├── device_callbacks.go      # DeviceCallbacks interface definition
│   │   ├── encoder.go               # ffmpeg encoder simulation
│   │   ├── sdk_client.go            # SDK REST client
│   │   ├── shim.go                  # TR-12 model shim
│   │   └── thumbnails.go            # Thumbnail image simulator
│   ├── device/                       # Device-specific callbacks (gitignored)
│   │   └── osprey_callbacks.go      # OspreyCallbacks — real device implementation
│   ├── cddlogger/logger.go          # JSON rotating file logger
│   ├── credentials/store.go         # X.509 cert persistence
│   ├── models/
│   │   └── tr12_models.go           # TR-12 protocol model aliases + state constants
│   ├── pairing/pairing.go           # Pairing/auth flow
│   ├── sdk/
│   │   ├── sdk.go                   # Core SDK struct and state machine
│   │   ├── connect.go               # Public API methods
│   │   └── mqtt.go                  # MQTT connection and callbacks
│   ├── thumbnails/manager.go        # Thumbnail upload manager
│   └── utils/utils.go               # TLS, upload, throttle, key gen
├── go.mod
└── go.sum
```

## TR-12 Protocol Reference

- Smithy Models: https://github.com/vsf-tv/TR-12-Models
- Draft Protocol: https://github.com/vsf-tv/TR-12-Models/blob/main/VSF_TR-12-ClientDeviceDiscoveryDraft.pdf

## License

Apache License, Version 2.0
