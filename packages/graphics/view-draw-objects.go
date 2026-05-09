package graphics

func (v *View) DrawObjects(objects ...*Object) {
	v.begin()
	for _, t := range objects {
		if t == nil || !v.IsAreaVisible(t.Bounds()) {
			continue
		}

		t.tryRegenerateText()
		// for _, s := range t.chars {

		// }
	}
	batcher.Draw()
	v.end()
}
