package gui

import (
	"encoding/xml"
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
	k "pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	m "pure-game-kit/input/mouse"
	b "pure-game-kit/input/mouse/button"
	"pure-game-kit/internal"
	"pure-game-kit/utility/text"
)

type root struct {
	XmlName       xml.Name     `xml:"GUI"`
	XmlContainers []*container `xml:"Container"`
	XmlScale      float32      `xml:"scale,attr"`
	XmlVolume     float32      `xml:"volume,attr"`

	Volume       float32
	ContainerIds []string
	Themes       map[string]*theme
	Containers   map[string]*container
	Widgets      map[string]*widget

	wHovered, wWasHovered, wFocused, wPressedOn                 *widget
	cHovered, cWasHovered, cFocused, cScrolledOn, cWasScrolling *container
	wPressedAt                                                  float32

	cMiddlePressed, cPressedOnScrollH, cPressedOnScrollV *container // for container slider

	cam *graphics.Camera
}

func (r *root) IsButtonJustClicked(buttonId string) bool {
	var widget, exists = r.Widgets[buttonId]
	if !exists {
		return false
	}

	var owner = r.Containers[widget.OwnerId]
	var hotkey = key.FromName(r.themedField(field.ButtonHotkey, owner, widget))
	var focus = widget.isFocused() && r.wPressedOn == widget
	var input = k.IsKeyJustPressed(hotkey) || (focus && m.IsButtonJustReleased(b.Left))
	return input
}
func (r *root) IsButtonClickedAndHeld(buttonId string) bool {
	var widget, exists = r.Widgets[buttonId]
	if !exists {
		return false
	}

	var focus = widget.isFocused()
	var owner = r.Containers[widget.OwnerId]
	var hotkey = key.FromName(r.themedField(field.ButtonHotkey, owner, widget))
	var first = k.IsKeyJustPressed(hotkey) || (focus && m.IsButtonJustPressed(b.Left))
	var tick = internal.Runtime > r.wPressedAt+0.5
	var inputHold = k.IsKeyPressed(hotkey) || (focus && r.wPressedOn == widget && m.IsButtonPressed(b.Left))
	var hold = inputHold && condition.TrueEvery(0.1, text.New(";;hold-", buttonId, "-", hotkey)) && tick

	return first || hold
}
