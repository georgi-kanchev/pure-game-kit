package graphics

import (
	"pure-game-kit/packages/assets"
	"pure-game-kit/packages/internal"
)

type Sprite struct {
	Quad

	ImageId       assets.ImageId
	ImageCropArea Area // Zero value = entire asset
}

func NewSprite(imageId assets.ImageId, x, y float32) *Sprite {
	var sprite = Sprite{Quad: *NewQuad(x, y), ImageId: imageId}
	var img, has = internal.Images[int32(imageId)]
	if has {
		sprite.Width, sprite.Height = float32(img.Texture.Width), float32(img.Texture.Height)
	}
	return &sprite
}
