//go:build linux
// +build linux

package screen

import (
	"fmt"
	"os"
	"os/exec"
)

type linuxCapturer struct {
	config *Config
}

func newCapturer(cfg *Config) (Capturer, error) {
	if os.Getenv("DISPLAY") == "" {
		return nil, fmt.Errorf("DISPLAY environment variable not set")
	}

	if _, err := exec.LookPath("scrot"); err != nil {
		return nil, fmt.Errorf("no screenshot tool found (need scrot)")
	}

	return &linuxCapturer{config: cfg}, nil
}

func (c *linuxCapturer) Capture() (*image.RGBA, error) {
	// For MVP, use external tools
	return nil, fmt.Errorf("not implemented - use scrot or x11 API")
}

func (c *linuxCapturer) Bounds() (int, int) {
	return 1920, 1080
}
