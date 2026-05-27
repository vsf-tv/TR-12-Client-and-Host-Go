
# DeviceRegistration

Device registration payload — advertises device capabilities to the host.  Uses a template/assignment pattern to avoid redundant channel definitions: - channelTemplates: unique capability definitions (max 5) - channelAssignments: maps each channel ID to a template (max 50)  This keeps the registration payload under 90 kB for MQTT transport.

## Properties

Name | Type
------------ | -------------
`channelTemplates` | [Array&lt;ChannelTemplate&gt;](ChannelTemplate.md)
`channelAssignments` | [Array&lt;ChannelAssignment&gt;](ChannelAssignment.md)
`settings` | [Array&lt;Setting&gt;](Setting.md)

## Example

```typescript
import type { DeviceRegistration } from ''

// TODO: Update the object below with actual values
const example = {
  "channelTemplates": null,
  "channelAssignments": null,
  "settings": null,
} satisfies DeviceRegistration

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DeviceRegistration
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


