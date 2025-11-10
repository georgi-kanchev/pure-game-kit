package tiled

import (
	"pure-game-kit/debug"
	"pure-game-kit/internal"
	"pure-game-kit/utility/collection"
)

type Project struct {
	Properties map[string]any
	Classes    map[string]any
}

func NewProject(projectId string) *Project {
	var data, _ = internal.TiledProjects[projectId]
	if data == nil {
		debug.LogError("Failed to create project: \"", projectId, "\"\nNo data is loaded with this project id.")
		return nil
	}

	var result = Project{}
	result.initClasses(data)
	result.initProperties(data)
	return &result
}

//=================================================================
// private

func (p *Project) initProperties(data *internal.Project) {
	p.Properties = make(map[string]any)

	// for _, prop := range data.Properties {
	// p.Properties[prop.Name] = parseProperty(propToProject(prop), p)
	// }
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
