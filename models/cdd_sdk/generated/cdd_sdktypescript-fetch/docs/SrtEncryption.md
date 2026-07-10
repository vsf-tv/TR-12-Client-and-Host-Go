
# SrtEncryption

SRT encryption. Passphrase is 10-80 characters (protocol-enforced). Key length selects the AES key size for the derived Stream Encrypting Key. All SRT versions support all three key lengths.

## Properties

Name | Type
------------ | -------------
`passphrase` | string
`keyLength` | [SrtEncryptionKeyLength](SrtEncryptionKeyLength.md)

## Example

```typescript
import type { SrtEncryption } from ''

// TODO: Update the object below with actual values
const example = {
  "passphrase": null,
  "keyLength": null,
} satisfies SrtEncryption

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as SrtEncryption
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


