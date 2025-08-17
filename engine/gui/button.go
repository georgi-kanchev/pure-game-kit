package gui

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/input/mouse"
)

func Button(id string, properties ...string) string {
	return newWidget("button", id, properties...)
}

// #region private

func button(w, h float32, cam *graphics.Camera, root *root, widget *widget, owner *container) {
	var prev = widget.ThemeId
	var hover = themedProp(property.ButtonHoverThemeId, root, owner, widget)
	var press = themedProp(property.ButtonPressThemeId, root, owner, widget)

	if widget.IsHovered(root, owner, cam) {
		mouse.SetCursor(mouse.CursorHand)

		if hover != "" {
			widget.ThemeId = hover
		}

		if press != "" && mouse.IsButtonPressed(mouse.ButtonLeft) {
			widget.ThemeId = press
		}
	}

	visual(w, h, cam, root, widget, owner)
	widget.ThemeId = prev
}

// #endregion
