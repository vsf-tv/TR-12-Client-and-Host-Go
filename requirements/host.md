# TR-12 Host Service: Requirements

## Overview

Build a TR-12 Host Service — the cloud-side counterpart to the existing Go CDD SDK and ARD. The host service is a self-contained Go binary that manages a device registry, handles pairing/authentication, runs an embedded MQTT broker for device communication, stores device state (registration, configuration, status, thumbnails), and exposes a management REST API. It must run on a Mac, Linux instance, or any machine with no external service dependencies (no AWS, no cloud infrastructure).

The existing Go SDK (`client/`) and ARD are complete and tested against the VSF test endpoint. This host service replaces that test endpoint with a production-grade, self-hosted implementation.

The management API must be designed to support both programmatic access and a future model-driven web console that dynamically generates UI from TR-12 Smithy models.

### Reference Materials
- TR-12 Smithy Models: `models/TR-12-Models/smithy/common.smithy`
- TR-12 Configuration Models: `models/TR-12-Models/smithy/configuration.smithy`
- TR-12 Registration Models: `models/TR-12-Models/smithy/registration.smithy`
- TR-12 Service API: `models/TR-12-Models/smithy/service_api.smithy`
- TR-12 Status Models: `models/TR-12-Models/smithy/status.smithy`
- Go SDK README: `client/README.md`
- Go SDK Client Requirements: `requirements/client.md`
- VSF Test Host Config: `client/host_configuration/vsf_test_host.json`

---

## Requirement 0: Use Smithy-Generated TR-12 Models

The host service must use the official Smithy-generated Go models for all TR-12 protocol types. These models are generated from the TR-12 Smithy definitions (`models/TR-12-Models/smithy/`) using the OpenAPI generator (`generate-tr12-models.sh`), producing Go structs in the `tr12go` package. The Go SDK already consumes these generated models via the shared submodule.

The host service must NOT define hand-rolled Go structs that duplicate the generated protocol types. Instead, it must import the generated models from the shared submodule and use them directly in API handlers, service logic, MQTT message serialization/deserialization, and database I/O.

### Model Source
- Smithy definitions: `models/TR-12-Models/smithy/*.smithy`
- Generator: `models/TR-12-Models/generate-tr12-models.sh`
- Generated Go output: `models/TR-12-Models/generated/tr12go/` (package `openapi`)

### Generated Protocol Types (used by the host service)
| Generated Type | Smithy Source | Usage |
|---|---|---|
| `PairRequestContent` | `PairRequest` | `POST /pair` request body |
| `PairResponseContent` | `PairResponse` | `POST /pair` response body |
| `PairResult` | `PairResult` (union) | Success/failure wrapper in pair response |
| `PairSuccessData` | `PairSuccessData` | Pairing success payload |
| `PairFailureData` | `PairFailureData` | Pairing failure payload |
| `PairFailureReason` | `PairFailureReason` (enum) | HOST_ID_MISMATCH, VERSION_NOT_SUPPORTED, etc. |
| `AuthenticateRequestContent` | `AuthRequest` | `POST /authenticate` request body |
| `AuthenticateResponseContent` | `AuthResponse` | `POST /authenticate` response body |
| `AuthStatus` | `AuthStatus` (enum) | STANDBY, CLAIMED |
| `HostSettings` | `HostSettings` | MQTT topics, keepalive, throttle intervals |
| `RotateCertificatesRequestContent` | `CertRotate` | Certificate rotation MQTT payload |
| `DeprovisionDeviceRequestContent` | `DeprovisionMessage` | Deprovision MQTT payload |
| `DeprovisionReason` | `DeprovisionReason` (enum) | DEPROVISIONED, EXPIRED, UNKNOWN |
| `RequestThumbnailRequestContent` | `ThumbnailSubscription` | Thumbnail subscription MQTT payload |
| `ThumbnailRequest` | `ThumbnailRequest` | Per-source thumbnail parameters |
| `RequestLogRequestContent` | `LogRequest` | Log upload request MQTT payload |
| `GetHostConfigResponseContent` | `HostConfig` | Host configuration response |

### Integration Pattern
The host service imports the generated models from the shared submodule (`models/TR-12-Models/generated/tr12go/`) and re-exports them via type aliases in `host/internal/models/tr12.go`. This allows service code to reference `models.PairRequestContent`, `models.AuthStatus`, etc. while the underlying types are the canonical Smithy-generated structs. Host-service-only types (e.g., `Device`, `Account`, `ClaimRequest`) that have no Smithy equivalent are defined directly in the `host/internal/models/` package.

### Acceptance Criteria
- [ ] All TR-12 protocol types used by the host service are the Smithy-generated Go structs, not hand-rolled duplicates
- [ ] Generated models are imported from the shared submodule at `models/TR-12-Models/generated/tr12go/`
- [ ] `host/internal/models/tr12.go` re-exports generated types as type aliases so service code uses `models.XYZ` consistently
- [ ] Host-service-only types (Device, Account, etc.) are defined separately in `host/internal/models/` and do not duplicate any generated type
- [ ] JSON serialization of protocol messages matches the SDK's serialization exactly (both use the same generated structs)
- [ ] Enum constants (e.g., `AuthStatus` values, `PairFailureReason` values, `DeprovisionReason` values) are used from the generated package, not redefined
- [ ] When the upstream Smithy models are updated and regenerated, the host service can update by pulling the submodule — no hand-rolled struct changes needed
- [ ] The generated models' constructor functions (e.g., `NewHostSettings()`) and accessor methods (e.g., `GetMqttUri()`) are used where appropriate instead of direct struct literal construction

---

## Requirement 1: Device Pairing and Authentication

The host service must implement the TR-12 pairing protocol (server side).

### Pairing Flow (Host Side)
1. Device SDK calls `POST /pair` with device type, host ID, CSR, and protocol version
2. Host validates the request (device type supported, version compatible, host ID matches)
3. Host generates a unique device ID, pairing code, and access code
4. Host signs the CSR to produce a device certificate using the service's local CA
5. Host stores the pending pairing record with a configurable timeout (`pairingTimeoutSeconds`)
6. Host returns `PairResult` with success data (device ID, pairing code, access code, timeout)

### Authentication Flow (Host Side)
1. Device SDK polls `POST /authenticate` with device ID, pairing code, and access code
2. If device is not yet claimed: return `AuthResponse` with `status: STANDBY`
3. If device is claimed: return `AuthResponse` with `status: CLAIMED`, CA cert, device cert, MQTT URI (pointing to the embedded broker), region, and `HostSettings`
4. `HostSettings` includes all MQTT topic templates (scoped per device ID), keepalive, throttle intervals

### Claim Flow (Management API)
1. Authenticated user calls `PUT /authorize/{pairingCode}` with a valid JWT
2. Request body accepts an optional `expiration_days` parameter (default: 730 days / 2 years) that sets the device's overall registration expiration — this is the maximum lifetime of the device's registration, not the individual certificate expiration
3. Host validates the pairing code exists and hasn't expired
4. Host marks the device as claimed, associating it with the caller's `account_id` from the JWT, and records the registration expiration date

### Acceptance Criteria
- [ ] `POST /pair` validates device type against the service's supported types
- [ ] `POST /pair` returns `PairFailureReason` (HOST_ID_MISMATCH, VERSION_NOT_SUPPORTED, DEVICE_TYPE_NOT_SUPPORTED) for invalid requests
- [ ] Pairing codes expire after `pairingTimeoutSeconds` (configurable, default 1800s)
- [ ] `POST /authenticate` returns STANDBY until the device is claimed via the management API
- [ ] `PUT /authorize/{pairingCode}` accepts an optional `expiration_days` body parameter (default 730) that sets the device registration expiration
- [ ] Device certificates are signed by the service's local CA with a configurable expiration
- [ ] MQTT topics in `HostSettings` are scoped per device (e.g., `cdd/{deviceId}/config/update`)
- [ ] Pairing records are cleaned up after expiration (background goroutine)
- [ ] Devices whose registration expiration has passed are automatically deprovisioned (background goroutine)

---

## Requirement 2: Device Registry and State Storage

The host service must maintain a persistent registry of all paired devices and their current state using an embedded database.

### Per-Device State
- **Registration** — Full `DeviceRegistration` payload (channels, settings, profiles, thumbnails) as received from the device via MQTT
- **Desired Configuration** — The latest `DeviceConfiguration` pushed to the device by the host
- **Actual Configuration** — The latest `DeviceConfiguration` reported back by the device
- **Status** — The latest `DeviceStatus` reported by the device (device-level and per-channel)
- **Online State** — Whether the device's MQTT connection is active, with last-seen timestamp
- **Certificate Metadata** — Cert expiration, last rotation timestamp
- **Device Metadata** — Device ID, device type, owner account_id, pairing timestamp, source IP

### Acceptance Criteria
- [ ] All device state is persisted in an embedded database (SQLite via `modernc.org/sqlite` — pure Go, no CGO)
- [ ] Registration, configuration, and status payloads are stored as full structured JSON (no lossy transformations)
- [ ] Online state is updated via MQTT broker connect/disconnect events
- [ ] Device state is queryable by device ID
- [ ] Database file path is configurable via CLI flag `--db-path` (default: `./tr12-host.db`)
- [ ] Database is created automatically on first run with schema migrations
- [ ] The SQLite database is the ONLY filesystem artifact — all persistence (device state, CA certs, thumbnails, logs, JWT secret, server cert) is stored in the database

---

## Requirement 3: Embedded MQTT Broker

The host service must run an embedded MQTT broker for device communication, eliminating the need for an external MQTT service.

### Broker Requirements
1. Embedded Go MQTT broker (e.g., `mochi-mqtt/server`) running within the host service process
2. TLS listener using the service's CA cert and server cert for mutual TLS authentication
3. Device connections authenticated via client certificate (the device cert issued during pairing)
4. Per-device topic ACLs enforced by the broker — each device can only publish/subscribe to its own topics

### MQTT Topic Structure
Per-device topics (templates from `HostSettings`):
- `cdd/{deviceId}/config/update` — Host → device: desired configuration
- `cdd/{deviceId}/certs/update` — Host → device: certificate rotation
- `cdd/{deviceId}/thumbnail/subscription` — Host → device: thumbnail requests
- `cdd/{deviceId}/deprovision` — Bidirectional: deprovision commands
- `cdd/{deviceId}/log/subscription` — Host → device: log upload requests
- `cdd/{deviceId}/registration/report` — Device → host: registration payload
- `cdd/{deviceId}/status/report` — Device → host: status updates
- `cdd/{deviceId}/config/actual/report` — Device → host: actual configuration

### Internal Message Handling
The host service subscribes to all device report topics (`cdd/+/registration/report`, `cdd/+/status/report`, `cdd/+/config/actual/report`, `cdd/+/deprovision`) using wildcard subscriptions. Incoming messages are processed by internal handlers that update the database.

### Acceptance Criteria
- [ ] MQTT broker runs in-process (no external MQTT server required)
- [ ] TLS with mutual authentication using the service CA
- [ ] Device connections are authenticated by validating the client certificate against the service CA
- [ ] Per-device topic ACLs prevent cross-device access
- [ ] Connect/disconnect events update device online state in the database
- [ ] MQTT port is configurable (default: 8883 for TLS)
- [ ] Wildcard subscriptions route device-published messages to internal handlers
- [ ] Broker supports configurable keepalive and handles client reconnection gracefully

---

## Requirement 4: Management REST API

The host service must expose a REST API for managing devices. The API paths must match the VSF test host endpoint, ensuring existing tooling and documentation remain compatible. This API is consumed by operators, the future web console, and programmatic clients.

All management API endpoints (except `/pair`, `/authenticate`, `/host-config`, `/version`) require authentication via a Bearer token (JWT) obtained from the account system (see Requirement 11). The authenticated user's account ID scopes all device operations.

### Device Management Endpoints (VSF-compatible)

| Operation | Method | Path | Auth | Description |
|-----------|--------|------|------|-------------|
| ListDevices | GET | `/devices` | Yes | List all registered devices for the caller's account |
| DescribeDevice | GET | `/device/{deviceId}` | Yes | Full device state: registration, config, status, online details, cert expiration |
| UpdateConfiguration | PUT | `/device/{deviceId}` | Yes | Push desired configuration to device via MQTT |
| Authorize/Claim | PUT | `/authorize/{pairingCode}` | Yes | Claim a pairing device into the caller's account |
| Deprovision | PUT | `/deprovision/{deviceId}` | Yes | Set device to DEPROVISIONED state, notify device via MQTT (Phase 1 of two-phase deprovision) |
| GetThumbnail | GET | `/thumbnail/{deviceId}?source={sourceId}` | Yes | Retrieve latest thumbnail for a source |
| RotateCredentials | PUT | `/credentials/{deviceId}` | Yes | Trigger certificate rotation |

### Device-Facing Endpoints (unauthenticated — called by SDK)

| Operation | Method | Path | Auth | Description |
|-----------|--------|------|------|-------------|
| Pair | POST | `/pair` | No | Device pairing request |
| Authenticate | POST | `/authenticate` | No | Device authentication polling |

### DescribeDevice Response Shape
Matches the VSF test host response format:
```json
{
  "device_id": "string",
  "message": "string",
  "errors": [],
  "registration": { "DeviceRegistration (full Smithy structure)" },
  "configuration": { "DeviceConfiguration (desired)" },
  "actual_configuration": { "DeviceConfiguration (actual, from device)" },
  "status": { "DeviceStatus (full Smithy structure)" },
  "online": true,
  "online_details": "online: 0d 2h 15m",
  "cert_expiration": "23d 20h 63m",
  "device_metadata": {
    "online": true,
    "online_details": "online: 0d 2h 15m",
    "cert_expiration": "23d 20h 63m",
    "source_ip": "203.0.113.42",
    "device_type": "SOURCE",
    "account_id": "acc_a1b2c3d4",
    "paired_at": "2025-01-15T12:00:00Z"
  }
}
```

### ListDevices Response Shape
```json
[
  {
    "device_id": "001XI02IJ2FtSIirk01",
    "message": "",
    "errors": [],
    "online_details": "online: 0d 0h 11m",
    "online": true
  }
]
```

### UpdateConfiguration Response Shape
```json
{
  "device_id": "string",
  "message": "Device updated",
  "error": ""
}
```

### Acceptance Criteria
- [ ] Management API runs on a configurable HTTP port (default: 8080)
- [ ] Management endpoints require a valid Bearer token (JWT) in the `Authorization` header
- [ ] Device-facing endpoints (`/pair`, `/authenticate`) are unauthenticated
- [ ] All device operations are scoped to the authenticated user's account ID
- [ ] API paths match the VSF test host endpoint (`/devices`, `/device/{id}`, `/authorize/{code}`, `/deprovision/{id}`, `/thumbnail/{id}`, `/credentials/{id}`)
- [ ] `ListDevices` returns only devices belonging to the caller's account
- [ ] `DescribeDevice` returns the complete registration, desired configuration, actual configuration, and status as structured JSON matching Smithy models
- [ ] `UpdateConfiguration` validates the payload against the device's registration before publishing to MQTT
- [ ] `Deprovision` sets device state to DEPROVISIONED (Phase 1), full cleanup occurs after device acknowledgment (Phase 2)
- [ ] `RotateCredentials` generates a new cert, publishes via MQTT as a retained message, and updates cert metadata
- [ ] `GetThumbnail` returns the latest stored thumbnail as a base64-encoded image with metadata (timestamp, type, size)
- [ ] API returns appropriate HTTP status codes (200, 400, 401, 404, 500) with consistent error response format
- [ ] CORS is enabled for all origins (to support the future web console)

---

## Requirement 5: Configuration Validation

The host service must validate desired configuration payloads against the device's registration before pushing them to the device.

### Validation Rules
1. Each `ChannelConfiguration.id` in the desired config must match a `Channel.id` in the device's registration
2. If `SettingsChoice` is `simpleSettings`, each `IdAndValue.key` must match a `Setting.id` in the channel's registration
3. If `SettingsChoice` is `profile`, the `SettingProfile.id` must match a `ProfileDefinition.id` in the channel's registration
4. `ChannelState` must be a valid enum value (`ACTIVE` or `IDLE`)
5. If a `Connection.transportProtocol` is specified, the protocol variant (e.g., `srtCaller`, `ristListener`) must correspond to a `SupportedProtocol` in the channel's `connectionProtocols` list
6. Transport protocol fields must satisfy Smithy constraints (e.g., port ranges, Hex32/Hex64 patterns for encryption passcodes)

### Acceptance Criteria
- [ ] `UpdateConfiguration` (PUT `/device/{deviceId}`) rejects payloads with unknown channel IDs (HTTP 400)
- [ ] `UpdateConfiguration` rejects payloads with unknown setting IDs or profile IDs (HTTP 400)
- [ ] `UpdateConfiguration` rejects payloads with unsupported transport protocol types for the channel (HTTP 400)
- [ ] Validation error responses include a descriptive message identifying the invalid field and the valid options from the registration
- [ ] Valid payloads are published to the device's MQTT config update topic and stored as the desired configuration in the database

---

## Requirement 6: Certificate Management

The host service must operate as a local Certificate Authority (CA) for device certificates and support manual and automatic rotation. All CA materials are stored in the SQLite database — no filesystem paths for certificates.

### CA Operations
1. On first run, generate a self-signed root CA certificate (RSA 4096-bit, 10-year validity) and private key, and store them in the SQLite database
2. On subsequent runs, load the CA cert and key from the database into memory
3. During pairing, sign the device's CSR to produce a device certificate with a configurable expiration (default 30 days)
4. The CA cert is returned to the device in the `AuthResponse` and used for mutual TLS on the MQTT connection
5. Go's `tls.Certificate` is constructed from PEM bytes loaded from the database into memory — no cert files on disk

### Manual Rotation
1. `PUT /credentials/{deviceId}` generates a new device certificate from the device's stored CSR
2. `SignCSR` uses a cryptographically random 128-bit serial number (`crypto/rand.Int`) for each signing — guarantees unique certs even for the same CSR
3. The new cert, MQTT URI, and region are published to the device's cert update MQTT topic as a retained message using the `RotateCertificatesRequestContent` struct
4. The MQTT URI uses the `tls://` scheme (e.g., `tls://127.0.0.1:8883`), matching the authentication response format
5. Database is updated: current cert → previous cert, new cert → current cert

### Auto-Rotation
1. A background goroutine runs every 30 days (configurable via `--rotation-interval-days`) and rotates certificates for all active devices
2. Rotation messages are published as retained MQTT messages so offline devices receive them on reconnect
3. Credential overlap: both current and previous cert sets are valid simultaneously until the device connects with the new cert
4. When a device connects using the current cert, the previous cert is revoked
5. When a device connects using the previous cert, the connection is allowed — the device will pick up the retained rotation message and reconnect with the new cert
6. The CSR stays the same between rotations — the cert differs due to random serial and timestamps

### Acceptance Criteria
- [ ] CA certificate, CA private key, server certificate, server private key, and JWT signing secret are all stored in the SQLite database
- [ ] On startup, CA and server cert materials are loaded from the database into memory for the MQTT broker's TLS listener
- [ ] CA is auto-generated on first run (RSA 4096-bit, 10-year validity) if not present in the database
- [ ] Device certificates are signed with the service CA and include the device ID in the subject
- [ ] Certificate expiration is configurable (default 30 days)
- [ ] `SignCSR` uses a cryptographically random 128-bit serial number so each call produces a unique certificate
- [ ] `PUT /credentials/{deviceId}` generates a new cert and publishes a `RotateCertificatesRequestContent` message via MQTT
- [ ] The rotation message `mqttUri` uses the `tls://` scheme matching the authentication response format
- [ ] The database is updated (current → previous, new → current) before the MQTT publish
- [ ] Auto-rotation runs every 30 days (configurable) as a background goroutine
- [ ] Rotation messages are published as retained MQTT messages
- [ ] Credential overlap: after rotation, both current and previous cert sets are valid simultaneously
- [ ] The database tracks `current_cert`, `previous_cert`, and their expiration timestamps per device
- [ ] No filesystem paths are used for any certificate or key material

---

## Requirement 7: Thumbnail Storage and Retrieval

The host service must store device thumbnails in the SQLite database and serve them to API consumers. No filesystem paths are used for thumbnail storage.

### Upload Flow (Device → Host)
1. The host generates an HTTP upload URL pointing to itself (e.g., `http://<host>:<httpPort>/upload/thumbnail/{deviceId}/{sourceId}`)
2. The upload URL uses `http://` (not `https://`) since the host service's HTTP listener does not use TLS
3. The upload URL, along with `localPath`, `period`, `expires`, and `maxSizeKilobyte`, is published to the device via MQTT as a `RequestThumbnailRequestContent` message
4. The `localPath` for each thumbnail source is resolved from the device's registration payload — the host reads the stored `DeviceRegistration` JSON, finds the `thumbnails` array entry whose `id` matches the requested `sourceId`, and includes that `localPath` in the subscription message
5. The `maxSizeKilobyte` must accommodate typical device thumbnails (ARD sample images are ~131 KB); default 500 KB
6. The device SDK reads the image from `localPath`, validates freshness (< 10 seconds old) and size, then uploads via HTTP PUT
7. The host stores the image as a BLOB in the SQLite database

### Storage Schema
- `device_id` TEXT, `source_id` TEXT — composite primary key
- `image_data` BLOB — raw image bytes
- `timestamp` TEXT — ISO 8601 upload timestamp
- `image_type` TEXT — image format (e.g., `jpg`, `png`)
- `image_size_kb` INTEGER — image size in kilobytes
- Only the latest thumbnail per source is kept (UPSERT)

### Retrieval Flow (API Consumer → Host)
1. `GET /thumbnail/{deviceId}?source={sourceId}` checks for an active subscription
2. If no subscription exists, creates one (publishes to MQTT) and returns a "pending" response
3. If a subscription exists and a thumbnail is available, returns the image as base64 with metadata

### Acceptance Criteria
- [ ] Thumbnails are stored as BLOBs in the SQLite `thumbnails` table — no filesystem paths
- [ ] Upload endpoint accepts HTTP PUT with the image body and stores it in the database
- [ ] Upload URL uses `http://` scheme (not `https://`)
- [ ] The `localPath` in the thumbnail subscription is resolved from the device's registration `thumbnails` array
- [ ] `maxSizeKilobyte` is set to a value that accommodates typical device images (default 500 KB)
- [ ] `GET /thumbnail` returns base64-encoded image with `timestamp`, `image_type`, `max_size_KB`, and `image_size_KB` metadata
- [ ] Old thumbnails are overwritten via UPSERT (only latest is kept per source)

---

## Requirement 8: Log Collection

The host service must support requesting and storing device log uploads in the SQLite database. Only the most recent log upload per device is retained.

### Flow
1. An operator triggers a log request for a device
2. The host generates an HTTP upload URL and publishes a `LogRequest` message to the device's log subscription MQTT topic
3. The device uploads rotated log files to the upload URL
4. The uploaded log is stored in the database, overwriting any previous log for that device

### Storage Schema
- `device_id` TEXT PRIMARY KEY
- `log_data` BLOB — raw log file content
- `uploaded_at` TEXT — ISO 8601 upload timestamp
- `log_size_kb` INTEGER — log size in kilobytes

### Acceptance Criteria
- [ ] Log upload URLs point to the host service's own HTTP upload endpoint
- [ ] Only the last log upload per device is kept — stored as a BLOB, overwriting via UPSERT
- [ ] The `logFileMaxSizeKB` from `HostConfig` is communicated to the device and enforced on upload
- [ ] Log requests are published to the device's log subscription MQTT topic
- [ ] No filesystem paths are used for log storage

---

## Requirement 9: Deprovision Workflow

The host service must support both API-initiated and device-initiated deprovisioning. API-initiated deprovision is a two-phase process.

### API-Initiated Deprovision (Two-Phase)

**Phase 1: Mark as DEPROVISIONED**
1. `PUT /deprovision/{deviceId}` sets the device state to `DEPROVISIONED` in the database
2. The host publishes a `DeprovisionMessage` (reason: `DEPROVISIONED`) to the device's deprovision MQTT topic
3. The device record is NOT deleted — it remains in the database with `state = DEPROVISIONED`
4. Status and configuration fields are cleared from the device's record
5. The update API rejects configuration changes for this device
6. The device still appears in `ListDevices` with `state = DEPROVISIONED`

**Phase 2: Full Cleanup (after device acknowledgment)**
1. The device receives the deprovision message and publishes an acknowledgment back on the deprovision topic
2. Only after the host receives this acknowledgment does it perform full cleanup:
   - Revoke the device certificate
   - Clean up MQTT resources
   - Remove the device record from the database
   - Delete thumbnail and log rows for this device

### Device-Initiated Deprovision
1. The device publishes a `DeprovisionMessage` to its deprovision topic first
2. The host performs immediate full cleanup since the device initiated the action

### Acceptance Criteria
- [ ] `PUT /deprovision/{deviceId}` sets device state to `DEPROVISIONED` and publishes via MQTT (Phase 1)
- [ ] Device record is NOT deleted in Phase 1 — remains visible in `ListDevices`
- [ ] Status and configuration fields are cleared on Phase 1
- [ ] The update API rejects changes for devices in `DEPROVISIONED` state (HTTP 409)
- [ ] Full cleanup only happens after device acknowledgment (Phase 2)
- [ ] Device-initiated deprovision triggers immediate full cleanup
- [ ] Deprovisioning an already-deprovisioned device is idempotent
- [ ] Deprovisioning an unknown device returns HTTP 404

---

## Requirement 10: Self-Contained Binary and Configuration

The host service must be a single Go binary with no external dependencies.

### CLI Arguments
| Argument | Required | Default | Description |
|----------|----------|---------|-------------|
| `--http-port` | No | `8080` | Management API HTTP port |
| `--mqtt-port` | No | `8883` | MQTT broker TLS port |
| `--db-path` | No | `./tr12-host.db` | Path to the SQLite database file (the ONLY filesystem artifact) |
| `--service-id` | No | `tr12-host` | Service identifier returned in `HostConfig` |
| `--service-name` | No | `TR-12 Host Service` | Human-readable service name |
| `--host-address` | Yes | — | Externally reachable address for constructing MQTT URI and upload URLs |
| `--cert-expiry-days` | No | `30` | Device certificate validity period |
| `--rotation-interval-days` | No | `30` | Auto-rotation interval in days |
| `--pairing-timeout` | No | `1800` | Pairing code timeout in seconds |
| `--jwt-expiry-hours` | No | `24` | JWT token expiration in hours |
| `--log-level` | No | `info` | Log verbosity: `debug`, `info`, `warn`, `error` |

### Filesystem Footprint
The only filesystem artifact is the single SQLite database file specified by `--db-path`. Everything is stored inside it:
- Account data (usernames, password hashes)
- Device registry and state
- CA certificate, CA private key, server certificate, server private key
- JWT signing secret
- Device certificates (current and previous per device)
- Thumbnails (as BLOBs)
- Device logs (as BLOBs)

To move or back up the service, copy the single `.db` file.

### Acceptance Criteria
- [ ] `go build -o tr12-host ./cmd/tr12-host/` produces a single static binary (no CGO)
- [ ] Binary runs on macOS and Linux without external dependencies
- [ ] All state is stored in the single SQLite database file at `--db-path`
- [ ] The SQLite DB file is the ONLY filesystem artifact
- [ ] CA, server cert, and JWT secret are auto-generated on first run and stored in the database
- [ ] Graceful shutdown on SIGINT/SIGTERM
- [ ] Startup logs print the HTTP API URL, MQTT broker URL, and database file path
- [ ] Go 1.24+ required, dependencies managed via `go.mod`

---

## Requirement 11: Account System and Multi-Tenant Isolation

The host service must provide a built-in account system for user registration, login, and device ownership scoping.

### Account Model
- Each account has: `account_id` (auto-generated, e.g., `acc_a1b2c3d4`), `username` (unique), `password_hash`, `display_name`, `created_at`
- Passwords are hashed with bcrypt

### Account API Endpoints

| Operation | Method | Path | Auth | Description |
|-----------|--------|------|------|-------------|
| Register | POST | `/account/register` | No | Create a new account |
| Login | POST | `/account/login` | No | Authenticate and receive a JWT token |
| GetAccount | GET | `/account` | Yes | Return current account info |

### Device Ownership
- `PUT /authorize/{pairingCode}` associates the device with the caller's `account_id`
- All device queries are scoped by `account_id` — users cannot access other accounts' devices
- The `account_id` is stored on the device record and included in `device_metadata`

### Acceptance Criteria
- [ ] `POST /account/register` creates a new account with a unique username and returns a JWT
- [ ] `POST /account/login` validates credentials and returns a JWT
- [ ] Passwords are stored as bcrypt hashes
- [ ] JWT tokens include `account_id` and `username` with configurable expiration
- [ ] JWT signing secret is auto-generated on first run and persisted in the SQLite database
- [ ] All management API endpoints require a valid JWT in the `Authorization: Bearer <token>` header
- [ ] Device-facing endpoints (`/pair`, `/authenticate`) and account endpoints are unauthenticated
- [ ] All device queries are scoped by `account_id`
- [ ] Duplicate username registration returns HTTP 409 Conflict
- [ ] Invalid or expired JWT returns HTTP 401 Unauthorized

---

## Requirement 12: Model-Driven Console Support

The host service API must preserve full Smithy model fidelity to support a future model-driven web console (separate TypeScript/React project).

### Console Contract
The console will dynamically generate UI pages from the TR-12 Smithy models. The host service API must return data in a structure that enables this without transformation:

1. `DeviceRegistration` serves as the UI schema:
   - `Channel` entries → tabbed or sectioned device management pages
   - `Setting` with `EnumValues` → dropdown/select widgets (values list + defaultValue)
   - `Setting` with `RangeValues` → slider/numeric input widgets (min, max, defaultValue)
   - `Setting` with neither → free-text input
   - `ProfileDefinition` entries → profile selector (name, id, info)
   - `connectionProtocols` list → connection configuration form variants
   - `Thumbnail` entries → thumbnail viewer panels (name, id)

2. `DeviceConfiguration` (desired) vs actual configuration → side-by-side comparison view
3. `DeviceStatus` with `StatusValue` entries → read-only status panels
4. `DescribeDevice` must return all of the above in a single response

### Acceptance Criteria
- [ ] `DescribeDevice` returns `registration`, `configuration` (desired), `actual_configuration`, and `status` as top-level structured JSON fields
- [ ] Registration `Setting` structures include `enums` and `ranges` objects with all fields
- [ ] Registration `ProfileDefinition` entries include `name`, `id`, and `info`
- [ ] Registration `Channel.connectionProtocols` lists are preserved
- [ ] Registration `Thumbnail` entries include `name`, `id`, and `localPath`
- [ ] No lossy transformations are applied to any Smithy-modeled payload between device report and API response
- [ ] The `updateId` on configuration is preserved end-to-end for stale-state detection


---

## Unit Test Coverage

Unit tests are placed alongside source files (idiomatic Go convention) and run with standard `go test`:

```bash
go test ./host/... -count=1
```

### Covered Packages

| Package | Test File | Tests | What's Covered |
|---------|-----------|-------|----------------|
| `internal/db` | `db_test.go` | 23 | Full CRUD for devices, accounts, thumbnails, logs, config; state updates, claim, expiry queries, rotation queries, upsert behavior |
| `internal/service` | `account_service_test.go` | 12 | Register success/duplicate/validation, login success/wrong password/nonexistent, JWT validate/expired/tampered/wrong secret, account ID format |
| `internal/service` | `device_service_test.go` | 22 | Pair success/host mismatch/bad type/empty version, authenticate standby/claimed/wrong creds/not found, claim success/not found/already claimed, list devices, describe success/not found/forbidden, update config success/validation error, deprovision/forbidden/idempotent, full cleanup, rotate credentials/not active, validateConfiguration 5 cases, helper format functions |

All tests use the Go standard library only — no third-party test frameworks. `device_service_test.go` uses a `mockMQTT` struct implementing the `MQTTPublisher` interface to isolate service logic from the MQTT broker.
