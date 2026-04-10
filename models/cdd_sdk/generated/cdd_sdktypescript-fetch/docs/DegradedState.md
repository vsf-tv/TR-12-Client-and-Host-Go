
# DegradedState

Degraded or critical state — all diagnostic fields required.

## Properties

Name | Type
------------ | -------------
`messages` | Array&lt;string&gt;
`timestamp` | Date
`componentName` | string

## Example

```typescript
import type { DegradedState } from ''

// TODO: Update the object below with actual values
const example = {
  "messages": null,
  "timestamp": null,
  "componentName": null,
} satisfies DegradedState

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as DegradedState
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


