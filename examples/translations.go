package example

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/graphics"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/window"
)

func Translations() {
	window.Create("example - translations", true, true)
	var view = graphics.NewView(1)

	var lang = assets.LoadTranslations("examples/data/english.yaml")
	var bulgarian = assets.LoadTranslations("examples/data/bulgarian.yaml")
	var tag = "intro_cutscene_dialog"

	bulgarian.Unload()

	for window.KeepOpen() {
		view.DrawText(0, 0, 100, 0, palette.White, lang.Translate(tag))
		// view.DrawText(0, 100, 100, 0, palette.White, bulgarian.Translate(tag))
	}
}
