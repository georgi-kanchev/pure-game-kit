package assets

import (
	"pure-game-kit/data/storage"
	"pure-game-kit/internal"
	"pure-game-kit/utility/text"

	_ "embed"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LoadDefaultFont() (fontId string) {
	loadFont(defaultFont, 49, storage.DecompressGZIP(font))
	return defaultFont
}
func LoadDefaultTexture() string {
	return loadTexture(defaultTexture, texture)
}
func LoadDefaultAtlasCursors() (atlasId string, tileIds []string) {
	var tex = loadTexture(defaultCursors, cursors)
	var id = SetTextureAtlas(tex, 32, 32, 0)
	var ids = []string{
		"pointer1", "pointer2", "pointer3", "pointer4", "pointer5", "pointer6", "pointer7", "pointer8",
		"pointer9", "pointer10", "pointer11", "pointer12", "pointer13", "pointer14", "pointer15", "pointer16",
		"pointer17", "pointer18", "pointer19",
		"hand1", "hand2", "hand3", "hand4", "hand5", "hand6", "hand7", "hand8", "hand9", "hand10",
		"hand11", "hand12", "hand13", "hand14", "hand15", "hand16", "hand17", "hand18", "hand19",
		"clock1", "clock2", "clock3", "clock4", "clock5", "clock6", "clock7", "clock8",
		"hourglass1", "hourglass2", "hourglass3", "hourglass4", "spinner1", "spinner2", "lock1", "lock2",
		"arrow1", "arrow2", "", "x1", "x2", "disabled1", "disabled2", "door1", "door2", "door3", "door4",
		"eye1", "eye2", "eye3", "eye4", "zoom", "zoom+", "zoom-", "zoom=", "arrow3", "arrow4", "",
		"refresh1", "refresh2", "refresh3", "refresh4", "bubble1", "bubble2", "bubble3", "bubble4",
		"pointer1-open", "pointer1-menu", "pointer1-gears", "pointer1+", "?", "!", "pointer1?", "pointer1!",
		"pointer7?", "pointer7!", "", "art-pencil", "art-pen", "art-bucket", "art-eraser", "art-pick",
		"art-brush1", "art-brush2", "water-can", "art-spray", "art-wand", "art-wrench",
		"crosshair1", "crosshair2", "crosshair3", "crosshair4", "crosshair5", "crosshair6", "crosshair7", "",
		"move1", "move2", "stairs1", "stairs2", "stairs3", "item-axe1", "item-axe2", "item-bomb", "item-bow",
		"item-hammer1", "item-hammer2", "item-pickaxe", "item-spade", "item-sword1", "item-sword2", "item-flashlight",
		"input1", "input2", "input3", "move1", "move2", "resize1", "resize2", "move3", "move4", "resize3", "resize4",
		"move5", "move6", "resize5", "resize6", "move7", "move8", "resize7", "resize8", "split1", "split2", "split3",
	}

	for i := range ids {
		if ids[i] != "" {
			ids[i] = defaultCursors + ids[i]
		}
	}
	var tiles = SetTextureAtlasTiles(id, 0, 0, ids...)
	return id, tiles
}
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
func LoadDefaultAtlasRetro() (atlasId string, tileIds []string) {
	var tex = loadTexture(defaultRetroAtlas, retro)
	var id = SetTextureAtlas(tex, 8, 8, 1)
	var ids = []string{
		"empty", "shade1", "shade2", "shade3", "shade4", "shade5", "shade6", "shade7", "shade8", "shade9", "full",
		"tile1", "tile2", "tile3", "tile4", "tile5", "tile6", "tile7", "tile8", "tile9",
		"tile10", "tile11", "tile12", "tile13", "tile14", "tile15", "tile16", "tile17", "tile18", "tile19",
		"tile20", "tile21", "tile22", "tile23", "tile24", "tile25", "tile26", "tile27", "tile28", "tile29",
		"tile30", "tile31", "tile32", "tile33", "tile34", "tile35", "tile36", "tile37", "tile38", "tile39",
		"tile40", "tile41", "tile42", "tile43", "tile44", "tile45", "tile46", "tile47", "tile48", "tile49",
		"tile50", "tile51", "tile52", "tile53", "tile54", "tile55", "tile56", "tile57", "tile58", "tile59",
		"tile60", "tile61", "tile62", "tile63", "tile64", "tile65", "tile66", "tile67",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V",
		"W", "X", "Y", "Z", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r",
		"s", "t", "u", "v", "w", "x", "y", "z", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "1/8", "1/7",
		"1/6", "1/5", "1/4", "1/3", "3/8", "2/5", "1/2", "3/5", "5/8", "2/3", "3/4", "4/5", "5/6", "7/8",
		"sub0", "sub1", "sub2", "sub3", "sub4", "sub5", "sub6", "sub7", "sub8", "sub9", "sub10", "sub11", "sub12",
		"sup0", "sup1", "sup2", "sup3", "sup4", "sup5", "sup6", "sup7", "sup8", "sup9", "sup10", "sup11", "sup12",
		"-", "+", "multiply", "over", "divide", "%", "=", "!=", "approximately", "sqrt", "func", "integral", "sum",
		"epsilon", "euler", "gold-ratio", "pi", "silver-ratio", "infinity", "<<", ">>", "<=", ">=", "shape-line",
		"", "", "<", ">", "(", ")", "[", "]", "{", "}", "perpendicular", "parallel", "angle", "angle-right", "~",
		"degree", "celsius", "fahrenheit", "*", "^", "#", "number", "$", "euro", "pound", "yen", "cent", "currency",
		"!", "?", ".", ",", "...", ":", ";", "\"", "'", "`", "-", "_", "|", "/", "\\", "@", "&",
		"registered", "copyright-audio", "copyright", "trademark", "", "", "", "", "",
		"pipe1-straight", "pipe1-corner", "pipe1-t-shaped", "pipe1-cross",
		"pipe2-straight", "pipe2-corner", "pipe2-t-shaped", "pipe2-cross",
		"pipe3-straight", "pipe3-corner", "pipe3-t-shaped", "pipe3-cross",
		"pipe4-straight", "pipe4-corner", "pipe4-t-shaped", "pipe4-cross",
		"pipe5-straight", "pipe5-corner", "pipe5-t-shaped", "pipe5-cross",
		"pipe6-straight", "pipe6-corner", "pipe6-t-shaped", "pipe6-cross", "", "",
		"bar1-edge", "bar1-straight", "bar2-edge", "bar2-straight", "bar3-edge", "bar3-straight",
		"bar4-edge", "bar4-straight", "bar5-edge", "bar5-straight", "bar6-edge", "bar6-straight",
		"bar7-edge", "bar7-straight", "bar8-edge", "bar8-straight", "bar9-edge", "bar9-straight",
		"bar10-edge", "bar10-straight", "bar11-edge", "bar11-straight", "bar12-edge", "bar12-straight",
		"bar13-edge", "bar-spike-straight",
		"box1-corner", "box1-edge", "box2-corner", "box2-edge", "box3-corner", "box3-edge", "box4-corner",
		"box4-edge", "box5-corner", "box5-edge", "box6-corner", "box6-edge", "box7-corner", "box7-edge",
		"box8-corner", "box8-edge",
		"box9-corner", "box9-edge", "box10-corner", "box10-edge", "box11-corner", "box11-edge",
		"box12-corner", "box12-edge", "box13-corner", "box13-edge",
		"home", "settings", "save-load", "info", "wait", "file", "folder", "trash", "lock", "key", "pin", "mark",
		"globe", "talk", "letter", "bell", "calendar", "signal-low", "signal-high", "person", "people", "trophy",
		"star1", "star2", "eye", "eye-closed", "bright", "sun", "moon1", "moon2", "stars", "grid", "shut-down",
		"book", "cloud-rain", "cloud", "flag1", "flag2", "pick", "camera-movie", "camera-portable", "microphone",
		"door", "pen", "banner1", "banner2", "filter", "zoom", "stack1", "stack2", "loading1", "loading2", "picture",
		"zap", "mouse", "keyboard", "controller", "check", "x", "cancel", "icon+", "", "back", "loop",
		"reduce", "increase", "list", "grid", "align1", "align2", "align3", "align4", "align5", "align6", "mirror",
		"flip", "bucket", "palette", "previous", "backtrack", "reverse", "play", "forward", "skip", "pause", "record",
		"stop", "mute", "low", "high", "quarter", "seight", "beamed-eight", "beamed-sixteenth", "flat", "natural",
		"sharp", "arrow1", "arrow2", "arrow3", "arrow4-diagonal", "arrow5", "arrow6", "arrow7-diagonal",
		"mountain", "water", "wind", "tree", "pine", "flower", "fish", "animal",
		"time1", "time2", "time3", "time4", "time5", "time6", "time7", "time8",
		"barrier1", "barrier2", "arrow8", "arrow9", "arrow10", "arrow11-diagonal", "arrow12-diagonal",
		"arrow13-diagonal", "arrow14", "arrow15-diagonal", "dice1", "dice2", "dice3", "dice4", "dice5", "dice6",
		"spade1", "heart1", "club1", "diamond1", "spade2", "heart2", "club2", "diamond2",
		"pawn1", "rook1", "knight1", "bishop1", "queen1", "king1",
		"pawn2", "rook2", "knight2", "bishop2", "queen2", "king2",
		"face-smiling", "face-laughing", "face-sad", "face-scared", "face-angry", "face-no-emotion", "face-bored",
		"face-happy", "face-in-love", "face-relieved", "face-unhappy", "face-ego", "face-annoyed", "face-surprised",
		"face-sleepy", "face-kissing", "face-aww", "face-wholesome", "face-crying", "face-tantrum", "face-interested",
		"face-evil", "face-winking", "face-confident", "face-suspicious", "face-mustache",
		"square1", "square2", "square3", "square4", "square5", "square6",
		"circle1", "circle2", "circle3", "circle4", "circle5", "circle6",
		"triangle1", "triangle2", "triangle3", "triangle4",
		"pointer", "wait", "input", "hand", "resize1", "resize2", "resize3", "resize4", "move", "crosshair",
	}

	for i := range ids {
		if ids[i] != "" {
			ids[i] = defaultRetroAtlas + ids[i]
		}
	}
	var tiles = SetTextureAtlasTiles(id, 0, 0, ids...)

	return id, tiles
}
func LoadDefaultAtlasPatterns() (atlasId string, tileIds []string) {
	var tex = loadTexture(defaultPatterns, patterns)
	var id = SetTextureAtlas(tex, 64, 64, 1)
	var ids = []string{}

	for i := range 84 {
		ids = append(ids, defaultPatterns+text.New(i))
	}
	var tiles = SetTextureAtlasTiles(id, 0, 0, ids...)
	return id, tiles
}
func LoadDefaultAtlasInput() (atlasId string, tileIds []string) {
	var tex = loadTexture(defaultInputLeft+defaultInputRight, input)
	var id = SetTextureAtlas(tex, 50, 50, 0)
	var ids = []string{
		"escape", "f1", "f2", "f3", "f4", "f5", "f6", "f7", "f8", "f9", "f10", "f11", "f12", "print", "pause",
		"num-lock", "*",
		"`", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-", "=", "backspace", "insert", "home", "page-up",
		"tab", "q", "w", "e", "r", "t", "y", "u", "i", "o", "p", "[", "]", "enter", "delete", "end", "page-down",
		"caps-lock", "a", "s", "d", "f", "g", "h", "j", "k", "l", ";", "'", "\\", "any", "key+key", "up", "+",
		"shift", "z", "x", "c", "v", "b", "n", "m", ",", ".", "/", "space", "?", "!", "left", "down", "right",
		"control", "alt", "arrows", "arrows1", "arrows2", "arrows3", "arrows4", "keyboard", "mouse", "mouse1",
		"mouse2", "mouse3", "scroll1", "scroll2", "mouse-y", "mouse-x", "mouse-xy",
	}

	for i := range ids {
		if ids[i] != "" {
			ids[i] = defaultInputLeft + ids[i] + defaultInputRight
		}
	}
	var tiles = SetTextureAtlasTiles(id, 0, 0, ids...)
	return id, tiles
}
func LoadDefaultAtlasIcons() (atlasId string, tileIds []string) {
	var tex = loadTexture(defaultIcons, icons)
	var id = SetTextureAtlas(tex, 50, 50, 0)
	var ids = []string{
		"club", "spade", "diamond1", "heart", "heart-broken", "shape",
		"pawn", "knight", "bishop", "rook", "queen", "king", "crown1", "crown2",
		"tower1", "tower2", "tower3", "tower4", "house1", "house2", "house3", "house3",
		"dice0", "dice1", "dice2", "dice3", "dice4", "dice5", "dice6", "dice3D", "dice-hand1", "dice-hand2",
		"dice-arrow", "dices1", "dices2", "dices3", "dices4", "", "", "", "", "bag1", "bag2", "bag3",
		"triangle1", "triangle2", "square1", "square2", "diamond2", "diamond3", "diamond4", "diamond5",
		"pentagon1", "pentagon2", "hexagon1", "hexagon2", "octagon1", "octagon2", "pentagon3", "pentagon4",
		"pentagon5", "pentagon6", "", "", "basket", "cart", "card-hand",
		"card1", "card2", "card3", "card4", "card5", "card6", "card7", "card8", "card9", "card10", "card11",
		"card12", "card13", "card14", "card15", "card16", "card17", "card18", "card19", "card20", "card21", "card22",
		"card23", "card24", "card25", "card26", "card27", "card28", "card29", "card30", "cards1", "cards2",
		"cards3", "cards4", "cards5", "cards6", "cards7", "cards8", "cards9", "cards10", "cards11", "cards12",
		"cards13", "hand", "hand-x", "hand-coin1", "hand-coin2", "coin1", "coin2", "coin3", "coin4", "coin5", "coin6",
		"coin7", "coin8", "coin9", "coin10", "coin11", "coin12", "coins1", "coins2", "coins3", "", "", "shut-down",
		"person1", "people1", "people2", "people3", "person2", "person3", "person4", "person5", "person6", "person7",
		"person8", "person9", "people4", "people5", "people6", "people7", "people8", "rank1", "rank2", "trophy",
		"star", "puzzle", "sword", "shield", "bow", "skull", "flag1", "flag2", "flag3", "bonfire", "fire", "apple",
		"brick", "ingot", "log", "planks", "crops", "victory", "coin", "diamond", "cloud1", "book1", "book2", "pc",
		"gear", "wrench", "save", "movie1", "movie2", "connect", "door", "pen", "eraser", "brush", "bucket", "page",
		"+", "-", "!", "?", "i", "!2", "audio1", "audio2", "audio3", "audio4",
		"clock1", "clock2", "clock3", "clock4", "clock5", "hourglass1", "hourglass2", "potion1", "potion2", "potion3",
		"pause", "stop", "skip", "forward", "play", "left-right1", "arrow1", "arrow2", "arrow3", "pointer",
		"pointer-hand", "cancel", "east", "north", "south", "west", "export", "import", "exit1", "exit2", "arrow4",
		"arrow5", "resize1", "resize2", "increase", "decrease", "move1", "move2", "move3", "move4", "left-right2",
		"refresh1", "refresh2", "loop", "grid", "list", "menu1", "menu2", "signal1", "signal2", "signal3", "accept",
		"x", "trash1", "trash2", "lock1", "lock2", "lock3", "lock4", "key", "zoom", "zoom=", "zoom+", "zoom-",
		"ribbon1", "ribbon2", "bubble1", "bubble2", "bubble3", "bubble4", "stars", "cloud2", "spiral", "water-drop",
		"water-drops", "emotion1", "emotion2", "emotion3", "zz", "haha", "bulb", "hearts", "angry", "happy", "sad",
		".", "..", "...",
	}

	for i := range ids {
		if ids[i] != "" {
			ids[i] = defaultIcons + ids[i]
		}
	}
	var tiles = SetTextureAtlasTiles(id, 0, 0, ids...)
	return id, tiles
}
func LoadDefaultSoundsUI() []string {
	var soundIds = []string{
		loadSound(defaultSoundsUI+"press", press),
		loadSound(defaultSoundsUI+"release", release),
		loadSound(defaultSoundsUI+"on", on),
		loadSound(defaultSoundsUI+"off", off),
		loadSound(defaultSoundsUI+"write", write),
		loadSound(defaultSoundsUI+"erase", erase),
		loadSound(defaultSoundsUI+"slider", slider),
		loadSound(defaultSoundsUI+"popup", popup),
		loadSound(defaultSoundsUI+"error", oops),
	}

	return soundIds
}

//=================================================================
// private

//go:embed default/texture.png
var texture []byte

//go:embed default/cursors.png
var cursors []byte

//go:embed default/retro.png
var retro []byte

//go:embed default/patterns.png
var patterns []byte

//go:embed default/icons.png
var icons []byte

//go:embed default/input.png
var input []byte

//go:embed default/ui.png
var ui []byte

//go:embed default/font.ttf.gz
var font []byte

//go:embed default/press.ogg.gz
var press []byte

//go:embed default/release.ogg.gz
var release []byte

//go:embed default/on.ogg.gz
var on []byte

//go:embed default/off.ogg.gz
var off []byte

//go:embed default/type.ogg.gz
var write []byte

//go:embed default/erase.ogg.gz
var erase []byte

//go:embed default/slider.ogg.gz
var slider []byte

//go:embed default/popup.ogg.gz
var popup []byte

//go:embed default/error.ogg.gz
var oops []byte

func loadTexture(id string, bytes []byte) string {
	tryCreateWindow()

	var _, has = internal.Textures[id]
	if has {
		UnloadTexture(id)
	}

	var image = rl.LoadImageFromMemory(".png", bytes, int32(len(bytes)))
	var tex = rl.LoadTextureFromImage(image)
	internal.Textures[id] = &tex
	rl.UnloadImage(image)
	return id
}
func loadSound(id string, bytes []byte) string {
	tryCreateWindow()
	tryInitAudio()

	var _, has = internal.Sounds[id]
	if has {
		UnloadSound(id)
	}

	var decompressed = storage.DecompressGZIP(bytes)
	var wave = rl.LoadWaveFromMemory(".ogg", decompressed, int32(len(decompressed)))
	var sound = rl.LoadSoundFromWave(wave)
	internal.Sounds[id] = &sound
	rl.UnloadWave(wave)
	return id
}
