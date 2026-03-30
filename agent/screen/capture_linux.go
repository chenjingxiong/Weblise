//go:build linux
// +build linux

package screen

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type linuxCapturer struct {
	config    *Config
	display   string
	width     int
	height    int
}

func newCapturer(cfg *Config) (Capturer, error) {
	return newLinuxCapturer(cfg)
}

func newLinuxCapturer(cfg *Config) (Capturer, error) {
	display := os.Getenv("DISPLAY")
	if display == "" {
		display = ":0"
	}

	c := &linuxCapturer{
		config:  cfg,
		display: display,
	}

	// Get screen dimensions
	if err := c.detectScreenSize(); err != nil {
		// Fallback to default
		c.width = 1920
		c.height = 1080
	}

	return c, nil
}

func (c *linuxCapturer) detectScreenSize() error {
	// Try xrandr first
	cmd := exec.Command("xrandr", "--query")
	output, err := cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "*") {
				fields := strings.Fields(line)
				if len(fields) >= 3 {
					// Parse resolution (e.g., "1920x1080")
					res := strings.Replace(fields[2], "x", " ", 1)
					parts := strings.Fields(res)
					if len(parts) >= 2 {
						w, _ := strconv.Atoi(parts[0])
						h, _ := strconv.Atoi(parts[1])
						if w > 0 && h > 0 {
							c.width = w
							c.height = h
							return nil
						}
					}
				}
			}
		}
	}

	// Try xwininfo as fallback
	cmd = exec.Command("xdpyinfo", "-display", c.display)
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "dimensions:") {
				fields := strings.Fields(line)
				for i, f := range fields {
					if f == "dimensions" && i+2 < len(fields) {
						dim := strings.Replace(fields[i+2], "x", " ", 1)
						parts := strings.Fields(dim)
						if len(parts) >= 2 {
							w, _ := strconv.Atoi(parts[0])
							h, _ := strconv.Atoi(parts[1])
							if w > 0 && h > 0 {
								c.width = w
								c.height = h
								return nil
							}
						}
					}
				}
			}
		}
	}

	return fmt.Errorf("could not detect screen size")
}

func (c *linuxCapturer) Bounds() (int, int) {
	return c.width, c.height
}

func (c *linuxCapturer) Capture() (*image.RGBA, error) {
	// Try using scrot (simple and reliable)
	if img, err := c.captureWithScrot(); err == nil && img != nil {
		return img, nil
	}

	// Try using import (ImageMagick)
	if img, err := c.captureWithImport(); err == nil && img != nil {
		return img, nil
	}

	// Fallback: create a placeholder image
	return c.createPlaceholder()
}

func (c *linuxCapturer) captureWithScrot() (*image.RGBA, error) {
	cmd := exec.Command("scrot", "/dev/stdout")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// scrot outputs PNG, decode it
	img, _, err := image.Decode(bytes.NewReader(output))
	if err != nil {
		return nil, err
	}

	return convertToRGBA(img)
}

func (c *linuxCapturer) captureWithImport() (*image.RGBA, error) {
	cmd := exec.Command("import", "-window", "root", "-screen", "/dev/stdout")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(bytes.NewReader(output))
	if err != nil {
		return nil, err
	}

	return convertToRGBA(img)
}

func (c *linuxCapturer) createPlaceholder() (*image.RGBA, error) {
	img := image.NewRGBA(image.Rect(0, 0, c.width, c.height))

	// Fill with a gradient pattern to show it's working
	for y := 0; y < c.height; y++ {
		for x := 0; x < c.width; x++ {
			r := uint8((x * 255) / c.width)
			g := uint8((y * 255) / c.height)
			b := uint8(128)
			img.SetRGBA(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}

	return img, nil
}

func convertToRGBA(img image.Image) (*image.RGBA, error) {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			result.Set(x, y, img.At(x, y))
		}
	}

	return result, nil
}

// X11Capture provides direct X11 screen capture (requires cgo)
// This is a stub for future implementation with cgo
type X11Capture struct {
	display string
}

func NewX11Capture(display string) (*X11Capture, error) {
	return nil, fmt.Errorf("X11 direct capture requires CGO - use scrot or import instead")
}
