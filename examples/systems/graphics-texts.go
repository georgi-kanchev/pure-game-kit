package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/graphics/tag"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/text"
	"pure-game-kit/utility/time"
	"pure-game-kit/window"
)

func Texts() {
	var cam = graphics.NewCamera(1)
	var _, tiles = assets.LoadDefaultAtlasIcons()
	var textBox = graphics.NewTextBox("", 0, 0, "")
	textBox.PivotX, textBox.PivotY = 0.5, 0.5
	textBox.AlignmentX, textBox.AlignmentY = 0, 1
	textBox.Angle = 0
	textBox.LineHeight = 160
	textBox.Width, textBox.Height = 1000, 500
	textBox.ShadowOffsetX, textBox.ShadowOffsetY = 0.5, 0.5
	textBox.Text = text.New(
		"Lorem ",
		tag.Color("ipsum dolor", palette.Red),
		" sit amet, ",
		tag.ShadowBold("consectetur"),
		" adipiscing elit, sed do\n",
		tag.BackColor("eiusmod tempor", palette.Green),
		" aliqua.\nUt ",
		tag.Bold("enim ad minim"),
		" veniam, quis\n",
		tag.Asset(tiles[162]),
		" in voluptate ",
		tag.Crossout("velit esse incididunt ut labore"),
		" incididunt ut labore et\n",
		tag.Underline("dolore magna"),
		" exercitation ",
		tag.ShadowBlur("ullamco laboris", 3),
		" nisi ut ",
		tag.OutlineColor("aliquip ex ea commodo", palette.Red),
		" consequat.\n",
		tag.Color("Duis aute", palette.Blue),
		" irure dolor ",
		tag.ShadowColor("in reprehenderit", palette.Green),
		" cillum ",
		tag.Thin("doloreeu"),
		" fugiat nulla pariatur.")
	textBox.WordWrap = false
	// textBox.Text = "Hello, World! Hello, World! Hello, World! Hello, World! Hello, World!"
	textBox.Angle = 5
	textBox.ScaleX, textBox.ScaleY = 2, 2

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		textBox.Tint = palette.DarkGray
		cam.DrawQuads(&textBox.Quad)
		textBox.Tint = palette.White

		if keyboard.IsKeyPressed(key.D) {
			textBox.AlignmentX += time.FrameDelta() * 0.3
		}
		if keyboard.IsKeyPressed(key.A) {
			textBox.AlignmentX -= time.FrameDelta() * 0.3
		}
		if keyboard.IsKeyPressed(key.S) {
			textBox.AlignmentY += time.FrameDelta() * 0.5
		}
		if keyboard.IsKeyPressed(key.W) {
			textBox.AlignmentY -= time.FrameDelta() * 0.5
		}

		textBox.AlignmentX = number.Limit(textBox.AlignmentX, 0, 1)
		textBox.AlignmentY = number.Limit(textBox.AlignmentY, 0, 1)

		// textBox.Angle += time.FrameDelta() * 10

		cam.DrawTextBoxes(textBox)
		cam.DrawTextFPS()
	}
}
