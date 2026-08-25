package internal

import "encoding/xml"

type GUILayoutVars struct {
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
type GUILayout struct {
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
		Vars                        GUILayoutVars
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
		Vars       GUILayoutVars
	} `xml:"items>item"`
}

type GUIImage struct {
	ImgId  *int     `xml:"imageId,attr"`
	Rnds   *float32 `xml:"roundness,attr"`
	Col    *string  `xml:"color,attr"`
	BorSz  *float32 `xml:"borderSize,attr"`
	BorCol *string  `xml:"borderColor,attr"`
}
type GUIText struct {
	FontId *int     `xml:"fontId,attr"`
	LineH  *float32 `xml:"lineHeight,attr"`
	Gap    *string  `xml:"gap,attr"`
	Margin *string  `xml:"margin,attr"`
	Align  *string  `xml:"align,attr"`
	Wgt    *float32 `xml:"weight,attr"`
	Col    *string  `xml:"color,attr"`
	OutSz  *float32 `xml:"outlineSize,attr"`
	OutCol *string  `xml:"outlineColor,attr"`
	ShWgt  *float32 `xml:"shadowWeight,attr"`
	ShCol  *string  `xml:"shadowColor,attr"`
	ShBlur *float32 `xml:"shadowBlur,attr"`
	ShOff  *string  `xml:"shadowOffset,attr"`
}
type GUITheme struct {
	XMLName xml.Name `xml:"theme"`
	Image   GUIImage `xml:"image"`
	Text    GUIText  `xml:"text"`
	Label   GUIText  `xml:"label"`
	Button  struct {
		Body struct {
			GUIImage
			Disabled GUIImage `xml:"disabled"`
			Focused  GUIImage `xml:"focused"`
			Clicked  GUIImage `xml:"clicked"`
		} `xml:"body"`
		Value struct {
			GUIText
			Disabled GUIText `xml:"disabled"`
			Focused  GUIText `xml:"focused"`
			Clicked  GUIText `xml:"clicked"`
		} `xml:"value"`
	} `xml:"button"`
	Scroll struct {
		Body struct {
			GUIImage
			Size *float32 `xml:"size,attr"`
		} `xml:"body"`
		Handle struct {
			GUIImage
			Speed   *float32 `xml:"speed,attr"`
			Focused GUIImage `xml:"focused"`
			Clicked GUIImage `xml:"clicked"`
		} `xml:"handle"`
	} `xml:"scroll"`
	Slider struct {
		Body struct {
			GUIImage
			Disabled GUIImage `xml:"disabled"`
			Focused  GUIImage `xml:"focused"`
			Clicked  GUIImage `xml:"clicked"`
		} `xml:"body"`
		Hnd struct {
			GUIImage
			Disabled GUIImage `xml:"disabled"`
			Focused  GUIImage `xml:"focused"`
			Clicked  GUIImage `xml:"clicked"`
		} `xml:"handle"`
		Step GUIImage `xml:"step"`
	} `xml:"slider"`
	Inputbox struct {
		Body struct {
			GUIImage
			Disabled GUIImage `xml:"disabled"`
			Focused  GUIImage `xml:"focused"`
			Typing   GUIImage `xml:"typing"`
		} `xml:"body"`
		Value struct {
			GUIText
			Disabled GUIText `xml:"disabled"`
			Focused  GUIText `xml:"focused"`
			Typing   GUIText `xml:"typing"`
		} `xml:"value"`
		Placeholder GUIText  `xml:"placeholder"`
		Selection   GUIImage `xml:"selection"`
		Cursor      struct {
			GUIImage
			Width *float32 `xml:"width,attr"`
		} `xml:"cursor"`
	} `xml:"inputbox"`
}

var GUILayouts = make(map[uint16]GUILayout)
var GUIThemes = make(map[uint16]GUITheme)
var NextGUILayoutId, NextGUIThemeId uint16
