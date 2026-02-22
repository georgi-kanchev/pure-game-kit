package tag

import (
	col "pure-game-kit/utility/color"
	t "pure-game-kit/utility/text"
)

func Asset(assetId string) string { return t.New("{assetId=", assetId, "}") }
func Underline(text string, size float32) string {
	return t.New("{underline=", size, "}", text, "{underline}")
}
func Strikethrough(text string, size float32) string { return t.New("{-", size, "}", text, "{-}") }
func Box(text string, size float32) string           { return t.New("{=", size, "}", text, "{=}") }

func Color(text string, color uint) string {
	var r, g, b, a = col.Channels(color)
	return ColorRGBA(text, r, g, b, a)
}
func ColorRGB(text string, r, g, b, a byte) string {
	return ColorRGBA(text, r, g, b, 255)
}
func ColorRGBA(text string, r, g, b, a byte) string {
	return t.New("{color=", r, " ", g, " ", b, " ", a, "}", text, "{color}")
}
