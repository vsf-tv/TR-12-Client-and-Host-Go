# DesiredChannelConfiguration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Version** | **string** |  | 
**State** | [**ChannelState**](ChannelState.md) |  | 
**ChannelSettings** | Pointer to [**ChannelSettings**](ChannelSettings.md) |  | [optional] 
**Protocol** | Pointer to [**TransportProtocol**](TransportProtocol.md) |  | [optional] 

## Methods

### NewDesiredChannelConfiguration

`func NewDesiredChannelConfiguration(id string, version string, state ChannelState, ) *DesiredChannelConfiguration`

NewDesiredChannelConfiguration instantiates a new DesiredChannelConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDesiredChannelConfigurationWithDefaults

`func NewDesiredChannelConfigurationWithDefaults() *DesiredChannelConfiguration`

NewDesiredChannelConfigurationWithDefaults instantiates a new DesiredChannelConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *DesiredChannelConfiguration) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *DesiredChannelConfiguration) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *DesiredChannelConfiguration) SetId(v string)`

SetId sets Id field to given value.


### GetVersion

`func (o *DesiredChannelConfiguration) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *DesiredChannelConfiguration) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *DesiredChannelConfiguration) SetVersion(v string)`

SetVersion sets Version field to given value.


### GetState

`func (o *DesiredChannelConfiguration) GetState() ChannelState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *DesiredChannelConfiguration) GetStateOk() (*ChannelState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *DesiredChannelConfiguration) SetState(v ChannelState)`

SetState sets State field to given value.


### GetChannelSettings

`func (o *DesiredChannelConfiguration) GetChannelSettings() ChannelSettings`

GetChannelSettings returns the ChannelSettings field if non-nil, zero value otherwise.

### GetChannelSettingsOk

`func (o *DesiredChannelConfiguration) GetChannelSettingsOk() (*ChannelSettings, bool)`

GetChannelSettingsOk returns a tuple with the ChannelSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannelSettings

`func (o *DesiredChannelConfiguration) SetChannelSettings(v ChannelSettings)`

SetChannelSettings sets ChannelSettings field to given value.

### HasChannelSettings

`func (o *DesiredChannelConfiguration) HasChannelSettings() bool`

HasChannelSettings returns a boolean if a field has been set.

### GetProtocol

`func (o *DesiredChannelConfiguration) GetProtocol() TransportProtocol`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *DesiredChannelConfiguration) GetProtocolOk() (*TransportProtocol, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *DesiredChannelConfiguration) SetProtocol(v TransportProtocol)`

SetProtocol sets Protocol field to given value.

### HasProtocol

`func (o *DesiredChannelConfiguration) HasProtocol() bool`

HasProtocol returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


