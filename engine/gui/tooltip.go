package gui

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/input/mouse"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/seconds"
)

func Tooltip(id string, properties ...string) string {
	return newWidget("tooltip", id, properties...)
}

// #region private

var tooltip *widget
var tooltipForWidget *widget
var tooltipAt float32

func tryShowTooltip(widget *widget, root *root, c *container, cam *graphics.Camera) {
	var hov = widget.IsHovered(c, cam)
	if condition.TrueOnce(hov, ";;hoverrr-"+widget.Id) {
		tooltipForWidget = widget
		tooltipAt = seconds.RealRuntime()
		var tooltipId = themedProp(property.TooltipId, root, c, widget)
		tooltip = root.Widgets[tooltipId]

		if tooltip != nil {
			var text = widget.Properties[property.TooltipText]
			tooltip.Properties[property.Text] = text

			if text != "" {
				mouse.SetCursor(mouse.CursorHand)
			}
		}
	}
	if widget == tooltipForWidget && condition.TrueOnce(!hov, ";;unhoverrr-"+widget.Id) {
		tooltipForWidget = nil
		tooltip = nil
	}
}
func drawTooltip(root *root, c *container, cam *graphics.Camera) {
	if tooltip == nil || tooltipForWidget == nil || seconds.RealRuntime() < tooltipAt+0.5 {
		return
	}
	if tooltip.Properties[property.Text] == "" {
		return
	}

	var camW, camH = cam.Size()
	var width = parseNum(dyn(c, tooltip.Properties[property.Width], "500"), 500)
	tooltip.Width, tooltip.Height = width, camH
	setupVisualsText(root, tooltip, c)

	var lines = reusableTextBox.TextLines()
	var lh = reusableTextBox.LineHeight
	var textH = float32(len(lines)*int(lh+reusableTextBox.LineGap)) + lh
	reusableTextBox.Height = textH
	reusableTextBox.X = tooltipForWidget.X + tooltipForWidget.Width/2 - reusableTextBox.Width/2
	reusableTextBox.Y = tooltipForWidget.Y - textH
	reusableTextBox.X = number.Limit(reusableTextBox.X, -camW/2, camW/2-width)
	reusableTextBox.Y = number.Limit(reusableTextBox.Y, -camH/2, camH/2-textH)
	tooltip.X, tooltip.Y, tooltip.Width, tooltip.Height = reusableTextBox.X, reusableTextBox.Y, width, textH

	if tooltip.Y+tooltip.Height > tooltipForWidget.Y+2 { // margin of error 2 pixels
		tooltip.Y = tooltipForWidget.Y + tooltipForWidget.Height
		setupVisualsText(root, tooltip, c)
	}

	setupVisualsTextured(root, tooltip, c)
	drawVisuals(cam, root, tooltip, c)
}

// #endregion
