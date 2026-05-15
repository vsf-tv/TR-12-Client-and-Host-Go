
# SettingConstraint

A setting constraint is either an enumerated list of valid values or a numeric range — never both.

## Properties

Name | Type
------------ | -------------
`enums` | [EnumValues](EnumValues.md)
`ranges` | [RangeValues](RangeValues.md)

## Example

```typescript
import type { SettingConstraint } from ''

// TODO: Update the object below with actual values
const example = {
  "enums": null,
  "ranges": null,
} satisfies SettingConstraint

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as SettingConstraint
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


