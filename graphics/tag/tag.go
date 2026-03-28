package tag

import (
	col "pure-game-kit/utility/color"
	t "pure-game-kit/utility/text"
)

func Underline(text string) string { return t.New("{_}", text, "{_}") }
func Crossout(text string) string  { return t.New("{-}", text, "{-}") }

func Thin(text string) string     { return t.New("{weight=thin}", text, "{weight}") }
func Regular(text string) string  { return t.New("{weight=regular}", text, "{weight}") }
func SemiBold(text string) string { return t.New("{weight=semiBold}", text, "{weight}") }
func Bold(text string) string     { return t.New("{weight=bold}", text, "{weight}") }

func Asset(assetId string) string {
	return t.New("{assetId=", assetId, "}@") // @ is a placeholder symbol that gets replaced with the asset
}

func Color(text string, color uint) string {
	var r, g, b, a = col.Channels(color)
	return t.New("{color=", r, " ", g, " ", b, " ", a, "}", text, "{color}")
}
func BackColor(text string, color uint) string {
	var r, g, b, a = col.Channels(color)
	return t.New("{backColor=", r, " ", g, " ", b, " ", a, "}", text, "{backColor}")
}
func OutlineColor(text string, color uint) string {
	var r, g, b, a = col.Channels(color)
	return t.New("{outlineColor=", r, " ", g, " ", b, " ", a, "}", text, "{outlineColor}")
}
func ShadowColor(text string, color uint) string {
	var r, g, b, a = col.Channels(color)
	return t.New("{shadowColor=", r, " ", g, " ", b, " ", a, "}", text, "{shadowColor}")
}

func OutlineThin(text string) string {
	return t.New("{outlineWeight=thin}", text, "{outlineWeight}")
}
func OutlineRegular(text string) string {
	return t.New("{outlineWeight=regular}", text, "{outlineWeight}")
}
func OutlineSemiBold(text string) string {
	return t.New("{outlineWeight=semiBold}", text, "{outlineWeight}")
}
func OutlineBold(text string) string {
	return t.New("{outlineWeight=bold}", text, "{outlineWeight}")
}

func ShadowThin(text string) string {
	return t.New("{shadowWeight=thin}", text, "{shadowWeight}")
}
func ShadowRegular(text string) string {
	return t.New("{shadowWeight=regular}", text, "{shadowWeight}")
}
func ShadowSemiBold(text string) string {
	return t.New("{shadowWeight=semiBold}", text, "{shadowWeight}")
}
func ShadowBold(text string) string {
	return t.New("{shadowWeight=bold}", text, "{shadowWeight}")
}

func ShadowBlurNone(text string) string {
	return t.New("{shadowBlur=0}", text, "{shadowBlur}")
}
func ShadowBlurLow(text string) string {
	return t.New("{shadowBlur=1}", text, "{shadowBlur}")
}
func ShadowBlurMedium(text string) string {
	return t.New("{shadowBlur=2}", text, "{shadowBlur}")
}
func ShadowBlurHigh(text string) string {
	return t.New("{shadowBlur=3}", text, "{shadowBlur}")
}
