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

	if w.handle == nil {
		w.handle = graphics.NewSprite("", 0, 0)
	}

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
		var totalSteps = int(number.RoundUp(1/step)) - 1
		var stepAssetId = w.root.themedField(field.SliderStepAssetId, owner, w)

		if len(w.steps) < totalSteps {
			w.steps = make([]*graphics.Sprite, totalSteps)
			for i := range w.steps {
				w.steps[i] = graphics.NewSprite("", 0, 0)
			}
		}

		for i := 1; i <= totalSteps; i++ {
			var stepX = w.X + float32(i)*stepPx
			var step = w.steps[i-1]
			step.X, step.Y = stepX, w.Y
			step.Width, step.Height = w.Height, w.Height
			step.TextureId, step.Tint = stepAssetId, buttonColor
			step.Mask = owner.mask
			step.PivotX, step.PivotY = 0, 0

			if stepAssetId == "" {
				step.X += handleWidth / 2
				step.Width = 4
			}
		}
		w.root.sprites = append(w.root.sprites, w.steps[:totalSteps]...)
	}

	if w.root.wPressedOn == w {
		var mx, _ = w.root.cam.MousePosition()
		value = number.Map(mx, w.X+handleWidth/2, w.X+w.Width-handleWidth/2, 0, 1)
		value = w.setSliderValue(value)
	}

	var x = number.Map(value, 0, 1, w.X, w.X+w.Width-handleWidth)
	buttonColor = color.Brighten(buttonColor, 0.5)

	w.handle.X, w.handle.Y = x, handleY
	w.handle.Width, w.handle.Height = handleWidth, handleHeight
	w.handle.TextureId, w.handle.Tint = handleAssetId, buttonColor
	w.handle.PivotX, w.handle.PivotY = 0, 0
	w.handle.Mask = owner.mask

	if handleAssetId == "" {
		w.handle.Y = w.Y
	}

	w.root.sprites = append(w.root.sprites, w.handle)
}

func (w *widget) setSliderValue(value float32) float32 {
	var owner = w.root.Containers[w.OwnerId]
	var step = parseNum(w.root.themedField(field.SliderStep, owner, w), 0)
	value = number.Snap(value, number.Unsign(step))
	value = number.Limit(value, 0, 1)
	w.Fields[field.Value] = text.New(value)
	return value
}
