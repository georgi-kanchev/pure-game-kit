package tiled

import (
	"pure-game-kit/internal"
	"pure-game-kit/utility/color"
)

func defaultValueText(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
func parseProperty(prop *internal.Property, project *Project) any {
	switch prop.Type {
	case "color":
		return color.Hex(prop.Value)
	case "class":
		var class, hasClass = project.Classes[prop.CustomType]
		if !hasClass {
			return prop.Value
		}

		var classMembers = class.(map[string]any)
		var result = make(map[string]any, len(classMembers))
		for n, v := range classMembers {
			result[n] = v

			for _, prop := range prop.Properties {
				if prop.Name == n {
					result[prop.Name] = parseProperty(prop, project)
				}
			}
		}
		return result
	}
	return prop.Value
}
