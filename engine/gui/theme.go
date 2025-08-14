package gui

import (
	"encoding/xml"
	"pure-kit/engine/gui/property"
)

type theme struct {
	XmlProps []xml.Attr `xml:",any,attr"`

	Properties map[string]string
}

func Theme(id string, properties ...string) string {
	return "<Theme " + property.Id + "=\"" + id + "\"" + extraProps(properties...) + " />"
}
