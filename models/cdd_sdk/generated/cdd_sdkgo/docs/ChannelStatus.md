# ChannelStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | An identifier string: 1–12 alphanumeric characters (letters and digits only, no special characters). Used for channelId, templateId, setting id, profile id, and channel status id. See limits.smithy: MAX_ID_LENGTH | 
**State** | [**ChannelState**](ChannelState.md) |  | 
**Status** | [**[]StatusValue**](StatusValue.md) | See limits.smithy: MAX_STATUS_VALUES | 
**Health** | Pointer to [**Health**](Health.md) |  | [optional] 

## Methods

### NewChannelStatus

`func NewChannelStatus(id string, state ChannelState, status []StatusValue, ) *ChannelStatus`

NewChannelStatus instantiates a new ChannelStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewChannelStatusWithDefaults

`func NewChannelStatusWithDefaults() *ChannelStatus`

NewChannelStatusWithDefaults instantiates a new ChannelStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *ChannelStatus) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *ChannelStatus) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *ChannelStatus) SetId(v string)`

SetId sets Id field to given value.


### GetState

`func (o *ChannelStatus) GetState() ChannelState`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *ChannelStatus) GetStateOk() (*ChannelState, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *ChannelStatus) SetState(v ChannelState)`

SetState sets State field to given value.


### GetStatus

`func (o *ChannelStatus) GetStatus() []StatusValue`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ChannelStatus) GetStatusOk() (*[]StatusValue, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ChannelStatus) SetStatus(v []StatusValue)`

SetStatus sets Status field to given value.


### GetHealth

`func (o *ChannelStatus) GetHealth() Health`

GetHealth returns the Health field if non-nil, zero value otherwise.

### GetHealthOk

`func (o *ChannelStatus) GetHealthOk() (*Health, bool)`

GetHealthOk returns a tuple with the Health field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHealth

`func (o *ChannelStatus) SetHealth(v Health)`

SetHealth sets Health field to given value.

### HasHealth

`func (o *ChannelStatus) HasHealth() bool`

HasHealth returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


