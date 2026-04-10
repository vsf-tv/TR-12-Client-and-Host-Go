# ConnectRequestContent

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Registration** | [**DeviceRegistration**](DeviceRegistration.md) |  | 
**HostId** | **string** |  | 

## Methods

### NewConnectRequestContent

`func NewConnectRequestContent(registration DeviceRegistration, hostId string, ) *ConnectRequestContent`

NewConnectRequestContent instantiates a new ConnectRequestContent object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConnectRequestContentWithDefaults

`func NewConnectRequestContentWithDefaults() *ConnectRequestContent`

NewConnectRequestContentWithDefaults instantiates a new ConnectRequestContent object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRegistration

`func (o *ConnectRequestContent) GetRegistration() DeviceRegistration`

GetRegistration returns the Registration field if non-nil, zero value otherwise.

### GetRegistrationOk

`func (o *ConnectRequestContent) GetRegistrationOk() (*DeviceRegistration, bool)`

GetRegistrationOk returns a tuple with the Registration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistration

`func (o *ConnectRequestContent) SetRegistration(v DeviceRegistration)`

SetRegistration sets Registration field to given value.


### GetHostId

`func (o *ConnectRequestContent) GetHostId() string`

GetHostId returns the HostId field if non-nil, zero value otherwise.

### GetHostIdOk

`func (o *ConnectRequestContent) GetHostIdOk() (*string, bool)`

GetHostIdOk returns a tuple with the HostId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHostId

`func (o *ConnectRequestContent) SetHostId(v string)`

SetHostId sets HostId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


