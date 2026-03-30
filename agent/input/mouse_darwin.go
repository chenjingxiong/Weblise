//go:build darwin
// +build darwin

package input

import (
	"fmt"
	"os/exec"
)

type darwinMouse struct{}

func (m *darwinMouse) Move(x, y int) error {
	cmd := exec.Command("cliclick", "m", fmt.Sprintf("%d,%d", x, y))
	return cmd.Run()
}

func (m *darwinMouse) Click(button string) {
	m.Press(button)
	m.Release(button)
}

func (m *darwinMouse) Press(button string) error {
	b := "c"
	if button == "right" {
		b = "r"
	} else if button == "middle" {
		b = "m"
	}
	cmd := exec.Command("cliclick", "kd:"+b)
	return cmd.Run()
}

func (m *darwinMouse) Release(button string) error {
	b := "c"
	if button == "right" {
		b = "r"
	} else if button == "middle" {
		b = "m"
	}
	cmd := exec.Command("cliclick", "ku:"+b)
	return cmd.Run()
}
