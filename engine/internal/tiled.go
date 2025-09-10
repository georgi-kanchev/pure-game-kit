package internal

import (
	"encoding/xml"
)

type TiledData struct {
	Directory       string
	XMLName         xml.Name       `xml:"map"`
	Version         string         `xml:"version,attr"`
	TiledVersion    string         `xml:"tiledversion,attr"`
	Class           string         `xml:"class,attr"`
	Width           int            `xml:"width,attr"`
	Height          int            `xml:"height,attr"`
	Orientation     string         `xml:"orientation,attr"`
	RenderOrder     string         `xml:"renderorder,attr"`
	TileWidth       int            `xml:"tilewidth,attr"`
	TileHeight      int            `xml:"tileheight,attr"`
	ParallaxOriginX float32        `xml:"parallaxoriginx,attr"`
	ParallaxOriginY float32        `xml:"parallaxoriginy,attr"`
	Infinite        bool           `xml:"infinite,attr"`
	NextLayerID     int            `xml:"nextlayerid,attr"`
	NextObjectID    int            `xml:"nextobjectid,attr"`
	OutputChunkSize ChunkSize      `xml:"editorsettings>chunksize"`
	BackgroundColor string         `xml:"backgroundcolor,attr"`
	Tilesets        []Tileset      `xml:"tileset"`
	Groups          []LayerGroup   `xml:"group"`
	LayersTiles     []LayerTiles   `xml:"layer"`
	LayersObjects   []LayerObjects `xml:"objectgroup"`
	LayersImages    []LayerImage   `xml:"imagelayer"`
}

type Tileset struct {
	Identity
	FirstGID        int             `xml:"firstgid,attr"`
	TileWidth       int             `xml:"tilewidth,attr"`
	TileHeight      int             `xml:"tileheight,attr"`
	TileCount       int             `xml:"tilecount,attr"`
	Columns         int             `xml:"columns,attr"`
	Spacing         int             `xml:"spacing,attr"`
	Margin          int             `xml:"margin,attr"`
	ObjectAlignment string          `xml:"objectalignment,attr"`
	TileRenderSize  string          `xml:"tilerendersize,attr"`
	BackgroundColor string          `xml:"backgroundcolor,attr"`
	FillMode        string          `xml:"fillmode,attr"`
	TileOffset      TileOffset      `xml:"tileoffset"`
	Grid            Grid            `xml:"grid"`
	Transformations Transformations `xml:"transformations"`
	PerTileData     []TilesetTile   `xml:"tile"`
	Image           Image           `xml:"image"`
	Properties      []Property      `xml:"properties>property"`
}

type TilesetTile struct {
	Id             string         `xml:"id,attr"`
	TilesCollision []LayerObjects `xml:"objectgroup"`
}
type ChunkSize struct {
	Width  int `xml:"width,attr"`
	Height int `xml:"height,attr"`
}

type TileOffset struct {
	X int `xml:"x,attr"`
	Y int `xml:"y,attr"`
}
type Grid struct {
	Orientation string `xml:"orientation,attr"`
	Width       int    `xml:"width,attr"`
	Height      int    `xml:"height,attr"`
}
type Transformations struct {
	FlipH                    bool `xml:"hflip,attr"`
	FlipV                    bool `xml:"vflip,attr"`
	Rotate                   bool `xml:"rotate,attr"`
	PreferUntransformedTiles bool `xml:"preferuntransformed,attr"`
}

type Image struct {
	Source           string `xml:"source,attr"`
	Width            int    `xml:"width,attr"`
	Height           int    `xml:"height,attr"`
	TransparentColor string `xml:"trans,attr"`
}

type Property struct {
	Name       string     `xml:"name,attr"`
	Type       string     `xml:"type,attr"`
	CustomType string     `xml:"propertytype,attr"`
	Value      string     `xml:"value,attr"`
	Properties []Property `xml:"properties>property"`
}

type Identity struct {
	ID    int    `xml:"id,attr"`
	Name  string `xml:"name,attr"`
	Class string `xml:"class,attr"`
}

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

type LayerTiles struct {
	Layer
	Data Data `xml:"data"`
}

type Data struct {
	Encoding string `xml:"encoding,attr"`
	Tiles    string `xml:",chardata"`
}

type LayerObjects struct {
	Layer
	DrawOrder string   `xml:"draworder,attr"`
	Color     string   `xml:"color,attr"`
	Objects   []Object `xml:"object"`
}
type LayerImage struct {
	Layer
	RepeatX bool  `xml:"repeatx,attr"`
	RepeatY bool  `xml:"repeaty,attr"`
	Image   Image `xml:"image"`
}

type Object struct {
	Identity
	Width      float64    `xml:"width,attr"`
	Height     float64    `xml:"height,attr"`
	X          float64    `xml:"x,attr"`
	Y          float64    `xml:"y,attr"`
	Rotation   float64    `xml:"rotation,attr"`
	Visible    string     `xml:"visible,attr"`
	Text       Text       `xml:"text"`
	Polygon    Polygon    `xml:"polygon"`
	Properties []Property `xml:"properties>property"`
}
type Polygon struct {
	Points string `xml:"points,attr"`
}

type Text struct {
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

func LoadTextureAtlas(textureId, atlasId string, cellWidth, cellHeight, cellGap int) string {
	var _, has = Textures[textureId]

	if has {
		var atlas = Atlas{TextureId: textureId, CellWidth: cellWidth, CellHeight: cellHeight, Gap: cellGap}
		Atlases[atlasId] = atlas
		return atlasId
	}

	return ""
}
