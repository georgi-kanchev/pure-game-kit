package graphics

import txt "pure-game-kit/packages/utility/text"

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

	Text, FontId string
	AlignX, AlignY,
	LineHeight, SymbolGap, LineGap float32
	WordWrap bool

	//=================================================================

	cache    textboxCache // tracks when to regenerate chars
	chars    []TextSymbol
	newLines int
}

func NewTextbox(fontId string, x, y float32, text ...any) Textbox {
	var textbox = Textbox{FontId: fontId, Text: txt.New(text...), Quad: *NewQuad(x, y), LineHeight: 100, WordWrap: true}
	textbox.tryRegenerate()
	return textbox
}

//=================================================================

func (t *Textbox) Measure(text string) (width, height float32) {
	t.tryRegenerate()
	return 0, 0
}
func (t *Textbox) NewLines() int {
	t.tryRegenerate()
	return t.newLines
}
func (t *Textbox) Symbol(index int) *TextSymbol {
	t.tryRegenerate()
	return nil
}

// private ========================================================

type textboxCache struct {
	text, fontId string
	width, height, alignX, alignY,
	lineHeight, symbolGap, lineGap float32
	wordWrap bool
}

func (t *Textbox) tryRegenerate() {
	var w, h, ax, ay = t.Width, t.Height, t.AlignX, t.AlignY
	var state = textboxCache{t.Text, t.FontId, w, h, ax, ay, t.LineHeight, t.SymbolGap, t.LineGap, t.WordWrap}
	if state == t.cache {
		return
	}
	t.cache = state
}
