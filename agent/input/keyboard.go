package input

import (
	"fmt"
	"runtime"
)

// Keyboard provides keyboard control interface
type Keyboard struct {
	impl keyboardImpl
}

// keyboardImpl is the platform-specific implementation
type keyboardImpl interface {
	Press(key string) error
	Type(text string) error
}

// NewKeyboard creates a new keyboard controller
func NewKeyboard() *Keyboard {
	var impl keyboardImpl

	switch runtime.GOOS {
	case "linux":
		impl = &linuxKeyboard{}
	case "windows":
		impl = &windowsKeyboard{}
	case "darwin":
		impl = &darwinKeyboard{}
	default:
		impl = &mockKeyboard{}
	}

	return &Keyboard{impl: impl}
}

func (k *Keyboard) Press(key string) error {
	return k.impl.Press(key)
}

func (k *Keyboard) Type(text string) error {
	return k.impl.Type(text)
}

// mockKeyboard is a fallback implementation
type mockKeyboard struct{}

// Type declarations for all platforms (stubs for non-current platforms)
// linuxKeyboard is implemented in keyboard_linux.go
// windowsKeyboard is implemented in keyboard_windows.go
// darwinKeyboard is implemented in keyboard_darwin.go

// Stub implementations for non-current platforms
type windowsKeyboard struct{}
type darwinKeyboard struct{}

func (k *windowsKeyboard) Press(string) error {
	return fmt.Errorf("keyboard control not implemented on this platform")
}
func (k *windowsKeyboard) Type(string) error {
	return fmt.Errorf("keyboard control not implemented on this platform")
}

func (k *darwinKeyboard) Press(string) error {
	return fmt.Errorf("keyboard control not implemented on this platform")
}
func (k *darwinKeyboard) Type(string) error {
	return fmt.Errorf("keyboard control not implemented on this platform")
}

func (k *mockKeyboard) Press(key string) error {
	return fmt.Errorf("keyboard control not implemented on %s", runtime.GOOS)
}

func (k *mockKeyboard) Type(text string) error {
	return fmt.Errorf("keyboard type not implemented on %s", runtime.GOOS)
}
