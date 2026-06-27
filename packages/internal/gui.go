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

type Theme struct {
	XMLName xml.Name `xml:"theme"`
}

var Layouts map[uint32]*Layout = make(map[uint32]*Layout)
var NextLayoutId uint32
