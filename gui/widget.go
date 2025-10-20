package gui

import (
	"encoding/xml"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
)

type widget struct {
	XmlProps []xml.Attr `xml:",any,attr"`

	Id, Class, OwnerId, ThemeId string
	X, Y, Width, Height,
	DragX, DragY, PrevValue float32
	IsCulled      bool
	Properties    map[string]string
	UpdateAndDraw func(cam *graphics.Camera, root *root, widget *widget)
}

func newWidget(class, id string, properties ...string) string {
	return "<Widget " + field.Class + "=\"" + class + "\" " + field.Id + "=\"" + id + "\"" +
		extraProps(properties...) + " />"
}

//=================================================================

func (widget *widget) isHovered(owner *container, cam *graphics.Camera) bool {
	return isHovered(owner.X, owner.Y, owner.Width, owner.Height, cam) &&
		isHovered(widget.X, widget.Y, widget.Width, widget.Height, cam)
}
func (widget *widget) isFocused(root *root, cam *graphics.Camera) bool {
	return wFocused == widget && wWasHovered == widget && widget.isHovered(root.Containers[widget.OwnerId], cam)
}
func (widget *widget) isDisabled(owner *container) bool {
	var _, disabled = widget.Properties[field.Disabled]
	var ownerDisabled = false

	if owner != nil {
		_, ownerDisabled = owner.Properties[field.Disabled]
	}

	return disabled || ownerDisabled
}
