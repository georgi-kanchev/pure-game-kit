package example

import (
	"fmt"
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/graphics/tag"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/text"
	"pure-game-kit/window"
)

func Texts() {
	var cam = graphics.NewCamera(1)
	var font = assets.LoadDefaultFont()
	var _, tiles = assets.LoadDefaultAtlasIcons()
	var textBox = graphics.NewTextBox(font, 0, 0, "")
	textBox.PivotX, textBox.PivotY = 0.5, 0.5
	textBox.AlignmentX = 0
	textBox.AlignmentY = 1
	textBox.Angle = 0
	textBox.LineHeight = 200
	// textBox.Width, textBox.Height = 1090, 1496

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		textBox.Tint = palette.DarkGray
		cam.DrawNodes(&textBox.Node)
		textBox.Tint = palette.White

		textBox.Width, textBox.Height = textBox.MousePosition(cam)

		textBox.Text = text.New(
			"Lorem ",
			tag.Color("ipsum dolor", palette.Red),
			" sit amet, ",
			tag.ShadowBold("consectetur"),
			" adipiscing elit, sed do ",
			tag.BackColor("eiusmod tempor", palette.Green),
			" incididunt ut labore et ",
			tag.Underline("dolore magna"),
			" aliqua. Ut ",
			tag.Bold("enim ad minim"),
			" veniam, quis ",
			tag.Asset(tiles[162]),
			"exercitation ullamco laboris nisi ut ",
			tag.OutlineColor("aliquip ex ea commodo", palette.Red),
			" consequat. ",
			tag.Color("Duis aute", palette.Blue),
			" irure dolor ",
			tag.ShadowColor("in reprehenderit", palette.Green),
			" in voluptate ",
			tag.Crossout("velit esse incididunt ut labore"),
			" cillum ",
			tag.Thin("doloreeu"),
			" fugiat nulla pariatur.")

		cam.DrawTextBoxes(textBox)

		if keyboard.IsKeyJustPressed(key.A) {
			fmt.Printf("textBox.Text: %v\n", textBox.Text)
		}

		fmt.Printf("%v, %v\n", textBox.Width, textBox.Height)
	}
}
