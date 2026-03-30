//go:build windows
// +build windows

package screen

import (
	"fmt"
	"image"
	"image/color"
	"syscall"
	"unsafe"
)

type windowsCapturer struct {
	config   *Config
	hdc      syscall.Handle
	hdcMem   syscall.Handle
	hBitmap  syscall.Handle
	width    int
	height   int
}

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	gdi32            = syscall.NewLazyDLL("gdi32.dll")
	procGetDC        = user32.NewProc("GetDC")
	procReleaseDC    = user32.NewProc("ReleaseDC")
	procCreateCompatibleDC = gdi32.NewProc("CreateCompatibleDC")
	procCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
	procSelectObject = gdi32.NewProc("SelectObject")
	procBitBlt      = gdi32.NewProc("BitBlt")
	procGetDIBits    = gdi32.NewProc("GetDIBits")
	procDeleteObject = gdi32.NewProc("DeleteObject")
	procDeleteDC     = gdi32.NewProc("DeleteDC")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
)

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
	DIB_RGB_COLORS = 0
	BI_RGB = 0
)

type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression  uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors [1]uint32
}

func newCapturer(cfg *Config) (Capturer, error) {
	return newWindowsCapturer(cfg)
}

func newWindowsCapturer(cfg *Config) (Capturer, error) {
	c := &windowsCapturer{
		config: cfg,
	}

	// Get screen dimensions
	c.width = int(getSystemMetrics(SM_CXSCREEN))
	c.height = int(getSystemMetrics(SM_CYSCREEN))

	if c.width <= 0 || c.height <= 0 {
		return nil, fmt.Errorf("invalid screen dimensions: %dx%d", c.width, c.height)
	}

	return c, nil
}

func getSystemMetrics(index int) int32 {
	ret, _, _ := procGetSystemMetrics.Call(uintptr(index))
	return int32(ret)
}

func (c *windowsCapturer) Bounds() (int, int) {
	return c.width, c.height
}

func (c *windowsCapturer) Capture() (*image.RGBA, error) {
	// Get DC for entire screen
	hdc, _, _ := procGetDC.Call(0)
	if hdc == 0 {
		return nil, fmt.Errorf("failed to get DC")
	}
	defer procReleaseDC.Call(0, hdc)

	// Create compatible DC
	hdcMem, _, _ := procCreateCompatibleDC.Call(hdc)
	if hdcMem == 0 {
		return nil, fmt.Errorf("failed to create compatible DC")
	}
	defer procDeleteDC.Call(hdcMem)

	// Create bitmap
	hBitmap, _, _ := procCreateCompatibleBitmap.Call(hdc, uintptr(c.width), uintptr(c.height))
	if hBitmap == 0 {
		return nil, fmt.Errorf("failed to create bitmap")
	}
	defer procDeleteObject.Call(hBitmap)

	// Select bitmap into DC
	procSelectObject.Call(hdcMem, hBitmap)

	// Copy screen to bitmap
	procBitBlt.Call(hdcMem, 0, 0, uintptr(c.width), uintptr(c.height), hdc, 0, 0, 0x00CC0020) // SRCCOPY

	// Create image and get bitmap data
	img := image.NewRGBA(image.Rect(0, 0, c.width, c.height))

	bmi := BITMAPINFO{
		BmiHeader: BITMAPINFOHEADER{
			BiSize:     uint32(unsafe.Sizeof(BITMAPINFOHEADER{})),
			BiWidth:    int32(c.width),
			BiHeight:   -int32(c.height), // Negative for top-down bitmap
			BiPlanes:   1,
			BiBitCount: 32,
			BiCompression: BI_RGB,
		},
	}

	// Get bitmap bits
	_, _, _ = procGetDIBits.Call(
		hdc,
		hBitmap,
		0,
		uintptr(c.height),
		uintptr(unsafe.Pointer(&img.Pix[0])),
		uintptr(unsafe.Pointer(&bmi)),
		DIB_RGB_COLORS,
	)

	return img, nil
}
