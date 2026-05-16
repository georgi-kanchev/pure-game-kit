package internal

import rl "github.com/gen2brain/raylib-go/raylib"

type BatchManager struct {
	ActiveBatch  *Batch   // the batch currently being written to
	ReadyBatches []*Batch // batches ready to be sent to the GPU
	BatchPool    []*Batch // empty batches ready to be reused

	polygonBuf, clipResultBuf, clipTempBuf [12]vertex // reused working buffers; avoids per-call heap escapes
}

var ActiveBatchManager *BatchManager

func (b *BatchManager) ResetBatches() {
	if b.ActiveBatch != nil {
		b.BatchPool = append(b.BatchPool, b.ActiveBatch)
		b.ActiveBatch = nil
	}
	for _, rb := range b.ReadyBatches {
		b.BatchPool = append(b.BatchPool, rb)
	}
	b.ReadyBatches = b.ReadyBatches[:0]
}

func (b *BatchManager) CloseBatch() {
	if b.ActiveBatch != nil && b.ActiveBatch.vertCount > 0 {
		b.ReadyBatches = append(b.ReadyBatches, b.ActiveBatch)
		b.ActiveBatch = nil
	}
}

func (b *BatchManager) Draw() {
	for _, batch := range b.ReadyBatches {
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
