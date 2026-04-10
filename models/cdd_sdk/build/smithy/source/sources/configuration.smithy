$version: "2"

namespace com.example.cdd.configuration
use com.example.cdd.common#ChannelState
use com.example.cdd.common#Health
use com.example.cdd.common#IdAndValueList
use com.example.cdd.common#StringList

structure DeviceConfiguration {
    @required
    configurationId: String
    @required
    channels: ChannelConfigurationList
    simpleSettings: IdAndValueList
    health: Health
}

list ChannelConfigurationList {
    member: ChannelConfiguration
}

structure ChannelConfiguration {
    @required
    id: String
    @required
    configurationId: String
    @required
    state: ChannelState
    settings: SettingsChoice
    connection: Connection
    health: Health
}

union SettingsChoice {
    simpleSettings: IdAndValueList
    profile: SettingProfile
}

union RistStreamIdentifier {
    synchronizationSource: Integer
    streamId: String
}

structure SettingProfile {
    @required
    id: String
}

structure Connection {
    transportProtocol: TransportProtocol
}

@length(min: 32, max: 32)
@pattern("^[a-fA-F0-9]+$")
@documentation("A 32-character hexadecimal string.")
string Hex32

@length(min: 64, max: 64)
@pattern("^[a-fA-F0-9]+$")
@documentation("A 64-character hexadecimal string.")
string Hex64

structure EncryptionAes128 {
    @required
    passcode: Hex32
}

structure EncryptionAes256 {
    @required
    passcode: Hex64
}

union EncryptionAes {
    aes128: EncryptionAes128
    aes256: EncryptionAes256
}

union TransportProtocol {
    srtListener: SrtListenerTransportProtocol
    srtCaller: SrtCallerTransportProtocol
    ristListener: RistListenerTransportProtocol
    ristCaller: RistCallerTransportProtocol
    zixiListener: ZixiListenerTransportProtocol
    zixiCaller: ZixiCallerTransportProtocol
    rtp: RtpTransportProtocol
    webRtc: WebRtcTransportProtocol
}

structure SrtListenerTransportProtocol {
    streamId: String
    @required
    @range(min: 1024, max: 65535)
    port: Integer
    @required
    @default(3000)
    minimumLatencyMilliseconds: Integer
    encryption: EncryptionAes
    interface: String
}

structure SrtCallerTransportProtocol {
    streamId: String
    @required
    ip: String
    @required
    @range(min: 1, max: 65535)
    port: Integer
    @required
    @default(3000)
    minimumLatencyMilliseconds: Integer
    encryption: EncryptionAes
}

structure RistListenerTransportProtocol {
    streamId: RistStreamIdentifier
    @required
    port: Integer
    @required
    @default(3000)
    minimumLatencyMilliseconds: Integer
    encryption: EncryptionAes
    interface: String
}

structure RistCallerTransportProtocol {
    streamId: RistStreamIdentifier
    @required
    ip: String
    @required
    port: Integer
    @required
    @default(3000)
    minimumLatencyMilliseconds: Integer
    encryption: EncryptionAes
}

structure ZixiListenerTransportProtocol {
    @required
    streamId: String
    @required
    port: Integer
    @required
    @default(3000)
    minimumLatencyMilliseconds: Integer
    encryption: EncryptionAes
    interface: String
}

structure ZixiCallerTransportProtocol {
    @required
    streamId: String
    @required
    ip: String
    @required
    port: Integer
    @required
    @default(3000)
    minimumLatencyMilliseconds: Integer
    encryption: EncryptionAes
}

// RTP transport — covers unicast and multicast RTP streams including
// SMPTE ST 2022 multicast with optional SMPTE ST 2022-5 FEC.
// Also used when TR-12 is orchestrating an NMOS IS-05 device, since
// NMOS connection management ultimately configures RTP streams described
// by these same SDP parameters.
structure RtpTransportProtocol {
    // Maps to the SDP 'c=' line (Connection Data) — unicast or multicast IP
    @required
    ip: String
    // Maps to the SDP 'm=' line (Media Description) port
    @required
    port: Integer
    // Maps to 'a=source-filter' (IGMPv3 SSM source-specific multicast)
    sourceIpFilter: String
    // Maps to the RTP/AVP payload type value in the SDP 'm=' line
    rtpPayloadType: Integer
    // Optional SMPTE ST 2022-5 FEC configuration
    fecConfig: RtpFecConfiguration
}

// SMPTE ST 2022-5 FEC stream descriptor.
// FEC is transmitted out-of-band as a separate parallel multicast stream —
// it never shares a port with the primary media stream.
structure RtpFecStreamConfig {
    // Multicast IP for the out-of-band FEC stream
    ip: String
    // UDP port for the out-of-band FEC stream
    port: Integer
    // RTP payload type for this FEC stream
    rtpPayloadType: Integer
}

// SMPTE ST 2022-5 FEC configuration.
structure RtpFecConfiguration {
    // 1D column FEC stream (always present when FEC is enabled)
    @required
    columnFec: RtpFecStreamConfig
    // 2D row FEC stream (optional — only used with row+column FEC)
    rowFec: RtpFecStreamConfig
    // FEC matrix column dimension (L)
    matrixColumns: Integer
    // FEC matrix row dimension (D)
    matrixRows: Integer
}

// WebRTC transport — peer-to-peer media using ICE/DTLS/SRTP.
// Carries the explicit signaling parameters needed to establish the connection
// without requiring a full SDP exchange out of band.
structure WebRtcTransportProtocol {
    // Defines who initiates the DTLS handshake (maps to 'a=setup')
    @required
    dtlsSetupRole: DtlsSetupRole
    // ICE credentials for this session
    @required
    iceParameters: IceParameters
    // DTLS certificate fingerprints (one per certificate)
    @required
    dtlsFingerprints: DtlsFingerprintList
    // STUN/TURN servers — only needed when forcing traffic through a relay
    iceServers: IceServerList
    // Optional in-band FEC configuration
    fecConfig: WebRtcFecConfig
    // Flexible key/value escape hatch for media payload negotiation
    // (codec, payload types, extensions, etc.)
    simpleSettings: IdAndValueList
}

// DTLS setup role — maps to the SDP 'a=setup' line
enum DtlsSetupRole {
    ACTPASS
    DTLS_ACTIVE
    DTLS_PASSIVE
}

// FEC mechanism for WebRTC in-band network protection
enum WebRtcFecMechanism {
    ULPFEC
    RED
    FLEXFEC
}

// ICE parameters for session authentication
// Maps to the 'a=ice-ufrag' and 'a=ice-pwd' SDP lines
structure IceParameters {
    @required
    usernameFragment: String
    @required
    password: String
}

// DTLS certificate fingerprint for encryption
// Maps to the 'a=fingerprint' SDP line
structure DtlsFingerprint {
    @required
    algorithm: String
    @required
    value: String
}

list DtlsFingerprintList {
    member: DtlsFingerprint
}

// STUN/TURN server for ICE NAT traversal
structure IceServer {
    @required
    urls: StringList
    username: String
    credential: String
}

list IceServerList {
    member: IceServer
}

// WebRTC in-band FEC configuration
structure WebRtcFecConfig {
    @required
    fecMechanism: WebRtcFecMechanism
    // RTP payload type for RED (a=rtpmap:<pt> red/90000)
    redPayloadType: Integer
    // RTP payload type for ULPFEC (a=rtpmap:<pt> ulpfec/90000)
    ulpfecPayloadType: Integer
    // Target FEC overhead as a percentage of media bitrate
    targetOverheadPercentage: Integer
}
