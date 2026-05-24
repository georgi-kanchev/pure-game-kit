#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

in vec4 fragData0;
in vec4 fragData1;
in vec4 fragData2;
in vec4 fragData3;
in vec4 fragData4;
in vec4 fragData5;
in vec4 fragData6;
in vec4 fragData7;

out vec4 finalColor;

#define KIND_SHAPE 0
#define KIND_SPRITE 1
#define KIND_TEXT 2
#define KIND_TILEMAP 3

#define TIME 0

uniform sampler2D texture0;
uniform sampler2D tileData;
uniform float u[1];

// ========================================================================

float map(float value, float min1, float max1, float min2, float max2) {
    return min2 + (value - min1) * (max2 - min2) / (max1 - min1);
}
float shape_sdf(vec2 pLocal, vec2 halfExtents, float r, float roundness) {
    vec2 q = abs(pLocal) - halfExtents + r;
    float dShape = length(max(q, 0.0)) + min(max(q.x, q.y), 0.0) - r;
    
    // Handle inverted roundness
    if (roundness < 0.0) {
        dShape = max(r - length(max(q, 0.0)), max(abs(pLocal).x - halfExtents.x, abs(pLocal).y - halfExtents.y));
    }
    return dShape;
}
float median(vec3 rgb) {
    return max(min(rgb.r, rgb.g), min(max(rgb.r, rgb.g), rgb.b));
}

vec2 compute_pixelated_uv(vec2 uv, vec2 texSize, float pixelSize) {
    pixelSize /= 1.5;
    vec2 numBlocks = texSize / max(pixelSize, 0.001);
    vec2 pixelated = (floor(uv * numBlocks) + 0.5) / numBlocks;
    return pixelSize <= 1.0 ? uv : pixelated;
}
vec4 compute_blur(vec2 uv, vec2 texSize, vec2 blur) {
    if (blur.x == 0.0 && blur.y == 0.0)
        return texture(texture0, uv);
    
    blur /= 8.0; // adjust
    vec2 res = 1.0 / texSize;
    vec2 offset = (blur + 0.5) * res;
    vec4 sum = texture(texture0, uv + vec2(-offset.x, -offset.y));
    sum += texture(texture0, uv + vec2(offset.x, -offset.y));
    sum += texture(texture0, uv + vec2(-offset.x, offset.y));
    sum += texture(texture0, uv + vec2(offset.x, offset.y));
    return sum * 0.25;
}
vec4 compute_outline(vec4 color, vec2 uv, vec2 texSize, float outlineSize, vec4 outlineColor) {
    if (color.a > 0 || outlineSize == 0.0)
        return color;
    
    vec2 texel = 1.0 / texSize;

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
vec4 compute_silhouette(vec4 color, vec4 silhouetteColor) {
    color.rgb = mix(color.rgb, silhouetteColor.rgb, silhouetteColor.a);
    return color;
}
vec4 compute_color_adjust(vec4 color, vec4 colorAdjust) {
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
    color.rgb *= brightness;
    return color;
}
vec2 compute_tile(vec2 uv, vec2 texSize, float tileColumns, float tileRows, float tileW, float tileH) {
    if (tileColumns == 0.0)
        return uv;

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
    uint currentFrame = uint(mod(floor(u[TIME] * multiplier) + float(animOffset), float(frameRange)));
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

    localUV = mix(localUV, vec2(0.5), 1.0 / vec2(tileW, tileH));

    vec2 atlasSizeInTiles = vec2(texSize.x / tileW, texSize.y / tileH);
    return (coord + localUV) / atlasSizeInTiles;
}
vec4 compute_sdf_shape(vec2 uv, vec2 texSize, vec4 color, float roundness, float borderSize, vec4 borderColor) {
    if (abs(roundness) < 0.001 && abs(borderSize) < 0.001) { return color; }
    
    vec2 halfSize = texSize * 0.5;
    vec2 pLocal = (uv - 0.5) * texSize;

    // Shapes use a 1×1 texture with UV [0,1] — compensate screen-space aspect ratio
    // so corners are circular on screen. Skip for sprites/atlas textures.
    if (texSize.x < 2.0 && texSize.y < 2.0) {
        vec2 sd = fwidth(uv);
        float scaleX = max(sd.y / max(sd.x, 0.0001), 1.0);
        float scaleY = max(sd.x / max(sd.y, 0.0001), 1.0);
        pLocal *= vec2(scaleX, scaleY);
        halfSize *= vec2(scaleX, scaleY);
    }
    
    float maxRadius = min(halfSize.x, halfSize.y);
    float radius = abs(roundness) * maxRadius;
    
    // Base shape SDF
    float dShape = shape_sdf(pLocal, halfSize, radius, roundness);
    
    // Calculate anti-aliasing factor
    float af = fwidth(dShape) * 1.5;
    float absBorder = abs(borderSize);

    if (absBorder > 0.0) {
        float sShape = 1.0 - smoothstep(-af, af, dShape);
        float sRing;
        
        if (borderSize > 0.0) { 
            // OUTER RING: Generate a larger shape that scales its corner radius to match the same roundness ratio
            vec2 outerSize = halfSize + absBorder;
            float outerRadius = abs(roundness) * min(outerSize.x, outerSize.y);
            float dOuter = shape_sdf(pLocal, outerSize, outerRadius, roundness);
            
            float sOuter = 1.0 - smoothstep(-af, af, dOuter);
            sRing = max(sOuter - sShape, 0.0);
            
        } else { 
            // INNER RING: Generate a smaller shape
            vec2 innerSize = max(halfSize - absBorder, vec2(0.0));
            float innerRadius = abs(roundness) * min(innerSize.x, innerSize.y);
            float dInner = shape_sdf(pLocal, innerSize, innerRadius, roundness);
            
            float sInner = 1.0 - smoothstep(-af, af, dInner);
            sRing = max(sShape - sInner, 0.0);
        }
        
        vec4 fill = color * sShape;
        vec4 ring = borderColor * sRing;
        return ring + fill * (1.0 - ring.a);
    }

    float sShape = 1.0 - smoothstep(-af, af, dShape);
    return color * sShape;
}
vec4 compute_msdf_text(vec2 uv, vec4 baseColor, vec4 outlineColor) {
    vec4 shadowColor = fragData4;
    float weight = fragData5.x;
    float outlineWeight = fragData5.y;
    float shadowWeight = fragData5.z;
    float shadowX = fragData6.x;
    float shadowY = fragData6.y;
    float shadowBlur = fragData6.z;
    
    float pxRange = 8.0;
    vec2 unitRange = vec2(pxRange) / vec2(textureSize(texture0, 0));
    vec2 screenTexSize = vec2(1.0) / fwidth(uv);
    float screenPxRange = max(0.5 * dot(unitRange, screenTexSize), 1.0);
    
    float baseSample = median(texture(texture0, uv).rgb);
    float shadowSample = median(texture(texture0, uv + vec2(shadowX, shadowY) / vec2(textureSize(texture0, 0))).rgb);
    
    float basePxDist = screenPxRange * (baseSample - 0.5);
    float shadowPxDist = screenPxRange * (shadowSample - 0.5);
    
    float thickness = weight * screenPxRange * 0.25;
    float textPxDist = basePxDist + thickness;
    float sdfAlpha = baseColor.a * smoothstep(-0.5, 0.5, textPxDist);
    
    float outlinePxDist = textPxDist + map(outlineWeight, 0.0, 1.0, 0, screenPxRange*0.4);
    float outlineAlpha = outlineColor.a * smoothstep(-0.5, 0.5, outlinePxDist);
    outlineAlpha = max(0.0, outlineAlpha - sdfAlpha);
    
    float shadowThickness = shadowWeight * screenPxRange * 0.25;
    float shadowSmooth = 0.5 + shadowBlur/128 * screenPxRange * 0.25;
    float shadowAlpha = shadowColor.a * smoothstep(-shadowSmooth, shadowSmooth, shadowPxDist + shadowThickness);
    
    vec3 rgb = mix(shadowColor.rgb, outlineColor.rgb, outlineAlpha);
    rgb = mix(rgb, baseColor.rgb, sdfAlpha);
    
    float alpha = max(shadowAlpha, max(outlineAlpha, sdfAlpha));
    return vec4(rgb, alpha);
}

void main() {
    vec2 texSize = fragData0.xy;
    float depthZ = fragData0.z;
    int objKind = int(fragData0.w);

    vec4 colorAdjust1 = fragData1;
    vec4 data2 = fragData2;
    float roundness = data2.x;
    float pixelSize = data2.y;
    vec2 blur = data2.zw * 16.0;

    vec4 outlineColor = fragData3;
    vec4 silhouetteColor = fragData4;

    float outlineSize = fragData5.x;
    float borderSize = fragData5.y;
    vec4  borderColor = fragData7;

    float tileColumns = fragData6.x;
    float tileRows = fragData6.y;
    float tileSize = fragData6.z;
    
    // ========================================================================

    vec4 color;

    if (objKind == KIND_TEXT) { // Text: MSDF path (skip compute_tile: text reuses tile slots for shadow data)
        color = compute_msdf_text(fragTexCoord, fragColor, outlineColor);
        if (color.a < 0.004)
            discard;

        color = compute_color_adjust(color, colorAdjust1);
        finalColor = color;
    } else {
        vec2 uv = compute_tile(fragTexCoord, texSize, tileColumns, tileRows, tileSize, tileSize);
        if (objKind == KIND_SHAPE) { // Shape: white fill, skip pixelate/blur
            color = vec4(1.0);
        } else { // Sprite / Tilemap
            uv = compute_pixelated_uv(uv, texSize, pixelSize);
            color = compute_blur(uv, texSize, blur);
            color = compute_outline(color, uv, texSize, outlineSize, outlineColor);
        }

        color = compute_sdf_shape(uv, texSize, color, roundness, borderSize, borderColor);

        if (color.a * fragColor.a < 0.004)
            discard;

        if (objKind != KIND_SHAPE) {
            color = compute_color_adjust(color, colorAdjust1);
            color = compute_silhouette(color, silhouetteColor);
        }

        finalColor = color * fragColor;
    }

    gl_FragDepth = depthZ;
}