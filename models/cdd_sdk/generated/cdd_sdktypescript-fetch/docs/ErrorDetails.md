
# ErrorDetails


## Properties

Name | Type
------------ | -------------
`type` | string
`message` | string
`details` | string

## Example

```typescript
import type { ErrorDetails } from ''

// TODO: Update the object below with actual values
const example = {
  "type": null,
  "message": null,
  "details": null,
} satisfies ErrorDetails

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ErrorDetails
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


