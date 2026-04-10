# GetConfigurationResponseContent

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | **bool** |  | 
**State** | **string** |  | 
**Message** | **string** |  | 
**Error** | Pointer to [**ErrorDetails**](ErrorDetails.md) |  | [optional] 
**Configuration** | Pointer to [**ConfigurationData**](ConfigurationData.md) |  | [optional] 

## Methods

### NewGetConfigurationResponseContent

`func NewGetConfigurationResponseContent(success bool, state string, message string, ) *GetConfigurationResponseContent`

NewGetConfigurationResponseContent instantiates a new GetConfigurationResponseContent object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGetConfigurationResponseContentWithDefaults

`func NewGetConfigurationResponseContentWithDefaults() *GetConfigurationResponseContent`

NewGetConfigurationResponseContentWithDefaults instantiates a new GetConfigurationResponseContent object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSuccess

`func (o *GetConfigurationResponseContent) GetSuccess() bool`

GetSuccess returns the Success field if non-nil, zero value otherwise.

### GetSuccessOk

`func (o *GetConfigurationResponseContent) GetSuccessOk() (*bool, bool)`

GetSuccessOk returns a tuple with the Success field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuccess

`func (o *GetConfigurationResponseContent) SetSuccess(v bool)`

SetSuccess sets Success field to given value.


### GetState

`func (o *GetConfigurationResponseContent) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *GetConfigurationResponseContent) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *GetConfigurationResponseContent) SetState(v string)`

SetState sets State field to given value.


### GetMessage

`func (o *GetConfigurationResponseContent) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *GetConfigurationResponseContent) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *GetConfigurationResponseContent) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetError

`func (o *GetConfigurationResponseContent) GetError() ErrorDetails`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *GetConfigurationResponseContent) GetErrorOk() (*ErrorDetails, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *GetConfigurationResponseContent) SetError(v ErrorDetails)`

SetError sets Error field to given value.

### HasError

`func (o *GetConfigurationResponseContent) HasError() bool`

HasError returns a boolean if a field has been set.

### GetConfiguration

`func (o *GetConfigurationResponseContent) GetConfiguration() ConfigurationData`

GetConfiguration returns the Configuration field if non-nil, zero value otherwise.

### GetConfigurationOk

`func (o *GetConfigurationResponseContent) GetConfigurationOk() (*ConfigurationData, bool)`

GetConfigurationOk returns a tuple with the Configuration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConfiguration

`func (o *GetConfigurationResponseContent) SetConfiguration(v ConfigurationData)`

SetConfiguration sets Configuration field to given value.

### HasConfiguration

`func (o *GetConfigurationResponseContent) HasConfiguration() bool`

HasConfiguration returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


