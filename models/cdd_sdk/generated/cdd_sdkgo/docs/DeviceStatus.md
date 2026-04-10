# DeviceStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Status** | [**[]StatusValue**](StatusValue.md) |  | 
**Channels** | Pointer to [**[]ChannelStatus**](ChannelStatus.md) |  | [optional] 

## Methods

### NewDeviceStatus

`func NewDeviceStatus(status []StatusValue, ) *DeviceStatus`

NewDeviceStatus instantiates a new DeviceStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeviceStatusWithDefaults

`func NewDeviceStatusWithDefaults() *DeviceStatus`

NewDeviceStatusWithDefaults instantiates a new DeviceStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStatus

`func (o *DeviceStatus) GetStatus() []StatusValue`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *DeviceStatus) GetStatusOk() (*[]StatusValue, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *DeviceStatus) SetStatus(v []StatusValue)`

SetStatus sets Status field to given value.


### GetChannels

`func (o *DeviceStatus) GetChannels() []ChannelStatus`

GetChannels returns the Channels field if non-nil, zero value otherwise.

### GetChannelsOk

`func (o *DeviceStatus) GetChannelsOk() (*[]ChannelStatus, bool)`

GetChannelsOk returns a tuple with the Channels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannels

`func (o *DeviceStatus) SetChannels(v []ChannelStatus)`

SetChannels sets Channels field to given value.

### HasChannels

`func (o *DeviceStatus) HasChannels() bool`

HasChannels returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


