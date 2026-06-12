package internal

import "encoding/xml"

type Layout struct {
	XMLName xml.Name `xml:"layout"`
	Boxes   []struct {
		Id           uint32 `xml:"id,attr"`
		Name         string `xml:"name,attr"`
		NamePosition string `xml:"namePos,attr"`
		Color        string `xml:"col,attr"`
		Visible      int    `xml:"vis,attr"`
		Rectangle    string `xml:"rect,attr"`
		Expression   string `xml:"expr,attr"`
		Targets      string `xml:"tar,attr"`
		ItemSize     string `xml:"itSz,attr"`
		ItemSpacing  string `xml:"itSp,attr"`
		ItemGap      int    `xml:"itGap,attr"`
		ItemNewRow   int    `xml:"itNewRow,attr"`
		ItemAlign    string `xml:"itAl,attr"`
		Vars         map[string]float32
	} `xml:"boxes>box"`
	Items []struct {
		Id               uint32 `xml:"id,attr"`
		BoxId            uint32 `xml:"boxId,attr"`
		Name             string `xml:"name,attr"`
		Visible          int    `xml:"vis,attr"`
		Rectangle        string `xml:"rect,attr"`
		Expression       string `xml:"expr,attr"`
		NewRow           int    `xml:"newRow,attr"`
		NewRowExpression string `xml:"newRowExpr,attr"`
		Variables        map[string]float32
	} `xml:"items>item"`
}

var Layouts map[uint32]Layout = make(map[uint32]Layout)
var NextLayoutId uint32
