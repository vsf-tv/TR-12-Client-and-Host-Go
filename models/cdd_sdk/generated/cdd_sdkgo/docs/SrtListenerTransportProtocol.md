# SrtListenerTransportProtocol

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StreamId** | Pointer to **string** |  | [optional] 
**Port** | **float32** |  | 
**MinimumLatencyMilliseconds** | Pointer to **float32** |  | [optional] [default to 1000]
**Encryption** | Pointer to [**EncryptionAes**](EncryptionAes.md) |  | [optional] 
**Interface** | Pointer to **string** |  | [optional] 

## Methods

### NewSrtListenerTransportProtocol

`func NewSrtListenerTransportProtocol(port float32, ) *SrtListenerTransportProtocol`

NewSrtListenerTransportProtocol instantiates a new SrtListenerTransportProtocol object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSrtListenerTransportProtocolWithDefaults

`func NewSrtListenerTransportProtocolWithDefaults() *SrtListenerTransportProtocol`

NewSrtListenerTransportProtocolWithDefaults instantiates a new SrtListenerTransportProtocol object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStreamId

`func (o *SrtListenerTransportProtocol) GetStreamId() string`

GetStreamId returns the StreamId field if non-nil, zero value otherwise.

### GetStreamIdOk

`func (o *SrtListenerTransportProtocol) GetStreamIdOk() (*string, bool)`

GetStreamIdOk returns a tuple with the StreamId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamId

`func (o *SrtListenerTransportProtocol) SetStreamId(v string)`

SetStreamId sets StreamId field to given value.

### HasStreamId

`func (o *SrtListenerTransportProtocol) HasStreamId() bool`

HasStreamId returns a boolean if a field has been set.

### GetPort

`func (o *SrtListenerTransportProtocol) GetPort() float32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *SrtListenerTransportProtocol) GetPortOk() (*float32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *SrtListenerTransportProtocol) SetPort(v float32)`

SetPort sets Port field to given value.


### GetMinimumLatencyMilliseconds

`func (o *SrtListenerTransportProtocol) GetMinimumLatencyMilliseconds() float32`

GetMinimumLatencyMilliseconds returns the MinimumLatencyMilliseconds field if non-nil, zero value otherwise.

### GetMinimumLatencyMillisecondsOk

`func (o *SrtListenerTransportProtocol) GetMinimumLatencyMillisecondsOk() (*float32, bool)`

GetMinimumLatencyMillisecondsOk returns a tuple with the MinimumLatencyMilliseconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMinimumLatencyMilliseconds

`func (o *SrtListenerTransportProtocol) SetMinimumLatencyMilliseconds(v float32)`

SetMinimumLatencyMilliseconds sets MinimumLatencyMilliseconds field to given value.

### HasMinimumLatencyMilliseconds

`func (o *SrtListenerTransportProtocol) HasMinimumLatencyMilliseconds() bool`

HasMinimumLatencyMilliseconds returns a boolean if a field has been set.

### GetEncryption

`func (o *SrtListenerTransportProtocol) GetEncryption() EncryptionAes`

GetEncryption returns the Encryption field if non-nil, zero value otherwise.

### GetEncryptionOk

`func (o *SrtListenerTransportProtocol) GetEncryptionOk() (*EncryptionAes, bool)`

GetEncryptionOk returns a tuple with the Encryption field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryption

`func (o *SrtListenerTransportProtocol) SetEncryption(v EncryptionAes)`

SetEncryption sets Encryption field to given value.

### HasEncryption

`func (o *SrtListenerTransportProtocol) HasEncryption() bool`

HasEncryption returns a boolean if a field has been set.

### GetInterface

`func (o *SrtListenerTransportProtocol) GetInterface() string`

GetInterface returns the Interface field if non-nil, zero value otherwise.

### GetInterfaceOk

`func (o *SrtListenerTransportProtocol) GetInterfaceOk() (*string, bool)`

GetInterfaceOk returns a tuple with the Interface field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterface

`func (o *SrtListenerTransportProtocol) SetInterface(v string)`

SetInterface sets Interface field to given value.

### HasInterface

`func (o *SrtListenerTransportProtocol) HasInterface() bool`

HasInterface returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


