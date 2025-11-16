package assets

func LoadDefaultAtlasUI() (atlasId string, tileIds []string, boxIds []string) {
	var tex = loadTexture(defaultUI, ui)
	var id = SetTextureAtlas(tex, 16, 16, 0)
	var t = []string{
		"out1-tl", "out1-t", "out1-tr", "out2-tl", "out2-t", "out2-tr", "out3-tl", "out3-t", "out3-tr",
		"out1-l", "out1-c", "out1-r", "out2-l", "out2-c", "out2-r", "out3-l", "out3-c", "out3-r",
		"out1-bl", "out1-b", "out1-br", "out2-bl", "out2-b", "out2-br", "out3-bl", "out3-b", "out3-br",
		"out1+tl", "out1+t", "out1+tr", "out2+tl", "out2+t", "out2+tr", "out3+tl", "out3+t", "out3+tr",
		"out1+bl", "out1+b", "out1+br", "out2+bl", "out2+b", "out2+br", "out3+bl", "out3+b", "out3+br",
		"in-tl", "in-t", "in-tr", "step-l", "step-c", "step-r", "circle-tl", "circle-tr", "dot",
		"in-l", "in-c", "in-r", "bar-l", "bar-c", "bar-r", "circle-bl", "circle-br", "handle-t",
		"in-bl", "in-b", "in-br", "divider-l", "divider-c", "divider-r", "handle1", "handle2", "handle-b",
	}

	for i := range t {
		if t[i] != "" {
			t[i] = defaultUI + t[i]
		}
	}
	var tiles = SetTextureAtlasTiles(id, 0, 0, t...)
	boxIds = []string{
		SetTextureBox(defaultUI+"out1", [9]string{t[0], t[1], t[2], t[9], t[10], t[11], t[18], t[19], t[20]}),
		SetTextureBox(defaultUI+"out1-", [9]string{t[27], t[28], t[29], t[9], t[10], t[11], t[18], t[19], t[20]}),
		SetTextureBox(defaultUI+"out1+", [9]string{t[0], t[1], t[2], t[9], t[10], t[11], t[36], t[37], t[38]}),
		SetTextureBox(defaultUI+"out2", [9]string{t[3], t[4], t[5], t[12], t[13], t[14], t[21], t[22], t[23]}),
		SetTextureBox(defaultUI+"out2-", [9]string{t[30], t[31], t[32], t[12], t[13], t[14], t[21], t[22], t[23]}),
		SetTextureBox(defaultUI+"out2+", [9]string{t[3], t[4], t[5], t[12], t[13], t[14], t[39], t[40], t[41]}),
		SetTextureBox(defaultUI+"out3", [9]string{t[6], t[7], t[8], t[15], t[16], t[17], t[24], t[25], t[26]}),
		SetTextureBox(defaultUI+"out3-", [9]string{t[33], t[34], t[35], t[15], t[16], t[17], t[24], t[25], t[26]}),
		SetTextureBox(defaultUI+"out3+", [9]string{t[6], t[7], t[8], t[15], t[16], t[17], t[42], t[43], t[44]}),
		SetTextureBox(defaultUI+"in", [9]string{t[45], t[46], t[47], t[54], t[55], t[56], t[63], t[64], t[65]}),
		SetTextureBox(defaultUI+"step", [9]string{"", "", "", t[48], t[58], t[50], "", "", ""}),
		SetTextureBox(defaultUI+"bar", [9]string{"", "", "", t[57], t[58], t[59], "", "", ""}),
		SetTextureBox(defaultUI+"divider", [9]string{"", "", "", t[66], t[67], t[68], "", "", ""}),
		SetTextureBox(defaultUI+"circle", [9]string{t[51], "", t[52], "", "", "", t[60], "", t[61]}),
		SetTextureAtlasTile(id, defaultUI+"handle", 8, 6, 1, 2, 0, false),
	}

	return id, tiles, boxIds
}

const ui = `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`
