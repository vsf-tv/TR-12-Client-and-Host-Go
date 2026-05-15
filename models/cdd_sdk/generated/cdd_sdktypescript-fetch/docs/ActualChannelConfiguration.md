
# ActualChannelConfiguration

Actual channel configuration — reported by device to host. Extends desired fields with device-only reporting fields.

## Properties

Name | Type
------------ | -------------
`id` | string
`version` | string
`state` | [ChannelState](ChannelState.md)
`channelSettings` | [ChannelSettings](ChannelSettings.md)
`protocol` | [TransportProtocol](TransportProtocol.md)
`health` | [Health](Health.md)
`thumbnailLocalPath` | string

## Example

```typescript
import type { ActualChannelConfiguration } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "version": null,
  "state": null,
  "channelSettings": null,
  "protocol": null,
  "health": null,
  "thumbnailLocalPath": null,
} satisfies ActualChannelConfiguration

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ActualChannelConfiguration
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


