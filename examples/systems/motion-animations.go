package example

import (
	"pure-kit/engine/data/assets"
	"pure-kit/engine/graphics"
	"pure-kit/engine/input/keyboard"
	"pure-kit/engine/input/keyboard/key"
	"pure-kit/engine/motion"
	"pure-kit/engine/window"
)

func Animations() {
	var animation = motion.NewAnimation(2, true,
		"[h]", "[e]", "[l]", "[l]", "[o]", "[,]", "[space]",
		"[w]", "[o]", "[r]", "[l]", "[d]", "[!]", "")
	var cam = graphics.NewCamera(3)
	var sprite = graphics.NewSprite("", 0, 0)
	assets.LoadDefaultAtlasInput(true)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		sprite.AssetId = *animation.CurrentItem()
		cam.DrawSprites(&sprite)

		if keyboard.IsKeyPressedOnce(key.A) {
			animation.SetTime(3.8)
		}
	}
}
