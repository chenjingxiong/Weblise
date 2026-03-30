//go:build darwin
// +build darwin

package input

import (
	"fmt"
	"os/exec"
)

type darwinKeyboard struct{}

func (k *darwinKeyboard) Press(key string) error {
	cmd := exec.Command("cliclick", "k", key)
	return cmd.Run()
}

func (k *darwinKeyboard) Type(text string) error {
	cmd := exec.Command("cliclick", "t:", text)
	return cmd.Run()
}
