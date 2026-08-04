package internal

import "encoding/xml"

type AtlasData struct {
	XMLName xml.Name `xml:"data"`
	Crops   []struct {
		X int `xml:"x,attr"`
		Y int `xml:"y,attr"`
		W int `xml:"w,attr"`
		H int `xml:"h,attr"`
	} `xml:"crops>crop"`
	Groups []struct {
		Name        string `xml:"name,attr"`
		CropIndexes string `xml:"cropIndexes,attr"` // space-separated list of integers
	} `xml:"groups>group"`

	Map map[string][]int32
}

var Atlases = make(map[uint16]AtlasData)
var NextAtlasId uint16
