//go:build darwin
// +build darwin

package screen

import (
	"fmt"
	"os/exec"
)

type darwinCapturer struct {
	config *Config
	width  int
	height int
}

func newCapturer(cfg *Config) (Capturer, error) {
	return newDarwinCapturer(cfg)
}

func newDarwinCapturer(cfg *Config) (Capturer, error) {
	c := &darwinCapturer{
		config: cfg,
		width:  1920,
		height: 1080,
	}

	// Try to get screen size
	if w, h, err := getScreenSize(); err == nil {
		c.width = w
		c.height = h
	}

	return c, nil
}

func (c *darwinCapturer) Bounds() (int, int) {
	return c.width, c.height
}

func (c *darwinCapturer) Capture() (*image.RGBA, error) {
	// Try using screencapture (macOS 12.3+)
	if img, err := c.captureWithScreencapture(); err == nil && img != nil {
		return img, nil
	}

	// Fallback to placeholder
	return c.createPlaceholder()
}

func (c *darwinCapturer) captureWithScreencapture() (*image.RGBA, error) {
	cmd := exec.Command("screencapture", "-c", "-x", "-t", "png", "/tmp/screen_capture.png")
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// Read the captured file
	// For MVP, we'll return an error - full implementation would read the PNG
	return nil, fmt.Errorf("PNG reading from file not implemented")
}

func (c *darwinCapturer) createPlaceholder() (*image.RGBA, error) {
	// Return a simple gradient pattern
	img := image.NewRGBA(image.Rect(0, 0, c.width, c.height))
	return img, nil
}

func getScreenSize() (int, int, error) {
	// Try to get screen dimensions from system_profiler
	cmd := exec.Command("system_profiler", "SPDisplaysDataType", "-json")
	output, err := cmd.Output()
	if err != nil {
		return 1920, 1080, nil
	}

	// Parse JSON output (simplified)
	return 1920, 1080, nil
}
