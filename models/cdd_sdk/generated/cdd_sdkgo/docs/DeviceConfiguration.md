# DeviceConfiguration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ConfigurationId** | **string** |  | 
**Channels** | [**[]ChannelConfiguration**](ChannelConfiguration.md) |  | 
**SimpleSettings** | Pointer to [**[]IdAndValue**](IdAndValue.md) |  | [optional] 
**Health** | Pointer to [**Health**](Health.md) |  | [optional] 

## Methods

### NewDeviceConfiguration

`func NewDeviceConfiguration(configurationId string, channels []ChannelConfiguration, ) *DeviceConfiguration`

NewDeviceConfiguration instantiates a new DeviceConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeviceConfigurationWithDefaults

`func NewDeviceConfigurationWithDefaults() *DeviceConfiguration`

NewDeviceConfigurationWithDefaults instantiates a new DeviceConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetConfigurationId

`func (o *DeviceConfiguration) GetConfigurationId() string`

GetConfigurationId returns the ConfigurationId field if non-nil, zero value otherwise.

### GetConfigurationIdOk

`func (o *DeviceConfiguration) GetConfigurationIdOk() (*string, bool)`

GetConfigurationIdOk returns a tuple with the ConfigurationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfigurationId

`func (o *DeviceConfiguration) SetConfigurationId(v string)`

SetConfigurationId sets ConfigurationId field to given value.


### GetChannels

`func (o *DeviceConfiguration) GetChannels() []ChannelConfiguration`

GetChannels returns the Channels field if non-nil, zero value otherwise.

### GetChannelsOk

`func (o *DeviceConfiguration) GetChannelsOk() (*[]ChannelConfiguration, bool)`

GetChannelsOk returns a tuple with the Channels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannels

`func (o *DeviceConfiguration) SetChannels(v []ChannelConfiguration)`

SetChannels sets Channels field to given value.


### GetSimpleSettings

`func (o *DeviceConfiguration) GetSimpleSettings() []IdAndValue`

GetSimpleSettings returns the SimpleSettings field if non-nil, zero value otherwise.

### GetSimpleSettingsOk

`func (o *DeviceConfiguration) GetSimpleSettingsOk() (*[]IdAndValue, bool)`

GetSimpleSettingsOk returns a tuple with the SimpleSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSimpleSettings

`func (o *DeviceConfiguration) SetSimpleSettings(v []IdAndValue)`

SetSimpleSettings sets SimpleSettings field to given value.

### HasSimpleSettings

`func (o *DeviceConfiguration) HasSimpleSettings() bool`

HasSimpleSettings returns a boolean if a field has been set.

### GetHealth

`func (o *DeviceConfiguration) GetHealth() Health`

GetHealth returns the Health field if non-nil, zero value otherwise.

### GetHealthOk

`func (o *DeviceConfiguration) GetHealthOk() (*Health, bool)`

GetHealthOk returns a tuple with the Health field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHealth

`func (o *DeviceConfiguration) SetHealth(v Health)`

SetHealth sets Health field to given value.

### HasHealth

`func (o *DeviceConfiguration) HasHealth() bool`

HasHealth returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


