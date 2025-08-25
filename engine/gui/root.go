package gui

import (
	"encoding/xml"
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/input/mouse"
	"pure-kit/engine/utility/seconds"
)

type root struct {
	XmlName       xml.Name     `xml:"GUI"`
	XmlContainers []*container `xml:"Container"`

	Themes     map[string]*theme
	Containers map[string]*container
	Widgets    map[string]*widget
}

func (root *root) ButtonClickedOnce(buttonId string, camera *graphics.Camera) bool {
	var widget, exists = root.Widgets[buttonId]
	var owner = root.Containers[widget.OwnerId]

	return exists && mouse.IsButtonReleasedOnce(mouse.ButtonLeft) &&
		widget.IsHovered(owner, camera) &&
		pressedOn == widget
}
func (root *root) ButtonClickedAndHeld(buttonId string, camera *graphics.Camera) bool {
	var widget, exists = root.Widgets[buttonId]
	if !exists {
		return false
	}

	var owner = root.Containers[widget.OwnerId]
	var hover = widget.IsHovered(owner, camera)
	var first = condition.TrueOnce(hover && mouse.IsButtonPressedOnce(mouse.ButtonLeft), ";;first-"+buttonId)

	return first || (hover && pressedOn == widget &&
		condition.TrueEvery(0.1, ";;hold-"+buttonId) &&
		seconds.RealRuntime() > pressedAt+0.5)
}
