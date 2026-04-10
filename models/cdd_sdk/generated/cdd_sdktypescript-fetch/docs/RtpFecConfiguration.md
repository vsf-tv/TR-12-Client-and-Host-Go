
# RtpFecConfiguration


## Properties

Name | Type
------------ | -------------
`columnFec` | [RtpFecStreamConfig](RtpFecStreamConfig.md)
`rowFec` | [RtpFecStreamConfig](RtpFecStreamConfig.md)
`matrixColumns` | number
`matrixRows` | number

## Example

```typescript
import type { RtpFecConfiguration } from ''

// TODO: Update the object below with actual values
const example = {
  "columnFec": null,
  "rowFec": null,
  "matrixColumns": null,
  "matrixRows": null,
} satisfies RtpFecConfiguration

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as RtpFecConfiguration
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


