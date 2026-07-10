# SrtEncryption

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Passphrase** | **string** |  | 
**KeyLength** | Pointer to [**SrtEncryptionKeyLength**](SrtEncryptionKeyLength.md) |  | [optional] 

## Methods

### NewSrtEncryption

`func NewSrtEncryption(passphrase string, ) *SrtEncryption`

NewSrtEncryption instantiates a new SrtEncryption object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSrtEncryptionWithDefaults

`func NewSrtEncryptionWithDefaults() *SrtEncryption`

NewSrtEncryptionWithDefaults instantiates a new SrtEncryption object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPassphrase

`func (o *SrtEncryption) GetPassphrase() string`

GetPassphrase returns the Passphrase field if non-nil, zero value otherwise.

### GetPassphraseOk

`func (o *SrtEncryption) GetPassphraseOk() (*string, bool)`

GetPassphraseOk returns a tuple with the Passphrase field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassphrase

`func (o *SrtEncryption) SetPassphrase(v string)`

SetPassphrase sets Passphrase field to given value.


### GetKeyLength

`func (o *SrtEncryption) GetKeyLength() SrtEncryptionKeyLength`

GetKeyLength returns the KeyLength field if non-nil, zero value otherwise.

### GetKeyLengthOk

`func (o *SrtEncryption) GetKeyLengthOk() (*SrtEncryptionKeyLength, bool)`

GetKeyLengthOk returns a tuple with the KeyLength field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyLength

`func (o *SrtEncryption) SetKeyLength(v SrtEncryptionKeyLength)`

SetKeyLength sets KeyLength field to given value.

### HasKeyLength

`func (o *SrtEncryption) HasKeyLength() bool`

HasKeyLength returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


