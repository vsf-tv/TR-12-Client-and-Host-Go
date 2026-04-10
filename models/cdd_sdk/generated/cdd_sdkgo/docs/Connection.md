# Connection

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TransportProtocol** | Pointer to [**TransportProtocol**](TransportProtocol.md) |  | [optional] 

## Methods

### NewConnection

`func NewConnection() *Connection`

NewConnection instantiates a new Connection object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConnectionWithDefaults

`func NewConnectionWithDefaults() *Connection`

NewConnectionWithDefaults instantiates a new Connection object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTransportProtocol

`func (o *Connection) GetTransportProtocol() TransportProtocol`

GetTransportProtocol returns the TransportProtocol field if non-nil, zero value otherwise.

### GetTransportProtocolOk

`func (o *Connection) GetTransportProtocolOk() (*TransportProtocol, bool)`

GetTransportProtocolOk returns a tuple with the TransportProtocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransportProtocol

`func (o *Connection) SetTransportProtocol(v TransportProtocol)`

SetTransportProtocol sets TransportProtocol field to given value.

### HasTransportProtocol

`func (o *Connection) HasTransportProtocol() bool`

HasTransportProtocol returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


