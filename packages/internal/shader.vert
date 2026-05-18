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

out vec4 fragData0; // texSize.xy + depthZ + objectType
out vec4 fragData1; // colorAdjust1 (gamma, saturation, contrast, brightness)
out vec4 fragData2; // colorAdjust2 (grayscale, inversion, blurX, blurY)
out vec4 fragData3; // outlineColor RGBA
out vec4 fragData4; // silhouetteColor RGBA
out vec4 fragData5; // outlineSize + roundness + pixelSize
out vec4 fragData6; // tileColumns + tileRows + tileSize

vec4 unpack_6_6_6_6(float packedFloat) {
    uint bits = floatBitsToUint(packedFloat);
    float x = float((bits >> 18u) & 0x3Fu) / 63.0;
    float y = float((bits >> 12u) & 0x3Fu) / 63.0;
    float z = float((bits >> 6u)  & 0x3Fu) / 63.0;
    float w = float(bits & 0x3Fu) / 63.0;
    return vec4(x, y, z, w);
}
vec2 unpack_12_12(float packedFloat) {
    uint bits = floatBitsToUint(packedFloat);
    float w = float((bits >> 12u) & 0xFFFu);
    float h = float(bits & 0xFFFu);
    return vec2(w, h);
}
void unpack_11_11_2(float packedFloat, out float depthZ, out float borderSize, out int objType) {
    uint bits = floatBitsToUint(packedFloat);
    depthZ     = float((bits >> 13u) & 0x7FFu) / 2047.0;
    borderSize = float((bits >> 2u)  & 0x7FFu);
    objType    = int(bits & 0x3u);
}
void unpack_16_8(float packedFloat, out float roundness, out float pixelSize) {
    uint bits = floatBitsToUint(packedFloat);
    roundness = float((bits >> 8u) & 0xFFFFu) / 65535.0;
    pixelSize = float(bits & 0xFFu);
}
void unpack_8_8_8(float packedFloat, out float outlineSize, out float tileSize, out float roundness) {
    uint bits = floatBitsToUint(packedFloat);
    outlineSize = float((bits >> 16u) & 0xFFu);
    tileSize    = float((bits >> 8u)  & 0xFFu);
    roundness   = float(bits & 0xFFu) / 255.0;
}
vec3 unpack_10_10_4(float packedFloat) {
    uint bits = floatBitsToUint(packedFloat);
    float cols = float((bits >> 14u) & 0x3FFu);
    float rows = float((bits >> 4u)  & 0x3FFu);
    float ps   = float(bits & 0xFu);
    return vec3(cols, rows, ps);
}
void unpack_6_6_8_4(float packedFloat, out float shadowX, out float shadowY, out float roundness, out float pixelSize) {
    uint bits = floatBitsToUint(packedFloat);
    int rawX = int((bits >> 18u) & 0x3Fu);
    int rawY = int((bits >> 12u) & 0x3Fu);
    shadowX = float(rawX >= 32 ? rawX - 64 : rawX) / 32.0;
    shadowY = float(rawY >= 32 ? rawY - 64 : rawY) / 32.0;
    roundness = float((bits >> 4u) & 0xFFu) / 255.0;
    pixelSize = float(bits & 0xFu);
}

void main() {
    fragTexCoord = vertTexCoord;
    fragColor = vertColor;
    
    vec2 texSize = unpack_12_12(vertTexCoord2.x);
    vec4 colorAdjust1 = unpack_6_6_6_6(vertNormal.x);
    vec4 colorAdjust2 = unpack_6_6_6_6(vertNormal.y);
    float depthZ, borderSize;
    int   objectType;
    unpack_11_11_2(vertNormal.z, depthZ, borderSize, objectType);
    float roundness  = 0.0;
    float pixelSize  = 0.0;
    float outlineSize = 0.0;
    vec4  outlineColor   = vec4(0.0);
    vec4  silhouetteColor = vec4(0.0);
    float tileColumns = 0.0;
    float tileRows    = 0.0;
    float tileSize    = 0.0;

    if (objectType == 0) { // Shape
        roundness = vertTangent.x;
        pixelSize = vertTangent.y;
    }
    else if (objectType == 1) { // Sprite
        outlineColor    = unpack_6_6_6_6(vertTangent.x);
        silhouetteColor = unpack_6_6_6_6(vertTangent.y);
        outlineSize     = vertTangent.z;
        unpack_16_8(vertTangent.w, roundness, pixelSize);
    }
    else if (objectType == 2) { // Text
        outlineColor = unpack_6_6_6_6(vertTangent.x);
        vec4 shadowColor = unpack_6_6_6_6(vertTangent.y);
        vec4 textWeights = unpack_6_6_6_6(vertTangent.z);
        float shadowX, shadowY;
        unpack_6_6_8_4(vertTangent.w, shadowX, shadowY, roundness, pixelSize);
    }
    else if (objectType == 3) { // Tilemap
        outlineColor    = unpack_6_6_6_6(vertTangent.x);
        silhouetteColor = unpack_6_6_6_6(vertTangent.y);
        vec3 tileInfo   = unpack_10_10_4(vertTangent.z);
        tileColumns = tileInfo.x;
        tileRows    = tileInfo.y;
        pixelSize   = tileInfo.z;
        unpack_8_8_8(vertTangent.w, outlineSize, tileSize, roundness);
    }
    
    fragData0 = vec4(texSize, depthZ, float(objectType));
    fragData1 = colorAdjust1;
    fragData2 = colorAdjust2;
    fragData3 = outlineColor;
    fragData4 = silhouetteColor;
    fragData5 = vec4(outlineSize, roundness, pixelSize, 0.0);
    fragData6 = vec4(tileColumns, tileRows, tileSize, 0.0);

    gl_Position = mvp * vec4(vertPosition, 1.0);
}
