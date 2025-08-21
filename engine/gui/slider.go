package gui

import (
	"math"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/input/mouse"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/symbols"
)

func Slider(id string, properties ...string) string {
	return newWidget("slider", id, properties...)
}

// #region private

var handleWidget *widget = &widget{Properties: map[string]string{}}

func slider(cam *graphics.Camera, root *root, widget *widget, owner *container) {
	var assetId = themedProp(property.AssetId, root, owner, widget)
	if assetId == "" {
		widget.Height /= 2
		widget.Y += widget.Height / 2
	}
	button(cam, root, widget, owner)
	if assetId == "" {
		widget.Y -= widget.Height / 2
		widget.Height *= 2
	}

	var handleWidth = (widget.Height * 1.2) / 2
	var x = widget.X + handleWidth
	var value = parseNum(widget.Properties[property.SliderValue], 0)
	var step = parseNum(themedProp(property.SliderStep, root, owner, widget), 0)
	var handleAssetId = themedProp(property.SliderHandleAssetId, root, owner, widget)

	if step > 0 {
		var stepPx = (widget.Width - handleWidth*2) * step
		var totalSteps = int(math.Ceil(float64((1 - step) / step)))
		var stepAssetId = themedProp(property.SliderStepAssetId, root, owner, widget)

		for i := 1; i <= totalSteps; i++ {
			var stepX = (widget.X + handleWidth) + float32(i)*stepPx
			if stepAssetId != "" && stepPx > widget.Height {
				handleWidget.Width, handleWidget.Height = widget.Height, widget.Height
				drawReusableWidget(buttonColor, stepAssetId, stepX-widget.Height/2, widget.Y, root, owner, cam)
			} else {
				cam.DrawRectangle(stepX, widget.Y, 5, widget.Height, 0, buttonColor)
			}
		}
	}

	if pressedOn == widget {
		var mx, _ = cam.MousePosition()
		value = number.Map(mx, widget.X+handleWidth, widget.X+widget.Width-handleWidth, 0, 1)
		value = widget.setSliderValue(value, root, owner)
	}

	if widget.IsHovered(owner, cam) && mouse.Scroll() != 0 {
		step = number.Limit(float32(math.Abs(float64(step))), 0.05, 1)
		value -= step * float32(mouse.Scroll())
		value = widget.setSliderValue(value, root, owner)
	}

	x = number.Map(value, 0, 1, widget.X+handleWidth, widget.X+widget.Width-handleWidth)

	if handleAssetId == "" {
		var smaller = handleWidth * 0.75
		cam.DrawCircle(x, widget.Y+smaller, handleWidth, color.Gray)
		cam.DrawCircle(x, widget.Y+smaller, smaller, buttonColor)
		return
	}

	handleWidget.Width, handleWidget.Height = widget.Height*1.2, widget.Height*1.2
	drawReusableWidget(buttonColor, handleAssetId, x-handleWidth, widget.Y, root, owner, cam)
}

func (widget *widget) setSliderValue(value float32, root *root, owner *container) float32 {
	var step = parseNum(themedProp(property.SliderStep, root, owner, widget), 0)
	value = number.Snap(value, float32(math.Abs(float64(step))))
	value = number.Limit(value, 0, 1)
	widget.Properties[property.SliderValue] = symbols.New(value)
	return value
}

func drawReusableWidget(col uint, assetId string, x, y float32, root *root, owner *container, cam *graphics.Camera) {
	var r, g, b, a = color.Channels(col)
	clear(handleWidget.Properties)
	handleWidget.Properties[property.AssetId] = assetId
	handleWidget.Properties[property.Color] = symbols.New(r, " ", g, " ", b, " ", a)
	handleWidget.X, handleWidget.Y = x, y

	setupVisualsTextured(root, handleWidget, owner)
	drawVisuals(cam, root, handleWidget, owner)
}

// #endregion
