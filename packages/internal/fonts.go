package internal

import (
	_ "embed"
)

type Font struct {
	AtlasId int32 // see assets.ImageId
	Chars   map[rune]Glyph

	Ascender, Descender, LineHeight, EmSize, Size float32
}
type FontJSON struct {
	Atlas struct {
		Type                string  `json:"type"` // "psdf", "msdf" etc
		DistanceRange       float32 `json:"distanceRange"`
		DistanceRangeMiddle float32 `json:"distanceRangeMiddle"`
		Size                float32 `json:"size"`
		Width               int     `json:"width"`
		Height              int     `json:"height"`
		YOrigin             string  `json:"yOrigin"` // "bottom", "top" etc
	} `json:"atlas"`
	Metrics struct {
		EmSize             float32 `json:"emSize"`
		LineHeight         float32 `json:"lineHeight"`
		Ascender           float32 `json:"ascender"`
		Descender          float32 `json:"descender"`
		UnderlineY         float32 `json:"underlineY"`
		UnderlineThickness float32 `json:"underlineThickness"`
	} `json:"metrics"`
	Kernings []struct {
		Unicode1 rune    `json:"unicode1"`
		Unicode2 rune    `json:"unicode2"`
		Advance  float32 `json:"advance"`
	} `json:"kerning"`
	Glyphs []Glyph `json:"glyphs"`
}
type Bounds struct {
	Left   float32 `json:"left"`
	Bottom float32 `json:"bottom"`
	Right  float32 `json:"right"`
	Top    float32 `json:"top"`
}
type Glyph struct {
	Unicode        rune    `json:"unicode"`
	Advance        float32 `json:"advance"`
	PlaneBounds    Bounds  `json:"planeBounds"`
	AtlasBounds    Bounds  `json:"atlasBounds"`
	Kernings       map[rune]float32
	EmbededImageId int32
}

var Fonts = make(map[uint8]Font) // 0 = default
var FontNextId uint8

const Crossout, Underline = '\uE000', '\uE001'

func LoadFont(fontData *FontJSON, imageId int32) uint8 {
	var id = FontNextId
	var font = Font{AtlasId: imageId, Chars: make(map[rune]Glyph),
		Ascender: fontData.Metrics.Ascender, Descender: fontData.Metrics.Descender,
		LineHeight: fontData.Metrics.LineHeight, EmSize: fontData.Metrics.EmSize, Size: fontData.Atlas.Size,
	}

	for _, glyph := range fontData.Glyphs {
		glyph.Kernings = make(map[rune]float32)
		font.Chars[glyph.Unicode] = glyph
	}
	for _, kern := range fontData.Kernings {
		font.Chars[kern.Unicode1].Kernings[kern.Unicode2] = kern.Advance
	}

	var dash = font.Chars['-']
	if dash.Unicode != 0 {
		var center = (dash.AtlasBounds.Left + dash.AtlasBounds.Right) / 2
		dash.Unicode = Underline
		dash.AtlasBounds.Left, dash.AtlasBounds.Right = center-1, center+1
		font.Chars[Crossout] = dash
	}

	var underscore = font.Chars['_']
	if underscore.Unicode != 0 {
		var center = (underscore.AtlasBounds.Left + underscore.AtlasBounds.Right) / 2
		underscore.Unicode = Underline
		underscore.AtlasBounds.Left, underscore.AtlasBounds.Right = center-1, center+1
		font.Chars[Underline] = underscore
	}

	var space = font.Chars[' ']
	space.Advance = 0.35
	font.Chars[' '] = space

	Fonts[id] = font
	FontNextId++
	return id
}

// private ========================================================

//go:embed font.json.gz
var defaultFont []byte

//go:embed font.png
var defaultFontAtlas []byte

// ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz
// ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÑÒÓÔÕÖØÙÚÛÜÝÞßŒŠŽŁŃŚŹŻĆČĐŐŰàáâãäåæçèéêëìíîïñòóôõöøùúûüýþßœšžłńśźżćčđőűáéíóúüñãõçâêôÁÉÍÓÚÜÑÃÕÇÂÊÔßẞøØåÅþÞðÐœŒ
// 0123456789⁰¹²³⁴⁵⁶⁷⁸⁹₀₁₂₃₄₅₆₇₈₉¼½¾⅐⅑⅒⅓⅔⅕⅖⅗⅘⅙⅚⅛⅜⅝⅞
// \t\n .,;:!?¡¿\"'()[]{}<>-/\\@#$€£%^&*_+=|~`…•™§©®°–—‑′″‰ˆ˜“”‘’
// ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩαβγδεζηθικλμνξοπρστυφχψως
// АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюяҐЄІЇґєії
// აბგდევზთიკლმნოპჟრსტუფქღყშჩცძწჭხჯჰ
// ԱԲԳԴԵԶԷԸԹԺԻԼԽԾԿՀՁՂՃՄՅՆՇՈՉՊՋՌՍՎՏՐՑՒՓՔՕՖաբգդեզէըթժիլխծկհձղճմյնշոչպջռսվտրցւփքօֆ
// $€£₴₽₲₵₡₢₣₤₥₦₧₨₩₪₫₭₮₯₰₱₲₳₴₸₺₼₽¢
// ºª«»¶±×÷=≠<>≤≥∞∑∏√∫∆∂∇≈≡∈∉∪∩∧∨¬⇒⇔∀∃⊂⊆∅←↑→↓↔↕♠♥♦♣☺☹░▒▓│┤╡╢╖╕╣║╗╝┐└┴┬├─┼ˉ˙·
