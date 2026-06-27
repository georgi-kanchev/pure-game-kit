package internal

import "encoding/xml"

type Vars struct {
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

type Layout struct {
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
		Vars                        Vars
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
		Vars       Vars
	} `xml:"items>item"`
}

type Image struct {
	ImageId     int     `xml:"imageId,attr"`
	Roundness   float32 `xml:"roundness,attr"`
	Color       string  `xml:"color,attr"`
	BorderSize  float32 `xml:"borderSize,attr"`
	BorderColor string  `xml:"borderColor,attr"`
}
type Text struct {
	FontId        int     `xml:"fontId,attr"`
	LineHeight    float32 `xml:"lineHeight,attr"`
	Gap           string  `xml:"gap,attr"`
	Margin        string  `xml:"margin,attr"`
	Align         string  `xml:"align,attr"`
	Weight        int8    `xml:"weight,attr"`
	FillColor     string  `xml:"fillColor,attr"`
	OutlineWeight int8    `xml:"outlineWeight,attr"`
	OutlineColor  string  `xml:"outlineColor,attr"`
	ShadowWeight  int8    `xml:"shadowWeight,attr"`
	ShadowColor   string  `xml:"shadowColor,attr"`
	ShadowBlur    int8    `xml:"shadowBlur,attr"`
	ShadowOffset  string  `xml:"shadowOffset,attr"`
}
type Theme struct {
	XMLName xml.Name `xml:"theme"`
	Image   Image    `xml:"image"`
	Text    Text     `xml:"text"`
	Label   Text     `xml:"label"`
	Button  struct {
		Body struct {
			Image
			Disabled Image `xml:"disabled"`
			Focused  Image `xml:"focused"`
			Clicked  Image `xml:"clicked"`
		} `xml:"body"`
		Value struct {
			Text
			Disabled Text `xml:"disabled"`
			Focused  Text `xml:"focused"`
			Clicked  Text `xml:"clicked"`
		} `xml:"value"`
	} `xml:"button"`
	Scroll struct {
		Body   Image `xml:"body"`
		Handle struct {
			Image
			Focused Image `xml:"focused"`
			Clicked Image `xml:"clicked"`
		} `xml:"handle"`
	} `xml:"scroll"`
	Slider struct {
		Body struct {
			Image
			Disabled Image `xml:"disabled"`
			Focused  Image `xml:"focused"`
			Clicked  Image `xml:"clicked"`
		} `xml:"body"`
		Handle struct {
			Image
			Disabled Image `xml:"disabled"`
			Focused  Image `xml:"focused"`
			Clicked  Image `xml:"clicked"`
		} `xml:"handle"`
		Step Image `xml:"step"`
	} `xml:"slider"`
	InputBox struct {
		Body struct {
			Image
			Disabled Image `xml:"disabled"`
			Typing   Image `xml:"typing"`
		} `xml:"body"`
		Value struct {
			Text
			ShadowOffset string `xml:"shadowOffset,attr"`
			Disabled     Text   `xml:"disabled"`
			Typing       Text   `xml:"typing"`
		} `xml:"value"`
		Placeholder Text  `xml:"placeholder"`
		Selection   Image `xml:"selection"`
		Cursor      struct {
			Image
			Width int `xml:"width,attr"`
		} `xml:"cursor"`
	} `xml:"inputbox"`
}

var Layouts map[uint32]Layout = make(map[uint32]Layout)
var Themes map[uint32]Theme = make(map[uint32]Theme)
var NextLayoutId, NextThemeId uint32
