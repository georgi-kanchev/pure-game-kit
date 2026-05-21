package internal

import (
	_ "embed"
	"pure-game-kit/packages/utility/angle"
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
	Gamma, Saturation, Contrast, Brightness float32 // Ranged 0..1

	OutlineSize, BorderSize float32

	PixelSize    byte // Ranged 0..15
	BlurX, BlurY byte // Ranged 0..31

	OutlineColor, BorderColor, SilhouetteColor uint

	// Ranged 0..1.
	//
	// Requires semi-transparent pixels to be drawn last to avoid artifacts. Fully opaque/transparent pixels work in any sorting.
	DepthZ float32
}

var DefaultMaterial rl.Material
var DefaultMatrix rl.Matrix
var Shader rl.Shader
var ShaderLoc int32 // uniform location, all properties are packed in one uniform for speed
var ShaderTileDataLoc int32

var Images = make(map[int32]ImageData) // negative = crops; 0 = White1x1; positive = full images
var NextImageId int16
var NextImageCropId int16

var ActiveBatch *Batch    // the batch currently being written to
var ReadyBatches []*Batch // batches ready to be drawn
var BatchPool []*Batch    // empty batches ready to be reused

//=================================================================

func QueueTexture(tex rl.Texture2D, src, dst rl.Rectangle, ang float32, col rl.Color, mask Area, eff *Effects) {
	if dst.Width < 0 {
		dst.Width = -dst.Width
	}
	if dst.Height < 0 {
		dst.Height = -dst.Height
	}

	var invTexW, invTexH = 1.0 / float32(tex.Width), 1.0 / float32(tex.Height)
	var u1, v1 = src.X * invTexW, src.Y * invTexH
	var u2, v2 = (src.X + src.Width) * invTexW, (src.Y + src.Height) * invTexH
	var ww, wh = float32(WindowWidth) / 2, float32(WindowHeight) / 2
	var dx = [4]float32{ww, ww, dst.Width + ww, dst.Width + ww}
	var dy = [4]float32{wh, dst.Height + wh, dst.Height + wh, wh}
	var uvs = [8]float32{u1, v1, u1, v2, u2, v2, u2, v1}
	var vCount int32
	if eff == nil {
		eff = defaultEffects
	}

	if eff.BorderSize > 0 {
		var padX = eff.BorderSize * (dst.Width / src.Width)
		var padY = eff.BorderSize * (dst.Height / src.Height)
		dx[0] -= padX
		dx[1] -= padX
		dx[2] += padX
		dx[3] += padX
		dy[0] -= padY
		dy[3] -= padY
		dy[1] += padY
		dy[2] += padY
		var padU = eff.BorderSize * invTexW
		var padV = eff.BorderSize * invTexH
		uvs[0] -= padU
		uvs[2] -= padU
		uvs[4] += padU
		uvs[6] += padU
		uvs[1] -= padV
		uvs[7] -= padV
		uvs[3] += padV
		uvs[5] += padV
	}

	for i := range len(polygonBuf) {
		polygonBuf[i].U2 = packU2(uint16(tex.Width), uint16(tex.Height))
		polygonBuf[i].V2 = packV2(eff.BorderColor)
		polygonBuf[i].NX = packNormalX(eff.Gamma, eff.Saturation, eff.Contrast, eff.Brightness)
		polygonBuf[i].NY = packNormalY(0.05, number.Limit(eff.PixelSize, 0, 16), eff.BlurX, eff.BlurY)
		polygonBuf[i].NZ = packNormalZ(eff.DepthZ, eff.BorderSize, 1)

		if true { // sprite
			polygonBuf[i].TX = packTangentXSprite(eff.OutlineColor)
			polygonBuf[i].TY = packTangentYSprite(eff.SilhouetteColor)
			polygonBuf[i].TZ = eff.OutlineSize
		}
		// else if false { // text
		// 	polygonBuf[i].TX = packTangentXText(eff.OutlineColor)
		// } else if false {
		// 	polygonBuf[i].TX = packTangentXSprite(eff.OutlineColor)
		// 	polygonBuf[i].TY = packTangentYSprite(eff.SilhouetteColor)
		// }
	}

	if mask == (Area{}) {
		vCount = 4
		if ang == 0 {
			for i := range 4 {
				polygonBuf[i].X = dx[i] + dst.X
				polygonBuf[i].Y = dy[i] + dst.Y
				polygonBuf[i].U = uvs[i*2]
				polygonBuf[i].V = uvs[i*2+1]
			}
		} else {
			var sinA, cosA = SinCos(ang)
			var cx = ww + dst.Width/2
			var cy = wh + dst.Height/2
			for i := range 4 {
				var rx = dx[i] - cx
				var ry = dy[i] - cy
				polygonBuf[i].X = (rx*cosA - ry*sinA) + cx + dst.X
				polygonBuf[i].Y = (rx*sinA + ry*cosA) + cy + dst.Y
				polygonBuf[i].U = uvs[i*2]
				polygonBuf[i].V = uvs[i*2+1]
			}
		}
		queueVertices(polygonBuf[:4], vCount, tex, col)
	} else { // CLIPPED PATH logic...
		var sinA, cosA = SinCos(ang)
		var cx = ww + dst.Width/2
		var cy = wh + dst.Height/2
		for i := range 4 {
			var rx = dx[i] - cx
			var ry = dy[i] - cy
			polygonBuf[i].X = (rx*cosA - ry*sinA) + cx + dst.X
			polygonBuf[i].Y = (rx*sinA + ry*cosA) + cy + dst.Y
			polygonBuf[i].U = uvs[i*2]
			polygonBuf[i].V = uvs[i*2+1]
		}
		vCount = clipPolygonAABB(polygonBuf[:4], clipResultBuf[:], clipTempBuf[:], mask)
		if vCount >= 3 {
			queueVertices(clipResultBuf[:vCount], vCount, tex, col)
		}
	}
}
func QueueQuad(x, y, width, height, angle float32, color rl.Color, mask Area) {
	var rect = rl.NewRectangle(x, y, width, height)
	QueueTexture(Images[0].Texture, rl.NewRectangle(0, 0, 1, 1), rect, angle, color, mask, nil)
}
func QueueLine(x1, y1, x2, y2, thickness float32, color rl.Color, mask Area) {
	if thickness <= 0 {
		return
	}

	var ang = angle.BetweenPoints(x1, y1, x2, y2)
	var dx, dy = x2 - x1, y2 - y1
	var length = number.SquareRoot(dx*dx + dy*dy)
	var perpAngle = ang - 90
	var startX, startY = moveAtAngle(x1, y1, perpAngle, thickness*0.5)
	QueueQuad(startX, startY, length, thickness, ang, color, mask)
}

//=================================================================

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
var defaultEffects = &Effects{Gamma: 0.5, Saturation: 0.5, Contrast: 0.5, Brightness: 0.5}

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
	var minX, maxX = mask.X, mask.X + mask.Width
	var minY, maxY = mask.Y, mask.Y + mask.Height
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

			out[outCount] = Vertex{
				X: prev.X + t*(curr.X-prev.X),
				Y: prev.Y + t*(curr.Y-prev.Y),
				U: prev.U + t*(curr.U-prev.U),
				V: prev.V + t*(curr.V-prev.V),
			}
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
