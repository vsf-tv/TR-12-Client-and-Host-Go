# ZixiEncryption

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Passphrase** | **string** |  | 
**KeyLength** | Pointer to [**ZixiEncryptionKeyLength**](ZixiEncryptionKeyLength.md) |  | [optional] 

## Methods

### NewZixiEncryption

`func NewZixiEncryption(passphrase string, ) *ZixiEncryption`

NewZixiEncryption instantiates a new ZixiEncryption object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewZixiEncryptionWithDefaults

`func NewZixiEncryptionWithDefaults() *ZixiEncryption`

NewZixiEncryptionWithDefaults instantiates a new ZixiEncryption object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPassphrase

`func (o *ZixiEncryption) GetPassphrase() string`

GetPassphrase returns the Passphrase field if non-nil, zero value otherwise.

### GetPassphraseOk

`func (o *ZixiEncryption) GetPassphraseOk() (*string, bool)`

GetPassphraseOk returns a tuple with the Passphrase field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassphrase

`func (o *ZixiEncryption) SetPassphrase(v string)`

SetPassphrase sets Passphrase field to given value.


### GetKeyLength

`func (o *ZixiEncryption) GetKeyLength() ZixiEncryptionKeyLength`

GetKeyLength returns the KeyLength field if non-nil, zero value otherwise.

### GetKeyLengthOk

`func (o *ZixiEncryption) GetKeyLengthOk() (*ZixiEncryptionKeyLength, bool)`

GetKeyLengthOk returns a tuple with the KeyLength field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeyLength

`func (o *ZixiEncryption) SetKeyLength(v ZixiEncryptionKeyLength)`

SetKeyLength sets KeyLength field to given value.

### HasKeyLength

`func (o *ZixiEncryption) HasKeyLength() bool`

HasKeyLength returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


