package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/storage"
)

type ThemeId uint32

func LoadTheme(xmlPath string) ThemeId {
	var theme = internal.Theme{}
	storage.FromXML(file.LoadText(xmlPath), &theme)
	if theme.XMLName.Local == "" {
		return 0
	}

	internal.NextThemeId++
	var id = internal.NextThemeId

	internal.Themes[id] = theme
	return ThemeId(id)
}
