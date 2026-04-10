
# ChannelConfiguration


## Properties

Name | Type
------------ | -------------
`id` | string
`configurationId` | string
`state` | [ChannelState](ChannelState.md)
`settings` | [SettingsChoice](SettingsChoice.md)
`connection` | [Connection](Connection.md)
`health` | [Health](Health.md)

## Example

```typescript
import type { ChannelConfiguration } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "configurationId": null,
  "state": null,
  "settings": null,
  "connection": null,
  "health": null,
} satisfies ChannelConfiguration

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ChannelConfiguration
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


