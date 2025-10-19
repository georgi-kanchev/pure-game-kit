package gui

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
)

func Checkbox(id string, properties ...string) string {
	return newWidget("checkbox", id, properties...)
}

//=================================================================
// private

func checkbox(cam *graphics.Camera, root *root, widget *widget) {
	var owner = root.Containers[widget.OwnerId]
	var on = themedProp(field.CheckboxThemeId, root, owner, widget)
	var off = themedProp(field.ThemeId, root, owner, widget)
	widget.ThemeId = condition.If(widget.Properties[field.Value] == "", off, on)

	button(cam, root, widget)

	if root.IsButtonClickedOnce(widget.Id, cam) {
		widget.Properties[field.Value] = condition.If(widget.Properties[field.Value] == "", "v", "")

		var group = themedProp(field.CheckboxGroup, root, owner, widget)
		if group == "" {
			return
		}

		for _, w := range root.Widgets {
			var wOwner = root.Containers[w.OwnerId]
			var wGroup = themedProp(field.CheckboxGroup, root, wOwner, w)
			if wGroup == group {
				w.Properties[field.Value] = ""
			}
		}

		widget.Properties[field.Value] = "v"
	}
}
