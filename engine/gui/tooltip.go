package gui

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/input/mouse"
	"pure-kit/engine/input/mouse/cursor"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/time"
)

func Tooltip(id string, properties ...string) string {
	return newWidget("tooltip", id, properties...)
}

//=================================================================
// private

var tooltip *widget
var tooltipForWidget *widget
var tooltipAt float32

func tryShowTooltip(widget *widget, root *root, c *container, cam *graphics.Camera) {
	var hov = widget.isFocused(root, cam)

	if condition.TrueOnce(hov, ";;hoverrr-"+widget.Id) {
		tooltipForWidget = widget
		tooltipAt = time.RealRuntime()
		var tooltipId = themedProp(property.TooltipId, root, c, widget)
		tooltip = root.Widgets[tooltipId]

		if tooltip != nil {
			var text = themedProp(property.TooltipText, root, c, widget)
			tooltip.Properties[property.Text] = text

			if text != "" {
				mouse.SetCursor(cursor.Hand)
			}
		}
	}
	if widget == tooltipForWidget && condition.TrueOnce(!hov, ";;unhoverrr-"+widget.Id) {
		tooltipForWidget = nil
		tooltip = nil
	}
}
func drawTooltip(root *root, c *container, cam *graphics.Camera) {
	if tooltip == nil || tooltipForWidget == nil || time.RealRuntime() < tooltipAt+0.5 {
		return
	}
	if tooltip.Properties[property.Text] == "" || !tooltipForWidget.isFocused(root, cam) {
		return
	}

	var camW, camH = cam.Size()
	var width = parseNum(dyn(c, tooltip.Properties[property.Width], "500"), 500)
	var margin = parseNum(themedProp(property.TooltipMargin, root, c, tooltip), 50)
	tooltip.Width, tooltip.Height = width-margin, camH

	setupVisualsText(root, tooltip, true)

	var lines = reusableTextBox.TextLines()
	var lh = reusableTextBox.LineHeight
	var textH = float32(len(lines)*int(lh+reusableTextBox.LineGap)) + lh
	reusableTextBox.Height = textH
	reusableTextBox.X = tooltipForWidget.X + tooltipForWidget.Width/2 - reusableTextBox.Width/2
	reusableTextBox.Y = tooltipForWidget.Y - textH
	reusableTextBox.X = number.Limit(reusableTextBox.X, -camW/2, camW/2-width)
	reusableTextBox.Y = number.Limit(reusableTextBox.Y, -camH/2, camH/2-textH)
	tooltip.X, tooltip.Y, tooltip.Width, tooltip.Height = reusableTextBox.X, reusableTextBox.Y, width, textH

	reusableTextBox.X += margin / 2

	if tooltip.Y+tooltip.Height > tooltipForWidget.Y+2 { // margin of error 2 pixels
		tooltip.Y = tooltipForWidget.Y + tooltipForWidget.Height
		reusableTextBox.Y = tooltip.Y
	}

	setupVisualsTextured(root, tooltip)
	drawVisuals(cam, root, tooltip, false, nil)
}
