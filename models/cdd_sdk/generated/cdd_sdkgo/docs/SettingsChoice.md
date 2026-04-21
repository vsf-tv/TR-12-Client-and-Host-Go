# SettingsChoice

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StandardSettings** | [**[]IdAndValue**](IdAndValue.md) |  | 
**Profile** | [**ChannelProfile**](ChannelProfile.md) |  | 

## Methods

### NewSettingsChoice

`func NewSettingsChoice(standardSettings []IdAndValue, profile ChannelProfile, ) *SettingsChoice`

NewSettingsChoice instantiates a new SettingsChoice object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSettingsChoiceWithDefaults

`func NewSettingsChoiceWithDefaults() *SettingsChoice`

NewSettingsChoiceWithDefaults instantiates a new SettingsChoice object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStandardSettings

`func (o *SettingsChoice) GetStandardSettings() []IdAndValue`

GetStandardSettings returns the StandardSettings field if non-nil, zero value otherwise.

### GetStandardSettingsOk

`func (o *SettingsChoice) GetStandardSettingsOk() (*[]IdAndValue, bool)`

GetStandardSettingsOk returns a tuple with the StandardSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStandardSettings

`func (o *SettingsChoice) SetStandardSettings(v []IdAndValue)`

SetStandardSettings sets StandardSettings field to given value.


### GetProfile

`func (o *SettingsChoice) GetProfile() ChannelProfile`

GetProfile returns the Profile field if non-nil, zero value otherwise.

### GetProfileOk

`func (o *SettingsChoice) GetProfileOk() (*ChannelProfile, bool)`

GetProfileOk returns a tuple with the Profile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProfile

`func (o *SettingsChoice) SetProfile(v ChannelProfile)`

SetProfile sets Profile field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


