package assets

func LoadDefaultAtlasUI(smooth bool) (atlasId string, tileIds []string, boxIds []string) {
	const symbol = "!"
	var tex = loadTexture(symbol, ui, smooth)
	var id = SetTextureAtlas(tex, 16, 16, 0)
	var t = []string{
		"out1-tl", "out1-t", "out1-tr", "out2-tl", "out2-t", "out2-tr", "out3-tl", "out3-t", "out3-tr",
		"out1-l", "out1-c", "out1-r", "out2-l", "out2-c", "out2-r", "out3-l", "out3-c", "out3-r",
		"out1-bl", "out1-b", "out1-br", "out2-bl", "out2-b", "out2-br", "out3-bl", "out3-b", "out3-br",
		"out1+tl", "out1+t", "out1+tr", "out2+tl", "out2+t", "out2+tr", "out3+tl", "out3+t", "out3+tr",
		"out1+bl", "out1+b", "out1+br", "out2+bl", "out2+b", "out2+br", "out3+bl", "out3+b", "out3+br",
		"in-tl", "in-t", "in-tr", "step-l", "step-c", "step-r", "bar-l", "bar-c", "bar-r",
		"in-l", "in-c", "in-r", "circle-tl", "circle-tr", "dot", "divider-l", "divider-c", "divider-r",
		"in-bl", "in-b", "in-br", "circle-bl", "circle-br", "handle1", "handle2", "", "",
	}

	for i := range t {
		if t[i] != "" {
			t[i] = symbol + t[i]
		}
	}
	var tiles = SetTextureAtlasTiles(id, 0, 0, t...)
	boxIds = []string{
		SetTextureBox(symbol+"out1", [9]string{t[0], t[1], t[2], t[9], t[10], t[11], t[18], t[19], t[20]}),
		SetTextureBox(symbol+"out1-", [9]string{t[27], t[28], t[29], t[9], t[10], t[11], t[18], t[19], t[20]}),
		SetTextureBox(symbol+"out1+", [9]string{t[0], t[1], t[2], t[9], t[10], t[11], t[36], t[37], t[38]}),
		SetTextureBox(symbol+"out2", [9]string{t[3], t[4], t[5], t[12], t[13], t[14], t[21], t[22], t[23]}),
		SetTextureBox(symbol+"out2-", [9]string{t[30], t[31], t[32], t[12], t[13], t[14], t[21], t[22], t[23]}),
		SetTextureBox(symbol+"out2+", [9]string{t[3], t[4], t[5], t[12], t[13], t[14], t[39], t[40], t[41]}),
		SetTextureBox(symbol+"out3", [9]string{t[6], t[7], t[8], t[15], t[16], t[17], t[24], t[25], t[26]}),
		SetTextureBox(symbol+"out3-", [9]string{t[33], t[34], t[35], t[15], t[16], t[17], t[24], t[25], t[26]}),
		SetTextureBox(symbol+"out3+", [9]string{t[6], t[7], t[8], t[15], t[16], t[17], t[42], t[43], t[44]}),
		SetTextureBox(symbol+"in", [9]string{t[45], t[46], t[47], t[54], t[55], t[56], t[63], t[64], t[65]}),
		SetTextureBox(symbol+"step", [9]string{t[48], t[52], t[50], t[48], t[52], t[50], t[48], t[52], t[50]}),
		SetTextureBox(symbol+"bar", [9]string{t[51], t[52], t[53], t[51], t[52], t[53], t[51], t[52], t[53]}),
		SetTextureBox(symbol+"divider", [9]string{t[60], t[61], t[62], t[60], t[61], t[62], t[60], t[61], t[62]}),
	}

	return id, tiles, boxIds
}

const ui = `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`
