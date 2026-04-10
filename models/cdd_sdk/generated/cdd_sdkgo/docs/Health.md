# Health

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Healthy** | **map[string]interface{}** | Healthy state — only the state indicator, no additional fields needed. | 
**Degraded** | [**DegradedState**](DegradedState.md) |  | 
**Critical** | [**CriticalState**](CriticalState.md) |  | 

## Methods

### NewHealth

`func NewHealth(healthy map[string]interface{}, degraded DegradedState, critical CriticalState, ) *Health`

NewHealth instantiates a new Health object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewHealthWithDefaults

`func NewHealthWithDefaults() *Health`

NewHealthWithDefaults instantiates a new Health object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHealthy

`func (o *Health) GetHealthy() map[string]interface{}`

GetHealthy returns the Healthy field if non-nil, zero value otherwise.

### GetHealthyOk

`func (o *Health) GetHealthyOk() (*map[string]interface{}, bool)`

GetHealthyOk returns a tuple with the Healthy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHealthy

`func (o *Health) SetHealthy(v map[string]interface{})`

SetHealthy sets Healthy field to given value.


### GetDegraded

`func (o *Health) GetDegraded() DegradedState`

GetDegraded returns the Degraded field if non-nil, zero value otherwise.

### GetDegradedOk

`func (o *Health) GetDegradedOk() (*DegradedState, bool)`

GetDegradedOk returns a tuple with the Degraded field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDegraded

`func (o *Health) SetDegraded(v DegradedState)`

SetDegraded sets Degraded field to given value.


### GetCritical

`func (o *Health) GetCritical() CriticalState`

GetCritical returns the Critical field if non-nil, zero value otherwise.

### GetCriticalOk

`func (o *Health) GetCriticalOk() (*CriticalState, bool)`

GetCriticalOk returns a tuple with the Critical field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCritical

`func (o *Health) SetCritical(v CriticalState)`

SetCritical sets Critical field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


