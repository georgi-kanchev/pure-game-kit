package assets

func LoadDefaultAtlasUI() (atlasId string, tileIds []string, boxIds []string) {
	const symbol = "!"
	var tex = loadTexture(symbol, ui)
	var id = LoadTextureAtlas(tex, 16, 16, 0)
	var t = []string{
		"out1-tl", "out1-t", "out1-tr", "out2-tl", "out2-t", "out2-tr", "out3-tl", "out3-t", "out3-tr",
		"out1-l", "out1-c", "out1-r", "out2-l", "out2-c", "out2-r", "out3-l", "out3-c", "out3-r",
		"out1-bl", "out1-b", "out1-br", "out2-bl", "out2-b", "out2-br", "out3-bl", "out3-b", "out3-br",
		"out1+tl", "out1+t", "out1+tr", "out2+tl", "out2+t", "out2+tr", "out3+tl", "out3+t", "out3+tr",
		"out1+bl", "out1+b", "out1+br", "out2+bl", "out2+b", "out2+br", "out3+bl", "out3+b", "out3+br",
		"in-tl", "in-t", "in-tr", "step-l", "step-c", "step-r", "bar-l", "bar-c", "bar-r",
		"in-l", "in-c", "in-r", "circle-tl", "circle-tr", "dot", "divider-l", "divider-c", "divider-r",
		"in-bl", "in-b", "in-br", "circle-bl", "circle-br", "disable", "handle1", "handle2", "",
	}

	for i := range t {
		if t[i] != "" {
			t[i] = symbol + t[i]
		}
	}
	var tiles = LoadTextureAtlasTiles(id, 0, 0, t...)
	boxIds = []string{
		LoadTextureBox(symbol+"out1", [9]string{t[0], t[1], t[2], t[9], t[10], t[11], t[18], t[19], t[20]}),
		LoadTextureBox(symbol+"out1-", [9]string{t[27], t[28], t[29], t[9], t[10], t[11], t[18], t[19], t[20]}),
		LoadTextureBox(symbol+"out1+", [9]string{t[0], t[1], t[2], t[9], t[10], t[11], t[36], t[37], t[38]}),
		LoadTextureBox(symbol+"out2", [9]string{t[3], t[4], t[5], t[12], t[13], t[14], t[21], t[22], t[23]}),
		LoadTextureBox(symbol+"out2-", [9]string{t[30], t[31], t[32], t[12], t[13], t[14], t[21], t[22], t[23]}),
		LoadTextureBox(symbol+"out2+", [9]string{t[3], t[4], t[5], t[12], t[13], t[14], t[39], t[40], t[41]}),
		LoadTextureBox(symbol+"out3", [9]string{t[6], t[7], t[8], t[15], t[16], t[17], t[24], t[25], t[26]}),
		LoadTextureBox(symbol+"out3-", [9]string{t[33], t[34], t[35], t[15], t[16], t[17], t[24], t[25], t[26]}),
		LoadTextureBox(symbol+"out3+", [9]string{t[6], t[7], t[8], t[15], t[16], t[17], t[42], t[43], t[44]}),
		LoadTextureBox(symbol+"in", [9]string{t[36], t[37], t[38], t[45], t[46], t[47], t[54], t[55], t[56]}),
		LoadTextureBox(symbol+"step", [9]string{t[39], t[43], t[41], t[39], t[43], t[41], t[39], t[43], t[41]}),
		LoadTextureBox(symbol+"bar", [9]string{t[42], t[43], t[44], t[42], t[43], t[44], t[42], t[43], t[44]}),
		LoadTextureBox(symbol+"divider", [9]string{t[51], t[52], t[53], t[51], t[52], t[53], t[51], t[52], t[53]}),
	}

	return id, tiles, boxIds
}

const ui = `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`
