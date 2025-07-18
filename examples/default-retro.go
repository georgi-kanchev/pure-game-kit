package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/text"
	"pure-kit/engine/window"
)

func DefaultRetro() {
	var camera = graphics.NewCamera(1)
	var assetId, tileIds = assets.LoadDefaultAtlasRetro()
	var sprite = graphics.NewSprite(assetId)

	var textBox = graphics.NewTextBox("", 5, 5, "")
	textBox.LineHeight, textBox.GapSymbols, textBox.Color = 10, 0.5, color.Cyan

	for window.KeepOpen() {
		camera.SetScreenAreaToWindow()
		camera.PivotX, camera.PivotY = 0, 0
		sprite.Fit(&camera)
		sprite.ScaleX *= 0.9
		sprite.ScaleY *= 0.9

		var mx, my = sprite.MousePosition(&camera)
		var mmx, mmy = sprite.PointToCamera(&camera, number.Snap(mx-4.5, 9), number.Snap(my-4.5, 9))
		var imx, imy = int(mx / 9), int(my / 9)
		var index = number.LimitInt(number.Indexes2DToIndex1D(imy, imx, 26, 21), 0, len(tileIds)-1)

		camera.DrawSprites(&sprite)

		if !sprite.IsHovered(&camera) {
			continue
		}

		var info = text.New("id: ", tileIds[index], "\ncoords: ", imx, ", ", imy, "\nindex: ", index)
		camera.DrawFrame(mmx, mmy, 8*sprite.ScaleX, 8*sprite.ScaleY, 0, 6, color.Cyan)

		textBox.Value = info
		camera.DrawTextBoxes(&textBox)
	}
}
