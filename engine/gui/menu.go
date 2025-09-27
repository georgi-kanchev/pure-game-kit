package gui

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/field"
	"pure-kit/engine/input/keyboard"
	"pure-kit/engine/input/keyboard/key"
	"pure-kit/engine/input/mouse"
	"pure-kit/engine/utility/text"
	"pure-kit/engine/window"
)

func Menu(id string, properties ...string) string {
	return newWidget("menu", id, properties...)
}

//=================================================================
// private

func menu(cam *graphics.Camera, root *root, widget *widget) {
	var owner = root.Containers[widget.OwnerId]
	button(cam, root, widget)

	var escape = keyboard.IsKeyPressedOnce(key.Escape)
	if mouse.IsAnyButtonPressedOnce() || mouse.Scroll() != 0 || !window.IsHovered() || escape {
		var containerId = themedProp(field.MenuContainerId, root, owner, widget)
		var c, has = root.Containers[containerId]
		if escape || (has && !c.isFocused(cam)) {
			c.Properties[field.Hidden] = "+"
		}
	}

	if root.ButtonClickedOnce(widget.Id, cam) {
		var containerId = themedProp(field.MenuContainerId, root, owner, widget)
		var c, has = root.Containers[containerId]
		if !has {
			return
		}

		c.Properties[field.Hidden] = condition.If(c.Properties[field.Hidden] == "", "+", "")
		c.Properties[field.X] = text.New(widget.X)
		c.Properties[field.Y] = text.New(widget.Y + widget.Height)

		c.X = widget.X
		c.Y = widget.Y + widget.Height

		var _, camH = cam.Size()
		var h = parseNum(themedProp(field.Height, root, c, nil), 0)
		if c.Y+h > camH/2 {
			c.Properties[field.Y] = text.New(widget.Y - h)
			c.Y = widget.Y - h
		}
	}
}
