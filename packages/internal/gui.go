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
		Id           uint32  `xml:"id,attr"`
		Name         string  `xml:"name,attr"`
		NamePosition string  `xml:"namePos,attr"`
		Color        string  `xml:"col,attr"`
		Visible      int     `xml:"vis,attr"`
		Rectangle    string  `xml:"rect,attr"`
		Expression   string  `xml:"expr,attr"`
		Targets      string  `xml:"tar,attr"`
		ItemSize     string  `xml:"itSz,attr"`
		ItemSpacing  string  `xml:"itSp,attr"`
		ItemGap      float32 `xml:"itGap,attr"`
		ItemNewRow   float32 `xml:"itNewRow,attr"`
		ItemAlign    string  `xml:"itAl,attr"`
		Vars         Vars
	} `xml:"boxes>box"`
	Items []struct {
		Id               uint32  `xml:"id,attr"`
		BoxId            uint32  `xml:"boxId,attr"`
		Name             string  `xml:"name,attr"`
		Visible          int     `xml:"vis,attr"`
		Size             string  `xml:"size,attr"`
		Expression       string  `xml:"expr,attr"`
		NewRow           float32 `xml:"newRow,attr"`
		NewRowExpression string  `xml:"newRowExpr,attr"`
		Vars             Vars
	} `xml:"items>item"`
}

var Layouts map[uint32]Layout = make(map[uint32]Layout)
var NextLayoutId uint32
