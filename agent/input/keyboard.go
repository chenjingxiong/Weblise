package input

import "github.com/go-vgo/robotgo"

type Keyboard struct{}

func NewKeyboard() *Keyboard {
	return &Keyboard{}
}

func (k *Keyboard) Press(key string) error {
	return robotgo.KeyTap(key)
}

func (k *Keyboard) Type(text string) error {
	return robotgo.TypeStr(text)
}
