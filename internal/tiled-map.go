package internal

type Map struct {
	Name, Directory string
	WorldX, WorldY  float32
	FirstTileIds    []uint32
	Version         string        `xml:"version,attr"`
	TiledVersion    string        `xml:"tiledversion,attr"`
	Class           string        `xml:"class,attr"`
	Width           int           `xml:"width,attr"`
	Height          int           `xml:"height,attr"`
	Orientation     string        `xml:"orientation,attr"`
	RenderOrder     string        `xml:"renderorder,attr"`
	TileWidth       int           `xml:"tilewidth,attr"`
	TileHeight      int           `xml:"tileheight,attr"`
	ParallaxOriginX float32       `xml:"parallaxoriginx,attr"`
	ParallaxOriginY float32       `xml:"parallaxoriginy,attr"`
	Infinite        bool          `xml:"infinite,attr"`
	NextLayerID     int           `xml:"nextlayerid,attr"`
	NextObjectID    int           `xml:"nextobjectid,attr"`
	BackgroundColor string        `xml:"backgroundcolor,attr"`
	OutputChunkSize *MapChunkSize `xml:"editorsettings>chunksize"`
	Tilesets        []*Tileset    `xml:"tileset"`
	Layers
	Properties []*Property `xml:"properties>property"`
}

type MapChunkSize struct {
	Width  int `xml:"width,attr"`
	Height int `xml:"height,attr"`
}
