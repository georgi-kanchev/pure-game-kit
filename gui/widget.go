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
	Fields        map[string]string
	UpdateAndDraw func(widget *widget)

	root *root

	sprite, handle, cursor1, cursor2,
	top, left, right, bottom *graphics.Sprite
	steps          []*graphics.Sprite
	textBox        *graphics.TextBox
	box, highlight *graphics.Box
}

func newWidget(class, id string, properties ...string) string {
	return "<Widget " + field.Class + "=\"" + class + "\" " + field.Id + "=\"" + id + "\"" +
		extraProps(properties...) + " />"
}

//=================================================================

func (w *widget) isHovered(owner *container) bool {
	return isHovered(owner.X, owner.Y, owner.Width, owner.Height, w.root.cam) &&
		isHovered(w.X, w.Y, w.Width, w.Height, w.root.cam)
}
func (w *widget) isFocused() bool {
	return w.root.wFocused == w &&
		w.root.wWasHovered == w &&
		w.isHovered(w.root.Containers[w.OwnerId])
}
func (w *widget) isDisabled(owner *container) bool {
	var disabled = w.Fields[field.Disabled] != ""
	var ownerDisabled = false

	if owner != nil {
		ownerDisabled = dyn(owner, owner.Fields[field.Disabled], "") != ""
	}

	return disabled || ownerDisabled
}
func (w *widget) isHidden(owner *container) bool {
	var toggleParentId = w.root.themedField(field.ToggleButtonId, owner, w)
	var toggleParent = w.root.Widgets[toggleParentId]
	var hidden = w.Fields[field.Hidden] != ""
	var parentHidden = toggleParent != nil && toggleParent.isHidden(owner)
	return hidden || parentHidden
}
func (w *widget) isSkipped(owner *container) bool {
	return w.isHidden(owner) || w.Class == "tooltip"
}
