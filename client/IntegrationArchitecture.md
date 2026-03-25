# TR-12 Client Integration Architecture

A device integrator building TR-12 support needs to implement three components on the device alongside the CDD SDK process: a **TR-12 Shim** that bridges the device's native control plane to TR-12 models, a **Native Device Control API** that the shim drives, and a **Web Console** that gives the device user visibility into both the TR-12 connection state and the underlying device settings. The diagram below shows how these pieces interact.

```
  ┌──────────────────────────────┐
  │                              │         Runs in browser.
  │  4. Web Console              │         User enables/disables TR-12,
  │  (Device UI)                 │         sees connected status,
  │                              │         pairing code, etc.
  └──────────┬───────────────────┘
             │
             │  start() / stop() /
             │  should_console_refresh()
             │
┌────────────┼─────────────────────────────────────────────────────────────┐
│            │             DEVICE (on host hardware)                       │
│            │                                                             │
│            ▼                                                             │
│  ┌───────────────────────┐  connect / disconnect /                      │
│  │                       │  get_state / deprovision /                   │
│  │  2. TR-12 Shim        │  get_configuration /                         │
│  │                       │  report_status /                             │
│  │                       │  report_actual_config                        │
│  │                       │──────▶┌───────────────────────────┐          │
│  │                       │       │                           │          │
│  └───┬───────────▲───────┘       │  3. CDD SDK               │          │
│      │           │               │  (localhost REST API)      │          │
│      │ apply     │ read          │                           │          │
│      │ config    │ back          │  State Machine · MQTT     │          │
│      ▼           │               │  Pairing · Certs          │          │
│  ┌───────────────┴───────┐       │                           │          │
│  │                       │       └─────────────┬─────────────┘          │
│  │  1. Native Device     │                     │                        │
│  │  Control API          │                     │                        │
│  │                       │                     │                        │
│  └───────────────────────┘                     │                        │
│                                                │                        │
└────────────────────────────────────────────────┼────────────────────────┘
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

---

## Key Integration Points

### 1. Persistent Connections

The Web Console must persist the user-selected TR-12 host across reboots and software updates. When a device powers on, the console should automatically call `connect` with the previously configured `hostId` and `registration` payload. Users expect a paired and connected device to reconnect without manual intervention after a restart.

### 2. Configuration Refresh

The console should refresh its display whenever the TR-12 Shim receives an updated desired configuration from the host service via `get_configuration`. This ensures the device user always sees the complete, current state of the underlying device — including any remote changes made through the TR-12 host — without needing to manually reload the page or poll for updates.
