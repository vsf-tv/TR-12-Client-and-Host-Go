# IdAndValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | An identifier string: 1–12 alphanumeric characters (letters and digits only, no special characters). Used for channelId, templateId, setting id, profile id, and channel status id. See limits.smithy: MAX_ID_LENGTH | 
**Value** | **string** |  | 

## Methods

### NewIdAndValue

`func NewIdAndValue(id string, value string, ) *IdAndValue`

NewIdAndValue instantiates a new IdAndValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIdAndValueWithDefaults

`func NewIdAndValueWithDefaults() *IdAndValue`

NewIdAndValueWithDefaults instantiates a new IdAndValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *IdAndValue) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *IdAndValue) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *IdAndValue) SetId(v string)`

SetId sets Id field to given value.


### GetValue

`func (o *IdAndValue) GetValue() string`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *IdAndValue) GetValueOk() (*string, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *IdAndValue) SetValue(v string)`

SetValue sets Value field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


