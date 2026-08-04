package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/storage"
)

type GUIThemeId uint16

func LoadGUITheme(xmlPath string) GUIThemeId {
	var theme = internal.GUITheme{}
	storage.FromXML(file.LoadText(xmlPath), &theme)
	if theme.XMLName.Local == "" {
		return 0
	}

	internal.NextGUIThemeId++
	internal.GUIThemes[internal.NextGUIThemeId] = theme
	return GUIThemeId(internal.NextGUIThemeId)
}

func (t GUIThemeId) Unload() {
	delete(internal.GUIThemes, uint16(t))
}
