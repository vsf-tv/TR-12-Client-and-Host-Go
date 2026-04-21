# ZixiPushTransportProtocol

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StreamId** | Pointer to **string** |  | [optional] 
**Address** | **string** |  | 
**Port** | Pointer to **float32** |  | [optional] [default to 2088]
**MinimumLatencyMilliseconds** | Pointer to **float32** |  | [optional] [default to 1000]
**Encryption** | Pointer to [**EncryptionAes**](EncryptionAes.md) |  | [optional] 

## Methods

### NewZixiPushTransportProtocol

`func NewZixiPushTransportProtocol(address string, ) *ZixiPushTransportProtocol`

NewZixiPushTransportProtocol instantiates a new ZixiPushTransportProtocol object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewZixiPushTransportProtocolWithDefaults

`func NewZixiPushTransportProtocolWithDefaults() *ZixiPushTransportProtocol`

NewZixiPushTransportProtocolWithDefaults instantiates a new ZixiPushTransportProtocol object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStreamId

`func (o *ZixiPushTransportProtocol) GetStreamId() string`

GetStreamId returns the StreamId field if non-nil, zero value otherwise.

### GetStreamIdOk

`func (o *ZixiPushTransportProtocol) GetStreamIdOk() (*string, bool)`

GetStreamIdOk returns a tuple with the StreamId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamId

`func (o *ZixiPushTransportProtocol) SetStreamId(v string)`

SetStreamId sets StreamId field to given value.

### HasStreamId

`func (o *ZixiPushTransportProtocol) HasStreamId() bool`

HasStreamId returns a boolean if a field has been set.

### GetAddress

`func (o *ZixiPushTransportProtocol) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *ZixiPushTransportProtocol) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *ZixiPushTransportProtocol) SetAddress(v string)`

SetAddress sets Address field to given value.


### GetPort

`func (o *ZixiPushTransportProtocol) GetPort() float32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *ZixiPushTransportProtocol) GetPortOk() (*float32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *ZixiPushTransportProtocol) SetPort(v float32)`

SetPort sets Port field to given value.

### HasPort

`func (o *ZixiPushTransportProtocol) HasPort() bool`

HasPort returns a boolean if a field has been set.

### GetMinimumLatencyMilliseconds

`func (o *ZixiPushTransportProtocol) GetMinimumLatencyMilliseconds() float32`

GetMinimumLatencyMilliseconds returns the MinimumLatencyMilliseconds field if non-nil, zero value otherwise.

### GetMinimumLatencyMillisecondsOk

`func (o *ZixiPushTransportProtocol) GetMinimumLatencyMillisecondsOk() (*float32, bool)`

GetMinimumLatencyMillisecondsOk returns a tuple with the MinimumLatencyMilliseconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMinimumLatencyMilliseconds

`func (o *ZixiPushTransportProtocol) SetMinimumLatencyMilliseconds(v float32)`

SetMinimumLatencyMilliseconds sets MinimumLatencyMilliseconds field to given value.

### HasMinimumLatencyMilliseconds

`func (o *ZixiPushTransportProtocol) HasMinimumLatencyMilliseconds() bool`

HasMinimumLatencyMilliseconds returns a boolean if a field has been set.

### GetEncryption

`func (o *ZixiPushTransportProtocol) GetEncryption() EncryptionAes`

GetEncryption returns the Encryption field if non-nil, zero value otherwise.

### GetEncryptionOk

`func (o *ZixiPushTransportProtocol) GetEncryptionOk() (*EncryptionAes, bool)`

GetEncryptionOk returns a tuple with the Encryption field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryption

`func (o *ZixiPushTransportProtocol) SetEncryption(v EncryptionAes)`

SetEncryption sets Encryption field to given value.

### HasEncryption

`func (o *ZixiPushTransportProtocol) HasEncryption() bool`

HasEncryption returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


