package graphics

import (
	"image/color"
	"pure-game-kit/internal"
	col "pure-game-kit/utility/color"
	"pure-game-kit/utility/color/palette"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

//=================================================================
// objects

var reusableSprite = NewSprite("", 0, 0)
var defaultTextPack = &symbol{Color: palette.White, Weight: 1, OutlineColor: 255, OutlineWeight: 1}

func drawBoxPart(parent *Sprite, camera *Camera, x, y, w, h float32, id string, color uint) {
	reusableSprite.AssetId = id
	reusableSprite.X, reusableSprite.Y = parent.PointToGlobal(x, y)
	reusableSprite.Angle = parent.Angle
	reusableSprite.ScaleX, reusableSprite.ScaleY = parent.ScaleX, parent.ScaleY
	reusableSprite.Width, reusableSprite.Height = w, h
	reusableSprite.PivotX, reusableSprite.PivotY = 0, 0
	reusableSprite.Tint = color
	camera.DrawSprites(reusableSprite)
}

//=================================================================
// primitives

func triangulate(points []float32) []float32 {
	n := len(points) / 2
	if n < 3 {
		return nil
	}

	var triangles []float32
	var verts = make([]int, n)
	for i := 0; i < n; i++ {
		verts[i] = i
	}

	ccw := area(points) > 0
	for len(verts) > 3 {
		earFound := false
		for i := 0; i < len(verts); i++ {
			prev := verts[(i+len(verts)-1)%len(verts)]
			curr := verts[i]
			next := verts[(i+1)%len(verts)]

			if !isEar(points, verts, prev, curr, next, ccw) {
				continue
			}

			// Add triangle vertices to flat slice
			triangles = append(triangles,
				points[prev*2], points[prev*2+1],
				points[curr*2], points[curr*2+1],
				points[next*2], points[next*2+1],
			)
			verts = append(verts[:i], verts[i+1:]...)
			earFound = true
			break
		}
		if !earFound {
			break
		}
	}

	if len(verts) == 3 {
		for _, v := range verts {
			triangles = append(triangles, points[v*2], points[v*2+1])
		}
	}

	return triangles
}
func area(points []float32) float32 {
	var a float32
	n := len(points) / 2
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		a += points[i*2]*points[j*2+1] - points[j*2]*points[i*2+1]
	}
	return a / 2
}
func isEar(points []float32, verts []int, i1, i2, i3 int, ccw bool) bool {
	p1x, p1y := points[i1*2], points[i1*2+1]
	p2x, p2y := points[i2*2], points[i2*2+1]
	p3x, p3y := points[i3*2], points[i3*2+1]

	cross := (p2x-p1x)*(p3y-p1y) - (p2y-p1y)*(p3x-p1x)
	if ccw && cross <= 0 {
		return false
	}
	if !ccw && cross >= 0 {
		return false
	}

	for _, vi := range verts {
		if vi == i1 || vi == i2 || vi == i3 {
			continue
		}
		if pointInTriangle(points[vi*2], points[vi*2+1], p1x, p1y, p2x, p2y, p3x, p3y) {
			return false
		}
	}
	return true
}
func pointInTriangle(px, py, ax, ay, bx, by, cx, cy float32) bool {
	v0x, v0y := cx-ax, cy-ay
	v1x, v1y := bx-ax, by-ay
	v2x, v2y := px-ax, py-ay

	dot00 := v0x*v0x + v0y*v0y
	dot01 := v0x*v1x + v0y*v1y
	dot02 := v0x*v2x + v0y*v2y
	dot11 := v1x*v1x + v1y*v1y
	dot12 := v1x*v2x + v1y*v2y

	invDenom := 1 / (dot00*dot11 - dot01*dot01)
	u := (dot11*dot02 - dot01*dot12) * invDenom
	v := (dot00*dot12 - dot01*dot02) * invDenom

	return (u >= 0) && (v >= 0) && (u+v <= 1)
}
func isClockwiseFlat(tri []float32) bool {
	area := (tri[2]-tri[0])*(tri[5]-tri[1]) - (tri[4]-tri[0])*(tri[3]-tri[1])
	return area < 0
}

func separateShapes(points [][2]float32) (flatPoints []float32, shapeCounts []int) {
	var currentCount int
	for _, p := range points {
		if number.IsNaN(p[0]) || number.IsNaN(p[1]) {
			if currentCount >= 3 {
				shapeCounts = append(shapeCounts, currentCount)
			}
			currentCount = 0
			continue
		}
		flatPoints = append(flatPoints, p[0], p[1])
		currentCount++
	}
	if currentCount >= 3 {
		shapeCounts = append(shapeCounts, currentCount)
	}
	return flatPoints, shapeCounts
}

//=================================================================
// camera

var rlCam = rl.Camera2D{}
var maskX, maskY, maskW, maskH int

const placeholderCharAsset = '@'

// call before draw to update camera but use screen space instead of camera space
func (c *Camera) update() {
	tryRecreateWindow()

	rlCam.Target.X = float32(c.X)
	rlCam.Target.Y = float32(c.Y)
	rlCam.Rotation = float32(c.Angle)
	rlCam.Zoom = float32(c.Zoom)
	rlCam.Offset.X = float32(c.ScreenX) + float32(c.ScreenWidth/2)
	rlCam.Offset.Y = float32(c.ScreenY) + float32(c.ScreenHeight/2)

	var mx = number.Biggest(c.MaskX, c.ScreenX)
	var my = number.Biggest(c.MaskY, c.ScreenY)
	var maxW = c.ScreenX + c.ScreenWidth - mx
	var maxH = c.ScreenY + c.ScreenHeight - my
	var mw = number.Smallest(c.MaskWidth, maxW)
	var mh = number.Smallest(c.MaskHeight, maxH)

	maskX, maskY, maskW, maskH = mx, my, mw, mh
}

// call before draw to update camera and use camera space
func (c *Camera) begin() {
	c.update()
	if skipStartEnd {
		return
	}

	rl.BeginMode2D(rlCam)
	rl.BeginScissorMode(int32(maskX), int32(maskY), int32(maskW), int32(maskH))

	if c.Blend != 0 {
		rl.BeginBlendMode(rl.BlendMode(c.Blend))
	}
	if c.Effects != nil {
		rl.BeginShaderMode(internal.Shader)
		rl.EnableDepthTest()
	}
}

// call after draw to get back to using screen space
func (c *Camera) end() {
	if skipStartEnd {
		return
	}

	if c.Blend != 0 {
		rl.EndBlendMode()
	}
	if c.Effects != nil {
		rl.DisableDepthTest()
		rl.EndShaderMode()
	}

	rl.EndScissorMode()
	rl.EndMode2D()
}

//=================================================================
// other

func tryRecreateWindow() {
	if internal.WindowReady {
		return
	}

	if !rl.IsWindowReady() {
		window.Recreate()
	}
}

func getColor(value uint) color.RGBA {
	var r, g, b, a = col.Channels(value)
	return color.RGBA{R: r, G: g, B: b, A: a}
}
func packSymbolColor(s *symbol) rl.Color {
	var packLayer = func(c rl.Color) uint8 {
		var r = (c.R >> 6) & 0x03
		var g = (c.G >> 6) & 0x03
		var b = (c.B >> 6) & 0x03
		var a = (c.A >> 6) & 0x03
		return (r << 6) | (g << 4) | (b << 2) | a
	}

	var thick, out, sh, shSmooth byte = s.Weight, s.OutlineWeight, s.ShadowWeight, s.ShadowBlur
	var r = packLayer(getColor(s.Color))
	var g = packLayer(getColor(s.OutlineColor))
	var b = packLayer(getColor(s.ShadowColor))
	var a = ((thick & 0x03) << 6) | ((out & 0x03) << 4) | ((sh & 0x03) << 2) | (shSmooth & 0x03)
	return rl.NewColor(r, g, b, a)
}

func editAssetRects(src, dst *rl.Rectangle, ang float32, rotations int, flip bool) {
	if dst.Width < 0 { // raylib doesn't seem to support negative width/height???
		dst.X, dst.Y = point.MoveAtAngle(dst.X, dst.Y, ang+180, -dst.Width)
		src.Width *= -1
	}
	if dst.Height < 0 {
		dst.X, dst.Y = point.MoveAtAngle(dst.X, dst.Y, ang+270, -dst.Height)
		src.Height *= -1
	}

	if flip {
		src.Width *= -1
	}
	switch rotations % 4 {
	case 1: // 90
		dst.X, dst.Y = point.MoveAtAngle(dst.X, dst.Y, ang, dst.Height)
	case 2: // 180
		src.Height *= -1
		dst.X, dst.Y = point.MoveAtAngle(dst.X, dst.Y, ang, dst.Width)
		dst.X, dst.Y = point.MoveAtAngle(dst.X, dst.Y, ang+90, dst.Width)
	case 3: // 270
		dst.X, dst.Y = point.MoveAtAngle(dst.X, dst.Y, ang+90, dst.Width)
	}
}
func asset(assetId string) (tex *rl.Texture2D, src rl.Rectangle, rotations int, flip bool) {
	var texture, hasTexture = internal.Textures[assetId]
	src = rl.NewRectangle(0, 0, 0, 0)
	if !hasTexture {
		var rect, hasArea = internal.AtlasRects[assetId]
		if hasArea {
			var atlas, _ = internal.Atlases[rect.AtlasId]
			var tex, _ = internal.Textures[atlas.TextureId]

			texture = tex
			src.X = rect.CellX * float32(atlas.CellWidth+atlas.Gap)
			src.Y = rect.CellY * float32(atlas.CellHeight+atlas.Gap)
			src.Width = float32(atlas.CellWidth * int(rect.CountX))
			src.Height = float32(atlas.CellHeight * int(rect.CountY))
			rotations, flip = rect.Rotations, rect.Flip
		} else {
			var font, hasFont = internal.Fonts[assetId]
			if hasFont {
				texture = &font.Texture
				src.Width, src.Height = float32(texture.Width), float32(texture.Height)
			}
		}
	} else {
		src.Width, src.Height = float32(texture.Width), float32(texture.Height)
	}
	tex = texture
	return
}
