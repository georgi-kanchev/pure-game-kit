package gui

type button struct {
	widget
}

func NewButton(id, x, y, width, height string, properties ...string) string {
	return newWidget("Button", id, x, y, width, height) + extraProps(properties...) + " />"
}
