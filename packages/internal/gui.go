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
	ImgId  int     `xml:"imageId,attr"`
	Rnds   float32 `xml:"roundness,attr"`
	Col    string  `xml:"color,attr"`
	BorSz  float32 `xml:"borderSize,attr"`
	BorCol string  `xml:"borderColor,attr"`
}
type GuiText struct {
	FontId int     `xml:"fontId,attr"`
	LineH  float32 `xml:"lineHeight,attr"`
	Gap    string  `xml:"gap,attr"`
	Margin string  `xml:"margin,attr"`
	Align  string  `xml:"align,attr"`
	Wgt    float32 `xml:"weight,attr"`
	Col    string  `xml:"color,attr"`
	OutSz  float32 `xml:"outlineSize,attr"`
	OutCol string  `xml:"outlineColor,attr"`
	ShWgt  float32 `xml:"shadowWeight,attr"`
	ShCol  string  `xml:"shadowColor,attr"`
	ShBlur float32 `xml:"shadowBlur,attr"`
	ShOff  string  `xml:"shadowOffset,attr"`
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
		Body struct {
			GuiImage
			Size float32 `xml:"size,attr"`
		} `xml:"body"`
		Handle struct {
			GuiImage
			Speed   float32  `xml:"speed,attr"`
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
		Hnd struct {
			GuiImage
			Disabled GuiImage `xml:"disabled"`
			Focused  GuiImage `xml:"focused"`
			Clicked  GuiImage `xml:"clicked"`
		} `xml:"handle"`
		Step GuiImage `xml:"step"`
	} `xml:"slider"`
	Inputbox struct {
		Body struct {
			GuiImage
			Disabled GuiImage `xml:"disabled"`
			Focused  GuiImage `xml:"focused"`
			Typing   GuiImage `xml:"typing"`
		} `xml:"body"`
		Value struct {
			GuiText
			ShadowOffset string  `xml:"shadowOffset,attr"`
			Disabled     GuiText `xml:"disabled"`
			Focused      GuiText `xml:"focused"`
			Typing       GuiText `xml:"typing"`
		} `xml:"value"`
		Placeholder GuiText  `xml:"placeholder"`
		Selection   GuiImage `xml:"selection"`
		Cursor      struct {
			GuiImage
			Width float32 `xml:"width,attr"`
		} `xml:"cursor"`
	} `xml:"inputbox"`
}

var Layouts = make(map[uint16]GuiLayout)
var Themes = make(map[uint16]GuiTheme)
var NextLayoutId, NextThemeId uint16
