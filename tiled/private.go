package tiled

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/color"
	"pure-game-kit/utility/text"
)

func defaultText(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
func parseProperty(prop *internal.Property) any {
	switch prop.Type {
	case "float":
		return text.ToNumber[float32](prop.Value)
	case "int", "object":
		return text.ToNumber[int](prop.Value)
	case "bool":
		return prop.Value == "true"
	case "color":
		return color.Hex(prop.Value)
	case "class":
		var result = make(map[string]any, len(prop.Properties))
		for _, p := range prop.Properties {
			result[p.Name] = parseProperty(p)
		}
		return ""
	}
	return prop.Value
}
