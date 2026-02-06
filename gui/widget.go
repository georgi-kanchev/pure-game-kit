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
	UpdateAndDraw func(cam *graphics.Camera, root *root, widget *widget)

	// it is better to have many textbox instances instead of reusing a single one
	// because of its internal field caching for the symbols formatting
	textBox *graphics.TextBox
}

func newWidget(class, id string, properties ...string) string {
	return "<Widget " + field.Class + "=\"" + class + "\" " + field.Id + "=\"" + id + "\"" +
		extraProps(properties...) + " />"
}

//=================================================================

func (w *widget) isHovered(owner *container, cam *graphics.Camera) bool {
	return isHovered(owner.X, owner.Y, owner.Width, owner.Height, cam) &&
		isHovered(w.X, w.Y, w.Width, w.Height, cam)
}
func (w *widget) isFocused(root *root, cam *graphics.Camera) bool {
	return root.wFocused == w &&
		root.wWasHovered == w &&
		w.isHovered(root.Containers[w.OwnerId], cam)
}
func (w *widget) isDisabled(owner *container) bool {
	var disabled = w.Fields[field.Disabled] != ""
	var ownerDisabled = false

	if owner != nil {
		ownerDisabled = dyn(owner, owner.Fields[field.Disabled], "") != ""
	}

	return disabled || ownerDisabled
}
func (w *widget) isHidden(root *root, owner *container) bool {
	var toggleParentId = root.themedField(field.ToggleButtonId, owner, w)
	var toggleParent = root.Widgets[toggleParentId]
	var hidden = w.Fields[field.Hidden] != ""
	var parentHidden = toggleParent != nil && toggleParent.isHidden(root, owner)
	return hidden || parentHidden
}

func (w *widget) isSkipped(root *root, owner *container) bool {
	return w.isHidden(root, owner) || w.Class == "tooltip"
}
