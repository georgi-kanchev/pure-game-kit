package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/utility/number"
	"pure-kit/engine/utility/symbols"
	"pure-kit/engine/window"
)

func DefaultAssetRetro() {
	runDefaultAssetDisplay(0.9, 8, 1, 26, 21, assets.LoadDefaultAtlasRetro)
}
func DefaultAssetIcons() {
	var function = func() (string, []string) { return assets.LoadDefaultAtlasIcons(true) }
	runDefaultAssetDisplay(0.7, 50, 0, 22, 13, function)
}
func DefaultAssetCursors() {
	var function = func() (string, []string) { return assets.LoadDefaultAtlasCursors(true) }
	runDefaultAssetDisplay(0.9, 32, 0, 19, 8, function)
}
func DefaultAssetInput() {
	var function = func() (string, []string) { return assets.LoadDefaultAtlasInput(true) }
	runDefaultAssetDisplay(0.9, 50, 0, 17, 6, function)
}
func DefaultAssetPatterns() {
	var function = func() (string, []string) { return assets.LoadDefaultAtlasPatterns(true) }
	runDefaultAssetDisplay(0.7, 64, 0, 12, 7, function)
}
func DefaultAssetFont() {
	var function = func() (string, []string) { return "", []string{} }
	runDefaultAssetDisplay(0.7, 1024, 0, 0, 0, function)
}
func DefaultAssetTexture() {
	assets.LoadDefaultTexture()
	var function = func() (string, []string) { return "", []string{} }
	runDefaultAssetDisplay(0.7, 256, 0, 0, 0, function)
}
func DefaultAssetUI() {
	var function = func() (string, []string) {
		var a, b, _ = assets.LoadDefaultAtlasUI(false)
		return a, b
	}
	runDefaultAssetDisplay(0.7, 16, 0, 9, 8, function)
}

// #region private

func runDefaultAssetDisplay(scale float32, tileSize, gap, w, h float32, load func() (string, []string)) {
	var camera = graphics.NewCamera(1)
	var assetId, tileIds = load()
	var sprite = graphics.NewSprite(assetId, 0, 0)
	var textBox = graphics.NewTextBox(assets.LoadDefaultFont(), 5, 5, "")
	textBox.LineGap, textBox.Color = -1, color.Cyan
	textBox.PivotX, textBox.PivotY = 0, 0
	textBox.EmbeddedAssetsTag, textBox.EmbeddedColorsTag = 0, 0
	var fullSz = tileSize + gap
	var txt = ""
	var aw, ah = assets.Size(txt)

	for window.KeepOpen() {
		camera.SetScreenAreaToWindow()
		textBox.Width, textBox.Height = camera.Size()
		camera.PivotX, camera.PivotY = 0.5, 0.5
		sprite.CameraFit(camera)
		sprite.ScaleX *= scale
		sprite.ScaleY *= scale

		if w == 0 && h == 0 {
			sprite.Width, sprite.Height = aw, ah
		} else {
			aw, ah = tileSize, tileSize
		}

		var mx, my = sprite.MousePosition(camera)
		var sx, sy = number.Snap(mx-fullSz/2, fullSz), number.Snap(my-fullSz/2, fullSz)
		var mmx, mmy = sprite.PointToCamera(camera, sx, sy)
		var imx, imy = int(mx / fullSz), int(my / fullSz)
		var index = number.LimitInt(number.Indexes2DToIndex1D(imy, imx, int(w), int(h)), 0, len(tileIds)-1)

		camera.DrawSprites(&sprite)

		if !sprite.IsHovered(camera) {
			continue
		}

		if index < len(tileIds) {
			txt = tileIds[index]
			aw, ah = assets.Size(txt)

			if w != 0 || h != 0 {
				aw, ah = tileSize, tileSize
			}
		}

		camera.DrawFrame(mmx, mmy, aw*sprite.ScaleX, ah*sprite.ScaleY, 0, 6, color.Cyan)

		var info = symbols.New(
			"id: '", txt, "'",
			"\ncell: ", imx, ", ", imy,
			"\nindex: ", index,
			"\ncoords: ", imx*int(tileSize+gap), ", ", imy*int(tileSize+gap),
			"\nsize:", tileSize, "x", tileSize)

		if txt == "" && len(tileIds) == 0 && imx == 0 && imy == 0 { // display default texture & font
			info = symbols.New("id: '", txt, "'", "\nsize:", aw, "x", ah)
		}

		textBox.Text = info
		camera.PivotX, camera.PivotY = 0, 0
		camera.DrawTextBoxes(&textBox)
	}
}

// #endregion
