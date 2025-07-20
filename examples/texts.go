package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/utility/color"
	"pure-kit/engine/window"
)

func Texts() {
	var cam = graphics.NewCamera(1)
	var font = assets.LoadFonts(32, "font.ttf")[0]
	var _, tiles = assets.LoadDefaultAtlasIcons()
	var textBox = graphics.NewTextBox(font, 0, 0, "Lorem `ipsum ^^ dolor` sit amet, consectetur ^^ adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. `Ut enim ^^ ad minim` veniam, quis nostrud exercitation ^^\r^^ ullamco laboris nisi ut aliquip ex ea commodo consequat. `Duis aute irure dolor in reprehenderit in voluptate velit esse cillum doloreeu fugiat nulla pariatur.")
	textBox.PivotX, textBox.PivotY = 0, 0
	textBox.AlignmentY = 1
	textBox.LineGap = -1
	textBox.Angle = 5
	textBox.EmbeddedColors = []uint{color.Red, color.Green, color.Black, color.White, color.Blue}
	textBox.EmbeddedAssetIds = []string{tiles[162], tiles[256], tiles[156], tiles[154], tiles[157]}

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.PivotX, cam.PivotY = 0.05, 0
		textBox.Color = color.Gray
		cam.DrawNodes(&textBox.Node)
		textBox.Color = color.White
		cam.DrawTextBoxes(&textBox)

		textBox.Width, textBox.Height = textBox.MousePosition(&cam)

	}
}
