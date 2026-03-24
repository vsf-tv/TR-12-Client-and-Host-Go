# TR-12 Client and Host — Go: Project Requirements

## Overview

This repository is a complete Go implementation of the VSF TR-12 Client Device Discovery protocol. It contains both sides of the protocol — the device-side client SDK with an Application Reference Design (ARD), and a self-contained host service — in a single monorepo with shared Smithy-generated protocol models.

TR-12 solves the problem of discovering, managing, and monitoring professional streaming video devices from a scalable service. While many transport protocols exist for securely streaming video (SRT, RIST, TR-07, etc.), production workflows involving distributed sources and destinations still require manual management. TR-12 provides a mechanism to securely pair devices into a registry, after which they can be managed and monitored anywhere in the world across the open internet.

The protocol uses modern, cloud-first solutions for security (mTLS, X.509 certificates, JWT), modeling (Smithy), validation, and resource lifecycle management — applied specifically for video streaming devices. A device installs the SDK, integrates with its native control plane using the provided models, and pairs with a host service of the operator's choice. The protocol accommodates the widely differentiated settings (codec, channels, transport protocols, etc.) available across different device types (encoders, decoders, cameras, playout devices) from different manufacturers. Devices expose completely customized settings within the protocol's structure.

### Protocol Reference
- TR-12 Smithy Models: https://github.com/vsf-tv/TR-12-Models
- Draft Protocol Specification: https://github.com/vsf-tv/TR-12-Models/blob/main/VSF_TR-12-ClientDeviceDiscoveryDraft.pdf
- VSF TR-12 Working Group: Contact Brad Gilmer or Brian Rundle for access to the VSF Bi-Weekly Forum

---

## Requirement 1: Monorepo Structure

The project must be organized as a single Git repository containing both the client and host implementations, with shared protocol models provided via a Git submodule.

### Components
| Directory | Component | Description |
|---|---|---|
| `client/` | CDD SDK + ARD | Device-side SDK daemon and Application Reference Design (simulated 1-channel encoder) |
| `host/` | TR-12 Host Service | Self-contained host service with REST API, embedded MQTT broker, and SQLite persistence |
| `models/TR-12-Models/` | Shared Models | Git submodule containing Smithy-generated Go protocol types (`package openapi`) |

### Go Workspace
A root-level `go.work` file links all three Go modules (`client`, `host`, `models/TR-12-Models/generated/tr12go`) so standard `go build` works from each component directory without manual `replace` directive management.

### Acceptance Criteria
- [ ] Repository contains `client/`, `host/`, and `models/` directories at the top level
- [ ] `go.work` at the repo root links all three Go modules
- [ ] Each component (`client/`, `host/`) has its own `go.mod` with a `replace` directive pointing to the local submodule path
- [ ] `git clone --recurse-submodules` pulls the complete repo including the TR-12 models submodule
- [ ] All three binaries (`cdd-sdk`, `ard`, `tr12-host`) compile successfully from their respective directories
- [ ] The submodule content (`models/TR-12-Models/`) is never modified by this repo — it is a 100% external dependency

---

## Requirement 2: Shared Smithy-Generated Protocol Models

Both client and host must consume the same set of Smithy-generated Go types for all TR-12 protocol messages. No hand-rolled duplicates of protocol types are permitted.

### Model Source
- Smithy definitions: `models/TR-12-Models/smithy/*.smithy` (5 files: common, configuration, registration, service_api, status)
- Generator: `models/TR-12-Models/generate-tr12-models.sh`
- Generated Go output: `models/TR-12-Models/generated/tr12go/` (package `openapi`)

### Import Pattern
Both client and host import the generated types with an alias:
```go
import tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
```

### Key Protocol Types
| Type | Smithy Source | Used By |
|---|---|---|
| `PairRequestContent` / `PairResponseContent` | `PairRequest` / `PairResponse` | Client (sends), Host (receives/responds) |
| `AuthenticateRequestContent` / `AuthenticateResponseContent` | `AuthRequest` / `AuthResponse` | Client (polls), Host (responds) |
| `HostSettings` | `HostSettings` | Host (generates), Client (consumes) |
| `RotateCertificatesRequestContent` | `CertRotate` | Host (publishes via MQTT), Client (receives) |
| `DeprovisionDeviceRequestContent` | `DeprovisionMessage` | Bidirectional via MQTT |
| `RequestThumbnailRequestContent` | `ThumbnailSubscription` | Host (publishes), Client (receives) |
| `RequestLogRequestContent` | `LogRequest` | Host (publishes), Client (receives) |

### Acceptance Criteria
- [ ] Both client and host import all TR-12 protocol types from the shared submodule — zero vendored copies
- [ ] JSON serialization of protocol messages is identical between client and host (same structs, same field names)
- [ ] Updating the upstream Smithy models requires only pulling the submodule and rebuilding — no code changes in client or host
- [ ] The client additionally has its own `pkg/cddmodels/` for CDD SDK-specific models (registration, configuration, status) that are separate from the TR-12 protocol types

---

## Requirement 3: Build and Distribution

All components must produce single static Go binaries with no CGO dependencies, enabling cross-compilation for macOS and Linux.

### Binaries
| Binary | Build Command | Directory | Description |
|---|---|---|---|
| `cdd-sdk` | `go build -o cdd-sdk ./cmd/cdd-sdk` | `client/` | SDK daemon (~13 MB) |
| `ard` | `go build -o ard ./cmd/ard` | `client/` | Application Reference Design |
| `tr12-host` | `go build -o tr12-host ./cmd/tr12-host` | `host/` | Host service |

### Acceptance Criteria
- [ ] All three binaries compile with `go build` from their respective directories
- [ ] No CGO dependencies — pure Go for cross-compilation
- [ ] Client requires Go 1.23+; Host requires Go 1.24+
- [ ] Dependencies managed via `go.mod` per module
- [ ] Built binaries are excluded from version control via `.gitignore`

---

## Requirement 4: End-to-End Protocol Flow

The client and host must implement the complete TR-12 protocol lifecycle when running together. This is the primary integration requirement — both sides must interoperate correctly.

### Lifecycle Phases

1. **Account Setup** — Operator registers an account on the host service and receives a JWT token
2. **Device Pairing** — SDK calls `POST /pair` on the host, receives a pairing code; operator claims the device via `PUT /authorize/{code}` with their JWT
3. **Authentication** — SDK polls `POST /authenticate` until the host returns `CLAIMED` with CA cert, device cert, MQTT URI, and host settings
4. **MQTT Connection** — SDK establishes mTLS connection to the host's embedded MQTT broker using the issued certificates
5. **Registration** — SDK publishes the device's `DeviceRegistration` payload (channels, settings, profiles, thumbnails) to the host via MQTT
6. **Operational Loop** — Host pushes desired configuration via MQTT; device reports actual configuration and status back; thumbnails and logs are exchanged on demand
7. **Certificate Rotation** — Host auto-rotates device certificates every 30 days (configurable) with credential overlap to prevent connection drops
8. **Deprovision** — Two-phase process: host marks device as DEPROVISIONED and notifies via MQTT; device acknowledges; host performs full cleanup

### Acceptance Criteria
- [ ] A fresh client can pair with a fresh host, complete authentication, establish MQTT, and enter the operational loop without manual intervention beyond claiming the pairing code
- [ ] Configuration updates pushed from the host reach the SDK via MQTT and are applied by the ARD
- [ ] Status and actual configuration reports from the SDK reach the host and are visible via `GET /device/{id}`
- [ ] Thumbnail requests from the host result in image uploads from the SDK, retrievable via `GET /thumbnail/{id}?source={sourceId}`
- [ ] Certificate rotation completes without dropping the device connection (credential overlap)
- [ ] Deprovision from the host cleanly disconnects the device and cleans up all state
- [ ] The SDK can reconnect after restart using persisted credentials (no re-pairing needed)

---

## Requirement 5: Security Model

The protocol uses a layered security model with no shared secrets transmitted in plaintext.

### Certificate Authority
- The host service operates as a local CA (RSA 4096-bit, 10-year validity)
- Device certificates are signed by the CA during pairing (RSA, configurable expiry, default 30 days)
- The CA cert, CA key, server cert, server key, and JWT secret are all stored in the SQLite database — no filesystem paths for secrets
- Each certificate signing uses a cryptographically random 128-bit serial number

### Mutual TLS
- The embedded MQTT broker requires client certificate authentication
- Device connections are validated against the service CA
- Per-device topic ACLs prevent cross-device access

### Credential Rotation
- Auto-rotation every 30 days (configurable) as a background goroutine
- Rotation messages are published as retained MQTT messages so offline devices receive them on reconnect
- Credential overlap: both current and previous cert sets are valid simultaneously until the device connects with the new cert
- The CSR stays the same between rotations — the cert differs due to random serial and timestamps

### Account Authentication
- Management API endpoints require JWT Bearer tokens
- Passwords stored as bcrypt hashes
- JWT signing secret auto-generated on first run and persisted in the database

### Acceptance Criteria
- [ ] Private keys never leave the device (only CSR is sent during pairing)
- [ ] All secrets (CA, server cert, JWT secret) are stored in the SQLite database, not on the filesystem
- [ ] mTLS is enforced on the MQTT broker
- [ ] Certificate rotation does not cause device disconnection (overlap period)
- [ ] Management API rejects requests without valid JWT tokens (HTTP 401)
- [ ] Device-facing endpoints (`/pair`, `/authenticate`) are unauthenticated by design

---

## Requirement 6: Persistence Model

All state is stored in a single SQLite database file. This is the ONLY filesystem artifact produced by the host service.

### What's in the Database
- Account data (usernames, bcrypt password hashes)
- Device registry and state (registration, configuration, status, online state, cert metadata)
- CA certificate, CA private key, server certificate, server private key
- JWT signing secret
- Device certificates (current and previous per device)
- Thumbnails (as BLOBs, latest per device per source)
- Device logs (as BLOBs, latest per device)

### Client-Side Persistence
The SDK persists credentials on the filesystem at `<certs_path>/<internal_device_id>/<host_id>/`:
- `ca_cert` — CA certificate from host
- `device_cert` — Device certificate
- `priv_key` — RSA private key (generated locally)
- `connection_settings` — JSON with device ID, MQTT URI, region
- `host_settings` — JSON with MQTT topics, keepalive, throttle intervals

### Acceptance Criteria
- [ ] Host service: single `.db` file is the only filesystem artifact — portable (copy the file, move the service)
- [ ] Host service: SQLite via `modernc.org/sqlite` (pure Go, no CGO)
- [ ] Client: credentials stored at `<certs_path>/<device_id>/<host_id>/` with restrictive file permissions
- [ ] Client: different `--internal_device_id` values create separate credential sets with no overlap

---

## Requirement 7: Model-Driven Console Support

The host service API must preserve full Smithy model fidelity to support a future model-driven web console. The `DeviceRegistration` payload serves as the UI schema for dynamically generating console pages.

### Console Contract
- `Channel` entries → tabbed/sectioned device management pages
- `Setting` with `EnumValues` → dropdown/select widgets
- `Setting` with `RangeValues` → slider/numeric input widgets
- `ProfileDefinition` entries → profile selector
- `connectionProtocols` list → connection configuration form variants
- `Thumbnail` entries → thumbnail viewer panels
- `DeviceConfiguration` (desired) vs actual → side-by-side comparison view
- `DeviceStatus` with `StatusValue` entries → read-only status panels

### Acceptance Criteria
- [ ] `DescribeDevice` returns `registration`, `configuration` (desired), `actual_configuration`, and `status` as top-level structured JSON fields
- [ ] No lossy transformations are applied to any Smithy-modeled payload between device report and API response
- [ ] The `updateId` on configuration is preserved end-to-end for stale-state detection
- [ ] Registration `Setting` structures preserve `enums`, `ranges`, and `defaultValue` fields

---

## Requirement 8: Development Workflow

The repository must support a straightforward local development workflow where all three components run on the same machine.

### Local Development Setup
1. Clone with `--recurse-submodules`
2. Build all three binaries from their respective directories
3. Start the host service (Terminal 1): `./tr12-host --host-address 127.0.0.1`
4. Start the SDK (Terminal 2): `./cdd-sdk --internal_device_id test001 --certs_path ~/certs --log_path /tmp --ip 127.0.0.1 --port 8603 --device_type SOURCE`
5. Start the ARD (Terminal 3): `./ard --host_id local_go_host`
6. Register an account, claim the pairing code, and interact via the management API

### Host Configuration
The SDK loads host service connection details from `host_configuration/<host_id>.json`. The included `local_go_host.json` points to `127.0.0.1:8080` for local development. The `vsf_test_host.json` points to the VSF cloud test endpoint for interop testing.

### Acceptance Criteria
- [ ] A developer can go from `git clone` to a fully operational local system in under 5 minutes
- [ ] The `local_go_host.json` config file is tracked in version control for out-of-the-box local development
- [ ] User-specific config files (e.g., `dev_*.json`, `vsf_*.json`) are gitignored
- [ ] Built binaries, SQLite databases, and Go workspace sum files are gitignored

---

## Implementation Notes: Integration Fixes

During integration testing of the client and host, several protocol-level issues were discovered and resolved. These are documented here so future implementors understand the edge cases:

1. **MQTT URI Scheme** — The SDK's MQTT client library expects `tls://host:port` for TLS connections. The host service must return `tls://` (not `mqtts://` or `ssl://`) in both the `AuthResponse.mqttUri` and `RotateCertificatesRequestContent.mqttUri` fields. The SDK's `parseBrokerAddress()` function strips the scheme and passes the raw `host:port` to the MQTT client.

2. **NULL Scan on Optional Columns** — The host's SQLite device table has many nullable TEXT columns (pairing_code, access_code, registration, etc.). Go's `database/sql` requires `sql.NullString` for scanning nullable columns — using plain `string` causes a NULL scan error. All nullable device columns use `sql.NullString` in the scan.

3. **Pairing Field Preservation on Claim** — When a device is claimed via `PUT /authorize/{code}`, the host must preserve the `pairing_code`, `access_code`, and `pairing_expires_at` fields. The SDK continues polling `POST /authenticate` with these values after the claim. Clearing them on claim causes the SDK to get a 401 on its next authenticate poll.

4. **Thumbnail localPath Resolution** — The host must resolve the `localPath` for each thumbnail source from the device's registration payload (the `thumbnails[].localPath` field). Without this, the SDK doesn't know which file to read and upload. The `maxSizeKilobyte` must also be large enough to accommodate typical images (the ARD's samples are ~131 KB, so 100 KB is too small — 500 KB is the default).

5. **Credential Rotation Overlap** — During rotation, both the current and previous certificate sets must be valid simultaneously. The host tracks `current_cert_pem` and `previous_cert_pem` per device. When a device connects with the new cert, the old cert is revoked. This prevents dropping devices that connect with slightly stale credentials during a rotation window.

6. **Configuration Update via MQTT** — The host publishes desired configuration to `cdd/{deviceId}/config/update`. The SDK subscribes to this topic and expects the payload to be a raw `DeviceConfiguration` JSON object (not wrapped in `{"message": ...}`). The host must publish the configuration directly, not wrapped in a `ReportMessage` envelope.

---

## Testing Infrastructure

The repository includes comprehensive test coverage at two levels:

### Unit Tests
Unit tests are placed alongside source files (idiomatic Go convention) using the standard `testing` package. No third-party test frameworks.

```bash
# Run all unit tests
go test ./client/... ./host/... -count=1
```

Coverage spans 6 test files with ~80 test cases across both client and host:
- Client: `internal/utils`, `internal/credentials`, `internal/cddlogger`
- Host: `internal/db`, `internal/service` (account + device)

See `requirements/client.md` and `requirements/host.md` for per-package details.

### Integration Tests
A separate Go module at `test/integration/` runs the full TR-12 protocol lifecycle as a black-box test — starting the host service and SDK as real OS processes and interacting via HTTP.

```bash
# Run integration tests (requires -tags integration)
go test -tags integration -v -timeout 120s ./test/integration/
```

The `TestFullLifecycle` test covers 11 phases: account setup, pairing, claim, MQTT connection, device listing, describe, configuration push (3 negative + 1 positive), status reporting, thumbnails, credential rotation, and deprovision. Typical runtime ~24 seconds.

See `test/integration/README.md` and `requirements/integration_test.md` for full details.
