package gui

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/gui/field"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

func Slider(id string, properties ...string) string {
	return newWidget("slider", id, properties...)
}

//=================================================================
// private

func slider(cam *graphics.Camera, root *root, widget *widget) {
	var owner = root.Containers[widget.OwnerId]
	var assetId = root.themedField(field.AssetId, owner, widget)
	if assetId == "" {
		widget.Height /= 2
		widget.Y += widget.Height / 2
	}
	btnSounds = false
	button(cam, root, widget)
	btnSounds = true
	if assetId == "" {
		widget.Y -= widget.Height / 2
		widget.Height *= 2
	}

	var _, h = assets.Size(assetId)
	var ratio = widget.Height / float32(h)
	var handleAssetId = root.themedField(field.SliderHandleAssetId, owner, widget)
	var hw, hh = assets.Size(handleAssetId)
	var handleWidth, handleHeight = float32(hw), float32(hh)
	handleWidth *= ratio
	handleHeight *= ratio
	var handleY = widget.Y - (handleWidth)/3
	var value = parseNum(widget.Fields[field.Value], 0)
	var step = parseNum(root.themedField(field.SliderStep, owner, widget), 0)

	if value != widget.PrevValue && !sound.IsPlaying() {
		sound.AssetId = defaultValue(root.themedField(field.SliderSound, owner, widget), "~slider")
		sound.Volume = root.Volume
		sound.Play()
	}

	widget.PrevValue = value

	if step > 0 {
		var stepPx = (widget.Width - handleWidth) * step
		var totalSteps = int(number.RoundUp((1-step)/step, -1))
		var stepAssetId = root.themedField(field.SliderStepAssetId, owner, widget)

		for i := 1; i <= totalSteps; i++ {
			var stepX = (widget.X + handleWidth/2) + float32(i)*stepPx
			if stepAssetId != "" && stepPx > widget.Height {
				reusableWidget.Width, reusableWidget.Height = widget.Height, widget.Height
				drawReusableWidget(buttonColor, stepAssetId, stepX-widget.Height/2, widget.Y, root, cam)
			} else {
				cam.DrawQuad(stepX, widget.Y, 5, widget.Height, 0, buttonColor)
			}
		}
	}

	if root.wPressedOn == widget {
		var mx, _ = cam.MousePosition()
		value = number.Map(mx, widget.X+handleWidth/2, widget.X+widget.Width-handleWidth/2, 0, 1)
		value = widget.setSliderValue(value, root)
	}

	var x = number.Map(value, 0, 1, widget.X, widget.X+widget.Width-handleWidth)
	buttonColor = color.Brighten(buttonColor, 0.5)

	if handleAssetId == "" {
		cam.DrawCircle(x, handleY+handleWidth*0.8, handleWidth, color.Darken(buttonColor, 0.5))
		cam.DrawCircle(x, handleY+handleWidth*0.8, handleWidth*0.75, buttonColor)
	} else {
		reusableWidget.Width, reusableWidget.Height = handleWidth, handleHeight
		drawReusableWidget(buttonColor, handleAssetId, x, handleY, root, cam)
	}
}

func (widget *widget) setSliderValue(value float32, root *root) float32 {
	var owner = root.Containers[widget.OwnerId]
	var step = parseNum(root.themedField(field.SliderStep, owner, widget), 0)
	value = number.Snap(value, number.Unsign(step))
	value = number.Limit(value, 0, 1)
	widget.Fields[field.Value] = text.New(value)
	return value
}
func drawReusableWidget(col uint, assetId string, x, y float32, root *root, cam *graphics.Camera) {
	var r, g, b, a = color.Channels(col)
	clear(reusableWidget.Fields)
	reusableWidget.Fields[field.AssetId] = assetId
	reusableWidget.Fields[field.Color] = text.New(r, " ", g, " ", b, " ", a)
	reusableWidget.X, reusableWidget.Y = x, y

	setupVisualsTextured(root, reusableWidget)
	drawVisuals(cam, root, reusableWidget, false, nil)
}
