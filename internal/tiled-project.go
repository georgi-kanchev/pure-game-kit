package internal

type Project struct {
	Properties  []*ProjectProperty     `json:"properties"`
	CustomTypes []*ProjectPropertyType `json:"propertyTypes"`
}

type ProjectProperty struct {
	Name       string `json:"name"`
	CustomType string `json:"propertytype"`
	Type       string `json:"type"`
	Value      any    `json:"value"` // can be int, string, or object
}

type ProjectPropertyType struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	StorageType string `json:"storageType"`
	Type        string `json:"type"`

	EnumValues        []string `json:"values"`
	EnumValuesAsFlags bool     `json:"valuesAsFlags"`

	ClassColor    string                `json:"color"`
	ClassDrawFill bool                  `json:"drawFill"`
	ClassMembers  []*ProjectClassMember `json:"members"`
	ClassUseAs    []string              `json:"useAs"`
}

type ProjectClassMember struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value any    `json:"value"`
}
