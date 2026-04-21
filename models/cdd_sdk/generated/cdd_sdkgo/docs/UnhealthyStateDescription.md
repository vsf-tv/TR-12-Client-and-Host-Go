# UnhealthyStateDescription

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Messages** | **[]string** |  | 
**Timestamp** | **time.Time** |  | 
**ComponentName** | **string** |  | 

## Methods

### NewUnhealthyStateDescription

`func NewUnhealthyStateDescription(messages []string, timestamp time.Time, componentName string, ) *UnhealthyStateDescription`

NewUnhealthyStateDescription instantiates a new UnhealthyStateDescription object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewUnhealthyStateDescriptionWithDefaults

`func NewUnhealthyStateDescriptionWithDefaults() *UnhealthyStateDescription`

NewUnhealthyStateDescriptionWithDefaults instantiates a new UnhealthyStateDescription object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMessages

`func (o *UnhealthyStateDescription) GetMessages() []string`

GetMessages returns the Messages field if non-nil, zero value otherwise.

### GetMessagesOk

`func (o *UnhealthyStateDescription) GetMessagesOk() (*[]string, bool)`

GetMessagesOk returns a tuple with the Messages field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessages

`func (o *UnhealthyStateDescription) SetMessages(v []string)`

SetMessages sets Messages field to given value.


### GetTimestamp

`func (o *UnhealthyStateDescription) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *UnhealthyStateDescription) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *UnhealthyStateDescription) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.


### GetComponentName

`func (o *UnhealthyStateDescription) GetComponentName() string`

GetComponentName returns the ComponentName field if non-nil, zero value otherwise.

### GetComponentNameOk

`func (o *UnhealthyStateDescription) GetComponentNameOk() (*string, bool)`

GetComponentNameOk returns a tuple with the ComponentName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetComponentName

`func (o *UnhealthyStateDescription) SetComponentName(v string)`

SetComponentName sets ComponentName field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


