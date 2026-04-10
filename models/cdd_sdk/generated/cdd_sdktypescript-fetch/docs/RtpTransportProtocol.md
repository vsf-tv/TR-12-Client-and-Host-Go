
# RtpTransportProtocol


## Properties

Name | Type
------------ | -------------
`ip` | string
`port` | number
`sourceIpFilter` | string
`rtpPayloadType` | number
`fecConfig` | [RtpFecConfiguration](RtpFecConfiguration.md)

## Example

```typescript
import type { RtpTransportProtocol } from ''

// TODO: Update the object below with actual values
const example = {
  "ip": null,
  "port": null,
  "sourceIpFilter": null,
  "rtpPayloadType": null,
  "fecConfig": null,
} satisfies RtpTransportProtocol

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as RtpTransportProtocol
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


