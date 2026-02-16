package graphics

import (
	"math"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Batch struct {
	mesh *rl.Mesh

	quadCountCurrent, quadCountCapacity   int32
	vertices, texCoords, colors, indices  []byte
	dirtyVerts, dirtyTexCoords, dirtyCols bool
}

// var batches = make(map[*rl.Texture2D]*Batch)

func (b *Batch) Init(quadCountCapacity int32) {
	if b.mesh != nil {
		rl.UnloadMesh(b.mesh)
	}

	b.mesh = &rl.Mesh{VertexCount: 4 * quadCountCapacity, TriangleCount: 2 * quadCountCapacity}
	b.quadCountCurrent, b.quadCountCapacity = 0, quadCountCapacity

	b.vertices = make([]byte, b.mesh.VertexCount*3*4)  // 4 verts * 3 floats (xyz) * 4 bytes per float
	b.texCoords = make([]byte, b.mesh.VertexCount*2*4) // 4 verts * 2 floats (uv) * 4 bytes per float
	b.colors = make([]byte, b.mesh.VertexCount*4)      // 4 verts * 4 bytes (rgba)
	b.indices = make([]byte, b.mesh.TriangleCount*3*2) // 6 indices * 2 bytes (uint16)

	var rawIndices = unsafe.Slice((*uint16)(unsafe.Pointer(&b.indices[0])), len(b.indices)/2)
	for i := range quadCountCapacity {
		var baseVert, baseInd = uint16(i * 4), i * 6
		rawIndices[baseInd+0] = baseVert + 0
		rawIndices[baseInd+1] = baseVert + 1
		rawIndices[baseInd+2] = baseVert + 2
		rawIndices[baseInd+3] = baseVert + 0
		rawIndices[baseInd+4] = baseVert + 2
		rawIndices[baseInd+5] = baseVert + 3
	}

	b.mesh.Vertices = (*float32)(unsafe.Pointer(&b.vertices[0]))
	b.mesh.Texcoords = (*float32)(unsafe.Pointer(&b.texCoords[0]))
	b.mesh.Colors = (*uint8)(unsafe.Pointer(&b.colors[0]))
	b.mesh.Indices = (*uint16)(unsafe.Pointer(&b.indices[0]))

	rl.UploadMesh(b.mesh, true)
}

func (b *Batch) Queue(tex rl.Texture2D, src, dst rl.Rectangle, origin rl.Vector2, rotation float32, color rl.Color) {
	var id = b.quadCountCurrent
	var vertices = unsafe.Slice((*float32)(unsafe.Pointer(&b.vertices[id*48])), 12)
	var rad = rotation * rl.Deg2rad
	var cosA, sinA = float32(math.Cos(float64(rad))), float32(math.Sin(float64(rad)))
	var dx, dy = [4]float32{0, 0, dst.Width, dst.Width}, [4]float32{0, dst.Height, dst.Height, 0}

	for i := range 4 {
		rx := dx[i] - origin.X
		ry := dy[i] - origin.Y

		vertices[i*3+0] = (rx*cosA - ry*sinA) + dst.X
		vertices[i*3+1] = (rx*sinA + ry*cosA) + dst.Y
		vertices[i*3+2] = 0 // z always 0
	}

	var texCoords = unsafe.Slice((*float32)(unsafe.Pointer(&b.texCoords[id*32])), 8)
	var invTexW, invTexH = 1.0 / float32(tex.Width), 1.0 / float32(tex.Height)
	var u1, v1 = src.X * invTexW, src.Y * invTexH
	var u2, v2 = (src.X + src.Width) * invTexW, (src.Y + src.Height) * invTexH

	texCoords[0], texCoords[1] = u1, v1 // tl
	texCoords[2], texCoords[3] = u1, v2 // bl
	texCoords[4], texCoords[5] = u2, v2 // br
	texCoords[6], texCoords[7] = u2, v1 // tr

	var cOffset = id * 16
	var colors = b.colors[cOffset : cOffset+16]

	for i := range 4 {
		colors[i*4+0] = color.R
		colors[i*4+1] = color.G
		colors[i*4+2] = color.B
		colors[i*4+3] = color.A
	}

	b.dirtyVerts = true
	b.dirtyTexCoords = true
	b.dirtyCols = true
	b.quadCountCurrent++
}

func (b *Batch) Draw(material rl.Material) {
	if b.quadCountCurrent == 0 {
		return
	}

	if b.dirtyVerts {
		rl.UpdateMeshBuffer(*b.mesh, 0, b.vertices, 0)
	}
	if b.dirtyTexCoords {
		rl.UpdateMeshBuffer(*b.mesh, 1, b.texCoords, 0)
	}
	if b.dirtyCols {
		rl.UpdateMeshBuffer(*b.mesh, 3, b.colors, 0)
	}

	b.mesh.TriangleCount = b.quadCountCurrent * 2
	rl.DrawMesh(*b.mesh, material, rl.MatrixIdentity())

	if b.quadCountCurrent >= b.quadCountCapacity {
		b.Init(b.quadCountCapacity * 2)
	}

	b.quadCountCurrent = 0
	b.dirtyVerts = false
	b.dirtyTexCoords = false
	b.dirtyCols = false
}
