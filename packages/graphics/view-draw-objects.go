package graphics

import (
	"pure-game-kit/packages/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (v *View) DrawObjects(objects ...*Object) {
	for _, o := range objects {
		if o == nil || !v.IsAreaVisible(o.Bounds()) {
			continue
		}

		var tex = internal.Images[int32(o.ImageId)]
		var src = rl.NewRectangle(tex.CropX, tex.CropY, tex.CropWidth, tex.CropHeight)
		var dst = rl.NewRectangle(o.X, o.Y, o.Width*o.ScaleX, o.Height*o.ScaleY)
		internal.ActiveBatchManager.QueueTexture(tex.Texture, src, dst, o.Angle, getColor(o.Color), internal.Area(o.Mask))

		o.tryRegenerateText()
		// for _, s := range t.chars {

		// }
	}
}
