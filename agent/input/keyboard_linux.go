//go:build linux
// +build linux

package input

import (
	"fmt"
	"os/exec"
)

type linuxKeyboard struct{}

func (k *linuxKeyboard) Press(key string) error {
	if _, err := exec.LookPath("xdotool"); err == nil {
		cmd := exec.Command("xdotool", "key", key)
		return cmd.Run()
	}

	if _, err := exec.LookPath("xte"); err == nil {
		cmd := exec.Command("xte", "key", key)
		return cmd.Run()
	}

	return fmt.Errorf("no keyboard control tool found (need xdotool or xte)")
}

func (k *linuxKeyboard) Type(text string) error {
	if _, err := exec.LookPath("xdotool"); err == nil {
		cmd := exec.Command("xdotool", "type", "--window", "$(xdotool getactivewindow)", text)
		return cmd.Run()
	}

	if _, err := exec.LookPath("xte"); err == nil {
		cmd := exec.Command("xte", "str", text)
		return cmd.Run()
	}

	return fmt.Errorf("no keyboard control tool found")
}
