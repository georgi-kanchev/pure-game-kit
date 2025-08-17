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

func (widget *widget) IsHovered(root *root, owner *container, cam *graphics.Camera) bool {
	if !owner.IsHovered(root, cam) {
		return false
	}

	var x, y = cam.PointToScreen(widget.X, widget.Y)
	var mx, my = cam.PointToScreen(cam.MousePosition())
	return mx > x && mx < x+int(widget.Width) && my > y && my < y+int(widget.Height)
}
