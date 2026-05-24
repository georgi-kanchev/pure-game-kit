package graphics

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/color/palette"
)

type Effects internal.Effects

func NewEffects() *Effects {
	return &Effects{TextColor: palette.White, TextShadowColor: palette.Black, TextLineHeight: 40, TextWordWrap: true}
}
