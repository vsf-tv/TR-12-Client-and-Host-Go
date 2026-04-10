# IceServer

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Urls** | **[]string** |  | 
**Username** | Pointer to **string** |  | [optional] 
**Credential** | Pointer to **string** |  | [optional] 

## Methods

### NewIceServer

`func NewIceServer(urls []string, ) *IceServer`

NewIceServer instantiates a new IceServer object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIceServerWithDefaults

`func NewIceServerWithDefaults() *IceServer`

NewIceServerWithDefaults instantiates a new IceServer object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUrls

`func (o *IceServer) GetUrls() []string`

GetUrls returns the Urls field if non-nil, zero value otherwise.

### GetUrlsOk

`func (o *IceServer) GetUrlsOk() (*[]string, bool)`

GetUrlsOk returns a tuple with the Urls field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUrls

`func (o *IceServer) SetUrls(v []string)`

SetUrls sets Urls field to given value.


### GetUsername

`func (o *IceServer) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *IceServer) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *IceServer) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *IceServer) HasUsername() bool`

HasUsername returns a boolean if a field has been set.

### GetCredential

`func (o *IceServer) GetCredential() string`

GetCredential returns the Credential field if non-nil, zero value otherwise.

### GetCredentialOk

`func (o *IceServer) GetCredentialOk() (*string, bool)`

GetCredentialOk returns a tuple with the Credential field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCredential

`func (o *IceServer) SetCredential(v string)`

SetCredential sets Credential field to given value.

### HasCredential

`func (o *IceServer) HasCredential() bool`

HasCredential returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


