package example

import (
	"pure-game-kit/data/assets"
	"pure-game-kit/debug"
	"pure-game-kit/graphics"
	"pure-game-kit/utility/time"
	"pure-game-kit/window"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func Effects() {
	var cam = graphics.NewCamera(4)
	var sh = rl.LoadShaderFromMemory(VERT, FRAG)
	defer rl.UnloadShader(sh)

	var tex = assets.LoadTexture("examples/data/objects.png")
	var spr = graphics.NewSprite(tex, 0, 0)
	assets.SetTextureSmoothness(tex, true)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmoothly()

		var w, h = assets.Size("")
		rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "texSize"), []float32{float32(w), float32(h)}, rl.ShaderUniformVec2)

		rl.BeginShaderMode(sh)
		cam.DrawSprites(spr)
		rl.EndShaderMode()

		debug.Print(time.FrameRate())
	}
}

const FRAG = `#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

out vec4 finalColor;

uniform sampler2D texture0;
uniform vec2 texSize; 

//=================================================================

uniform vec2 blur;

uniform float gam = 0.5;
uniform float sat = 0.5;
uniform float con = 0.5;
uniform float bri = 0.5;
uniform float gra = 0.0;
uniform float inv = 0.0;

uniform vec4 overlay;

//=================================================================

float map(float value, float min1, float max1, float min2, float max2) {
    return min2 + (value - min1) * (max2 - min2) / (max1 - min1);
}

//=================================================================

vec4 compute_color_adjust(vec4 color) {
	float luminance = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
	float gamma = gam < 0.5 ? map(gam, 0.0, 0.5, 6.0, 1.0) : map(gam, 0.5, 1.0, 1.0, 0.0);
	float saturation = sat < 0.5 ? map(sat, 0.0, 0.5, 0.0, 1.0) : map(sat, 0.5, 1.0, 1.0, 10.0);
	float contrast = con < 0.5 ? map(con, 0.0, 0.5, 0.0, 1.0) : map(con, 0.5, 1.0, 1.0, 3.0);
	float brightness = bri < 0.5 ? map(bri, 0.0, 0.5, 0.0, 1.0) : map(bri, 0.5, 1.0, 1.0, 4.0);

	color.rgb = pow(color.rgb, vec3(gamma));
	color.rgb = mix(vec3(luminance), color.rgb, saturation);
	color.rgb = mix(vec3(0.5), color.rgb, contrast);
	color.rgb = mix(color.rgb, vec3(luminance), gra);
	color.rgb = mix(color.rgb, 1.0 - color.rgb, inv);
	color.rgb *= brightness;
	return color;
}
vec4 compute_blur(vec4 color) {
	vec2 texel = 1.0 / texSize;
    vec2 blurValue = texel * blur;
    color *= 4.0;
    color += texture(texture0, fragTexCoord + vec2(-blurValue.x, 0.0)) * 2.0;
    color += texture(texture0, fragTexCoord + vec2( blurValue.x, 0.0)) * 2.0;
    color += texture(texture0, fragTexCoord + vec2(0.0, -blurValue.y)) * 2.0;
    color += texture(texture0, fragTexCoord + vec2(0.0,  blurValue.y)) * 2.0;
    color += texture(texture0, fragTexCoord + vec2(-blurValue.x, -blurValue.y)) * 1.0;
    color += texture(texture0, fragTexCoord + vec2(-blurValue.x,  blurValue.y)) * 1.0;
    color += texture(texture0, fragTexCoord + vec2( blurValue.x, -blurValue.y)) * 1.0;
    color += texture(texture0, fragTexCoord + vec2( blurValue.x,  blurValue.y)) * 1.0;
	return color / 16.0;
}
vec4 compute_overlay(vec4 color) {
	return vec4(mix(color.rgb, overlay.rgb, overlay.a), color.a);
}
void main()
{
    vec4 color = texture(texture0, fragTexCoord);
	if (blur.x > 0 || blur.y > 0)
		color = compute_blur(color);
	if (gam != 0.5 || sat != 0.5 || con != 0.5 || bri != 0.5 || gra != 0.0 || inv != 0.0)
		color = compute_color_adjust(color);
	if (overlay.a > 0)
		color = compute_overlay(color);
    finalColor = color * fragColor;
}`

const VERT = `#version 330

in vec3 vertexPosition;
in vec2 vertexTexCoord;
in vec4 vertexColor;

uniform mat4 mvp;

out vec2 fragTexCoord;
out vec4 fragColor;

void main()
{
    fragTexCoord = vertexTexCoord;
    fragColor = vertexColor;
    gl_Position = mvp * vec4(vertexPosition, 1.0);
}`
