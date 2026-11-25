package model

type SystemAction struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	Icon     string `json:"icon"`
	Command  string `json:"command"`
	Category string `json:"category"`
}

type SystemVolume struct {
	Volume int  `json:"volume"`
	Muted  bool `json:"muted"`
}

type SystemBrightness struct {
	Brightness float64 `json:"brightness"`
}

// SystemScroll es un desplazamiento en píxeles de la ventana bajo el puntero.
type SystemScroll struct {
	DX int `json:"dx"`
	DY int `json:"dy"`
}

type SystemActions []SystemAction