#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

in vec2 fragTexCoord2;
in vec3 fragNormal;
in vec4 fragTangent;

// texCoord2.x = TextureWidth(12) + TextureHeight(12)
// texCoord2.y = BorderColor(24)

// normal.x = Gamma(6) + Saturation(6) + Contrast(6) + Brightness(6)
// normal.y = Grayscale(6) + Inversion(6) + BlurX(6) + BlurY(6)
// normal.z = DepthZ(11) + BorderSize(11) + Type(2)

// Shape:
//  tangent.x = Roundness(32)
//  tangent.y = PixelSize(32)

// Sprite:
//  tangent.x = OutlineColor(24)
//  tangent.y = SilhouetteColor(24)
//  tangent.z = OutlineSize(32)
//  tangent.w = Roundness(16) + PixelSize(8)

// Tilemap:
//  tangent.x = OutlineColor(24)
//  tangent.y = SilhouetteColor(24)
//  tangent.z = TileColumns(10) + TileRows(10) + PixelSize(4)
//  tangent.w = OutlineSize(8) + TileSize(8) + Roundness(8)

// Text:
//  tangent.x = OutlineColor(24)
//  tangent.y = ShadowColor(24)
//  tangent.z = Weight(6) + OutlineWeight(6) + ShadowWeight(6) + ShadowBlur(6)
//  tangent.w = TextShadowX(6) + TextShadowY(6) + Roundness(8) + PixelSize(4)

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

// --- unpack helpers (all read from bits 23-0 of a safe normal float32) ---

vec3 unpack_color24(float packedFloat) {
    uint bits = floatBitsToUint(packedFloat);
    float r = float((bits >> 16u) & 0xFFu) / 255.0;
    float g = float((bits >> 8u)  & 0xFFu) / 255.0;
    float b = float(bits & 0xFFu) / 255.0;
    return vec3(r, g, b);
}

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
    depthZ     = float((bits >> 13u) & 0x7FFu) / 2047.0; // 11 bits normalized
    borderSize = float((bits >> 2u)  & 0x7FFu);          // 11 bits absolute
    objType    = int(bits & 0x3u);                        //  2 bits
}

void unpack_16_8(float packedFloat, out float roundness, out float pixelSize) {
    uint bits = floatBitsToUint(packedFloat);
    roundness = float((bits >> 8u) & 0xFFFFu) / 65535.0; // 16 bits normalized
    pixelSize = float(bits & 0xFFu);                      //  8 bits absolute
}

void unpack_8_8_8(float packedFloat, out float outlineSize, out float tileSize, out float roundness) {
    uint bits = floatBitsToUint(packedFloat);
    outlineSize = float((bits >> 16u) & 0xFFu);         // 8 bits absolute
    tileSize    = float((bits >> 8u)  & 0xFFu);         // 8 bits absolute
    roundness   = float(bits & 0xFFu) / 255.0;          // 8 bits normalized
}

vec3 unpack_10_10_4(float packedFloat) {
    uint bits = floatBitsToUint(packedFloat);
    float cols = float((bits >> 14u) & 0x3FFu); // 10 bits absolute
    float rows = float((bits >> 4u)  & 0x3FFu); // 10 bits absolute
    float ps   = float(bits & 0xFu);            //  4 bits absolute
    return vec3(cols, rows, ps);
}

void unpack_6_6_8_4(float packedFloat, out float shadowX, out float shadowY, out float roundness, out float pixelSize) {
    uint bits = floatBitsToUint(packedFloat);
    // shadowX/Y: 6-bit two's complement stored raw
    int rawX = int((bits >> 18u) & 0x3Fu);
    int rawY = int((bits >> 12u) & 0x3Fu);
    shadowX = float(rawX >= 32 ? rawX - 64 : rawX) / 32.0; // maps to ~[-1, 1]
    shadowY = float(rawY >= 32 ? rawY - 64 : rawY) / 32.0;
    roundness = float((bits >> 4u) & 0xFFu) / 255.0; // 8 bits normalized
    pixelSize = float(bits & 0xFu);                   // 4 bits absolute
}

// --- effect functions ---

vec2 compute_pixelated_uv(vec2 uv, vec2 texSize, float pixelSize) {
    vec2 numBlocks = texSize / max(pixelSize, 0.001);
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
    
    float gamma = exp2((0.5 - gam) * 5.0);
    float saturation = 10.0 * pow(sat, log2(10.0));
    float contrast = 3.0  * pow(con, log2(3.0));
    float brightness = 4.0  * bri * bri;

    color.rgb = pow(max(color.rgb, vec3(0.0)), vec3(gamma));
    float lum_pre_sat = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
    color.rgb = mix(vec3(lum_pre_sat), color.rgb, saturation);
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
//     
//     uint thickIdx = (c.a >> 6) & 0x03u;
//     uint outlIdx = (c.a >> 4) & 0x03u;
//     uint shadIdx = (c.a >> 2) & 0x03u;
//     uint smoothIdx = (c.a) & 0x03u;
//     
//     float thick[4] = float[](0.35, 0.50, 0.65, 0.80);
//     float smooths[4] = float[](0.50, 4.00, 8.00, 12.0);
//     
//     vec2 shadowOffset = vec2(u[TEXT_SHADOW_X], u[TEXT_SHADOW_Y]);
//     float shadowDistance = texture(texture0, uv - shadowOffset).a - (1.0 - thick[shadIdx]);
//     float shadowSmooth = smooths[smoothIdx] * length(vec2(dFdx(shadowDistance), dFdy(shadowDistance)));
//     float shadowAlpha = shadowColor.a * smoothstep(-shadowSmooth, shadowSmooth, shadowDistance);
//     
//     float distance = texture(texture0, uv).a - (1.0 - thick[thickIdx]);
//     float baseSmooth = 0.5 * length(vec2(dFdx(distance), dFdy(distance)));
//     float sdfAlpha = base.a * smoothstep(-baseSmooth, baseSmooth, distance);
//     
//     float compressedOutlIdx = map(float(outlIdx), 0.0, 3.0, 0.7, 2.9);
//     float outlineThick = (1.0 - thick[thickIdx]) * (compressedOutlIdx / 3.0);
//     float outlineAlpha = outlineColor.a * smoothstep(-baseSmooth, baseSmooth, distance + outlineThick);
//     
//     vec3 mixedRGB = mix(shadowColor.rgb, outlineColor.rgb, outlineAlpha);
//     mixedRGB = mix(mixedRGB, base.rgb, sdfAlpha);
//     float mixedAlpha = max(shadowAlpha, max(outlineAlpha, sdfAlpha));
//     
//     vec3 finalRGB = distance > sdfAlpha ? base.rgb : mixedRGB;
//     float finalAlpha = distance > sdfAlpha ? base.a : mixedAlpha;
//     
//     return vec4(finalRGB, finalAlpha);
// }

void main() {
    // --- unpack vertex attributes (24-bit safe float32 in bits 23-0) ---

    vec2 texSize = unpack_12_12(fragTexCoord2.x);  // TextureWidth(12) + TextureHeight(12)

    vec3 borderColor = unpack_color24(fragTexCoord2.y); // BorderColor(24)

    vec4 colorAdjust1 = unpack_6_6_6_6(fragNormal.x); // Gamma(6)+Saturation(6)+Contrast(6)+Brightness(6)
    vec4 colorAdjust2 = unpack_6_6_6_6(fragNormal.y); // Grayscale(6)+Inversion(6)+BlurX(6)+BlurY(6)
    vec2 blur = colorAdjust2.zw * 16.0;

    float depthZ, borderSize;
    int objectType;
    unpack_11_11_2(fragNormal.z, depthZ, borderSize, objectType); // DepthZ(11)+BorderSize(11)+Type(2)

    // --- per-type tangent unpacking ---

    float roundness = 0.0;
    float pixelSize = 0.0;
    // outlineColor, silhouetteColor, shadowColor, outlineSize, tileCols, tileRows, tileSize
    // are unpacked below per type; compute functions still use u[] uniforms for now.

    if (objectType == 0) { // Shape
        roundness = fragTangent.x; // full float
        pixelSize = fragTangent.y; // full float
    }
    else if (objectType == 1) { // Sprite
        vec3 outlineColor   = unpack_color24(fragTangent.x); // OutlineColor(24)
        vec3 silhouetteColor = unpack_color24(fragTangent.y); // SilhouetteColor(24)
        float outlineSize   = fragTangent.z;                  // full float
        unpack_16_8(fragTangent.w, roundness, pixelSize);     // Roundness(16)+PixelSize(8)
    }
    else if (objectType == 2) { // Text
        vec3 outlineColor = unpack_color24(fragTangent.x);             // OutlineColor(24)
        vec3 shadowColor  = unpack_color24(fragTangent.y);             // ShadowColor(24)
        vec4 textWeights  = unpack_6_6_6_6(fragTangent.z);             // Weight+OutlineWeight+ShadowWeight+ShadowBlur
        float shadowX, shadowY;
        unpack_6_6_8_4(fragTangent.w, shadowX, shadowY, roundness, pixelSize); // ShadowX+ShadowY+Roundness+PixelSize
    }
    else if (objectType == 3) { // Tilemap
        vec3 outlineColor   = unpack_color24(fragTangent.x); // OutlineColor(24)
        vec3 silhouetteColor = unpack_color24(fragTangent.y); // SilhouetteColor(24)
        vec3 tileInfo = unpack_10_10_4(fragTangent.z);        // TileColumns(10)+TileRows(10)+PixelSize(4)
        float tileColumns = tileInfo.x;
        float tileRows    = tileInfo.y;
        pixelSize         = tileInfo.z;
        float outlineSize, tileSize;
        unpack_8_8_8(fragTangent.w, outlineSize, tileSize, roundness); // OutlineSize(8)+TileSize(8)+Roundness(8)
    }

    // ========================================================================

    vec2 uv = fragTexCoord;
    // if (u[CALCULATE_SDF_TEXT] > 0.5) {
    //     uv = compute_pixelated_uv(uv, texSize, pixelSize);
    //     
    //     vec4 color = texture(texture0, uv);
    //     color = compute_outline(color, uv, texSize);
    //     color = compute_color_adjust(color, colorAdjust1, colorAdjust2.xy);
    //     color = compute_silhouette(color);
    //     finalColor = color;
    //     gl_FragDepth = depthZ;
    //     return;
    // }

    uv = compute_tile(uv, texSize);
    uv = compute_pixelated_uv(uv, texSize, pixelSize);

    vec4 color = compute_blur(uv, texSize, blur);
    color = compute_outline(color, uv, texSize);

    if (color.a * fragColor.a < 0.004)
        discard;

    color = compute_color_adjust(color, colorAdjust1, colorAdjust2.xy);
    color = compute_silhouette(color);

    finalColor = color * fragColor;
    gl_FragDepth = depthZ;
}
