package screen

import (
	"fmt"
	"image"
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
	return &mockCapturer{}, nil
}

type mockCapturer struct{}

func (m *mockCapturer) Capture() (*image.RGBA, error) {
	return nil, fmt.Errorf("screen capture not implemented - requires platform-specific libraries")
}

func (m *mockCapturer) Bounds() (int, int) {
	return 1920, 1080
}

func EncodeJPEG(img *image.RGBA, quality int) ([]byte, error) {
	return nil, fmt.Errorf("JPEG encoding not implemented")
}
