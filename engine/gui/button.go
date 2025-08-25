package gui

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/input/mouse"
	"pure-kit/engine/utility/seconds"
)

func Button(id string, properties ...string) string {
	return newWidget("button", id, properties...)
}

func (gui *GUI) ButtonClickedOnce(buttonId string, camera *graphics.Camera) bool {
	return gui.root.ButtonClickedOnce(buttonId, camera)
}
func (gui *GUI) ButtonClickedAndHeld(buttonId string, camera *graphics.Camera) bool {
	return gui.root.ButtonClickedAndHeld(buttonId, camera)
}

// #region private

var pressedOn *widget
var pressedAt float32
var buttonColor uint

func button(cam *graphics.Camera, root *root, widget *widget, owner *container) {
	var prev = widget.ThemeId
	var _, ownerDisabled = owner.Properties[property.Disabled]
	var _, disabled = widget.Properties[property.Disabled]
	var hover = themedProp(property.ButtonThemeIdHover, root, owner, widget)
	var press = themedProp(property.ButtonThemeIdPress, root, owner, widget)

	if widget.IsHovered(owner, cam) {
		mouse.SetCursor(mouse.CursorHand)

		if disabled || ownerDisabled {
			mouse.SetCursor(mouse.CursorNotAllowed)
		} else {
			if hover != "" {
				widget.ThemeId = hover
			}
			if press != "" && pressedOn == widget && mouse.IsButtonPressed(mouse.ButtonLeft) {
				widget.ThemeId = press
			}
			if mouse.IsButtonPressedOnce(mouse.ButtonLeft) {
				pressedOn = widget
				pressedAt = seconds.RealRuntime()
			}
		}
	}

	setupVisualsTextured(root, widget, owner)
	setupVisualsText(root, widget, owner)
	drawVisuals(cam, root, widget, owner)
	buttonColor = parseColor(themedProp(property.Color, root, owner, widget), widget.IsDisabled(owner))
	widget.ThemeId = prev
}

// #endregion
