# ReportActualConfigurationResponseContent

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | **bool** |  | 
**State** | **string** |  | 
**Message** | **string** |  | 
**Error** | Pointer to [**ErrorDetails**](ErrorDetails.md) |  | [optional] 

## Methods

### NewReportActualConfigurationResponseContent

`func NewReportActualConfigurationResponseContent(success bool, state string, message string, ) *ReportActualConfigurationResponseContent`

NewReportActualConfigurationResponseContent instantiates a new ReportActualConfigurationResponseContent object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReportActualConfigurationResponseContentWithDefaults

`func NewReportActualConfigurationResponseContentWithDefaults() *ReportActualConfigurationResponseContent`

NewReportActualConfigurationResponseContentWithDefaults instantiates a new ReportActualConfigurationResponseContent object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSuccess

`func (o *ReportActualConfigurationResponseContent) GetSuccess() bool`

GetSuccess returns the Success field if non-nil, zero value otherwise.

### GetSuccessOk

`func (o *ReportActualConfigurationResponseContent) GetSuccessOk() (*bool, bool)`

GetSuccessOk returns a tuple with the Success field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuccess

`func (o *ReportActualConfigurationResponseContent) SetSuccess(v bool)`

SetSuccess sets Success field to given value.


### GetState

`func (o *ReportActualConfigurationResponseContent) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *ReportActualConfigurationResponseContent) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *ReportActualConfigurationResponseContent) SetState(v string)`

SetState sets State field to given value.


### GetMessage

`func (o *ReportActualConfigurationResponseContent) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ReportActualConfigurationResponseContent) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ReportActualConfigurationResponseContent) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetError

`func (o *ReportActualConfigurationResponseContent) GetError() ErrorDetails`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *ReportActualConfigurationResponseContent) GetErrorOk() (*ErrorDetails, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *ReportActualConfigurationResponseContent) SetError(v ErrorDetails)`

SetError sets Error field to given value.

### HasError

`func (o *ReportActualConfigurationResponseContent) HasError() bool`

HasError returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


