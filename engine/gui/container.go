package gui

type container struct {
	widget
	Widgets []widget `xml:"Widget"`
}

func (c *container) findWidget(id string) *widget {
	if len(c.widget.Properties) > 0 && c.widget.Properties[0].Value == id {
		return &c.widget
	}

	for _, v := range c.Widgets {
		if len(v.Properties) > 0 && v.Properties[0].Value == id {
			return &v
		}
	}
	return nil
}
