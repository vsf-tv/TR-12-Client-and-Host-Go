# DesiredDeviceConfiguration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Version** | **string** |  | 
**Channels** | [**[]DesiredChannelConfiguration**](DesiredChannelConfiguration.md) | See limits.smithy: MAX_CHANNELS | 
**StandardSettings** | Pointer to [**[]IdAndValue**](IdAndValue.md) | See limits.smithy: MAX_SETTINGS_PER_SCOPE | [optional] 

## Methods

### NewDesiredDeviceConfiguration

`func NewDesiredDeviceConfiguration(version string, channels []DesiredChannelConfiguration, ) *DesiredDeviceConfiguration`

NewDesiredDeviceConfiguration instantiates a new DesiredDeviceConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDesiredDeviceConfigurationWithDefaults

`func NewDesiredDeviceConfigurationWithDefaults() *DesiredDeviceConfiguration`

NewDesiredDeviceConfigurationWithDefaults instantiates a new DesiredDeviceConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetVersion

`func (o *DesiredDeviceConfiguration) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *DesiredDeviceConfiguration) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *DesiredDeviceConfiguration) SetVersion(v string)`

SetVersion sets Version field to given value.


### GetChannels

`func (o *DesiredDeviceConfiguration) GetChannels() []DesiredChannelConfiguration`

GetChannels returns the Channels field if non-nil, zero value otherwise.

### GetChannelsOk

`func (o *DesiredDeviceConfiguration) GetChannelsOk() (*[]DesiredChannelConfiguration, bool)`

GetChannelsOk returns a tuple with the Channels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannels

`func (o *DesiredDeviceConfiguration) SetChannels(v []DesiredChannelConfiguration)`

SetChannels sets Channels field to given value.


### GetStandardSettings

`func (o *DesiredDeviceConfiguration) GetStandardSettings() []IdAndValue`

GetStandardSettings returns the StandardSettings field if non-nil, zero value otherwise.

### GetStandardSettingsOk

`func (o *DesiredDeviceConfiguration) GetStandardSettingsOk() (*[]IdAndValue, bool)`

GetStandardSettingsOk returns a tuple with the StandardSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStandardSettings

`func (o *DesiredDeviceConfiguration) SetStandardSettings(v []IdAndValue)`

SetStandardSettings sets StandardSettings field to given value.

### HasStandardSettings

`func (o *DesiredDeviceConfiguration) HasStandardSettings() bool`

HasStandardSettings returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


