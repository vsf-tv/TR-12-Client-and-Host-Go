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
package thumbnails

import (
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/cddlogger"
	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/internal/models"
	tr12models "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/TR-12-Models/generated/tr12go"
)

func newTestLogger(t *testing.T) *cddlogger.CDDLogger {
	t.Helper()
	log, err := cddlogger.NewWithName(t.TempDir(), "test.log", "test", nil)
	if err != nil {
		t.Fatalf("logger: %v", err)
	}
	t.Cleanup(func() { log.Close() })
	return log
}

func writeTempJPEG(t *testing.T, sizeBytes int) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "thumb*.jpg")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	data := make([]byte, sizeBytes)
	f.Write(data)
	f.Close()
	return f.Name()
}

func makeRequest(localPath, remotePath string, period, expires, maxSizeKB float32) tr12models.ThumbnailRequest {
	p, e, m, lp, rp := period, expires, maxSizeKB, localPath, remotePath
	return tr12models.ThumbnailRequest{
		Period:          &p,
		Expires:         &e,
		MaxSizeKilobyte: &m,
		LocalPath:       &lp,
		RemotePath:      &rp,
	}
}

func makeSub(requests map[string]tr12models.ThumbnailRequest) *models.RequestThumbnailRequestContent {
	return &models.RequestThumbnailRequestContent{Requests: requests}
}

// TestUploadHappyPath: uploader fires and uploads the file to the server.
func TestUploadHappyPath(t *testing.T) {
	var uploadCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			atomic.AddInt32(&uploadCount, 1)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	imgPath := writeTempJPEG(t, 1024) // 1KB
	expires := float32(time.Now().Add(10 * time.Second).Unix())
	m := NewManager(newTestLogger(t))

	m.UpdateThumbnail(makeSub(map[string]tr12models.ThumbnailRequest{
		"SDI-1": makeRequest(imgPath, srv.URL+"/upload", 1, expires, 500),
	}))

	// Wait up to 3 seconds for at least one upload
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadInt32(&uploadCount) > 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	m.StopAll()

	if atomic.LoadInt32(&uploadCount) == 0 {
		t.Fatal("expected at least one upload, got none")
	}
}

// TestExpiredSubscriptionNotUploaded: uploader exits immediately when expires is in the past.
func TestExpiredSubscriptionNotUploaded(t *testing.T) {
	var uploadCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&uploadCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	imgPath := writeTempJPEG(t, 1024)
	expires := float32(time.Now().Add(-5 * time.Second).Unix()) // already expired
	m := NewManager(newTestLogger(t))

	m.UpdateThumbnail(makeSub(map[string]tr12models.ThumbnailRequest{
		"SDI-1": makeRequest(imgPath, srv.URL+"/upload", 1, expires, 500),
	}))

	time.Sleep(500 * time.Millisecond)
	m.StopAll()

	if atomic.LoadInt32(&uploadCount) > 0 {
		t.Fatalf("expected no uploads for expired subscription, got %d", uploadCount)
	}
}

// TestOversizedImageNotUploaded: file exceeding maxSizeKilobyte is skipped.
func TestOversizedImageNotUploaded(t *testing.T) {
	var uploadCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&uploadCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	imgPath := writeTempJPEG(t, 600*1024) // 600KB
	expires := float32(time.Now().Add(10 * time.Second).Unix())
	m := NewManager(newTestLogger(t))

	m.UpdateThumbnail(makeSub(map[string]tr12models.ThumbnailRequest{
		"SDI-1": makeRequest(imgPath, srv.URL+"/upload", 1, expires, 500), // max 500KB
	}))

	time.Sleep(1500 * time.Millisecond)
	m.StopAll()

	if atomic.LoadInt32(&uploadCount) > 0 {
		t.Fatalf("expected no uploads for oversized image, got %d", uploadCount)
	}
}

// TestMissingFileNotUploaded: uploader skips gracefully when localPath doesn't exist.
func TestMissingFileNotUploaded(t *testing.T) {
	var uploadCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&uploadCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	expires := float32(time.Now().Add(10 * time.Second).Unix())
	m := NewManager(newTestLogger(t))

	m.UpdateThumbnail(makeSub(map[string]tr12models.ThumbnailRequest{
		"SDI-1": makeRequest("/nonexistent/path/thumb.jpg", srv.URL+"/upload", 1, expires, 500),
	}))

	time.Sleep(1500 * time.Millisecond)
	m.StopAll()

	if atomic.LoadInt32(&uploadCount) > 0 {
		t.Fatalf("expected no uploads for missing file, got %d", uploadCount)
	}
}

// TestSubscriptionReplacement: sending a new subscription stops the old uploader
// and starts a new one. Old goroutine does not continue uploading.
func TestSubscriptionReplacement(t *testing.T) {
	var uploadCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			atomic.AddInt32(&uploadCount, 1)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	imgPath := writeTempJPEG(t, 1024)
	expires := float32(time.Now().Add(30 * time.Second).Unix())
	m := NewManager(newTestLogger(t))

	// Start first subscription with period=1s
	m.UpdateThumbnail(makeSub(map[string]tr12models.ThumbnailRequest{
		"SDI-1": makeRequest(imgPath, srv.URL+"/upload", 1, expires, 500),
	}))
	time.Sleep(1500 * time.Millisecond)
	countAfterFirst := atomic.LoadInt32(&uploadCount)
	if countAfterFirst == 0 {
		t.Fatal("expected uploads from first subscription")
	}

	// Replace with a new subscription pointing to a different (nonexistent) path
	// — uploads should stop
	m.UpdateThumbnail(makeSub(map[string]tr12models.ThumbnailRequest{
		"SDI-1": makeRequest("/nonexistent/replaced.jpg", srv.URL+"/upload", 1, expires, 500),
	}))
	time.Sleep(1500 * time.Millisecond)
	countAfterReplace := atomic.LoadInt32(&uploadCount)
	m.StopAll()

	// Count should not have increased significantly after replacement
	// (allow 1 in-flight upload that was already in progress)
	if countAfterReplace > countAfterFirst+1 {
		t.Fatalf("old uploader kept running after replacement: before=%d after=%d", countAfterFirst, countAfterReplace)
	}
	t.Logf("uploads before replacement=%d, after=%d", countAfterFirst, countAfterReplace)
}

// TestMultipleSources: two sources run concurrently and both upload.
func TestMultipleSources(t *testing.T) {
	uploads := map[string]*int32{"SDI-1": new(int32), "HDMI-1": new(int32)}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			src := r.URL.Query().Get("src")
			if c, ok := uploads[src]; ok {
				atomic.AddInt32(c, 1)
			}
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	imgPath := writeTempJPEG(t, 1024)
	expires := float32(time.Now().Add(10 * time.Second).Unix())
	m := NewManager(newTestLogger(t))

	m.UpdateThumbnail(makeSub(map[string]tr12models.ThumbnailRequest{
		"SDI-1":  makeRequest(imgPath, srv.URL+"/upload?src=SDI-1", 1, expires, 500),
		"HDMI-1": makeRequest(imgPath, srv.URL+"/upload?src=HDMI-1", 1, expires, 500),
	}))

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadInt32(uploads["SDI-1"]) > 0 && atomic.LoadInt32(uploads["HDMI-1"]) > 0 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	m.StopAll()

	if atomic.LoadInt32(uploads["SDI-1"]) == 0 {
		t.Error("SDI-1 never uploaded")
	}
	if atomic.LoadInt32(uploads["HDMI-1"]) == 0 {
		t.Error("HDMI-1 never uploaded")
	}
}

// TestStopAll: StopAll halts all uploaders promptly.
func TestStopAll(t *testing.T) {
	var uploadCount int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&uploadCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	imgPath := writeTempJPEG(t, 1024)
	expires := float32(time.Now().Add(60 * time.Second).Unix())
	m := NewManager(newTestLogger(t))

	m.UpdateThumbnail(makeSub(map[string]tr12models.ThumbnailRequest{
		"SDI-1":  makeRequest(imgPath, srv.URL+"/upload", 1, expires, 500),
		"HDMI-1": makeRequest(imgPath, srv.URL+"/upload", 1, expires, 500),
	}))
	time.Sleep(1500 * time.Millisecond)
	m.StopAll()

	countAtStop := atomic.LoadInt32(&uploadCount)
	time.Sleep(2 * time.Second) // wait to see if uploads continue
	countAfterStop := atomic.LoadInt32(&uploadCount)

	// Allow at most 1 in-flight upload per source after StopAll
	if countAfterStop > countAtStop+2 {
		t.Fatalf("uploaders kept running after StopAll: at stop=%d, 2s later=%d", countAtStop, countAfterStop)
	}
}
