package screen

import (
	"image"
	"time"
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

type Frame struct {
	Data      []byte
	Width     int
	Height    int
	Timestamp time.Time
}
