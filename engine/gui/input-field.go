package gui

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/input/mouse"
)

func InputField(id string, properties ...string) string {
	return newWidget("inputField", id, properties...)
}

// #region private

func inputField(cam *graphics.Camera, root *root, widget *widget, owner *container) {
	setupVisualsText(root, widget, owner)
	setupVisualsTextured(root, widget, owner)
	drawVisuals(cam, root, widget, owner)

	if widget.IsFocused(root, cam) {
		mouse.SetCursor(mouse.CursorInput)
	}
}

// #endregion
