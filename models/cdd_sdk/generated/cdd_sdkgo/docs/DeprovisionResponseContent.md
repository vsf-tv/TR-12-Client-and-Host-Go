# DeprovisionResponseContent

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | **bool** |  | 
**State** | **string** |  | 
**Message** | **string** |  | 
**Error** | Pointer to [**ErrorDetails**](ErrorDetails.md) |  | [optional] 

## Methods

### NewDeprovisionResponseContent

`func NewDeprovisionResponseContent(success bool, state string, message string, ) *DeprovisionResponseContent`

NewDeprovisionResponseContent instantiates a new DeprovisionResponseContent object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDeprovisionResponseContentWithDefaults

`func NewDeprovisionResponseContentWithDefaults() *DeprovisionResponseContent`

NewDeprovisionResponseContentWithDefaults instantiates a new DeprovisionResponseContent object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSuccess

`func (o *DeprovisionResponseContent) GetSuccess() bool`

GetSuccess returns the Success field if non-nil, zero value otherwise.

### GetSuccessOk

`func (o *DeprovisionResponseContent) GetSuccessOk() (*bool, bool)`

GetSuccessOk returns a tuple with the Success field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuccess

`func (o *DeprovisionResponseContent) SetSuccess(v bool)`

SetSuccess sets Success field to given value.


### GetState

`func (o *DeprovisionResponseContent) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *DeprovisionResponseContent) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *DeprovisionResponseContent) SetState(v string)`

SetState sets State field to given value.


### GetMessage

`func (o *DeprovisionResponseContent) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *DeprovisionResponseContent) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *DeprovisionResponseContent) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetError

`func (o *DeprovisionResponseContent) GetError() ErrorDetails`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *DeprovisionResponseContent) GetErrorOk() (*ErrorDetails, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *DeprovisionResponseContent) SetError(v ErrorDetails)`

SetError sets Error field to given value.

### HasError

`func (o *DeprovisionResponseContent) HasError() bool`

HasError returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


