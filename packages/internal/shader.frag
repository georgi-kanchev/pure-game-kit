#version 330

in vec2 fragTexCoord;
in vec4 fragColor;
in vec2 fragLocalPos;

out vec4 finalColor;

uniform sampler2D texture0;
uniform sampler2D tileData;
uniform float u[33];

#define KIND_SHAPE 0
#define KIND_SPRITE 1
#define KIND_TEXT 2
#define KIND_TILEMAP 3

#define TIME            u[0]
#define TEX_W           u[1]
#define TEX_H           u[2]
#define OBJ_KIND        u[4]
#define GAM             u[5]
#define SAT             u[6]
#define CON             u[7]
#define BRI             u[8]
#define ROUNDNESS       u[9]
#define PIXEL_SIZE      u[10]
#define BLUR_X          u[11]
#define BLUR_Y          u[12]
#define OUTLINE_R       u[13]
#define OUTLINE_G       u[14]
#define OUTLINE_B       u[15]
#define OUTLINE_A       u[16]
#define SIL_R           u[17]
#define SIL_G           u[18]
#define SIL_B           u[19]
#define SIL_A           u[20]
#define OUTLINE_SIZE    u[21]
#define BORDER_SIZE     u[22]
#define SHADOW_WEIGHT   u[23]
#define NO_COLOR_ADJUST u[24]
#define SHADOW_X        u[25]
#define SHADOW_Y        u[26]
#define SHADOW_BLUR     u[27]
#define CROP_V          u[28]
#define BORDER_R        u[29]
#define BORDER_G        u[30]
#define BORDER_B        u[31]
#define BORDER_A        u[32]

// ========================================================================

float map(float value, float min1, float max1, float min2, float max2) {
    return min2 + (value - min1) * (max2 - min2) / (max1 - min1);
}
float shape_sdf(vec2 pLocal, vec2 halfExtents, float r, float roundness) {
    vec2 q = abs(pLocal) - halfExtents + r;
    float dShape = length(max(q, 0.0)) + min(max(q.x, q.y), 0.0) - r;
    
    if (roundness < 0.0) { // Handle inverted roundness
        dShape = max(r - length(max(q, 0.0)), max(abs(pLocal).x - halfExtents.x, abs(pLocal).y - halfExtents.y));
    }
    return dShape;
}
float median(vec3 rgb) {
    return max(min(rgb.r, rgb.g), min(max(rgb.r, rgb.g), rgb.b));
}

vec2 do_pixelated_uv(vec2 uv) {
    float pixelSize = PIXEL_SIZE / 1.5;
    vec2 texSize = vec2(TEX_W, TEX_H);
    vec2 numBlocks = texSize / max(pixelSize, 0.001);
    vec2 pixelated = (floor(uv * numBlocks) + 0.5) / numBlocks;
    return pixelSize <= 1.0 ? uv : pixelated;
}
vec4 do_blur(vec2 uv) {
    vec2 blur = vec2(BLUR_X, BLUR_Y) * 16.0;
    if (blur.x == 0.0 && blur.y == 0.0)
        return texture(texture0, uv);
    
    blur /= 8.0; // adjust
    vec2 texSize = vec2(TEX_W, TEX_H);
    vec2 res = 1.0 / texSize;
    vec2 offset = (blur + 0.5) * res;
    vec4 sum = texture(texture0, uv + vec2(-offset.x, -offset.y));
    sum += texture(texture0, uv + vec2(offset.x, -offset.y));
    sum += texture(texture0, uv + vec2(-offset.x, offset.y));
    sum += texture(texture0, uv + vec2(offset.x, offset.y));
    return sum * 0.25;
}
vec4 do_outline(vec4 color, vec2 uv) {
    float outlineSize = OUTLINE_SIZE;
    if (color.a > 0 || outlineSize == 0.0)
        return color;
    
    vec2 texSize = vec2(TEX_W, TEX_H);
    vec2 texel = 1.0 / texSize;
    vec4 outlineColor = vec4(OUTLINE_R, OUTLINE_G, OUTLINE_B, OUTLINE_A);

    if (texture(texture0, uv + vec2(texel.x * outlineSize, 0.0)).a > 0.0 ||
            texture(texture0, uv + vec2(-texel.x * outlineSize, 0.0)).a > 0.0 ||
            texture(texture0, uv + vec2(0.0, texel.y * outlineSize)).a > 0.0 ||
            texture(texture0, uv + vec2(0.0, -texel.y * outlineSize)).a > 0.0 ||
            texture(texture0, uv + vec2(texel.x * outlineSize, texel.y * outlineSize)).a > 0.0 ||
            texture(texture0, uv + vec2(-texel.x * outlineSize, texel.y * outlineSize)).a > 0.0 ||
            texture(texture0, uv + vec2(-texel.x * outlineSize, -texel.y * outlineSize)).a > 0.0 ||
            texture(texture0, uv + vec2(texel.x * outlineSize, -texel.y * outlineSize)).a > 0.0)
        return outlineColor;

    return color;
}
vec4 do_silhouette(vec4 color) {
    color.rgb = mix(color.rgb, vec3(SIL_R, SIL_G, SIL_B), SIL_A);
    return color;
}
vec4 do_color_adjust(vec4 color) {
    float gam = GAM;
    float sat = SAT;
    float con = CON;
    float bri = BRI;
    
    float gamma = exp2((0.5 - gam) * 5.0);
    float saturation = 10.0 * pow(sat, log2(10.0));
    float contrast = 3.0  * pow(con, log2(3.0));
    float brightness = 4.0  * bri * bri;

    color.rgb = pow(max(color.rgb, vec3(0.0)), vec3(gamma));
    float lum_pre_sat = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
    color.rgb = mix(vec3(lum_pre_sat), color.rgb, saturation);
    color.rgb = mix(vec3(0.5), color.rgb, contrast);
    color.rgb *= brightness;
    return color;
}
vec2 do_tile(float tileColumns, float tileRows, float tileW, float tileH) {
    vec2 uv = fragTexCoord;
    if (tileColumns == 0.0)
        return uv;
    
    vec2 texSize = vec2(TEX_W, TEX_H);
    ivec2 mapSize = ivec2(int(tileColumns), int(tileRows));
    ivec2 tile = ivec2(int(uv.x * float(mapSize.x)), int(uv.y * float(mapSize.y)));
    tile = clamp(tile, ivec2(0), mapSize - 1);
    int linearTileID = tile.y * mapSize.x + tile.x;
    ivec2 dataUv = ivec2(linearTileID % mapSize.x, linearTileID / mapSize.x);
    vec4 texColor = texelFetch(tileData, dataUv, 0);
    uvec4 b = uvec4(texColor * 255.0 + 0.5);
    uint gid = (b.r << 24) | (b.g << 16) | (b.b << 8) | b.a;

    bool flip = (gid & 0x80000000u) != 0u;
    uint rot = (gid & 0x60000000u) >> 29;
    uint animCount = (gid & 0x1E000000u) >> 25;
    uint animOffset = (gid & 0x01E00000u) >> 21;
    uint speedRaw = (gid & 0x001F0000u) >> 16;
    uint atlasBase = gid & 0xFFFFu;

    float s = float(speedRaw);
    float multiplier = (s <= 10.0) ? (s * 0.1) : (1.0 + (s - 10.0) * 0.45);
    uint frameRange = animCount + 1u;
    uint currentFrame = uint(mod(floor(TIME * multiplier) + float(animOffset), float(frameRange)));
    uint atlasIndex = atlasBase + currentFrame;

    float atlasCols = floor(texSize.x / tileW);
    vec2 coord = vec2(mod(float(atlasIndex), atlasCols), floor(float(atlasIndex) / atlasCols));
    vec2 localUV = fract(uv * vec2(float(mapSize.x), float(mapSize.y)));
    localUV -= 0.5;
    localUV.x = flip ? -localUV.x : localUV.x;
    localUV = rot == 1u ? vec2(localUV.y, -localUV.x) : localUV;
    localUV = rot == 2u ? vec2(-localUV.x, -localUV.y) : localUV;
    localUV = rot == 3u ? vec2(-localUV.y, localUV.x) : localUV;
    localUV += 0.5;
    
    vec2 atlasSizeInTiles = vec2(texSize.x / tileW, texSize.y / tileH);
    return (coord + localUV) / atlasSizeInTiles;
}
vec4 do_sdf_shape(vec4 color, vec2 cropBoundsU, vec2 cropBoundsV) {
    float roundness = ROUNDNESS;
    float borderSize = BORDER_SIZE;
    vec4 borderColor = vec4(BORDER_R, BORDER_G, BORDER_B, BORDER_A);
    if (abs(roundness) < 0.001 && abs(borderSize) < 0.001)
        return color;
    
    vec2 cropRange = max(vec2(cropBoundsU.y - cropBoundsU.x, cropBoundsV.y - cropBoundsV.x), 0.001);
    vec2 pLocal = (fragLocalPos - vec2(cropBoundsU.x, cropBoundsV.x)) / cropRange - 0.5;
    vec2 pixPerUnit = 1.0 / max(fwidth(pLocal), 0.0001);
    pLocal *= pixPerUnit;
    vec2 halfSize = pixPerUnit * 0.5;
    float maxRadius = min(halfSize.x, halfSize.y);
    float radius = abs(roundness) * maxRadius;
    float dShape = shape_sdf(pLocal, halfSize, radius, roundness);
    float af = fwidth(dShape) * 1.5;
    float absBorder = abs(borderSize);

    if (absBorder > 0.0) {
        float sShape = 1.0 - smoothstep(-af, af, dShape);
        
        if (borderSize > 0.0) { 
            vec2 outerSize = halfSize + absBorder;
            float outerRadius = abs(roundness) * min(outerSize.x, outerSize.y);
            float dOuter = shape_sdf(pLocal, outerSize, outerRadius, roundness);
            float sOuter = 1.0 - smoothstep(-af, af, dOuter);
            float sRing = max(sOuter - sShape, 0.0);
            
            vec4 fill = color * sShape;
            vec4 ring = borderColor * sRing;
            return vec4(fill.rgb + ring.rgb, fill.a + ring.a);
        } else { 
            vec2 innerSize = max(halfSize - absBorder, vec2(0.0));
            float innerRadius = abs(roundness) * min(innerSize.x, innerSize.y);
            float dInner = shape_sdf(pLocal, innerSize, innerRadius, roundness);
            float sInner = 1.0 - smoothstep(-af, af, dInner);
            float sRing = max(sShape - sInner, 0.0);
            
            vec4 fill = color * sInner;
            vec4 ring = borderColor * sRing;
            return vec4(fill.rgb + ring.rgb, fill.a + ring.a);
        }
    }

    float sShape = 1.0 - smoothstep(-af, af, dShape);
    return color * sShape;
}
vec4 do_msdf_text() {
    vec2 uv = fragTexCoord;
    vec4 outlineColor = vec4(OUTLINE_R, OUTLINE_G, OUTLINE_B, OUTLINE_A);
    vec4 shadowColor = vec4(SIL_R, SIL_G, SIL_B, SIL_A);
    float weight = OUTLINE_SIZE;
    float outlineWeight = BORDER_SIZE;
    float shadowWeight = SHADOW_WEIGHT;
    vec2 shadow = vec2(SHADOW_X, SHADOW_Y);
    float shadowBlur = SHADOW_BLUR;
    
    float pxRange = 8.0;
    vec2 unitRange = vec2(pxRange) / vec2(textureSize(texture0, 0));
    vec2 screenTexSize = vec2(1.0) / fwidth(uv);
    float screenPxRange = max(0.5 * dot(unitRange, screenTexSize), 1.0);
    
    float baseSample = median(texture(texture0, uv).rgb);
    float shadowSample = median(texture(texture0, uv + vec2(shadow.x, shadow.y) / vec2(textureSize(texture0, 0))).rgb);
    
    float basePxDist = screenPxRange * (baseSample - 0.5);
    float shadowPxDist = screenPxRange * (shadowSample - 0.5);
    
    float thickness = weight * screenPxRange * 0.25;
    float textPxDist = basePxDist + thickness;
    float sdfAlpha = fragColor.a * smoothstep(-0.5, 0.5, textPxDist);
    
    float outlinePxDist = textPxDist + map(outlineWeight, 0.0, 1.0, 0, screenPxRange*0.4);
    float outlineAlpha = outlineColor.a * smoothstep(-0.5, 0.5, outlinePxDist);
    outlineAlpha = max(0.0, outlineAlpha - sdfAlpha);
    
    float shadowThickness = shadowWeight * screenPxRange * 0.25;
    float shadowSmooth = 0.5 + shadowBlur/128 * screenPxRange * 0.25;
    float shadowAlpha = shadowColor.a * smoothstep(-shadowSmooth, shadowSmooth, shadowPxDist + shadowThickness);
    
    vec3 rgb = mix(shadowColor.rgb, outlineColor.rgb, outlineAlpha);
    rgb = mix(rgb, fragColor.rgb, sdfAlpha);
    
    float alpha = max(shadowAlpha, max(outlineAlpha, sdfAlpha));
    // if (alpha < 0.004) // debug character bounding box
    //     return vec4(1.0, 0.0, 0.0, 0.5);

    return vec4(rgb, alpha);
}

void main() {
    int objKind = int(OBJ_KIND);
    bool hasCrop = objKind == KIND_SHAPE || objKind == KIND_SPRITE;
    float tileColumns = hasCrop ? 0.0 : SHADOW_X;
    float tileRows    = hasCrop ? 0.0 : SHADOW_Y;
    float tileSize    = hasCrop ? 0.0 : SHADOW_BLUR;
    vec2  cropBoundsU = hasCrop ? vec2(SHADOW_X, SHADOW_Y) / 4095.0 - 0.5 : vec2(-0.5, 0.5);
    vec2  cropBoundsV = hasCrop ? vec2(SHADOW_BLUR, CROP_V) / 4095.0 - 0.5 : vec2(-0.5, 0.5);
    
    // ========================================================================

    vec4 color;
    
    if (objKind == KIND_TEXT) { // Text: MSDF path (skip do_tile: text reuses tile slots for shadow data)
        color = do_msdf_text();
        if (color.a < 0.004)
            discard;

        if (NO_COLOR_ADJUST < 0.5) color = do_color_adjust(color);
        finalColor = color;
        return;
    }

    vec2 uv = do_tile(tileColumns, tileRows, tileSize, tileSize);
    if (objKind == KIND_SHAPE) { // Shape: use vertex color as fill, skip pixelate/blur
        color = fragColor;
    } else { // Sprite / Tilemap
        uv = do_pixelated_uv(uv);
        color = do_blur(uv);
        color = do_outline(color, uv);
        color *= fragColor;
    }
    
    if (objKind != KIND_SHAPE)
        color = do_silhouette(color);
    
    color = do_sdf_shape(color, cropBoundsU, cropBoundsV);
    
    if (objKind != KIND_SHAPE && NO_COLOR_ADJUST < 0.5)
        color = do_color_adjust(color);
    
    finalColor = color;
}