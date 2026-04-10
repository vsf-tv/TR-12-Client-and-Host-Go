
# DeviceStatus


## Properties

Name | Type
------------ | -------------
`status` | [Array&lt;StatusValue&gt;](StatusValue.md)
`channels` | [Array&lt;ChannelStatus&gt;](ChannelStatus.md)

## Example

```typescript
import type { DeviceStatus } from ''

// TODO: Update the object below with actual values
const example = {
  "status": null,
  "channels": null,
} satisfies DeviceStatus

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DeviceStatus
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


