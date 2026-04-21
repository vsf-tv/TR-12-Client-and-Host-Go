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
package models

import (
	tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
)

// --- Re-exported generated TR-12 protocol types ---

type CreatePairingCodeRequestContent = tr12models.CreatePairingCodeRequestContent
type CreatePairingCodeResponseContent = tr12models.CreatePairingCodeResponseContent
type CreatePairingCodeResult = tr12models.CreatePairingCodeResult
type CreatePairingCodeSuccessData = tr12models.CreatePairingCodeSuccessData
type CreatePairingCodeFailureData = tr12models.CreatePairingCodeFailureData
type CreatePairingCodeFailureReason = tr12models.CreatePairingCodeFailureReason
type Success = tr12models.Success
type Failure = tr12models.Failure
type AuthenticatePairingCodeRequestContent = tr12models.AuthenticatePairingCodeRequestContent
type AuthenticatePairingCodeResponseContent = tr12models.AuthenticatePairingCodeResponseContent
type AuthStatus = tr12models.PairingCodeAuthorizedStatus
type HostSettings = tr12models.HostSettings
type RotateCertificatesRequestContent = tr12models.RotateCertificatesRequestContent
type DeprovisionRequest = tr12models.DeprovisionDeviceRequestContent
type DeprovisionReason = tr12models.DeprovisionReason
type RequestThumbnailRequestContent = tr12models.RequestThumbnailRequestContent
type ThumbnailRequest = tr12models.ThumbnailRequest
type RequestLogRequestContent = tr12models.RequestLogRequestContent
type GetHostConfigResponseContent = tr12models.GetHostConfigResponseContent

// Re-export enum constants
var (
	AuthStatusSTANDBY = tr12models.STANDBY
	AuthStatusCLAIMED = tr12models.CLAIMED

	PairFailureHostIDMismatch         = tr12models.HOST_ID_MISMATCH
	PairFailureVersionNotSupported    = tr12models.VERSION_NOT_SUPPORTED
	PairFailureDeviceTypeNotSupported = tr12models.DEVICE_TYPE_NOT_SUPPORTED
)

// Re-export helper functions
var PtrString = tr12models.PtrString
var PtrFloat32 = tr12models.PtrFloat32

// --- Host-service-only types ---

type ClaimRequest struct {
	ExpirationDays       int    `json:"expiration_days,omitempty"`
	LocationName         string `json:"location_name,omitempty"`
	DeviceName           string `json:"device_name,omitempty"`
	RotationIntervalDays int    `json:"rotation_interval_days,omitempty"`
}

type UpdateDeviceMetadata struct {
	Name                 string `json:"name,omitempty"`
	Location             string `json:"location,omitempty"`
	RotationIntervalDays int    `json:"rotation_interval_days,omitempty"`
}

type UpdateDeviceRequest struct {
	Metadata            *UpdateDeviceMetadata `json:"metadata,omitempty"`
	DeviceConfiguration interface{}           `json:"deviceConfiguration,omitempty"`
}
