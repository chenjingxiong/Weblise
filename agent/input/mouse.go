package input

import "fmt"

type Mouse struct{}

func NewMouse() *Mouse {
	return &Mouse{}
}

func (m *Mouse) Move(x, y int) error {
	// TODO: Implement platform-specific mouse control
	return fmt.Errorf("mouse control not implemented")
}

func (m *Mouse) Click(button string) error {
	return fmt.Errorf("mouse click not implemented")
}

func (m *Mouse) Press(button string) error {
	return fmt.Errorf("mouse press not implemented")
}

func (m *Mouse) Release(button string) error {
	return fmt.Errorf("mouse release not implemented")
}
