package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/graphics/tag"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/text"
	"pure-game-kit/utility/time"
	"pure-game-kit/window"
)

func Texts() {
	var cam = graphics.NewCamera(1)
	var _, tiles = assets.LoadDefaultAtlasIcons()
	var textBox = graphics.NewTextBox("", 0, 0, "")
	textBox.PivotX, textBox.PivotY = 0, 0.5
	textBox.AlignmentX, textBox.AlignmentY = 0, 1
	textBox.Angle = 0
	textBox.LineHeight = 300
	textBox.Width, textBox.Height = 2000, 500
	textBox.ShadowOffsetX, textBox.ShadowOffsetY = 0.5, 0.5
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
		"exercitation ",
		tag.ShadowBlur("ullamco laboris", 3),
		" nisi ut ",
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
	textBox.WordWrap = false
	// textBox.Text = "Hello, World! Hello, World! Hello, World! Hello, World! Hello, World!"
	textBox.Angle = 0
	textBox.X = -1000

	var dir float32 = 1
	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		textBox.Tint = palette.DarkGray
		cam.DrawQuads(&textBox.Quad)
		textBox.Tint = palette.DarkGreen

		if keyboard.IsKeyPressed(key.A) {
			textBox.AlignmentX += time.FrameDelta() * 0.05 * dir
		} else {
			textBox.AlignmentX += time.FrameDelta() * 0.005 * dir
		}
		if keyboard.IsKeyPressed(key.D) {
			textBox.AlignmentX -= time.FrameDelta() * 0.05 * dir
		}

		if textBox.AlignmentX > 1 {
			dir = -1
		}
		if textBox.AlignmentX < 0 {
			dir = 1
		}

		// textBox.Angle += time.FrameDelta() * 10

		cam.DrawTextBoxes(textBox)
		cam.DrawTextFPS()
	}
}
