# ChannelTemplate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**ChannelType** | [**ChannelType**](ChannelType.md) |  | 
**Settings** | Pointer to [**[]Setting**](Setting.md) | Settings this channel type supports. | [optional] 
**Profiles** | Pointer to [**[]ProfileDefinition**](ProfileDefinition.md) | Profiles this channel type supports (mutually exclusive with settings in config). | [optional] 
**Protocols** | Pointer to [**[]TransportProtocolName**](TransportProtocolName.md) | Transport protocols this channel type can use. | [optional] 

## Methods

### NewChannelTemplate

`func NewChannelTemplate(id string, channelType ChannelType, ) *ChannelTemplate`

NewChannelTemplate instantiates a new ChannelTemplate object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewChannelTemplateWithDefaults

`func NewChannelTemplateWithDefaults() *ChannelTemplate`

NewChannelTemplateWithDefaults instantiates a new ChannelTemplate object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ChannelTemplate) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ChannelTemplate) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ChannelTemplate) SetId(v string)`

SetId sets Id field to given value.


### GetChannelType

`func (o *ChannelTemplate) GetChannelType() ChannelType`

GetChannelType returns the ChannelType field if non-nil, zero value otherwise.

### GetChannelTypeOk

`func (o *ChannelTemplate) GetChannelTypeOk() (*ChannelType, bool)`

GetChannelTypeOk returns a tuple with the ChannelType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannelType

`func (o *ChannelTemplate) SetChannelType(v ChannelType)`

SetChannelType sets ChannelType field to given value.


### GetSettings

`func (o *ChannelTemplate) GetSettings() []Setting`

GetSettings returns the Settings field if non-nil, zero value otherwise.

### GetSettingsOk

`func (o *ChannelTemplate) GetSettingsOk() (*[]Setting, bool)`

GetSettingsOk returns a tuple with the Settings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSettings

`func (o *ChannelTemplate) SetSettings(v []Setting)`

SetSettings sets Settings field to given value.

### HasSettings

`func (o *ChannelTemplate) HasSettings() bool`

HasSettings returns a boolean if a field has been set.

### GetProfiles

`func (o *ChannelTemplate) GetProfiles() []ProfileDefinition`

GetProfiles returns the Profiles field if non-nil, zero value otherwise.

### GetProfilesOk

`func (o *ChannelTemplate) GetProfilesOk() (*[]ProfileDefinition, bool)`

GetProfilesOk returns a tuple with the Profiles field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProfiles

`func (o *ChannelTemplate) SetProfiles(v []ProfileDefinition)`

SetProfiles sets Profiles field to given value.

### HasProfiles

`func (o *ChannelTemplate) HasProfiles() bool`

HasProfiles returns a boolean if a field has been set.

### GetProtocols

`func (o *ChannelTemplate) GetProtocols() []TransportProtocolName`

GetProtocols returns the Protocols field if non-nil, zero value otherwise.

### GetProtocolsOk

`func (o *ChannelTemplate) GetProtocolsOk() (*[]TransportProtocolName, bool)`

GetProtocolsOk returns a tuple with the Protocols field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocols

`func (o *ChannelTemplate) SetProtocols(v []TransportProtocolName)`

SetProtocols sets Protocols field to given value.

### HasProtocols

`func (o *ChannelTemplate) HasProtocols() bool`

HasProtocols returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


