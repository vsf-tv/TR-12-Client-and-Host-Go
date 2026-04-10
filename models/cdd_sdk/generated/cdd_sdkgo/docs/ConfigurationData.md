# ConfigurationData

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**UpdateId** | Pointer to **string** |  | [optional] 
**Payload** | Pointer to [**DeviceConfiguration**](DeviceConfiguration.md) |  | [optional] 

## Methods

### NewConfigurationData

`func NewConfigurationData() *ConfigurationData`

NewConfigurationData instantiates a new ConfigurationData object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigurationDataWithDefaults

`func NewConfigurationDataWithDefaults() *ConfigurationData`

NewConfigurationDataWithDefaults instantiates a new ConfigurationData object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUpdateId

`func (o *ConfigurationData) GetUpdateId() string`

GetUpdateId returns the UpdateId field if non-nil, zero value otherwise.

### GetUpdateIdOk

`func (o *ConfigurationData) GetUpdateIdOk() (*string, bool)`

GetUpdateIdOk returns a tuple with the UpdateId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpdateId

`func (o *ConfigurationData) SetUpdateId(v string)`

SetUpdateId sets UpdateId field to given value.

### HasUpdateId

`func (o *ConfigurationData) HasUpdateId() bool`

HasUpdateId returns a boolean if a field has been set.

### GetPayload

`func (o *ConfigurationData) GetPayload() DeviceConfiguration`

GetPayload returns the Payload field if non-nil, zero value otherwise.

### GetPayloadOk

`func (o *ConfigurationData) GetPayloadOk() (*DeviceConfiguration, bool)`

GetPayloadOk returns a tuple with the Payload field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPayload

`func (o *ConfigurationData) SetPayload(v DeviceConfiguration)`

SetPayload sets Payload field to given value.

### HasPayload

`func (o *ConfigurationData) HasPayload() bool`

HasPayload returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


