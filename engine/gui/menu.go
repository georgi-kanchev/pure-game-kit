package gui

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/input/mouse"
	"pure-kit/engine/utility/symbols"
)

func Menu(id string, properties ...string) string {
	return newWidget("menu", id, properties...)
}

func menu(cam *graphics.Camera, root *root, widget *widget, owner *container) {
	button(cam, root, widget, owner)

	var scrolledOutside = mouse.Scroll() != 0 && !widget.IsFocused(root, cam)
	// var draggedOutside = mouse.IsButtonPressed(mouse.ButtonMiddle) && cFocused != owner
	if mouse.IsButtonPressedOnce(mouse.ButtonLeft) ||
		mouse.IsButtonPressedOnce(mouse.ButtonMiddle) || scrolledOutside ||
		mouse.IsButtonPressedOnce(mouse.ButtonRight) {
		var containerId = themedProp(property.MenuContainerId, root, owner, widget)
		var c, has = root.Containers[containerId]
		if has && !c.IsFocused(root, cam) && !widget.IsFocused(root, cam) {
			c.Properties[property.Hidden] = "+"
		}
	}

	if root.ButtonClickedOnce(widget.Id, cam) {
		var containerId = themedProp(property.MenuContainerId, root, owner, widget)
		var c, has = root.Containers[containerId]
		if has {
			c.Properties[property.Hidden] = condition.If(c.Properties[property.Hidden] == "", "+", "")
			c.Properties[property.X] = symbols.New(widget.X)
			c.Properties[property.Y] = symbols.New(widget.Y + widget.Height)
		}
	}
}
