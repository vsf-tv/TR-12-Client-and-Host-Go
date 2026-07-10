
# TransportProtocol


## Properties

Name | Type
------------ | -------------
`srtListener` | [SrtListenerTransportProtocol](SrtListenerTransportProtocol.md)
`srtCaller` | [SrtCallerTransportProtocol](SrtCallerTransportProtocol.md)
`ristSimpleListener` | [RistSimpleListenerTransportProtocol](RistSimpleListenerTransportProtocol.md)
`ristSimpleCaller` | [RistSimpleCallerTransportProtocol](RistSimpleCallerTransportProtocol.md)
`zixiPushSender` | [ZixiPushSenderTransportProtocol](ZixiPushSenderTransportProtocol.md)
`zixiPushReceiver` | [ZixiPushReceiverTransportProtocol](ZixiPushReceiverTransportProtocol.md)
`zixiPullSender` | [ZixiPullSenderTransportProtocol](ZixiPullSenderTransportProtocol.md)
`zixiPullReceiver` | [ZixiPullReceiverTransportProtocol](ZixiPullReceiverTransportProtocol.md)
`rtp` | [RtpTransportProtocol](RtpTransportProtocol.md)

## Example

```typescript
import type { TransportProtocol } from ''

// TODO: Update the object below with actual values
const example = {
  "srtListener": null,
  "srtCaller": null,
  "ristSimpleListener": null,
  "ristSimpleCaller": null,
  "zixiPushSender": null,
  "zixiPushReceiver": null,
  "zixiPullSender": null,
  "zixiPullReceiver": null,
  "rtp": null,
} satisfies TransportProtocol

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as TransportProtocol
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


