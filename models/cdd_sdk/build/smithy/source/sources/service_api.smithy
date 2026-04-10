$version: "2"

namespace com.example.cdd.internal

use aws.protocols#restJson1
use com.example.cdd.common#HostSettings
use com.example.cdd.common#CreatePairingCodeRequest
use com.example.cdd.common#CreatePairingCodeResponse
use com.example.cdd.common#AuthenticatePairingCodeRequest
use com.example.cdd.common#AuthenticatePairingCodeResponse
use com.example.cdd.common#RotateCertificatesRequest
use com.example.cdd.common#DeprovisionRequest
use com.example.cdd.common#ThumbnailSubscription
use com.example.cdd.common#LogRequest
use com.example.cdd.common#HostConfig
use com.example.cdd.common#VersionResponse

@restJson1
service HostServiceApi {
    version: "1.0"
    operations: [
        CreatePairingCode,
        AuthenticatePairingCode,
        // Marker operations to generate MQTT payload models
        RotateCertificates,
        DeprovisionDevice,
        RequestThumbnail,
        RequestLog,
        GetHostConfig,
        GetVersion
    ]
}

@http(method: "POST", uri: "/pair")
operation CreatePairingCode {
    input: CreatePairingCodeRequest
    output: CreatePairingCodeResponse
}

@http(method: "POST", uri: "/authenticate")
operation AuthenticatePairingCode {
    input: AuthenticatePairingCodeRequest
    output: AuthenticatePairingCodeResponse
}

// Marker operations for MQTT payloads - not actual REST endpoints
@http(method: "POST", uri: "/internal/rotate-certs")
operation RotateCertificates {
    input: RotateCertificatesRequest
}

@http(method: "POST", uri: "/internal/deprovision")
operation DeprovisionDevice {
    input: DeprovisionRequest
}

@http(method: "POST", uri: "/internal/thumbnail")
operation RequestThumbnail {
    input: ThumbnailSubscription
}

@http(method: "POST", uri: "/internal/log")
operation RequestLog {
    input: LogRequest
}

@http(method: "GET", uri: "/internal/host-config")
operation GetHostConfig {
    output: HostConfig
}

@http(method: "GET", uri: "/internal/version")
operation GetVersion {
    output: VersionResponse
}
