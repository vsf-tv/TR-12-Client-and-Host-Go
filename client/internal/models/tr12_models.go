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
type CreatePairingCodeFailureReason = tr12models.CreatePairingCodeFailureReason
type CreatePairingCodeExceptionResponseContent = tr12models.CreatePairingCodeExceptionResponseContent
type AuthenticatePairingCodeRequestContent = tr12models.AuthenticatePairingCodeRequestContent
type AuthenticatePairingCodeResponseContent = tr12models.AuthenticatePairingCodeResponseContent
type AuthStatus = tr12models.PairingCodeAuthorizedStatus
type HostSettings = tr12models.HostSettings
type GetHostConfigResponseContent = tr12models.GetHostConfigResponseContent
type DeprovisionRequest = tr12models.DeviceSubscribesToDeprovisionResponseContent
type DeprovisionReason = tr12models.DeprovisionReason
type RequestLogSubscriptionContent = tr12models.DeviceSubscribesToLogSubscriptionResponseContent
type RotateCertificatesRequestContent = tr12models.DeviceSubscribesToCertificateRotationResponseContent
type ThumbnailSubscription = tr12models.ThumbnailSubscription

// RequestThumbnailRequestContent is the thumbnail subscription payload sent to devices.
// It maps directly to DeviceSubscribesToThumbnailSubscriptionResponseContent.
type RequestThumbnailRequestContent = tr12models.DeviceSubscribesToThumbnailSubscriptionResponseContent

// ---- Convenience constants ----

const (
	StateDisconnected = "DISCONNECTED"
	StatePairing      = "PAIRING"
	StateConnecting   = "CONNECTING"
	StateConnected    = "CONNECTED"
	StateReconnecting = "RECONNECTING"
)

var (
	AuthStatusSTANDBY = tr12models.PAIRINGCODEAUTHORIZEDSTATUS_STANDBY
	AuthStatusCLAIMED = tr12models.PAIRINGCODEAUTHORIZEDSTATUS_CLAIMED
)

var (
	PairFailureHostIDMismatch         = tr12models.CREATEPAIRINGCODEFAILUREREASON_HOST_ID_MISMATCH
	PairFailureVersionNotSupported    = tr12models.CREATEPAIRINGCODEFAILUREREASON_VERSION_NOT_SUPPORTED
	PairFailureDeviceTypeNotSupported = tr12models.CREATEPAIRINGCODEFAILUREREASON_DEVICE_TYPE_NOT_SUPPORTED
)

var (
	DeprovisionReasonDeprovisioned = tr12models.DEPROVISIONREASON_DEPROVISIONED
)

var (
	PtrString  = tr12models.PtrString
	PtrFloat32 = tr12models.PtrFloat32
)
