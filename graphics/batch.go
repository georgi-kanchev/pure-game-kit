package graphics

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/angle"
	"pure-game-kit/utility/number"
	"pure-game-kit/utility/point"
	"pure-game-kit/utility/text"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Batch struct {
	mesh     *rl.Mesh
	material rl.Material

	vertsCur, indCur, quadsCapacity int32

	vertices, texCoords, colors, indices []byte
}

var batch *Batch
var prevColor rl.Color
var skipStartEnd bool

func (b *Batch) Init(quadCountCapacity int32) {
	if b.mesh != nil {
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
}
func (b *Batch) QueueQuad(tex *rl.Texture2D, src, dst rl.Rectangle, ang float32, col rl.Color) {
	if b.vertsCur != 0 && (b.material.Maps.Texture.ID != tex.ID || b.vertsCur+4 > b.mesh.VertexCount) {
		b.Draw()
	}

	if b.vertsCur == 0 {
		var mat = b.material
		rl.SetMaterialTexture(&mat, rl.MapDiffuse, *tex)
		b.material = mat
	}

	dst.Width, dst.Height = number.Absolute(dst.Width), number.Absolute(dst.Height)

	var vOffset = b.vertsCur * 12
	var vertices = unsafe.Slice((*float32)(unsafe.Pointer(&b.vertices[vOffset])), 12)
	var sinA, cosA = internal.SinCos(ang)
	var dx, dy = [4]float32{0, 0, dst.Width, dst.Width}, [4]float32{0, dst.Height, dst.Height, 0}

	for i := range 4 {
		vertices[i*3+0] = (dx[i]*cosA - dy[i]*sinA) + dst.X
		vertices[i*3+1] = (dx[i]*sinA + dy[i]*cosA) + dst.Y
		vertices[i*3+2] = 0
	}

	var tOffset = b.vertsCur * 8
	var texCoords = unsafe.Slice((*float32)(unsafe.Pointer(&b.texCoords[tOffset])), 8)
	var invTexW, invTexH = 1.0 / float32(tex.Width), 1.0 / float32(tex.Height)
	var u1, v1 = src.X * invTexW, src.Y * invTexH
	var u2, v2 = (src.X + src.Width) * invTexW, (src.Y + src.Height) * invTexH

	texCoords[0], texCoords[1] = u1, v1
	texCoords[2], texCoords[3] = u1, v2
	texCoords[4], texCoords[5] = u2, v2
	texCoords[6], texCoords[7] = u2, v1

	var cOffset = b.vertsCur * 4
	var colors = b.colors[cOffset : cOffset+16]
	for i := range 4 {
		colors[i*4+0], colors[i*4+1], colors[i*4+2], colors[i*4+3] = col.R, col.G, col.B, col.A
	}

	var iOffset = b.indCur * 2
	var indices = unsafe.Slice((*uint16)(unsafe.Pointer(&b.indices[iOffset])), 6)
	var base = uint16(b.vertsCur)
	indices[0], indices[1], indices[2] = base+0, base+1, base+2
	indices[3], indices[4], indices[5] = base+0, base+2, base+3

	b.vertsCur += 4
	b.indCur += 6
}
func (b *Batch) QueueTriangles(points []float32, col rl.Color) {
	var vertCount = int32(len(points) / 2)
	if vertCount%3 != 0 {
		return
	}

	var whiteTex = internal.White
	if b.vertsCur != 0 && (b.material.Maps.Texture.ID != whiteTex.ID || b.vertsCur+vertCount > b.mesh.VertexCount) {
		b.Draw()
	}

	if b.vertsCur == 0 {
		var mat = b.material
		rl.SetMaterialTexture(&mat, rl.MapDiffuse, *whiteTex)
		b.material = mat
	}

	var vOffset = b.vertsCur * 12
	var vertices = unsafe.Slice((*float32)(unsafe.Pointer(&b.vertices[vOffset])), vertCount*3)
	for i := range vertCount {
		vertices[i*3+0] = points[i*2+0]
		vertices[i*3+1] = points[i*2+1]
		vertices[i*3+2] = 0
	}

	var tOffset = b.vertsCur * 8
	var texCoords = unsafe.Slice((*float32)(unsafe.Pointer(&b.texCoords[tOffset])), vertCount*2)
	for i := range texCoords {
		texCoords[i] = 0
	}

	var cOffset = b.vertsCur * 4
	var colors = b.colors[cOffset : cOffset+(vertCount*4)]
	for i := range vertCount {
		colors[i*4+0], colors[i*4+1], colors[i*4+2], colors[i*4+3] = col.R, col.G, col.B, col.A
	}

	var iOffset = b.indCur * 2
	var indices = unsafe.Slice((*uint16)(unsafe.Pointer(&b.indices[iOffset])), vertCount)
	var base = uint16(b.vertsCur)
	for i := uint16(0); i < uint16(vertCount); i++ {
		indices[i] = base + i
	}

	b.vertsCur += vertCount
	b.indCur += vertCount
}
func (b *Batch) QueueTriangle(x1, y1, x2, y2, x3, y3 float32, color rl.Color) {
	b.QueueTriangles([]float32{x1, y1, x2, y2, x3, y3}, color)
}
func (b *Batch) QueueTriangleFanFloats(points []float32, color rl.Color) {
	var count = len(points) / 2
	if count < 3 {
		return
	}

	var x0, y0 = points[0], points[1]
	for i := 1; i < count-1; i++ {
		b.QueueTriangle(x0, y0, points[i*2], points[i*2+1], points[(i+1)*2], points[(i+1)*2+1], color)
	}
}
func (b *Batch) QueueLine(x1, y1, x2, y2, thickness float32, color rl.Color) {
	if thickness <= 0 {
		return
	}

	var angle = angle.BetweenPoints(x1, y1, x2, y2)
	var perpAngle = angle + 90
	var halfThickness = thickness * 0.5
	var v1x, v1y = point.MoveAtAngle(x1, y1, perpAngle, halfThickness)
	var v2x, v2y = point.MoveAtAngle(x1, y1, perpAngle, -halfThickness)
	var v3x, v3y = point.MoveAtAngle(x2, y2, perpAngle, -halfThickness)
	var v4x, v4y = point.MoveAtAngle(x2, y2, perpAngle, halfThickness)

	b.QueueTriangle(v1x, v1y, v3x, v3y, v2x, v2y, color)
	b.QueueTriangle(v1x, v1y, v4x, v4y, v3x, v3y, color)
}
func (b *Batch) QueueSymbol(font *rl.Font, s *symbol, lineHeight, gapX float32) {
	if s.UnderlineSize > 0 {
		var x, y = float32(font.Texture.Width) - 0.75, float32(font.Texture.Height) - 0.75
		var src = rl.NewRectangle(x, y, 0.2, 0.2)
		var dst = rl.NewRectangle(s.Rect.X-gapX/2, s.Y+lineHeight, s.Rect.Width+gapX, s.UnderlineSize)
		var r, g, bb, a = packSymbolColor(getColor(s.Color), getColor(255), getColor(255), 0, 0, 0, 0)
		batch.QueueQuad(&font.Texture, src, dst, s.Angle, rl.NewColor(r, g, bb, a))
	}
	if s.AssetId != "" {
		var texture, src, rotations, flip = asset(s.AssetId)
		editAssetRects(&src, &s.Rect, s.Angle, rotations, flip)
		var r, g, bb, a = packSymbolColor(getColor(s.Color), getColor(255), getColor(255), 1, 3, 3, 0)
		batch.QueueQuad(texture, src, s.Rect, s.Angle, rl.NewColor(r, g, bb, a))
		return
	}
	if text.Trim(s.Value) != "" {
		var r, g, bb, a = packSymbolColor(getColor(s.Color), getColor(255), getColor(255), 1, 2, 3, 0)
		b.QueueQuad(&font.Texture, s.TexRect, s.Rect, s.Angle, rl.NewColor(r, g, bb, a))
	}
}

func (b *Batch) Draw() {
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
