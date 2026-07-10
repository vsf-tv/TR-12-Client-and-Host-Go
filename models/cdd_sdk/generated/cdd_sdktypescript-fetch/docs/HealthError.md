
# HealthError

Maximum 128 characters. Messages exceeding this limit are truncated before transmission to prevent MQTT payload bloat on devices with many channels.

## Properties

Name | Type
------------ | -------------
`message` | string
`timestamp` | Date

## Example

```typescript
import type { HealthError } from ''

// TODO: Update the object below with actual values
const example = {
  "message": null,
  "timestamp": null,
} satisfies HealthError

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as HealthError
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


