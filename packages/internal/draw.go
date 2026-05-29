package internal

import (
	_ "embed"
	"image/color"
	col "pure-game-kit/packages/utility/color"
	"pure-game-kit/packages/utility/color/palette"
	"pure-game-kit/packages/utility/number"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Area struct{ X, Y, Width, Height float32 }
type Batch struct {
	mesh         *rl.Mesh
	meshUploaded bool
	material     rl.Material

	vertCount, indexCount int32

	verts, texCoords, normals, cols, tangents, tex2s, indexes []byte
}
type Vertex struct {
	X, Y, U, V     float32
	NX, NY, NZ     float32 // Normals
	TX, TY, TZ, TW float32 // Tangents
	U2, V2         float32 // Texcoords2
}
type Effects struct {
	Gamma, Saturation, Contrast, Brightness int8 // Ranged -128..127, 0 = no change

	OutlineSize, BorderSize float32

	Tint, BorderColor uint

	FillColor    uint // Not used by Texts.
	OutlineColor uint // Not used by Shapes.

	PixelSize    uint8 // Ranged 0..15; Not used by Shapes & Texts.
	BlurX, BlurY uint8 // Ranged 0..31; Not used by Shapes & Texts.

	// Ranged 0..1.
	//
	// Requires semi-transparent pixels to be drawn last to avoid artifacts. Fully opaque/transparent pixels work in any sorting.
	DepthZ float32

	//=================================================================

	TextAlignX, TextAlignY                     float32 // Ranged 0..1
	TextLineHeight, TextSymbolGap, TextLineGap float32
	TextWordWrap, TextUnderline, TextCrossout  bool
	TextWeight, TextShadowWeight,
	TextShadowOffsetX, TextShadowOffsetY int8
	TextShadowBlur                            uint8
	TextColor, TextBackColor, TextShadowColor uint
}

var Shader rl.Shader
var ShaderLoc int32 // uniform location, all properties are packed in one uniform for speed
var ShaderTileDataLoc int32
var DefaultMaterial rl.Material
var DefaultMatrix rl.Matrix
var DefaultEffects = Effects{BorderColor: palette.White, Tint: palette.White,
	TextColor: palette.White, TextShadowColor: palette.Black, TextShadowOffsetX: 30, TextShadowOffsetY: 30,
	TextLineHeight: 40, TextWordWrap: true}

var Images = make(map[int32]ImageData) // negative = crops; 0 = Font+White1x1; positive = full images
var NextImageId int16
var NextImageCropId int16

var ActiveBatch *Batch    // the batch currently being written to
var ReadyBatches []*Batch // batches ready to be drawn
var BatchPool []*Batch    // empty batches ready to be reused

//=================================================================

func Queue(tex rl.Texture2D, src, dst rl.Rectangle, ang, round float32, mask Area, eff *Effects, kind uint8) {
	dst.Width, dst.Height = number.Absolute(dst.Width), number.Absolute(dst.Height)

	var invTexW, invTexH = 1.0 / float32(tex.Width), 1.0 / float32(tex.Height)
	var u1, v1, u2, v2 = src.X * invTexW, src.Y * invTexH, (src.X + src.Width) * invTexW, (src.Y + src.Height) * invTexH
	var ww, wh = float32(WindowWidth) / 2, float32(WindowHeight) / 2
	var dx, dy = [4]float32{ww, ww, dst.Width + ww, dst.Width + ww}, [4]float32{wh, dst.Height + wh, dst.Height + wh, wh}
	var uvs = [8]float32{u1, v1, u1, v2, u2, v2, u2, v1}
	var vCount int32
	if eff == nil {
		eff = &DefaultEffects
	}

	var padU, padV float32
	if kind != KindText && eff.BorderSize > 0 { // no border padding for text symbols
		var padX, padY = eff.BorderSize, eff.BorderSize
		padU, padV = eff.BorderSize*(u2-u1)/dst.Width, eff.BorderSize*(v2-v1)/dst.Height
		dx[0], dx[1], dx[2], dx[3] = dx[0]-padX, dx[1]-padX, dx[2]+padX, dx[3]+padX
		dy[0], dy[3], dy[1], dy[2] = dy[0]-padY, dy[3]-padY, dy[1]+padY, dy[2]+padY
		uvs[0], uvs[2], uvs[4], uvs[6] = uvs[0]-padU, uvs[2]-padU, uvs[4]+padU, uvs[6]+padU
		uvs[1], uvs[7], uvs[3], uvs[5] = uvs[1]-padV, uvs[7]-padV, uvs[3]+padV, uvs[5]+padV
	}

	var bx, by = uint8(number.Limit(eff.BlurX, 0, 31)), uint8(number.Limit(eff.BlurY, 0, 31))
	var ps, oc = number.Limit(eff.PixelSize, 0, 16), col.Tint(eff.OutlineColor, eff.Tint)
	var r, g, b, a = col.Channels(palette.White)
	var cropMinU, cropMaxU, cropMinV, cropMaxV = u1 - 0.5, u2 - 0.5, v1 - 0.5, v2 - 0.5
	var u, v = packU2(uint16(src.Width), uint16(src.Height)), packV2(col.Tint(eff.BorderColor, eff.Tint))
	var nx = packNormalX(eff.Gamma, eff.Saturation, eff.Contrast, eff.Brightness)
	var ny, nz = packNormalY(round, ps, bx, by), packNormalZ(eff.DepthZ, eff.BorderSize, kind)
	var tx, ty, tz, tw float32
	var os = uint8(number.Limit(eff.OutlineSize, 0, 255))
	switch kind {
	case KindText:
		var w, sc = eff.TextWeight, col.Tint(eff.TextShadowColor, eff.Tint)
		var ss, sb, sx, sy = eff.TextShadowWeight, eff.TextShadowBlur, eff.TextShadowOffsetX, eff.TextShadowOffsetY
		tx, ty, tz, tw = packTangentXText(oc), packTangentYText(sc), packTangentZText(w, os, ss), packTangentWText(sx, sy, sb)
		r, g, b, a = col.Channels(col.Tint(eff.TextColor, eff.Tint))
	case KindSprite, KindShape:
		tx, ty = packTangentXSprite(oc), packTangentYSprite(os, col.Tint(eff.FillColor, eff.Tint))
		tz, tw = packTangentZSprite(cropMinU, cropMaxU), packTangentWSprite(cropMinV, cropMaxV)
		if kind == KindShape {
			r, g, b, a = col.Channels(col.Tint(eff.FillColor, eff.Tint))
		}
	default:
		tx, ty = packColor24(oc), packColor24(col.Tint(eff.FillColor, eff.Tint))
	}

	for i := range len(polygonBuf) {
		polygonBuf[i].U2, polygonBuf[i].V2 = u, v
		polygonBuf[i].NX, polygonBuf[i].NY, polygonBuf[i].NZ = nx, ny, nz
		polygonBuf[i].TX, polygonBuf[i].TY, polygonBuf[i].TZ, polygonBuf[i].TW = tx, ty, tz, tw
	}

	var finalColor = color.RGBA{R: r, G: g, B: b, A: a}
	if mask != (Area{}) {
		var sinA, cosA = SinCos(ang)
		var cx, cy = ww + dst.Width/2, wh + dst.Height/2
		for i := range 4 {
			var rx, ry = dx[i] - cx, dy[i] - cy
			polygonBuf[i].X, polygonBuf[i].Y = (rx*cosA-ry*sinA)+cx+dst.X, (rx*sinA+ry*cosA)+cy+dst.Y
			polygonBuf[i].U, polygonBuf[i].V = uvs[i*2], uvs[i*2+1]
		}
		vCount = clipPolygonAABB(polygonBuf[:4], clipResultBuf[:], clipTempBuf[:], mask)
		if vCount >= 3 {
			queueVertices(clipResultBuf[:vCount], vCount, tex, finalColor)
		}
		return
	}

	vCount = 4
	if ang == 0 {
		for i := range 4 {
			polygonBuf[i].X, polygonBuf[i].Y = dx[i]+dst.X, dy[i]+dst.Y
			polygonBuf[i].U, polygonBuf[i].V = uvs[i*2], uvs[i*2+1]
		}
	} else {
		var sinA, cosA = SinCos(ang)
		var cx, cy = ww + dst.Width/2, wh + dst.Height/2
		for i := range 4 {
			var rx, ry = dx[i] - cx, dy[i] - cy
			polygonBuf[i].X, polygonBuf[i].Y = (rx*cosA-ry*sinA)+cx+dst.X, (rx*sinA+ry*cosA)+cy+dst.Y
			polygonBuf[i].U, polygonBuf[i].V = uvs[i*2], uvs[i*2+1]
		}
	}
	queueVertices(polygonBuf[:4], vCount, tex, finalColor)
}

func ResetBatches() {
	if ActiveBatch != nil {
		BatchPool = append(BatchPool, ActiveBatch)
		ActiveBatch = nil
	}
	for _, rb := range ReadyBatches {
		BatchPool = append(BatchPool, rb)
	}
	ReadyBatches = ReadyBatches[:0]
}
func CloseBatch() {
	if ActiveBatch != nil && ActiveBatch.vertCount > 0 {
		ReadyBatches = append(ReadyBatches, ActiveBatch)
		ActiveBatch = nil
	}
}
func Draw() {
	for _, batch := range ReadyBatches {
		if !batch.meshUploaded {
			rl.UploadMesh(batch.mesh, true)
			batch.meshUploaded = true
		}
		rl.UpdateMeshBuffer(*batch.mesh, 0, batch.verts[:batch.vertCount*12], 0)
		rl.UpdateMeshBuffer(*batch.mesh, 1, batch.texCoords[:batch.vertCount*8], 0)
		rl.UpdateMeshBuffer(*batch.mesh, 2, batch.normals[:batch.vertCount*12], 0)
		rl.UpdateMeshBuffer(*batch.mesh, 3, batch.cols[:batch.vertCount*4], 0)
		rl.UpdateMeshBuffer(*batch.mesh, 4, batch.tangents[:batch.vertCount*16], 0)
		rl.UpdateMeshBuffer(*batch.mesh, 5, batch.tex2s[:batch.vertCount*8], 0)
		rl.UpdateMeshBuffer(*batch.mesh, 6, batch.indexes[:batch.indexCount*2], 0)
		batch.mesh.TriangleCount = batch.indexCount / 3
		rl.DrawMesh(*batch.mesh, batch.material, DefaultMatrix)
	}
}

// private =================================================================

var polygonBuf, clipResultBuf, clipTempBuf [12]Vertex // reused working buffers; avoids per-call heap escapes

//go:embed shader.frag
var shaderFrag string

//go:embed shader.vert
var shaderVert string

func newBatch() *Batch {
	const quadCapacity = 4096 // fixed size for all batches

	var b = &Batch{}
	b.mesh = &rl.Mesh{VertexCount: 4 * quadCapacity, TriangleCount: 2 * quadCapacity}

	b.verts = make([]byte, b.mesh.VertexCount*3*4)
	b.texCoords = make([]byte, b.mesh.VertexCount*2*4)
	b.normals = make([]byte, b.mesh.VertexCount*3*4)
	b.cols = make([]byte, b.mesh.VertexCount*4)
	b.tangents = make([]byte, b.mesh.VertexCount*4*4)
	b.tex2s = make([]byte, b.mesh.VertexCount*2*4)
	b.indexes = make([]byte, b.mesh.TriangleCount*3*2)

	b.mesh.Vertices = (*float32)(unsafe.Pointer(&b.verts[0]))
	b.mesh.Texcoords = (*float32)(unsafe.Pointer(&b.texCoords[0]))
	b.mesh.Normals = (*float32)(unsafe.Pointer(&b.normals[0]))
	b.mesh.Colors = (*uint8)(unsafe.Pointer(&b.cols[0]))
	b.mesh.Tangents = (*float32)(unsafe.Pointer(&b.tangents[0]))
	b.mesh.Texcoords2 = (*float32)(unsafe.Pointer(&b.tex2s[0]))
	b.mesh.Indices = (*uint16)(unsafe.Pointer(&b.indexes[0]))

	b.material = DefaultMaterial
	b.material.Maps = &rl.MaterialMap{
		Texture: DefaultMaterial.Maps.Texture,
		Color:   DefaultMaterial.Maps.Color,
		Value:   DefaultMaterial.Maps.Value,
	}
	b.material.Shader = Shader
	return b
}
func queueVertices(verts []Vertex, vCount int32, tex rl.Texture2D, col rl.Color) {
	if ActiveBatch != nil {
		var texChanged = ActiveBatch.material.Maps.Texture.ID != tex.ID
		var outOfSpace = ActiveBatch.vertCount+vCount > ActiveBatch.mesh.VertexCount

		if texChanged || outOfSpace { // do we need to break the batch?
			if ActiveBatch.vertCount > 0 {
				ReadyBatches = append(ReadyBatches, ActiveBatch) // push to draw later
			}
			ActiveBatch = nil // clear active to force a new one
		}
	}

	if ActiveBatch == nil {
		if len(BatchPool) > 0 { // grab a fresh batch if we don't have an active one
			ActiveBatch = BatchPool[len(BatchPool)-1]
			BatchPool = BatchPool[:len(BatchPool)-1]
		} else {
			ActiveBatch = newBatch() // pool is empty, allocate a new one (will only happen as the game ramps up)
		}

		ActiveBatch.vertCount = 0 // reset counters and set material
		ActiveBatch.indexCount = 0
		ActiveBatch.material.Maps.Texture = tex
		ActiveBatch.material.Shader = Shader
	}

	// write data to the active batch
	var b = ActiveBatch
	var count = int32(len(verts))
	var v_slice = unsafe.Slice((*float32)(unsafe.Pointer(&b.verts[b.vertCount*12])), count*3)
	var t_slice = unsafe.Slice((*float32)(unsafe.Pointer(&b.texCoords[b.vertCount*8])), count*2)
	var n_slice = unsafe.Slice((*float32)(unsafe.Pointer(&b.normals[b.vertCount*12])), count*3)
	var c_slice = b.cols[b.vertCount*4 : (b.vertCount*4)+(count*4)]
	var tan_slice = unsafe.Slice((*float32)(unsafe.Pointer(&b.tangents[b.vertCount*16])), count*4)
	var t2_slice = unsafe.Slice((*float32)(unsafe.Pointer(&b.tex2s[b.vertCount*8])), count*2)

	for i, v := range verts {
		v_slice[i*3+0], v_slice[i*3+1], v_slice[i*3+2] = v.X, v.Y, 0
		t_slice[i*2+0], t_slice[i*2+1] = v.U, v.V
		n_slice[i*3+0], n_slice[i*3+1], n_slice[i*3+2] = v.NX, v.NY, v.NZ
		c_slice[i*4+0], c_slice[i*4+1], c_slice[i*4+2], c_slice[i*4+3] = col.R, col.G, col.B, col.A
		tan_slice[i*4+0], tan_slice[i*4+1], tan_slice[i*4+2], tan_slice[i*4+3] = v.TX, v.TY, v.TZ, v.TW
		t2_slice[i*2+0], t2_slice[i*2+1] = v.U2, v.V2
	}

	var trisCount = count - 2
	var indSlice = unsafe.Slice((*uint16)(unsafe.Pointer(&b.indexes[b.indexCount*2])), trisCount*3)
	var base = uint16(b.vertCount)

	for i := range trisCount {
		indSlice[i*3+0] = base
		indSlice[i*3+1] = base + uint16(i+1)
		indSlice[i*3+2] = base + uint16(i+2)
	}

	b.vertCount += count
	b.indexCount += trisCount * 3
}

func clipPolygonAABB(poly, outBuf, tempBuf []Vertex, mask Area) int32 {
	var minX, maxX, minY, maxY = mask.X, mask.X + mask.Width, mask.Y, mask.Y + mask.Height
	var count = clipPolyEdge(poly, tempBuf, true, minX, true)
	if count == 0 {
		return 0
	}
	count = clipPolyEdge(tempBuf[:count], outBuf, true, maxX, false)
	if count == 0 {
		return 0
	}
	count = clipPolyEdge(outBuf[:count], tempBuf, false, minY, true)
	if count == 0 {
		return 0
	}
	count = clipPolyEdge(tempBuf[:count], outBuf, false, maxY, false)
	return count
}
func clipPolyEdge(in, out []Vertex, isX bool, edgeVal float32, keepGreater bool) int32 {
	var outCount int32 = 0
	if len(in) == 0 {
		return 0
	}

	var prev = in[len(in)-1]
	var prevVal float32
	if isX {
		prevVal = prev.X
	} else {
		prevVal = prev.Y
	}
	var prevInside = (keepGreater && prevVal >= edgeVal) || (!keepGreater && prevVal <= edgeVal)

	for _, curr := range in {
		var currVal float32
		if isX {
			currVal = curr.X
		} else {
			currVal = curr.Y
		}
		var currInside = (keepGreater && currVal >= edgeVal) || (!keepGreater && currVal <= edgeVal)

		if currInside != prevInside {
			var t float32
			if isX {
				t = (edgeVal - prev.X) / (curr.X - prev.X)
			} else {
				t = (edgeVal - prev.Y) / (curr.Y - prev.Y)
			}

			out[outCount] = prev
			out[outCount].X, out[outCount].Y = prev.X+t*(curr.X-prev.X), prev.Y+t*(curr.Y-prev.Y)
			out[outCount].U, out[outCount].V = prev.U+t*(curr.U-prev.U), prev.V+t*(curr.V-prev.V)
			outCount++
		}

		if currInside {
			out[outCount] = curr
			outCount++
		}

		prev = curr
		prevInside = currInside
	}

	return outCount
}
