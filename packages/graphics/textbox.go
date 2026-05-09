package graphics

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/internal"
	txt "pure-game-kit/packages/utility/text"
)

type TextSymbol struct {
	Sprite
	Value rune

	Underline, Crossout bool

	Weight, OutlineWeight,
	ShadowWeight, ShadowBlur,
	ShadowOffsetX, ShadowOffsetY float32

	BackColor, OutlineColor, ShadowColor uint
}

type Textbox struct {
	Quad

	Text string
	Font assets.FontId
	AlignX, AlignY,
	LineHeight, SymbolGap, LineGap float32
	WordWrap bool

	//=================================================================

	cache     textboxCache // tracks when to regenerate chars
	chars     []TextSymbol
	lineCount int
}

func NewTextbox(font assets.FontId, x, y float32, text ...any) *Textbox {
	var textbox = &Textbox{Font: font, Text: txt.New(text...), Quad: *NewQuad(x, y), LineHeight: 100, WordWrap: true}
	textbox.tryRegenerate()
	return textbox
}

//=================================================================

func (t *Textbox) Measure(text string) (width, height float32) {
	t.tryRegenerate()
	return 0, 0
}
func (t *Textbox) LineCount() int {
	t.tryRegenerate()
	return t.lineCount
}
func (t *Textbox) Symbol(index int) *TextSymbol {
	t.tryRegenerate()
	return nil
}

// private ========================================================

type textboxCache struct {
	text string
	font assets.FontId
	width, height, alignX, alignY,
	lineHeight, symbolGap, lineGap float32
	wordWrap bool
}

func (t *Textbox) tryRegenerate() {
	var w, h, ax, ay = t.Width, t.Height, t.AlignX, t.AlignY
	var state = textboxCache{t.Text, t.Font, w, h, ax, ay, t.LineHeight, t.SymbolGap, t.LineGap, t.WordWrap}
	if state == t.cache {
		return
	}
	t.cache = state

	t.chars = t.chars[:]
	var fontData = internal.Fonts2[byte(t.Font)]
	for i, r := range t.Text {
		var sprite = *NewSprite(assets.ImageId(fontData.AtlasId), float32(30*i), 0)
		var symbol = TextSymbol{Sprite: sprite, Value: r}
		t.chars = append(t.chars, symbol)
	}
}
