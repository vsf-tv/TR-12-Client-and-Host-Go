# RtpTransportProtocol

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Address** | **string** |  | 
**Port** | **float32** |  | 
**SourceAddressFilter** | Pointer to **string** |  | [optional] 
**RtpPayloadType** | Pointer to **float32** |  | [optional] 
**FecConfig** | Pointer to [**RtpFecConfiguration**](RtpFecConfiguration.md) |  | [optional] 

## Methods

### NewRtpTransportProtocol

`func NewRtpTransportProtocol(address string, port float32, ) *RtpTransportProtocol`

NewRtpTransportProtocol instantiates a new RtpTransportProtocol object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRtpTransportProtocolWithDefaults

`func NewRtpTransportProtocolWithDefaults() *RtpTransportProtocol`

NewRtpTransportProtocolWithDefaults instantiates a new RtpTransportProtocol object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAddress

`func (o *RtpTransportProtocol) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *RtpTransportProtocol) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *RtpTransportProtocol) SetAddress(v string)`

SetAddress sets Address field to given value.


### GetPort

`func (o *RtpTransportProtocol) GetPort() float32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *RtpTransportProtocol) GetPortOk() (*float32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *RtpTransportProtocol) SetPort(v float32)`

SetPort sets Port field to given value.


### GetSourceAddressFilter

`func (o *RtpTransportProtocol) GetSourceAddressFilter() string`

GetSourceAddressFilter returns the SourceAddressFilter field if non-nil, zero value otherwise.

### GetSourceAddressFilterOk

`func (o *RtpTransportProtocol) GetSourceAddressFilterOk() (*string, bool)`

GetSourceAddressFilterOk returns a tuple with the SourceAddressFilter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceAddressFilter

`func (o *RtpTransportProtocol) SetSourceAddressFilter(v string)`

SetSourceAddressFilter sets SourceAddressFilter field to given value.

### HasSourceAddressFilter

`func (o *RtpTransportProtocol) HasSourceAddressFilter() bool`

HasSourceAddressFilter returns a boolean if a field has been set.

### GetRtpPayloadType

`func (o *RtpTransportProtocol) GetRtpPayloadType() float32`

GetRtpPayloadType returns the RtpPayloadType field if non-nil, zero value otherwise.

### GetRtpPayloadTypeOk

`func (o *RtpTransportProtocol) GetRtpPayloadTypeOk() (*float32, bool)`

GetRtpPayloadTypeOk returns a tuple with the RtpPayloadType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRtpPayloadType

`func (o *RtpTransportProtocol) SetRtpPayloadType(v float32)`

SetRtpPayloadType sets RtpPayloadType field to given value.

### HasRtpPayloadType

`func (o *RtpTransportProtocol) HasRtpPayloadType() bool`

HasRtpPayloadType returns a boolean if a field has been set.

### GetFecConfig

`func (o *RtpTransportProtocol) GetFecConfig() RtpFecConfiguration`

GetFecConfig returns the FecConfig field if non-nil, zero value otherwise.

### GetFecConfigOk

`func (o *RtpTransportProtocol) GetFecConfigOk() (*RtpFecConfiguration, bool)`

GetFecConfigOk returns a tuple with the FecConfig field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFecConfig

`func (o *RtpTransportProtocol) SetFecConfig(v RtpFecConfiguration)`

SetFecConfig sets FecConfig field to given value.

### HasFecConfig

`func (o *RtpTransportProtocol) HasFecConfig() bool`

HasFecConfig returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


