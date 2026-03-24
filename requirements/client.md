# CDD Go SDK — TR-12 Client Device Discovery: Requirements

## Overview

The TR-12 Client Device Discovery (CDD) SDK and Application Reference Design (ARD) are implemented in Go. The SDK exposes a localhost REST API and CLI interface for device applications, while communicating with any TR-12 compliant host service via HTTPS (pairing/auth) and MQTT (pub/sub).

### Reference Materials
- TR-12 Smithy Models: `models/TR-12-Models/smithy/common.smithy`
- TR-12 Configuration Models: `models/TR-12-Models/smithy/configuration.smithy`
- TR-12 Registration Models: `models/TR-12-Models/smithy/registration.smithy`
- TR-12 Service API: `models/TR-12-Models/smithy/service_api.smithy`
- TR-12 Status Models: `models/TR-12-Models/smithy/status.smithy`

---

## Requirement 1: SDK State Machine

The SDK must implement a connection state machine with the following states and transitions:

- **DISCONNECTED** → Initial state. Transitions to PAIRING (no certs) or CONNECTING (certs exist).
- **PAIRING** → Awaiting device claim. Transitions to CONNECTING (claimed) or DISCONNECTED (expired/error).
- **CONNECTING** → Establishing MQTT connection. Transitions to CONNECTED (success) or DISCONNECTED (failure).
- **CONNECTED** → Fully operational. Transitions to RECONNECTING (connection lost) or DISCONNECTED (explicit disconnect/deprovision).
- **RECONNECTING** → Auto-reconnect in progress. Transitions to CONNECTED (success) or DISCONNECTED (failure).

### Acceptance Criteria
- [ ] State transitions follow the defined state machine exactly
- [ ] Thread-safe state access via mutex
- [ ] `Connect()` is idempotent — repeated calls in CONNECTED state return current status without side effects
- [ ] State resets cleanly on `Disconnect()` or `Deprovision()`

---

## Requirement 2: Pairing and Authentication Flow

The SDK must implement the TR-12 pairing protocol as defined in the Smithy models.

### Pairing Flow
1. SDK calls `POST /pair` on the host service's pairing URL with device type, host ID, CSR, and protocol version
2. Host returns a `PairResult` (success with pairing code + access code, or failure with reason)
3. SDK exposes the pairing code to the calling application
4. Application user claims the device on the host service using the pairing code
5. SDK polls `POST /authenticate` with device ID, pairing code, and access code
6. On success, host returns CA cert, device cert, MQTT URI, region, and host settings
7. SDK persists credentials to filesystem

### Acceptance Criteria
- [ ] CSR generation uses RSA 2048-bit keys
- [ ] Pairing code has a configurable timeout (from host service `pairingTimeoutSeconds`)
- [ ] Expired pairing codes trigger a reset and require a new `/connect` call
- [ ] Credentials (CA cert, device cert, private key, connection settings, host settings) are persisted to `<certs_path>/<device_id>/<host_id>/`
- [ ] Pairing failure reasons (HOST_ID_MISMATCH, VERSION_NOT_SUPPORTED, DEVICE_TYPE_NOT_SUPPORTED) are surfaced in the API response

---

## Requirement 3: MQTT Communication

The SDK must establish and maintain a TLS-secured MQTT connection to the host service.

### Subscriptions (host → device)
- **Configuration updates** (`subUpdateTopic`) — Desired device configuration from host
- **Certificate rotation** (`subUpdateCertsTopic`) — New device certificates
- **Thumbnail subscriptions** (`subUpdateThumbnailSubscriptionTopic`) — Thumbnail request parameters
- **Deprovision** (`subDeprovisionTopic`) — Remote deprovision command
- **Log requests** (`subUpdateLogSubscriptionTopic`) — Log upload requests

### Publications (device → host)
- **Registration** (`pubReportRegistrationTopic`) — Device capabilities, sent on connect
- **Status** (`pubReportStatusTopic`) — Device status updates
- **Actual configuration** (`pubReportActualConfigurationTopic`) — Current device configuration
- **Deprovision** (`pubDeprovisionTopic`) — Client-initiated deprovision notification

### Acceptance Criteria
- [ ] TLS connection uses mutual TLS with the device cert and CA cert
- [ ] MQTT keepalive interval is configurable via host settings (`mqttKeepaliveSeconds`)
- [ ] Auto-reconnect is enabled with max 10-second reconnect interval
- [ ] Registration is published on every MQTT connect (including reconnects)
- [ ] Status and configuration publications are throttled per `minIntervalPubSeconds` from host settings
- [ ] All subscription callbacks handle malformed payloads gracefully (log error, don't crash)

---

## Requirement 4: REST API Surface

The SDK must expose a localhost REST API for device applications to interact with the TR-12 protocol.

### Endpoints
| Method | Path | Description |
|--------|------|-------------|
| PUT | `/connect` | Initiate/continue connection with `hostId` and `registration` payload |
| PUT | `/disconnect` | Disconnect from host service |
| GET | `/get_state` | Return current connection state |
| PUT | `/report_status` | Publish device status to host |
| PUT | `/report_actual_configuration` | Publish actual configuration to host |
| GET | `/get_configuration` | Return latest configuration from host |
| PUT | `/deprovision` | Remove device from host, delete credentials. `?force=true` for offline deprovision |

### Acceptance Criteria
- [ ] All endpoints return JSON with `success`, `state`, `message`, and optional `error` fields
- [ ] `/connect` response includes `pairingCode` and `expires` when in PAIRING state
- [ ] `/connect` response includes `deviceId` and `region` when CONNECTED
- [ ] `/get_configuration` returns the latest MQTT-received configuration with an `updateId` for change detection
- [ ] CORS is enabled for all origins
- [ ] Request validation returns HTTP 400 for missing required fields

---

## Requirement 5: Credential Management

The SDK must securely persist and manage X.509 credentials on the filesystem.

### Stored Files
- `ca_cert` — CA certificate from host service
- `device_cert` — Device certificate
- `priv_key` — RSA private key (generated locally)
- `connection_settings` — JSON with device ID, MQTT URI, region
- `host_settings` — JSON with MQTT topics, keepalive, throttle intervals

### Credential Rotation
1. Host sends a `RotateCertificates` message via MQTT with new device cert and MQTT URI
2. SDK compares new cert with current cert
3. If different: persist new cert, disconnect, reconnect with new credentials
4. If same: no action

### Acceptance Criteria
- [ ] Credentials stored at `<certs_path>/<device_id>/<host_id>/`
- [ ] Private key never leaves the device (only CSR is sent during pairing)
- [ ] Credential rotation triggers a graceful reconnect (disconnect → 1s delay → connect)
- [ ] `Deprovision()` deletes the credential directory for the specified host
- [ ] File permissions are restrictive (owner read/write only where OS supports it)

---

## Requirement 6: Thumbnail Management

The SDK must handle thumbnail upload subscriptions from the host service.

### Flow
1. Host sends a `ThumbnailSubscription` via MQTT with a map of thumbnail requests
2. Each request specifies: `localPath`, `remotePath` (pre-signed URL), `period` (seconds), `expires` (unix timestamp), `maxSizeKilobyte`
3. SDK reads the image from `localPath`, validates size and freshness (< 10 seconds old), and uploads to `remotePath`
4. Upload repeats at the specified `period` until `expires`

### Acceptance Criteria
- [ ] Images older than 10 seconds are considered stale and not uploaded
- [ ] Images exceeding `maxSizeKilobyte` are not uploaded
- [ ] Each thumbnail source runs on its own goroutine
- [ ] Subscriptions are replaced (not accumulated) when a new subscription arrives
- [ ] Expired subscriptions are automatically cleaned up

---

## Requirement 7: Log Upload

The SDK must support host-requested log uploads.

### Flow
1. Host sends a `LogRequest` via MQTT with `remotePath` (pre-signed URL) and `expires` (unix timestamp)
2. SDK writes JSON-formatted rotating log files (500 KB max, 3 backups)
3. On log rotation, if an active log request exists and hasn't expired, SDK uploads the rotated file to `remotePath`

### Acceptance Criteria
- [ ] Log files are JSON-formatted with timestamps
- [ ] Log rotation at 500 KB with up to 3 backup files
- [ ] Log spew detection prevents recursive upload loops
- [ ] Expired log requests are ignored
- [ ] Logs are also printed to stdout

---

## Requirement 8: CLI Interface

The SDK binary must accept the following CLI arguments.

| Argument | Required | Description |
|----------|----------|-------------|
| `--internal_device_id` | Yes | Unique device identifier |
| `--certs_path` | Yes | Persistent directory for credential storage |
| `--log_path` | Yes | Writable directory for log files |
| `--ip` | Yes | Listen IP for REST API |
| `--port` | Yes | Listen port for REST API |
| `--device_type` | Yes | `SOURCE`, `DESTINATION`, or `BOTH` |

### Acceptance Criteria
- [ ] Missing required arguments print usage and exit with code 1
- [ ] `certs_path` is validated as writable on startup
- [ ] Graceful shutdown on SIGINT/SIGTERM (disconnect MQTT, flush logs)
- [ ] Host configuration files are resolved relative to the binary location, falling back to working directory

---

## Requirement 9: Application Reference Design (ARD)

A companion program simulating a 1-channel encoder device that exercises the SDK REST API.

### Behavior
1. Load device registration from `payloads/1_channel_encoder/registration.json`
2. Call `PUT /connect` in a loop until CONNECTED
3. Display pairing code when in PAIRING state
4. Once connected, poll `GET /get_configuration` for updates
5. Apply desired configuration via the TR-12 shim (settings, connection, state)
6. Report device status via `PUT /report_status`
7. Report actual configuration via `PUT /report_actual_configuration`
8. Simulate thumbnail emission by cycling sample images to `/tmp/image_sdi.jpg` and `/tmp/image_hdmi.jpg`
9. Manage an ffmpeg subprocess for SRT streaming when channel state is ACTIVE

### TR-12 Shim
The shim bridges TR-12 model structures to device-specific callbacks:
- Walks `DeviceConfiguration` and dispatches to update callbacks (settings, profiles, connections, state)
- Reads back actual values using `DeviceRegistration` as a template via get callbacks
- State is applied last so settings/connection are in place before start/stop

### Device Callbacks
- **Update callbacks**: Receive desired values from host (key/value settings, profiles, connections, channel state)
- **Get callbacks**: Return current device values (settings defaults, connection info, channel state, status metrics)
- **Important**: TR-12 communicates desired configuration once. The device is responsible for retrying until the desired state is achieved. The host does not re-send configuration.

### Acceptance Criteria
- [ ] ARD accepts `--host_id` (required) and `--sdk_url` (default `http://127.0.0.1:8603`)
- [ ] Registration payload loaded from JSON file matches the TR-12 `DeviceRegistration` Smithy model format
- [ ] Configuration change detection uses `updateId` to avoid reprocessing
- [ ] Thumbnail simulators cycle through sample images at 2-second intervals
- [ ] ffmpeg process starts on ACTIVE state, stops on IDLE state
- [ ] Graceful shutdown stops ffmpeg, thumbnails, and disconnects from SDK

---

## Requirement 10: Model Compatibility

Go models must be serialization-compatible with the TR-12 Smithy-generated models.

### Two Model Layers
1. **TR-12 Protocol Models** — Imported from the shared submodule (`models/TR-12-Models/generated/tr12go/`, package `openapi`). Used for pairing, auth, cert rotation, deprovision, thumbnails, logs. Both client and host import these with the alias `tr12models`.
2. **CDD SDK Models** (`client/pkg/cddmodels/`) — Registration, configuration, status (device-facing). These are generated from the CDD SDK Smithy definitions, separate from the TR-12 protocol types.

### Acceptance Criteria
- [ ] JSON serialization matches the TR-12 Smithy model schemas for all model types
- [ ] Union types (e.g., `TransportProtocol`, `SettingsChoice`, `PairResult`) serialize as discriminated unions matching the Smithy `oneOf` pattern
- [ ] Optional fields are omitted from JSON when nil (not serialized as `null`)
- [ ] Enum values match the Smithy definitions exactly (e.g., `ACTIVE`, `IDLE`, `SRT_CALLER`)

---

## Requirement 11: Host Configuration

The SDK must load host service configuration from JSON files.

### Format
```json
{
  "serviceId": "string",
  "serviceName": "string",
  "deviceTypes": ["SOURCE", "DESTINATION"],
  "thumbnailMaxSizeKB": 250,
  "logFileMaxSizeKB": 500,
  "pairingUrl": "https://...",
  "authUrl": "https://..."
}
```

### Acceptance Criteria
- [ ] Configuration files are loaded from `host_configuration/<host_id>.json`
- [ ] Path resolution: binary directory first, then working directory
- [ ] Missing or invalid configuration files produce clear error messages
- [ ] Host configuration is immutable after loading (no runtime modification)

---

## Requirement 12: Build and Distribution

The Go SDK must produce single static binaries.

### Acceptance Criteria
- [ ] `go build -o cdd-sdk ./cmd/cdd-sdk/` produces a working binary
- [ ] `go build -o ard ./cmd/ard/` produces a working ARD binary
- [ ] No CGO dependencies (pure Go for cross-compilation)
- [ ] Go 1.23+ required
- [ ] Dependencies managed via `go.mod`

---

## Requirement 13: Host Service API Design (Console-Informed)

The Go SDK's host service integration and the TR-12 protocol models must be designed with awareness that a future TR-12 web console will consume the host service API. This requirement does not spec the console itself — it ensures the host service API surface and the SDK's data contracts support a model-driven management UI.

### Host Service API Operations (consumed by both SDK and future console)

| Operation | Method | Path | Description |
|-----------|--------|------|-------------|
| ListDevices | GET | `/devices` | List all registered devices with online status |
| DescribeDevice | GET | `/device/{deviceId}` | Full device state: registration, configuration, status, online details, cert expiration |
| Authorize/Claim | PUT | `/authorize/{pairingCode}` | Claim a pairing device into the registry |
| Deprovision | PUT | `/deprovision/{deviceId}` | Remove device from registry |
| UpdateConfiguration | PUT | `/device/{deviceId}` | Push desired configuration to device |
| GetThumbnail | GET | `/thumbnail/{deviceId}?source={sourceId}` | Retrieve latest thumbnail for a source |
| RotateCredentials | PUT | `/credentials/{deviceId}` | Trigger certificate rotation |

### Model-Driven UI Contract

The registration payload (`DeviceRegistration` from `registration.smithy`) serves as the schema for dynamically generating console UI pages. The SDK and host service must preserve the full fidelity of registration data so a console can:

1. **Enumerate channels** — Each `Channel` in the registration defines a manageable unit with its own settings, profiles, and connection protocols
2. **Render type-appropriate widgets** from `Setting` structures:
   - `EnumValues` (with `values` list and `defaultValue`) → dropdown/select widget
   - `RangeValues` (with `min`, `max`, `defaultValue`) → slider or numeric input widget
   - Settings with neither → free-text input
3. **Support profiles** — `ProfileDefinition` entries (with `name`, `id`, `info`) allow profile-based configuration as an alternative to individual settings
4. **Show connection protocols** — `ProtocolList` on each channel advertises supported transport protocols (SRT_LISTENER, SRT_CALLER, RIST, ZIXI variants), informing which connection configuration forms to present
5. **Display actual vs desired configuration** — The console will compare `DeviceConfiguration` (desired, from host) against the actual configuration reported by the device via `report_actual_configuration`
6. **Show device and channel status** — `DeviceStatus` and `ChannelStatus` with `StatusValue` entries (name, info, value) render as read-only status panels
7. **Show thumbnails** — `ThumbnailList` in registration defines available thumbnail sources with IDs and names

### Acceptance Criteria
- [ ] SDK serializes registration, configuration, and status payloads as structured JSON matching the Smithy model schemas (not flattened or transformed)
- [ ] Host service API operations listed above are compatible with the SDK's pairing, MQTT, and REST flows
- [ ] Registration `Setting` structures preserve `enums`, `ranges`, and `defaultValue` fields for downstream UI rendering
- [ ] Registration `ProfileDefinition` entries preserve `name`, `id`, and `info` fields
- [ ] Channel `connectionProtocols` list is included in registration and matches the `TransportProtocol` union variants in configuration
- [ ] Actual configuration reported by the SDK uses the same `DeviceConfiguration` schema as desired configuration, enabling field-by-field comparison


---

## Unit Test Coverage

Unit tests are placed alongside source files (idiomatic Go convention) and run with standard `go test`:

```bash
go test ./client/... -count=1
```

### Covered Packages

| Package | Test File | Tests | What's Covered |
|---------|-----------|-------|----------------|
| `internal/utils` | `utils_test.go` | 12 | Path validation, key generation, CSR creation, host config loading, throttle timing, update ID generation, error detail extraction |
| `internal/credentials` | `store_test.go` | 12 | Constructor, getters (nil/populated), key gen idempotency, filesystem read/write round-trip, cert rotation (changed/unchanged), deprovision cleanup |
| `internal/cddlogger` | `logger_test.go` | 10 | Creation, JSON format, Info/Error/Errorf/Infof/Exception levels, device ID update, log rotation at 500KB, Dump, Close |

All tests use the Go standard library only — no third-party test frameworks.
