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
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/host/internal/service"
)

// UploadHandlers handles device upload endpoints.
type UploadHandlers struct {
	thumbnailSvc *service.ThumbnailService
	logSvc       *service.LogService
}

// NewUploadHandlers creates upload handlers.
func NewUploadHandlers(thumbnailSvc *service.ThumbnailService, logSvc *service.LogService) *UploadHandlers {
	return &UploadHandlers{thumbnailSvc: thumbnailSvc, logSvc: logSvc}
}

// UploadThumbnail handles PUT /upload/thumbnail/:deviceId/:sourceId.
func (h *UploadHandlers) UploadThumbnail(c *gin.Context) {
	deviceID := c.Param("deviceId")
	sourceID := c.Param("sourceId")

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body", "code": 400})
		return
	}

	// Determine image type from content-type header or default to jpg
	imageType := "jpg"
	ct := c.GetHeader("Content-Type")
	if strings.Contains(ct, "png") {
		imageType = "png"
	} else if strings.Contains(ct, "jpeg") || strings.Contains(ct, "jpg") {
		imageType = "jpg"
	} else if ext := filepath.Ext(sourceID); ext != "" {
		imageType = strings.TrimPrefix(ext, ".")
	}

	if err := h.thumbnailSvc.StoreThumbnail(deviceID, sourceID, data, imageType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "thumbnail stored"})
}

// UploadLog handles PUT /upload/log/:deviceId.
func (h *UploadHandlers) UploadLog(c *gin.Context) {
	deviceID := c.Param("deviceId")

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read body", "code": 400})
		return
	}

	if err := h.logSvc.StoreLog(deviceID, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "code": 500})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "log stored"})
}
