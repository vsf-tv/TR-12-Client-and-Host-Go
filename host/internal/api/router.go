package api

import (
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
) *gin.Engine {
	r := gin.Default()
	r.Use(cors.Default())

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

	return r
}
