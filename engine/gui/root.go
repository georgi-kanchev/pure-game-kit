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

	ContainerIds []string
	Themes       map[string]*theme
	Containers   map[string]*container
	Widgets      map[string]*widget
}

func (root *root) ButtonClickedOnce(buttonId string, camera *graphics.Camera) bool {
	var widget, exists = root.Widgets[buttonId]

	return exists && mouse.IsButtonReleasedOnce(mouse.ButtonLeft) &&
		widget.IsFocused(root, camera) && wPressedOn == widget
}
func (root *root) ButtonClickedAndHeld(buttonId string, camera *graphics.Camera) bool {
	var widget, exists = root.Widgets[buttonId]
	if !exists {
		return false
	}

	var hover = widget.IsFocused(root, camera)
	var first = condition.TrueOnce(hover && mouse.IsButtonPressedOnce(mouse.ButtonLeft), ";;first-"+buttonId)

	return first || (hover && wPressedOn == widget &&
		condition.TrueEvery(0.1, ";;hold-"+buttonId) &&
		seconds.RealRuntime() > wPressedAt+0.5)
}
