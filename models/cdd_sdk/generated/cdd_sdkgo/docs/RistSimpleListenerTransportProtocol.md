# RistSimpleListenerTransportProtocol

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Port** | **float32** |  | 
**MinimumLatencyMilliseconds** | Pointer to **float32** |  | [optional] [default to 1000]
**Encryption** | Pointer to [**EncryptionAes**](EncryptionAes.md) |  | [optional] 
**Interface** | Pointer to **string** |  | [optional] 

## Methods

### NewRistSimpleListenerTransportProtocol

`func NewRistSimpleListenerTransportProtocol(port float32, ) *RistSimpleListenerTransportProtocol`

NewRistSimpleListenerTransportProtocol instantiates a new RistSimpleListenerTransportProtocol object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRistSimpleListenerTransportProtocolWithDefaults

`func NewRistSimpleListenerTransportProtocolWithDefaults() *RistSimpleListenerTransportProtocol`

NewRistSimpleListenerTransportProtocolWithDefaults instantiates a new RistSimpleListenerTransportProtocol object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPort

`func (o *RistSimpleListenerTransportProtocol) GetPort() float32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *RistSimpleListenerTransportProtocol) GetPortOk() (*float32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *RistSimpleListenerTransportProtocol) SetPort(v float32)`

SetPort sets Port field to given value.


### GetMinimumLatencyMilliseconds

`func (o *RistSimpleListenerTransportProtocol) GetMinimumLatencyMilliseconds() float32`

GetMinimumLatencyMilliseconds returns the MinimumLatencyMilliseconds field if non-nil, zero value otherwise.

### GetMinimumLatencyMillisecondsOk

`func (o *RistSimpleListenerTransportProtocol) GetMinimumLatencyMillisecondsOk() (*float32, bool)`

GetMinimumLatencyMillisecondsOk returns a tuple with the MinimumLatencyMilliseconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMinimumLatencyMilliseconds

`func (o *RistSimpleListenerTransportProtocol) SetMinimumLatencyMilliseconds(v float32)`

SetMinimumLatencyMilliseconds sets MinimumLatencyMilliseconds field to given value.

### HasMinimumLatencyMilliseconds

`func (o *RistSimpleListenerTransportProtocol) HasMinimumLatencyMilliseconds() bool`

HasMinimumLatencyMilliseconds returns a boolean if a field has been set.

### GetEncryption

`func (o *RistSimpleListenerTransportProtocol) GetEncryption() EncryptionAes`

GetEncryption returns the Encryption field if non-nil, zero value otherwise.

### GetEncryptionOk

`func (o *RistSimpleListenerTransportProtocol) GetEncryptionOk() (*EncryptionAes, bool)`

GetEncryptionOk returns a tuple with the Encryption field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryption

`func (o *RistSimpleListenerTransportProtocol) SetEncryption(v EncryptionAes)`

SetEncryption sets Encryption field to given value.

### HasEncryption

`func (o *RistSimpleListenerTransportProtocol) HasEncryption() bool`

HasEncryption returns a boolean if a field has been set.

### GetInterface

`func (o *RistSimpleListenerTransportProtocol) GetInterface() string`

GetInterface returns the Interface field if non-nil, zero value otherwise.

### GetInterfaceOk

`func (o *RistSimpleListenerTransportProtocol) GetInterfaceOk() (*string, bool)`

GetInterfaceOk returns a tuple with the Interface field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterface

`func (o *RistSimpleListenerTransportProtocol) SetInterface(v string)`

SetInterface sets Interface field to given value.

### HasInterface

`func (o *RistSimpleListenerTransportProtocol) HasInterface() bool`

HasInterface returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


