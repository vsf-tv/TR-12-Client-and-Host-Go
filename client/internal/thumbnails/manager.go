// Copyright 2025 Amazon.com Inc
// Licensed under the Apache License, Version 2.0
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

// Uploader handles a single thumbnail subscription in its own goroutine.
type Uploader struct {
	source  string
	request tr12models.ThumbnailRequest
	stopCh  chan struct{}
	logger  *cddlogger.CDDLogger
}

func newUploader(source string, req tr12models.ThumbnailRequest, logger *cddlogger.CDDLogger) *Uploader {
	return &Uploader{
		source:  source,
		request: req,
		stopCh:  make(chan struct{}),
		logger:  logger,
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
	u.logger.Infof("Thumbnails: Starting uploader for source: %s", u.source)
	period := int(u.request.GetPeriod())
	for {
		now := time.Now().Unix()
		if float64(now) >= float64(u.request.GetExpires()) {
			u.logger.Info("Thumbnails: Subscription Expired.")
			return
		}
		if validateRequestParams(&u.request, u.logger) {
			if err := utils.UploadFile(u.request.GetLocalPath(), u.request.GetRemotePath(), period, "thumbnail"); err != nil {
				u.logger.Errorf("Thumbnails: upload error: %v", err)
			}
		}
		// Wait period seconds, checking for stop every 100ms
		for i := 0; i < period*10; i++ {
			select {
			case <-u.stopCh:
				return
			case <-time.After(100 * time.Millisecond):
			}
		}
	}
}

func validateRequestParams(req *tr12models.ThumbnailRequest, logger *cddlogger.CDDLogger) bool {
	if float64(time.Now().Unix()) >= float64(req.GetExpires()) {
		logger.Info("Thumbnail: Request expired.")
		return false
	}
	localPath := req.GetLocalPath()
	info, err := os.Stat(localPath)
	if err != nil {
		logger.Infof("Thumbnails: Local path does not exist: %s", localPath)
		return false
	}
	fileSizeKB := float64(info.Size()) / 1024.0
	maxSize := float64(req.GetMaxSizeKilobyte())
	if fileSizeKB > maxSize {
		logger.Infof("Thumbnails: %.0fKB file exceeds %.0fKB limit", fileSizeKB, maxSize)
		return false
	}
	return true
}

// Manager manages all active thumbnail uploaders.
type Manager struct {
	mu        sync.Mutex
	uploaders map[string]*Uploader
	logger    *cddlogger.CDDLogger
}

// NewManager creates a new ThumbnailManager.
func NewManager(logger *cddlogger.CDDLogger) *Manager {
	return &Manager{
		uploaders: make(map[string]*Uploader),
		logger:    logger,
	}
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
	for key, req := range sub.Requests {
		if existing, ok := m.uploaders[key]; ok {
			existing.stop()
			time.Sleep(200 * time.Millisecond)
		}
		if !validateRequestParams(&req, m.logger) {
			continue
		}
		u := newUploader(key, req, m.logger)
		m.uploaders[key] = u
		u.start()
	}
	return nil
}
