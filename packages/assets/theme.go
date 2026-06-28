package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/storage"
)

type ThemeId uint32

func LoadTheme(xmlPath string) ThemeId {
	var theme = internal.GuiTheme{}
	storage.FromXML(file.LoadText(xmlPath), &theme)
	if theme.XMLName.Local == "" {
		return 0
	}

	internal.NextThemeId++
	internal.Themes[internal.NextThemeId] = theme
	return ThemeId(internal.NextThemeId)
}

func (t ThemeId) Unload() {
	delete(internal.Themes, uint16(t))
}
