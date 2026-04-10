
# ZixiListenerTransportProtocol


## Properties

Name | Type
------------ | -------------
`streamId` | string
`port` | number
`minimumLatencyMilliseconds` | number
`encryption` | [EncryptionAes](EncryptionAes.md)
`_interface` | string

## Example

```typescript
import type { ZixiListenerTransportProtocol } from ''

// TODO: Update the object below with actual values
const example = {
  "streamId": null,
  "port": null,
  "minimumLatencyMilliseconds": null,
  "encryption": null,
  "_interface": null,
} satisfies ZixiListenerTransportProtocol

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ZixiListenerTransportProtocol
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


