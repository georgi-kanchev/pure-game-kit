package gui

import (
	"pure-game-kit/execution/condition"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
	"pure-game-kit/input/mouse"
	"pure-game-kit/input/mouse/cursor"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/time"
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
		var tooltipId = themedProp(field.TooltipId, root, c, widget)
		tooltip = root.Widgets[tooltipId]

		if tooltip != nil {
			var text = themedProp(field.TooltipText, root, c, widget)
			tooltip.Properties[field.Text] = text

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
	if tooltip.Properties[field.Text] == "" || !tooltipForWidget.isFocused(root, cam) {
		return
	}

	var camW, camH = cam.Size()
	var width = parseNum(dyn(c, tooltip.Properties[field.Width], "500"), 500)
	var margin = parseNum(themedProp(field.TooltipMargin, root, c, tooltip), 50)
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
