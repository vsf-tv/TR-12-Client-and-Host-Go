
# WebRtcTransportProtocol


## Properties

Name | Type
------------ | -------------
`dtlsSetupRole` | [DtlsSetupRole](DtlsSetupRole.md)
`iceParameters` | [IceParameters](IceParameters.md)
`dtlsFingerprints` | [Array&lt;DtlsFingerprint&gt;](DtlsFingerprint.md)
`iceServers` | [Array&lt;IceServer&gt;](IceServer.md)
`fecConfig` | [WebRtcFecConfig](WebRtcFecConfig.md)
`simpleSettings` | [Array&lt;IdAndValue&gt;](IdAndValue.md)

## Example

```typescript
import type { WebRtcTransportProtocol } from ''

// TODO: Update the object below with actual values
const example = {
  "dtlsSetupRole": null,
  "iceParameters": null,
  "dtlsFingerprints": null,
  "iceServers": null,
  "fecConfig": null,
  "simpleSettings": null,
} satisfies WebRtcTransportProtocol

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as WebRtcTransportProtocol
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


