# RistListenerTransportProtocol

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StreamId** | Pointer to [**RistStreamIdentifier**](RistStreamIdentifier.md) |  | [optional] 
**Port** | **float32** |  | 
**MinimumLatencyMilliseconds** | **float32** |  | [default to 3000]
**Encryption** | Pointer to [**EncryptionAes**](EncryptionAes.md) |  | [optional] 
**Interface** | Pointer to **string** |  | [optional] 

## Methods

### NewRistListenerTransportProtocol

`func NewRistListenerTransportProtocol(port float32, minimumLatencyMilliseconds float32, ) *RistListenerTransportProtocol`

NewRistListenerTransportProtocol instantiates a new RistListenerTransportProtocol object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRistListenerTransportProtocolWithDefaults

`func NewRistListenerTransportProtocolWithDefaults() *RistListenerTransportProtocol`

NewRistListenerTransportProtocolWithDefaults instantiates a new RistListenerTransportProtocol object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStreamId

`func (o *RistListenerTransportProtocol) GetStreamId() RistStreamIdentifier`

GetStreamId returns the StreamId field if non-nil, zero value otherwise.

### GetStreamIdOk

`func (o *RistListenerTransportProtocol) GetStreamIdOk() (*RistStreamIdentifier, bool)`

GetStreamIdOk returns a tuple with the StreamId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamId

`func (o *RistListenerTransportProtocol) SetStreamId(v RistStreamIdentifier)`

SetStreamId sets StreamId field to given value.

### HasStreamId

`func (o *RistListenerTransportProtocol) HasStreamId() bool`

HasStreamId returns a boolean if a field has been set.

### GetPort

`func (o *RistListenerTransportProtocol) GetPort() float32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *RistListenerTransportProtocol) GetPortOk() (*float32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *RistListenerTransportProtocol) SetPort(v float32)`

SetPort sets Port field to given value.


### GetMinimumLatencyMilliseconds

`func (o *RistListenerTransportProtocol) GetMinimumLatencyMilliseconds() float32`

GetMinimumLatencyMilliseconds returns the MinimumLatencyMilliseconds field if non-nil, zero value otherwise.

### GetMinimumLatencyMillisecondsOk

`func (o *RistListenerTransportProtocol) GetMinimumLatencyMillisecondsOk() (*float32, bool)`

GetMinimumLatencyMillisecondsOk returns a tuple with the MinimumLatencyMilliseconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMinimumLatencyMilliseconds

`func (o *RistListenerTransportProtocol) SetMinimumLatencyMilliseconds(v float32)`

SetMinimumLatencyMilliseconds sets MinimumLatencyMilliseconds field to given value.


### GetEncryption

`func (o *RistListenerTransportProtocol) GetEncryption() EncryptionAes`

GetEncryption returns the Encryption field if non-nil, zero value otherwise.

### GetEncryptionOk

`func (o *RistListenerTransportProtocol) GetEncryptionOk() (*EncryptionAes, bool)`

GetEncryptionOk returns a tuple with the Encryption field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryption

`func (o *RistListenerTransportProtocol) SetEncryption(v EncryptionAes)`

SetEncryption sets Encryption field to given value.

### HasEncryption

`func (o *RistListenerTransportProtocol) HasEncryption() bool`

HasEncryption returns a boolean if a field has been set.

### GetInterface

`func (o *RistListenerTransportProtocol) GetInterface() string`

GetInterface returns the Interface field if non-nil, zero value otherwise.

### GetInterfaceOk

`func (o *RistListenerTransportProtocol) GetInterfaceOk() (*string, bool)`

GetInterfaceOk returns a tuple with the Interface field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterface

`func (o *RistListenerTransportProtocol) SetInterface(v string)`

SetInterface sets Interface field to given value.

### HasInterface

`func (o *RistListenerTransportProtocol) HasInterface() bool`

HasInterface returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


