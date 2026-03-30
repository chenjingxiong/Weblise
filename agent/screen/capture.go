package screen

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
)

type Capturer interface {
	Capture() (*image.RGBA, error)
	Bounds() (width, height int)
}

type Config struct {
	Quality int
	Format  string
}

func DefaultConfig() *Config {
	return &Config{
		Quality: 80,
		Format:  "jpeg",
	}
}

func New(cfg *Config) (Capturer, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}
	return newCapturer(cfg)
}

// newCapturer is defined in platform-specific files
// newLinuxCapturer for Linux
// newWindowsCapturer for Windows
// newDarwinCapturer for Darwin

type mockCapturer struct{}

func (m *mockCapturer) Capture() (*image.RGBA, error) {
	return nil, fmt.Errorf("screen capture not implemented - requires platform-specific libraries")
}

func (m *mockCapturer) Bounds() (int, int) {
	return 1920, 1080
}

// EncodeJPEG encodes an image.RGBA to JPEG bytes
func EncodeJPEG(img *image.RGBA, quality int) ([]byte, error) {
	if img == nil {
		return nil, fmt.Errorf("image is nil")
	}

	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, fmt.Errorf("JPEG encode failed: %w", err)
	}
	return buf.Bytes(), nil
}

// EncodePNG encodes an image.RGBA to PNG bytes
func EncodePNG(img *image.RGBA) ([]byte, error) {
	if img == nil {
		return nil, fmt.Errorf("image is nil")
	}

	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("PNG encode failed: %w", err)
	}
	return buf.Bytes(), nil
}
