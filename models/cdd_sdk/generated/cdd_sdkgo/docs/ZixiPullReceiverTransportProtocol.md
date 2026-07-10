# ZixiPullReceiverTransportProtocol

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MaximumLatencyMilliseconds** | Pointer to **float32** |  | [optional] [default to 1000]
**Encryption** | Pointer to [**ZixiEncryption**](ZixiEncryption.md) |  | [optional] 
**StreamId** | **string** |  | 
**Address** | **string** |  | 
**Port** | Pointer to **float32** |  | [optional] [default to 2088]

## Methods

### NewZixiPullReceiverTransportProtocol

`func NewZixiPullReceiverTransportProtocol(streamId string, address string, ) *ZixiPullReceiverTransportProtocol`

NewZixiPullReceiverTransportProtocol instantiates a new ZixiPullReceiverTransportProtocol object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewZixiPullReceiverTransportProtocolWithDefaults

`func NewZixiPullReceiverTransportProtocolWithDefaults() *ZixiPullReceiverTransportProtocol`

NewZixiPullReceiverTransportProtocolWithDefaults instantiates a new ZixiPullReceiverTransportProtocol object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMaximumLatencyMilliseconds

`func (o *ZixiPullReceiverTransportProtocol) GetMaximumLatencyMilliseconds() float32`

GetMaximumLatencyMilliseconds returns the MaximumLatencyMilliseconds field if non-nil, zero value otherwise.

### GetMaximumLatencyMillisecondsOk

`func (o *ZixiPullReceiverTransportProtocol) GetMaximumLatencyMillisecondsOk() (*float32, bool)`

GetMaximumLatencyMillisecondsOk returns a tuple with the MaximumLatencyMilliseconds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaximumLatencyMilliseconds

`func (o *ZixiPullReceiverTransportProtocol) SetMaximumLatencyMilliseconds(v float32)`

SetMaximumLatencyMilliseconds sets MaximumLatencyMilliseconds field to given value.

### HasMaximumLatencyMilliseconds

`func (o *ZixiPullReceiverTransportProtocol) HasMaximumLatencyMilliseconds() bool`

HasMaximumLatencyMilliseconds returns a boolean if a field has been set.

### GetEncryption

`func (o *ZixiPullReceiverTransportProtocol) GetEncryption() ZixiEncryption`

GetEncryption returns the Encryption field if non-nil, zero value otherwise.

### GetEncryptionOk

`func (o *ZixiPullReceiverTransportProtocol) GetEncryptionOk() (*ZixiEncryption, bool)`

GetEncryptionOk returns a tuple with the Encryption field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncryption

`func (o *ZixiPullReceiverTransportProtocol) SetEncryption(v ZixiEncryption)`

SetEncryption sets Encryption field to given value.

### HasEncryption

`func (o *ZixiPullReceiverTransportProtocol) HasEncryption() bool`

HasEncryption returns a boolean if a field has been set.

### GetStreamId

`func (o *ZixiPullReceiverTransportProtocol) GetStreamId() string`

GetStreamId returns the StreamId field if non-nil, zero value otherwise.

### GetStreamIdOk

`func (o *ZixiPullReceiverTransportProtocol) GetStreamIdOk() (*string, bool)`

GetStreamIdOk returns a tuple with the StreamId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamId

`func (o *ZixiPullReceiverTransportProtocol) SetStreamId(v string)`

SetStreamId sets StreamId field to given value.


### GetAddress

`func (o *ZixiPullReceiverTransportProtocol) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *ZixiPullReceiverTransportProtocol) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *ZixiPullReceiverTransportProtocol) SetAddress(v string)`

SetAddress sets Address field to given value.


### GetPort

`func (o *ZixiPullReceiverTransportProtocol) GetPort() float32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *ZixiPullReceiverTransportProtocol) GetPortOk() (*float32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *ZixiPullReceiverTransportProtocol) SetPort(v float32)`

SetPort sets Port field to given value.

### HasPort

`func (o *ZixiPullReceiverTransportProtocol) HasPort() bool`

HasPort returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


