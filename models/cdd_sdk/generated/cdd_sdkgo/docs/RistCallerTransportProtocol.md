# RistCallerTransportProtocol

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StreamId** | Pointer to [**RistStreamIdentifier**](RistStreamIdentifier.md) |  | [optional] 
**Address** | **string** |  | 
**Port** | **float32** |  | 
**MinimumLatencyMilliseconds** | Pointer to **float32** |  | [optional] [default to 1000]
**Encryption** | Pointer to [**EncryptionAes**](EncryptionAes.md) |  | [optional] 

## Methods

### NewRistCallerTransportProtocol

`func NewRistCallerTransportProtocol(address string, port float32, ) *RistCallerTransportProtocol`

NewRistCallerTransportProtocol instantiates a new RistCallerTransportProtocol object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRistCallerTransportProtocolWithDefaults

`func NewRistCallerTransportProtocolWithDefaults() *RistCallerTransportProtocol`

NewRistCallerTransportProtocolWithDefaults instantiates a new RistCallerTransportProtocol object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStreamId

`func (o *RistCallerTransportProtocol) GetStreamId() RistStreamIdentifier`

GetStreamId returns the StreamId field if non-nil, zero value otherwise.

### GetStreamIdOk

`func (o *RistCallerTransportProtocol) GetStreamIdOk() (*RistStreamIdentifier, bool)`

GetStreamIdOk returns a tuple with the StreamId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamId

`func (o *RistCallerTransportProtocol) SetStreamId(v RistStreamIdentifier)`

SetStreamId sets StreamId field to given value.

### HasStreamId

`func (o *RistCallerTransportProtocol) HasStreamId() bool`

HasStreamId returns a boolean if a field has been set.

### GetAddress

`func (o *RistCallerTransportProtocol) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *RistCallerTransportProtocol) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *RistCallerTransportProtocol) SetAddress(v string)`

SetAddress sets Address field to given value.


### GetPort

`func (o *RistCallerTransportProtocol) GetPort() float32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *RistCallerTransportProtocol) GetPortOk() (*float32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *RistCallerTransportProtocol) SetPort(v float32)`

SetPort sets Port field to given value.


### GetMinimumLatencyMilliseconds

`func (o *RistCallerTransportProtocol) GetMinimumLatencyMilliseconds() float32`

GetMinimumLatencyMilliseconds returns the MinimumLatencyMilliseconds field if non-nil, zero value otherwise.

### GetMinimumLatencyMillisecondsOk

`func (o *RistCallerTransportProtocol) GetMinimumLatencyMillisecondsOk() (*float32, bool)`

GetMinimumLatencyMillisecondsOk returns a tuple with the MinimumLatencyMilliseconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMinimumLatencyMilliseconds

`func (o *RistCallerTransportProtocol) SetMinimumLatencyMilliseconds(v float32)`

SetMinimumLatencyMilliseconds sets MinimumLatencyMilliseconds field to given value.

### HasMinimumLatencyMilliseconds

`func (o *RistCallerTransportProtocol) HasMinimumLatencyMilliseconds() bool`

HasMinimumLatencyMilliseconds returns a boolean if a field has been set.

### GetEncryption

`func (o *RistCallerTransportProtocol) GetEncryption() EncryptionAes`

GetEncryption returns the Encryption field if non-nil, zero value otherwise.

### GetEncryptionOk

`func (o *RistCallerTransportProtocol) GetEncryptionOk() (*EncryptionAes, bool)`

GetEncryptionOk returns a tuple with the Encryption field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryption

`func (o *RistCallerTransportProtocol) SetEncryption(v EncryptionAes)`

SetEncryption sets Encryption field to given value.

### HasEncryption

`func (o *RistCallerTransportProtocol) HasEncryption() bool`

HasEncryption returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


