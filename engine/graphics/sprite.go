package graphics

import "pure-kit/engine/internal"

type Sprite struct {
	Node
	RepeatX, RepeatY float32
}

func NewSprite(assetId string) Sprite {
	var sprite = Sprite{Node: NewNode(assetId, 0, 0), RepeatX: 1, RepeatY: 1}
	var tex, has = internal.Textures[assetId]
	if has {
		sprite.Width, sprite.Height = float32(tex.Width), float32(tex.Height)
	}

	return sprite
}
