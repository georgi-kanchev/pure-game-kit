package internal

type Tileset struct {
	Identity
	AtlasId         string
	Source          string                 `xml:"source,attr"`
	Version         string                 `xml:"version,attr"`
	TiledVersion    string                 `xml:"tiledversion,attr"`
	FirstTileId     int                    `xml:"firstgid,attr"`
	TileWidth       int                    `xml:"tilewidth,attr"`
	TileHeight      int                    `xml:"tileheight,attr"`
	TileCount       int                    `xml:"tilecount,attr"`
	Columns         int                    `xml:"columns,attr"`
	Spacing         int                    `xml:"spacing,attr"`
	Margin          int                    `xml:"margin,attr"`
	ObjectAlignment string                 `xml:"objectalignment,attr"`
	TileRenderSize  string                 `xml:"tilerendersize,attr"`
	BackgroundColor string                 `xml:"backgroundcolor,attr"`
	FillMode        string                 `xml:"fillmode,attr"`
	Offset          TilesetOffset          `xml:"tileoffset"`
	Grid            TilesetGrid            `xml:"grid"`
	Transformations TilesetTransformations `xml:"transformations"`
	Tiles           []TilesetTile          `xml:"tile"`
	Image           TilesetImage           `xml:"image"`
	Properties      []Property             `xml:"properties>property"`

	MappedTiles map[int]*TilesetTile
}
type TilesetOffset struct {
	X int `xml:"x,attr"`
	Y int `xml:"y,attr"`
}
type TilesetTransformations struct {
	FlipH                    bool `xml:"hflip,attr"`
	FlipV                    bool `xml:"vflip,attr"`
	Rotate                   bool `xml:"rotate,attr"`
	PreferUntransformedTiles bool `xml:"preferuntransformed,attr"`
}
type TilesetGrid struct {
	Orientation string `xml:"orientation,attr"`
	Width       int    `xml:"width,attr"`
	Height      int    `xml:"height,attr"`
}
type TilesetImage struct {
	Source           string `xml:"source,attr"`
	Width            int    `xml:"width,attr"`
	Height           int    `xml:"height,attr"`
	TransparentColor string `xml:"trans,attr"`
}

// =================================================================

type TilesetTile struct {
	Identity                             // no name
	Probability     float32              `xml:"probability,attr"`
	CollisionLayers []*LayerObjects      `xml:"objectgroup"`
	Animation       TilesetTileAnimation `xml:"animation"`
	Properties      []Property           `xml:"properties>property"`

	// animation sequence from engine/flow - any to avoid circular dependency with internal -> flow -> internal
	Sequence any
}
type TilesetTileAnimation struct {
	Frames []TilesetTileFrame `xml:"frame"`
}
type TilesetTileFrame struct {
	TileId   int `xml:"tileid,attr"`
	Duration int `xml:"duration,attr"`
}
