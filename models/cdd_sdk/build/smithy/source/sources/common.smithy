$version: "2"

namespace com.example.cdd.common

structure ProtocolVersion {
    @default("1.0.0")
    version: String
}

enum ChannelState {
    ACTIVE
    IDLE
}

enum ChannelType {
    SOURCE
    DESTINATION
}

enum AuthStatus {
    STANDBY
    CLAIMED
}

enum SupportedProtocol {
    SRT_LISTENER
    SRT_CALLER
    ZIXI_LISTENER
    ZIXI_CALLER
    RIST_CALLER
    RIST_LISTENER
    RTP
    WEBRTC
}

structure IdAndValue {
    @required
    key: String
    @required
    value: String
}

//TODO: make this a map
list IdAndValueList {
    member: IdAndValue
}

list StringList {
    member: String
}

structure HostSettings {
    @required
    // TODO: should this be an enum?
    // TODO: Do we want to expose iot or keep that as an implementation detail?
    iotProtocolName: String
    @required
    pairingTimeoutSeconds: Integer
    @required
    minIntervalPubSeconds: Integer
    @required
    mqttKeepaliveSeconds: Integer
    @required
    subUpdateTopic: String
    @required
    subUpdateThumbnailSubscriptionTopic: String
    @required
    pubReportSchemaTopic: String
    @required
    pubReportRegistrationTopic: String
    @required
    pubReportStatusTopic: String
    @required
    pubReportActualConfigurationTopic: String
    @required
    subUpdateCertsTopic: String
    @required
    pubDeprovisionTopic: String
    @required
    subDeprovisionTopic: String
    @required
    subUpdateLogSubscriptionTopic: String
}

structure CreatePairingCodeRequest {
    @required
    deviceType: String
    @required
    hostId: String
    @required
    certificateSigningRequest: String
    @required
    version: String
}

enum CreatePairingCodeFailureReason {
    HOST_ID_MISMATCH
    VERSION_NOT_SUPPORTED
    DEVICE_TYPE_NOT_SUPPORTED
}

union CreatePairingCodeResult {
    success: CreatePairingCodeSuccessData
    failure: CreatePairingCodeFailureData
}

structure CreatePairingCodeSuccessData {
    @required
    deviceId: String
    @required
    pairingCode: String
    @required
    accessCode: String
    @required
    pairingTimeoutSeconds: Integer
}

structure CreatePairingCodeFailureData {
    @required
    reason: CreatePairingCodeFailureReason
}

structure CreatePairingCodeResponse {
    @required
    result: CreatePairingCodeResult
}

structure AuthenticatePairingCodeRequest {
    @required
    deviceId: String
    @required
    pairingCode: String
    @required
    accessCode: String
}

structure AuthenticatePairingCodeResponse {
    @required
    status: AuthStatus
    caCertificate: String
    deviceCertificate: String
    mqttUri: String
    regionName: String
    hostSettings: HostSettings
}

structure RotateCertificatesRequest {
    @required
    mqttUri: String
    @required
    deviceCertificate: String
    @required
    regionName: String
}

enum DeprovisionReason {
    DEPROVISIONED
    EXPIRED
    UNKNOWN
}

structure DeprovisionRequest {
    reason: DeprovisionReason
    //TODO: use timestamp instead?
    time: Integer
}

structure ThumbnailRequest {
    periodSeconds: Integer
    // Unix epoch seconds at which the subscription expires
    expiresAtEpochSeconds: Integer
    maxSizeKilobyte: Integer
    localPath: String
    remotePath: String
    // HTTP headers to include when uploading thumbnails — required for non-AWS
    // storage backends that need authentication or content-type headers
    headers: StringMap
}

map StringMap {
    key: String
    value: String
}

structure ThumbnailSubscription {
    @required
    requests: ThumbnailRequestMap
}

map ThumbnailRequestMap {
    key: String
    value: ThumbnailRequest
}

// All fields are optional — a partial LogRequest instructs the device to upload
// any available logs to the provided remotePath before expiresAtEpochSeconds.
structure LogRequest {
    // Unix epoch seconds at which this upload request expires
    expiresAtEpochSeconds: Integer
    remotePath: String
}

structure HostConfig {
    @required
    serviceId: String
    @required
    serviceName: String
    @required
    deviceTypes: StringList
    @required
    thumbnailMaxSizeKB: Integer
    @required
    logFileMaxSizeKB: Integer
    @required
    createPairingCodeUrl: String
    @required
    authenticatePairingCodeUrl: String
}

union Health {
    healthy: HealthyState
    degraded: DegradedState
    critical: CriticalState
}

/// Healthy state — only the state indicator, no additional fields needed.
structure HealthyState {}

/// Degraded or critical state — all diagnostic fields required.
structure DegradedState {
    @required
    messages: StringList
    @required
    @timestampFormat("date-time")
    timestamp: Timestamp
    @required
    componentName: String
}

structure CriticalState {
    @required
    messages: StringList
    @required
    @timestampFormat("date-time")
    timestamp: Timestamp
    @required
    componentName: String
}

enum HealthLevel {
    HEALTHY
    DEGRADED
    CRITICAL
}

structure VersionResponse {
    @required
    version: ProtocolVersion
}
