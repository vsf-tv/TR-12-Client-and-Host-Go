# ChannelAssignment

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ChannelId** | **string** | An identifier string: 1–12 alphanumeric characters (letters and digits only, no special characters). Used for channelId, templateId, setting id, profile id, and channel status id. See limits.smithy: MAX_ID_LENGTH | 
**Name** | **string** |  | 
**TemplateId** | **string** | An identifier string: 1–12 alphanumeric characters (letters and digits only, no special characters). Used for channelId, templateId, setting id, profile id, and channel status id. See limits.smithy: MAX_ID_LENGTH | 

## Methods

### NewChannelAssignment

`func NewChannelAssignment(channelId string, name string, templateId string, ) *ChannelAssignment`

NewChannelAssignment instantiates a new ChannelAssignment object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewChannelAssignmentWithDefaults

`func NewChannelAssignmentWithDefaults() *ChannelAssignment`

NewChannelAssignmentWithDefaults instantiates a new ChannelAssignment object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetChannelId

`func (o *ChannelAssignment) GetChannelId() string`

GetChannelId returns the ChannelId field if non-nil, zero value otherwise.

### GetChannelIdOk

`func (o *ChannelAssignment) GetChannelIdOk() (*string, bool)`

GetChannelIdOk returns a tuple with the ChannelId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChannelId

`func (o *ChannelAssignment) SetChannelId(v string)`

SetChannelId sets ChannelId field to given value.


### GetName

`func (o *ChannelAssignment) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ChannelAssignment) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ChannelAssignment) SetName(v string)`

SetName sets Name field to given value.


### GetTemplateId

`func (o *ChannelAssignment) GetTemplateId() string`

GetTemplateId returns the TemplateId field if non-nil, zero value otherwise.

### GetTemplateIdOk

`func (o *ChannelAssignment) GetTemplateIdOk() (*string, bool)`

GetTemplateIdOk returns a tuple with the TemplateId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTemplateId

`func (o *ChannelAssignment) SetTemplateId(v string)`

SetTemplateId sets TemplateId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


