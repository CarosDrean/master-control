package model

type Hotkey struct {
	Keycode   int64    `json:"keycode"`
	Modifiers []string `json:"modifiers"`
	Command   string   `json:"command"`
}