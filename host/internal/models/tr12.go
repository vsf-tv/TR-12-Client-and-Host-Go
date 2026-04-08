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

// This file bridges the generated TR-12 Smithy models (from the shared submodule) into the
// host service's internal models package. Protocol types are re-exported as aliases
// so existing service code can keep using models.XYZ. Host-service-only types that
// have no generated equivalent are defined directly here.

import (
	tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
)

// --- Re-exported generated TR-12 protocol types ---

type PairRequestContent = tr12models.PairRequestContent
type PairResponseContent = tr12models.PairResponseContent
type PairResult = tr12models.PairResult
type PairSuccessData = tr12models.PairSuccessData
type PairFailureData = tr12models.PairFailureData
type PairFailureReason = tr12models.PairFailureReason
type Success = tr12models.Success
type Failure = tr12models.Failure
type AuthenticateRequestContent = tr12models.AuthenticateRequestContent
type AuthenticateResponseContent = tr12models.AuthenticateResponseContent
type AuthStatus = tr12models.AuthStatus
type HostSettings = tr12models.HostSettings
type RotateCertificatesRequestContent = tr12models.RotateCertificatesRequestContent
type DeprovisionDeviceRequestContent = tr12models.DeprovisionDeviceRequestContent
type DeprovisionReason = tr12models.DeprovisionReason
type RequestThumbnailRequestContent = tr12models.RequestThumbnailRequestContent
type ThumbnailRequest = tr12models.ThumbnailRequest
type RequestLogRequestContent = tr12models.RequestLogRequestContent
type GetHostConfigResponseContent = tr12models.GetHostConfigResponseContent

// Re-export enum constants
const (
	AuthStatusSTANDBY = tr12models.STANDBY
	AuthStatusCLAIMED = tr12models.CLAIMED

	PairFailureHostIDMismatch        = tr12models.HOST_ID_MISMATCH
	PairFailureVersionNotSupported   = tr12models.VERSION_NOT_SUPPORTED
	PairFailureDeviceTypeNotSupported = tr12models.DEVICE_TYPE_NOT_SUPPORTED
)

// Re-export helper functions
var PtrString = tr12models.PtrString
var PtrFloat32 = tr12models.PtrFloat32

// --- Host-service-only types (no generated equivalent) ---

// ClaimRequest body for PUT /authorize/{pairingCode}.
type ClaimRequest struct {
	ExpirationDays      int    `json:"expiration_days,omitempty"`
	LocationName        string `json:"location_name,omitempty"`
	DeviceName          string `json:"device_name,omitempty"`
	RotationIntervalDays int   `json:"rotation_interval_days,omitempty"`
}

// UpdateDeviceMetadata contains editable device metadata.
type UpdateDeviceMetadata struct {
	Name                 string `json:"name,omitempty"`
	Location             string `json:"location,omitempty"`
	RotationIntervalDays int    `json:"rotation_interval_days,omitempty"`
}

// UpdateDeviceRequest is the body for PUT /device/{deviceId}.
// Both fields are optional — omit either to leave it unchanged.
type UpdateDeviceRequest struct {
	Metadata            *UpdateDeviceMetadata `json:"metadata,omitempty"`
	DeviceConfiguration interface{}           `json:"deviceConfiguration,omitempty"`
}
