package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/graphics/tag"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/text"
	"pure-game-kit/window"
)

func Texts() {
	var cam = graphics.NewCamera(1)
	var font = assets.LoadDefaultFont()
	var _, tiles = assets.LoadDefaultAtlasIcons()
	var textBox = graphics.NewTextBox(font, 0, 0, "")
	textBox.PivotX, textBox.PivotY = 0, 0
	textBox.Angle = 0
	// textBox.EmbeddedAssetIds = []string{tiles[162], tiles[256], tiles[156], tiles[154], tiles[157]}
	textBox.LineHeight = 100
	textBox.Width, textBox.Height = 2000, 500

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		textBox.X, textBox.Y = cam.PointFromScreen(50, 50)
		textBox.Tint = palette.DarkGray
		cam.DrawNodes(&textBox.Node)
		textBox.Tint = palette.White

		textBox.Width, textBox.Height = textBox.MousePosition(cam)

		textBox.Text = text.New("Lorem ", tag.Color("ipsum dolor", palette.Red), " sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et ", tag.Underline("dolore magna", 1), " aliqua. Ut enim ad minim veniam, quis nostrud ", tag.Asset(tiles[162]), " exercitation ullamco laboris nisi ut ", tag.Box("aliquip", 1), " ex ea commodo consequat. ", tag.ColorOutline("Duis aute", palette.Blue), " irure dolor in reprehenderit in voluptate ", tag.Strikethrough("velit esse", 1), " cillum doloreeu fugiat nulla pariatur.")

		cam.DrawTextBoxes(textBox)
	}
}
