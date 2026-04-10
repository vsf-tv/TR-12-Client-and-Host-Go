# RistStreamIdentifier

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SynchronizationSource** | **float32** |  | 
**StreamId** | **string** |  | 

## Methods

### NewRistStreamIdentifier

`func NewRistStreamIdentifier(synchronizationSource float32, streamId string, ) *RistStreamIdentifier`

NewRistStreamIdentifier instantiates a new RistStreamIdentifier object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRistStreamIdentifierWithDefaults

`func NewRistStreamIdentifierWithDefaults() *RistStreamIdentifier`

NewRistStreamIdentifierWithDefaults instantiates a new RistStreamIdentifier object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSynchronizationSource

`func (o *RistStreamIdentifier) GetSynchronizationSource() float32`

GetSynchronizationSource returns the SynchronizationSource field if non-nil, zero value otherwise.

### GetSynchronizationSourceOk

`func (o *RistStreamIdentifier) GetSynchronizationSourceOk() (*float32, bool)`

GetSynchronizationSourceOk returns a tuple with the SynchronizationSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSynchronizationSource

`func (o *RistStreamIdentifier) SetSynchronizationSource(v float32)`

SetSynchronizationSource sets SynchronizationSource field to given value.


### GetStreamId

`func (o *RistStreamIdentifier) GetStreamId() string`

GetStreamId returns the StreamId field if non-nil, zero value otherwise.

### GetStreamIdOk

`func (o *RistStreamIdentifier) GetStreamIdOk() (*string, bool)`

GetStreamIdOk returns a tuple with the StreamId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStreamId

`func (o *RistStreamIdentifier) SetStreamId(v string)`

SetStreamId sets StreamId field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


