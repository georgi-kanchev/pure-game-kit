package graphics

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/color/palette"
)

type Effects internal.Effects

func NewEffects() *Effects {
	return &Effects{Gamma: 0.5, Saturation: 0.5, Contrast: 0.5, Brightness: 0.5,
		TextColor: palette.White, TextShadowColor: palette.Black, TextLineHeight: 40, TextWordWrap: true}
}
