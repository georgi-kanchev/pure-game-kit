package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/storage"
)

type LayoutId uint32

func LoadLayout(xmlPath string) LayoutId {
	var layout = internal.Layout{}
	storage.FromXML(file.LoadText(xmlPath), &layout)

	if len(layout.Boxes) == 0 {
		return 0
	}

	internal.NextLayoutId++
	var id = internal.NextLayoutId
	internal.Layouts[id] = layout
	return LayoutId(id)
}

func (l LayoutId) Unload() {
	delete(internal.Layouts, uint32(l))
}
