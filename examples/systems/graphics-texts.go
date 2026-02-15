package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/time"
	"pure-game-kit/window"
)

func Texts() {
	var cam = graphics.NewCamera(1)
	var font = assets.LoadDefaultFont()
	var _, tiles = assets.LoadDefaultAtlasIcons()
	var textBox = graphics.NewTextBox(font, 0, 0, "Lorem `ipsum ^^ dolor` sit amet, *consectetur* ^^ adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. `Ut enim ^^ ad minim` veniam, quis nostrud *exercitation* ^^\r^^ ullamco laboris nisi ut aliquip ex ea commodo consequat. `Duis aute irure dolor in reprehenderit in voluptate velit esse cillum doloreeu fugiat nulla pariatur.")
	textBox.PivotX, textBox.PivotY = 0, 0
	textBox.AlignmentY = 1
	textBox.Angle = 0
	textBox.EmbeddedColors = []uint{palette.Red, palette.Green, palette.Black, palette.White, palette.Blue}
	textBox.EmbeddedAssetIds = []string{tiles[162], tiles[256], tiles[156], tiles[154], tiles[157]}
	textBox.EmbeddedThicknesses = []float32{0.7, 0.5, 0.35}
	textBox.LineHeight = 100

	var a float32 = 0

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		textBox.X, textBox.Y = cam.PointFromScreen(50, 0)
		textBox.Tint = palette.DarkGray
		cam.DrawNodes(&textBox.Node)
		textBox.Tint = palette.White

		a = number.Sine(time.Runtime() / 5)

		textBox.Thickness = float32(0.5 + a/2)

		cam.DrawTextBoxes(textBox)
		textBox.Width, textBox.Height = textBox.MousePosition(cam)
	}
}
