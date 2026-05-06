package graphics

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/angle"
	"pure-game-kit/packages/utility/number"
	"pure-game-kit/packages/utility/point"
	"unicode"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (b *batch) QueueTexture(tex rl.Texture2D, src, dst rl.Rectangle, ang float32, col rl.Color) {
	if dst.Width < 0 {
		dst.Width = -dst.Width
	}
	if dst.Height < 0 {
		dst.Height = -dst.Height
	}

	var invTexW, invTexH = 1.0 / float32(tex.Width), 1.0 / float32(tex.Height)
	var u1, v1 = src.X * invTexW, src.Y * invTexH
	var u2, v2 = (src.X + src.Width) * invTexW, (src.Y + src.Height) * invTexH
	var dx = [4]float32{0, 0, dst.Width, dst.Width}
	var dy = [4]float32{0, dst.Height, dst.Height, 0}
	var uvs = [8]float32{u1, v1, u1, v2, u2, v2, u2, v1}
	var vCount int32

	if b.clipMask == (Area{}) {
		vCount = 4
		if ang == 0 {
			for i := range 4 {
				b.polygonBuf[i].X = dx[i] + dst.X
				b.polygonBuf[i].Y = dy[i] + dst.Y
				b.polygonBuf[i].U = uvs[i*2]
				b.polygonBuf[i].V = uvs[i*2+1]
			}
		} else {
			var sinA, cosA = internal.SinCos(ang)
			for i := range 4 {
				b.polygonBuf[i].X = (dx[i]*cosA - dy[i]*sinA) + dst.X
				b.polygonBuf[i].Y = (dx[i]*sinA + dy[i]*cosA) + dst.Y
				b.polygonBuf[i].U = uvs[i*2]
				b.polygonBuf[i].V = uvs[i*2+1]
			}
		}
		b.writeToBuffers(b.polygonBuf[:4], vCount, tex, col) // pass the buffer directly
	} else { // CLIPPED PATH logic...
		var sinA, cosA = internal.SinCos(ang)
		for i := range 4 {
			b.polygonBuf[i].X = (dx[i]*cosA - dy[i]*sinA) + dst.X
			b.polygonBuf[i].Y = (dx[i]*sinA + dy[i]*cosA) + dst.Y
			b.polygonBuf[i].U = uvs[i*2]
			b.polygonBuf[i].V = uvs[i*2+1]
		}
		vCount = clipPolygonAABB(b.polygonBuf[:4], b.clipResultBuf[:], b.clipTempBuf[:], b.clipMask)
		if vCount >= 3 {
			b.writeToBuffers(b.clipResultBuf[:vCount], vCount, tex, col)
		}
	}
}

func (b *batch) QueueQuad(x, y, width, height, angle float32, color rl.Color) {
	b.QueueTexture(internal.White1x1, rl.NewRectangle(0, 0, 1, 1), rl.NewRectangle(x, y, width, height), angle, color)
}
func (b *batch) QueueLine(x1, y1, x2, y2, thickness float32, color rl.Color) {
	if thickness <= 0 {
		return
	}

	var ang = angle.BetweenPoints(x1, y1, x2, y2)
	var dx, dy = x2 - x1, y2 - y1
	var length = number.SquareRoot(dx*dx + dy*dy)
	var perpAngle = ang - 90
	var startX, startY = point.MoveAtAngle(x1, y1, perpAngle, thickness*0.5)
	b.QueueQuad(startX, startY, length, thickness, ang, color)
}
func (b *batch) QueueSymbol(font rl.Font, s symbol, lineHeight, gapX float32) {
	var lineThickness = lineHeight / 15
	var tx, ty = float32(font.Texture.Width) - 0.75, float32(font.Texture.Height) - 0.75
	var fontSrc = rl.NewRectangle(tx, ty, 0.2, 0.2)

	if s.BackColor > 0 {
		var prevCol = s.Color
		var rect = rl.NewRectangle(s.Bounds.X, s.Bounds.Y, s.Bounds.Width+gapX, s.Bounds.Height)
		s.Color = s.BackColor
		batcher.QueueTexture(font.Texture, fontSrc, rect, s.Angle, packSymbolColor(s))
		s.Color = prevCol
	}

	if s.Underline {
		var offset = (lineHeight - lineThickness) - s.TopCrop
		if offset >= 0 && offset+lineThickness <= s.Bounds.Height {
			var x, y = point.MoveAtAngle(s.Bounds.X, s.Bounds.Y, s.Angle+90, offset)
			batcher.QueueTexture(font.Texture, fontSrc, rl.NewRectangle(x, y, s.Bounds.Width+gapX, lineThickness), s.Angle, packSymbolColor(s))
		}
	}

	if !unicode.IsSpace(s.Value) {
		b.QueueTexture(s.Texture, s.TexRect, s.Rect, s.Angle, packSymbolColor(s))
	}

	if s.Strikethrough {
		var offset = (lineHeight*0.55 - lineThickness/2) - s.TopCrop
		if offset >= 0 && offset+lineThickness <= s.Bounds.Height {
			var x, y = point.MoveAtAngle(s.Bounds.X, s.Bounds.Y, s.Angle+90, offset)
			batcher.QueueTexture(font.Texture, fontSrc, rl.NewRectangle(x, y, s.Bounds.Width+gapX, lineThickness), s.Angle, packSymbolColor(s))
		}
	}
}

func (b *batch) Init(quadCountCapacity int32) {
	if b.mesh != nil { // nilling these is needed, causing problems on windows (on linux it's fine)
		b.mesh.Vertices = nil
		b.mesh.Texcoords = nil
		b.mesh.Normals = nil
		b.mesh.Colors = nil
		b.mesh.Tangents = nil
		b.mesh.Texcoords2 = nil
		b.mesh.Indices = nil
		rl.UnloadMesh(b.mesh)
	}

	b.vertCount = 0
	b.indexCount = 0
	b.maxQuads = quadCountCapacity

	b.mesh = &rl.Mesh{VertexCount: 4 * quadCountCapacity, TriangleCount: 2 * quadCountCapacity}

	b.verts = make([]byte, b.mesh.VertexCount*3*4)     // vec3 (12 bytes)
	b.texCoords = make([]byte, b.mesh.VertexCount*2*4) // vec2 (8 bytes)
	b.normals = make([]byte, b.mesh.VertexCount*3*4)   // vec3 (12 bytes)
	b.cols = make([]byte, b.mesh.VertexCount*4)        // rgba (4 bytes)
	b.tangents = make([]byte, b.mesh.VertexCount*4*4)  // vec4 (16 bytes)
	b.tex2s = make([]byte, b.mesh.VertexCount*2*4)     // vec2 (8 bytes)
	b.indexes = make([]byte, b.mesh.TriangleCount*3*2) // uint16 (6 bytes per triangle)

	b.mesh.Vertices = (*float32)(unsafe.Pointer(&b.verts[0]))
	b.mesh.Texcoords = (*float32)(unsafe.Pointer(&b.texCoords[0]))
	b.mesh.Normals = (*float32)(unsafe.Pointer(&b.normals[0]))
	b.mesh.Colors = (*uint8)(unsafe.Pointer(&b.cols[0]))
	b.mesh.Tangents = (*float32)(unsafe.Pointer(&b.tangents[0]))
	b.mesh.Texcoords2 = (*float32)(unsafe.Pointer(&b.tex2s[0]))
	b.mesh.Indices = (*uint16)(unsafe.Pointer(&b.indexes[0]))

	rl.UploadMesh(b.mesh, true)
	b.material = rl.LoadMaterialDefault()
	b.material.Shader = internal.Shader
}
func (b *batch) Draw() {
	if b.vertCount == 0 {
		return
	}

	rl.UpdateMeshBuffer(*b.mesh, 0, b.verts, 0)
	rl.UpdateMeshBuffer(*b.mesh, 1, b.texCoords, 0)
	rl.UpdateMeshBuffer(*b.mesh, 2, b.normals, 0)
	rl.UpdateMeshBuffer(*b.mesh, 3, b.cols, 0)
	rl.UpdateMeshBuffer(*b.mesh, 4, b.tangents, 0)
	rl.UpdateMeshBuffer(*b.mesh, 5, b.tex2s, 0)
	rl.UpdateMeshBuffer(*b.mesh, 6, b.indexes, 0)

	b.mesh.TriangleCount = b.indexCount / 3
	rl.DrawMesh(*b.mesh, b.material, internal.MatrixDefault)

	if b.vertCount >= b.mesh.VertexCount {
		b.Init(b.maxQuads * 2)
	}

	b.vertCount = 0
	b.indexCount = 0
}

// private =================================================================

type vertex struct {
	X, Y, U, V float32
	NX, NY, NZ float32 // Normals
	TX, TY, TZ float32 // Tangents
	U2, V2     float32 // Texcoords2
}
type batch struct {
	mesh     *rl.Mesh
	material rl.Material

	clipMask Area

	vertCount, indexCount, maxQuads int32

	verts     []byte // 0: vec3
	texCoords []byte // 1: vec2
	normals   []byte // 2: vec3
	cols      []byte // 3: rgba byte
	tangents  []byte // 4: vec3
	tex2s     []byte // 5: vec2
	indexes   []byte // 6: uint16

	polygonBuf, clipResultBuf, clipTempBuf [12]vertex // reused working buffers; avoids per-call heap escapes
}

var batcher *batch

func (b *batch) writeToBuffers(verts []vertex, vCount int32, tex rl.Texture2D, col rl.Color) {
	if b.vertCount != 0 && (b.material.Maps.Texture.ID != tex.ID || b.vertCount+vCount > b.mesh.VertexCount) {
		b.Draw()
	}

	if b.vertCount == 0 {
		b.material.Maps.Texture = tex
		b.material.Shader = internal.Shader
	}

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
		tan_slice[i*4+0], tan_slice[i*4+1], tan_slice[i*4+2], tan_slice[i*4+3] = 1, 0, 0, 1
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

func clipPolygonAABB(poly, outBuf, tempBuf []vertex, mask Area) int32 {
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
func clipPolyEdge(in, out []vertex, isX bool, edgeVal float32, keepGreater bool) int32 {
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

			out[outCount] = vertex{
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
