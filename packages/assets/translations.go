package assets

import (
	"pure-game-kit/packages/internal"
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/storage"
)

type LanguageId uint32

func LoadTranslations(yamlPath string) LanguageId {
	var lang = internal.Lang{Tags: make(map[string]string)}
	storage.FromYAML(file.LoadText(yamlPath), &lang.Tags)
	if len(lang.Tags) == 0 {
		return 0
	}

	var id = LanguageId(len(internal.Translations) + 1)
	internal.Translations[uint32(id)] = lang
	return id
}

func (l LanguageId) Translate(tag string) string {
	var value, has = internal.Translations[uint32(l)]
	if !has {
		return ""
	}
	return value.Tags[tag]
}
func (l LanguageId) Unload() {
	delete(internal.Translations, uint32(l))
}
