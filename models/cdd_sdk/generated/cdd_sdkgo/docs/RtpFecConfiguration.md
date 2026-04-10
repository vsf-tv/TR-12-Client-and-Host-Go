# RtpFecConfiguration

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ColumnFec** | [**RtpFecStreamConfig**](RtpFecStreamConfig.md) |  | 
**RowFec** | Pointer to [**RtpFecStreamConfig**](RtpFecStreamConfig.md) |  | [optional] 
**MatrixColumns** | Pointer to **float32** |  | [optional] 
**MatrixRows** | Pointer to **float32** |  | [optional] 

## Methods

### NewRtpFecConfiguration

`func NewRtpFecConfiguration(columnFec RtpFecStreamConfig, ) *RtpFecConfiguration`

NewRtpFecConfiguration instantiates a new RtpFecConfiguration object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRtpFecConfigurationWithDefaults

`func NewRtpFecConfigurationWithDefaults() *RtpFecConfiguration`

NewRtpFecConfigurationWithDefaults instantiates a new RtpFecConfiguration object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetColumnFec

`func (o *RtpFecConfiguration) GetColumnFec() RtpFecStreamConfig`

GetColumnFec returns the ColumnFec field if non-nil, zero value otherwise.

### GetColumnFecOk

`func (o *RtpFecConfiguration) GetColumnFecOk() (*RtpFecStreamConfig, bool)`

GetColumnFecOk returns a tuple with the ColumnFec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetColumnFec

`func (o *RtpFecConfiguration) SetColumnFec(v RtpFecStreamConfig)`

SetColumnFec sets ColumnFec field to given value.


### GetRowFec

`func (o *RtpFecConfiguration) GetRowFec() RtpFecStreamConfig`

GetRowFec returns the RowFec field if non-nil, zero value otherwise.

### GetRowFecOk

`func (o *RtpFecConfiguration) GetRowFecOk() (*RtpFecStreamConfig, bool)`

GetRowFecOk returns a tuple with the RowFec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowFec

`func (o *RtpFecConfiguration) SetRowFec(v RtpFecStreamConfig)`

SetRowFec sets RowFec field to given value.

### HasRowFec

`func (o *RtpFecConfiguration) HasRowFec() bool`

HasRowFec returns a boolean if a field has been set.

### GetMatrixColumns

`func (o *RtpFecConfiguration) GetMatrixColumns() float32`

GetMatrixColumns returns the MatrixColumns field if non-nil, zero value otherwise.

### GetMatrixColumnsOk

`func (o *RtpFecConfiguration) GetMatrixColumnsOk() (*float32, bool)`

GetMatrixColumnsOk returns a tuple with the MatrixColumns field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatrixColumns

`func (o *RtpFecConfiguration) SetMatrixColumns(v float32)`

SetMatrixColumns sets MatrixColumns field to given value.

### HasMatrixColumns

`func (o *RtpFecConfiguration) HasMatrixColumns() bool`

HasMatrixColumns returns a boolean if a field has been set.

### GetMatrixRows

`func (o *RtpFecConfiguration) GetMatrixRows() float32`

GetMatrixRows returns the MatrixRows field if non-nil, zero value otherwise.

### GetMatrixRowsOk

`func (o *RtpFecConfiguration) GetMatrixRowsOk() (*float32, bool)`

GetMatrixRowsOk returns a tuple with the MatrixRows field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatrixRows

`func (o *RtpFecConfiguration) SetMatrixRows(v float32)`

SetMatrixRows sets MatrixRows field to given value.

### HasMatrixRows

`func (o *RtpFecConfiguration) HasMatrixRows() bool`

HasMatrixRows returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


