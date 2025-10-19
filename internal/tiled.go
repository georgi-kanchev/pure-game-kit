package internal

type World struct {
	Directory, Name      string
	Maps                 []WorldMap `json:"maps"`
	OnlyShowAdjacentMaps bool       `json:"onlyShowAdjacentMaps"`
	Type                 string     `json:"type"`
}

type WorldMap struct {
	FileName string `json:"fileName"`
	Height   int    `json:"height"`
	Width    int    `json:"width"`
	X        int    `json:"x"`
	Y        int    `json:"y"`
}

type Property struct {
	Name       string `xml:"name,attr"`
	Type       string `xml:"type,attr"`
	CustomType string `xml:"propertytype,attr"`
	Value      string `xml:"value,attr"`
	// Properties []Property `xml:"properties>property"`
}

type Identity struct {
	Id    int    `xml:"id,attr"`
	Name  string `xml:"name,attr"`
	Class string `xml:"class,attr"`
}
