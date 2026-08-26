package model

type Page struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Order int    `json:"order"`
}

type Pages []Page