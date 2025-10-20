package gui

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/input/mouse"
	"pure-game-kit/utility/text"
	"pure-game-kit/window"
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
	var anyButton = mouse.IsAnyButtonPressedOnce() && !widget.isHovered(owner, cam)
	var containerId = themedProp(field.MenuContainerId, root, owner, widget)
	var c, has = root.Containers[containerId]
	var visible = c.Properties[field.Hidden] == ""

	if has && root.IsButtonClickedOnce(widget.Id, cam) {
		c.Properties[field.Hidden] = condition.If(visible, "+", "")
		c.Properties[field.X] = text.New(widget.X)
		c.Properties[field.Y] = text.New(widget.Y + widget.Height)
		visible = !visible

		c.X = widget.X
		c.Y = widget.Y + widget.Height

		var _, camH = cam.Size()
		var h = parseNum(themedProp(field.Height, root, c, nil), 0)
		if c.Y+h > camH/2 {
			c.Properties[field.Y] = text.New(widget.Y - h)
			c.Y = widget.Y - h
		}
	}

	if anyButton || mouse.Scroll() != 0 || !window.IsHovered() || escape {
		if escape || (has && !c.isFocused(cam)) {
			c.Properties[field.Hidden] = "+"
			visible = false
		}
	}

	if c.WasHidden && visible {
		sound.AssetId = defaultValue(themedProp(field.MenuSound, root, owner, widget), "~popup")
		sound.Volume = root.Volume
		sound.Play()
	}

	c.WasHidden = !visible
}
