package gui

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/gui/field"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
)

func Slider(id string, properties ...string) string {
	return newWidget("slider", id, properties...)
}

func (g *GUI) IsSliderJustSlid(id string) bool {
	return sliderSlidId == id
}

//=================================================================
// private

func slider(w *widget) {
	var owner = w.root.Containers[w.OwnerId]
	var assetId = w.root.themedField(field.AssetId, owner, w)
	btnSounds = false
	button(w)
	btnSounds = true

	var _, h = assets.Size(assetId)
	var ratio = w.Height / float32(h)
	var handleAssetId = w.root.themedField(field.SliderHandleAssetId, owner, w)
	var hw, hh = assets.Size(handleAssetId)
	var handleWidth, handleHeight = float32(hw), float32(hh)
	handleWidth *= ratio
	handleHeight *= ratio
	if handleAssetId == "" {
		handleWidth, handleHeight = w.Height, w.Height
	}
	var handleY = w.Y - (handleWidth)/3
	var value = parseNum(w.Fields[field.Value], 0)
	var step = parseNum(w.root.themedField(field.SliderStep, owner, w), 0)
	var cam = w.root.cam

	if value != w.PrevValue && !sound.IsPlaying() {
		sound.AssetId = defaultValue(w.root.themedField(field.SliderSound, owner, w), "~slider")
		sound.Volume = w.root.Volume
		sound.Play()
	}

	if w.PrevValue != value {
		sliderSlidId = w.Id
	}
	w.PrevValue = value

	if step > 0 {
		var stepPx = (w.Width - handleWidth) * step
		var totalSteps = int(number.RoundUp((1 - step) / step))
		var stepAssetId = w.root.themedField(field.SliderStepAssetId, owner, w)

		for i := 1; i <= totalSteps; i++ {
			var stepX = (w.X + handleWidth/2) + float32(i)*stepPx
			if stepAssetId != "" && stepPx > w.Height {
				reusableWidget.Width, reusableWidget.Height = w.Height, w.Height
				drawReusableWidget(buttonColor, stepAssetId, stepX-w.Height/2, w.Y)
			} else {
				cam.DrawQuad(stepX-2.5, w.Y, 5, w.Height, 0, buttonColor)
			}
		}
	}

	if w.root.wPressedOn == w {
		var mx, _ = w.root.cam.MousePosition()
		value = number.Map(mx, w.X+handleWidth/2, w.X+w.Width-handleWidth/2, 0, 1)
		value = w.setSliderValue(value)
	}

	var x = number.Map(value, 0, 1, w.X, w.X+w.Width-handleWidth)
	buttonColor = color.Brighten(buttonColor, 0.5)

	if handleAssetId == "" {
		cam.DrawCircle(x+handleWidth/2, handleY+handleWidth*0.8, handleWidth/2, color.Darken(buttonColor, 0.5))
		cam.DrawCircle(x+handleWidth/2, handleY+handleWidth*0.8, handleWidth/3, buttonColor)
	} else {
		reusableWidget.Width, reusableWidget.Height = handleWidth, handleHeight
		drawReusableWidget(buttonColor, handleAssetId, x, handleY)
	}
}

func (w *widget) setSliderValue(value float32) float32 {
	var owner = w.root.Containers[w.OwnerId]
	var step = parseNum(w.root.themedField(field.SliderStep, owner, w), 0)
	value = number.Snap(value, number.Unsign(step))
	value = number.Limit(value, 0, 1)
	w.Fields[field.Value] = text.New(value)
	return value
}
func drawReusableWidget(col uint, assetId string, x, y float32) {
	var r, g, b, a = color.Channels(col)
	clear(reusableWidget.Fields)
	reusableWidget.Fields[field.AssetId] = assetId
	reusableWidget.Fields[field.Color] = text.New(r, " ", g, " ", b, " ", a)
	reusableWidget.X, reusableWidget.Y = x, y

	setupVisualsTextured(reusableWidget)
	drawVisuals(reusableWidget, false, nil)
}
