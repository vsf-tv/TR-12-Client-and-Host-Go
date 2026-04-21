
# Channel


## Properties

Name | Type
------------ | -------------
`name` | string
`id` | string
`channelType` | [ChannelType](ChannelType.md)
`standardSettings` | [Array&lt;Setting&gt;](Setting.md)
`profiles` | [Array&lt;ProfileDefinition&gt;](ProfileDefinition.md)
`connectionProtocols` | [Array&lt;TransportProtocolName&gt;](TransportProtocolName.md)

## Example

```typescript
import type { Channel } from ''

// TODO: Update the object below with actual values
const example = {
  "name": null,
  "id": null,
  "channelType": null,
  "standardSettings": null,
  "profiles": null,
  "connectionProtocols": null,
} satisfies Channel

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as Channel
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


