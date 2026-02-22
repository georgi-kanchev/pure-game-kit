#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

uniform sampler2D texture0;

out vec4 finalColor;

vec4 unpackRGB222(uint val) {
    float r = float((val >> 6u) & 0x03u) / 3.0;
    float g = float((val >> 4u) & 0x03u) / 3.0;
    float b = float((val >> 2u) & 0x03u) / 3.0;
    float a = float(val & 0x03u) / 3.0;
    return vec4(r, g, b, a);
}

void main() {
    uvec4 c = uvec4(fragColor * 255.0 + 0.5);
    vec4 base    = unpackRGB222(c.r);
    vec4 outlineColor = unpackRGB222(c.g);
    vec4 shadowColor = unpackRGB222(c.b);
    
    uint thickIdx  = (c.a >> 6) & 0x03u;
    uint outlIdx   = (c.a >> 4) & 0x03u;
    uint shadIdx   = (c.a >> 2) & 0x03u;
    uint smoothIdx = (c.a)      & 0x03u;
    
    float thicknesses[4] = float[](0.35, 0.50, 0.65, 0.80);
    float shadThicks[4]  = float[](0.20, 0.30, 0.40, 0.50);
    float smooths[4]     = float[](0.50, 4.00, 8.00, 12.0);

    vec2 shadowOffset = vec2(0.002, 0.005);
    float shadowDistance = texture(texture0, fragTexCoord - shadowOffset).a - (1.0 - shadThicks[shadIdx]);
    float shadowSmooth = smooths[smoothIdx] * length(vec2(dFdx(shadowDistance), dFdy(shadowDistance)));
    float shadowAlpha = shadowColor.a * smoothstep(-shadowSmooth, shadowSmooth, shadowDistance);
    
    float distance = texture(texture0, fragTexCoord).a - (1.0 - thicknesses[thickIdx]);
    float baseSmooth = 0.5 * length(vec2(dFdx(distance), dFdy(distance)));
    float sdfAlpha = base.a * smoothstep(-baseSmooth, baseSmooth, distance);
    
    float compressedOutlIdx = (outlIdx * 0.7) + (1.5 * 0.3);
    float outlineThick = (1.0 - thicknesses[thickIdx]) * (compressedOutlIdx / 3.0);
    float outlineAlpha = outlineColor.a * smoothstep(-baseSmooth, baseSmooth, distance + outlineThick);
    
    vec3 finalRGB = shadowColor.rgb;
    finalRGB = mix(finalRGB, outlineColor.rgb, outlineAlpha);
    finalRGB = mix(finalRGB, base.rgb, sdfAlpha);

    float finalAlpha = max(shadowAlpha, max(outlineAlpha, sdfAlpha));

    if (distance > sdfAlpha) {
        finalRGB = base.rgb;
        finalAlpha = base.a;
    }
    
    finalColor = vec4(finalRGB, finalAlpha);
}