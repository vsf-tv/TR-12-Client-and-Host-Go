# SettingsChoice

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SimpleSettings** | [**[]IdAndValue**](IdAndValue.md) |  | 
**Profile** | [**SettingProfile**](SettingProfile.md) |  | 

## Methods

### NewSettingsChoice

`func NewSettingsChoice(simpleSettings []IdAndValue, profile SettingProfile, ) *SettingsChoice`

NewSettingsChoice instantiates a new SettingsChoice object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSettingsChoiceWithDefaults

`func NewSettingsChoiceWithDefaults() *SettingsChoice`

NewSettingsChoiceWithDefaults instantiates a new SettingsChoice object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSimpleSettings

`func (o *SettingsChoice) GetSimpleSettings() []IdAndValue`

GetSimpleSettings returns the SimpleSettings field if non-nil, zero value otherwise.

### GetSimpleSettingsOk

`func (o *SettingsChoice) GetSimpleSettingsOk() (*[]IdAndValue, bool)`

GetSimpleSettingsOk returns a tuple with the SimpleSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSimpleSettings

`func (o *SettingsChoice) SetSimpleSettings(v []IdAndValue)`

SetSimpleSettings sets SimpleSettings field to given value.


### GetProfile

`func (o *SettingsChoice) GetProfile() SettingProfile`

GetProfile returns the Profile field if non-nil, zero value otherwise.

### GetProfileOk

`func (o *SettingsChoice) GetProfileOk() (*SettingProfile, bool)`

GetProfileOk returns a tuple with the Profile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProfile

`func (o *SettingsChoice) SetProfile(v SettingProfile)`

SetProfile sets Profile field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


