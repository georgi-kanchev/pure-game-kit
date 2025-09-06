package gui

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/gui/property"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/symbols"
)

func Slider(id string, properties ...string) string {
	return newWidget("slider", id, properties...)
}

//=================================================================
// private

var reusableWidget *widget = &widget{Properties: map[string]string{}}

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

	var _, h = assets.Size(assetId)
	var ratio = widget.Height / h
	var handleAssetId = themedProp(property.SliderHandleAssetId, root, owner, widget)
	var handleWidth, handleHeight = assets.Size(handleAssetId)
	handleWidth *= ratio
	handleHeight *= ratio
	var handleY = widget.Y - (handleWidth)/3
	var value = parseNum(widget.Properties[property.Value], 0)
	var step = parseNum(themedProp(property.SliderStep, root, owner, widget), 0)

	if step > 0 {
		var stepPx = (widget.Width - handleWidth) * step
		var totalSteps = int(number.RoundUp((1-step)/step, -1))
		var stepAssetId = themedProp(property.SliderStepAssetId, root, owner, widget)

		for i := 1; i <= totalSteps; i++ {
			var stepX = (widget.X + handleWidth/2) + float32(i)*stepPx
			if stepAssetId != "" && stepPx > widget.Height {
				reusableWidget.Width, reusableWidget.Height = widget.Height, widget.Height
				drawReusableWidget(buttonColor, stepAssetId, stepX-widget.Height/2, widget.Y, root, owner, cam)
			} else {
				cam.DrawRectangle(stepX, widget.Y, 5, widget.Height, 0, buttonColor)
			}
		}
	}

	if wPressedOn == widget {
		var mx, _ = cam.MousePosition()
		value = number.Map(mx, widget.X, widget.X+widget.Width-handleWidth, 0, 1)
		value = widget.setSliderValue(value, root, owner)
	}

	var x = number.Map(value, 0, 1, widget.X, widget.X+widget.Width-handleWidth)
	buttonColor = color.Brighten(buttonColor, 0.5)

	if handleAssetId == "" {
		cam.DrawCircle(x, handleY+handleWidth*0.8, handleWidth, color.Darken(buttonColor, 0.5))
		cam.DrawCircle(x, handleY+handleWidth*0.8, handleWidth*0.75, buttonColor)
	} else {
		reusableWidget.Width, reusableWidget.Height = handleWidth, handleHeight
		drawReusableWidget(buttonColor, handleAssetId, x, handleY, root, owner, cam)
	}
}

func (widget *widget) setSliderValue(value float32, root *root, owner *container) float32 {
	var step = parseNum(themedProp(property.SliderStep, root, owner, widget), 0)
	value = number.Snap(value, number.Unsign(step))
	value = number.Limit(value, 0, 1)
	widget.Properties[property.Value] = symbols.New(value)
	return value
}
func drawReusableWidget(col uint, assetId string, x, y float32, root *root, owner *container, cam *graphics.Camera) {
	var r, g, b, a = color.Channels(col)
	clear(reusableWidget.Properties)
	reusableWidget.Properties[property.AssetId] = assetId
	reusableWidget.Properties[property.Color] = symbols.New(r, " ", g, " ", b, " ", a)
	reusableWidget.X, reusableWidget.Y = x, y

	setupVisualsTextured(root, reusableWidget, owner)
	drawVisuals(cam, root, reusableWidget, owner)
}
