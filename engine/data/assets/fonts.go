package assets

import (
	"pure-kit/engine/data/file"
	"pure-kit/engine/internal"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const latin = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const latinPlus = "ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÑÒÓÔÕÖØÙÚÛÜÝÞßŒŠŽŁŃŚŹŻĆČĐŐŰ" +
	"àáâãäåæçèéêëìíîïñòóôõöøùúûüýþßœšžłńśźżćčđőű" + "áéíóúüñãõçâêôÁÉÍÓÚÜÑÃÕÇÂÊÔ" + "ßẞ" + "øØåÅ" + "þÞðÐ" + "œŒ"
const digits = "0123456789"
const punctuation = " \t\n.,;:!¡¿\"'()[]{}<>-/\\@#$€£%^&*_+=|~`…•™§©®°" + "–—‑′″‰ˆ˜“”‘’"
const greek = "ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩ" + "αβγδεζηθικλμνξοπρστυφχψω" + "ς"
const cyrillic = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯ" + "абвгдеёжзийклмнопрстуфхцчшщъыьэюя" + "ҐЄІЇґєії"
const georgian = "აბგდევზთიკლმნოპჟრსტუფქღყშჩცძწჭხჯჰ"
const armenian = "ԱԲԳԴԵԶԷԸԹԺԻԼԽԾԿՀՁՂՃՄՅՆՇՈՉՊՋՌՍՎՏՐՑՒՓՔՕՖ" + "աբգդեզէըթժիլխծկհձղճմյնշոչպջռսվտրցւփքօֆ"
const currencies = "$€£₴₽₲₵₡₢₣₤₥₦₧₨₩₪₫₭₮₯₰₱₲₳₴₸₺₼₽¢"
const extra = "ºª«»¶±×÷=≠<>≤≥∞∑∏√∫∆∂∇≈≡∈∉∪∩∧∨¬⇒⇔∀∃⊂⊆∅←↑→↓↔↕♠♥♦♣☺☹"
const all = punctuation + extra + currencies + digits + latin + latinPlus + cyrillic + greek + georgian + armenian

func LoadFonts(size int, filePaths ...string) []string {
	tryCreateWindow()
	var result = []string{}

	if internal.ShaderText.ID == 0 {
		internal.ShaderText = rl.LoadShaderFromMemory("", frag)
	}

	for _, path := range filePaths {
		var id, absolutePath = getIdPath(path)
		var _, has = internal.Fonts[id]

		if has || !file.Exists(absolutePath) {
			continue
		}

		var bytes = file.LoadBytes(path)
		var characters = uniqueRunes(all)
		var glyphs = rl.LoadFontData(bytes, int32(size), characters, int32(len(characters)), rl.FontSdf)
		var font = rl.Font{BaseSize: int32(size), CharsCount: int32(len(characters)), Chars: &glyphs[0]}
		var atlas = rl.GenImageFontAtlas(
			unsafe.Slice(font.Chars, font.CharsCount),
			unsafe.Slice(&font.Recs, font.CharsCount),
			int32(size), 0, 1,
		)
		font.Texture = rl.LoadTextureFromImage(&atlas)
		rl.UnloadImage(&atlas)
		rl.SetTextureFilter(font.Texture, rl.FilterBilinear)

		if font.BaseSize == 0 {
			continue
		}

		result = append(result, id)
		internal.Fonts[id] = &font
	}
	return result
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

const frag = `#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

uniform sampler2D texture0;
uniform float smoothness;
uniform float thickness;

out vec4 finalColor;

void main()
{
    float distance = texture(texture0, fragTexCoord).a - (1.0 - thickness);
    float baseSmooth = smoothness * length(vec2(dFdx(distance), dFdy(distance)));
    float alpha = smoothstep(-baseSmooth, baseSmooth, distance);
    vec4 fill = vec4(fragColor.rgb, fragColor.a * alpha);
    
    finalColor = fill;
}`

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
