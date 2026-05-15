# DeviceRegistration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Channels** | [**[]Channel**](Channel.md) |  | 
**DeviceRegistrationSettings** | Pointer to [**[]Setting**](Setting.md) |  | [optional] 

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


### GetDeviceRegistrationSettings

`func (o *DeviceRegistration) GetDeviceRegistrationSettings() []Setting`

GetDeviceRegistrationSettings returns the DeviceRegistrationSettings field if non-nil, zero value otherwise.

### GetDeviceRegistrationSettingsOk

`func (o *DeviceRegistration) GetDeviceRegistrationSettingsOk() (*[]Setting, bool)`

GetDeviceRegistrationSettingsOk returns a tuple with the DeviceRegistrationSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceRegistrationSettings

`func (o *DeviceRegistration) SetDeviceRegistrationSettings(v []Setting)`

SetDeviceRegistrationSettings sets DeviceRegistrationSettings field to given value.

### HasDeviceRegistrationSettings

`func (o *DeviceRegistration) HasDeviceRegistrationSettings() bool`

HasDeviceRegistrationSettings returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


