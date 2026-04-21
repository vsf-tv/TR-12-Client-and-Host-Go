# Channel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Id** | **string** |  | 
**ChannelType** | Pointer to [**ChannelType**](ChannelType.md) |  | [optional] 
**StandardSettings** | Pointer to [**[]Setting**](Setting.md) |  | [optional] 
**Profiles** | Pointer to [**[]ProfileDefinition**](ProfileDefinition.md) |  | [optional] 
**ConnectionProtocols** | Pointer to [**[]TransportProtocolName**](TransportProtocolName.md) |  | [optional] 

## Methods

### NewChannel

`func NewChannel(name string, id string, ) *Channel`

NewChannel instantiates a new Channel object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewChannelWithDefaults

`func NewChannelWithDefaults() *Channel`

NewChannelWithDefaults instantiates a new Channel object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *Channel) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Channel) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Channel) SetName(v string)`

SetName sets Name field to given value.


### GetId

`func (o *Channel) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Channel) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Channel) SetId(v string)`

SetId sets Id field to given value.


### GetChannelType

`func (o *Channel) GetChannelType() ChannelType`

GetChannelType returns the ChannelType field if non-nil, zero value otherwise.

### GetChannelTypeOk

`func (o *Channel) GetChannelTypeOk() (*ChannelType, bool)`

GetChannelTypeOk returns a tuple with the ChannelType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannelType

`func (o *Channel) SetChannelType(v ChannelType)`

SetChannelType sets ChannelType field to given value.

### HasChannelType

`func (o *Channel) HasChannelType() bool`

HasChannelType returns a boolean if a field has been set.

### GetStandardSettings

`func (o *Channel) GetStandardSettings() []Setting`

GetStandardSettings returns the StandardSettings field if non-nil, zero value otherwise.

### GetStandardSettingsOk

`func (o *Channel) GetStandardSettingsOk() (*[]Setting, bool)`

GetStandardSettingsOk returns a tuple with the StandardSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStandardSettings

`func (o *Channel) SetStandardSettings(v []Setting)`

SetStandardSettings sets StandardSettings field to given value.

### HasStandardSettings

`func (o *Channel) HasStandardSettings() bool`

HasStandardSettings returns a boolean if a field has been set.

### GetProfiles

`func (o *Channel) GetProfiles() []ProfileDefinition`

GetProfiles returns the Profiles field if non-nil, zero value otherwise.

### GetProfilesOk

`func (o *Channel) GetProfilesOk() (*[]ProfileDefinition, bool)`

GetProfilesOk returns a tuple with the Profiles field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProfiles

`func (o *Channel) SetProfiles(v []ProfileDefinition)`

SetProfiles sets Profiles field to given value.

### HasProfiles

`func (o *Channel) HasProfiles() bool`

HasProfiles returns a boolean if a field has been set.

### GetConnectionProtocols

`func (o *Channel) GetConnectionProtocols() []TransportProtocolName`

GetConnectionProtocols returns the ConnectionProtocols field if non-nil, zero value otherwise.

### GetConnectionProtocolsOk

`func (o *Channel) GetConnectionProtocolsOk() (*[]TransportProtocolName, bool)`

GetConnectionProtocolsOk returns a tuple with the ConnectionProtocols field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectionProtocols

`func (o *Channel) SetConnectionProtocols(v []TransportProtocolName)`

SetConnectionProtocols sets ConnectionProtocols field to given value.

### HasConnectionProtocols

`func (o *Channel) HasConnectionProtocols() bool`

HasConnectionProtocols returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


