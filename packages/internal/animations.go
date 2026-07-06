package internal

import "encoding/xml"

type AnimationsData struct {
	XMLName xml.Name `xml:"data"`
	Frames  []struct {
		X int `xml:"x,attr"`
		Y int `xml:"y,attr"`
		W int `xml:"w,attr"`
		H int `xml:"h,attr"`
	} `xml:"frames>frame"`
	Animations []struct {
		Name   string `xml:"name,attr"`
		Frames string `xml:"frames,attr"` // space-separated list of integers
	} `xml:"animations>animation"`

	Map map[string][]int32
}

var Animations = make(map[uint16]AnimationsData)
var NextAnimationsId uint16
