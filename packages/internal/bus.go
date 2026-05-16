package internal

type InputSnapshot struct {
	MouseX, MouseY float32
	Scroll         float32
	Input          string
	ActiveBtns     []int
	ActiveKeys     []int32
	KeyDurs        [350]float32
	WindowFocused  bool
}

type Bus struct {
	ActiveBatch  *Batch   // The batch currently being written to
	ReadyBatches []*Batch // Batches ready to be sent to the GPU
	BatchPool    []*Batch // Empty batches ready to be reused

	PendingWork   []byte              // Mainly for loading but can do other work on the Main thread too
	PendingImages map[int32]ImageData // Images loaded on the main thread, applied by the ticker

	InputSnap InputSnapshot // Input state captured on the main thread, consumed by the ticker
	Cursor    int           // Cursor state set by the ticker, consumed by the main thread

	polygonBuf, clipResultBuf, clipTempBuf [12]vertex // reused working buffers; avoids per-call heap escapes
}

var ActiveBus *Bus

var Work = make(map[byte]func())
var WorkQueue []byte
var Working = make(map[byte]bool)          // true if currently queued/running
var WorkJustFinished = make(map[byte]bool) // true for exactly one tick after finishing

// pendingImageLoads buffers images loaded on the main thread (via work).
// They are flushed onto the bus by FlushPendingImages and applied to the
// real Images map by the ticker in Reset(), avoiding map races.
var pendingImageLoads = make(map[int32]ImageData)

func AddPendingImage(id int32, img ImageData) {
	pendingImageLoads[id] = img
}

func (b *Bus) FlushPendingImages() {
	if len(pendingImageLoads) == 0 {
		return
	}
	if b.PendingImages == nil {
		b.PendingImages = make(map[int32]ImageData, len(pendingImageLoads))
	}
	for id, img := range pendingImageLoads {
		b.PendingImages[id] = img
	}
	clear(pendingImageLoads)
}

func (b *Bus) Reset() {
	if b.ActiveBatch != nil { // move all ready/active batches back to the local pool for this manager
		b.BatchPool = append(b.BatchPool, b.ActiveBatch)
		b.ActiveBatch = nil
	}
	for _, rb := range b.ReadyBatches {
		b.BatchPool = append(b.BatchPool, rb)
	}
	b.ReadyBatches = b.ReadyBatches[:0]

	// apply images loaded on the main thread to the real map (ticker-owned)
	for id, img := range b.PendingImages {
		Images[id] = img
	}
	clear(b.PendingImages)

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
	b.Cursor = Cursor // capture cursor state set by the ticker for the main thread

	if b.ActiveBatch != nil && b.ActiveBatch.vertCount > 0 {
		b.ReadyBatches = append(b.ReadyBatches, b.ActiveBatch)
		b.ActiveBatch = nil // the gameLoop finished and left the last batch open, move it to ready
	}

	if len(WorkQueue) > 0 {
		b.PendingWork = append(b.PendingWork, WorkQueue...) // copy to the bus so the ticker can safely:
		WorkQueue = WorkQueue[:0]                           // <- clear the global queue
	}
}
