# RegisterRequestContent

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Registration** | [**DeviceRegistration**](DeviceRegistration.md) |  | 

## Methods

### NewRegisterRequestContent

`func NewRegisterRequestContent(registration DeviceRegistration, ) *RegisterRequestContent`

NewRegisterRequestContent instantiates a new RegisterRequestContent object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRegisterRequestContentWithDefaults

`func NewRegisterRequestContentWithDefaults() *RegisterRequestContent`

NewRegisterRequestContentWithDefaults instantiates a new RegisterRequestContent object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRegistration

`func (o *RegisterRequestContent) GetRegistration() DeviceRegistration`

GetRegistration returns the Registration field if non-nil, zero value otherwise.

### GetRegistrationOk

`func (o *RegisterRequestContent) GetRegistrationOk() (*DeviceRegistration, bool)`

GetRegistrationOk returns a tuple with the Registration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistration

`func (o *RegisterRequestContent) SetRegistration(v DeviceRegistration)`

SetRegistration sets Registration field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


