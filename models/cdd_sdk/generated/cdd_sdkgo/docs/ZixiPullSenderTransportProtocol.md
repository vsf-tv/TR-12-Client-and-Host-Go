# ZixiPullSenderTransportProtocol

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MaximumLatencyMilliseconds** | Pointer to **float32** |  | [optional] [default to 1000]
**Encryption** | Pointer to [**ZixiEncryption**](ZixiEncryption.md) |  | [optional] 
**StreamId** | **string** |  | 
**Port** | **float32** |  | 
**Interface** | Pointer to **string** |  | [optional] 

## Methods

### NewZixiPullSenderTransportProtocol

`func NewZixiPullSenderTransportProtocol(streamId string, port float32, ) *ZixiPullSenderTransportProtocol`

NewZixiPullSenderTransportProtocol instantiates a new ZixiPullSenderTransportProtocol object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewZixiPullSenderTransportProtocolWithDefaults

`func NewZixiPullSenderTransportProtocolWithDefaults() *ZixiPullSenderTransportProtocol`

NewZixiPullSenderTransportProtocolWithDefaults instantiates a new ZixiPullSenderTransportProtocol object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMaximumLatencyMilliseconds

`func (o *ZixiPullSenderTransportProtocol) GetMaximumLatencyMilliseconds() float32`

GetMaximumLatencyMilliseconds returns the MaximumLatencyMilliseconds field if non-nil, zero value otherwise.

### GetMaximumLatencyMillisecondsOk

`func (o *ZixiPullSenderTransportProtocol) GetMaximumLatencyMillisecondsOk() (*float32, bool)`

GetMaximumLatencyMillisecondsOk returns a tuple with the MaximumLatencyMilliseconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaximumLatencyMilliseconds

`func (o *ZixiPullSenderTransportProtocol) SetMaximumLatencyMilliseconds(v float32)`

SetMaximumLatencyMilliseconds sets MaximumLatencyMilliseconds field to given value.

### HasMaximumLatencyMilliseconds

`func (o *ZixiPullSenderTransportProtocol) HasMaximumLatencyMilliseconds() bool`

HasMaximumLatencyMilliseconds returns a boolean if a field has been set.

### GetEncryption

`func (o *ZixiPullSenderTransportProtocol) GetEncryption() ZixiEncryption`

GetEncryption returns the Encryption field if non-nil, zero value otherwise.

### GetEncryptionOk

`func (o *ZixiPullSenderTransportProtocol) GetEncryptionOk() (*ZixiEncryption, bool)`

GetEncryptionOk returns a tuple with the Encryption field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryption

`func (o *ZixiPullSenderTransportProtocol) SetEncryption(v ZixiEncryption)`

SetEncryption sets Encryption field to given value.

### HasEncryption

`func (o *ZixiPullSenderTransportProtocol) HasEncryption() bool`

HasEncryption returns a boolean if a field has been set.

### GetStreamId

`func (o *ZixiPullSenderTransportProtocol) GetStreamId() string`

GetStreamId returns the StreamId field if non-nil, zero value otherwise.

### GetStreamIdOk

`func (o *ZixiPullSenderTransportProtocol) GetStreamIdOk() (*string, bool)`

GetStreamIdOk returns a tuple with the StreamId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamId

`func (o *ZixiPullSenderTransportProtocol) SetStreamId(v string)`

SetStreamId sets StreamId field to given value.


### GetPort

`func (o *ZixiPullSenderTransportProtocol) GetPort() float32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *ZixiPullSenderTransportProtocol) GetPortOk() (*float32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *ZixiPullSenderTransportProtocol) SetPort(v float32)`

SetPort sets Port field to given value.


### GetInterface

`func (o *ZixiPullSenderTransportProtocol) GetInterface() string`

GetInterface returns the Interface field if non-nil, zero value otherwise.

### GetInterfaceOk

`func (o *ZixiPullSenderTransportProtocol) GetInterfaceOk() (*string, bool)`

GetInterfaceOk returns a tuple with the Interface field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInterface

`func (o *ZixiPullSenderTransportProtocol) SetInterface(v string)`

SetInterface sets Interface field to given value.

### HasInterface

`func (o *ZixiPullSenderTransportProtocol) HasInterface() bool`

HasInterface returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


