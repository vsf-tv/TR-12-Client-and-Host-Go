# ConnectResponseContent

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Success** | **bool** |  | 
**State** | **string** |  | 
**Message** | **string** |  | 
**Error** | Pointer to [**ErrorDetails**](ErrorDetails.md) |  | [optional] 
**DeviceId** | Pointer to **string** |  | [optional] 
**RegionName** | Pointer to **string** |  | [optional] 
**PairingCode** | Pointer to **string** |  | [optional] 
**ExpiresSeconds** | Pointer to **float32** |  | [optional] 

## Methods

### NewConnectResponseContent

`func NewConnectResponseContent(success bool, state string, message string, ) *ConnectResponseContent`

NewConnectResponseContent instantiates a new ConnectResponseContent object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConnectResponseContentWithDefaults

`func NewConnectResponseContentWithDefaults() *ConnectResponseContent`

NewConnectResponseContentWithDefaults instantiates a new ConnectResponseContent object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSuccess

`func (o *ConnectResponseContent) GetSuccess() bool`

GetSuccess returns the Success field if non-nil, zero value otherwise.

### GetSuccessOk

`func (o *ConnectResponseContent) GetSuccessOk() (*bool, bool)`

GetSuccessOk returns a tuple with the Success field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuccess

`func (o *ConnectResponseContent) SetSuccess(v bool)`

SetSuccess sets Success field to given value.


### GetState

`func (o *ConnectResponseContent) GetState() string`

GetState returns the State field if non-nil, zero value otherwise.

### GetStateOk

`func (o *ConnectResponseContent) GetStateOk() (*string, bool)`

GetStateOk returns a tuple with the State field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetState

`func (o *ConnectResponseContent) SetState(v string)`

SetState sets State field to given value.


### GetMessage

`func (o *ConnectResponseContent) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ConnectResponseContent) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ConnectResponseContent) SetMessage(v string)`

SetMessage sets Message field to given value.


### GetError

`func (o *ConnectResponseContent) GetError() ErrorDetails`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *ConnectResponseContent) GetErrorOk() (*ErrorDetails, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *ConnectResponseContent) SetError(v ErrorDetails)`

SetError sets Error field to given value.

### HasError

`func (o *ConnectResponseContent) HasError() bool`

HasError returns a boolean if a field has been set.

### GetDeviceId

`func (o *ConnectResponseContent) GetDeviceId() string`

GetDeviceId returns the DeviceId field if non-nil, zero value otherwise.

### GetDeviceIdOk

`func (o *ConnectResponseContent) GetDeviceIdOk() (*string, bool)`

GetDeviceIdOk returns a tuple with the DeviceId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceId

`func (o *ConnectResponseContent) SetDeviceId(v string)`

SetDeviceId sets DeviceId field to given value.

### HasDeviceId

`func (o *ConnectResponseContent) HasDeviceId() bool`

HasDeviceId returns a boolean if a field has been set.

### GetRegionName

`func (o *ConnectResponseContent) GetRegionName() string`

GetRegionName returns the RegionName field if non-nil, zero value otherwise.

### GetRegionNameOk

`func (o *ConnectResponseContent) GetRegionNameOk() (*string, bool)`

GetRegionNameOk returns a tuple with the RegionName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegionName

`func (o *ConnectResponseContent) SetRegionName(v string)`

SetRegionName sets RegionName field to given value.

### HasRegionName

`func (o *ConnectResponseContent) HasRegionName() bool`

HasRegionName returns a boolean if a field has been set.

### GetPairingCode

`func (o *ConnectResponseContent) GetPairingCode() string`

GetPairingCode returns the PairingCode field if non-nil, zero value otherwise.

### GetPairingCodeOk

`func (o *ConnectResponseContent) GetPairingCodeOk() (*string, bool)`

GetPairingCodeOk returns a tuple with the PairingCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPairingCode

`func (o *ConnectResponseContent) SetPairingCode(v string)`

SetPairingCode sets PairingCode field to given value.

### HasPairingCode

`func (o *ConnectResponseContent) HasPairingCode() bool`

HasPairingCode returns a boolean if a field has been set.

### GetExpiresSeconds

`func (o *ConnectResponseContent) GetExpiresSeconds() float32`

GetExpiresSeconds returns the ExpiresSeconds field if non-nil, zero value otherwise.

### GetExpiresSecondsOk

`func (o *ConnectResponseContent) GetExpiresSecondsOk() (*float32, bool)`

GetExpiresSecondsOk returns a tuple with the ExpiresSeconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpiresSeconds

`func (o *ConnectResponseContent) SetExpiresSeconds(v float32)`

SetExpiresSeconds sets ExpiresSeconds field to given value.

### HasExpiresSeconds

`func (o *ConnectResponseContent) HasExpiresSeconds() bool`

HasExpiresSeconds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


