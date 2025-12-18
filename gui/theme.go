package gui

import (
	"encoding/xml"
	"pure-game-kit/gui/field"
)

type theme struct {
	XmlProps []xml.Attr `xml:",any,attr"`

	Fields map[string]string
}

func Theme(id string, properties ...string) string {
	return "<Theme " + field.Id + "=\"" + id + "\"" + extraProps(properties...) + " />"
}
