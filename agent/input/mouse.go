package input

import "github.com/go-vgo/robotgo"

type Mouse struct{}

func NewMouse() *Mouse {
	return &Mouse{}
}

func (m *Mouse) Move(x, y int) error {
	return robotgo.MoveMouse(x, y)
}

func (m *Mouse) Click(button string) error {
	return robotgo.MouseClick(button)
}

func (m *Mouse) Press(button string) error {
	return robotgo.MouseToggle("down", button)
}

func (m *Mouse) Release(button string) error {
	return robotgo.MouseToggle("up", button)
}
