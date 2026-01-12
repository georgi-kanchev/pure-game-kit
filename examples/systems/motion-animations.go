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
		"[w]", "[o]", "[r]", "[l]", "[d]", "[!]", "----")
	var cam = graphics.NewCamera(3)
	var sprite = graphics.NewSprite("", 0, 0)
	assets.LoadDefaultAtlasInput()

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()

		animation.Update()
		sprite.AssetId = *animation.Item()
		cam.DrawSprites(sprite)

		if keyboard.IsKeyJustPressed(key.A) {
			animation.SetIndex(5)
		}
	}
}
