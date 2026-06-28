package internal

import "encoding/xml"

type GuiLayoutVars struct {
	Mx, My, Mw, Mh     float32
	Mlx, Mly, Mrx, Mry float32
	Mux, Muy, Mdx, Mdy float32
	Sx, Sy, Sw, Sh     float32
	Slx, Sly, Srx, Sry float32
	Sux, Suy, Sdx, Sdy float32
	Tx, Ty, Tw, Th     float32
	Tlx, Tly, Trx, Try float32
	Tux, Tuy, Tdx, Tdy float32
	Ow, Oh, Ov         float32
	Osx, Osy, Og, Mnr  float32
}
type GuiLayout struct {
	XMLName xml.Name `xml:"layout"`
	Boxes   []struct {
		Id                          uint32  `xml:"id,attr"`
		Name                        string  `xml:"name,attr"`
		NamePosition                string  `xml:"namePosition,attr"`
		Color                       string  `xml:"color,attr"`
		Visible                     int     `xml:"visible,attr"`
		Rectangle                   string  `xml:"rectangle,attr"`
		Math                        string  `xml:"math,attr"`
		Targets                     string  `xml:"target,attr"`
		ItemSize                    string  `xml:"itemSize,attr"`
		ItemSpacing                 string  `xml:"itemSpacing,attr"`
		ItemGap                     float32 `xml:"itemGap,attr"`
		ItemNewRow                  float32 `xml:"itemNewRow,attr"`
		ItemAlign                   string  `xml:"itemAlign,attr"`
		Vars                        GuiLayoutVars
		ItemStart, ItemEnd          int // cache on load
		ItemRangeCalculated         bool
		ContentWidth, ContentHeight float32
	} `xml:"boxes>box"`
	Items []struct {
		Id         uint32 `xml:"id,attr"`
		BoxId      uint32 `xml:"boxId,attr"`
		Name       string `xml:"name,attr"`
		Visible    int    `xml:"visible,attr"`
		Size       string `xml:"size,attr"`
		Expression string `xml:"math,attr"`
		NewRowMath string `xml:"newRowMath,attr"`
		Vars       GuiLayoutVars
	} `xml:"items>item"`
}

type GuiImage struct {
	ImageId     int     `xml:"imageId,attr"`
	Roundness   float32 `xml:"roundness,attr"`
	Color       string  `xml:"color,attr"`
	BorderSize  float32 `xml:"borderSize,attr"`
	BorderColor string  `xml:"borderColor,attr"`
}
type GuiText struct {
	FontId       int     `xml:"fontId,attr"`
	LineHeight   float32 `xml:"lineHeight,attr"`
	Gap          string  `xml:"gap,attr"`
	Margin       string  `xml:"margin,attr"`
	Align        string  `xml:"align,attr"`
	Weight       int8    `xml:"weight,attr"`
	Color        string  `xml:"color,attr"`
	OutlineSize  int8    `xml:"outlineSize,attr"`
	OutlineColor string  `xml:"outlineColor,attr"`
	ShadowWeight int8    `xml:"shadowWeight,attr"`
	ShadowColor  string  `xml:"shadowColor,attr"`
	ShadowBlur   uint8   `xml:"shadowBlur,attr"`
	ShadowOffset string  `xml:"shadowOffset,attr"`
}
type GuiTheme struct {
	XMLName xml.Name `xml:"theme"`
	Image   GuiImage `xml:"image"`
	Text    GuiText  `xml:"text"`
	Label   GuiText  `xml:"label"`
	Button  struct {
		Body struct {
			GuiImage
			Disabled GuiImage `xml:"disabled"`
			Focused  GuiImage `xml:"focused"`
			Clicked  GuiImage `xml:"clicked"`
		} `xml:"body"`
		Value struct {
			GuiText
			Disabled GuiText `xml:"disabled"`
			Focused  GuiText `xml:"focused"`
			Clicked  GuiText `xml:"clicked"`
		} `xml:"value"`
	} `xml:"button"`
	Scroll struct {
		Body   GuiImage `xml:"body"`
		Handle struct {
			GuiImage
			Focused GuiImage `xml:"focused"`
			Clicked GuiImage `xml:"clicked"`
		} `xml:"handle"`
	} `xml:"scroll"`
	Slider struct {
		Body struct {
			GuiImage
			Disabled GuiImage `xml:"disabled"`
			Focused  GuiImage `xml:"focused"`
			Clicked  GuiImage `xml:"clicked"`
		} `xml:"body"`
		Handle struct {
			GuiImage
			Disabled GuiImage `xml:"disabled"`
			Focused  GuiImage `xml:"focused"`
			Clicked  GuiImage `xml:"clicked"`
		} `xml:"handle"`
		Step GuiImage `xml:"step"`
	} `xml:"slider"`
	InputBox struct {
		Body struct {
			GuiImage
			Disabled GuiImage `xml:"disabled"`
			Typing   GuiImage `xml:"typing"`
		} `xml:"body"`
		Value struct {
			GuiText
			ShadowOffset string  `xml:"shadowOffset,attr"`
			Disabled     GuiText `xml:"disabled"`
			Typing       GuiText `xml:"typing"`
		} `xml:"value"`
		Placeholder GuiText  `xml:"placeholder"`
		Selection   GuiImage `xml:"selection"`
		Cursor      struct {
			GuiImage
			Width int `xml:"width,attr"`
		} `xml:"cursor"`
	} `xml:"inputbox"`
}

var Layouts = make(map[uint16]GuiLayout)
var Themes = make(map[uint16]GuiTheme)
var NextLayoutId, NextThemeId uint16
