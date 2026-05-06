#version 330

layout(location = 0) in vec3 vertPosition;
layout(location = 1) in vec2 vertTexCoord;
layout(location = 2) in vec3 vertNormal;
layout(location = 3) in vec4 vertColor;
layout(location = 4) in vec4 vertTangent;
layout(location = 5) in vec2 vertTexCoord2;

uniform mat4 mvp;

out vec2 fragTexCoord;
out vec4 fragColor;
out vec3 fragNormal;
out vec4 fragTangent;
out vec2 fragTexCoord2;

void main() {
    fragTexCoord = vertTexCoord;
    fragColor = vertColor;
    fragNormal = vertNormal;
    fragTangent = vertTangent;
    fragTexCoord2 = vertTexCoord2;
    gl_Position = mvp * vec4(vertPosition, 1.0);
}
