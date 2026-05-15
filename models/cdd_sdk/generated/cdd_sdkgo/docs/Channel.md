# Channel

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | **string** |  | 
**Id** | **string** |  | 
**ChannelType** | [**ChannelType**](ChannelType.md) |  | 
**Settings** | Pointer to [**[]Setting**](Setting.md) |  | [optional] 
**Profiles** | Pointer to [**[]ProfileDefinition**](ProfileDefinition.md) |  | [optional] 
**Protocols** | Pointer to [**[]TransportProtocolName**](TransportProtocolName.md) |  | [optional] 

## Methods

### NewChannel

`func NewChannel(name string, id string, channelType ChannelType, ) *Channel`

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


### GetSettings

`func (o *Channel) GetSettings() []Setting`

GetSettings returns the Settings field if non-nil, zero value otherwise.

### GetSettingsOk

`func (o *Channel) GetSettingsOk() (*[]Setting, bool)`

GetSettingsOk returns a tuple with the Settings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSettings

`func (o *Channel) SetSettings(v []Setting)`

SetSettings sets Settings field to given value.

### HasSettings

`func (o *Channel) HasSettings() bool`

HasSettings returns a boolean if a field has been set.

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

### GetProtocols

`func (o *Channel) GetProtocols() []TransportProtocolName`

GetProtocols returns the Protocols field if non-nil, zero value otherwise.

### GetProtocolsOk

`func (o *Channel) GetProtocolsOk() (*[]TransportProtocolName, bool)`

GetProtocolsOk returns a tuple with the Protocols field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocols

`func (o *Channel) SetProtocols(v []TransportProtocolName)`

SetProtocols sets Protocols field to given value.

### HasProtocols

`func (o *Channel) HasProtocols() bool`

HasProtocols returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


