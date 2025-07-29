package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/motion"
	"pure-kit/engine/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Animations() {
	var animation = motion.NewAnimation(2, true,
		"[h]", "[e]", "[l]", "[l]", "[o]", "[,]", "[space]",
		"[w]", "[o]", "[r]", "[l]", "[d]", "[!]", "")
	var cam = graphics.NewCamera(3)
	var sprite = graphics.NewSprite("", 0, 0)
	assets.LoadDefaultAtlasInput()

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		sprite.AssetId = *animation.CurrentItem()
		cam.DrawSprites(&sprite)

		if rl.IsKeyPressed(rl.KeyA) {
			animation.SetTime(3.8)
		}
	}
}
