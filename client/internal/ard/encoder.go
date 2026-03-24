// Copyright 2025 Amazon.com Inc
// Licensed under the Apache License, Version 2.0
//
// Simple encoder simulation — mirrors the Python ARD's simple_encoder.py.
// Manages an ffmpeg subprocess for SRT streaming.
package ard

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"time"

	"github.com/vsf-tv/TR-12-Client-and-Host-Go/client/pkg/cddmodels"
)

// FfmpegPath can be overridden for different systems.
var FfmpegPath = "/opt/homebrew/bin/ffmpeg"

// Encoder manages a simple ffmpeg SRT caller process.
type Encoder struct {
	mu          sync.Mutex
	process     *exec.Cmd
	srtIP       string
	srtPort     int
	srtStreamID string
}

// NewEncoder creates a new encoder instance.
func NewEncoder() *Encoder {
	return &Encoder{}
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

// HandleTransportConfigChange processes a typed Connection update.
func (e *Encoder) HandleTransportConfigChange(channel string, connection *cddmodels.Connection) {
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
func (e *Encoder) HandleUpdateState(channel string, state cddmodels.ChannelState) {
	switch state {
	case cddmodels.IDLE:
		fmt.Println("Calling stop")
		e.Stop()
	case cddmodels.ACTIVE:
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

// GetChannelState returns ACTIVE or IDLE.
func (e *Encoder) GetChannelState(channel string) cddmodels.ChannelState {
	if e.Running() {
		return cddmodels.ACTIVE
	}
	return cddmodels.IDLE
}

// GetSimulatedBitrate returns a fake bitrate value.
func GetSimulatedBitrate() string {
	ms := time.Now().UnixMilli()
	return fmt.Sprintf("%d", int(math.Mod(float64(ms), 10000))+20000)
}
