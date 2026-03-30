//go:build windows
// +build windows

package screen

import (
	"bytes"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/kbinani/screenshot"
)

type windowsCapturer struct {
	config    *Config
	screenNum int
}

func newCapturer(cfg *Config) (Capturer, error) {
	c := &windowsCapturer{
		config:    cfg,
		screenNum: 0,
	}

	n := screenshot.NumActiveDisplays()
	if n == 0 {
		return nil, os.ErrNotExist
	}

	return c, nil
}

func (c *windowsCapturer) Capture() (*image.RGBA, error) {
	bounds := screenshot.GetDisplayBounds(c.screenNum)
	return screenshot.CaptureRect(bounds)
}

func (c *windowsCapturer) Bounds() (int, int) {
	bounds := screenshot.GetDisplayBounds(c.screenNum)
	return bounds.Dx(), bounds.Dy()
}

func EncodeJPEG(img *image.RGBA, quality int) ([]byte, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func EncodePNG(img *image.RGBA) ([]byte, error) {
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
