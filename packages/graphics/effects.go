package graphics

import "pure-game-kit/packages/internal"

type Effects internal.Effects

func NewEffects() *Effects {
	return &Effects{Gamma: 0.5, Saturation: 0.5, Contrast: 0.5, Brightness: 0.5}
}
