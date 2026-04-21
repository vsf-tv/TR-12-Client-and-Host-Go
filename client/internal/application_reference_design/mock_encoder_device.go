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
// Simple encoder simulation — mirrors the Python ARD's simple_encoder.py.
// Manages per-channel ffmpeg subprocesses for SRT streaming.
package application_reference_design

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	cddsdkgo "github.com/vsf-tv/TR-12-Client-and-Host-Go/models/cdd_sdk/generated/cdd_sdkgo"
)

// FfmpegPath can be overridden for different systems.
var FfmpegPath = "/opt/homebrew/bin/ffmpeg"

// channelState holds per-channel ffmpeg process, SRT config, and health.
type channelState struct {
	process            *exec.Cmd
	srtIP              string
	srtPort            int
	srtStreamID        string
	srtMinLatencyMs    int
	health             *cddsdkgo.Health
}

// Encoder manages per-channel ffmpeg SRT caller processes and device/channel settings.
type Encoder struct {
	mu              sync.Mutex
	channels        map[string]*channelState   // per-channel process + SRT state
	deviceSettings  map[string]string
	channelSettings map[string]map[string]string // channelID -> key -> value
}

// NewEncoder creates a new encoder instance with default settings.
// Channel settings are initialized lazily when first set.
func NewEncoder() *Encoder {
	return &Encoder{
		deviceSettings:  map[string]string{"sync_clock_source": "NTP"},
		channelSettings: make(map[string]map[string]string),
		channels:        make(map[string]*channelState),
	}
}

func (e *Encoder) getOrCreateChannel(channelID string) *channelState {
	if e.channels[channelID] == nil {
		e.channels[channelID] = &channelState{}
	}
	return e.channels[channelID]
}

// RunningChannel returns true if the given channel's ffmpeg process is active.
func (e *Encoder) RunningChannel(channelID string) bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	ch, ok := e.channels[channelID]
	if !ok || ch.process == nil {
		return false
	}
	if ch.process.ProcessState != nil {
		return false
	}
	if ch.process.Process != nil {
		return ch.process.Process.Signal(syscall.Signal(0)) == nil
	}
	return false
}

// Running returns true if ANY channel's ffmpeg process is active.
func (e *Encoder) Running() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, ch := range e.channels {
		if ch.process != nil && ch.process.ProcessState == nil && ch.process.Process != nil {
			if ch.process.Process.Signal(syscall.Signal(0)) == nil {
				return true
			}
		}
	}
	return false
}

// StartChannel launches ffmpeg for the given channel.
func (e *Encoder) StartChannel(channelID, ip string, port int, streamID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	ch := e.getOrCreateChannel(channelID)

	if ch.process != nil && ch.process.Process != nil {
		if ch.process.Process.Signal(syscall.Signal(0)) == nil {
			fmt.Printf("[%s] Already running\n", channelID)
			return
		}
	}

	fmt.Printf("[%s] ************* Starting *****************\n", channelID)
	srtURL := fmt.Sprintf("srt://%s:%d/%s", ip, port, streamID)
	cmd := exec.Command(FfmpegPath,
		"-f", "avfoundation", "-framerate", "30", "-video_size", "640x480",
		"-i", "0", "-vcodec", "libx264", "-f", "mpegts", srtURL,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	fmt.Printf("[%s] command: %s %v\n", channelID, FfmpegPath, cmd.Args[1:])
	if err := cmd.Start(); err != nil {
		fmt.Printf("[%s] Failed to start ffmpeg: %v\n", channelID, err)
		return
	}
	ch.process = cmd
	ch.srtIP = ip
	ch.srtPort = port
	ch.srtStreamID = streamID
	go func() { _ = cmd.Wait() }()
}

// StopChannel terminates the ffmpeg process for the given channel.
func (e *Encoder) StopChannel(channelID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	ch, ok := e.channels[channelID]
	if !ok || ch.process == nil || ch.process.Process == nil {
		fmt.Printf("[%s] Already stopped\n", channelID)
		return
	}

	fmt.Printf("[%s] ************* Stopping *****************\n", channelID)
	if err := ch.process.Process.Signal(syscall.SIGINT); err != nil {
		fmt.Printf("[%s] SIGINT failed: %v\n", channelID, err)
	} else {
		fmt.Printf("[%s] Sent SIGINT to process %d\n", channelID, ch.process.Process.Pid)
	}

	time.Sleep(3 * time.Second)
	if ch.process.Process != nil {
		if ch.process.Process.Signal(syscall.Signal(0)) == nil {
			fmt.Printf("[%s] Process didn't respond to SIGINT, trying SIGTERM...\n", channelID)
			_ = ch.process.Process.Signal(syscall.SIGTERM)
			time.Sleep(2 * time.Second)
		}
	}
	ch.process = nil
}

// Stop stops all channel processes (used on shutdown).
func (e *Encoder) Stop() {
	e.mu.Lock()
	channelIDs := make([]string, 0, len(e.channels))
	for id := range e.channels {
		channelIDs = append(channelIDs, id)
	}
	e.mu.Unlock()
	for _, id := range channelIDs {
		e.StopChannel(id)
	}
}

// SetDeviceSetting stores a device-level setting.
func (e *Encoder) SetDeviceSetting(key, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.deviceSettings[key] = value
}

// GetDeviceSetting returns a device-level setting value.
func (e *Encoder) GetDeviceSetting(key string) (string, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	v, ok := e.deviceSettings[key]
	return v, ok
}

// SetChannelSetting stores a channel-level simple setting.
func (e *Encoder) SetChannelSetting(channelID, key, value string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.channelSettings[channelID] == nil {
		e.channelSettings[channelID] = make(map[string]string)
	}
	e.channelSettings[channelID][key] = value
}

// GetChannelSetting returns a channel-level simple setting value.
func (e *Encoder) GetChannelSetting(channelID, key string) (string, bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if ch, ok := e.channelSettings[channelID]; ok {
		v, ok := ch[key]
		return v, ok
	}
	return "", false
}

// HandleTransportConfigChange stores SRT config for the given channel.
func (e *Encoder) HandleTransportConfigChange(channelID string, connection *cddsdkgo.Connection) {
	if connection == nil || connection.TransportProtocol == nil {
		fmt.Printf("[%s] Unsupported transport protocol format\n", channelID)
		return
	}
	tp := connection.TransportProtocol
	if tp.SrtCaller == nil {
		fmt.Printf("[%s] No srtCaller in transport protocol — stopping channel\n", channelID)
		e.StopChannel(channelID)
		return
	}
	srt := tp.SrtCaller.SrtCaller
	ip := srt.Address
	port := int(srt.Port)
	streamID := srt.GetStreamId()
	latencyMs := int(srt.GetMinimumLatencyMilliseconds())
	fmt.Printf("[%s] Got SRT config: ip=%s port=%d streamId=%s latencyMs=%d\n", channelID, ip, port, streamID, latencyMs)
	e.mu.Lock()
	ch := e.getOrCreateChannel(channelID)
	ch.srtIP = ip
	ch.srtPort = port
	ch.srtStreamID = streamID
	ch.srtMinLatencyMs = latencyMs
	e.mu.Unlock()
}

// HandleUpdateState processes a channel state change (ACTIVE/IDLE).
func (e *Encoder) HandleUpdateState(channelID string, state cddsdkgo.ChannelState) {
	switch state {
	case cddsdkgo.IDLE:
		fmt.Printf("[%s] Calling stop\n", channelID)
		e.StopChannel(channelID)
	case cddsdkgo.ACTIVE:
		e.mu.Lock()
		ch := e.getOrCreateChannel(channelID)
		ip, port, streamID := ch.srtIP, ch.srtPort, ch.srtStreamID
		e.mu.Unlock()
		if ip != "" && port > 0 {
			if e.RunningChannel(channelID) {
				fmt.Printf("[%s] Restarting with updated settings\n", channelID)
				e.StopChannel(channelID)
			}
			fmt.Printf("[%s] Calling Start\n", channelID)
			e.StartChannel(channelID, ip, port, streamID)
		} else {
			fmt.Printf("[%s] Cannot start: no SRT config available\n", channelID)
		}
	}
}

// GetChannelConnection returns the current SRT connection config for the given channel.
func (e *Encoder) GetChannelConnection(channelID string) *cddsdkgo.Connection {
	e.mu.Lock()
	ch := e.getOrCreateChannel(channelID)
	ip, port, streamID, latencyMs := ch.srtIP, ch.srtPort, ch.srtStreamID, ch.srtMinLatencyMs
	e.mu.Unlock()

	if ip == "" {
		ip = "127.0.0.1"
	}
	if port == 0 {
		port = 5000
	}
	if streamID == "" {
		streamID = "test_stream"
	}
	if latencyMs == 0 {
		latencyMs = 200
	}

	srtProto := cddsdkgo.SrtCallerTransportProtocol{
		Address: ip,
		Port:    float32(port),
	}
	srtProto.StreamId = &streamID
	latencyF := float32(latencyMs)
	srtProto.MinimumLatencyMilliseconds = &latencyF
	tp := cddsdkgo.SrtCallerAsTransportProtocol(cddsdkgo.NewSrtCaller(srtProto))
	conn := cddsdkgo.NewConnection()
	conn.SetTransportProtocol(tp)
	return conn
}

// SetChannelHealth stores health state for a channel.
// level: "DEGRADED" or "CRITICAL"
func (e *Encoder) SetChannelHealth(channelID string, level string, messages []string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	ch := e.getOrCreateChannel(channelID)
	now := time.Now().UTC()
	componentName := "encoder-" + channelID
	var h cddsdkgo.Health
	if level == "CRITICAL" {
		h = cddsdkgo.CriticalAsHealth(cddsdkgo.NewCritical(cddsdkgo.UnhealthyStateDescription{
			Messages:      messages,
			Timestamp:     now,
			ComponentName: componentName,
		}))
	} else {
		h = cddsdkgo.DegradedAsHealth(cddsdkgo.NewDegraded(cddsdkgo.UnhealthyStateDescription{
			Messages:      messages,
			Timestamp:     now,
			ComponentName: componentName,
		}))
	}
	ch.health = &h
}

// ClearChannelHealth resets a channel back to HEALTHY.
func (e *Encoder) ClearChannelHealth(channelID string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	ch := e.getOrCreateChannel(channelID)
	ch.health = nil
}

// GetChannelHealth returns the current health for a channel (HEALTHY if none set).
func (e *Encoder) GetChannelHealth(channelID string) *cddsdkgo.Health {
	e.mu.Lock()
	defer e.mu.Unlock()
	ch, ok := e.channels[channelID]
	if !ok || ch.health == nil {
		h := cddsdkgo.HealthyAsHealth(cddsdkgo.NewHealthy(map[string]interface{}{}))
		return &h
	}
	return ch.health
}

// GetChannelState returns ACTIVE or IDLE for the given channel.
func (e *Encoder) GetChannelState(channelID string) cddsdkgo.ChannelState {
	if e.RunningChannel(channelID) {
		return cddsdkgo.ACTIVE
	}
	return cddsdkgo.IDLE
}

// GetSimulatedBitrate returns a fake bitrate value.
func GetSimulatedBitrate() string {
	ms := time.Now().UnixMilli()
	return fmt.Sprintf("%d", int(math.Mod(float64(ms), 10000))+20000)
}
