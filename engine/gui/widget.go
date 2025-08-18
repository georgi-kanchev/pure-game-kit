package gui

import (
	"encoding/xml"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
)

type widget struct {
	XmlProps []xml.Attr `xml:",any,attr"`

	Class, Owner, ThemeId string
	X, Y, Width, Height   float32
	Properties            map[string]string
	UpdateAndDraw         func(w, h float32, cam *graphics.Camera, root *root, widget *widget, owner *container)
}

func newWidget(class, id string, properties ...string) string {
	return "<Widget " + property.Class + "=\"" + class + "\" " + property.Id + "=\"" + id + "\"" +
		extraProps(properties...) + " />"
}

func (widget *widget) IsHovered(owner *container, cam *graphics.Camera) bool {
	return isHovered(owner.X, owner.Y, owner.Width, owner.Height, cam) &&
		isHovered(widget.X, widget.Y, widget.Width, widget.Height, cam)
}
