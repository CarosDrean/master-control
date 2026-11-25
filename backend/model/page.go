package model

type Page struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	Order   int    `json:"order"`
	Columns int    `json:"columns"`
	Rows    int    `json:"rows"`
	// App vincula la página a una app de la Mac (nombre del bundle sin ".app"): la página
	// es el cajón de acciones que se despliega al tocar el botón de esa app.
	App string `json:"app,omitempty"`
}

type Pages []Page