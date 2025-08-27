package gui

import (
	"encoding/xml"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
)

type widget struct {
	XmlProps []xml.Attr `xml:",any,attr"`

	Id, Class, OwnerId, ThemeId string
	X, Y, Width, Height         float32
	Properties                  map[string]string
	UpdateAndDraw               func(cam *graphics.Camera, root *root, widget *widget, owner *container)
}

func newWidget(class, id string, properties ...string) string {
	return "<Widget " + property.Class + "=\"" + class + "\" " + property.Id + "=\"" + id + "\"" +
		extraProps(properties...) + " />"
}

func (widget *widget) IsHovered(owner *container, cam *graphics.Camera) bool {
	return isHovered(owner.X, owner.Y, owner.Width, owner.Height, cam) &&
		isHovered(widget.X, widget.Y, widget.Width, widget.Height, cam)
}
func (widget *widget) IsFocused(root *root, cam *graphics.Camera) bool {
	return focused == widget && wasHovered == widget && widget.IsHovered(root.Containers[widget.OwnerId], cam)
}

func (widget *widget) IsDisabled(owner *container) bool {
	var _, disabled = widget.Properties[property.Disabled]
	var ownerDisabled = false

	if owner != nil {
		_, ownerDisabled = owner.Properties[property.Disabled]
	}

	return disabled || ownerDisabled
}
