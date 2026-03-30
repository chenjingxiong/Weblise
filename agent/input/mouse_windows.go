//go:build windows
// +build windows

package input

import (
	"fmt"
	"syscall"
	"unsafe"
)

type windowsMouse struct{}

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procSetCursorPos = user32.NewProc("SetCursorPos")
	procSendMessage  = user32.NewProc("SendMessageW")
	procMouseEvents  = user32.NewProc("mouse_event")
)

const (
	MOUSEEVENTF_MOVE      = 0x0001
	MOUSEEVENTF_LEFTDOWN  = 0x0002
	MOUSEEVENTF_LEFTUP    = 0x0004
	MOUSEEVENTF_RIGHTDOWN = 0x0008
	MOUSEEVENTF_RIGHTUP   = 0x0010
	MOUSEEVENTF_MIDDLEDOWN = 0x0020
	MOUSEEVENTF_MIDDLEUP  = 0x0040
)

func (m *windowsMouse) Move(x, y int) error {
	ret, _, _ := procSetCursorPos.Call(uintptr(x), uintptr(y))
	if ret == 0 {
		return fmt.Errorf("SetCursorPos failed")
	}
	return nil
}

func (m *windowsMouse) Click(button string) {
	switch button {
	case "left":
		m.Press("left")
		m.Release("left")
	case "right":
		m.Press("right")
		m.Release("right")
	case "middle":
		m.Press("middle")
		m.Release("middle")
	}
}

func (m *windowsMouse) Press(button string) error {
	var flags uintptr
	switch button {
	case "left":
		flags = MOUSEEVENTF_LEFTDOWN
	case "right":
		flags = MOUSEEVENTF_RIGHTDOWN
	case "middle":
		flags = MOUSEEVENTF_MIDDLEDOWN
	default:
		flags = MOUSEEVENTF_LEFTDOWN
	}

	procMouseEvents.Call(flags, 0, 0, 0, 0, 0, 0)
	return nil
}

func (m *windowsMouse) Release(button string) error {
	var flags uintptr
	switch button {
	case "left":
		flags = MOUSEEVENTF_LEFTUP
	case "right":
		flags = MOUSEEVENTF_RIGHTUP
	case "middle":
		flags = MOUSEEVENTF_MIDDLEUP
	default:
		flags = MOUSEEVENTF_LEFTUP
	}

	procMouseEvents.Call(flags, 0, 0, 0, 0, 0, 0)
	return nil
}
