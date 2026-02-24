#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

uniform sampler2D texture0;
uniform sampler2D dataTex;

out vec4 finalColor;

void main() {
	vec4 col = texture(texture0, fragTexCoord);
	vec4 data = texture(dataTex, fragTexCoord);
	if (data.a > 0)
		col = vec4(0, 0, 1, 1);
	finalColor = col;
}