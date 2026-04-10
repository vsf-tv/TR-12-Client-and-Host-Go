# DeviceRegistration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Channels** | [**[]Channel**](Channel.md) |  | 
**SimpleSettings** | Pointer to [**[]Setting**](Setting.md) |  | [optional] 
**Thumbnails** | Pointer to [**[]Thumbnail**](Thumbnail.md) |  | [optional] 

## Methods

### NewDeviceRegistration

`func NewDeviceRegistration(channels []Channel, ) *DeviceRegistration`

NewDeviceRegistration instantiates a new DeviceRegistration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeviceRegistrationWithDefaults

`func NewDeviceRegistrationWithDefaults() *DeviceRegistration`

NewDeviceRegistrationWithDefaults instantiates a new DeviceRegistration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChannels

`func (o *DeviceRegistration) GetChannels() []Channel`

GetChannels returns the Channels field if non-nil, zero value otherwise.

### GetChannelsOk

`func (o *DeviceRegistration) GetChannelsOk() (*[]Channel, bool)`

GetChannelsOk returns a tuple with the Channels field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannels

`func (o *DeviceRegistration) SetChannels(v []Channel)`

SetChannels sets Channels field to given value.


### GetSimpleSettings

`func (o *DeviceRegistration) GetSimpleSettings() []Setting`

GetSimpleSettings returns the SimpleSettings field if non-nil, zero value otherwise.

### GetSimpleSettingsOk

`func (o *DeviceRegistration) GetSimpleSettingsOk() (*[]Setting, bool)`

GetSimpleSettingsOk returns a tuple with the SimpleSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSimpleSettings

`func (o *DeviceRegistration) SetSimpleSettings(v []Setting)`

SetSimpleSettings sets SimpleSettings field to given value.

### HasSimpleSettings

`func (o *DeviceRegistration) HasSimpleSettings() bool`

HasSimpleSettings returns a boolean if a field has been set.

### GetThumbnails

`func (o *DeviceRegistration) GetThumbnails() []Thumbnail`

GetThumbnails returns the Thumbnails field if non-nil, zero value otherwise.

### GetThumbnailsOk

`func (o *DeviceRegistration) GetThumbnailsOk() (*[]Thumbnail, bool)`

GetThumbnailsOk returns a tuple with the Thumbnails field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThumbnails

`func (o *DeviceRegistration) SetThumbnails(v []Thumbnail)`

SetThumbnails sets Thumbnails field to given value.

### HasThumbnails

`func (o *DeviceRegistration) HasThumbnails() bool`

HasThumbnails returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


