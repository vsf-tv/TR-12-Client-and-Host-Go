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
package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/service"
)

// NewRouter creates the gin router with all routes.
func NewRouter(
	deviceSvc *service.DeviceService,
	accountSvc *service.AccountService,
	thumbnailSvc *service.ThumbnailService,
	logSvc *service.LogService,
	consoleDir string,
) *gin.Engine {
	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	pairing := NewPairingHandlers(deviceSvc)
	accounts := NewAccountHandlers(accountSvc)
	devices := NewDeviceHandlers(deviceSvc, thumbnailSvc)
	uploads := NewUploadHandlers(thumbnailSvc, logSvc)

	// Unauthenticated endpoints
	r.POST("/pair", pairing.Pair)
	r.POST("/authenticate", pairing.Authenticate)
	r.POST("/account/register", accounts.Register)
	r.POST("/account/login", accounts.Login)

	// Upload endpoints (called by devices via generated URLs)
	r.PUT("/upload/thumbnail/:deviceId/:sourceId", uploads.UploadThumbnail)
	r.PUT("/upload/log/:deviceId", uploads.UploadLog)

	// Authenticated endpoints (JWT required)
	auth := r.Group("/")
	auth.Use(JWTAuth(accountSvc))
	{
		auth.GET("/devices", devices.ListDevices)
		auth.GET("/device/:deviceId", devices.DescribeDevice)
		auth.PUT("/device/:deviceId", devices.UpdateConfiguration)
		auth.PUT("/authorize/:pairingCode", devices.Claim)
		auth.PUT("/deprovision/:deviceId", devices.Deprovision)
		auth.GET("/thumbnail/:deviceId", devices.GetThumbnail)
		auth.PUT("/credentials/:deviceId", devices.RotateCredentials)
		auth.GET("/account", accounts.GetAccount)
	}

	// Serve console static files if configured
	if consoleDir != "" {
		r.Static("/console", consoleDir)
		r.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/console/")
		})
	}

	return r
}
