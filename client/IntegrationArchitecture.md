# TR-12 Client Integration Architecture

A device integrator building TR-12 support needs to implement three components on the device alongside the CDD SDK process: a **TR-12 Shim** that bridges the device's native control plane to TR-12 models, a **Native Device Control API** that the shim drives, and a **Web Console** that gives the device user visibility into both the TR-12 connection state and the underlying device settings.

## Integration Requirements

1. **User-facing controls** — The Web Console must let the device user enable/disable TR-12, select a host service, view the pairing code, and see the current connection state (DISCONNECTED, PAIRING, CONNECTED, RECONNECTING).

2. **Separation of concerns** — The Web Console handles user input and display only. The TR-12 Shim running on the device is responsible for driving the CDD SDK (connect, disconnect, report status, apply configuration) and translating between TR-12 models and the native device control API. The console never talks to the SDK directly.

3. **Persistent connections** — The shim must persist the user-selected TR-12 host ID across reboots and software updates. On power-up, the shim should automatically call `connect` with the previously configured `hostId` and `registration` payload. Users expect a paired and connected device to reconnect without manual intervention after a restart.

4. **Configuration refresh** — When the host service pushes a new desired configuration via MQTT, the shim applies it to the native device API and the console must refresh its display to reflect the updated state. The device user should always see the complete, current configuration without manually reloading.  The Web Console likely must call an API hosted by the Shim to see when an refresh is needed.

5. **Status reporting** — The shim must periodically read back actual device state from the native control API (bitrate, temperature, channel state, etc.) and publish it to the host service via `report_status`. The reporting interval is governed by the host settings received during authentication.

6. **Actual configuration reporting** — After applying a desired configuration, the shim must read back what the device actually applied and report it via `report_actual_configuration`. This lets the host service detect discrepancies between desired and actual state.

7. **Pairing flow visibility** — During initial pairing, the console must display the pairing code and expiration countdown so the user can claim the device on the host service. The shim polls `connect` until the device transitions from PAIRING to CONNECTED.  

8. **Deprovision support** — The console must provide a way for the user to deprovision the device from a host service, which deletes local credentials and informs the host. The shim calls `deprovision` on the SDK and resets to DISCONNECTED.

The diagram below shows how these pieces interact.

```
  ┌──────────────────────────────┐
  │                              │  Runs in browser.
  │  4. Web Console              │  User enables/disables TR-12, points to TR-12 Service
  │  (Device UI)                 │  Sees connected status,
  │                              │  pairing code, etc.
  └──────────────┬───────────────┘
                 │
                 │  Calls Shim:
                 |  start() / stop() /
                 │  should_console_refresh()
                 │
┌────────────────┼────────────────────────────────────────────────────────────────────────┐
│                │                    DEVICE (on host hardware)                           │
│                │                                                                        │
│                ▼                                                                        │
│  ┌───────────────────┐   connect / disconnect         ┌───────────────────────────┐     │
│  │                   │   get_state / deprovision      │                           │     │
│  │  2. TR-12 Shim    │   get_configuration     ──────▶│  3. CDD SDK               │     │
│  │                   │   report_status                │  (localhost REST API)     │     │
│  │                   │   report_actual_config         │                           │     │
│  └───┬───────────▲───┘                                │  State Machine · MQTT     │     │
│      │           │                                    │  Pairing · Certs          │     │
│      │ apply     │ read                               │                           │     │
│      │ config    │ back                               └─────────────┬─────────────┘     │
│      ▼           │                                                  │                   │
│  ┌───────────────┴───┐                                              │                   │
│  │                   │    ┌────────────┐                            │                   │
│  │  1. Native Device │    │ Persistence│                            │                   │
│  │  Control API      │    │            │                            │                   │
│  │                   │    │ tr-12 host │                            │                   │
│  └───────────────────┘    └────────────┘                            │                   │
│                                                                     │                   │
└─────────────────────────────────────────────────────────────────────┼───────────────────┘
                                                                      │ HTTPS / MQTT
                                                                      │ (port 443)
                                                                      ▼
                                                             ┌─────────────────┐
                                                             │  TR-12 Host     │
                                                             │  Service        │
                                                             └─────────────────┘
```

**Arrow summary**

| From | To | Operations |
|---|---|---|
| 4 → 2 | Web Console → TR-12 Shim | `start()`, `stop()`, `should_console_refresh()` |
| 2 → 3 | TR-12 Shim → CDD SDK | `connect`, `disconnect`, `get_state`, `deprovision`, `get_configuration`, `report_status`, `report_actual_configuration` |
| 2 → 1 | TR-12 Shim → Native API | Apply desired configuration (codec, resolution, channel state, transport, etc.) |
| 1 → 2 | Native API → TR-12 Shim | Read back actual device state for status and configuration reporting |
| 3 → Host | CDD SDK → Host Service | Pairing, authentication, MQTT pub/sub over TLS on port 443 |
