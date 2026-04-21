# TransportProtocol

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SrtListener** | [**SrtListenerTransportProtocol**](SrtListenerTransportProtocol.md) |  | 
**SrtCaller** | [**SrtCallerTransportProtocol**](SrtCallerTransportProtocol.md) |  | 
**RistListener** | [**RistListenerTransportProtocol**](RistListenerTransportProtocol.md) |  | 
**RistCaller** | [**RistCallerTransportProtocol**](RistCallerTransportProtocol.md) |  | 
**ZixiPush** | [**ZixiPushTransportProtocol**](ZixiPushTransportProtocol.md) |  | 
**ZixiPull** | [**ZixiPullTransportProtocol**](ZixiPullTransportProtocol.md) |  | 
**Rtp** | [**RtpTransportProtocol**](RtpTransportProtocol.md) |  | 

## Methods

### NewTransportProtocol

`func NewTransportProtocol(srtListener SrtListenerTransportProtocol, srtCaller SrtCallerTransportProtocol, ristListener RistListenerTransportProtocol, ristCaller RistCallerTransportProtocol, zixiPush ZixiPushTransportProtocol, zixiPull ZixiPullTransportProtocol, rtp RtpTransportProtocol, ) *TransportProtocol`

NewTransportProtocol instantiates a new TransportProtocol object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTransportProtocolWithDefaults

`func NewTransportProtocolWithDefaults() *TransportProtocol`

NewTransportProtocolWithDefaults instantiates a new TransportProtocol object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSrtListener

`func (o *TransportProtocol) GetSrtListener() SrtListenerTransportProtocol`

GetSrtListener returns the SrtListener field if non-nil, zero value otherwise.

### GetSrtListenerOk

`func (o *TransportProtocol) GetSrtListenerOk() (*SrtListenerTransportProtocol, bool)`

GetSrtListenerOk returns a tuple with the SrtListener field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSrtListener

`func (o *TransportProtocol) SetSrtListener(v SrtListenerTransportProtocol)`

SetSrtListener sets SrtListener field to given value.


### GetSrtCaller

`func (o *TransportProtocol) GetSrtCaller() SrtCallerTransportProtocol`

GetSrtCaller returns the SrtCaller field if non-nil, zero value otherwise.

### GetSrtCallerOk

`func (o *TransportProtocol) GetSrtCallerOk() (*SrtCallerTransportProtocol, bool)`

GetSrtCallerOk returns a tuple with the SrtCaller field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSrtCaller

`func (o *TransportProtocol) SetSrtCaller(v SrtCallerTransportProtocol)`

SetSrtCaller sets SrtCaller field to given value.


### GetRistListener

`func (o *TransportProtocol) GetRistListener() RistListenerTransportProtocol`

GetRistListener returns the RistListener field if non-nil, zero value otherwise.

### GetRistListenerOk

`func (o *TransportProtocol) GetRistListenerOk() (*RistListenerTransportProtocol, bool)`

GetRistListenerOk returns a tuple with the RistListener field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRistListener

`func (o *TransportProtocol) SetRistListener(v RistListenerTransportProtocol)`

SetRistListener sets RistListener field to given value.


### GetRistCaller

`func (o *TransportProtocol) GetRistCaller() RistCallerTransportProtocol`

GetRistCaller returns the RistCaller field if non-nil, zero value otherwise.

### GetRistCallerOk

`func (o *TransportProtocol) GetRistCallerOk() (*RistCallerTransportProtocol, bool)`

GetRistCallerOk returns a tuple with the RistCaller field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRistCaller

`func (o *TransportProtocol) SetRistCaller(v RistCallerTransportProtocol)`

SetRistCaller sets RistCaller field to given value.


### GetZixiPush

`func (o *TransportProtocol) GetZixiPush() ZixiPushTransportProtocol`

GetZixiPush returns the ZixiPush field if non-nil, zero value otherwise.

### GetZixiPushOk

`func (o *TransportProtocol) GetZixiPushOk() (*ZixiPushTransportProtocol, bool)`

GetZixiPushOk returns a tuple with the ZixiPush field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZixiPush

`func (o *TransportProtocol) SetZixiPush(v ZixiPushTransportProtocol)`

SetZixiPush sets ZixiPush field to given value.


### GetZixiPull

`func (o *TransportProtocol) GetZixiPull() ZixiPullTransportProtocol`

GetZixiPull returns the ZixiPull field if non-nil, zero value otherwise.

### GetZixiPullOk

`func (o *TransportProtocol) GetZixiPullOk() (*ZixiPullTransportProtocol, bool)`

GetZixiPullOk returns a tuple with the ZixiPull field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZixiPull

`func (o *TransportProtocol) SetZixiPull(v ZixiPullTransportProtocol)`

SetZixiPull sets ZixiPull field to given value.


### GetRtp

`func (o *TransportProtocol) GetRtp() RtpTransportProtocol`

GetRtp returns the Rtp field if non-nil, zero value otherwise.

### GetRtpOk

`func (o *TransportProtocol) GetRtpOk() (*RtpTransportProtocol, bool)`

GetRtpOk returns a tuple with the Rtp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRtp

`func (o *TransportProtocol) SetRtp(v RtpTransportProtocol)`

SetRtp sets Rtp field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


