
# ChannelStatus


## Properties

Name | Type
------------ | -------------
`id` | string
`state` | [ChannelState](ChannelState.md)
`status` | [Array&lt;StatusValue&gt;](StatusValue.md)

## Example

```typescript
import type { ChannelStatus } from ''

// TODO: Update the object below with actual values
const example = {
  "id": null,
  "state": null,
  "status": null,
} satisfies ChannelStatus

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ChannelStatus
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


