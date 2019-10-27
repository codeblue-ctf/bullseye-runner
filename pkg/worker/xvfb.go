package worker

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

var (
	XvfbPath   string
	FFmpegPath string
)

type XvfbManager struct {
	available []uint
	amut      sync.Mutex
	m         sync.Map
}

type XvfbWindow struct {
	XvfbCmd   *exec.Cmd
	FFmpegCmd *exec.Cmd
	display   uint
	width     uint
	height    uint
	depth     uint
}

func NewXvfbManager() *XvfbManager {
	xm := &XvfbManager{}

	for i := uint(0); i < 65536; i++ {
		_, err := os.Stat(fmt.Sprintf("/tmp/.X11-unix/X%d", i))
		if os.IsNotExist(err) {
			xm.available = append(xm.available, i)
		}
	}

	return xm
}

func (x *XvfbManager) AddDisplay(width, height, depth uint) *XvfbWindow {
	x.amut.Lock()
	display := x.available[0]
	x.available = x.available[1:]
	x.amut.Unlock()

	xw := NewXvfbWindow(display, width, height, depth)
	x.m.Store(display, xw)
	return xw
}

func (x *XvfbManager) GetAvailableDisplay(width, height, depth uint) *XvfbWindow {
	var hit uint
	x.m.Range(func(k, v interface{}) bool {
		xw := v.(*XvfbWindow)
		if xw.width == width && xw.height == height && xw.depth == depth &&
			xw.FFmpegCmd.ProcessState.Exited() {
			hit = k.(uint)
			return false
		}
		return true
	})
	xw, ok := x.m.Load(hit)
	if !ok {
		return nil
	}
	return xw.(*XvfbWindow)
}

func NewXvfbWindow(display, width, height, depth uint) *XvfbWindow {
	return &XvfbWindow{
		display: display,
		width:   width,
		height:  height,
		depth:   depth,
	}
}

// GetX11Path returns path to socket should be mounted on docker
func (x *XvfbWindow) GetX11Path() (string, error) {
	if x.XvfbCmd == nil {
		return "", fmt.Errorf("Xvfb process doesn't exist")
	}

	return fmt.Sprintf("/tmp/.X11-unix/X%d", x.display), nil
}

// CreateXvfb will execute Xvfb process in background.
func (x *XvfbWindow) CreateXvfb() error {
	cmd := exec.Command(XvfbPath, fmt.Sprintf(":%d", x.display), "-screen", "0", fmt.Sprintf("%dx%dx%d", x.width, x.height, x.depth))
	err := cmd.Start()

	if err != nil {
		return err
	}

	x.XvfbCmd = cmd

	return nil
}

// Capture will execute FFmpeg in background.
func (x *XvfbWindow) Capture(outfile string, duration time.Duration) error {
	cmd := exec.Command(FFmpegPath, "-f", "x11grab", "-video_size", fmt.Sprintf("%dx%d", x.width, x.height), "-i", fmt.Sprintf(":%d", x.display), "-t", fmt.Sprintf("%f", duration.Seconds()), outfile)
	err := cmd.Start()
	if err != nil {
		return err
	}

	x.FFmpegCmd = cmd

	return nil
}
