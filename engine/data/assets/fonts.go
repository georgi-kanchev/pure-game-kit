package assets

import (
	"pure-kit/engine/data/file"
	"pure-kit/engine/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SymbolsLatin = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const SymbolsLatinExtra = "ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÑÒÓÔÕÖØÙÚÛÜÝÞßŒŠŽŁŃŚŹŻĆČĐŐŰ" +
	"àáâãäåæçèéêëìíîïñòóôõöøùúûüýþßœšžłńśźżćčđőű" + "áéíóúüñãõçâêôÁÉÍÓÚÜÑÃÕÇÂÊÔ" + "ßẞ" + "øØåÅ" + "þÞðÐ" + "œŒ"
const SymbolsDigits = "0123456789"
const SymbolsPunctuation = " \t\n.,;:!¡¿\"'()[]{}<>-/\\@#$€£%^&*_+=|~`…•™§©®°" + "–—‑′″‰ˆ˜“”‘’"
const SymbolsGreek = "ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ" + "αβγδεζηθικλμνξοπρστυφχψω" + "ς"
const SymbolsCyrillic = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ" + "абвгдеёжзийклмнопрстуфхцчшщъыьэюя" + "ҐЄІЇґєії"
const SymbolsGeorgian = "აბგდევზთიკლმნოპჟრსტუფქღყშჩცძწჭხჯჰ"
const SymbolsArmenian = "ԱԲԳԴԵԶԷԸԹԺԻԼԽԾԿՀՁՂՃՄՅՆՇՈՉՊՋՌՍՎՏՐՑՒՓՔՕՖ" + "աբգդեզէըթժիլխծկհձղճմյնշոչպջռսվտրցւփքօֆ"
const SymbolsCurrencies = "$€£₴₽₲₵₡₢₣₤₥₦₧₨₩₪₫₭₮₯₰₱₲₳₴₸₺₼₽¢"
const SymbolsExtra = "ºª«»¶±×÷=≠<>≤≥∞∑∏√∫∆∂∇≈≡∈∉∪∩∧∨¬⇒⇔∀∃⊂⊆∅←↑→↓↔↕♠♥♦♣☺☹"
const symbolsDefault = SymbolsPunctuation + SymbolsDigits + SymbolsLatin

func LoadFonts(size int, filePaths ...string) []string {
	var result = []string{}
	for _, path := range filePaths {
		var id = LoadFontSymbols(size, path, symbolsDefault)
		if id != "" {
			result = append(result, id)
		}
	}
	return result
}
func LoadFontSymbols(size int, filePath string, symbols ...string) string {
	tryCreateWindow()

	var id, absolutePath = getIdPath(filePath)
	var _, has = internal.Fonts[id]

	if has || !file.Exists(absolutePath) {
		return ""
	}

	var allSymbols = "?" // it's good to have ? as first character in any case to account for missing symbols
	for _, s := range symbols {
		allSymbols += s
	}

	if len(symbols) == 0 {
		allSymbols += SymbolsPunctuation + SymbolsExtra + SymbolsCurrencies + SymbolsDigits + SymbolsLatin +
			SymbolsLatinExtra + SymbolsCyrillic + SymbolsGreek + SymbolsGeorgian + SymbolsArmenian
	}

	var font = rl.LoadFontEx(absolutePath, int32(size), uniqueRunes(allSymbols))
	if font.BaseSize == 0 {
		return ""
	}

	internal.Fonts[id] = &font
	return id
}

func UnloadFonts(fontIds ...string) {
	for _, v := range fontIds {
		var font, has = internal.Fonts[v]

		if has {
			delete(internal.Fonts, v)
			rl.UnloadFont(*font)
		}
	}
}

// #region private

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

// #endregion
