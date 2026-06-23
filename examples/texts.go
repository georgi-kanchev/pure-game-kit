package example

import (
	"fmt"
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/geometry"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func Texts() {
	window.Create("examples - texts", false, true)
	var font = assets.LoadFont("tools/msdf-atlas-gen/Libre.png", "tools/msdf-atlas-gen/Libre.json")
	var view = graphics.NewView(1)
	var textbox = graphics.NewTextbox(0, 0, 1600, 900, font)
	var img = assets.LoadImage("examples/data/flail.PNG")
	textbox.Effects.TextLineHeight = 100
	textbox.Effects.FillColor = palette.DarkGray
	textbox.Effects.OutlineSize = 50
	// textbox.Effects.TextAlignY = 0
	textbox.Text = "⚫⬜📢Lorem🪓🔉 ipsum dolor sit amet, 🟪consectetur🟪 ⬜adipiscing elit. Pellentesque imperdiet dignissim erat. Maecenas accumsan urna elit, ut porttitor⬜ 🟥ipsum🟥 ⬜ultricies et. ✅Nulla vel vulputate nisl.✅ Fusce lectus mauris, 🔵sagittis⚫ ac placerat eu, condimentum et nunc. ❎Cras pulvinar nisl ex.❎ Morbi et🪓 ultricies eros. \n\n📢Cras aliquet efficitur🔉 \n\nscelerisque. 🔇Suspendisse molestie🔉 finibus arcu, sed 🔼sagittis🔁 metus molestie a. 🔽Phasellus🔁 at fermentum massa, eget bibendum eros. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi ac lacus id enim dictum sodales at maximus tellus."
	textbox.Effects.TextMarginX = 100
	textbox.Effects.TextMarginY = 100
	textbox.Roundness = 0.2

	font.EmbedImage('🪓', img)

	// window.SetTargetFPS(0)

	for window.KeepOpen() {
		var x, y = view.MousePosition()
		// textbox.Effects.TextLineHeight = 70 + x/5

		fmt.Printf("textbox.TextIndexAtPoint(x, y): %v\n", textbox.TextIndexAtPoint(x, y))

		view.DrawObject(&textbox)
		view.DrawDebugInfo(false)

		for _, v := range textbox.TextSymbols {
			view.DrawShape(v.X, v.Y, v.Width, v.Height, 0, 0, color.RGBA(255, 0, 0, 100), geometry.Area{})
		}
	}
}
