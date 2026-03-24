# TR-12 Host Service

A self-contained Go binary implementing the server side of the VSF TR-12 Client Device Discovery protocol. Manages device pairing, an embedded MQTT broker for real-time device communication, a REST management API, and persists all state in a single SQLite database file.

This is the cloud-side counterpart to the [Go CDD SDK](../client/). It replaces the VSF test endpoint with a self-hosted, production-grade implementation that runs on any Mac or Linux machine with zero external dependencies.

## Architecture

```
┌──────────────┐   HTTPS (pair/auth)   ┌─────────────────────────────────────┐
│  Device SDK  │ ────────────────────── │         TR-12 Host Service          │
│  (Go/Python) │                        │                                     │
│              │ ◄── MQTT/TLS ────────► │  HTTP API :8080  │  MQTT :8883     │
└──────────────┘                        │       │          │      │          │
                                        │       ▼          │      ▼          │
┌──────────────┐   REST API :8080       │  ┌─────────┐    │  ┌──────────┐  │
│  Console /   │ ────────────────────── │  │ Service  │◄───┘  │ Embedded │  │
│  curl        │                        │  │ Layer    │       │ Broker   │  │
└──────────────┘                        │  └────┬─────┘       └──────────┘  │
                                        │       ▼                            │
                                        │  ┌──────────────────────────┐     │
                                        │  │   SQLite (single file)   │     │
                                        │  └──────────────────────────┘     │
                                        └─────────────────────────────────────┘
```

## Requirements

- Go 1.22 or newer

## Build

```bash
cd host
go build -o tr12-host ./cmd/tr12-host/
```

Produces a single ~18 MB static binary. No CGO, no external dependencies.

## Quick Start

```bash
# Start the service (first run auto-generates CA, server cert, and JWT secret)
./tr12-host --host-address 192.168.1.100

# Or with all options
./tr12-host \
  --host-address 192.168.1.100 \
  --http-port 8080 \
  --mqtt-port 8883 \
  --db-path ./tr12-host.db \
  --service-id my-host \
  --service-name "My TR-12 Host" \
  --cert-expiry-days 30 \
  --rotation-interval-days 30 \
  --pairing-timeout 1800 \
  --jwt-expiry-hours 24 \
  --log-level info
```

On first run the service:
1. Creates the SQLite database
2. Generates a self-signed CA (RSA 4096, 10-year validity)
3. Generates a server TLS certificate signed by the CA
4. Generates a random JWT signing secret
5. Stores everything in the database
6. Starts the MQTT broker and HTTP API

Subsequent runs load all credentials from the database.

## CLI Arguments

| Argument | Required | Default | Description |
|---|---|---|---|
| `--host-address` | Yes | — | Externally reachable IP or hostname |
| `--http-port` | No | `8080` | HTTP API port |
| `--mqtt-port` | No | `8883` | MQTT broker TLS port |
| `--db-path` | No | `./tr12-host.db` | SQLite database file path |
| `--service-id` | No | `tr12-host` | Service identifier |
| `--service-name` | No | `TR-12 Host Service` | Human-readable name |
| `--cert-expiry-days` | No | `30` | Device certificate validity (days) |
| `--rotation-interval-days` | No | `30` | Auto-rotation interval (days) |
| `--pairing-timeout` | No | `1800` | Pairing code timeout (seconds) |
| `--jwt-expiry-hours` | No | `24` | JWT token expiration (hours) |
| `--log-level` | No | `info` | `debug`, `info`, `warn`, `error` |

## API Reference

### Account Endpoints (unauthenticated)

#### Register

```bash
curl -X POST http://localhost:8080/account/register \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "changeme1", "display_name": "Admin User"}'
```

Response:
```json
{
  "account": {
    "account_id": "acc_a1b2c3d4",
    "username": "admin",
    "display_name": "Admin User",
    "created_at": "2025-01-15T12:00:00Z"
  },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

#### Login

```bash
curl -X POST http://localhost:8080/account/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "changeme1"}'
```

#### Get Account

```bash
curl http://localhost:8080/account \
  -H "Authorization: Bearer <token>"
```

### Device-Facing Endpoints (unauthenticated — called by SDK)

#### Pair

Called by the CDD SDK during the pairing flow.

```bash
curl -X POST http://localhost:8080/pair \
  -H "Content-Type: application/json" \
  -d '{
    "deviceType": "SOURCE",
    "hostId": "tr12-host",
    "csr": "-----BEGIN CERTIFICATE REQUEST-----\n...\n-----END CERTIFICATE REQUEST-----",
    "version": "1.0.0"
  }'
```

Response (success):
```json
{
  "result": {
    "success": {
      "deviceId": "001XI02IJ2FtSIirk01",
      "pairingCode": "A3B7K9",
      "accessCode": "a1b2c3d4e5f6...",
      "pairingTimeoutSeconds": 1800
    }
  }
}
```

Response (failure):
```json
{
  "result": {
    "failure": {
      "reason": "HOST_ID_MISMATCH"
    }
  }
}
```

#### Authenticate

Called by the CDD SDK polling for claim status.

```bash
curl -X POST http://localhost:8080/authenticate \
  -H "Content-Type: application/json" \
  -d '{
    "deviceId": "001XI02IJ2FtSIirk01",
    "pairingCode": "A3B7K9",
    "accessCode": "a1b2c3d4e5f6..."
  }'
```

Response (waiting for claim):
```json
{ "status": "STANDBY" }
```

Response (claimed):
```json
{
  "status": "CLAIMED",
  "caCert": "-----BEGIN CERTIFICATE-----\n...",
  "deviceCert": "-----BEGIN CERTIFICATE-----\n...",
  "mqttUri": "tls://192.168.1.100:8883",
  "region": "local",
  "hostSettings": {
    "iotProtocolName": "mqtt",
    "pairingTimeoutSeconds": 1800,
    "minIntervalPubSeconds": 1,
    "mqttKeepaliveSeconds": 30,
    "subUpdateTopic": "cdd/001XI02IJ2FtSIirk01/config/update",
    "pubReportRegistrationTopic": "cdd/001XI02IJ2FtSIirk01/registration/report",
    "pubReportStatusTopic": "cdd/001XI02IJ2FtSIirk01/status/report",
    "..."
  }
}
```

### Management Endpoints (JWT required)

All management endpoints require `Authorization: Bearer <token>` from login/register.

For brevity, examples below use `$TOKEN` as a placeholder:

```bash
export TOKEN="eyJhbGciOiJIUzI1NiIs..."
```

#### Claim a Device

After a device calls `/pair` and receives a pairing code, claim it into your account:

```bash
curl -X PUT http://localhost:8080/authorize/A3B7K9 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"expiration_days": 730}'
```

The `expiration_days` parameter is optional (default: 730 / 2 years). This sets the device's registration lifetime.

#### List Devices

```bash
curl http://localhost:8080/devices \
  -H "Authorization: Bearer $TOKEN"
```

Response:
```json
[
  {
    "device_id": "001XI02IJ2FtSIirk01",
    "message": "",
    "errors": [],
    "online_details": "online: 0d 2h 15m",
    "online": true
  }
]
```

#### Describe Device

```bash
curl http://localhost:8080/device/001XI02IJ2FtSIirk01 \
  -H "Authorization: Bearer $TOKEN"
```

Response:
```json
{
  "device_id": "001XI02IJ2FtSIirk01",
  "message": "",
  "errors": [],
  "registration": { "channels": [...], "thumbnails": [...] },
  "configuration": { "channels": [...] },
  "actual_configuration": { "channels": [...] },
  "status": { "status": [...], "channels": [...] },
  "online": true,
  "online_details": "online: 0d 2h 15m",
  "cert_expiration": "23d 20h 63m",
  "device_metadata": {
    "online": true,
    "online_details": "online: 0d 2h 15m",
    "cert_expiration": "23d 20h 63m",
    "source_ip": "192.168.1.50:54321",
    "device_type": "SOURCE",
    "account_id": "acc_a1b2c3d4",
    "paired_at": "2025-01-15T12:00:00Z"
  }
}
```

#### Update Device Configuration

Push a desired configuration to the device via MQTT:

```bash
curl -X PUT http://localhost:8080/device/001XI02IJ2FtSIirk01 \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "channels": [
      {
        "id": "SDI-1",
        "state": "ACTIVE",
        "settings": {
          "simpleSettings": [
            {"key": "resolution", "value": "1080p"}
          ]
        },
        "connection": {
          "transportProtocol": {
            "srtCaller": {
              "ip": "10.0.0.50",
              "port": 9000,
              "minimumLatencyMilliseconds": 3000
            }
          }
        }
      }
    ]
  }'
```

The service validates the configuration against the device's registration before publishing. Invalid channel IDs, setting keys, or protocol types are rejected with a descriptive 400 error.

#### Deprovision Device

Two-phase deprovision: marks the device as DEPROVISIONED and notifies it via MQTT. Full cleanup happens after the device acknowledges.

```bash
curl -X PUT http://localhost:8080/deprovision/001XI02IJ2FtSIirk01 \
  -H "Authorization: Bearer $TOKEN"
```

#### Get Thumbnail

```bash
curl "http://localhost:8080/thumbnail/001XI02IJ2FtSIirk01?source=SDI-1" \
  -H "Authorization: Bearer $TOKEN"
```

If no active subscription exists, the service creates one (publishes to the device via MQTT) and returns a pending message. On subsequent calls, returns the base64-encoded image.

#### Rotate Credentials

Trigger a certificate rotation for a device:

```bash
curl -X PUT http://localhost:8080/credentials/001XI02IJ2FtSIirk01 \
  -H "Authorization: Bearer $TOKEN"
```

Generates a new device certificate, publishes it as a retained MQTT message, and maintains credential overlap (previous cert stays valid until the device connects with the new one).

### Upload Endpoints (called by devices)

These URLs are generated by the service and sent to devices via MQTT. Devices call them directly.

```
PUT /upload/thumbnail/:deviceId/:sourceId   — thumbnail image upload
PUT /upload/log/:deviceId                   — log file upload
```

## End-to-End Walkthrough

Here's the full flow from starting the service to having a connected device:

```bash
# 1. Start the host service
./tr12-host --host-address 127.0.0.1

# 2. Register an account
curl -s -X POST http://127.0.0.1:8080/account/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"changeme1","display_name":"Admin"}' | jq .
# Save the token from the response
export TOKEN="<token from response>"

# 3. Start the CDD SDK and ARD (see client/README.md for details).
#    The ARD will display a pairing code, e.g.:
#      "Device is not paired. Pairing Code: A3B7K9 Expires in: 1800s."

# 4. Claim the device using the pairing code
curl -X PUT http://127.0.0.1:8080/authorize/A3B7K9 \
  -H "Authorization: Bearer $TOKEN"

# 5. The SDK automatically polls /authenticate. Once claimed it receives
#    certs and connects to the MQTT broker. The ARD will show
#    "State: CONNECTED" within a few seconds.

# 6. List connected devices
curl http://127.0.0.1:8080/devices \
  -H "Authorization: Bearer $TOKEN" | jq .

# 7. Describe a device
curl http://127.0.0.1:8080/device/<DEVICE_ID> \
  -H "Authorization: Bearer $TOKEN" | jq .

# 8. Push a configuration update
curl -X PUT http://127.0.0.1:8080/device/<DEVICE_ID> \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"channels":[{"id":"SDI-1","state":"ACTIVE"}]}'

# 9. Request a thumbnail
curl "http://127.0.0.1:8080/thumbnail/<DEVICE_ID>?source=SDI-1" \
  -H "Authorization: Bearer $TOKEN"

# 10. Rotate device credentials
curl -X PUT http://127.0.0.1:8080/credentials/<DEVICE_ID> \
  -H "Authorization: Bearer $TOKEN"

# 11. Deprovision a device
curl -X PUT http://127.0.0.1:8080/deprovision/<DEVICE_ID> \
  -H "Authorization: Bearer $TOKEN"
```

## Data Storage

Everything lives in a single SQLite file (default `./tr12-host.db`):

- Account credentials (bcrypt hashed passwords)
- Device registry and state (registration, configuration, status)
- CA certificate and private key
- Server TLS certificate and private key
- JWT signing secret
- Device certificates (current + previous for rotation overlap)
- Thumbnails (as BLOBs)
- Device logs (as BLOBs)

To back up or migrate: copy the `.db` file.

## Background Processes

Three background goroutines run automatically:

- **Certificate rotation** — checks hourly, rotates device certs older than `--rotation-interval-days`
- **Pairing cleanup** — checks every 60s, removes expired pairing records
- **Registration expiry** — checks hourly, deprovisions devices past their registration expiration

## Multi-Tenant Isolation

Each account sees only its own devices. When a device is claimed via `/authorize/{pairingCode}`, it's bound to the caller's account. All subsequent queries (list, describe, configure, deprovision, thumbnail, credentials) are scoped by account ID.

## Security

- Passwords stored as bcrypt hashes
- JWT tokens signed with HMAC-SHA256 (secret auto-generated, stored in DB)
- Device certificates signed by the service CA with device ID in subject CN
- MQTT broker enforces mutual TLS — only devices with valid certs can connect
- Per-device topic ACLs prevent cross-device MQTT access
- No secrets on the filesystem — everything in the SQLite DB

## Project Structure

```
host/
├── cmd/tr12-host/main.go              # Entry point, CLI flags, wiring
├── internal/
│   ├── api/
│   │   ├── router.go                 # Gin router, CORS
│   │   ├── middleware.go             # JWT auth middleware
│   │   ├── pairing_handlers.go       # /pair, /authenticate
│   │   ├── account_handlers.go       # /account/*
│   │   ├── device_handlers.go        # /devices, /device/:id, etc.
│   │   └── upload_handlers.go        # /upload/thumbnail, /upload/log
│   ├── broker/
│   │   ├── broker.go                 # Embedded mochi-mqtt broker
│   │   ├── auth_hook.go             # Client cert authentication
│   │   └── acl_hook.go              # Per-device topic ACLs
│   ├── ca/ca.go                      # Certificate Authority
│   ├── config/config.go              # CLI config struct
│   ├── db/
│   │   ├── db.go                     # SQLite connection, migrations
│   │   ├── accounts.go               # Account CRUD
│   │   ├── devices.go                # Device CRUD
│   │   ├── thumbnails.go             # Thumbnail BLOB storage
│   │   ├── logs.go                   # Log BLOB storage
│   │   └── config.go                 # Key-value config store
│   ├── models/
│   │   ├── tr12.go                   # TR-12 protocol models
│   │   ├── device.go                 # Device registry models
│   │   └── account.go                # Account models
│   ├── mqtt/handlers.go              # Internal MQTT subscription handlers
│   └── service/
│       ├── device_service.go         # Device business logic
│       ├── account_service.go        # Account + JWT logic
│       ├── thumbnail_service.go      # Thumbnail subscriptions
│       ├── log_service.go            # Log requests
│       ├── rotation_service.go       # Background cert rotation
│       └── errors.go                 # Sentinel errors
├── go.mod
└── go.sum
```

## Dependencies

| Package | Purpose |
|---|---|
| `github.com/gin-gonic/gin` | HTTP router |
| `github.com/gin-contrib/cors` | CORS middleware |
| `github.com/mochi-mqtt/server/v2` | Embedded MQTT broker |
| `modernc.org/sqlite` | Pure Go SQLite (no CGO) |
| `github.com/golang-jwt/jwt/v5` | JWT tokens |
| `golang.org/x/crypto` | bcrypt password hashing |

## TR-12 Protocol Reference

- Smithy Models: https://github.com/vsf-tv/TR-12-Models
- Draft Protocol: https://github.com/vsf-tv/TR-12-Models/blob/main/VSF_TR-12-ClientDeviceDiscoveryDraft.pdf

## License

Apache License, Version 2.0
