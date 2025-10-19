package graphics

import "pure-game-kit/internal"

type Sprite struct {
	Node
	TextureRepeat                  bool
	TextureScrollX, TextureScrollY float32
	AssetId                        string
}

func NewSprite(assetId string, x, y float32) Sprite {
	var sprite = Sprite{Node: NewNode(x, y), AssetId: assetId, TextureRepeat: false}
	var tex, has = internal.Textures[assetId]
	if has {
		sprite.Width, sprite.Height = float32(tex.Width), float32(tex.Height)
	}

	return sprite
}
