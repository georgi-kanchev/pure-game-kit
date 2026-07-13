#version 330

layout(location = 0) in vec3 vertPosition;
layout(location = 1) in vec2 vertTexCoord;
layout(location = 3) in vec4 vertColor;

uniform mat4 mvp;

out vec2 fragTexCoord;
out vec4 fragColor;
out vec2 fragLocalPos;

void main() {
    fragTexCoord = vertTexCoord;
    fragColor = vertColor;
    fragLocalPos = vertTexCoord - 0.5;
    gl_Position = mvp * vec4(vertPosition, 1.0);
}
