// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// TR-12 protocol models — re-exports from the Smithy-generated tr12models package.
// This file provides type aliases and convenience constants so that consuming code
// can continue to import "internal/models" while using the generated types underneath.
package models

import (
	tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
)

// ---- Type aliases for generated TR-12 models ----

type PairRequestContent = tr12models.PairRequestContent
type PairResponseContent = tr12models.PairResponseContent
type PairResult = tr12models.PairResult
type PairSuccessData = tr12models.PairSuccessData
type PairFailureData = tr12models.PairFailureData
type Success = tr12models.Success
type Failure = tr12models.Failure
type PairFailureReason = tr12models.PairFailureReason
type AuthenticateRequestContent = tr12models.AuthenticateRequestContent
type AuthenticateResponseContent = tr12models.AuthenticateResponseContent
type AuthStatus = tr12models.AuthStatus
type HostSettings = tr12models.HostSettings
type GetHostConfigResponseContent = tr12models.GetHostConfigResponseContent
type DeprovisionDeviceRequestContent = tr12models.DeprovisionDeviceRequestContent
type DeprovisionReason = tr12models.DeprovisionReason
type RequestLogRequestContent = tr12models.RequestLogRequestContent
type RotateCertificatesRequestContent = tr12models.RotateCertificatesRequestContent
type ThumbnailRequest = tr12models.ThumbnailRequest

// RequestThumbnailRequestContent wraps a map of thumbnail subscriptions.
// This type is not in the generated models (it's a container for the MQTT message).
type RequestThumbnailRequestContent struct {
	Requests map[string]ThumbnailRequest `json:"requests"`
}

// ---- Convenience constants ----

const ProtocolVersion = "1.0.0"

// SDK state machine constants
const (
	StateDisconnected = "DISCONNECTED"
	StatePairing      = "PAIRING"
	StateConnecting   = "CONNECTING"
	StateConnected    = "CONNECTED"
	StateReconnecting = "RECONNECTING"
)

// Auth status enum values
var (
	AuthStatusSTANDBY = tr12models.STANDBY
	AuthStatusCLAIMED = tr12models.CLAIMED
)

// Pair failure reason enum values
var (
	FailureHostIDMismatch         = tr12models.HOST_ID_MISMATCH
	FailureVersionNotSupported    = tr12models.VERSION_NOT_SUPPORTED
	FailureDeviceTypeNotSupported = tr12models.DEVICE_TYPE_NOT_SUPPORTED
)

// Deprovision reason enum values
var (
	DeprovisionReasonDeprovisioned = tr12models.DEPROVISIONED
)

// Pointer helpers re-exported from generated utils
var (
	PtrString  = tr12models.PtrString
	PtrFloat32 = tr12models.PtrFloat32
)
