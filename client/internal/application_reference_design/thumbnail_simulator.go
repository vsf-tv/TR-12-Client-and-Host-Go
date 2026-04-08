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
// ThumbnailSimulator — cycles through source images and copies them to a
// destination path, simulating an encoder emitting preview thumbnails.
package application_reference_design

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

// ThumbnailSimulator copies images from a source directory to a destination file on a timer.
type ThumbnailSimulator struct {
	sourceDir  string
	dest       string
	interval   time.Duration
	name       string
	files      []string
	imageIndex int
	running    bool
	stopCh     chan struct{}
	mu         sync.Mutex
}

// NewThumbnailSimulator creates a new simulator.
func NewThumbnailSimulator(sourceDir, dest string, intervalSec int, name string) (*ThumbnailSimulator, error) {
	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		return nil, fmt.Errorf("source directory does not exist: %s: %w", sourceDir, err)
	}
	var files []string
	for _, e := range entries {
		if !e.IsDir() {
			files = append(files, filepath.Join(sourceDir, e.Name()))
		}
	}
	sort.Strings(files)
	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in source directory: %s", sourceDir)
	}

	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("destination directory not writable: %s: %w", destDir, err)
	}

	return &ThumbnailSimulator{
		sourceDir: sourceDir,
		dest:      dest,
		interval:  time.Duration(intervalSec) * time.Second,
		name:      name,
		files:     files,
		stopCh:    make(chan struct{}),
	}, nil
}

// Start begins the thumbnail simulation in a goroutine.
func (t *ThumbnailSimulator) Start() {
	t.mu.Lock()
	t.running = true
	t.mu.Unlock()
	go t.run()
}

// Stop signals the simulator to stop.
func (t *ThumbnailSimulator) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.running {
		t.running = false
		close(t.stopCh)
	}
}

func (t *ThumbnailSimulator) run() {
	for {
		select {
		case <-t.stopCh:
			return
		default:
		}

		image := t.pickImage()
		tempFile := filepath.Join(filepath.Dir(t.dest), fmt.Sprintf("temp_%s", t.name))

		if err := copyFile(image, tempFile); err != nil {
			fmt.Printf("Error writing TN image to disk: %v\n", err)
			os.Remove(tempFile)
		} else {
			// Atomic move
			if err := os.Rename(tempFile, t.dest); err != nil {
				fmt.Printf("Error moving TN image: %v\n", err)
			}
			// Touch to update mtime
			now := time.Now()
			os.Chtimes(t.dest, now, now)
		}

		select {
		case <-t.stopCh:
			return
		case <-time.After(t.interval):
		}
	}
}

func (t *ThumbnailSimulator) pickImage() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	path := t.files[t.imageIndex]
	t.imageIndex = (t.imageIndex + 1) % len(t.files)
	return path
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}
