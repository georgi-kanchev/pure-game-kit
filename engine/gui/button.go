package gui

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
)

func WidgetButton(id string, properties ...string) string {
	return newWidget("button", id, properties...)
}

func buttonUpdateAndDraw(cam *graphics.Camera, root *root, widget *widget, owner *container) {
	var x = parseNum(dyn(cam, owner, widget.Properties[property.X], "0"), 0)
	var y = parseNum(dyn(cam, owner, widget.Properties[property.Y], "0"), 0)
	var w = parseNum(dyn(cam, owner, themedProp(property.Width, root, owner, widget), "0"), 0)
	var h = parseNum(dyn(cam, owner, themedProp(property.Height, root, owner, widget), "0"), 0)
	var col = parseColor(themedProp(property.Color, root, owner, widget))
	cam.DrawRectangle(x, y, w, h, 0, col)
}
