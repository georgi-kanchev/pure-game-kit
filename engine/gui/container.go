package gui

type container struct {
	widget
	Buttons []button `xml:"Button"`
}

func Container(id, x, y, width, height string, properties ...string) string {
	return newWidget("Container", id, x, y, width, height) + extraProps(properties...) + ">"
}
func ContainerEnd() string { return "</Container>" }

func (c *container) FindWidget(id string) *widget {
	for _, v := range c.Buttons {
		if v.Properties[0].Value == id {
			return &v.widget
		}
	}
	return nil
}
