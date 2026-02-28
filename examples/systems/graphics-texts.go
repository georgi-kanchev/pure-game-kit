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
	var _, tiles = assets.LoadDefaultAtlasIcons()
	var textBox = graphics.NewTextBox("", 0, 0, "")
	textBox.PivotX, textBox.PivotY = 0.5, 0.5
	textBox.AlignmentX, textBox.AlignmentY = 0.5, 1
	textBox.Angle = 0
	textBox.LineHeight = 120
	textBox.Width, textBox.Height = 2020, 1500
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

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		textBox.Tint = palette.DarkGray
		cam.DrawNodes(&textBox.Node)
		textBox.Tint = palette.White

		cam.DrawTextBoxes(textBox)
	}
}
