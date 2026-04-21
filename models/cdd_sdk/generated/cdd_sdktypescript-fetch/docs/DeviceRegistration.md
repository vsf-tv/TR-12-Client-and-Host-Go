
# DeviceRegistration


## Properties

Name | Type
------------ | -------------
`channels` | [Array&lt;Channel&gt;](Channel.md)
`standardSettings` | [Array&lt;Setting&gt;](Setting.md)
`thumbnails` | [Array&lt;Thumbnail&gt;](Thumbnail.md)

## Example

```typescript
import type { DeviceRegistration } from ''

// TODO: Update the object below with actual values
const example = {
  "channels": null,
  "standardSettings": null,
  "thumbnails": null,
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


