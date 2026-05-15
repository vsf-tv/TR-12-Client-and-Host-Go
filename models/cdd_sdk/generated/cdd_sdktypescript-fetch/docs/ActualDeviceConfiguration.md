
# ActualDeviceConfiguration

Actual device configuration — reported by device to host. Extends desired fields with device-only reporting fields.

## Properties

Name | Type
------------ | -------------
`version` | string
`channels` | [Array&lt;ActualChannelConfiguration&gt;](ActualChannelConfiguration.md)
`standardSettings` | [Array&lt;IdAndValue&gt;](IdAndValue.md)
`health` | [Health](Health.md)

## Example

```typescript
import type { ActualDeviceConfiguration } from ''

// TODO: Update the object below with actual values
const example = {
  "version": null,
  "channels": null,
  "standardSettings": null,
  "health": null,
} satisfies ActualDeviceConfiguration

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ActualDeviceConfiguration
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


