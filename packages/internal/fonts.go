package internal

import (
	_ "embed"
	"encoding/xml"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Font struct {
	AtlasId  int32 // see assets.ImageId
	Chars    map[rune]Char
	Kernings map[rune]Kerning
}
type Char struct {
	ID       rune   `xml:"id,attr"`
	Index    int    `xml:"index,attr"`
	Char     string `xml:"char,attr"`
	Width    int    `xml:"width,attr"`
	Height   int    `xml:"height,attr"`
	XOffset  int    `xml:"xoffset,attr"`
	YOffset  int    `xml:"yoffset,attr"`
	XAdvance int    `xml:"xadvance,attr"`
	Chnl     int    `xml:"chnl,attr"`
	X        int    `xml:"x,attr"`
	Y        int    `xml:"y,attr"`
	Page     int    `xml:"page,attr"`
}
type Kerning struct {
	First  rune `xml:"first,attr"`
	Second rune `xml:"second,attr"`
	Amount int  `xml:"amount,attr"`
}
type FontXML struct {
	XMLName xml.Name `xml:"font"`

	Info struct {
		Face     string `xml:"face,attr"`
		Size     int    `xml:"size,attr"`
		Bold     int    `xml:"bold,attr"`
		Italic   int    `xml:"italic,attr"`
		Charset  string `xml:"charset,attr"`
		Unicode  int    `xml:"unicode,attr"`
		StretchH int    `xml:"stretchH,attr"`
		Smooth   int    `xml:"smooth,attr"`
		AA       int    `xml:"aa,attr"`
		Padding  string `xml:"padding,attr"`
		Spacing  string `xml:"spacing,attr"`
		Outline  int    `xml:"outline,attr"`
	} `xml:"info"`

	Common struct {
		LineHeight int `xml:"lineHeight,attr"`
		Base       int `xml:"base,attr"`
		ScaleW     int `xml:"scaleW,attr"`
		ScaleH     int `xml:"scaleH,attr"`
		Pages      int `xml:"pages,attr"`
		Packed     int `xml:"packed,attr"`
		AlphaChnl  int `xml:"alphaChnl,attr"`
		RedChnl    int `xml:"redChnl,attr"`
		GreenChnl  int `xml:"greenChnl,attr"`
		BlueChnl   int `xml:"blueChnl,attr"`
	} `xml:"common"`

	Pages []struct {
		ID   int    `xml:"id,attr"`
		File string `xml:"file,attr"`
	} `xml:"pages>page"`

	DistanceField struct {
		FieldType     string `xml:"fieldType,attr"`
		DistanceRange int    `xml:"distanceRange,attr"`
	} `xml:"distanceField"`

	Chars struct {
		Count int    `xml:"count,attr"`
		Chars []Char `xml:"char"`
	} `xml:"chars"`

	Kernings struct {
		Count    int       `xml:"count,attr"`
		Kernings []Kerning `xml:"kerning"`
	} `xml:"kernings"`
}

var Fonts = make(map[byte]Font) // 0 = default
var FontNextId byte

func LoadFont(fontData *FontXML, imageId int32, isDefault bool) byte {
	if !isDefault {
		FontNextId++
	}
	var id = FontNextId
	var font = Font{AtlasId: imageId, Chars: make(map[rune]Char), Kernings: make(map[rune]Kerning)}

	for _, char := range fontData.Chars.Chars {
		font.Chars[char.ID] = char
	}
	for _, kern := range fontData.Kernings.Kernings {
		font.Kernings[kern.First] = kern
	}

	Fonts[id] = font
	return id
}

// private ========================================================

//go:embed font.ttf.gz
var defaultFont []byte

const punct = " \t\n.,;:!?¡¿\"'()[]{}<>-/\\@#$€£%^&*_+=|~`"
const extra = "…•™§©®°–—‑′″‰ˆ˜“”‘’ºª«»¶±×÷≠≤≥∞∑∏√∫∆∂∇≈≡∈∉∪∩∧∨¬⇒⇔∀∃⊂⊆∅←↑→↓↔↕♠♥♦♣☺☹░▒▓│┤╡╢╖╕╣║╗╝┐└┴┬├─┼ˉ˙·"
const currencies = "₴₽₲₵₡₢₣₤₥₦₧₨₩₪₫₭₮₯₰₱₳₸₺₼¢"
const digits = "0123456789⁰¹²³⁴⁵⁶⁷⁸⁹₀₁₂₃₄₅₆₇₈₉¼½¾⅐⅑⅒⅓⅔⅕⅖⅗⅘⅙⅚⅛⅜⅝⅞"
const latin = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
const latinPlus = "ÀÁÂÃÄÅÆÇÈÉÊËÌÍÎÏÑÒÓÔÕÖØÙÚÛÜÝÞßŒŠŽŁŃŚŹŻĆČĐŐŰàáâãäåæçèéêëìíîïñòóôõöøùúûüýþœšžłńśźżćčđőűẞðÐ"
const cyrillic = "АБВГДЕЁЖЗИЙКЛМНОПРСТУФХЦЧШЩЪЫЬЭЮЯабвгдеёжзийклмнопрстуфхцчшщъыьэюяҐЄІЇґєії"
const greek = "ΑΒΓΔΕΖΗΘΙΚΛΜΝΞΟΠΡΣΤΥΦΧΨΩαβγδεζηθικλμνξοπρστυφχψως"
const georgian = "აბგდევზთიკლმნოპჟრსტუფქღყშჩცძწჭხჯჰ"
const armenian = "ԱԲԳԴԵԶԷԸԹԺԻԼԽԾԿՀՁՂՃՄՅՆՇՈՉՊՋՌՍՎՏՐՑՒՓՔՕՖաբգդեզէըթժիլխծկհձղճմյնշոչպջռսվտրցւփքօֆ"

const all = punct + extra + currencies + digits + latin + latinPlus + cyrillic + greek + georgian + armenian

func loadFont(id string, size int, bytes []byte) bool {
	var characters = []rune(all)
	var glyphs = rl.LoadFontData(bytes, int32(size), characters, int32(len(characters)), rl.FontSdf)
	var recs = make([]*rl.Rectangle, len(glyphs))
	var atlas = rl.GenImageFontAtlas(glyphs, recs, int32(size), 0, 1)
	var font = rl.Font{BaseSize: int32(size), CharsCount: int32(len(glyphs)), Chars: &glyphs[0], Recs: recs[0]}

	font.Texture = rl.LoadTextureFromImage(&atlas)
	rl.UnloadImage(&atlas)
	rl.SetTextureFilter(font.Texture, rl.FilterBilinear)

	// if font.BaseSize != 0 {
	// 	internal.Fonts[id] = font
	// }

	return font.BaseSize != 0
}
