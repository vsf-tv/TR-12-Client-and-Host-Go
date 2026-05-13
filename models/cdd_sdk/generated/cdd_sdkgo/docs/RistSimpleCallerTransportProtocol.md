# RistSimpleCallerTransportProtocol

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Address** | **string** |  | 
**Port** | **float32** |  | 
**MinimumLatencyMilliseconds** | Pointer to **float32** |  | [optional] [default to 1000]
**Encryption** | Pointer to [**EncryptionAes**](EncryptionAes.md) |  | [optional] 

## Methods

### NewRistSimpleCallerTransportProtocol

`func NewRistSimpleCallerTransportProtocol(address string, port float32, ) *RistSimpleCallerTransportProtocol`

NewRistSimpleCallerTransportProtocol instantiates a new RistSimpleCallerTransportProtocol object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRistSimpleCallerTransportProtocolWithDefaults

`func NewRistSimpleCallerTransportProtocolWithDefaults() *RistSimpleCallerTransportProtocol`

NewRistSimpleCallerTransportProtocolWithDefaults instantiates a new RistSimpleCallerTransportProtocol object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAddress

`func (o *RistSimpleCallerTransportProtocol) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *RistSimpleCallerTransportProtocol) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *RistSimpleCallerTransportProtocol) SetAddress(v string)`

SetAddress sets Address field to given value.


### GetPort

`func (o *RistSimpleCallerTransportProtocol) GetPort() float32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *RistSimpleCallerTransportProtocol) GetPortOk() (*float32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *RistSimpleCallerTransportProtocol) SetPort(v float32)`

SetPort sets Port field to given value.


### GetMinimumLatencyMilliseconds

`func (o *RistSimpleCallerTransportProtocol) GetMinimumLatencyMilliseconds() float32`

GetMinimumLatencyMilliseconds returns the MinimumLatencyMilliseconds field if non-nil, zero value otherwise.

### GetMinimumLatencyMillisecondsOk

`func (o *RistSimpleCallerTransportProtocol) GetMinimumLatencyMillisecondsOk() (*float32, bool)`

GetMinimumLatencyMillisecondsOk returns a tuple with the MinimumLatencyMilliseconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMinimumLatencyMilliseconds

`func (o *RistSimpleCallerTransportProtocol) SetMinimumLatencyMilliseconds(v float32)`

SetMinimumLatencyMilliseconds sets MinimumLatencyMilliseconds field to given value.

### HasMinimumLatencyMilliseconds

`func (o *RistSimpleCallerTransportProtocol) HasMinimumLatencyMilliseconds() bool`

HasMinimumLatencyMilliseconds returns a boolean if a field has been set.

### GetEncryption

`func (o *RistSimpleCallerTransportProtocol) GetEncryption() EncryptionAes`

GetEncryption returns the Encryption field if non-nil, zero value otherwise.

### GetEncryptionOk

`func (o *RistSimpleCallerTransportProtocol) GetEncryptionOk() (*EncryptionAes, bool)`

GetEncryptionOk returns a tuple with the Encryption field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryption

`func (o *RistSimpleCallerTransportProtocol) SetEncryption(v EncryptionAes)`

SetEncryption sets Encryption field to given value.

### HasEncryption

`func (o *RistSimpleCallerTransportProtocol) HasEncryption() bool`

HasEncryption returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


