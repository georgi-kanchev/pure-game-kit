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

	var escape = keyboard.IsKeyJustPressed(key.Escape)
	var anyButton = mouse.IsAnyButtonJustPressed() && !widget.isHovered(owner, cam)
	var containerId = root.themedField(field.MenuContainerId, owner, widget)
	var c, has = root.Containers[containerId]
	var visible = c.Fields[field.Hidden] == ""

	if has && root.IsButtonJustClicked(widget.Id, cam) {
		c.Fields[field.Hidden] = condition.If(visible, "+", "")
		c.Fields[field.X] = text.New(widget.X)
		c.Fields[field.Y] = text.New(widget.Y + widget.Height)
		visible = !visible

		c.X = widget.X
		c.Y = widget.Y + widget.Height

		var _, camH = cam.Size()
		var h = parseNum(root.themedField(field.Height, c, nil), 0)
		if c.Y+h > camH/2 {
			c.Fields[field.Y] = text.New(widget.Y - h)
			c.Y = widget.Y - h
		}
	}

	if anyButton || mouse.Scroll() != 0 || !window.IsHovered() || escape {
		if escape || (has && !c.isFocused(cam)) {
			c.Fields[field.Hidden] = "+"
			visible = false
		}
	}

	if c.WasHidden && visible {
		sound.AssetId = defaultValue(root.themedField(field.MenuSound, owner, widget), "~popup")
		sound.Volume = root.Volume
		sound.Play()
	}

	c.WasHidden = !visible
}
