package model

type ActionType struct {
	ID              string `json:"id"`
	Key             string `json:"key"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	CommandTemplate string `json:"command_template"` // e.g., "bash -c", "osascript -e"
}
