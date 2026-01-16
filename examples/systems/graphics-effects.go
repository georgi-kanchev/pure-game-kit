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

	var spr = graphics.NewSprite(assets.LoadDefaultTexture(), 0, 0)

	for window.KeepOpen() {
		cam.SetScreenAreaToWindow()
		cam.MouseDragAndZoomSmoothly()

		var w, h = assets.Size("")
		rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "renderSize"), []float32{float32(w), float32(h)}, rl.ShaderUniformVec2)
		rl.SetShaderValue(sh, rl.GetShaderLocation(sh, "strength"), []float32{2}, rl.ShaderUniformFloat)

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
uniform vec4 colDiffuse;

// Replaces 'tileCount * tileSize' (The full size of the texture/screen)
uniform vec2 renderSize; 
// Replaces 's' (Blur spread/scale)
uniform vec2 blurScale; 

void main()
{
    // Calculate the size of one pixel in UV space (0.0 to 1.0)
    vec2 texel = 1.0 / renderSize;
    
    // Calculate the offset vector
    vec2 blur = texel * blurScale;

    // --- The 3x3 Kernel Math from your snippet ---
    
    // Center Pixel (Weight: 4)
    vec4 sum = texture(texture0, fragTexCoord) * 4.0;

    // Horizontal Neighbors (Weight: 2)
    // Interpret 'coord - blur.x' as moving Left
    sum += texture(texture0, fragTexCoord + vec2(-blur.x, 0.0)) * 2.0;
    sum += texture(texture0, fragTexCoord + vec2( blur.x, 0.0)) * 2.0;

    // Vertical Neighbors (Weight: 2)
    // Interpret 'coord - blur.y' as moving Up
    sum += texture(texture0, fragTexCoord + vec2(0.0, -blur.y)) * 2.0;
    sum += texture(texture0, fragTexCoord + vec2(0.0,  blur.y)) * 2.0;

    // Diagonal Neighbors (Weight: 1)
    // Top-Left, Top-Right, Bottom-Left, Bottom-Right
    sum += texture(texture0, fragTexCoord + vec2(-blur.x, -blur.y)) * 1.0;
    sum += texture(texture0, fragTexCoord + vec2(-blur.x,  blur.y)) * 1.0;
    sum += texture(texture0, fragTexCoord + vec2( blur.x, -blur.y)) * 1.0;
    sum += texture(texture0, fragTexCoord + vec2( blur.x,  blur.y)) * 1.0;

    // Normalize: Total weight is 4 + 2*4 + 1*4 = 16
    vec4 blurResult = sum / 16.0;

    finalColor = blurResult * colDiffuse * fragColor;
}`

const VERT = `#version 330

// Input vertex attributes
in vec3 vertexPosition;
in vec2 vertexTexCoord;
in vec4 vertexColor;

// Input uniform values
uniform mat4 mvp; // Model-View-Projection matrix

// Output vertex attributes (to fragment shader)
out vec2 fragTexCoord;
out vec4 fragColor;

void main()
{
    // Send the texture coordinates and color to the fragment shader
    fragTexCoord = vertexTexCoord;
    fragColor = vertexColor;

    // Calculate the final vertex position on screen
    gl_Position = mvp * vec4(vertexPosition, 1.0);
}`
