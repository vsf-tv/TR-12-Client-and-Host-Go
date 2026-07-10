
# RistEncryption

RIST encryption (Main Profile, DTLS-PSK mode). Passphrase constraints are implementation-defined.

## Properties

Name | Type
------------ | -------------
`passphrase` | string

## Example

```typescript
import type { RistEncryption } from ''

// TODO: Update the object below with actual values
const example = {
  "passphrase": null,
} satisfies RistEncryption

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as RistEncryption
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


