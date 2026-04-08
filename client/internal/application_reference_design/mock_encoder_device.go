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
// Manages an ffmpeg subprocess for SRT streaming.
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

// Encoder manages a simple ffmpeg SRT caller process and holds the current
// device/channel configuration state — the source of truth for actual config.
type Encoder struct {
	mu              sync.Mutex
	process         *exec.Cmd
	srtIP           string
	srtPort         int
	srtStreamID     string
	deviceSettings  map[string]string
	channelSettings map[string]map[string]string // channelID -> key -> value
}

// NewEncoder creates a new encoder instance with default settings.
func NewEncoder() *Encoder {
	return &Encoder{
		deviceSettings: map[string]string{
			"sync_clock_source": "NTP",
		},
		channelSettings: map[string]map[string]string{
			"CH01": {
				"RS01": "1920x1080",
				"FR01": "30",
				"MB01": "10000",
				"RC01": "CBR",
				"CO01": "H.264",
				"GP01": "60",
				"IN01": "SDI1",
			},
		},
	}
}

// Running returns true if the ffmpeg process is active.
func (e *Encoder) Running() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.process == nil {
		return false
	}
	if e.process.ProcessState != nil {
		return false
	}
	if e.process.Process != nil {
		err := e.process.Process.Signal(syscall.Signal(0))
		return err == nil
	}
	return false
}

// Start launches ffmpeg with the given SRT caller parameters.
func (e *Encoder) Start(ip string, port int, streamID string) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.process != nil && e.process.Process != nil {
		if err := e.process.Process.Signal(syscall.Signal(0)); err == nil {
			fmt.Println("Already running")
			return
		}
	}

	fmt.Println("************* Starting *****************")
	srtURL := fmt.Sprintf("srt://%s:%d/%s", ip, port, streamID)
	cmd := exec.Command(FfmpegPath,
		"-f", "avfoundation", "-framerate", "30", "-video_size", "640x480",
		"-i", "0", "-vcodec", "libx264", "-f", "mpegts", srtURL,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	fmt.Printf("command: %s %v\n", FfmpegPath, cmd.Args[1:])
	if err := cmd.Start(); err != nil {
		fmt.Printf("Failed to start ffmpeg: %v\n", err)
		return
	}
	e.process = cmd
	e.srtIP = ip
	e.srtPort = port
	e.srtStreamID = streamID

	go func() { _ = cmd.Wait() }()
}

// Stop terminates the ffmpeg process.
func (e *Encoder) Stop() {
	e.mu.Lock()
	defer e.mu.Unlock()

	fmt.Println("************* Stopping *****************")
	if e.process == nil || e.process.Process == nil {
		fmt.Println("Already stopped")
		return
	}

	if err := e.process.Process.Signal(syscall.SIGINT); err != nil {
		fmt.Printf("SIGINT failed: %v\n", err)
	} else {
		fmt.Printf("Sent SIGINT to process %d\n", e.process.Process.Pid)
	}

	time.Sleep(3 * time.Second)
	if e.process.Process != nil {
		if err := e.process.Process.Signal(syscall.Signal(0)); err == nil {
			fmt.Println("Process didn't respond to SIGINT, trying SIGTERM...")
			_ = e.process.Process.Signal(syscall.SIGTERM)
			time.Sleep(2 * time.Second)
		}
	}
	e.process = nil
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

// HandleTransportConfigChange processes a typed Connection update.
func (e *Encoder) HandleTransportConfigChange(channel string, connection *cddsdkgo.Connection) {
	if connection == nil || connection.TransportProtocol == nil {
		fmt.Println("Unsupported transport protocol format")
		return
	}
	tp := connection.TransportProtocol
	if tp.SrtCaller == nil {
		fmt.Println("No srtCaller in transport protocol — stopping encoder")
		e.Stop()
		return
	}
	srt := tp.SrtCaller.SrtCaller
	ip := srt.Ip
	port := int(srt.Port)
	streamID := srt.GetStreamId()
	fmt.Printf("Got SRT config update: ip=%s port=%d streamId=%s\n", ip, port, streamID)
	e.srtIP = ip
	e.srtPort = port
	e.srtStreamID = streamID
}

// HandleUpdateState processes a channel state change (ACTIVE/IDLE).
func (e *Encoder) HandleUpdateState(channel string, state cddsdkgo.ChannelState) {
	switch state {
	case cddsdkgo.IDLE:
		fmt.Println("Calling stop")
		e.Stop()
	case cddsdkgo.ACTIVE:
		if e.srtIP != "" && e.srtPort > 0 {
			running := e.Running()
			if running {
				fmt.Println("Restarting with updated settings")
				e.Stop()
			}
			fmt.Println("Calling Start")
			e.Start(e.srtIP, e.srtPort, e.srtStreamID)
		} else {
			fmt.Println("Cannot start: no SRT config available")
		}
	}
}

// GetChannelConnection returns the current SRT connection config from stored state.
func (e *Encoder) GetChannelConnection(channelID string) *cddsdkgo.Connection {
	e.mu.Lock()
	ip := e.srtIP
	port := e.srtPort
	streamID := e.srtStreamID
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

	srtProto := cddsdkgo.SrtCallerTransportProtocol{
		Ip:                         ip,
		Port:                       float32(port),
		MinimumLatencyMilliseconds: 200,
	}
	srtProto.StreamId = &streamID
	tp := cddsdkgo.SrtCallerAsTransportProtocol(cddsdkgo.NewSrtCaller(srtProto))
	conn := cddsdkgo.NewConnection()
	conn.SetTransportProtocol(tp)
	return conn
}

// GetChannelState returns ACTIVE or IDLE.
func (e *Encoder) GetChannelState(channel string) cddsdkgo.ChannelState {
	if e.Running() {
		return cddsdkgo.ACTIVE
	}
	return cddsdkgo.IDLE
}

// GetSimulatedBitrate returns a fake bitrate value.
func GetSimulatedBitrate() string {
	ms := time.Now().UnixMilli()
	return fmt.Sprintf("%d", int(math.Mod(float64(ms), 10000))+20000)
}
