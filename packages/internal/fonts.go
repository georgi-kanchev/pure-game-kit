package internal

import (
	"encoding/xml"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Font struct {
	AtlasId  int32 // see assets.ImageId
	Chars    map[rune]Char
	Kernings map[rune]Kerning
}
type Char struct {
	ID       rune   `xml:"id,attr"`
	Index    int    `xml:"index,attr"`
	Char     string `xml:"char,attr"`
	Width    int    `xml:"width,attr"`
	Height   int    `xml:"height,attr"`
	XOffset  int    `xml:"xoffset,attr"`
	YOffset  int    `xml:"yoffset,attr"`
	XAdvance int    `xml:"xadvance,attr"`
	Chnl     int    `xml:"chnl,attr"`
	X        int    `xml:"x,attr"`
	Y        int    `xml:"y,attr"`
	Page     int    `xml:"page,attr"`
}
type Kerning struct {
	First  rune `xml:"first,attr"`
	Second rune `xml:"second,attr"`
	Amount int  `xml:"amount,attr"`
}
type FontData struct {
	XMLName xml.Name `xml:"font"`

	Info struct {
		Face     string `xml:"face,attr"`
		Size     int    `xml:"size,attr"`
		Bold     int    `xml:"bold,attr"`
		Italic   int    `xml:"italic,attr"`
		Charset  string `xml:"charset,attr"`
		Unicode  int    `xml:"unicode,attr"`
		StretchH int    `xml:"stretchH,attr"`
		Smooth   int    `xml:"smooth,attr"`
		AA       int    `xml:"aa,attr"`
		Padding  string `xml:"padding,attr"`
		Spacing  string `xml:"spacing,attr"`
		Outline  int    `xml:"outline,attr"`
	} `xml:"info"`

	Common struct {
		LineHeight int `xml:"lineHeight,attr"`
		Base       int `xml:"base,attr"`
		ScaleW     int `xml:"scaleW,attr"`
		ScaleH     int `xml:"scaleH,attr"`
		Pages      int `xml:"pages,attr"`
		Packed     int `xml:"packed,attr"`
		AlphaChnl  int `xml:"alphaChnl,attr"`
		RedChnl    int `xml:"redChnl,attr"`
		GreenChnl  int `xml:"greenChnl,attr"`
		BlueChnl   int `xml:"blueChnl,attr"`
	} `xml:"common"`

	Pages []struct {
		ID   int    `xml:"id,attr"`
		File string `xml:"file,attr"`
	} `xml:"pages>page"`

	DistanceField struct {
		FieldType     string `xml:"fieldType,attr"`
		DistanceRange int    `xml:"distanceRange,attr"`
	} `xml:"distanceField"`

	Chars struct {
		Count int    `xml:"count,attr"`
		Chars []Char `xml:"char"`
	} `xml:"chars"`

	Kernings struct {
		Count    int       `xml:"count,attr"`
		Kernings []Kerning `xml:"kerning"`
	} `xml:"kernings"`
}

var Fonts = make(map[string]rl.Font)

var Fonts2 = make(map[byte]Font) // 0 = default
var Font2NextId byte
