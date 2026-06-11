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

	nextId++
	var id = LanguageId(nextId)
	allTranslations[id] = lang
	return id
}

func (l LanguageId) Translate(tag string) string {
	var value, has = allTranslations[l]
	if !has {
		return ""
	}
	return value.tags[tag]
}
func (l LanguageId) Unload() {
	delete(allTranslations, l)
}

// private ========================================================

type lang struct{ tags map[string]string }

var allTranslations map[LanguageId]lang
var nextId uint32
