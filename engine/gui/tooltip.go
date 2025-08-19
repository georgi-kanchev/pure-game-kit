package gui

import (
	"pure-kit/engine/execution/condition"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/seconds"
)

func Tooltip(id string, properties ...string) string {
	return newWidget("tooltip", id, properties...)
}

func tryShowTooltip(
	wId string, root *root, c *container, widget *widget, cam *graphics.Camera, x, y int, w, h float32) {
	var hov = widget.IsHovered(c, cam)
	if condition.TrueOnce(hov, ";;hoverrr-"+wId) {
		tooltipForWidget = widget
		tooltipAt = seconds.RealRuntime()
	}
	if condition.TrueOnce(!hov, ";;unhoverrr-"+wId) {
		tooltipForWidget = nil
	}

	var tooltipId = themedProp(property.TooltipId, root, c, widget)
	var tooltipText = widget.Properties[property.TooltipText]
	var tooltip = root.Widgets[tooltipId]
	if tooltipId == "" || tooltipText == "" || tooltip == nil || tooltipForWidget != widget ||
		tooltipAt+1 > seconds.RealRuntime() {
		return
	}

	var camW, camH = cam.Size()
	cam.Mask(cam.ScreenX, cam.ScreenY, cam.ScreenWidth, cam.ScreenHeight)
	tooltip.Properties[property.Text] = tooltipText
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

	setupVisualsTextured(root, tooltip, c)
	drawVisuals(cam, root, tooltip, c)

	cam.Mask(x, y, int(w), int(h))
}
