# WebRtcFecConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FecMechanism** | [**WebRtcFecMechanism**](WebRtcFecMechanism.md) |  | 
**RedPayloadType** | Pointer to **float32** |  | [optional] 
**UlpfecPayloadType** | Pointer to **float32** |  | [optional] 
**TargetOverheadPercentage** | Pointer to **float32** |  | [optional] 

## Methods

### NewWebRtcFecConfig

`func NewWebRtcFecConfig(fecMechanism WebRtcFecMechanism, ) *WebRtcFecConfig`

NewWebRtcFecConfig instantiates a new WebRtcFecConfig object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebRtcFecConfigWithDefaults

`func NewWebRtcFecConfigWithDefaults() *WebRtcFecConfig`

NewWebRtcFecConfigWithDefaults instantiates a new WebRtcFecConfig object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFecMechanism

`func (o *WebRtcFecConfig) GetFecMechanism() WebRtcFecMechanism`

GetFecMechanism returns the FecMechanism field if non-nil, zero value otherwise.

### GetFecMechanismOk

`func (o *WebRtcFecConfig) GetFecMechanismOk() (*WebRtcFecMechanism, bool)`

GetFecMechanismOk returns a tuple with the FecMechanism field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFecMechanism

`func (o *WebRtcFecConfig) SetFecMechanism(v WebRtcFecMechanism)`

SetFecMechanism sets FecMechanism field to given value.


### GetRedPayloadType

`func (o *WebRtcFecConfig) GetRedPayloadType() float32`

GetRedPayloadType returns the RedPayloadType field if non-nil, zero value otherwise.

### GetRedPayloadTypeOk

`func (o *WebRtcFecConfig) GetRedPayloadTypeOk() (*float32, bool)`

GetRedPayloadTypeOk returns a tuple with the RedPayloadType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedPayloadType

`func (o *WebRtcFecConfig) SetRedPayloadType(v float32)`

SetRedPayloadType sets RedPayloadType field to given value.

### HasRedPayloadType

`func (o *WebRtcFecConfig) HasRedPayloadType() bool`

HasRedPayloadType returns a boolean if a field has been set.

### GetUlpfecPayloadType

`func (o *WebRtcFecConfig) GetUlpfecPayloadType() float32`

GetUlpfecPayloadType returns the UlpfecPayloadType field if non-nil, zero value otherwise.

### GetUlpfecPayloadTypeOk

`func (o *WebRtcFecConfig) GetUlpfecPayloadTypeOk() (*float32, bool)`

GetUlpfecPayloadTypeOk returns a tuple with the UlpfecPayloadType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUlpfecPayloadType

`func (o *WebRtcFecConfig) SetUlpfecPayloadType(v float32)`

SetUlpfecPayloadType sets UlpfecPayloadType field to given value.

### HasUlpfecPayloadType

`func (o *WebRtcFecConfig) HasUlpfecPayloadType() bool`

HasUlpfecPayloadType returns a boolean if a field has been set.

### GetTargetOverheadPercentage

`func (o *WebRtcFecConfig) GetTargetOverheadPercentage() float32`

GetTargetOverheadPercentage returns the TargetOverheadPercentage field if non-nil, zero value otherwise.

### GetTargetOverheadPercentageOk

`func (o *WebRtcFecConfig) GetTargetOverheadPercentageOk() (*float32, bool)`

GetTargetOverheadPercentageOk returns a tuple with the TargetOverheadPercentage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargetOverheadPercentage

`func (o *WebRtcFecConfig) SetTargetOverheadPercentage(v float32)`

SetTargetOverheadPercentage sets TargetOverheadPercentage field to given value.

### HasTargetOverheadPercentage

`func (o *WebRtcFecConfig) HasTargetOverheadPercentage() bool`

HasTargetOverheadPercentage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


