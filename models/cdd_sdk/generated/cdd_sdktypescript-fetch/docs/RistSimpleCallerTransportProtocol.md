
# RistSimpleCallerTransportProtocol


## Properties

Name | Type
------------ | -------------
`address` | string
`port` | number
`minimumLatencyMilliseconds` | number
`encryption` | [EncryptionAes](EncryptionAes.md)

## Example

```typescript
import type { RistSimpleCallerTransportProtocol } from ''

// TODO: Update the object below with actual values
const example = {
  "address": null,
  "port": null,
  "minimumLatencyMilliseconds": null,
  "encryption": null,
} satisfies RistSimpleCallerTransportProtocol

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as RistSimpleCallerTransportProtocol
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


