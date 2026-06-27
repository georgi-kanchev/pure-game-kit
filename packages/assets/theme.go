package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/storage"
)

type ThemeId uint32

func LoadTheme(xmlPath string) ThemeId {
	var layout = internal.Layout{}
	storage.FromXML(file.LoadText(xmlPath), &layout)

	if len(layout.Boxes) == 0 {
		return 0
	}

	internal.NextLayoutId++
	var id = internal.NextLayoutId

	for i := range layout.Items { // pre-calculate item range indexes for each box
		var b = int(layout.Items[i].BoxId)
		if !layout.Boxes[b].ItemRangeCalculated {
			layout.Boxes[b].ItemStart = len(layout.Items)
			layout.Boxes[b].ItemEnd = 0
			layout.Boxes[b].ItemRangeCalculated = true
		}
		if i < layout.Boxes[b].ItemStart {
			layout.Boxes[b].ItemStart = i
		}
		if i+1 > layout.Boxes[b].ItemEnd {
			layout.Boxes[b].ItemEnd = i + 1
		}
	}

	internal.Layouts[id] = &layout
	return ThemeId(id)
}
