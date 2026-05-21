package internal

import (
	_ "embed"
)

type Font struct {
	AtlasId  int32 // see assets.ImageId
	Chars    map[rune]Glyph
	Kernings map[rune]Kerning
}
type FontJSON struct {
	Atlas struct {
		Type                string  `json:"type"` // "psdf", "msdf" etc
		DistanceRange       float64 `json:"distanceRange"`
		DistanceRangeMiddle float64 `json:"distanceRangeMiddle"`
		Size                float64 `json:"size"`
		Width               int     `json:"width"`
		Height              int     `json:"height"`
		YOrigin             string  `json:"yOrigin"` // "bottom", "top" etc
	} `json:"atlas"`
	Metrics struct {
		EmSize             float64 `json:"emSize"`
		LineHeight         float64 `json:"lineHeight"`
		Ascender           float64 `json:"ascender"`
		Descender          float64 `json:"descender"`
		UnderlineY         float64 `json:"underlineY"`
		UnderlineThickness float64 `json:"underlineThickness"`
	} `json:"metrics"`
	Glyphs   []Glyph   `json:"glyphs"`
	Kernings []Kerning `json:"kerning"`
}
type Bounds struct {
	Left   float64 `json:"left"`
	Bottom float64 `json:"bottom"`
	Right  float64 `json:"right"`
	Top    float64 `json:"top"`
}
type Glyph struct {
	Unicode     rune    `json:"unicode"`
	Advance     float64 `json:"advance"`
	PlaneBounds Bounds  `json:"planeBounds"`
	AtlasBounds Bounds  `json:"atlasBounds"`
}
type Kerning struct {
	Unicode1 rune    `json:"unicode1"`
	Unicode2 rune    `json:"unicode2"`
	Advance  float64 `json:"advance"`
}

var Fonts = make(map[byte]Font) // 0 = default
var FontNextId byte

func LoadFont(fontData *FontJSON, imageId int32, isDefault bool) byte {
	if !isDefault {
		FontNextId++
	}
	var id = FontNextId
	var font = Font{AtlasId: imageId, Chars: make(map[rune]Glyph), Kernings: make(map[rune]Kerning)}

	for _, glyph := range fontData.Glyphs {
		font.Chars[glyph.Unicode] = glyph
	}
	for _, kern := range fontData.Kernings {
		font.Kernings[kern.Unicode1] = kern
	}

	Fonts[id] = font
	return id
}
