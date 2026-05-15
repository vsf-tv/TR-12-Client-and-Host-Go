
# DesiredDeviceConfiguration

Desired device configuration — sent from host to device. Contains only fields the host controls. No device-reported fields.

## Properties

Name | Type
------------ | -------------
`version` | string
`channels` | [Array&lt;DesiredChannelConfiguration&gt;](DesiredChannelConfiguration.md)
`standardSettings` | [Array&lt;IdAndValue&gt;](IdAndValue.md)

## Example

```typescript
import type { DesiredDeviceConfiguration } from ''

// TODO: Update the object below with actual values
const example = {
  "version": null,
  "channels": null,
  "standardSettings": null,
} satisfies DesiredDeviceConfiguration

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DesiredDeviceConfiguration
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


