package input

import "fmt"

type Keyboard struct{}

func NewKeyboard() *Keyboard {
	return &Keyboard{}
}

func (k *Keyboard) Press(key string) error {
	return fmt.Errorf("keyboard control not implemented")
}

func (k *Keyboard) Type(text string) error {
	return fmt.Errorf("keyboard type not implemented")
}
