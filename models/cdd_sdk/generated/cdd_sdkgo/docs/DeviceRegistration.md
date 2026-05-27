# DeviceRegistration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ChannelTemplates** | [**[]ChannelTemplate**](ChannelTemplate.md) | Unique channel capability definitions. Most devices have a small number of distinct channel configurations shared across many channel IDs. See limits.smithy: MAX_CHANNEL_TEMPLATES | 
**ChannelAssignments** | [**[]ChannelAssignment**](ChannelAssignment.md) | Maps individual channel IDs to their template definition. Total channel count across the device — max 50. See limits.smithy: MAX_CHANNELS | 
**Settings** | Pointer to [**[]Setting**](Setting.md) | Device-level settings (not channel-specific). | [optional] 

## Methods

### NewDeviceRegistration

`func NewDeviceRegistration(channelTemplates []ChannelTemplate, channelAssignments []ChannelAssignment, ) *DeviceRegistration`

NewDeviceRegistration instantiates a new DeviceRegistration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeviceRegistrationWithDefaults

`func NewDeviceRegistrationWithDefaults() *DeviceRegistration`

NewDeviceRegistrationWithDefaults instantiates a new DeviceRegistration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChannelTemplates

`func (o *DeviceRegistration) GetChannelTemplates() []ChannelTemplate`

GetChannelTemplates returns the ChannelTemplates field if non-nil, zero value otherwise.

### GetChannelTemplatesOk

`func (o *DeviceRegistration) GetChannelTemplatesOk() (*[]ChannelTemplate, bool)`

GetChannelTemplatesOk returns a tuple with the ChannelTemplates field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannelTemplates

`func (o *DeviceRegistration) SetChannelTemplates(v []ChannelTemplate)`

SetChannelTemplates sets ChannelTemplates field to given value.


### GetChannelAssignments

`func (o *DeviceRegistration) GetChannelAssignments() []ChannelAssignment`

GetChannelAssignments returns the ChannelAssignments field if non-nil, zero value otherwise.

### GetChannelAssignmentsOk

`func (o *DeviceRegistration) GetChannelAssignmentsOk() (*[]ChannelAssignment, bool)`

GetChannelAssignmentsOk returns a tuple with the ChannelAssignments field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannelAssignments

`func (o *DeviceRegistration) SetChannelAssignments(v []ChannelAssignment)`

SetChannelAssignments sets ChannelAssignments field to given value.


### GetSettings

`func (o *DeviceRegistration) GetSettings() []Setting`

GetSettings returns the Settings field if non-nil, zero value otherwise.

### GetSettingsOk

`func (o *DeviceRegistration) GetSettingsOk() (*[]Setting, bool)`

GetSettingsOk returns a tuple with the Settings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSettings

`func (o *DeviceRegistration) SetSettings(v []Setting)`

SetSettings sets Settings field to given value.

### HasSettings

`func (o *DeviceRegistration) HasSettings() bool`

HasSettings returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


