package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/debug"
	"pure-game-kit/packages/utility/file"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type ImageId int32

func LoadImage(imagePath string) ImageId {
	if !file.Exists(imagePath) {
		debug.LogError("Failed to find image file: \"", imagePath, "\"")
		return 0
	}

	var texture = rl.LoadTexture(imagePath)
	if texture.Width == 0 {
		debug.LogError("Failed to load image file: \"", imagePath, "\"")
		return 0
	}
	internal.NextImageId++
	var id = internal.NextImageId
	var w, h = float32(texture.Width), float32(texture.Height)
	internal.AddPendingImage(int32(id), internal.ImageData{Texture: texture, CropWidth: w, CropHeight: h})
	return ImageId(id)
}
func (i ImageId) UnloadImage() {
	if i == 0 {
		return
	}
	var img = internal.Images[int32(i)]
	rl.UnloadTexture(img.Texture)
	delete(internal.Images, int32(i))
}

func (i ImageId) SetSmoothness(smooth bool) {
	if i == 0 {
		return
	}

	var img, has = internal.Images[int32(i)]
	if has && smooth {
		rl.SetTextureFilter(img.Texture, rl.FilterBilinear)
	}
	if has && !smooth {
		rl.SetTextureFilter(img.Texture, rl.FilterPoint)
	}
}
