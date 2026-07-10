
# ZixiPullReceiverTransportProtocol


## Properties

Name | Type
------------ | -------------
`maximumLatencyMilliseconds` | number
`encryption` | [ZixiEncryption](ZixiEncryption.md)
`streamId` | string
`address` | string
`port` | number

## Example

```typescript
import type { ZixiPullReceiverTransportProtocol } from ''

// TODO: Update the object below with actual values
const example = {
  "maximumLatencyMilliseconds": null,
  "encryption": null,
  "streamId": null,
  "address": null,
  "port": null,
} satisfies ZixiPullReceiverTransportProtocol

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ZixiPullReceiverTransportProtocol
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


