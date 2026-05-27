package main

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func main() {
	window.Create("game", false, false)
	var view = graphics.NewView(1)
	var font = assets.LoadFont("tools/msdf-atlas-gen/Libre.png", "tools/msdf-atlas-gen/Libre.json")
	// var font = assets.LoadFont("tools/msdf-atlas-gen/font.png", "tools/msdf-atlas-gen/font.json")
	var textbox = graphics.NewTextbox(-200, 0, 2000, 1500, font, "^&%#@!*_Wtyg aWAY AVATAR WAVE")
	textbox.Effects.FillColor = palette.DarkGray
	textbox.Effects.TextLineHeight = 90
	textbox.Effects.TextBackColor = palette.Red
	textbox.Effects.TextAlignX = 1
	textbox.Effects.TextAlignY = 1

	// window.SetTargetFPS(60)
	textbox.Text = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque imperdiet dignissim erat. Maecenas accumsan urna elit, ut porttitor ipsum ultricies et. Nulla vel vulputate nisl. Fusce lectus mauris, sagittis ac placerat eu, condimentum et nunc. Cras pulvinar nisl ex. Morbi et ultricies eros. Cras aliquet efficitur scelerisque. Suspendisse molestie finibus arcu, sed sagittis metus molestie a. Phasellus at fermentum massa, eget bibendum eros. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi ac lacus id enim dictum sodales at maximus tellus." + "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque imperdiet dignissim erat. Maecenas accumsan urna elit, ut porttitor ipsum ultricies et. Nulla vel vulputate nisl. Fusce lectus mauris, sagittis ac placerat eu, condimentum et nunc. Cras pulvinar nisl ex. Morbi et ultricies eros. Cras aliquet efficitur scelerisque. Suspendisse molestie finibus arcu, sed sagittis metus molestie a. Phasellus at fermentum massa, eget bibendum eros. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi ac lacus id enim dictum sodales at maximus tellus." + "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque imperdiet dignissim erat. Maecenas accumsan urna elit, ut porttitor ipsum ultricies et. Nulla vel vulputate nisl. Fusce lectus mauris, sagittis ac placerat eu, condimentum et nunc. Cras pulvinar nisl ex. Morbi et ultricies eros. Cras aliquet efficitur scelerisque. Suspendisse molestie finibus arcu, sed sagittis metus molestie a. Phasellus at fermentum massa, eget bibendum eros. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi ac lacus id enim dictum sodales at maximus tellus."
	textbox.Text += "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque imperdiet dignissim erat. Maecenas accumsan urna elit, ut porttitor ipsum ultricies et. Nulla vel vulputate nisl. Fusce lectus mauris, sagittis ac placerat eu, condimentum et nunc. Cras pulvinar nisl ex. Morbi et ultricies eros. Cras aliquet efficitur scelerisque. Suspendisse molestie finibus arcu, sed sagittis metus molestie a. Phasellus at fermentum massa, eget bibendum eros. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi ac lacus id enim dictum sodales at maximus tellus." + "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque imperdiet dignissim erat. Maecenas accumsan urna elit, ut porttitor ipsum ultricies et. Nulla vel vulputate nisl. Fusce lectus mauris, sagittis ac placerat eu, condimentum et nunc. Cras pulvinar nisl ex. Morbi et ultricies eros. Cras aliquet efficitur scelerisque. Suspendisse molestie finibus arcu, sed sagittis metus molestie a. Phasellus at fermentum massa, eget bibendum eros. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi ac lacus id enim dictum sodales at maximus tellus." + "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Pellentesque imperdiet dignissim erat. Maecenas accumsan urna elit, ut porttitor ipsum ultricies et. Nulla vel vulputate nisl. Fusce lectus mauris, sagittis ac placerat eu, condimentum et nunc. Cras pulvinar nisl ex. Morbi et ultricies eros. Cras aliquet efficitur scelerisque. Suspendisse molestie finibus arcu, sed sagittis metus molestie a. Phasellus at fermentum massa, eget bibendum eros. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Morbi ac lacus id enim dictum sodales at maximus tellus."

	for window.KeepOpen() {
		// textbox.Text = debug.MemoryUsage()
		var x, _ = view.MousePosition()
		textbox.Effects.TextLineHeight = 70 + x/10

		view.DrawObjects(&textbox)
	}
}
