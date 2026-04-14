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
	var frames = [14]string{"[h]", "[e]", "[l]", "[l]", "[o]", "[,]", "[space]", "[w]", "[o]", "[r]", "[l]", "[d]", "[!]", "----"}
	var animation = motion.NewAnimation(len(frames), 2, true)
	var cam = graphics.NewCamera(3)
	var sprite = graphics.NewSprite("", 0, 0)
	assets.LoadDefaultAtlasInput()

	for window.KeepOpen() {
		animation.Update()
		sprite.TextureId = frames[animation.Index()]
		cam.DrawSprites(sprite)

		if keyboard.IsKeyJustPressed(key.A) {
			animation.SetIndex(5)
		}
	}
}
