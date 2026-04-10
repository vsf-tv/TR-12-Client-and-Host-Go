
# ReportActualConfigurationResponseContent


## Properties

Name | Type
------------ | -------------
`success` | boolean
`state` | string
`message` | string
`error` | [ErrorDetails](ErrorDetails.md)

## Example

```typescript
import type { ReportActualConfigurationResponseContent } from ''

// TODO: Update the object below with actual values
const example = {
  "success": null,
  "state": null,
  "message": null,
  "error": null,
} satisfies ReportActualConfigurationResponseContent

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ReportActualConfigurationResponseContent
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


