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
// ThumbnailManager handles periodic upload of device preview images
// based on subscription requests from the host service.
package thumbnails

import (
	"os"
	"sync"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/cddlogger"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/models"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/utils"
	tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
)

// PathResolver returns the local filesystem path for a given channelId.
type PathResolver func(channelID string) string

// Uploader handles a single thumbnail subscription in its own goroutine.
type Uploader struct {
	channelID    string
	request      tr12models.ThumbnailSubscription
	pathResolver PathResolver
	stopCh       chan struct{}
	logger       *cddlogger.CDDLogger
}

func newUploader(channelID string, req tr12models.ThumbnailSubscription, resolver PathResolver, logger *cddlogger.CDDLogger) *Uploader {
	return &Uploader{
		channelID:    channelID,
		request:      req,
		pathResolver: resolver,
		stopCh:       make(chan struct{}),
		logger:       logger,
	}
}

func (u *Uploader) start() {
	go u.run()
}

func (u *Uploader) stop() {
	select {
	case <-u.stopCh:
	default:
		close(u.stopCh)
	}
}

func (u *Uploader) run() {
	u.logger.Infof("Thumbnails: Starting uploader for channel: %s", u.channelID)
	period := int(u.request.GetPeriodSeconds())
	for {
		now := time.Now()
		expires := u.request.GetExpiresAt()
		if u.request.HasExpiresAt() && !now.Before(expires) {
			u.logger.Info("Thumbnails: Subscription Expired.")
			return
		}
		localPath := u.pathResolver(u.channelID)
		if localPath != "" && validateRequestParams(localPath, &u.request, u.logger) {
			if err := utils.UploadFile(localPath, u.request.GetRemotePath(), period, "thumbnail", u.request.GetHeaders()); err != nil {
				u.logger.Errorf("Thumbnails: upload error: %v", err)
			}
		} else if localPath == "" {
			u.logger.Infof("Thumbnails: no thumbnailLocalPath for channel %s yet, skipping cycle", u.channelID)
		}
		for i := 0; i < period*10; i++ {
			select {
			case <-u.stopCh:
				return
			case <-time.After(100 * time.Millisecond):
			}
		}
	}
}

func validateRequestParams(localPath string, req *tr12models.ThumbnailSubscription, logger *cddlogger.CDDLogger) bool {
	if req.HasExpiresAt() && !time.Now().Before(req.GetExpiresAt()) {
		logger.Info("Thumbnail: Request expired.")
		return false
	}
	info, err := os.Stat(localPath)
	if err != nil {
		logger.Infof("Thumbnails: Local path does not exist: %s", localPath)
		return false
	}
	fileSizeKB := float64(info.Size()) / 1024.0
	maxSize := float64(req.GetMaxSizeKB())
	if fileSizeKB > maxSize {
		logger.Infof("Thumbnails: %.0fKB file exceeds %.0fKB limit", fileSizeKB, maxSize)
		return false
	}
	return true
}

// Manager manages all active thumbnail uploaders.
type Manager struct {
	mu           sync.Mutex
	uploaders    map[string]*Uploader
	pathResolver PathResolver
	logger       *cddlogger.CDDLogger
}

// NewManager creates a new ThumbnailManager.
func NewManager(logger *cddlogger.CDDLogger) *Manager {
	return &Manager{
		uploaders: make(map[string]*Uploader),
		logger:    logger,
	}
}

// SetPathResolver sets the function used to resolve thumbnail local paths per channel.
func (m *Manager) SetPathResolver(resolver PathResolver) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pathResolver = resolver
}

// StopAll stops all active thumbnail uploaders.
func (m *Manager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, u := range m.uploaders {
		u.stop()
	}
	m.uploaders = make(map[string]*Uploader)
}

// UpdateThumbnail processes a new thumbnail subscription request.
func (m *Manager) UpdateThumbnail(sub *models.RequestThumbnailRequestContent) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for channelID, req := range sub.Requests {
		if existing, ok := m.uploaders[channelID]; ok {
			existing.stop()
		}
		resolver := m.pathResolver
		if resolver == nil {
			m.logger.Infof("Thumbnails: no path resolver set, skipping channel %s", channelID)
			continue
		}
		u := newUploader(channelID, req, resolver, m.logger)
		m.uploaders[channelID] = u
		u.start()
	}
	return nil
}
