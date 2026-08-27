package model

type Action struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Icon     string `json:"icon"`
	Type     string `json:"type"`
	Command  string `json:"command"`
	Page     string `json:"page"`
	Position int    `json:"position"`
}

type Actions []Action
