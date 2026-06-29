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
type ImageData struct {
	Texture rl.Texture2D

	CropX, CropY, CropWidth, CropHeight,
	Top, Left, Right, Bottom float32 // edge offsets for 9patch
}
type Batch struct {
	mesh     *rl.Mesh
	material rl.Material
	meshUploaded,
	isRecord, IsMeshDirty bool

	vertCount, indexCount int32

	verts, texCoords, normals, cols, tangents, tex2s, indexes []byte

	tileDataTex rl.Texture2D // tile layer texture (set only for KindTilemap batches)
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

	FillColor    uint
	OutlineColor uint // Not used by Shapes.

	PixelSize    uint8 // Ranged 0..15; Not used by Shapes & Texts.
	BlurX, BlurY uint8 // Ranged 0..31; Not used by Shapes & Texts.

	// Ranged 0..1.
	//
	// Requires semi-transparent pixels to be drawn last to avoid artifacts. Fully opaque/transparent pixels work in any sorting.
	DepthZ float32

	//=================================================================

	TextAlignX, TextAlignY float32 // Ranged 0..1
	TextLineHeight, TextSymbolGap, TextLineGap,
	TextMarginX, TextMarginY float32
	TextWordWrap bool

	TextIsInput bool // No new lines; no effects; caches the cursor positions from the last draw.

	// Caches the text visuals across frames. Call object.TextUpdateBatch when visual changes are needed.
	// Useful for a huge static texts that change rarely.
	TextBatch bool

	TextUnderline, TextCrossout bool
	TextWeight, TextShadowWeight,
	TextShadowOffsetX, TextShadowOffsetY int8
	TextShadowBlur             uint8
	TextColor, TextShadowColor uint
}

var Shader rl.Shader
var ShaderLoc int32 // uniform location, all properties are packed in one uniform for speed
var ShaderTileDataLoc int32
var DefaultMaterial rl.Material
var DefaultMatrix rl.Matrix
var DefaultEffects = Effects{
	BorderColor: palette.White, Tint: palette.White,
	TextColor: palette.White, TextShadowColor: palette.Black, TextShadowOffsetX: 30, TextShadowOffsetY: 30,
	TextLineHeight: 40, TextWordWrap: true, TextShadowBlur: 20, TextMarginX: 20}

var Images = make(map[int32]ImageData) // negative = crops; 0 = Font+White1x1; positive = full images
var NextImageId int16
var NextImageCropId int16

var ActiveBatch *Batch          // the batch currently being written to
var ReadyBatches []*Batch       // batches ready to be drawn
var BatchPool []*Batch          // empty batches ready to be reused
var CurrentBatchRecord []*Batch // batches being recorded, see IsRecording
var IsRecording bool            // when true, batches are accumulated into CurrentBatchRecord instead of ReadyBatches
var DrawCalls int               // used for debug info, no functional purpose - shows the amount of ReadyBatches

var ViewArea Area // zero value = entire window
var ViewX, ViewY, ViewZoom, ViewAngle float32

//=================================================================

func Queue(tex, tiles rl.Texture2D, src, dst rl.Rectangle, ang, round float32, mask Area, eff *Effects, kind, tileSz uint8, cols, rows uint16) {
	var flipU, flipV = dst.Width < 0 && kind != KindText, dst.Height < 0 && kind != KindText
	dst.Width, dst.Height = number.Absolute(dst.Width), number.Absolute(dst.Height)
	if flipU {
		dst.X -= dst.Width
	}
	if flipV {
		dst.Y -= dst.Height
	}
	var borderSz = eff.BorderSize * ViewZoom

	if ViewAngle != 0 || ViewZoom != 1 || ViewX != 0 || ViewY != 0 {
		var cx, cy = dst.X + dst.Width/2, dst.Y + dst.Height/2
		var rx, ry = cx - ViewX, cy - ViewY
		if ViewAngle != 0 {
			var sin, cos = SinCos(ViewAngle)
			cx, cy = ViewZoom*(rx*cos+ry*sin), ViewZoom*(-rx*sin+ry*cos)
		} else {
			cx, cy = ViewZoom*rx, ViewZoom*ry
		}
		dst.Width, dst.Height = dst.Width*ViewZoom, dst.Height*ViewZoom
		dst.X, dst.Y, ang = cx-dst.Width/2, cy-dst.Height/2, ang-ViewAngle
	}

	if ViewArea != (Area{}) {
		dst.X += ViewArea.X - float32(WindowWidth)/2
		dst.Y += ViewArea.Y - float32(WindowHeight)/2
	}

	var invTexW, invTexH = 1.0 / float32(tex.Width), 1.0 / float32(tex.Height)
	var u1, v1, u2, v2 = src.X * invTexW, src.Y * invTexH, (src.X + src.Width) * invTexW, (src.Y + src.Height) * invTexH
	if flipU {
		u1, u2 = u2, u1
	}
	if flipV {
		v1, v2 = v2, v1
	}
	var ww, wh = float32(WindowWidth) / 2, float32(WindowHeight) / 2
	var dx, dy = [4]float32{ww, ww, dst.Width + ww, dst.Width + ww}, [4]float32{wh, dst.Height + wh, dst.Height + wh, wh}
	var uvs = [8]float32{u1, v1, u1, v2, u2, v2, u2, v1}
	var vCount int32
	if eff == nil {
		eff = &DefaultEffects
	}

	var padU, padV float32
	if kind != KindText && borderSz > 0 { // no border padding for text symbols
		padU, padV = borderSz*(u2-u1)/dst.Width, borderSz*(v2-v1)/dst.Height
		dx[0], dx[1], dx[2], dx[3] = dx[0]-borderSz, dx[1]-borderSz, dx[2]+borderSz, dx[3]+borderSz
		dy[0], dy[3], dy[1], dy[2] = dy[0]-borderSz, dy[3]-borderSz, dy[1]+borderSz, dy[2]+borderSz
		uvs[0], uvs[2], uvs[4], uvs[6] = uvs[0]-padU, uvs[2]-padU, uvs[4]+padU, uvs[6]+padU
		uvs[1], uvs[7], uvs[3], uvs[5] = uvs[1]-padV, uvs[7]-padV, uvs[3]+padV, uvs[5]+padV
	}

	var bx, by = uint8(number.Limit(eff.BlurX, 0, 31)), uint8(number.Limit(eff.BlurY, 0, 31))
	var ps, oc = number.Limit(eff.PixelSize, 0, 16), col.Tint(eff.OutlineColor, eff.Tint)
	var r, g, b, a = col.Channels(palette.White)
	var cropMinU, cropMaxU, cropMinV, cropMaxV = u1 - 0.5, u2 - 0.5, v1 - 0.5, v2 - 0.5
	var u, v = packU2(uint16(src.Width), uint16(src.Height)), packV2(col.Tint(eff.BorderColor, eff.Tint))
	var nx = packNormalX(eff.Gamma, eff.Saturation, eff.Contrast, eff.Brightness)
	var ny, nz = packNormalY(round, ps, bx, by), packNormalZ(eff.DepthZ, borderSz, kind)
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
		} else {
			r, g, b, a = col.Channels(eff.Tint)
		}
	case KindTilemap:
		tx, ty = packTangentXTilemap(oc), packTangentYTilemap(col.Tint(eff.FillColor, eff.Tint))
		tz, tw = packTangentZTilemap(cols, rows), packTangentWTilemap(uint16(os), tileSz)
	default:
		tx, ty = packColor24(oc), packColor24(col.Tint(eff.FillColor, eff.Tint))
	}

	for i := range len(polygonBuf) {
		polygonBuf[i].U2, polygonBuf[i].V2 = u, v
		polygonBuf[i].NX, polygonBuf[i].NY, polygonBuf[i].NZ = nx, ny, nz
		polygonBuf[i].TX, polygonBuf[i].TY, polygonBuf[i].TZ, polygonBuf[i].TW = tx, ty, tz, tw
	}

	var finalColor = color.RGBA{R: r, G: g, B: b, A: a}
	var clipMask = areaIntersection(mask, ViewArea)
	if clipMask != (Area{}) {
		var sinA, cosA = SinCos(ang)
		var cx, cy = ww + dst.Width/2, wh + dst.Height/2
		for i := range 4 {
			var rx, ry = dx[i] - cx, dy[i] - cy
			polygonBuf[i].X, polygonBuf[i].Y = (rx*cosA-ry*sinA)+cx+dst.X, (rx*sinA+ry*cosA)+cy+dst.Y
			polygonBuf[i].U, polygonBuf[i].V = uvs[i*2], uvs[i*2+1]
		}
		vCount = clipPolygonAABB(polygonBuf[:4], clipResultBuf[:], clipTempBuf[:], clipMask)
		if vCount >= 3 {
			queueVertices(clipResultBuf[:vCount], tex, finalColor, tiles)
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
	queueVertices(polygonBuf[:4], tex, finalColor, tiles)
}

func ResetBatches() {

}
func CloseBatch() {
	if ActiveBatch != nil && ActiveBatch.vertCount > 0 {
		if IsRecording {
			ActiveBatch.isRecord = true
			CurrentBatchRecord = append(CurrentBatchRecord, ActiveBatch)
		} else {
			ReadyBatches = append(ReadyBatches, ActiveBatch)
		}
		ActiveBatch = nil
	}
}

var uniforms [1]float32 // reused to avoid per-frame []float32 allocation

func Draw() {
	CloseBatch()

	uniforms[0] = Runtime
	rl.SetShaderValue(Shader, ShaderLoc, uniforms[:], rl.ShaderUniformFloat)

	DrawCalls = len(ReadyBatches)
	for _, b := range ReadyBatches {
		if !b.meshUploaded {
			rl.UploadMesh(b.mesh, true)
			b.meshUploaded = true
		}
		if !b.isRecord || (b.isRecord && b.IsMeshDirty) {
			b.IsMeshDirty = false
			rl.UpdateMeshBuffer(*b.mesh, 0, b.verts[:b.vertCount*12], 0)
			rl.UpdateMeshBuffer(*b.mesh, 1, b.texCoords[:b.vertCount*8], 0)
			rl.UpdateMeshBuffer(*b.mesh, 2, b.normals[:b.vertCount*12], 0)
			rl.UpdateMeshBuffer(*b.mesh, 3, b.cols[:b.vertCount*4], 0)
			rl.UpdateMeshBuffer(*b.mesh, 4, b.tangents[:b.vertCount*16], 0)
			rl.UpdateMeshBuffer(*b.mesh, 5, b.tex2s[:b.vertCount*8], 0)
			rl.UpdateMeshBuffer(*b.mesh, 6, b.indexes[:b.indexCount*2], 0)
			b.mesh.TriangleCount = b.indexCount / 3
		}
		if b.tileDataTex.ID != 0 {
			rl.DrawRenderBatchActive()         // flush raylib's internal batch to mess texture slots
			rl.ActiveTextureSlot(1)            // switch to slot 1
			rl.EnableTexture(b.tileDataTex.ID) // bind data texture there
			rl.SetShaderValueTexture(Shader, ShaderTileDataLoc, b.tileDataTex)
		}
		rl.DrawMesh(*b.mesh, b.material, DefaultMatrix)
	}

	if ActiveBatch != nil {
		BatchPool = append(BatchPool, ActiveBatch)
		ActiveBatch = nil
	}
	for _, rb := range ReadyBatches {
		if !rb.isRecord { // text batches live on the Object, never return to pool
			BatchPool = append(BatchPool, rb)
		}
	}
	ReadyBatches = ReadyBatches[:0]
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
		Texture: DefaultMaterial.Maps.Texture, Color: DefaultMaterial.Maps.Color, Value: DefaultMaterial.Maps.Value}
	b.material.Shader = Shader
	return b
}
func queueVertices(verts []Vertex, tex rl.Texture2D, col rl.Color, tileTex rl.Texture2D) {
	if ActiveBatch != nil {
		var texChanged = ActiveBatch.material.Maps.Texture.ID != tex.ID
		var tileTexChanged = ActiveBatch.tileDataTex.ID != tileTex.ID
		var outOfSpace = ActiveBatch.vertCount+int32(len(verts)) > ActiveBatch.mesh.VertexCount

		if texChanged || tileTexChanged || outOfSpace { // do we need to break the batch?
			if ActiveBatch.vertCount > 0 {
				if IsRecording {
					ActiveBatch.isRecord = true
					CurrentBatchRecord = append(CurrentBatchRecord, ActiveBatch)
				} else {
					ReadyBatches = append(ReadyBatches, ActiveBatch) // push to draw later
				}
			}
			ActiveBatch = nil // clear active to force a new one
		}
	}

	if ActiveBatch == nil {
		if IsRecording {
			ActiveBatch = newBatch() // text batches are never pooled
			ActiveBatch.isRecord = true
		} else if len(BatchPool) > 0 { // grab a fresh batch if we don't have an active one
			ActiveBatch = BatchPool[len(BatchPool)-1]
			BatchPool = BatchPool[:len(BatchPool)-1]
		} else {
			ActiveBatch = newBatch() // pool is empty, allocate a new one (will only happen as the game ramps up)
		}

		ActiveBatch.vertCount = 0 // reset counters and set material
		ActiveBatch.indexCount = 0
		ActiveBatch.material.Maps.Texture = tex
		ActiveBatch.material.Shader = Shader
		ActiveBatch.tileDataTex = tileTex
	}

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
	var minX, maxX, minY, maxY = mask.X - mask.Width/2, mask.X + mask.Width/2, mask.Y - mask.Height/2, mask.Y + mask.Height/2
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

func areaIntersection(a, b Area) Area {
	if a == (Area{}) {
		return b
	} else if b == (Area{}) {
		return a
	}
	var al, at, ar, ab = a.X - a.Width/2, a.Y - a.Height/2, a.X + a.Width/2, a.Y + a.Height/2
	var bl, bt, br, bb = b.X - b.Width/2, b.Y - b.Height/2, b.X + b.Width/2, b.Y + b.Height/2
	var il, it, ir, ib = max(al, bl), max(at, bt), min(ar, br), min(ab, bb)
	if il >= ir || it >= ib {
		return Area{}
	}
	return Area{X: (il + ir) / 2, Y: (it + ib) / 2, Width: ir - il, Height: ib - it}
}
