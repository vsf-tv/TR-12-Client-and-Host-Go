# ZixiListenerTransportProtocol

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**StreamId** | **string** |  | 
**Port** | **float32** |  | 
**MinimumLatencyMilliseconds** | **float32** |  | [default to 3000]
**Encryption** | Pointer to [**EncryptionAes**](EncryptionAes.md) |  | [optional] 
**Interface** | Pointer to **string** |  | [optional] 

## Methods

### NewZixiListenerTransportProtocol

`func NewZixiListenerTransportProtocol(streamId string, port float32, minimumLatencyMilliseconds float32, ) *ZixiListenerTransportProtocol`

NewZixiListenerTransportProtocol instantiates a new ZixiListenerTransportProtocol object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewZixiListenerTransportProtocolWithDefaults

`func NewZixiListenerTransportProtocolWithDefaults() *ZixiListenerTransportProtocol`

NewZixiListenerTransportProtocolWithDefaults instantiates a new ZixiListenerTransportProtocol object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetStreamId

`func (o *ZixiListenerTransportProtocol) GetStreamId() string`

GetStreamId returns the StreamId field if non-nil, zero value otherwise.

### GetStreamIdOk

`func (o *ZixiListenerTransportProtocol) GetStreamIdOk() (*string, bool)`

GetStreamIdOk returns a tuple with the StreamId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamId

`func (o *ZixiListenerTransportProtocol) SetStreamId(v string)`

SetStreamId sets StreamId field to given value.


### GetPort

`func (o *ZixiListenerTransportProtocol) GetPort() float32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *ZixiListenerTransportProtocol) GetPortOk() (*float32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *ZixiListenerTransportProtocol) SetPort(v float32)`

SetPort sets Port field to given value.


### GetMinimumLatencyMilliseconds

`func (o *ZixiListenerTransportProtocol) GetMinimumLatencyMilliseconds() float32`

GetMinimumLatencyMilliseconds returns the MinimumLatencyMilliseconds field if non-nil, zero value otherwise.

### GetMinimumLatencyMillisecondsOk

`func (o *ZixiListenerTransportProtocol) GetMinimumLatencyMillisecondsOk() (*float32, bool)`

GetMinimumLatencyMillisecondsOk returns a tuple with the MinimumLatencyMilliseconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMinimumLatencyMilliseconds

`func (o *ZixiListenerTransportProtocol) SetMinimumLatencyMilliseconds(v float32)`

SetMinimumLatencyMilliseconds sets MinimumLatencyMilliseconds field to given value.


### GetEncryption

`func (o *ZixiListenerTransportProtocol) GetEncryption() EncryptionAes`

GetEncryption returns the Encryption field if non-nil, zero value otherwise.

### GetEncryptionOk

`func (o *ZixiListenerTransportProtocol) GetEncryptionOk() (*EncryptionAes, bool)`

GetEncryptionOk returns a tuple with the Encryption field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryption

`func (o *ZixiListenerTransportProtocol) SetEncryption(v EncryptionAes)`

SetEncryption sets Encryption field to given value.

### HasEncryption

`func (o *ZixiListenerTransportProtocol) HasEncryption() bool`

HasEncryption returns a boolean if a field has been set.

### GetInterface

`func (o *ZixiListenerTransportProtocol) GetInterface() string`

GetInterface returns the Interface field if non-nil, zero value otherwise.

### GetInterfaceOk

`func (o *ZixiListenerTransportProtocol) GetInterfaceOk() (*string, bool)`

GetInterfaceOk returns a tuple with the Interface field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterface

`func (o *ZixiListenerTransportProtocol) SetInterface(v string)`

SetInterface sets Interface field to given value.

### HasInterface

`func (o *ZixiListenerTransportProtocol) HasInterface() bool`

HasInterface returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


