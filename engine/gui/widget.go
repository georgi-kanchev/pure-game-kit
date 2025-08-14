package gui

import (
	"encoding/xml"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
)

type widget struct {
	XmlProps []xml.Attr `xml:",any,attr"`

	Properties    map[string]string
	Owner         string
	UpdateAndDraw func(cam *graphics.Camera, root *root, widget *widget, owner *container)
}

func newWidget(class, id string, properties ...string) string {
	return "<Widget " + property.Class + "=\"" + class + "\" " + property.Id + "=\"" + id + "\"" +
		extraProps(properties...) + " />"
}
