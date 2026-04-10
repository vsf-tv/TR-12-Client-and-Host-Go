
# SettingsChoice


## Properties

Name | Type
------------ | -------------
`simpleSettings` | [Array&lt;IdAndValue&gt;](IdAndValue.md)
`profile` | [SettingProfile](SettingProfile.md)

## Example

```typescript
import type { SettingsChoice } from ''

// TODO: Update the object below with actual values
const example = {
  "simpleSettings": null,
  "profile": null,
} satisfies SettingsChoice

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as SettingsChoice
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


