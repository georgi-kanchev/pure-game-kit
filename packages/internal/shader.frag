#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

in vec2 fragTexCoord2;
in vec3 fragNormal;
in vec4 fragTangent;

// texCoord2.x = TextureWidth + TextureHeight + Type | 12+12+8 = 32 bits
// texCoord2.y = BorderColor + Roundness             | 24+8    = 32 bits

// normal.x = Gamma + Saturation + Contrast + Brightness | 8+8+8+8 = 32 bits
// normal.y = Grayscale + Inversion + BlurX + BlurY      | 8+8+8+8 = 32 bits
// normal.z = DepthZ + BorderSize + PixelSize            | 12+12+8 = 32 bits

// Shape:
//  tangent = free

// Sprite:
//  tangent.x = OutlineColor + OutlineSize | 24+8 = 32 bits
//  tangent.y = SilhouetteColor            | 32 bits
//  tangent.z = free
//  tangent.w = free

// Tilemap:
//  tangent.x = OutlineColor + OutlineSize        | 24+8 = 32 bits
//  tangent.y = SilhouetteColor                   | 32 bits
//  tangent.z = TileColumns + TileRows + TileSize | 12+12+8 = 32 bits
//  tangent.w = free

// Text:
//  tangent.x = OutlineColor + TextShadowX                         | 24+8 = 32 bits
//  tangent.y = ShadowColor + TextShadowY                          | 24+8 = 32 bits
//  tangent.z = Weight + OutlineWeight + ShadowWeight + ShadowBlur | 8+8+8+8 = 32 bits
//  tangent.w = free

out vec4 finalColor;

#define OUTLINE_SIZE 12
#define OUTLINE_R 13
#define OUTLINE_G 14
#define OUTLINE_B 15
#define OUTLINE_A 16
#define SILHOUETTE_R 17
#define SILHOUETTE_G 18
#define SILHOUETTE_B 19
#define SILHOUETTE_A 20
#define TILE_COLUMNS 21
#define TILE_ROWS 22
#define TILE_W 23
#define TILE_H 24
#define TIME 25

uniform sampler2D texture0;
uniform sampler2D tileData;
uniform float u[32];

void unpack_24_8(float packedFloat, out vec3 rgb, out float extra8) {
    uint bitString = floatBitsToUint(packedFloat);
    float r = float((bitString >> 24u) & 0xFFu) / 255.0; // 8 bits
    float g = float((bitString >> 16u) & 0xFFu) / 255.0; // 8 bits
    float b = float((bitString >> 8u) & 0xFFu) / 255.0;  // 8 bits
    rgb = vec3(r, g, b);
    extra8 = float(bitString & 0xFFu) / 255.0; // normalized 0.0 to 1.0
}
vec4 unpack_8_8_8_8(float packedFloat) {
    uint bitString = floatBitsToUint(packedFloat);
    float x = float((bitString >> 24u) & 0xFFu) / 255.0;
    float y = float((bitString >> 16u) & 0xFFu) / 255.0;
    float z = float((bitString >> 8u) & 0xFFu) / 255.0;
    float w = float(bitString & 0xFFu) / 255.0;
    return vec4(x, y, z, w); // normalized 0.0 to 1.0
}
vec3 unpack_12_12_8(float packedFloat) {
    uint bitString = floatBitsToUint(packedFloat);
    float val1 = float((bitString >> 20u) & 0xFFFu);
    float val2 = float((bitString >> 8u) & 0xFFFu);
    float val3 = float(bitString & 0xFFu);
    return vec3(val1, val2, val3); // absolute values, not normalized
}
vec4 unpack_RGBA32(float packedFloat) {
    uint bitString = floatBitsToUint(packedFloat);
    float r = float((bitString >> 24u) & 0xFFu) / 255.0;
    float g = float((bitString >> 16u) & 0xFFu) / 255.0;
    float b = float((bitString >> 8u) & 0xFFu) / 255.0;
    float a = float(bitString & 0xFFu) / 255.0;
    return vec4(r, g, b, a);
}

float map(float value, float min1, float max1, float min2, float max2) {
    return min2 + (value - min1) * (max2 - min2) / (max1 - min1);
}

vec2 compute_pixelated_uv(vec2 uv, vec2 texSize, float pixelSize) {
    vec2 numBlocks = texSize / max(pixelSize, 0.001); // avoid division by zero
    vec2 pixelated = (floor(uv * numBlocks) + 0.5) / numBlocks;
    return pixelSize <= 1.0 ? uv : pixelated;
}
vec4 compute_blur(vec2 uv, vec2 texSize, vec2 blur) {
    if (blur.x == 0.0 && blur.y == 0.0)
        return texture(texture0, uv);
    
    vec2 res = 1.0 / texSize;
    vec2 offset = (blur + 0.5) * res;
    vec4 sum = texture(texture0, uv + vec2(-offset.x, -offset.y));
    sum += texture(texture0, uv + vec2(offset.x, -offset.y));
    sum += texture(texture0, uv + vec2(-offset.x, offset.y));
    sum += texture(texture0, uv + vec2(offset.x, offset.y));
    return sum * 0.25;
}
vec4 compute_outline(vec4 color, vec2 uv, vec2 texSize) {
    float outline = u[OUTLINE_SIZE];
    if (color.a > 0 || outline == 0.0)
        return color;
    
    vec2 texel = 1.0 / texSize;

    if (texture(texture0, uv + vec2(texel.x * outline, 0.0)).a > 0.0 ||
            texture(texture0, uv + vec2(-texel.x * outline, 0.0)).a > 0.0 ||
            texture(texture0, uv + vec2(0.0, texel.y * outline)).a > 0.0 ||
            texture(texture0, uv + vec2(0.0, -texel.y * outline)).a > 0.0 ||
            texture(texture0, uv + vec2(texel.x * outline, texel.y * outline)).a > 0.0 ||
            texture(texture0, uv + vec2(-texel.x * outline, texel.y * outline)).a > 0.0 ||
            texture(texture0, uv + vec2(-texel.x * outline, -texel.y * outline)).a > 0.0 ||
            texture(texture0, uv + vec2(texel.x * outline, -texel.y * outline)).a > 0.0)
        return vec4(u[OUTLINE_R], u[OUTLINE_G], u[OUTLINE_B], u[OUTLINE_A]);

    return color;
}
vec4 compute_silhouette(vec4 color) {
    vec4 silColor = vec4(u[SILHOUETTE_R], u[SILHOUETTE_G], u[SILHOUETTE_B], u[SILHOUETTE_A]);
    color.rgb = mix(color.rgb, silColor.rgb, silColor.a);
    return color;
}
vec4 compute_color_adjust(vec4 color, vec4 colorAdjust, vec2 colorAdjust2) {
    float gam = colorAdjust.x;
    float sat = colorAdjust.y;
    float con = colorAdjust.z;
    float bri = colorAdjust.w;
    float luminance = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
    float gamma = gam < 0.5 ? map(gam, 0.0, 0.5, 6.0, 1.0) : map(gam, 0.5, 1.0, 1.0, 0.0);
    // float saturation = sat < 0.5 ? map(sat, 0.0, 0.5, 0.0, 1.0) : map(sat, 0.5, 1.0, 1.0, 3.0);
    float contrast = con < 0.5 ? map(con, 0.0, 0.5, 0.0, 1.0) : map(con, 0.5, 1.0, 1.0, 3.0);
    float brightness = bri < 0.5 ? map(bri, 0.0, 0.5, 0.0, 1.0) : map(bri, 0.5, 1.0, 1.0, 4.0);
    color.rgb = pow(max(color.rgb, vec3(0.0)), vec3(gamma));
    float lum_pre_sat = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
    color.rgb = mix(vec3(lum_pre_sat), color.rgb, sat);
    color.rgb = mix(vec3(0.5), color.rgb, contrast);
    float lum_post_con = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
    color.rgb = mix(color.rgb, vec3(lum_post_con), colorAdjust2.x);
    color.rgb = mix(color.rgb, 1.0 - color.rgb, colorAdjust2.y);
    color.rgb *= brightness;
    return color;
}

vec2 compute_tile(vec2 uv, vec2 texSize) {
    ivec2 mapSize = ivec2(int(u[TILE_COLUMNS]), int(u[TILE_ROWS]));
    if (mapSize.x == 0)
        return uv; // this is a regular sprite, not a tilemap

    ivec2 tile = ivec2(int(uv.x * float(mapSize.x)), int(uv.y * float(mapSize.y)));
    tile = clamp(tile, ivec2(0), mapSize - 1);
    int linearTileID = tile.y * mapSize.x + tile.x;
    ivec2 dataUv = ivec2(linearTileID % mapSize.x, linearTileID / mapSize.x);
    vec4 texColor = texelFetch(tileData, dataUv, 0);
    uvec4 b = uvec4(texColor * 255.0 + 0.5);
    uint gid = (b.r << 24) | (b.g << 16) | (b.b << 8) | b.a;

    bool flip = (gid & 0x80000000u) != 0u; // bits 31..31
    uint rot = (gid & 0x60000000u) >> 29; // bits 30..29
    uint animCount = (gid & 0x1E000000u) >> 25; // bits 28..25
    uint animOffset = (gid & 0x01E00000u) >> 21; // bits 24..21
    uint speedRaw = (gid & 0x001F0000u) >> 16; // bits 20..16
    uint atlasBase = gid & 0xFFFFu; // bits 15..00

    float s = float(speedRaw); // multiplier logic: 0..10 maps to 0.00..1.00; 11..31 maps to 1.33..10.00
    float multiplier = (s <= 10.0) ? (s * 0.1) : (1.0 + (s - 10.0) * 0.45);
    uint frameRange = animCount + 1u;
    uint currentFrame = uint(mod(floor(u[TIME] * multiplier) + float(animOffset), float(frameRange)));
    uint atlasIndex = atlasBase + currentFrame;

    float atlasCols = floor(texSize.x / u[TILE_W]);
    vec2 coord = vec2(mod(float(atlasIndex), atlasCols), floor(float(atlasIndex) / atlasCols));
    vec2 localUV = fract(uv * vec2(float(mapSize.x), float(mapSize.y)));
    localUV -= 0.5;
    localUV.x = flip ? -localUV.x : localUV.x;
    localUV = rot == 1u ? vec2(localUV.y, -localUV.x) : localUV; // 90 degrees
    localUV = rot == 2u ? vec2(-localUV.x, -localUV.y) : localUV; // 180 degrees
    localUV = rot == 3u ? vec2(-localUV.y, localUV.x) : localUV; // 270 degrees
    localUV += 0.5;

    localUV = mix(localUV, vec2(0.5), 1.0 / vec2(u[TILE_W], u[TILE_H])); // prevents texture bleeding artifacts

    vec2 atlasSizeInTiles = vec2(texSize.x / u[TILE_W], texSize.y / u[TILE_H]);
    return (coord + localUV) / atlasSizeInTiles;
}
// vec4 compute_sdf_text(vec2 uv) {
//     uvec4 c = uvec4(fragColor * 255.0 + 0.5);
//     vec4 base = unpackRGB222(c.r);
//     vec4 outlineColor = unpackRGB222(c.g);
//     vec4 shadowColor = unpackRGB222(c.b);
    
//     uint thickIdx = (c.a >> 6) & 0x03u;
//     uint outlIdx = (c.a >> 4) & 0x03u;
//     uint shadIdx = (c.a >> 2) & 0x03u;
//     uint smoothIdx = (c.a) & 0x03u;
    
//     float thick[4] = float[](0.35, 0.50, 0.65, 0.80);
//     float smooths[4] = float[](0.50, 4.00, 8.00, 12.0);
    
//     vec2 shadowOffset = vec2(u[TEXT_SHADOW_X], u[TEXT_SHADOW_Y]);
//     float shadowDistance = texture(texture0, uv - shadowOffset).a - (1.0 - thick[shadIdx]);
//     float shadowSmooth = smooths[smoothIdx] * length(vec2(dFdx(shadowDistance), dFdy(shadowDistance)));
//     float shadowAlpha = shadowColor.a * smoothstep(-shadowSmooth, shadowSmooth, shadowDistance);
    
//     float distance = texture(texture0, uv).a - (1.0 - thick[thickIdx]);
//     float baseSmooth = 0.5 * length(vec2(dFdx(distance), dFdy(distance)));
//     float sdfAlpha = base.a * smoothstep(-baseSmooth, baseSmooth, distance);
    
//     float compressedOutlIdx = map(float(outlIdx), 0.0, 3.0, 0.7, 2.9);
//     float outlineThick = (1.0 - thick[thickIdx]) * (compressedOutlIdx / 3.0);
//     float outlineAlpha = outlineColor.a * smoothstep(-baseSmooth, baseSmooth, distance + outlineThick);
    
//     vec3 mixedRGB = mix(shadowColor.rgb, outlineColor.rgb, outlineAlpha);
//     mixedRGB = mix(mixedRGB, base.rgb, sdfAlpha);
//     float mixedAlpha = max(shadowAlpha, max(outlineAlpha, sdfAlpha));
    
//     vec3 finalRGB = distance > sdfAlpha ? base.rgb : mixedRGB;
//     float finalAlpha = distance > sdfAlpha ? base.a : mixedAlpha;
    
//     return vec4(finalRGB, finalAlpha);
// }

void main() {
    vec3 texInfo = unpack_12_12_8(fragTexCoord2.x); // TextureWidth(12) + TextureHeight(12) + Type(8)
    vec2 texSize = texInfo.xy;
    int objectType = int(texInfo.z); // 0=Shape, 1=Sprite, 2=Text, 3=Tilemap
    
    vec3 borderColor;
    float roundness;
    unpack_24_8(fragTexCoord2.y, borderColor, roundness); // BorderColor(24) + Roundness(8)

    vec4 colorAdjust1 = unpack_8_8_8_8(fragNormal.x); // Gamma(8) + Saturation(8) + Contrast(8) + Brightness(8)
    vec4 colorAdjust2 = unpack_8_8_8_8(fragNormal.y); // Grayscale(8) + Inversion(8) + BlurX(8) + BlurY(8)
    vec2 blur = colorAdjust2.zw * 16;

    vec3 depthAndSizes = unpack_12_12_8(fragNormal.z); // DepthZ(12) + BorderSize(12) + PixelSize(8)
    float depthZ     = depthAndSizes.x / 4095.0; // normalized 0.0 to 1.0
    float borderSize = depthAndSizes.y;
    float pixelSize  = depthAndSizes.z;
    
    if (objectType == 1) { // Sprite
        vec3 outlineColor;
        float outlineSize;
        unpack_24_8(fragTangent.x, outlineColor, outlineSize);
        
        vec4 silhouetteColor = unpack_RGBA32(fragTangent.y);
    }
    else if (objectType == 2) { // Text
        vec3 outlineColor;
        float rawShadowX;
        unpack_24_8(fragTangent.x, outlineColor, rawShadowX);

        vec3 shadowColor;
        float rawShadowY;
        unpack_24_8(fragTangent.y, shadowColor, rawShadowY);

        float signedShadowX = rawShadowX * 2.0 - 1.0; // map from [0.0, 1.0] to [-1.0, 1.0]
        float signedShadowY = rawShadowY * 2.0 - 1.0; // map from [0.0, 1.0] to [-1.0, 1.0]
        
        vec4 textWeights = unpack_8_8_8_8(fragTangent.z); // Weight(8) + OutlineWeight(8) + ShadowWeight(8) + ShadowBlur(8)
    }
    else if (objectType == 3) { // Tilemap
        vec3 outlineColor;
        float outlineSize;
        unpack_24_8(fragTangent.x, outlineColor, outlineSize);
        
        vec4 silhouetteColor = unpack_RGBA32(fragTangent.y);
        
        vec3 tileInfo = unpack_12_12_8(fragTangent.z); // TileColumns(12) + TileRows(12) + TileSize(8)
        float tileColumns = tileInfo.x;
        float tileRows    = tileInfo.y;
        float tileSize    = tileInfo.z;
    }
    
    //========================================================================

    vec2 uv = fragTexCoord;
    // if (u[CALCULATE_SDF_TEXT] > 0.5) {
    //     uv = compute_pixelated_uv(uv, texSize, colorAdjust2.z);
        
    //     vec4 color = texture(texture0, uv);
    //     color = compute_outline(color, uv, texSize);
    //     color = compute_color_adjust(color, colorAdjust1, colorAdjust2.xy);
    //     color = compute_silhouette(color);
    //     finalColor = color;
    //     gl_FragDepth = colorAdjust2.w;
    //     return;
    // }
    
    uv = compute_tile(uv, texSize);
    uv = compute_pixelated_uv(uv, texSize, colorAdjust2.z);

    vec4 color = compute_blur(uv, texSize, blur);
    color = compute_outline(color, uv, texSize);

    if (color.a * fragColor.a < 0.004)
        discard;
    
    color = compute_color_adjust(color, colorAdjust1, colorAdjust2.xy);
    color = compute_silhouette(color);
    
    finalColor = color * fragColor;
    gl_FragDepth = colorAdjust2.w;
}
