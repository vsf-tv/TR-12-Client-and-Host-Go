
# DesiredChannelConfiguration

Desired channel configuration — sent from host to device. Contains only fields the host controls. No device-reported fields.

## Properties

Name | Type
------------ | -------------
`id` | string
`version` | string
`state` | [ChannelState](ChannelState.md)
`channelSettings` | [ChannelSettings](ChannelSettings.md)
`protocol` | [TransportProtocol](TransportProtocol.md)

## Example

```typescript
import type { DesiredChannelConfiguration } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "version": null,
  "state": null,
  "channelSettings": null,
  "protocol": null,
} satisfies DesiredChannelConfiguration

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DesiredChannelConfiguration
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


