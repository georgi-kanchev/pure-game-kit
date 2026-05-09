package example

import (
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/input/keyboard"
	"pure-game-kit/packages/input/keyboard/key"
	"pure-game-kit/packages/motion"
	"pure-game-kit/packages/window"
)

func Animations() {
	var animation = motion.NewAnimation(2, true,
		"[h]", "[e]", "[l]", "[l]", "[o]", "[,]", "[space]", "[w]", "[o]", "[r]", "[l]", "[d]", "[!]", "----")
	var view = graphics.NewView(3)
	var sprite = graphics.NewSprite(0, 0, 0)

	for window.KeepOpen() {
		animation.Update()
		view.DrawSprites(sprite)

		if keyboard.IsKeyJustPressed(key.A) {
			animation.SetIndex(5)
		}
	}
}
