package tiled

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/color"
)

func defaultText(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
func parseProperty(prop *internal.Property, project *Project) any {
	var value = prop.Value

	switch prop.Type {
	case "color":
		return color.Hex(prop.Value)
	case "class":
		// var class, hasClass = project.Classes[prop.CustomType]
		// if hasClass {
		// 	for n, v := range class.(map[string]any) {
		// 		var _, hasMember = value.(map[string]any)[n]
		// 		if !hasMember { // fill default member values, skip overwritten ones
		// 			value.(map[string]any)[n] = v
		// 		}
		// 	}
		// }
		return value
	}
	return value
}
