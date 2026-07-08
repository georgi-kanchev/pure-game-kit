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
out vec2 fragLocalPos;

out vec4 fragData0;
out vec4 fragData1;
out vec4 fragData2;
out vec4 fragData3;
out vec4 fragData4;
out vec4 fragData5;
out vec4 fragData6;
out vec4 fragData7;

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
    int rawB   = int((bits >> 2u) & 0x7FFu);
    borderSize = float(rawB >= 1024 ? rawB - 2048 : rawB);
    objType    = int(bits & 0x3u);
}
void unpack_16_8(float packedFloat, out float outlineSize, out float tileSize) {
    uint bits = floatBitsToUint(packedFloat);
    outlineSize = float((bits >> 8u) & 0xFFFFu);
    tileSize    = float(bits & 0xFFu);
}
void unpack_8_8_8(float packedFloat, out float shadowX, out float shadowY, out float shadowBlur) {
    uint bits = floatBitsToUint(packedFloat);
    int rawX = int((bits >> 16u) & 0xFFu);
    int rawY = int((bits >> 8u) & 0xFFu);
    shadowX = -float(rawX >= 128 ? rawX - 256 : rawX) / 32.0;
    shadowY = -float(rawY >= 128 ? rawY - 256 : rawY) / 32.0;
    shadowBlur = float(bits & 0xFFu);
}
vec4 unpack_10_4_5_5(float packedFloat) {
    uint bits = floatBitsToUint(packedFloat);
    float roundness = float((bits >> 14u) & 0x3FFu) / 1023.0;
    float pixelSize = float((bits >> 10u) & 0xFu);
    float blurX     = float((bits >> 5u)  & 0x1Fu) / 31.0;
    float blurY     = float(bits & 0x1Fu) / 31.0;
    return vec4(roundness, pixelSize, blurX, blurY);
}
void unpack_8_4_4_4_4(float packedFloat, out float outlineSize, out vec4 fillColor) {
    uint bits = floatBitsToUint(packedFloat);
    outlineSize = float((bits >> 16u) & 0xFFu);
    fillColor.x = float((bits >> 12u) & 0xFu) / 15.0;
    fillColor.y = float((bits >> 8u)  & 0xFu) / 15.0;
    fillColor.z = float((bits >> 4u)  & 0xFu) / 15.0;
    fillColor.w = float(bits & 0xFu) / 15.0;
}
vec3 unpack_8_8_8_text_weights(float packedFloat) {
    uint bits = floatBitsToUint(packedFloat);
    int rawX = int((bits >> 16u) & 0xFFu);
    int rawZ = int(bits & 0xFFu);
    float x = float(rawX >= 128 ? rawX - 256 : rawX); // signed: weight
    float y = float((bits >> 8u) & 0xFFu);             // unsigned: outlineWeight
    float z = float(rawZ >= 128 ? rawZ - 256 : rawZ); // signed: shadowWeight
    return vec3(x, y, z);
}

void main() {
    fragTexCoord = vertTexCoord;
    fragColor = vertColor;
    
    vec2 texSize = unpack_12_12(vertTexCoord2.x);
    vec4 borderColor = unpack_6_6_6_6(vertTexCoord2.y);
    vec4 colorAdjust1 = unpack_6_6_6_6(vertNormal.x);
    vec4 rgbAdjust2 = unpack_10_4_5_5(vertNormal.y);
    float roundness = rgbAdjust2.x;
    float pixelSize = rgbAdjust2.y;
    float blurX = rgbAdjust2.z;
    float blurY = rgbAdjust2.w;
    float depthZ, borderSize;
    int objectType;
    unpack_11_11_2(vertNormal.z, depthZ, borderSize, objectType);
    float outlineSize = 0.0;
    vec4 outlineColor   = vec4(0.0);
    vec4 fillColor = vec4(0.0);
    float tileColumns = 0.0;
    float tileRows    = 0.0;
    float tileSize    = 0.0;
    vec4 shadowColor_text = vec4(0.0);
    vec3 textWeights      = vec3(0.0);
    float shadowX = 0.0, shadowY = 0.0, shadowBlur = 0.0;
    vec2 cropBoundsU = vec2(0.0);
    vec2 cropBoundsV = vec2(0.0);
    
    if (objectType == 0) { // Shape
        cropBoundsU = unpack_12_12(vertTangent.z);
        cropBoundsV = unpack_12_12(vertTangent.w);
    }
    else if (objectType == 1) { // Sprite
        outlineColor = unpack_6_6_6_6(vertTangent.x);
        unpack_8_4_4_4_4(vertTangent.y, outlineSize, fillColor);
        cropBoundsU = unpack_12_12(vertTangent.z);
        cropBoundsV = unpack_12_12(vertTangent.w);
    }
    else if (objectType == 2) { // Text
        outlineColor     = unpack_6_6_6_6(vertTangent.x);
        shadowColor_text = unpack_6_6_6_6(vertTangent.y);
        textWeights      = unpack_8_8_8_text_weights(vertTangent.z);
        unpack_8_8_8(vertTangent.w, shadowX, shadowY, shadowBlur);
    }
    else if (objectType == 3) { // Tilemap
        outlineColor    = unpack_6_6_6_6(vertTangent.x);
        fillColor = unpack_6_6_6_6(vertTangent.y);
        vec2 tileInfo   = unpack_12_12(vertTangent.z);
        tileColumns = tileInfo.x;
        tileRows    = tileInfo.y;
        unpack_16_8(vertTangent.w, outlineSize, tileSize);
    }
    
    fragData0 = vec4(texSize, depthZ, float(objectType));
    fragData1 = colorAdjust1;
    fragData2 = vec4(roundness, pixelSize, blurX, blurY);
    fragData3 = outlineColor;
    fragData7 = borderColor;

    // Pack a flag into fragData5.w: 1.0 = color adjustments are neutral (skip expensive pow/exp2)
    bool neutralCA = (floatBitsToUint(vertNormal.x) == 0x3F8800A0u);

    if (objectType == 2) { // Text
        fragData4 = shadowColor_text;
        fragData5 = vec4(textWeights.x / 127.0, textWeights.y / 255.0, textWeights.z / 127.0, neutralCA ? 1.0 : 0.0);
        fragData6 = vec4(shadowX, shadowY, shadowBlur, 0.0);
    } else if (objectType == 0 || objectType == 1) { // Shape or Sprite
        fragData4 = fillColor;
        fragData5 = vec4(outlineSize, borderSize, 0.0, neutralCA ? 1.0 : 0.0);
        fragData6 = vec4(cropBoundsU, cropBoundsV);
    } else { // Tilemap
        fragData4 = fillColor;
        fragData5 = vec4(outlineSize, borderSize, 0.0, neutralCA ? 1.0 : 0.0);
        fragData6 = vec4(tileColumns, tileRows, tileSize, 0.0);
    }

    fragLocalPos = vertTexCoord - 0.5;

    gl_Position = mvp * vec4(vertPosition, 1.0);
}
