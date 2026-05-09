package graphics

import (
	"pure-game-kit/packages/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (v *View) DrawQuads(quads ...*Quad) {
	v.begin()

	var lastEffects *Effects
	for _, q := range quads {
		if q == nil || !v.IsAreaVisible(q.Bounds()) {
			continue
		}

		var x, y = q.PointFromEdge(0, 0) // applying pivot
		var ang = q.Angle
		var src = rl.NewRectangle(0, 0, 1, 1)
		var dst = rl.NewRectangle(x, y, q.Width*q.ScaleX, q.Height*q.ScaleY)
		var mask = v.Mask
		var effects = q.Effects
		if q.Mask != (Area{}) {
			mask = q.Mask
		}
		if effects == nil {
			effects = v.Effects
		}
		if lastEffects != effects {
			batcher.Draw() // effects are different & break the batch
			effects.updateUniforms(int(src.Width), int(src.Height), nil, nil, false)
		}
		batcher.QueueTexture(internal.White1x1, src, dst, ang, getColor(q.Color), mask)
		lastEffects = effects
	}
	batcher.Draw()
	v.end()
}
func (v *View) DrawSprites(sprites ...*Sprite) {
	v.begin()

	var lastEffects *Effects
	for _, s := range sprites {
		if s == nil || !v.IsAreaVisible(s.Bounds()) {
			continue
		}

		var img = internal.Images[int32(s.ImageId)]
		var src = rl.NewRectangle(img.CropX, img.CropY, img.CropWidth, img.CropHeight)

		if s.ImageCropArea != (Area{}) {
			src.Width, src.Height = s.ImageCropArea.Width, s.ImageCropArea.Height
			src.X, src.Y = s.ImageCropArea.X, s.ImageCropArea.Y
		}
		var x, y = s.PointFromEdge(0, 0) // applying pivot
		var ang = s.Angle
		var dst = rl.NewRectangle(x, y, s.Width*s.ScaleX, s.Height*s.ScaleY)

		var mask = v.Mask
		var effects = s.Effects
		if s.Mask != (Area{}) {
			mask = s.Mask
		}
		if effects == nil {
			effects = v.Effects
		}
		if lastEffects != effects {
			batcher.Draw() // effects are different & break the batch
			effects.updateUniforms(int(src.Width), int(src.Height), nil, nil, false)
		}
		batcher.QueueTexture(img.Texture, src, dst, ang, getColor(s.Color), mask)
		lastEffects = effects
	}
	batcher.Draw()
	v.end()
}
func (v *View) DrawNinePatches(ninePatches ...*NinePatch) {
	// v.begin()
	// defer v.end()

	// skipStartAndEnd = true
	// var lastEffects *Effects
	// for _, n := range ninePatches {
	// 	if n == nil || !v.IsAreaVisible(n.Bounds()) {
	// 		continue
	// 	}

	// 	var w, h = n.Width, n.Height
	// 	var assetIds, has = internal.Boxes[n.BoxId]

	// 	if !has { // fallback to quad if no 9-slice exists
	// 		v.DrawQuads(&n.Quad)
	// 		continue
	// 	}

	// 	var _, uh = internal.AssetSize(assetIds[1])
	// 	var lw, _ = internal.AssetSize(assetIds[3])
	// 	var rw, _ = internal.AssetSize(assetIds[4])
	// 	var _, dh = internal.AssetSize(assetIds[6])
	// 	var u = condition.If(assetIds[1] == "", 0, float32(uh)) * n.EdgeScale
	// 	var l = condition.If(assetIds[3] == "", 0, float32(lw)) * n.EdgeScale
	// 	var r = condition.If(assetIds[4] == "", 0, float32(rw)) * n.EdgeScale
	// 	var d = condition.If(assetIds[6] == "", 0, float32(dh)) * n.EdgeScale
	// 	var errX, errY float32 = 2, 2
	// 	var col = getColor(n.Color)

	// 	if w < 0 {
	// 		r *= -1
	// 		l *= -1
	// 		errX *= -1
	// 	}
	// 	if h < 0 {
	// 		u *= -1
	// 		d *= -1
	// 		errY *= -1
	// 	}

	// 	if number.IsBetween(w, -(l + r), l+r, false, false) {
	// 		var total = l + r
	// 		if total != 0 {
	// 			var scale = w / total
	// 			l *= scale
	// 			r *= scale
	// 		}
	// 	}
	// 	if number.IsBetween(h, -(u + d), u+d, false, false) {
	// 		var total = u + d
	// 		if total != 0 {
	// 			var scale = h / total
	// 			u *= scale
	// 			d *= scale
	// 		}
	// 	}

	// 	var mask = v.Mask
	// 	if n.Mask != (Area{}) {
	// 		mask = n.Mask
	// 	}
	// 	var effects = n.Effects
	// 	if effects == nil {
	// 		effects = v.Effects
	// 	}

	// 	var parts = [9]struct {
	// 		x, y, w, h float32
	// 		id         string
	// 	}{
	// 		{0, 0, l, u, assetIds[0]},                                                 // Top Left
	// 		{l, 0, w - l - r, u, assetIds[1]},                                         // Top
	// 		{w - r, 0, r, u, assetIds[2]},                                             // Top Right
	// 		{0, u, l, h - u - d, assetIds[3]},                                         // Left
	// 		{l - errX/2, u - errY/2, w - l - r + errX, h - u - d + errY, assetIds[4]}, // Center
	// 		{w - r, u, r, h - u - d, assetIds[5]},                                     // Right
	// 		{0, h - d, l, d, assetIds[6]},                                             // Bottom Left
	// 		{l, h - d, w - l - r, d, assetIds[7]},                                     // Bottom
	// 		{w - r, h - d, r, d, assetIds[8]},                                         // Bottom Right
	// 	}

	// 	var ang = n.Angle

	// 	for _, p := range parts {
	// 		var texture, src, rotations, flip = internal.AssetData(p.id)
	// 		if texture.Width == 0 {
	// 			texture = internal.White1x1
	// 		}

	// 		var globalX, globalY = n.PointToGlobal(p.x, p.y)
	// 		var dst = rl.NewRectangle(globalX, globalY, p.w*n.ScaleX, p.h*n.ScaleY)
	// 		var partAng = ang

	// 		internal.EditAssetRects(&src, &dst, partAng, rotations, flip)
	// 		partAng += float32(rotations * 90)

	// 		if lastEffects != effects {
	// 			batcher.Draw() // effects are different & break the batch
	// 			effects.updateUniforms(int(src.Width), int(src.Height), nil, nil, false)
	// 		}
	// 		batcher.QueueTexture(texture, src, dst, partAng, col, mask)
	// 		lastEffects = effects
	// 	}
	// }

	// batcher.Draw()
	// skipStartAndEnd = false
}
func (v *View) DrawTextBoxes(textBoxes ...*TextBox) {
	v.begin()
	for _, t := range textBoxes {
		if t == nil || !v.IsAreaVisible(t.Bounds()) {
			continue
		}

		var _, symbols = t.formatSymbols()
		var font = t.font()
		var gapX = t.gapSymbols() * t.ScaleX
		var effects = t.Effects
		var mask = v.Mask
		if t.Mask != (Area{}) {
			mask = t.Mask
		}
		if effects == nil {
			effects = v.Effects
		}
		effects.updateUniforms(int(font.Texture.Width), int(font.Texture.Height), nil, t, false)

		for _, s := range symbols {
			s.Rect.X, s.Rect.Y = t.PointToGlobal(s.Rect.X, s.Rect.Y)
			s.Rect.Width *= t.ScaleX
			s.Rect.Height *= t.ScaleY
			s.Bounds.X, s.Bounds.Y = t.PointToGlobal(s.Bounds.X, s.Bounds.Y)
			s.Bounds.Width *= t.ScaleX
			s.Bounds.Height *= t.ScaleY
			s.Angle = t.Angle
			batcher.QueueSymbol(font, s, t.LineHeight*t.ScaleY, gapX, mask)
		}
	}
	batcher.Draw()
	v.end()
}
func (v *View) DrawTileMaps(tileMaps ...*TileMap) {
	v.begin()
	for _, t := range tileMaps {
		if !v.IsAreaVisible(t.Bounds()) {
			continue
		}

		var atlas = internal.TileSets[t.TileSetId]
		var data = internal.TileLayers[t.TileLayerId]
		if atlas == nil && data == nil {
			continue
		}
		if data.Texture == nil {
			continue // only object points, no tile data
		}

		var texture = internal.Images[atlas.ImageId].Texture
		var x, y = t.PointFromEdge(0, 0) // applying pivot
		var src = rl.NewRectangle(0, 0, float32(texture.Width), float32(texture.Height))
		var dst = rl.NewRectangle(x, y, t.Width*t.ScaleX, t.Height*t.ScaleY)
		var mask = v.Mask
		if t.Mask != (Area{}) {
			mask = t.Mask
		}
		var effects = t.Effects
		if effects == nil {
			effects = v.Effects
		}
		effects.updateUniforms(int(texture.Width), int(texture.Height), t, nil, false)
		batcher.QueueTexture(texture, src, dst, t.Angle, getColor(t.Color), mask)
		batcher.Draw()
	}
	v.end()
}
func (v *View) DrawObjects(objects ...*Object) {
	v.begin()
	for _, t := range objects {
		if t == nil || !v.IsAreaVisible(t.Bounds()) {
			continue
		}

		t.tryRegenerateText()
		// for _, s := range t.chars {

		// }
	}
	batcher.Draw()
	v.end()
}
