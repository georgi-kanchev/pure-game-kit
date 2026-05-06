#version 330

layout(location = 0) in vec3 vertPosition;
layout(location = 1) in vec2 vertTexCoord;
layout(location = 3) in vec4 vertColor;
layout(location = 7) in vec4 vertCustom;

uniform mat4 mvp;

out vec2 fragTexCoord;
out vec4 fragColor;
out vec4 fragCustom;

void main() {
    fragTexCoord = vertTexCoord;
    fragColor = vertColor;
    gl_Position = mvp * vec4(vertPosition, 1.0);
}
