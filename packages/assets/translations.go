package assets

import (
	"pure-game-kit/packages/utility/file"
	"pure-game-kit/packages/utility/storage"
)

type LanguageId uint32

func LoadTranslations(yamlPath string) LanguageId {
	var lang = lang{tags: make(map[string]string)}
	storage.FromYAML(file.LoadText(yamlPath), &lang.tags)

	if len(lang.tags) == 0 {
		return 0
	}

	allTranslations = append(allTranslations, lang)
	return LanguageId(len(allTranslations))
}

func (l LanguageId) Translate(tag string) string {
	if l == 0 || int(l) > len(allTranslations) {
		return ""
	}
	return allTranslations[l-1].tags[tag]
}

func (l LanguageId) Unload() {
	clear(allTranslations[l-1].tags)
}

// private ========================================================

type lang struct{ tags map[string]string }

var allTranslations []lang
