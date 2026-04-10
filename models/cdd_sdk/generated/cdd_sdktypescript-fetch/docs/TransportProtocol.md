
# TransportProtocol


## Properties

Name | Type
------------ | -------------
`srtListener` | [SrtListenerTransportProtocol](SrtListenerTransportProtocol.md)
`srtCaller` | [SrtCallerTransportProtocol](SrtCallerTransportProtocol.md)
`ristListener` | [RistListenerTransportProtocol](RistListenerTransportProtocol.md)
`ristCaller` | [RistCallerTransportProtocol](RistCallerTransportProtocol.md)
`zixiListener` | [ZixiListenerTransportProtocol](ZixiListenerTransportProtocol.md)
`zixiCaller` | [ZixiCallerTransportProtocol](ZixiCallerTransportProtocol.md)
`rtp` | [RtpTransportProtocol](RtpTransportProtocol.md)
`webRtc` | [WebRtcTransportProtocol](WebRtcTransportProtocol.md)

## Example

```typescript
import type { TransportProtocol } from ''

// TODO: Update the object below with actual values
const example = {
  "srtListener": null,
  "srtCaller": null,
  "ristListener": null,
  "ristCaller": null,
  "zixiListener": null,
  "zixiCaller": null,
  "rtp": null,
  "webRtc": null,
} satisfies TransportProtocol

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as TransportProtocol
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


