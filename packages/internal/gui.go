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

func (v *Vars) Lookup() func(string) float32 {
	return func(name string) float32 { // better for performance but not readability
		switch name {
		case "mx":
			return v.Mx
		case "my":
			return v.My
		case "mw":
			return v.Mw
		case "mh":
			return v.Mh
		case "mlx":
			return v.Mlx
		case "mly":
			return v.Mly
		case "mrx":
			return v.Mrx
		case "mry":
			return v.Mry
		case "mux":
			return v.Mux
		case "muy":
			return v.Muy
		case "mdx":
			return v.Mdx
		case "mdy":
			return v.Mdy
		case "sx":
			return v.Sx
		case "sy":
			return v.Sy
		case "sw":
			return v.Sw
		case "sh":
			return v.Sh
		case "slx":
			return v.Slx
		case "sly":
			return v.Sly
		case "srx":
			return v.Srx
		case "sry":
			return v.Sry
		case "sux":
			return v.Sux
		case "suy":
			return v.Suy
		case "sdx":
			return v.Sdx
		case "sdy":
			return v.Sdy
		case "tx":
			return v.Tx
		case "ty":
			return v.Ty
		case "tw":
			return v.Tw
		case "th":
			return v.Th
		case "tlx":
			return v.Tlx
		case "tly":
			return v.Tly
		case "trx":
			return v.Trx
		case "try":
			return v.Try
		case "tux":
			return v.Tux
		case "tuy":
			return v.Tuy
		case "tdx":
			return v.Tdx
		case "tdy":
			return v.Tdy
		case "ow":
			return v.Ow
		case "oh":
			return v.Oh
		case "ov":
			return v.Ov
		case "osx":
			return v.Osx
		case "osy":
			return v.Osy
		case "og":
			return v.Og
		case "mnr":
			return v.Mnr
		}
		return 0 // NaN
	}
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
