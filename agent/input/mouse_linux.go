//go:build linux
// +build linux

package input

import (
	"fmt"
	"os/exec"
	"strconv"
)

type linuxMouse struct{}

func (m *linuxMouse) Move(x, y int) error {
	// Try xdotool first
	if _, err := exec.LookPath("xdotool"); err == nil {
		cmd := exec.Command("xdotool", "mousemove", strconv.Itoa(x), strconv.Itoa(y))
		return cmd.Run()
	}

	// Try xte
	if _, err := exec.LookPath("xte"); err == nil {
		cmd := exec.Command("xte", "mousermove", strconv.Itoa(x), strconv.Itoa(y))
		return cmd.Run()
	}

	return fmt.Errorf("no mouse control tool found (need xdotool or xte)")
}

func (m *linuxMouse) Click(button string) {
	_ = m.Press(button)
	_ = m.Release(button)
}

func (m *linuxMouse) Press(button string) error {
	btn := map[string]string{
		"left":   "1",
		"middle": "2",
		"right":  "3",
	}[button]

	if btn == "" {
		btn = "1"
	}

	if _, err := exec.LookPath("xdotool"); err == nil {
		cmd := exec.Command("xdotool", "mousedown", btn)
		return cmd.Run()
	}

	if _, err := exec.LookPath("xte"); err == nil {
		cmd := exec.Command("xte", "mousedown", btn)
		return cmd.Run()
	}

	return fmt.Errorf("no mouse control tool found")
}

func (m *linuxMouse) Release(button string) error {
	btn := map[string]string{
		"left":   "1",
		"middle": "2",
		"right":  "3",
	}[button]

	if btn == "" {
		btn = "1"
	}

	if _, err := exec.LookPath("xdotool"); err == nil {
		cmd := exec.Command("xdotool", "mouseup", btn)
		return cmd.Run()
	}

	if _, err := exec.LookPath("xte"); err == nil {
		cmd := exec.Command("xte", "mouseup", btn)
		return cmd.Run()
	}

	return fmt.Errorf("no mouse control tool found")
}
