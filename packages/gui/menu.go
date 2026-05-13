package gui

import (
	"pure-game-kit/packages/execution/condition"
	"pure-game-kit/packages/gui/field"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/input/mouse"
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/text"
)

func Menu(id string, properties ...string) string {
	return newWidget("menu", id, properties...)
}

// private ========================================================

func menu(w *widget) {
	var owner = w.root.Containers[w.OwnerId]
	button(w)

	var escape = keyboard.IsKeyJustPressed(key.Escape)
	var anyButton = mouse.IsAnyButtonJustPressed() && !w.isHovered(owner)
	var containerId = w.root.themedField(field.MenuContainerId, owner, w)
	var c, has = w.root.Containers[containerId]

	if !has {
		return
	}

	var visible = c.Fields[field.Hidden] == ""
	if w.root.IsButtonJustClicked(w.Id) {
		c.Fields[field.Hidden] = condition.If(visible, "1", "")
		visible = !visible
	}
	c.Fields[field.X] = text.New(w.X)
	c.Fields[field.Y] = text.New(w.Y + w.Height)

	c.X = w.X
	c.Y = w.Y + w.Height

	var _, camH = w.root.cam.Size()
	var h = parseNum(w.root.themedField(field.Height, c, nil), 0)
	if c.Y+h > camH/2 {
		c.Fields[field.Y] = text.New(w.Y - h)
		c.Y = w.Y - h
	}

	if anyButton || mouse.Scroll() != 0 || !internal.WindowHovered || escape {
		if escape || (has && !c.isFocused()) {
			c.Fields[field.Hidden] = "1"
			visible = false
		}
	}

	if c.WasHidden && visible {
		// sound.AssetId = defaultValue(w.root.themedField(field.MenuSound, owner, w), "~popup")
		sound.Volume = w.root.Volume
		sound.Play()
	}

	c.WasHidden = !visible
}
