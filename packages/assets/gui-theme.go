package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/storage"
)

type GUIThemeId uint16

func LoadGUITheme(xmlPath string) GUIThemeId {
	var theme = &internal.GUIThemeVariants{}
	storage.FromXML(file.LoadText(xmlPath), &theme)

	var id uint8
	var variantsCount = 0
	var variants = []*internal.GUITheme{
		theme.Variant0, theme.Variant1, theme.Variant2, theme.Variant3, theme.Variant4,
		theme.Variant5, theme.Variant6, theme.Variant7, theme.Variant8, theme.Variant9,
	}
	for i := range 10 {
		if variants[i] != nil {
			variantsCount++
			internal.NextGUIThemeId++
			internal.GUIThemes[internal.NextGUIThemeId] = variants[i]

			if id == 0 {
				id = internal.NextGUIThemeId
			}
		}
	}
	if id == 0 {
		return 0
	}

	internal.GUIThemes[id].VariantsCount = variantsCount
	return GUIThemeId(id)
}

func (t GUIThemeId) Unload() {
	delete(internal.GUIThemes, uint8(t))
}

func (t GUIThemeId) Variant(index int) GUIThemeId {
	var th, has = internal.GUIThemes[uint8(t)]
	if !has {
		return 0
	}
	var count = th.VariantsCount - 1
	if count+index < int(t) || index > count {
		return t
	}
	return GUIThemeId(int(t) + index)
}
