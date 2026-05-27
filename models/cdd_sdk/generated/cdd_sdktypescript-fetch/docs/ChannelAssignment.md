
# ChannelAssignment

Associates a channel ID with a template. The channel inherits all capabilities (settings, profiles, protocols) from the referenced template.

## Properties

Name | Type
------------ | -------------
`channelId` | string
`name` | string
`templateId` | string

## Example

```typescript
import type { ChannelAssignment } from ''

// TODO: Update the object below with actual values
const example = {
  "channelId": null,
  "name": null,
  "templateId": null,
} satisfies ChannelAssignment

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ChannelAssignment
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


