# WebRtcTransportProtocol

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DtlsSetupRole** | [**DtlsSetupRole**](DtlsSetupRole.md) |  | 
**IceParameters** | [**IceParameters**](IceParameters.md) |  | 
**DtlsFingerprints** | [**[]DtlsFingerprint**](DtlsFingerprint.md) |  | 
**IceServers** | Pointer to [**[]IceServer**](IceServer.md) |  | [optional] 
**FecConfig** | Pointer to [**WebRtcFecConfig**](WebRtcFecConfig.md) |  | [optional] 
**SimpleSettings** | Pointer to [**[]IdAndValue**](IdAndValue.md) |  | [optional] 

## Methods

### NewWebRtcTransportProtocol

`func NewWebRtcTransportProtocol(dtlsSetupRole DtlsSetupRole, iceParameters IceParameters, dtlsFingerprints []DtlsFingerprint, ) *WebRtcTransportProtocol`

NewWebRtcTransportProtocol instantiates a new WebRtcTransportProtocol object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewWebRtcTransportProtocolWithDefaults

`func NewWebRtcTransportProtocolWithDefaults() *WebRtcTransportProtocol`

NewWebRtcTransportProtocolWithDefaults instantiates a new WebRtcTransportProtocol object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDtlsSetupRole

`func (o *WebRtcTransportProtocol) GetDtlsSetupRole() DtlsSetupRole`

GetDtlsSetupRole returns the DtlsSetupRole field if non-nil, zero value otherwise.

### GetDtlsSetupRoleOk

`func (o *WebRtcTransportProtocol) GetDtlsSetupRoleOk() (*DtlsSetupRole, bool)`

GetDtlsSetupRoleOk returns a tuple with the DtlsSetupRole field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDtlsSetupRole

`func (o *WebRtcTransportProtocol) SetDtlsSetupRole(v DtlsSetupRole)`

SetDtlsSetupRole sets DtlsSetupRole field to given value.


### GetIceParameters

`func (o *WebRtcTransportProtocol) GetIceParameters() IceParameters`

GetIceParameters returns the IceParameters field if non-nil, zero value otherwise.

### GetIceParametersOk

`func (o *WebRtcTransportProtocol) GetIceParametersOk() (*IceParameters, bool)`

GetIceParametersOk returns a tuple with the IceParameters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIceParameters

`func (o *WebRtcTransportProtocol) SetIceParameters(v IceParameters)`

SetIceParameters sets IceParameters field to given value.


### GetDtlsFingerprints

`func (o *WebRtcTransportProtocol) GetDtlsFingerprints() []DtlsFingerprint`

GetDtlsFingerprints returns the DtlsFingerprints field if non-nil, zero value otherwise.

### GetDtlsFingerprintsOk

`func (o *WebRtcTransportProtocol) GetDtlsFingerprintsOk() (*[]DtlsFingerprint, bool)`

GetDtlsFingerprintsOk returns a tuple with the DtlsFingerprints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDtlsFingerprints

`func (o *WebRtcTransportProtocol) SetDtlsFingerprints(v []DtlsFingerprint)`

SetDtlsFingerprints sets DtlsFingerprints field to given value.


### GetIceServers

`func (o *WebRtcTransportProtocol) GetIceServers() []IceServer`

GetIceServers returns the IceServers field if non-nil, zero value otherwise.

### GetIceServersOk

`func (o *WebRtcTransportProtocol) GetIceServersOk() (*[]IceServer, bool)`

GetIceServersOk returns a tuple with the IceServers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIceServers

`func (o *WebRtcTransportProtocol) SetIceServers(v []IceServer)`

SetIceServers sets IceServers field to given value.

### HasIceServers

`func (o *WebRtcTransportProtocol) HasIceServers() bool`

HasIceServers returns a boolean if a field has been set.

### GetFecConfig

`func (o *WebRtcTransportProtocol) GetFecConfig() WebRtcFecConfig`

GetFecConfig returns the FecConfig field if non-nil, zero value otherwise.

### GetFecConfigOk

`func (o *WebRtcTransportProtocol) GetFecConfigOk() (*WebRtcFecConfig, bool)`

GetFecConfigOk returns a tuple with the FecConfig field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFecConfig

`func (o *WebRtcTransportProtocol) SetFecConfig(v WebRtcFecConfig)`

SetFecConfig sets FecConfig field to given value.

### HasFecConfig

`func (o *WebRtcTransportProtocol) HasFecConfig() bool`

HasFecConfig returns a boolean if a field has been set.

### GetSimpleSettings

`func (o *WebRtcTransportProtocol) GetSimpleSettings() []IdAndValue`

GetSimpleSettings returns the SimpleSettings field if non-nil, zero value otherwise.

### GetSimpleSettingsOk

`func (o *WebRtcTransportProtocol) GetSimpleSettingsOk() (*[]IdAndValue, bool)`

GetSimpleSettingsOk returns a tuple with the SimpleSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSimpleSettings

`func (o *WebRtcTransportProtocol) SetSimpleSettings(v []IdAndValue)`

SetSimpleSettings sets SimpleSettings field to given value.

### HasSimpleSettings

`func (o *WebRtcTransportProtocol) HasSimpleSettings() bool`

HasSimpleSettings returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


