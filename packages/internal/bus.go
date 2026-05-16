package internal

type Bus struct {
	ActiveBatch  *Batch   // The batch currently being written to
	ReadyBatches []*Batch // Batches ready to be sent to the GPU
	BatchPool    []*Batch // Empty batches ready to be reused

	PendingWork []byte // Mainly for loading but can do other work on the Main thread too

	polygonBuf, clipResultBuf, clipTempBuf [12]vertex // reused working buffers; avoids per-call heap escapes
}

var ActiveBus *Bus

var Work = make(map[byte]func())
var WorkQueue []byte
var Working = make(map[byte]bool)          // true if currently queued/running
var WorkJustFinished = make(map[byte]bool) // true for exactly one tick after finishing

func (b *Bus) Reset() {
	if b.ActiveBatch != nil { // move all ready/active batches back to the local pool for this manager
		b.BatchPool = append(b.BatchPool, b.ActiveBatch)
		b.ActiveBatch = nil
	}
	for _, rb := range b.ReadyBatches {
		b.BatchPool = append(b.BatchPool, rb)
	}
	b.ReadyBatches = b.ReadyBatches[:0]

	//=================================================================

	clear(WorkJustFinished) // clear the finished flags from the previous tick
	if len(b.PendingWork) > 0 {
		for _, id := range b.PendingWork {
			Working[id] = false         // no longer working
			WorkJustFinished[id] = true // just finished for this gameLoop tick
		}
		b.PendingWork = b.PendingWork[:0]
	}
}
func (b *Bus) Finalize() {
	if b.ActiveBatch != nil && b.ActiveBatch.vertCount > 0 {
		b.ReadyBatches = append(b.ReadyBatches, b.ActiveBatch)
		b.ActiveBatch = nil // the gameLoop finished and left the last batch open, move it to ready
	}

	if len(WorkQueue) > 0 {
		b.PendingWork = append(b.PendingWork, WorkQueue...) // copy to the bus so the ticker can safely:
		WorkQueue = WorkQueue[:0]                           // <- clear the global queue
	}
}
