package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("game", false, true)
	var font = assets.LoadFont("tools/msdf-atlas-gen/Libre.png", "tools/msdf-atlas-gen/Libre.json")
	var view = graphics.NewView(1)
	var textbox = graphics.NewTextbox(0, 0, 2100, 2100, font)
	var img = assets.LoadImage("examples/data/flail.PNG")
	textbox.Effects.TextLineHeight = 110
	textbox.Effects.FillColor = palette.DarkGray
	textbox.Effects.OutlineSize = 50
	textbox.Text = "⚫⬜Lorem 🪓 ipsum dolor sit amet, 🟪consectetur🟪 ⬜adipiscing elit. Pellentesque imperdiet dignissim erat. Maecenas accumsan urna elit, ut porttitor⬜ 🟥ipsum🟥 ⬜ultricies et. ✅Nulla vel vulputate nisl.✅ Fusce lectus mauris, 🔵sagittis⚫ ac placerat eu, condimentum et nunc. ❎Cras pulvinar nisl ex.❎ Morbi et🪓 ultricies eros. Cras aliquet efficitur scelerisque. Suspendisse molestie finibus arcu, sed 🔼sagittis🔁 metus molestie a. 🔽Phasellus🔁 at fermentum massa, eget bibendum eros. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi ac lacus id enim dictum sodales at maximus tellus."

	font.EmbedImage('🪓', img)

	window.SetTargetFPS(60)

	for window.KeepOpen() {
		view.DrawObjects(&textbox)
	}
}
