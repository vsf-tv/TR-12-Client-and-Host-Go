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
package models

import (
	tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
)

// ---- Type aliases for generated TR-12 models ----

type CreatePairingCodeRequestContent = tr12models.CreatePairingCodeRequestContent
type CreatePairingCodeResponseContent = tr12models.CreatePairingCodeResponseContent
type CreatePairingCodeResult = tr12models.CreatePairingCodeResult
type CreatePairingCodeSuccessData = tr12models.CreatePairingCodeSuccessData
type CreatePairingCodeFailureData = tr12models.CreatePairingCodeFailureData
type CreatePairingCodeFailureReason = tr12models.CreatePairingCodeFailureReason
type AuthenticatePairingCodeRequestContent = tr12models.AuthenticatePairingCodeRequestContent
type AuthenticatePairingCodeResponseContent = tr12models.AuthenticatePairingCodeResponseContent
type AuthStatus = tr12models.AuthStatus
type HostSettings = tr12models.HostSettings
type GetHostConfigResponseContent = tr12models.GetHostConfigResponseContent
type DeprovisionRequest = tr12models.DeprovisionDeviceRequestContent
type DeprovisionReason = tr12models.DeprovisionReason
type RequestLogRequestContent = tr12models.RequestLogRequestContent
type RotateCertificatesRequestContent = tr12models.RotateCertificatesRequestContent
type ThumbnailRequest = tr12models.ThumbnailRequest
type Success = tr12models.Success
type Failure = tr12models.Failure

// RequestThumbnailRequestContent wraps a map of thumbnail subscriptions.
type RequestThumbnailRequestContent struct {
	Requests map[string]ThumbnailRequest `json:"requests"`
}

// ---- Convenience constants ----

const ProtocolVersion = "1.0.0"

const (
	StateDisconnected = "DISCONNECTED"
	StatePairing      = "PAIRING"
	StateConnecting   = "CONNECTING"
	StateConnected    = "CONNECTED"
	StateReconnecting = "RECONNECTING"
)

var (
	AuthStatusSTANDBY = tr12models.STANDBY
	AuthStatusCLAIMED = tr12models.CLAIMED
)

var (
	PairFailureHostIDMismatch         = tr12models.HOST_ID_MISMATCH
	PairFailureVersionNotSupported    = tr12models.VERSION_NOT_SUPPORTED
	PairFailureDeviceTypeNotSupported = tr12models.DEVICE_TYPE_NOT_SUPPORTED
)

var (
	DeprovisionReasonDeprovisioned = tr12models.DEPROVISIONED
)

var (
	PtrString  = tr12models.PtrString
	PtrFloat32 = tr12models.PtrFloat32
)
