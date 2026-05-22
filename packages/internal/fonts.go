package internal

import (
	_ "embed"
)

type Font struct {
	AtlasId int32 // see assets.ImageId
	Chars   map[rune]Glyph
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
	Kernings []struct {
		Unicode1 rune    `json:"unicode1"`
		Unicode2 rune    `json:"unicode2"`
		Advance  float64 `json:"advance"`
	} `json:"kerning"`
	Glyphs []Glyph `json:"glyphs"`
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
	Kernings    map[rune]float64
}

var Fonts = make(map[byte]Font) // 0 = default
var FontNextId byte

func LoadFont(fontData *FontJSON, imageId int32, isDefault bool) byte {
	if !isDefault {
		FontNextId++
	}
	var id = FontNextId
	var font = Font{AtlasId: imageId, Chars: make(map[rune]Glyph)}

	for _, glyph := range fontData.Glyphs {
		glyph.Kernings = make(map[rune]float64)
		font.Chars[glyph.Unicode] = glyph
	}
	for _, kern := range fontData.Kernings {
		font.Chars[kern.Unicode1].Kernings[kern.Unicode2] = kern.Advance
	}

	Fonts[id] = font
	return id
}
