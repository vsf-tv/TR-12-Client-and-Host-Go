
# ZixiEncryption

Zixi encryption. Passphrase constraints are implementation-defined. Key length selects the AES key size.

## Properties

Name | Type
------------ | -------------
`passphrase` | string
`keyLength` | [ZixiEncryptionKeyLength](ZixiEncryptionKeyLength.md)

## Example

```typescript
import type { ZixiEncryption } from ''

// TODO: Update the object below with actual values
const example = {
  "passphrase": null,
  "keyLength": null,
} satisfies ZixiEncryption

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ZixiEncryption
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


