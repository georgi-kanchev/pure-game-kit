package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/angle"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"unicode"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type batchData struct {
	mesh     *rl.Mesh
	material rl.Material

	skipStartEnd, skipDraw bool

	mask Area

	vertsCur, indCur, quadsCapacity int32

	vertices, texCoords, colors, indices []byte

	polyBuf, clipBuf, clipTempBuf [12]batchVertex // reused working buffers; avoids per-call heap escapes
}

var batch *batchData
var prevColor rl.Color

func (b *batchData) Init(quadCountCapacity int32) {
	if b.mesh != nil {
		b.mesh.Vertices = nil
		b.mesh.Texcoords = nil
		b.mesh.Colors = nil
		b.mesh.Indices = nil
		rl.UnloadMesh(b.mesh)
	}

	b.vertsCur = 0
	b.indCur = 0
	b.quadsCapacity = quadCountCapacity

	b.mesh = &rl.Mesh{VertexCount: 4 * quadCountCapacity, TriangleCount: 2 * quadCountCapacity}

	b.vertices = make([]byte, b.mesh.VertexCount*3*4)
	b.texCoords = make([]byte, b.mesh.VertexCount*2*4)
	b.colors = make([]byte, b.mesh.VertexCount*4)
	b.indices = make([]byte, b.mesh.TriangleCount*3*2)

	b.mesh.Vertices = (*float32)(unsafe.Pointer(&b.vertices[0]))
	b.mesh.Texcoords = (*float32)(unsafe.Pointer(&b.texCoords[0]))
	b.mesh.Colors = (*uint8)(unsafe.Pointer(&b.colors[0]))
	b.mesh.Indices = (*uint16)(unsafe.Pointer(&b.indices[0]))

	rl.UploadMesh(b.mesh, true)
	b.material = rl.LoadMaterialDefault()
	b.material.Shader = internal.Shader
}
func (b *batchData) QueueTex(tex rl.Texture2D, src, dst rl.Rectangle, ang float32, col rl.Color) {
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
	var vCount int

	if b.mask == (Area{}) {
		vCount = 4
		if ang == 0 {
			for i := range 4 {
				// Direct assignment to pre-allocated buffer on the struct
				b.polyBuf[i].X = dx[i] + dst.X
				b.polyBuf[i].Y = dy[i] + dst.Y
				b.polyBuf[i].U = uvs[i*2]
				b.polyBuf[i].V = uvs[i*2+1]
			}
		} else {
			var sinA, cosA = internal.SinCos(ang)
			for i := range 4 {
				b.polyBuf[i].X = (dx[i]*cosA - dy[i]*sinA) + dst.X
				b.polyBuf[i].Y = (dx[i]*sinA + dy[i]*cosA) + dst.Y
				b.polyBuf[i].U = uvs[i*2]
				b.polyBuf[i].V = uvs[i*2+1]
			}
		}
		// Instead of verts = b.polyBuf[:vCount], pass the buffer directly
		b.prepareAndWrite(b.polyBuf[:4], vCount, tex, col)
	} else {
		// CLIPPED PATH logic...
		var sinA, cosA = internal.SinCos(ang)
		for i := 0; i < 4; i++ {
			b.polyBuf[i].X = (dx[i]*cosA - dy[i]*sinA) + dst.X
			b.polyBuf[i].Y = (dx[i]*sinA + dy[i]*cosA) + dst.Y
			b.polyBuf[i].U = uvs[i*2]
			b.polyBuf[i].V = uvs[i*2+1]
		}
		vCount = clipPolygonAABB(b.polyBuf[:4], b.clipBuf[:], b.clipTempBuf[:], b.mask)
		if vCount >= 3 {
			b.prepareAndWrite(b.clipBuf[:vCount], vCount, tex, col)
		}
	}
}

func (b *batchData) QueueQuad(x, y, width, height, angle float32, color rl.Color) {
	b.QueueTex(internal.White, rl.NewRectangle(0, 0, 1, 1), rl.NewRectangle(x, y, width, height), angle, color)
}
func (b *batchData) QueueTriangles(points []float32, col rl.Color) {
	var totalTriangles = int(len(points) / 6)
	if totalTriangles == 0 {
		return
	}
	var white = internal.White

	for i := range totalTriangles {
		var pOffset = i * 6
		b.polyBuf[0] = batchVertex{X: points[pOffset+0], Y: points[pOffset+1]}
		b.polyBuf[1] = batchVertex{X: points[pOffset+2], Y: points[pOffset+3]}
		b.polyBuf[2] = batchVertex{X: points[pOffset+4], Y: points[pOffset+5]}

		var verts []batchVertex
		var vCount int

		if b.mask == (Area{}) {
			vCount = 3
			verts = b.polyBuf[:3]
		} else {
			vCount = clipPolygonAABB(b.polyBuf[:3], b.clipBuf[:], b.clipTempBuf[:], b.mask)
			if vCount < 3 {
				continue
			}
			verts = b.clipBuf[:vCount]
		}

		var vCount32 = int32(vCount)
		if b.vertsCur != 0 && (b.material.Maps.Texture.ID != white.ID || b.vertsCur+vCount32 > b.mesh.VertexCount) {
			b.Draw()
		}
		if b.vertsCur == 0 {
			var mat = b.material
			rl.SetMaterialTexture(&mat, rl.MapDiffuse, white)
			b.material = mat
			b.material.Shader = internal.Shader
		}

		b.writeToBuffers(verts, col)
	}
}
func (b *batchData) QueueTriangle(x1, y1, x2, y2, x3, y3 float32, color rl.Color) {
	b.polyBuf[0] = batchVertex{X: x1, Y: y1}
	b.polyBuf[1] = batchVertex{X: x2, Y: y2}
	b.polyBuf[2] = batchVertex{X: x3, Y: y3}

	var verts []batchVertex
	var vCount int

	if b.mask == (Area{}) {
		vCount = 3
		verts = b.polyBuf[:3]
	} else {
		vCount = clipPolygonAABB(b.polyBuf[:3], b.clipBuf[:], b.clipTempBuf[:], b.mask)
		if vCount < 3 {
			return
		}
		verts = b.clipBuf[:vCount]
	}

	var white = internal.White
	var vCount32 = int32(vCount)
	if b.vertsCur != 0 && (b.material.Maps.Texture.ID != white.ID || b.vertsCur+vCount32 > b.mesh.VertexCount) {
		b.Draw()
	}
	if b.vertsCur == 0 {
		var mat = b.material
		rl.SetMaterialTexture(&mat, rl.MapDiffuse, white)
		b.material = mat
		b.material.Shader = internal.Shader
	}

	b.writeToBuffers(verts, color)
}
func (b *batchData) QueueTriangleFanFloats(points []float32, color rl.Color) {
	var count = len(points) / 2
	if count < 3 {
		return
	}

	var x0, y0 = points[0], points[1]
	for i := 1; i < count-1; i++ {
		b.QueueTriangle(x0, y0, points[i*2], points[i*2+1], points[(i+1)*2], points[(i+1)*2+1], color)
	}
}
func (b *batchData) QueueLine(x1, y1, x2, y2, thickness float32, color rl.Color) {
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
func (b *batchData) QueueSymbol(font rl.Font, s symbol, lineHeight, gapX float32) {
	var lineThickness = lineHeight / 15
	var tx, ty = float32(font.Texture.Width) - 0.75, float32(font.Texture.Height) - 0.75
	var fontSrc = rl.NewRectangle(tx, ty, 0.2, 0.2)

	if s.BackColor > 0 {
		var prevCol = s.Color
		s.Color = s.BackColor
		batch.QueueTex(font.Texture, fontSrc, rl.NewRectangle(s.Bounds.X, s.Bounds.Y, s.Bounds.Width+gapX, s.Bounds.Height), s.Angle, packSymbolColor(s))
		s.Color = prevCol
	}

	if s.Underline {
		var offset = (lineHeight - lineThickness) - s.TopCrop
		if offset >= 0 && offset+lineThickness <= s.Bounds.Height {
			var x, y = point.MoveAtAngle(s.Bounds.X, s.Bounds.Y, s.Angle+90, offset)
			batch.QueueTex(font.Texture, fontSrc, rl.NewRectangle(x, y, s.Bounds.Width+gapX, lineThickness), s.Angle, packSymbolColor(s))
		}
	}

	if !unicode.IsSpace(s.Value) {
		b.QueueTex(s.Texture, s.TexRect, s.Rect, s.Angle, packSymbolColor(s))
	}

	if s.Strikethrough {
		var offset = (lineHeight*0.55 - lineThickness/2) - s.TopCrop
		if offset >= 0 && offset+lineThickness <= s.Bounds.Height {
			var x, y = point.MoveAtAngle(s.Bounds.X, s.Bounds.Y, s.Angle+90, offset)
			batch.QueueTex(font.Texture, fontSrc, rl.NewRectangle(x, y, s.Bounds.Width+gapX, lineThickness), s.Angle, packSymbolColor(s))
		}
	}
}

func (b *batchData) Draw() {
	if b.vertsCur == 0 {
		return
	}

	rl.UpdateMeshBuffer(*b.mesh, 0, b.vertices, 0)
	rl.UpdateMeshBuffer(*b.mesh, 1, b.texCoords, 0)
	rl.UpdateMeshBuffer(*b.mesh, 3, b.colors, 0)
	rl.UpdateMeshBuffer(*b.mesh, 6, b.indices, 0)

	b.mesh.TriangleCount = b.indCur / 3
	rl.DrawMesh(*b.mesh, b.material, internal.MatrixDefault)

	if b.vertsCur >= b.mesh.VertexCount {
		b.Init(b.quadsCapacity * 2)
	}

	b.vertsCur = 0
	b.indCur = 0
}

// private =================================================================

type batchVertex struct {
	X, Y, U, V float32
}

func (b *batchData) prepareAndWrite(verts []batchVertex, vCount int, tex rl.Texture2D, col rl.Color) {
	vCount32 := int32(vCount)

	// Check if we need to flush
	if b.vertsCur != 0 && (b.material.Maps.Texture.ID != tex.ID || b.vertsCur+vCount32 > b.mesh.VertexCount) {
		b.Draw()
	}

	if b.vertsCur == 0 {
		// Direct assignment to the texture map to avoid SetMaterialTexture's overhead
		// Note: MapDiffuse index is 0
		b.material.Maps.Texture = tex
		b.material.Shader = internal.Shader
	}

	b.writeToBuffers(verts, col)
}
func (b *batchData) writeToBuffers(vertices []batchVertex, col rl.Color) {
	var vCount = int32(len(vertices))
	var vOffset = b.vertsCur * 12
	var vSlice = unsafe.Slice((*float32)(unsafe.Pointer(&b.vertices[vOffset])), vCount*3)
	var tOffset = b.vertsCur * 8
	var tSlice = unsafe.Slice((*float32)(unsafe.Pointer(&b.texCoords[tOffset])), vCount*2)
	var cOffset = b.vertsCur * 4
	var cSlice = b.colors[cOffset : cOffset+(vCount*4)]

	for i, v := range vertices {
		vSlice[i*3+0], vSlice[i*3+1], vSlice[i*3+2] = v.X, v.Y, 0
		tSlice[i*2+0], tSlice[i*2+1] = v.U, v.V
		cSlice[i*4+0], cSlice[i*4+1], cSlice[i*4+2], cSlice[i*4+3] = col.R, col.G, col.B, col.A
	}

	var trisCount = vCount - 2
	var iOffset = b.indCur * 2
	var iSlice = unsafe.Slice((*uint16)(unsafe.Pointer(&b.indices[iOffset])), trisCount*3)
	var base = uint16(b.vertsCur)

	for i := int32(0); i < trisCount; i++ {
		iSlice[i*3+0] = base
		iSlice[i*3+1] = base + uint16(i+1)
		iSlice[i*3+2] = base + uint16(i+2)
	}

	b.vertsCur += vCount
	b.indCur += trisCount * 3
}

func clipPolygonAABB(poly, outBuf, tempBuf []batchVertex, mask Area) int {
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
func clipPolyEdge(in, out []batchVertex, isX bool, edgeVal float32, keepGreater bool) int {
	var outCount = 0
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

			out[outCount] = batchVertex{
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
