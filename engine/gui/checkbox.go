package gui

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
)

func Checkbox(id string, properties ...string) string {
	return newWidget("checkbox", id, properties...)
}

//=================================================================
// private

func checkbox(cam *graphics.Camera, root *root, widget *widget, owner *container) {
	var on = themedProp(property.CheckboxThemeId, root, owner, widget)
	var off = themedProp(property.ThemeId, root, owner, widget)
	widget.ThemeId = condition.If(widget.Properties[property.Value] == "", off, on)

	button(cam, root, widget, owner)

	if root.ButtonClickedOnce(widget.Id, cam) {
		widget.Properties[property.Value] = condition.If(widget.Properties[property.Value] == "", "v", "")

		var group = themedProp(property.CheckboxGroup, root, owner, widget)
		if group == "" {
			return
		}

		for _, w := range root.Widgets {
			var wOwner = root.Containers[w.OwnerId]
			var wGroup = themedProp(property.CheckboxGroup, root, wOwner, w)
			if wGroup == group {
				w.Properties[property.Value] = ""
			}
		}

		widget.Properties[property.Value] = "v"
	}
}
