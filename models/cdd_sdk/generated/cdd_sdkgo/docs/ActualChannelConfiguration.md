# ActualChannelConfiguration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Version** | **string** |  | 
**State** | [**ChannelState**](ChannelState.md) |  | 
**ChannelSettings** | Pointer to [**ChannelSettings**](ChannelSettings.md) |  | [optional] 
**Protocol** | Pointer to [**TransportProtocol**](TransportProtocol.md) |  | [optional] 
**Health** | Pointer to [**Health**](Health.md) |  | [optional] 
**ThumbnailLocalPath** | Pointer to **string** | Informs the TR12 Client so it may service Thumbnail Subscriptions.  Host service can ignore. | [optional] 

## Methods

### NewActualChannelConfiguration

`func NewActualChannelConfiguration(id string, version string, state ChannelState, ) *ActualChannelConfiguration`

NewActualChannelConfiguration instantiates a new ActualChannelConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewActualChannelConfigurationWithDefaults

`func NewActualChannelConfigurationWithDefaults() *ActualChannelConfiguration`

NewActualChannelConfigurationWithDefaults instantiates a new ActualChannelConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ActualChannelConfiguration) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ActualChannelConfiguration) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ActualChannelConfiguration) SetId(v string)`

SetId sets Id field to given value.


### GetVersion

`func (o *ActualChannelConfiguration) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *ActualChannelConfiguration) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *ActualChannelConfiguration) SetVersion(v string)`

SetVersion sets Version field to given value.


### GetState

`func (o *ActualChannelConfiguration) GetState() ChannelState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *ActualChannelConfiguration) GetStateOk() (*ChannelState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *ActualChannelConfiguration) SetState(v ChannelState)`

SetState sets State field to given value.


### GetChannelSettings

`func (o *ActualChannelConfiguration) GetChannelSettings() ChannelSettings`

GetChannelSettings returns the ChannelSettings field if non-nil, zero value otherwise.

### GetChannelSettingsOk

`func (o *ActualChannelConfiguration) GetChannelSettingsOk() (*ChannelSettings, bool)`

GetChannelSettingsOk returns a tuple with the ChannelSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannelSettings

`func (o *ActualChannelConfiguration) SetChannelSettings(v ChannelSettings)`

SetChannelSettings sets ChannelSettings field to given value.

### HasChannelSettings

`func (o *ActualChannelConfiguration) HasChannelSettings() bool`

HasChannelSettings returns a boolean if a field has been set.

### GetProtocol

`func (o *ActualChannelConfiguration) GetProtocol() TransportProtocol`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *ActualChannelConfiguration) GetProtocolOk() (*TransportProtocol, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *ActualChannelConfiguration) SetProtocol(v TransportProtocol)`

SetProtocol sets Protocol field to given value.

### HasProtocol

`func (o *ActualChannelConfiguration) HasProtocol() bool`

HasProtocol returns a boolean if a field has been set.

### GetHealth

`func (o *ActualChannelConfiguration) GetHealth() Health`

GetHealth returns the Health field if non-nil, zero value otherwise.

### GetHealthOk

`func (o *ActualChannelConfiguration) GetHealthOk() (*Health, bool)`

GetHealthOk returns a tuple with the Health field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHealth

`func (o *ActualChannelConfiguration) SetHealth(v Health)`

SetHealth sets Health field to given value.

### HasHealth

`func (o *ActualChannelConfiguration) HasHealth() bool`

HasHealth returns a boolean if a field has been set.

### GetThumbnailLocalPath

`func (o *ActualChannelConfiguration) GetThumbnailLocalPath() string`

GetThumbnailLocalPath returns the ThumbnailLocalPath field if non-nil, zero value otherwise.

### GetThumbnailLocalPathOk

`func (o *ActualChannelConfiguration) GetThumbnailLocalPathOk() (*string, bool)`

GetThumbnailLocalPathOk returns a tuple with the ThumbnailLocalPath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThumbnailLocalPath

`func (o *ActualChannelConfiguration) SetThumbnailLocalPath(v string)`

SetThumbnailLocalPath sets ThumbnailLocalPath field to given value.

### HasThumbnailLocalPath

`func (o *ActualChannelConfiguration) HasThumbnailLocalPath() bool`

HasThumbnailLocalPath returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


