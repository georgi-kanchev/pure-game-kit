package gui

type container struct {
	widget
	Buttons []button `xml:"Button"`
}

func (c *container) findWidget(id string) *widget {
	if len(c.widget.Properties) > 0 && c.widget.Properties[0].Value == id {
		return &c.widget
	}

	for _, v := range c.Buttons {
		if len(v.widget.Properties) > 0 && v.Properties[0].Value == id {
			return &v.widget
		}
	}
	return nil
}
