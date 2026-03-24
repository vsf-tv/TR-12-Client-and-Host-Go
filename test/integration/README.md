# TR-12 Integration Tests

End-to-end integration tests that exercise the full TR-12 protocol lifecycle by running the host service and SDK as real OS processes and interacting with them via their HTTP APIs.

No mocks — these are true black-box tests against real binaries.

## Prerequisites

- Go 1.24+
- Both `client/` and `host/` modules must compile (the test harness builds them automatically)

## Running

```bash
cd test/integration
go test -tags integration -v -timeout 120s -count=1 ./...
```

Or from the repo root:

```bash
go test -tags integration -v -timeout 120s -count=1 ./test/integration/
```

The `-count=1` flag disables Go's test caching, which is important for integration tests that depend on external process state.

The `-tags integration` flag is required. Without it, the tests are skipped (all files use `//go:build integration`).

## What It Tests

`TestFullLifecycle` runs 11 sequential phases covering the complete TR-12 protocol:

| Phase | What happens |
|-------|-------------|
| 1. Account Setup | Register an account on the host, get a JWT |
| 2. Device Pairing | SDK calls `/connect`, receives a 6-char pairing code |
| 3. Device Claim | Operator claims the device via `PUT /authorize/{code}` |
| 4. SDK Connects | SDK polls `/authenticate`, gets certs, establishes MQTT |
| 5. List Devices | Verify the device appears in `GET /devices` |
| 6. Describe Device | Verify full registration, metadata, cert expiration |
| 7. Update Configuration | 3 negative cases (bad channel, bad setting, bad profile) + 1 full positive config push with SRT caller, verified on SDK side |
| 8. Report Status | SDK reports status and actual config, verified on host side |
| 9. Thumbnail Request | Host requests a thumbnail, SDK uploads it, host returns base64 image |
| 10. Credential Rotation | Host rotates device certs, SDK reconnects with new certs |
| 11. Deprovision | Host deprovisions device, SDK transitions to DISCONNECTED |

Typical runtime: ~24 seconds on a modern Mac.

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `TR12_HOST_BINARY` | *(auto-built)* | Path to a pre-built `tr12-host` binary |
| `TR12_SDK_BINARY` | *(auto-built)* | Path to a pre-built `cdd-sdk` binary |

If not set, the harness builds both binaries from source using `go build`.

## How It Works

- Each test run gets a fresh SQLite database, fresh certs directory, and fresh SDK state
- Ephemeral ports are allocated at runtime to avoid conflicts
- The host service and SDK are started as child processes and cleaned up via `t.Cleanup`
- A dynamic `host_configuration/tr12-host.json` is written to a temp directory so the SDK can discover the host's ephemeral ports
- A minimal test JPEG is created at `/tmp/image_sdi.jpg` for thumbnail testing

## File Layout

```
test/integration/
├── go.mod                # Standalone module (no imports from client/ or host/)
├── go.sum
├── helpers_test.go       # Process management, HTTP helpers, polling utilities
├── lifecycle_test.go     # TestFullLifecycle — the 11-phase integration test
├── testdata/
│   └── registration.json # 1-channel encoder registration payload
└── README.md             # This file
```
