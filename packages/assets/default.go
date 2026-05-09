package assets

import (
	"pure-game-kit/packages/utility/storage"

	_ "embed"
)

func LoadDefaultFont() (fontId string) {
	loadFont("", 49, storage.DecompressGZIP(font))
	return ""
}

// private ========================================================

//go:embed default/font.ttf.gz
var font []byte
