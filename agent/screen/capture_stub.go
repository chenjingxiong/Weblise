//go:build !linux && !windows && !darwin
// +build !linux,!windows,!darwin

package screen

import (
	"fmt"
)

func newCapturer(cfg *Config) (Capturer, error) {
	return &mockCapturer{}, fmt.Errorf("screen capture not implemented on this platform")
}
