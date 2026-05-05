package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/window"
)

func Animations() {
	var animation = motion.NewAnimation(2, true,
		"[h]", "[e]", "[l]", "[l]", "[o]", "[,]", "[space]", "[w]", "[o]", "[r]", "[l]", "[d]", "[!]", "----")
	var cam = graphics.NewCamera(3)
	var sprite = graphics.NewSprite("", 0, 0)
	assets.LoadDefaultAtlasInput()

	for window.KeepOpen() {
		animation.Update()
		sprite.TextureId = animation.Item()
		cam.DrawSprites(sprite)

		if keyboard.IsKeyJustPressed(key.A) {
			animation.SetIndex(5)
		}
	}
}
