package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/color"
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
	textBox.EmbeddedColors = []uint{color.Red, color.Green, color.Black, color.White, color.Blue}
	textBox.EmbeddedAssetIds = []string{tiles[162], tiles[256], tiles[156], tiles[154], tiles[157]}
	textBox.EmbeddedThicknesses = []float32{0.7, 0.5, 0.35}
	textBox.LineHeight = 100

	var a float32 = 0

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.PivotX, cam.PivotY = 0.05, 0
		textBox.Color = color.Darken(color.Gray, 0.5)
		cam.DrawNodes(&textBox.Node)
		textBox.Color = color.White

		a = number.Sine(time.Runtime() / 5)

		textBox.Thickness = float32(0.5 + a/2)

		var mx, my = textBox.MousePosition(cam)
		// var symbols = symbols.Count(textBox.TextSymbols())
		// for i := range symbols {
		// 	var x, y, w, h, a = textBox.TextSymbol(&cam, i)
		// 	cam.DrawRectangle(x, y, w, h, a, color.Brown)
		// }
		cam.DrawTextBoxes(&textBox)
		textBox.Width, textBox.Height = mx, my
	}
}
