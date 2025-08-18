package gui

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/input/mouse"
	"pure-kit/engine/utility/seconds"
)

func Button(id string, properties ...string) string {
	return newWidget("button", id, properties...)
}

func (gui *GUI) ButtonClickedOnce(buttonId string, camera *graphics.Camera) bool {
	var widget, exists = gui.root.Widgets[buttonId]
	var owner = gui.root.Containers[widget.Owner]

	return exists && mouse.IsButtonReleasedOnce(mouse.ButtonLeft) &&
		widget.IsHovered(owner, camera) &&
		pressedOn == widget
}
func (gui *GUI) ButtonClickedAndHeld(buttonId string, camera *graphics.Camera) bool {
	var widget, exists = gui.root.Widgets[buttonId]
	if !exists {
		return false
	}

	var owner = gui.root.Containers[widget.Owner]
	var hover = widget.IsHovered(owner, camera)
	var first = condition.TrueOnce(hover && mouse.IsButtonPressedOnce(mouse.ButtonLeft), ";;first-"+buttonId)

	return first || (hover && pressedOn == widget &&
		condition.TrueEvery(0.1, ";;hold-"+buttonId) &&
		seconds.RealRuntime() > pressedAt+0.5)
}

// #region private

var pressedOn *widget
var pressedAt float32

func button(w, h float32, cam *graphics.Camera, root *root, widget *widget, owner *container) {
	var prev = widget.ThemeId
	var _, ownerDisabled = owner.Properties[property.Disabled]
	var _, disabled = widget.Properties[property.Disabled]
	var hover = themedProp(property.ButtonHoverThemeId, root, owner, widget)
	var press = themedProp(property.ButtonPressThemeId, root, owner, widget)

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

	visual(w, h, cam, root, widget, owner)
	widget.ThemeId = prev
}

// #endregion
