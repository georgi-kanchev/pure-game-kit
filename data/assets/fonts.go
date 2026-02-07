package assets

import (
	"pure-game-kit/data/file"
	"pure-game-kit/debug"
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadedFontIds() []string {
	return collection.MapKeys(internal.Fonts)
}

func LoadFont(size int, filePath string) string {
	tryCreateWindow()

	var _, has = internal.Fonts[filePath]
	if has {
		return filePath
	}

	if !file.Exists(filePath) {
		debug.LogError("Failed to find font file: \"", filePath, "\"")
		return ""
	}

	var bytes = file.LoadBytes(filePath)
	var success = loadFont(filePath, size, bytes)

	if !success {
		debug.LogError("Failed to load font file: \"", filePath, "\"")
		return ""
	}
	return filePath
}
func UnloadFont(fontId string) {
	var font, has = internal.Fonts[fontId]

	if has && !isDefault(fontId) {
		delete(internal.Fonts, fontId)
		rl.UnloadFont(*font)
	}
}

//=================================================================
// private

const latin = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const latinPlus = "ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÑÒÓÔÕÖØÙÚÛÜÝÞßŒŠŽŁŃŚŹŻĆČĐŐŰàáâãäåæçèéêëìíîïñòóôõöøùúûüýþßœšžłńśźżćčđőűáéíóúüñãõçâêôÁÉÍÓÚÜÑÃÕÇÂÊÔßẞøØåÅþÞðÐœŒ"
const digits = "0123456789⁰¹²³⁴⁵⁶⁷⁸⁹₀₁₂₃₄₅₆₇₈₉¼½¾⅐⅑⅒⅓⅔⅕⅖⅗⅘⅙⅚⅛⅜⅝⅞"
const punct = " \t\n.,;:!?¡¿\"'()[]{}<>-/\\@#$€£%^&*_+=|~`…•™§©®°" + "–—‑′″‰ˆ˜“”‘’"
const greek = "ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ" + "αβγδεζηθικλμνξοπρστυφχψω" + "ς"
const cyrillic = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ" + "абвгдеёжзийклмнопрстуфхцчшщъыьэюя" + "ҐЄІЇґєії"
const georgian = "აბგდევზთიკლმნოპჟრსტუფქღყშჩცძწჭხჯჰ"
const armenian = "ԱԲԳԴԵԶԷԸԹԺԻԼԽԾԿՀՁՂՃՄՅՆՇՈՉՊՋՌՍՎՏՐՑՒՓՔՕՖ" + "աբգդեզէըթժիլխծկհձղճմյնշոչպջռսվտրցւփքօֆ"
const currencies = "$€£₴₽₲₵₡₢₣₤₥₦₧₨₩₪₫₭₮₯₰₱₲₳₴₸₺₼₽¢"
const extra = "ºª«»¶±×÷=≠<>≤≥∞∑∏√∫∆∂∇≈≡∈∉∪∩∧∨¬⇒⇔∀∃⊂⊆∅←↑→↓↔↕♠♥♦♣☺☹░▒▓│┤╡╢╖╕╣║╗╝┐└┴┬├─┼ˉ˙·"
const all = punct + extra + currencies + digits + latin + latinPlus + cyrillic + greek + georgian + armenian

func loadFont(id string, size int, bytes []byte) bool {
	tryCreateWindow()

	var characters = uniqueRunes(all)
	var glyphs = rl.LoadFontData(bytes, int32(size), characters, int32(len(characters)), rl.FontSdf)
	var recs = make([]*rl.Rectangle, len(glyphs))
	var atlas = rl.GenImageFontAtlas(glyphs, recs, int32(size), 0, 1)
	var font = rl.Font{BaseSize: int32(size), CharsCount: int32(len(glyphs)), Chars: &glyphs[0], Recs: recs[0]}

	font.Texture = rl.LoadTextureFromImage(&atlas)
	rl.UnloadImage(&atlas)
	rl.SetTextureFilter(font.Texture, rl.FilterBilinear)

	if font.BaseSize != 0 {
		internal.Fonts[id] = &font
	}

	return font.BaseSize != 0
}

func uniqueRunes(str string) []rune {
	var seen = make(map[rune]bool)
	var unique []rune

	for _, r := range str {
		if !seen[r] {
			seen[r] = true
			unique = append(unique, r)
		}
	}
	return unique
}
