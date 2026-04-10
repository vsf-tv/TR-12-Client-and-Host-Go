# CriticalState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Messages** | **[]string** |  | 
**Timestamp** | **time.Time** |  | 
**ComponentName** | **string** |  | 

## Methods

### NewCriticalState

`func NewCriticalState(messages []string, timestamp time.Time, componentName string, ) *CriticalState`

NewCriticalState instantiates a new CriticalState object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCriticalStateWithDefaults

`func NewCriticalStateWithDefaults() *CriticalState`

NewCriticalStateWithDefaults instantiates a new CriticalState object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMessages

`func (o *CriticalState) GetMessages() []string`

GetMessages returns the Messages field if non-nil, zero value otherwise.

### GetMessagesOk

`func (o *CriticalState) GetMessagesOk() (*[]string, bool)`

GetMessagesOk returns a tuple with the Messages field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessages

`func (o *CriticalState) SetMessages(v []string)`

SetMessages sets Messages field to given value.


### GetTimestamp

`func (o *CriticalState) GetTimestamp() time.Time`

GetTimestamp returns the Timestamp field if non-nil, zero value otherwise.

### GetTimestampOk

`func (o *CriticalState) GetTimestampOk() (*time.Time, bool)`

GetTimestampOk returns a tuple with the Timestamp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimestamp

`func (o *CriticalState) SetTimestamp(v time.Time)`

SetTimestamp sets Timestamp field to given value.


### GetComponentName

`func (o *CriticalState) GetComponentName() string`

GetComponentName returns the ComponentName field if non-nil, zero value otherwise.

### GetComponentNameOk

`func (o *CriticalState) GetComponentNameOk() (*string, bool)`

GetComponentNameOk returns a tuple with the ComponentName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetComponentName

`func (o *CriticalState) SetComponentName(v string)`

SetComponentName sets ComponentName field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


