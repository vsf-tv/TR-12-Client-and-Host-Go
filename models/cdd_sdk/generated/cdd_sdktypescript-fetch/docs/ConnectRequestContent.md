
# ConnectRequestContent


## Properties

Name | Type
------------ | -------------
`registration` | [DeviceRegistration](DeviceRegistration.md)
`hostId` | string

## Example

```typescript
import type { ConnectRequestContent } from ''

// TODO: Update the object below with actual values
const example = {
  "registration": null,
  "hostId": null,
} satisfies ConnectRequestContent

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ConnectRequestContent
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


