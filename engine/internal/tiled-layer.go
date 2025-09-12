package internal

type Layer struct {
	Identity
	Tint       string     `xml:"tintcolor,attr"`
	Opacity    float32    `xml:"opacity,attr"`
	Visible    string     `xml:"visible,attr"`
	Locked     bool       `xml:"locked,attr"`
	OffsetX    float32    `xml:"offsetx,attr"`
	OffsetY    float32    `xml:"offsety,attr"`
	ParallaxX  float32    `xml:"parallaxx,attr"`
	ParallaxY  float32    `xml:"parallaxy,attr"`
	Properties []Property `xml:"properties>property"`
}

type LayerGroup struct {
	Layer
	LayersTiles   []LayerTiles   `xml:"layer"`
	LayersObjects []LayerObjects `xml:"objectgroup"`
	LayersImages  []LayerImage   `xml:"imagelayer"`
}

// =================================================================
type LayerTiles struct {
	Layer
	TileData LayerTilesData `xml:"data"`
}
type LayerTilesData struct {
	Encoding string `xml:"encoding,attr"`
	Tiles    string `xml:",chardata"`
}

// =================================================================
type LayerObjects struct {
	Layer
	DrawOrder string        `xml:"draworder,attr"`
	Color     string        `xml:"color,attr"`
	Objects   []LayerObject `xml:"object"`
}
type LayerObject struct {
	Identity
	Width      float64            `xml:"width,attr"`
	Height     float64            `xml:"height,attr"`
	X          float64            `xml:"x,attr"`
	Y          float64            `xml:"y,attr"`
	Rotation   float64            `xml:"rotation,attr"`
	Visible    string             `xml:"visible,attr"`
	Text       LayerObjectText    `xml:"text"`
	Polygon    LayerObjectPolygon `xml:"polygon"`
	Properties []Property         `xml:"properties>property"`
}
type LayerObjectPolygon struct {
	Points string `xml:"points,attr"`
}
type LayerObjectText struct {
	FontFamily string `xml:"fontfamily,attr"`
	FontSize   int    `xml:"pixelsize,attr"`
	WordWrap   bool   `xml:"wrap,attr"`
	Italic     bool   `xml:"italic,attr"`
	Bold       bool   `xml:"bold,attr"`
	Underline  bool   `xml:"underline,attr"`
	Strikeout  bool   `xml:"strikeout,attr"`
	Color      string `xml:"color,attr"`
	AlignX     string `xml:"halign,attr"`
	AlignY     string `xml:"valign,attr"`
	Value      string `xml:",chardata"`
}

// =================================================================
type LayerImage struct {
	Layer
	RepeatX bool         `xml:"repeatx,attr"`
	RepeatY bool         `xml:"repeaty,attr"`
	Image   TilesetImage `xml:"image"`
}
