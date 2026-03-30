//go:build windows
// +build windows

package input

import (
	"fmt"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

type windowsKeyboard struct{}

var (
	procKeybdEvent = user32.NewProc("keybd_event")
)

const (
	KEYEVENTF_KEYUP = 0x0002
	VK_SHIFT       = 0x10
	VK_CONTROL     = 0x11
	VK_MENU        = 0x12
)

func (k *windowsKeyboard) Press(key string) error {
	// Map common keys to virtual key codes
	vk := map[string]uint16{
		"return": 0x0D,
		"enter":   0x0D,
		"tab":     0x09,
		"escape":  0x1B,
		"esc":     0x1B,
		"backspace": 0x08,
		"delete":  0x2E,
		"space":   0x20,
		"up":      0x26,
		"down":    0x28,
		"left":    0x25,
		"right":   0x27,
	}[key]

	code, ok := vk[key]
	if !ok {
		// Single character
		if len(key) == 1 {
			r := []rune(key)[0]
			if r >= 'A' && r <= 'Z' {
				code = uint16(r) - 'A' + 0x41
			} else if r >= 'a' && r <= 'z' {
				code = uint16(r) - 'a' + 0x41
			} else if r >= '0' && r <= '9' {
				code = uint16(r) - '0' + 0x30
			} else {
				return fmt.Errorf("unsupported key: %s", key)
			}
		} else {
			return fmt.Errorf("unsupported key: %s", key)
		}
	}

	// Key down
	procKeybdEvent.Call(uintptr(code), 0, 0, 0)
	// Key up
	procKeybdEvent.Call(uintptr(code), 0, KEYEVENTF_KEYUP, 0)

	return nil
}

func (k *windowsKeyboard) Type(text string) error {
	runes := []rune(text)
	utf16Str := utf16.Encode(runes)

	for i, r := range utf16Str {
		// Simulate Unicode character input via SendMessage
		// This is a simplified version
		procKeybdEvent.Call(uintptr(r), 0, 0, 0)
		procKeybdEvent.Call(uintptr(r), 0, KEYEVENTF_KEYUP, 0)
	}

	return nil
}
