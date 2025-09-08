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

//=================================================================
// getters

func (gui *GUI) ButtonOnClickOnce(buttonId string, camera *graphics.Camera) bool {
	return gui.root.ButtonClickedOnce(buttonId, camera)
}
func (gui *GUI) ButtonOnClickAndHold(buttonId string, camera *graphics.Camera) bool {
	return gui.root.ButtonClickedAndHeld(buttonId, camera)
}

//=================================================================
// private

var wPressedOn *widget
var wPressedAt float32
var buttonColor uint

func button(cam *graphics.Camera, root *root, widget *widget) {
	var owner = root.Containers[widget.OwnerId]
	var prev = widget.ThemeId
	var _, ownerDisabled = owner.Properties[property.Disabled]
	var _, disabled = widget.Properties[property.Disabled]
	var hover = themedProp(property.ButtonThemeIdHover, root, owner, widget)
	var press = themedProp(property.ButtonThemeIdPress, root, owner, widget)

	if widget.isFocused(root, cam) {
		mouse.SetCursor(mouse.CursorHand)

		if disabled || ownerDisabled {
			mouse.SetCursor(mouse.CursorNotAllowed)
		} else {
			if hover != "" {
				widget.ThemeId = hover
			}
			if press != "" && wPressedOn == widget && mouse.IsButtonPressed(mouse.ButtonLeft) {
				widget.ThemeId = press
			}
			if mouse.IsButtonPressedOnce(mouse.ButtonLeft) {
				wPressedOn = widget
				wPressedAt = seconds.RealRuntime()
			}
		}
	}

	setupVisualsTextured(root, widget)
	setupVisualsText(root, widget)
	drawVisuals(cam, root, widget)
	buttonColor = parseColor(themedProp(property.Color, root, owner, widget), widget.isDisabled(owner))
	widget.ThemeId = prev
}
