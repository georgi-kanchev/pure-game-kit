package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/graphics"
	"pure-game-kit/input/keyboard"
	"pure-game-kit/input/keyboard/key"
	"pure-game-kit/motion"
	"pure-game-kit/window"
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
		cam.DrawSprites(sprite)

		if keyboard.IsKeyPressedOnce(key.A) {
			animation.SetTime(3.8)
		}
	}
}
