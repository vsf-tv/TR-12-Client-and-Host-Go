
# WebRtcFecConfig


## Properties

Name | Type
------------ | -------------
`fecMechanism` | [WebRtcFecMechanism](WebRtcFecMechanism.md)
`redPayloadType` | number
`ulpfecPayloadType` | number
`targetOverheadPercentage` | number

## Example

```typescript
import type { WebRtcFecConfig } from ''

// TODO: Update the object below with actual values
const example = {
  "fecMechanism": null,
  "redPayloadType": null,
  "ulpfecPayloadType": null,
  "targetOverheadPercentage": null,
} satisfies WebRtcFecConfig

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as WebRtcFecConfig
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


