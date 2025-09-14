package gui

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
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

	if mouse.IsAnyButtonPressedOnce() {
		print()
	}

	if mouse.IsAnyButtonPressedOnce() || mouse.Scroll() != 0 || !window.IsHovered() {
		var containerId = themedProp(property.MenuContainerId, root, owner, widget)
		var c, has = root.Containers[containerId]
		if has && !c.isFocused(cam) {
			c.Properties[property.Hidden] = "+"
		}
	}

	if root.ButtonClickedOnce(widget.Id, cam) {
		var containerId = themedProp(property.MenuContainerId, root, owner, widget)
		var c, has = root.Containers[containerId]
		if has {
			c.Properties[property.Hidden] = condition.If(c.Properties[property.Hidden] == "", "+", "")
			c.Properties[property.X] = text.New(widget.X)
			c.Properties[property.Y] = text.New(widget.Y + widget.Height)

			c.X = widget.X
			c.Y = widget.Y + widget.Height

			var _, camH = cam.Size()
			var h = parseNum(themedProp(property.Height, root, c, nil), 0)
			if c.Y+h > camH/2 {
				c.Properties[property.Y] = text.New(widget.Y - h)
				c.Y = widget.Y - h
			}
		}
	}
}
