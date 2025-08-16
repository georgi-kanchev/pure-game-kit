package gui

import (
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/input/mouse"
)

func WidgetButton(id string, properties ...string) string {
	return newWidget("button", id, properties...)
}

// #region private

func button(w, h float32, cam *graphics.Camera, root *root, widget *widget, owner *container) {
	var prev = widget.AssetId
	var hover = themedProp(property.ButtonAssetIdHover, root, owner, widget)
	var press = themedProp(property.ButtonAssetIdPress, root, owner, widget)

	if widget.IsHovered(root, cam) {
		if hover != "" {
			widget.AssetId = hover
		}

		if press != "" && mouse.IsButtonPressed(mouse.ButtonLeft) {
			widget.AssetId = press
		}
	}

	visual(w, h, cam, root, widget, owner)
	widget.AssetId = prev
}

// #endregion
