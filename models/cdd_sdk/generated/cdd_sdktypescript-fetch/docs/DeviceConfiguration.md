
# DeviceConfiguration


## Properties

Name | Type
------------ | -------------
`configurationId` | string
`channels` | [Array&lt;ChannelConfiguration&gt;](ChannelConfiguration.md)
`simpleSettings` | [Array&lt;IdAndValue&gt;](IdAndValue.md)
`health` | [Health](Health.md)

## Example

```typescript
import type { DeviceConfiguration } from ''

// TODO: Update the object below with actual values
const example = {
  "configurationId": null,
  "channels": null,
  "simpleSettings": null,
  "health": null,
} satisfies DeviceConfiguration

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DeviceConfiguration
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


