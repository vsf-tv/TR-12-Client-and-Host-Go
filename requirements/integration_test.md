# TR-12 Integration Tests: Requirements

## Overview

Automated integration tests that exercise the full TR-12 protocol lifecycle by running the host service and SDK as real processes and interacting with them via their HTTP APIs. The test harness acts as both the device application (calling SDK endpoints, like the ARD does) and the host operator (calling host management API endpoints, like a user with curl/Postman).

Tests use Go's standard `testing` package with `go test`. Each test starts a fresh host service and SDK process, runs the protocol lifecycle, asserts on responses, and tears everything down. No mocks — these are true end-to-end integration tests against real binaries.

### Location
- Test code: `test/integration/`
- Go module: `github.com/vsf-tv/TR-12-Client-and-Host-Go/test/integration`
- Added to the root `go.work` so `go test ./...` works from the test directory

### Reference Materials
- Host API endpoints: `host/internal/api/router.go`
- SDK API endpoints: `client/internal/api/server.go`
- Registration payload: `client/payloads/1_channel_encoder/registration.json`
- Host CLI flags: `host/internal/config/config.go`
- SDK CLI flags: `client/cmd/cdd-sdk/main.go`
- Project requirements: `requirements/project.md`

---

## Requirement 1: Test Infrastructure — Process Management

The test harness must start and stop the host service and SDK as real OS processes, using pre-built binaries.

### Setup (per test)
1. Build the `tr12-host` and `cdd-sdk` binaries (or skip if already built) using `go build` from their respective module directories
2. Start the host service process with a unique `--db-path` (temp file), a unique `--http-port`, a unique `--mqtt-port`, and `--host-address 127.0.0.1`
3. Wait for the host service to be ready by polling its HTTP port (e.g., `GET /host-config` or a TCP connect check) with a timeout
4. Start the SDK process with a unique `--internal_device_id`, a temp `--certs_path`, a temp `--log_path`, `--ip 127.0.0.1`, a unique `--port`, and `--device_type SOURCE`
5. Wait for the SDK to be ready by polling its HTTP port (e.g., `GET /get_state`) with a timeout

### Teardown (per test)
1. Send SIGTERM to the SDK process
2. Send SIGTERM to the host service process
3. Wait for both processes to exit (with a kill timeout)
4. Remove temp directories (certs, logs, DB file)

### Port Allocation
Each test run uses ephemeral ports to avoid conflicts with other tests or local development. The harness picks free ports at runtime (e.g., by binding to `:0` and reading the assigned port, or by using a port range starting from a high base like 19000 + random offset).

### Acceptance Criteria
- [ ] Host service and SDK are started as real OS processes (not in-process)
- [ ] Each test gets a fresh database, fresh certs directory, and fresh SDK state — no shared state between tests
- [ ] Process startup is verified with a health check poll (max 10 seconds)
- [ ] Processes are reliably cleaned up even if the test panics (use `t.Cleanup`)
- [ ] Temp files (DB, certs, logs) are removed after each test
- [ ] Port conflicts between parallel test runs are avoided
- [ ] Build step is skipped if binaries already exist and are up to date (use `go build` which is a no-op if unchanged)

---

## Requirement 2: Test Infrastructure — HTTP Client Helpers

The test harness must provide helper functions for calling both the SDK and host service APIs, with JSON marshaling/unmarshaling and error assertion.

### SDK Client Helpers
Wrappers around the SDK's localhost REST API:
- `sdkConnect(hostID string, registration map[string]interface{}) ConnectResponse` — `PUT /connect`
- `sdkGetState() StateResponse` — `GET /get_state`
- `sdkGetConfiguration() ConfigResponse` — `GET /get_configuration`
- `sdkReportStatus(status map[string]interface{}) Response` — `PUT /report_status`
- `sdkReportActualConfig(config map[string]interface{}) Response` — `PUT /report_actual_configuration`
- `sdkDisconnect() Response` — `PUT /disconnect`
- `sdkDeprovision(hostID string) Response` — `PUT /deprovision`

### Host Client Helpers
Wrappers around the host service's management API:
- `hostRegisterAccount(username, password, displayName string) AccountResponse` — `POST /account/register`
- `hostLogin(username, password string) AccountResponse` — `POST /account/login`
- `hostClaim(pairingCode string, token string) Response` — `PUT /authorize/{code}`
- `hostListDevices(token string) []DeviceSummary` — `GET /devices`
- `hostDescribeDevice(deviceID, token string) DeviceDetail` — `GET /device/{id}`
- `hostUpdateConfig(deviceID, token string, config json.RawMessage) Response` — `PUT /device/{id}`
- `hostGetThumbnail(deviceID, sourceID, token string) ThumbnailResponse` — `GET /thumbnail/{id}?source={sourceId}`
- `hostRotateCredentials(deviceID, token string) Response` — `PUT /credentials/{id}`
- `hostDeprovision(deviceID, token string) Response` — `PUT /deprovision/{id}`

### Polling Helpers
- `waitForSDKState(desiredState string, timeout time.Duration) error` — polls `GET /get_state` until the SDK reaches the desired state or times out
- `waitForSDKConnected(timeout time.Duration) error` — convenience wrapper for `waitForSDKState("CONNECTED", ...)`
- `waitForSDKPairing(timeout time.Duration) (pairingCode string, err error)` — polls `PUT /connect` until the SDK returns a pairing code

### Acceptance Criteria
- [ ] All helpers return typed Go structs (not raw `[]byte`)
- [ ] All helpers accept `*testing.T` and call `t.Fatalf` on unexpected HTTP errors (non-2xx when 2xx expected)
- [ ] Helpers include configurable timeouts (default 5 seconds for single calls, 30 seconds for polling)
- [ ] Polling helpers use exponential backoff or fixed interval (500ms) with a max timeout
- [ ] Response structs match the actual JSON shapes returned by the SDK and host APIs

---

## Requirement 3: Full Lifecycle Test

A single test function (`TestFullLifecycle`) that exercises the complete TR-12 protocol from pairing through deprovision. This is the primary integration test.

### Test Flow and Assertions

#### Phase 1: Account Setup
1. Call `POST /account/register` on the host with username `testuser`, password `testpass123`, display name `Test User`
2. Assert: response contains a non-empty `token` and `account_id` starting with `acc_`

#### Phase 2: Device Pairing
1. Load the 1-channel encoder registration payload from `client/payloads/1_channel_encoder/registration.json`
2. Call `PUT /connect` on the SDK with `hostId` = the host's service ID (default `tr12-host`) and the registration payload
3. Assert: response `state` is `PAIRING` and `pairingCode` is a 6-character uppercase alphanumeric string
4. Assert: `expires` is a positive number (pairing timeout in seconds)

#### Phase 3: Device Claim
1. Call `PUT /authorize/{pairingCode}` on the host with the JWT token from Phase 1
2. Assert: response is HTTP 200

#### Phase 4: SDK Connects
1. Poll `GET /get_state` on the SDK until state is `CONNECTED` (timeout: 30 seconds)
2. Assert: state reaches `CONNECTED`

#### Phase 5: List Devices
1. Call `GET /devices` on the host with the JWT token
2. Assert: response is a JSON array with exactly 1 entry
3. Assert: the entry has a non-empty `device_id` (21-char alphanumeric), `online` is `true`

#### Phase 6: Describe Device
1. Call `GET /device/{deviceId}` on the host with the device ID from Phase 5
2. Assert: `registration` is non-null and contains `channels` array with 1 entry where `id` = `CH01`
3. Assert: `registration.channels[0].simpleSettings` has 7 entries (resolution, framerate, max_bitrate, rate_control, codec, gop_size, selected_input)
4. Assert: `registration.thumbnails` has 2 entries (`SDI-1` and `HDMI-1`)
5. Assert: `online` is `true`
6. Assert: `device_metadata.device_type` is `SOURCE`
7. Assert: `device_metadata.account_id` matches the account ID from Phase 1
8. Assert: `cert_expiration` is a non-empty string (e.g., `29d ...`)

#### Phase 7: Update Configuration

##### 7a: Negative — Unknown Channel ID
1. Send a configuration with a non-existent channel ID:
   ```json
   {
     "channels": [
       { "id": "BOGUS_CHANNEL", "state": "ACTIVE" }
     ]
   }
   ```
2. Call `PUT /device/{deviceId}` on the host
3. Assert: HTTP 400 with error message containing `"unknown channel ID"` and listing `CH01` as a valid option

##### 7b: Negative — Unknown Setting Key
1. Send a configuration referencing a setting ID that doesn't exist in the registration:
   ```json
   {
     "channels": [
       {
         "id": "CH01",
         "state": "ACTIVE",
         "settings": {
           "simpleSettings": [
             {"key": "NONEXISTENT_SETTING", "value": "foo"}
           ]
         }
       }
     ]
   }
   ```
2. Call `PUT /device/{deviceId}` on the host
3. Assert: HTTP 400 with error message containing `"unknown setting key"` and `"NONEXISTENT_SETTING"`

##### 7c: Negative — Unknown Profile ID
1. Send a configuration referencing a profile ID that doesn't exist in the registration:
   ```json
   {
     "channels": [
       {
         "id": "CH01",
         "state": "ACTIVE",
         "settings": {
           "profile": { "id": "nonexistent_profile" }
         }
       }
     ]
   }
   ```
2. Call `PUT /device/{deviceId}` on the host
3. Assert: HTTP 400 with error message containing `"unknown profile ID"` and `"nonexistent_profile"`

##### 7d: Positive — Full Configuration with Device-Level Settings, Channel Settings, and SRT Caller Connection
1. Build a comprehensive valid configuration payload that exercises device-level settings, channel-level simple settings, and an SRT caller connection:
   ```json
   {
     "simpleSettings": [
       {"key": "sync_clock_source", "value": "PTP"}
     ],
     "channels": [
       {
         "id": "CH01",
         "state": "ACTIVE",
         "settings": {
           "simpleSettings": [
             {"key": "RS01", "value": "1920x1080"},
             {"key": "FR01", "value": "60"},
             {"key": "MB01", "value": "20000"},
             {"key": "RC01", "value": "CBR"},
             {"key": "CO01", "value": "H.264"},
             {"key": "GP01", "value": "60"},
             {"key": "IN01", "value": "SDI1"}
           ]
         },
         "connection": {
           "transportProtocol": {
             "srtCaller": {
               "address": "192.168.1.100",
               "port": 9000,
               "latency": 200
             }
           }
         }
       }
     ]
   }
   ```
2. Call `PUT /device/{deviceId}` on the host with this payload
3. Assert: HTTP 200 with `message` containing `"updated"` (case-insensitive)
4. Wait 2 seconds for MQTT delivery
5. Call `GET /get_configuration` on the SDK
6. Assert: the returned configuration contains `channels[0].id` = `CH01` and `channels[0].state` = `ACTIVE`
7. Assert: the returned configuration contains the SRT caller connection with `address` = `192.168.1.100` and `port` = `9000`
8. Assert: the returned configuration contains channel simple settings including `RS01` = `1920x1080` and `FR01` = `60`

#### Phase 8: Report Status and Actual Configuration
1. Build a status payload:
   ```json
   {
     "channels": [
       {
         "id": "CH01",
         "state": "ACTIVE",
         "statusValues": [
           {"name": "bitrate", "value": "9500", "info": "Current output bitrate (Kbps)"}
         ]
       }
     ]
   }
   ```
2. Call `PUT /report_status` on the SDK
3. Assert: response `success` is `true`
4. Build an actual configuration payload matching the desired config from Phase 7d
5. Call `PUT /report_actual_configuration` on the SDK
6. Assert: response `success` is `true`
7. Wait 2 seconds for MQTT delivery to host
8. Call `GET /device/{deviceId}` on the host
9. Assert: `status` is non-null and contains the reported status values
10. Assert: `actual_configuration` is non-null and contains the reported configuration

#### Phase 9: Thumbnail Request
1. Before the test starts the SDK, create a small test JPEG file at `/tmp/image_sdi.jpg` (a valid JPEG, can be a 1x1 pixel image)
2. Call `GET /thumbnail/{deviceId}?source=SDI-1` on the host
3. The first call may return a "pending" response (subscription just created) — if so, wait 5 seconds and retry
4. Assert: eventually the response contains `image_data` (base64-encoded, non-empty) and `image_type` = `jpg`

#### Phase 10: Credential Rotation
1. Call `GET /device/{deviceId}` on the host and record the current `cert_expiration`
2. Call `PUT /credentials/{deviceId}` on the host
3. Assert: response is HTTP 200
4. Wait 5 seconds for the SDK to receive the rotation message, reconnect, and re-establish
5. Call `GET /get_state` on the SDK
6. Assert: state is `CONNECTED` (SDK reconnected successfully after rotation)
7. Call `GET /device/{deviceId}` on the host
8. Assert: `cert_expiration` has changed (new cert has a fresh expiry)
9. Assert: `online` is `true`

#### Phase 11: Deprovision
1. Call `PUT /deprovision/{deviceId}` on the host
2. Assert: response is HTTP 200
3. Wait 3 seconds for MQTT delivery
4. Call `GET /get_state` on the SDK
5. Assert: state is `DISCONNECTED` (SDK received deprovision and disconnected)

### Acceptance Criteria
- [ ] Test runs with `go test -v -timeout 120s ./...` from the `test/integration/` directory
- [ ] Test passes on a clean machine with no pre-existing state (fresh DB, fresh certs)
- [ ] All 11 phases execute sequentially within a single test function
- [ ] Each phase has clear assertion messages identifying what failed and the actual vs expected values
- [ ] Test completes in under 90 seconds on a typical development machine
- [ ] Test creates and cleans up a test JPEG file for thumbnail testing
- [ ] Test is tagged with `//go:build integration` so it doesn't run during normal `go test` — requires `-tags integration` or an explicit `go test -run TestFullLifecycle`

---

## Requirement 4: Test Configuration and Build Tags

The integration tests must be isolated from normal unit test runs and configurable for different environments.

### Build Tags
- All integration test files use the build constraint `//go:build integration`
- Normal `go test ./...` from the repo root does NOT run integration tests
- Integration tests are run explicitly: `go test -tags integration -v -timeout 120s ./test/integration/`

### Environment Variables (optional overrides)
| Variable | Default | Description |
|---|---|---|
| `TR12_HOST_BINARY` | (auto-built) | Path to pre-built `tr12-host` binary |
| `TR12_SDK_BINARY` | (auto-built) | Path to pre-built `cdd-sdk` binary |
| `TR12_TEST_TIMEOUT` | `90s` | Overall test timeout |
| `TR12_HOST_ADDRESS` | `127.0.0.1` | Host address for both services |

### Acceptance Criteria
- [ ] Integration tests do not run with plain `go test ./...`
- [ ] Integration tests run with `go test -tags integration ./test/integration/`
- [ ] Pre-built binaries can be provided via environment variables (useful in CI)
- [ ] If binaries are not provided, the test harness builds them automatically using `go build`
- [ ] Test timeout is configurable but defaults to 90 seconds

---

## Requirement 5: Test Module Structure

The integration test module must be a separate Go module within the monorepo, added to the Go workspace.

### File Layout
```
test/
└── integration/
    ├── go.mod                    # module github.com/vsf-tv/TR-12-Client-and-Host-Go/test/integration
    ├── go.sum
    ├── helpers_test.go           # Process management, HTTP client helpers, polling utilities
    ├── lifecycle_test.go         # TestFullLifecycle — the main integration test
    ├── README.md                 # Instructions for running the integration tests
    └── testdata/
        └── registration.json     # Copy of client/payloads/1_channel_encoder/registration.json
```

### Module Dependencies
The test module does NOT import the client or host Go packages directly. It interacts with them only via HTTP and process management. Dependencies are minimal:
- Go standard library (`testing`, `net/http`, `os/exec`, `encoding/json`, `time`, etc.)
- No third-party test frameworks (no testify, no gomega) — use standard `t.Errorf`/`t.Fatalf`

### Go Workspace
The root `go.work` is updated to include the test module:
```
use (
    ./client
    ./host
    ./models/TR-12-Models/generated/tr12go
    ./test/integration
)
```

### Acceptance Criteria
- [x] Test module has its own `go.mod` with no dependencies on client or host modules
- [x] Test module is added to the root `go.work`
- [x] Test files use `package integration_test` (external test package)
- [x] Registration payload is copied into `testdata/` so the test is self-contained
- [x] No third-party test dependencies — standard library only
- [x] `go vet` and `go build` pass on the test module

---

## Implementation Status

All 5 requirements are fully implemented and passing. The integration test (`TestFullLifecycle`) completes all 11 phases in ~24 seconds on a typical Mac.

```bash
# Run from the test/integration/ directory:
go test -tags integration -v -timeout 120s ./...

# Or from the repo root:
go test -tags integration -v -timeout 120s ./test/integration/
```

See `test/integration/README.md` for full details on running the tests, environment variables, and how the test harness works.
