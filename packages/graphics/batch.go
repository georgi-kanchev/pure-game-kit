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
	if b.mesh != nil {
		b.mesh.Vertices = nil
		b.mesh.Texcoords = nil
		b.mesh.Colors = nil
		b.mesh.Indices = nil
		rl.UnloadMesh(b.mesh)
	}

	b.vertCount = 0
	b.indexCount = 0
	b.maxQuads = quadCountCapacity

	b.mesh = &rl.Mesh{VertexCount: 4 * quadCountCapacity, TriangleCount: 2 * quadCountCapacity}

	b.vertData = make([]byte, b.mesh.VertexCount*3*4)      // Vec3 (x,y,z) * float32
	b.texCoordsData = make([]byte, b.mesh.VertexCount*2*4) // Vec2 (u,v) * float32
	b.colData = make([]byte, b.mesh.VertexCount*4)         // Color (r,g,b,a) * byte
	b.indexData = make([]byte, b.mesh.TriangleCount*3*2)   // Triangle (i1,i2,i3) * uint16

	b.mesh.Vertices = (*float32)(unsafe.Pointer(&b.vertData[0]))
	b.mesh.Texcoords = (*float32)(unsafe.Pointer(&b.texCoordsData[0]))
	b.mesh.Colors = (*uint8)(unsafe.Pointer(&b.colData[0]))
	b.mesh.Indices = (*uint16)(unsafe.Pointer(&b.indexData[0]))

	rl.UploadMesh(b.mesh, true)
	b.material = rl.LoadMaterialDefault()
	b.material.Shader = internal.Shader
}
func (b *batch) Draw() {
	if b.vertCount == 0 {
		return
	}

	rl.UpdateMeshBuffer(*b.mesh, 0, b.vertData, 0)
	rl.UpdateMeshBuffer(*b.mesh, 1, b.texCoordsData, 0)
	rl.UpdateMeshBuffer(*b.mesh, 3, b.colData, 0)
	rl.UpdateMeshBuffer(*b.mesh, 6, b.indexData, 0)

	b.mesh.TriangleCount = b.indexCount / 3
	rl.DrawMesh(*b.mesh, b.material, internal.MatrixDefault)

	if b.vertCount >= b.mesh.VertexCount {
		b.Init(b.maxQuads * 2)
	}

	b.vertCount = 0
	b.indexCount = 0
}

// private =================================================================

type vertex struct{ X, Y, U, V float32 }
type batch struct {
	mesh     *rl.Mesh
	material rl.Material

	clipMask Area

	vertCount, indexCount, maxQuads int32

	vertData      []byte // Vec3 (x,y,z) * float32
	texCoordsData []byte // Vec2 (u,v) * float32
	colData       []byte // Color (r,g,b,a) * byte
	indexData     []byte // Triangle (i1,i2,i3) * uint16

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

	var vertCount = int32(len(verts))
	var vertSlice = unsafe.Slice((*float32)(unsafe.Pointer(&b.vertData[b.vertCount*12])), vertCount*3)    // 3 floats (12 bytes)
	var texSlice = unsafe.Slice((*float32)(unsafe.Pointer(&b.texCoordsData[b.vertCount*8])), vertCount*2) // 2 floats (8 bytes)
	var colSlice = b.colData[(b.vertCount * 4) : (b.vertCount*4)+(vertCount*4)]                           // 4 bytes (rgba)

	for i, v := range verts {
		vertSlice[i*3+0], vertSlice[i*3+1], vertSlice[i*3+2] = v.X, v.Y, 0
		texSlice[i*2+0], texSlice[i*2+1] = v.U, v.V
		colSlice[i*4+0], colSlice[i*4+1], colSlice[i*4+2], colSlice[i*4+3] = col.R, col.G, col.B, col.A
	}

	var trisCount = vertCount - 2
	var indSlice = unsafe.Slice((*uint16)(unsafe.Pointer(&b.indexData[b.indexCount*2])), trisCount*3)
	var base = uint16(b.vertCount)

	for i := range trisCount {
		indSlice[i*3+0] = base
		indSlice[i*3+1] = base + uint16(i+1)
		indSlice[i*3+2] = base + uint16(i+2)
	}

	b.vertCount += vertCount
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
