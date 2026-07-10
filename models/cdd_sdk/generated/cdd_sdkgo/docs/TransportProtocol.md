# TransportProtocol

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**SrtListener** | [**SrtListenerTransportProtocol**](SrtListenerTransportProtocol.md) |  | 
**SrtCaller** | [**SrtCallerTransportProtocol**](SrtCallerTransportProtocol.md) |  | 
**RistSimpleListener** | [**RistSimpleListenerTransportProtocol**](RistSimpleListenerTransportProtocol.md) |  | 
**RistSimpleCaller** | [**RistSimpleCallerTransportProtocol**](RistSimpleCallerTransportProtocol.md) |  | 
**ZixiPushSender** | [**ZixiPushSenderTransportProtocol**](ZixiPushSenderTransportProtocol.md) |  | 
**ZixiPushReceiver** | [**ZixiPushReceiverTransportProtocol**](ZixiPushReceiverTransportProtocol.md) |  | 
**ZixiPullSender** | [**ZixiPullSenderTransportProtocol**](ZixiPullSenderTransportProtocol.md) |  | 
**ZixiPullReceiver** | [**ZixiPullReceiverTransportProtocol**](ZixiPullReceiverTransportProtocol.md) |  | 
**Rtp** | [**RtpTransportProtocol**](RtpTransportProtocol.md) |  | 

## Methods

### NewTransportProtocol

`func NewTransportProtocol(srtListener SrtListenerTransportProtocol, srtCaller SrtCallerTransportProtocol, ristSimpleListener RistSimpleListenerTransportProtocol, ristSimpleCaller RistSimpleCallerTransportProtocol, zixiPushSender ZixiPushSenderTransportProtocol, zixiPushReceiver ZixiPushReceiverTransportProtocol, zixiPullSender ZixiPullSenderTransportProtocol, zixiPullReceiver ZixiPullReceiverTransportProtocol, rtp RtpTransportProtocol, ) *TransportProtocol`

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


### GetRistSimpleListener

`func (o *TransportProtocol) GetRistSimpleListener() RistSimpleListenerTransportProtocol`

GetRistSimpleListener returns the RistSimpleListener field if non-nil, zero value otherwise.

### GetRistSimpleListenerOk

`func (o *TransportProtocol) GetRistSimpleListenerOk() (*RistSimpleListenerTransportProtocol, bool)`

GetRistSimpleListenerOk returns a tuple with the RistSimpleListener field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRistSimpleListener

`func (o *TransportProtocol) SetRistSimpleListener(v RistSimpleListenerTransportProtocol)`

SetRistSimpleListener sets RistSimpleListener field to given value.


### GetRistSimpleCaller

`func (o *TransportProtocol) GetRistSimpleCaller() RistSimpleCallerTransportProtocol`

GetRistSimpleCaller returns the RistSimpleCaller field if non-nil, zero value otherwise.

### GetRistSimpleCallerOk

`func (o *TransportProtocol) GetRistSimpleCallerOk() (*RistSimpleCallerTransportProtocol, bool)`

GetRistSimpleCallerOk returns a tuple with the RistSimpleCaller field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRistSimpleCaller

`func (o *TransportProtocol) SetRistSimpleCaller(v RistSimpleCallerTransportProtocol)`

SetRistSimpleCaller sets RistSimpleCaller field to given value.


### GetZixiPushSender

`func (o *TransportProtocol) GetZixiPushSender() ZixiPushSenderTransportProtocol`

GetZixiPushSender returns the ZixiPushSender field if non-nil, zero value otherwise.

### GetZixiPushSenderOk

`func (o *TransportProtocol) GetZixiPushSenderOk() (*ZixiPushSenderTransportProtocol, bool)`

GetZixiPushSenderOk returns a tuple with the ZixiPushSender field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZixiPushSender

`func (o *TransportProtocol) SetZixiPushSender(v ZixiPushSenderTransportProtocol)`

SetZixiPushSender sets ZixiPushSender field to given value.


### GetZixiPushReceiver

`func (o *TransportProtocol) GetZixiPushReceiver() ZixiPushReceiverTransportProtocol`

GetZixiPushReceiver returns the ZixiPushReceiver field if non-nil, zero value otherwise.

### GetZixiPushReceiverOk

`func (o *TransportProtocol) GetZixiPushReceiverOk() (*ZixiPushReceiverTransportProtocol, bool)`

GetZixiPushReceiverOk returns a tuple with the ZixiPushReceiver field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZixiPushReceiver

`func (o *TransportProtocol) SetZixiPushReceiver(v ZixiPushReceiverTransportProtocol)`

SetZixiPushReceiver sets ZixiPushReceiver field to given value.


### GetZixiPullSender

`func (o *TransportProtocol) GetZixiPullSender() ZixiPullSenderTransportProtocol`

GetZixiPullSender returns the ZixiPullSender field if non-nil, zero value otherwise.

### GetZixiPullSenderOk

`func (o *TransportProtocol) GetZixiPullSenderOk() (*ZixiPullSenderTransportProtocol, bool)`

GetZixiPullSenderOk returns a tuple with the ZixiPullSender field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZixiPullSender

`func (o *TransportProtocol) SetZixiPullSender(v ZixiPullSenderTransportProtocol)`

SetZixiPullSender sets ZixiPullSender field to given value.


### GetZixiPullReceiver

`func (o *TransportProtocol) GetZixiPullReceiver() ZixiPullReceiverTransportProtocol`

GetZixiPullReceiver returns the ZixiPullReceiver field if non-nil, zero value otherwise.

### GetZixiPullReceiverOk

`func (o *TransportProtocol) GetZixiPullReceiverOk() (*ZixiPullReceiverTransportProtocol, bool)`

GetZixiPullReceiverOk returns a tuple with the ZixiPullReceiver field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZixiPullReceiver

`func (o *TransportProtocol) SetZixiPullReceiver(v ZixiPullReceiverTransportProtocol)`

SetZixiPullReceiver sets ZixiPullReceiver field to given value.


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


