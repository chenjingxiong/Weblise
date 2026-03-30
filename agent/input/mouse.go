package input

import (
	"fmt"
	"runtime"
)

// Mouse provides mouse control interface
type Mouse struct {
	impl mouseImpl
}

// mouseImpl is the platform-specific implementation
type mouseImpl interface {
	Move(x, y int) error
	Click(button string)
	Press(button string) error
	Release(button string) error
}

// NewMouse creates a new mouse controller
func NewMouse() *Mouse {
	var impl mouseImpl

	switch runtime.GOOS {
	case "linux":
		impl = &linuxMouse{}
	case "windows":
		impl = &windowsMouse{}
	case "darwin":
		impl = &darwinMouse{}
	default:
		impl = &mockMouse{}
	}

	return &Mouse{impl: impl}
}

func (m *Mouse) Move(x, y int) error {
	return m.impl.Move(x, y)
}

func (m *Mouse) Click(button string) error {
	m.impl.Click(button)
	return nil
}

func (m *Mouse) Press(button string) error {
	return m.impl.Press(button)
}

func (m *Mouse) Release(button string) error {
	return m.impl.Release(button)
}

// mockMouse is a fallback implementation
type mockMouse struct{}

// Type declarations for all platforms (stubs for non-current platforms)
// linuxMouse is implemented in mouse_linux.go
// windowsMouse is implemented in mouse_windows.go
// darwinMouse is implemented in mouse_darwin.go

// Stub implementations for non-current platforms
type windowsMouse struct{}
type darwinMouse struct{}

func (m *windowsMouse) Move(x, y int) error {
	return fmt.Errorf("mouse control not implemented on this platform")
}
func (m *windowsMouse) Click(string) {}
func (m *windowsMouse) Press(string) error {
	return fmt.Errorf("mouse control not implemented")
}
func (m *windowsMouse) Release(string) error {
	return fmt.Errorf("mouse control not implemented")
}

func (m *darwinMouse) Move(x, y int) error {
	return fmt.Errorf("mouse control not implemented on this platform")
}
func (m *darwinMouse) Click(string) {}
func (m *darwinMouse) Press(string) error {
	return fmt.Errorf("mouse control not implemented")
}
func (m *darwinMouse) Release(string) error {
	return fmt.Errorf("mouse control not implemented")
}

func (m *mockMouse) Move(x, y int) error {
	return fmt.Errorf("mouse control not implemented on %s", runtime.GOOS)
}

func (m *mockMouse) Click(button string) {}
func (m *mockMouse) Press(button string) error {
	return fmt.Errorf("mouse control not implemented")
}
func (m *mockMouse) Release(button string) error {
	return fmt.Errorf("mouse control not implemented")
}
