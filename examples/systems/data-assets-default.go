package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/internal"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
	"pure-game-kit/window"
)

func DefaultAssetRetro() {
	runDefaultAssetDisplay(0.9, 8, 1, 26, 21, assets.LoadDefaultAtlasRetro)
}
func DefaultAssetIcons() {
	runDefaultAssetDisplay(0.7, 50, 0, 22, 13, assets.LoadDefaultAtlasIcons)
}
func DefaultAssetCursors() {
	runDefaultAssetDisplay(0.75, 32, 0, 19, 8, assets.LoadDefaultAtlasCursors)
}
func DefaultAssetInput() {
	runDefaultAssetDisplay(0.85, 50, 0, 17, 6, assets.LoadDefaultAtlasInput)
}
func DefaultAssetPatterns() {
	runDefaultAssetDisplay(0.7, 64, 1, 12, 7, assets.LoadDefaultAtlasPatterns)
}
func DefaultAssetFont() {
	runDefaultAssetDisplay(0.7, 1024, 0, 0, 0, func() (string, []string) { return "", []string{} })
}
func DefaultAssetUI() {
	runDefaultAssetDisplay(0.7, 16, 0, 9, 8, func() (string, []string) {
		var a, b, _ = assets.LoadDefaultAtlasUI()
		return a, b
	})
}

//=================================================================
// private

func runDefaultAssetDisplay(scale float32, tileSize, gap, w, h float32, load func() (string, []string)) {
	var camera = graphics.NewCamera(1)
	var assetId, tileIds = load()
	var sprite = graphics.NewSprite(assetId, 0, 0)
	var textBox = graphics.NewTextBox("", 5, 5, "")
	textBox.LineGap, textBox.Tint = -1, palette.Red
	textBox.PivotX, textBox.PivotY = 0, 0
	var fullSz = tileSize + gap
	var txt = ""
	var aw, ah = assets.Size(txt)

	for window.KeepOpen() {
		camera.SetScreenAreaToWindow()
		textBox.Width, textBox.Height = camera.Size()
		sprite.CameraFit(camera)
		sprite.ScaleX *= scale
		sprite.ScaleY *= scale

		if w == 0 && h == 0 {
			sprite.Width, sprite.Height = float32(aw), float32(ah)
		} else {
			aw, ah = int(tileSize), int(tileSize)
		}

		var cmx, cmy = camera.MousePosition()
		var mx, my = sprite.PointToLocal(cmx, cmy)
		var sx, sy = number.Snap(mx-fullSz/2, fullSz), number.Snap(my-fullSz/2, fullSz)
		var mmx, mmy = sprite.PointToGlobal(sx, sy)
		var imx, imy = int(mx / fullSz), int(my / fullSz)
		var index = number.Limit(number.Indexes2DToIndex1D(imy, imx, int(w), int(h)), 0, len(tileIds)-1)

		if sprite.AssetId == "" { // default font display hacky wacky
			var tex = internal.Fonts[""].Texture
			var id = ";;;default-font-texture"
			internal.Textures[id] = &tex
			sprite.AssetId = id
		}

		camera.DrawSprites(sprite)

		if !sprite.ContainsPoint(cmx, cmy) {
			continue
		}

		if index < len(tileIds) {
			txt = tileIds[index]
			aw, ah = assets.Size(txt)

			if w != 0 || h != 0 {
				aw, ah = int(tileSize), int(tileSize)
			}
		}

		camera.DrawQuadFrame(mmx, mmy, float32(aw)*sprite.ScaleX, float32(ah)*sprite.ScaleY, 0, 6, palette.Cyan)

		var info = text.New(
			"id: '", txt, "'",
			"\ncell: ", imx, ", ", imy,
			"\nindex: ", index,
			"\ncoords: ", imx*int(tileSize+gap), ", ", imy*int(tileSize+gap),
			"\nsize:", tileSize, "x", tileSize)

		if txt == "" && len(tileIds) == 0 && imx == 0 && imy == 0 { // display default font
			info = text.New("id: '", txt, "'", "\nsize:", aw, "x", ah)
		}

		textBox.Text = info
		textBox.X, textBox.Y = camera.PointFromScreen(0, 0)
		camera.DrawTextBoxes(textBox)
	}
}
