# Setting

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** |  | 
**Name** | **string** |  | 
**Info** | **string** |  | 
**Enums** | Pointer to [**EnumValues**](EnumValues.md) |  | [optional] 
**Ranges** | Pointer to [**RangeValues**](RangeValues.md) |  | [optional] 

## Methods

### NewSetting

`func NewSetting(id string, name string, info string, ) *Setting`

NewSetting instantiates a new Setting object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSettingWithDefaults

`func NewSettingWithDefaults() *Setting`

NewSettingWithDefaults instantiates a new Setting object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *Setting) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Setting) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Setting) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *Setting) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Setting) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Setting) SetName(v string)`

SetName sets Name field to given value.


### GetInfo

`func (o *Setting) GetInfo() string`

GetInfo returns the Info field if non-nil, zero value otherwise.

### GetInfoOk

`func (o *Setting) GetInfoOk() (*string, bool)`

GetInfoOk returns a tuple with the Info field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfo

`func (o *Setting) SetInfo(v string)`

SetInfo sets Info field to given value.


### GetEnums

`func (o *Setting) GetEnums() EnumValues`

GetEnums returns the Enums field if non-nil, zero value otherwise.

### GetEnumsOk

`func (o *Setting) GetEnumsOk() (*EnumValues, bool)`

GetEnumsOk returns a tuple with the Enums field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnums

`func (o *Setting) SetEnums(v EnumValues)`

SetEnums sets Enums field to given value.

### HasEnums

`func (o *Setting) HasEnums() bool`

HasEnums returns a boolean if a field has been set.

### GetRanges

`func (o *Setting) GetRanges() RangeValues`

GetRanges returns the Ranges field if non-nil, zero value otherwise.

### GetRangesOk

`func (o *Setting) GetRangesOk() (*RangeValues, bool)`

GetRangesOk returns a tuple with the Ranges field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRanges

`func (o *Setting) SetRanges(v RangeValues)`

SetRanges sets Ranges field to given value.

### HasRanges

`func (o *Setting) HasRanges() bool`

HasRanges returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


