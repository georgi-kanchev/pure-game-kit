package tiled

import (
	"pure-game-kit/debug"
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"
	"pure-game-kit/utility/color"
)

type Project struct {
	Properties     map[string]any
	Classes        map[string]any      // a collection of custom properties that anything in the project can use
	UniqueTilesets map[string]*Tileset // maps in the same project will try to reuse these instead of loading them
}

func NewProject(projectId string) *Project {
	var data, _ = internal.TiledProjects[projectId]
	if data == nil {
		debug.LogError("Failed to create project: \"", projectId, "\"\nNo data is loaded with this project id.")
		return nil
	}

	var result = Project{UniqueTilesets: map[string]*Tileset{}}
	result.initClasses(data)
	result.initProperties(data)
	return &result
}

//=================================================================
// private

func (p *Project) initProperties(data *internal.Project) {
	p.Properties = make(map[string]any)

	for _, prop := range data.Properties {
		p.Properties[prop.Name] = p.parseProjectProperty(prop)
	}
}
func (p *Project) initClasses(data *internal.Project) {
	p.Classes = make(map[string]any, len(data.CustomTypes))

	for _, t := range data.CustomTypes {
		if t.Type == "enum" {
			p.Classes[t.Name] = collection.Clone(t.EnumValues)
			continue
		}

		if t.Type == "class" {
			var classes = make(map[string]any)
			for _, m := range t.ClassMembers {
				classes[m.Name] = m.Value
			}
			p.Classes[t.Name] = classes
			continue
		}

	}
}

func (p *Project) parseProjectProperty(prop *internal.ProjectProperty) any {
	switch prop.Type {
	case "color":
		return color.Hex(prop.Value.(string))
	case "class":
		var class, hasClass = p.Classes[prop.CustomType]
		if !hasClass {
			return prop.Value
		}

		var classMembers = class.(map[string]any)
		var result = make(map[string]any, len(classMembers))
		for n, v := range classMembers {
			var members = prop.Value.(map[string]any)
			var m, hasMember = members[n]
			if hasMember {
				result[n] = m
				continue
			}

			result[n] = v
		}
		return result
	}
	return prop.Value
}
