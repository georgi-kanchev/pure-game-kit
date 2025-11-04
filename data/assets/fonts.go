package assets

import (
	"pure-game-kit/data/file"
	"pure-game-kit/debug"
	"pure-game-kit/internal"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadFont(size int, filePath string) []string {
	var result = []string{}
	var id = getIdPath(filePath)

	if !file.IsExisting(filePath) {
		debug.LogError("Failed to find font file: \"", filePath, "\"")
		return result
	}

	var bytes = file.LoadBytes(filePath)
	var success = loadFont(id, size, bytes)

	if success {
		result = append(result, id)
	} else {
		debug.LogError("Failed to load font file: \"", filePath, "\"")
	}

	return result
}
func UnloadFont(fontId string) {
	var font, has = internal.Fonts[fontId]

	if has {
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

func loadFont(id string, size int, bytes []byte) bool {
	tryCreateWindow()
	tryInitShader()

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
