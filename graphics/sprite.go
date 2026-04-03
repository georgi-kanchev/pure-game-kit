package graphics

import "pure-game-kit/internal"

type Sprite struct {
	Quad

	TextureId   string // "" = white
	TextureArea *Area  // nil = entire asset
}

func NewSprite(textureId string, x, y float32) *Sprite {
	var sprite = Sprite{Quad: *NewQuad(x, y), TextureId: textureId}
	var tex, has = internal.Textures[textureId]
	if has {
		sprite.Width, sprite.Height = float32(tex.Width), float32(tex.Height)
	}
	return &sprite
}
