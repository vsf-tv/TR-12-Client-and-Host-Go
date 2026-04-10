# ChannelConfiguration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**ConfigurationId** | **string** |  | 
**State** | [**ChannelState**](ChannelState.md) |  | 
**Settings** | Pointer to [**SettingsChoice**](SettingsChoice.md) |  | [optional] 
**Connection** | Pointer to [**Connection**](Connection.md) |  | [optional] 
**Health** | Pointer to [**Health**](Health.md) |  | [optional] 

## Methods

### NewChannelConfiguration

`func NewChannelConfiguration(id string, configurationId string, state ChannelState, ) *ChannelConfiguration`

NewChannelConfiguration instantiates a new ChannelConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewChannelConfigurationWithDefaults

`func NewChannelConfigurationWithDefaults() *ChannelConfiguration`

NewChannelConfigurationWithDefaults instantiates a new ChannelConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ChannelConfiguration) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ChannelConfiguration) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ChannelConfiguration) SetId(v string)`

SetId sets Id field to given value.


### GetConfigurationId

`func (o *ChannelConfiguration) GetConfigurationId() string`

GetConfigurationId returns the ConfigurationId field if non-nil, zero value otherwise.

### GetConfigurationIdOk

`func (o *ChannelConfiguration) GetConfigurationIdOk() (*string, bool)`

GetConfigurationIdOk returns a tuple with the ConfigurationId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfigurationId

`func (o *ChannelConfiguration) SetConfigurationId(v string)`

SetConfigurationId sets ConfigurationId field to given value.


### GetState

`func (o *ChannelConfiguration) GetState() ChannelState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *ChannelConfiguration) GetStateOk() (*ChannelState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *ChannelConfiguration) SetState(v ChannelState)`

SetState sets State field to given value.


### GetSettings

`func (o *ChannelConfiguration) GetSettings() SettingsChoice`

GetSettings returns the Settings field if non-nil, zero value otherwise.

### GetSettingsOk

`func (o *ChannelConfiguration) GetSettingsOk() (*SettingsChoice, bool)`

GetSettingsOk returns a tuple with the Settings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSettings

`func (o *ChannelConfiguration) SetSettings(v SettingsChoice)`

SetSettings sets Settings field to given value.

### HasSettings

`func (o *ChannelConfiguration) HasSettings() bool`

HasSettings returns a boolean if a field has been set.

### GetConnection

`func (o *ChannelConfiguration) GetConnection() Connection`

GetConnection returns the Connection field if non-nil, zero value otherwise.

### GetConnectionOk

`func (o *ChannelConfiguration) GetConnectionOk() (*Connection, bool)`

GetConnectionOk returns a tuple with the Connection field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnection

`func (o *ChannelConfiguration) SetConnection(v Connection)`

SetConnection sets Connection field to given value.

### HasConnection

`func (o *ChannelConfiguration) HasConnection() bool`

HasConnection returns a boolean if a field has been set.

### GetHealth

`func (o *ChannelConfiguration) GetHealth() Health`

GetHealth returns the Health field if non-nil, zero value otherwise.

### GetHealthOk

`func (o *ChannelConfiguration) GetHealthOk() (*Health, bool)`

GetHealthOk returns a tuple with the Health field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHealth

`func (o *ChannelConfiguration) SetHealth(v Health)`

SetHealth sets Health field to given value.

### HasHealth

`func (o *ChannelConfiguration) HasHealth() bool`

HasHealth returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


