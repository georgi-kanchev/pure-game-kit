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
	runDefaultAssetDisplay(0.7, 50, 0, 22, 13, assets.LoadDefaultAtlasIcons)
}
func DefaultAssetCursors() {
	runDefaultAssetDisplay(0.9, 32, 0, 19, 8, assets.LoadDefaultAtlasCursors)
}
func DefaultAssetInput() {
	runDefaultAssetDisplay(0.9, 50, 0, 17, 6, assets.LoadDefaultAtlasInput)
}
func DefaultAssetPatterns() {
	runDefaultAssetDisplay(0.7, 64, 0, 12, 7, assets.LoadDefaultAtlasPatterns)
}
func DefaultAssetFont() {
	var function = func() (string, []string) { return "", []string{} }
	runDefaultAssetDisplay(0.7, 100, 0, 0, 0, function)
}
func DefaultAssetTexture() {
	assets.LoadDefaultTexture()
	var function = func() (string, []string) { return "", []string{} }
	runDefaultAssetDisplay(0.7, 256, 0, 12, 7, function)
}
func DefaultAssetUI() {
	runDefaultAssetDisplay(0.7, 16, 0, 6, 6, assets.LoadDefaultAtlasUI)
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

	for window.KeepOpen() {
		camera.SetScreenAreaToWindow()
		textBox.Width, textBox.Height = camera.Size()
		camera.PivotX, camera.PivotY = 0.5, 0.5
		sprite.CameraFit(&camera)
		sprite.ScaleX *= scale
		sprite.ScaleY *= scale

		var mx, my = sprite.MousePosition(&camera)
		var sx, sy = number.Snap(mx-fullSz/2, fullSz), number.Snap(my-fullSz/2, fullSz)
		var mmx, mmy = sprite.PointToCamera(&camera, sx, sy)
		var imx, imy = int(mx / fullSz), int(my / fullSz)
		var index = number.LimitInt(number.Indexes2DToIndex1D(imy, imx, int(w), int(h)), 0, len(tileIds)-1)

		camera.DrawSprites(&sprite)

		if !sprite.MouseIsHovering(&camera) {
			continue
		}
		camera.DrawFrame(mmx, mmy, tileSize*sprite.ScaleX, tileSize*sprite.ScaleY, 0, 6, color.Cyan)

		var txt = ""
		if index < len(tileIds) {
			txt = tileIds[index]
		}

		var w, h = assets.Size(txt)
		var info = symbols.New(
			"id: '", txt, "'",
			"\ncoords: ", imx, ", ", imy,
			"\nindex: ", index,
			"\nsize:", tileSize, "x", tileSize)

		if txt == "" && len(tileIds) == 0 && imx == 0 && imy == 0 { // display default texture & font
			info = symbols.New("id: '", txt, "'", "\nsize:", w, "x", h)
		}

		textBox.Text = info
		camera.PivotX, camera.PivotY = 0, 0
		camera.DrawTextBoxes(&textBox)
	}
}

// #endregion
