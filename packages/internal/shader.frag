#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

in vec4 fragData0; // texSize.xy + depthZ + objectType
in vec4 fragData1; // colorAdjust1 (gamma, saturation, contrast, brightness)
in vec4 fragData2; // rgbAdjust2 (roundness, pixelSize, blurX, blurY)
in vec4 fragData3; // outlineColor RGBA
in vec4 fragData4; // silhouetteColor RGBA
in vec4 fragData5; // outlineSize
in vec4 fragData6; // tileColumns + tileRows + tileSize

out vec4 finalColor;

#define TIME 0

uniform sampler2D texture0;
uniform sampler2D tileData;
uniform float u[1];

// ========================================================================

float map(float value, float min1, float max1, float min2, float max2) {
    return min2 + (value - min1) * (max2 - min2) / (max1 - min1);
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

void main() {
    vec2 texSize    = fragData0.xy;
    float depthZ    = fragData0.z;
    int objectType  = int(fragData0.w);

    vec4 colorAdjust1 = fragData1;
    vec4 rgbAdjust2   = fragData2;
    float roundness   = rgbAdjust2.x;
    float pixelSize   = rgbAdjust2.y;
    vec2 blur = rgbAdjust2.zw * 16.0;

    vec4 outlineColor    = fragData3;
    vec4 silhouetteColor = fragData4;

    float outlineSize = fragData5.x;

    float tileColumns = fragData6.x;
    float tileRows    = fragData6.y;
    float tileSize    = fragData6.z;

    // ========================================================================

    vec2 uv = fragTexCoord;

    uv = compute_tile(uv, texSize, tileColumns, tileRows, tileSize, tileSize);
    uv = compute_pixelated_uv(uv, texSize, pixelSize);

    vec4 color = compute_blur(uv, texSize, blur);
    color = compute_outline(color, uv, texSize, outlineSize, outlineColor);

    if (color.a * fragColor.a < 0.004)
        discard;

    color = compute_color_adjust(color, colorAdjust1);
    color = compute_silhouette(color, silhouetteColor);
    
    finalColor = color * fragColor;
    gl_FragDepth = depthZ;
}

// func shape(texCoord vec2, color vec4, custom vec4) vec4 {
//     var width, height = color.r, color.g
//     var fillColor = unpack(color.b)
//     var roundness, outlineSize = custom.x, custom.z
//     var outlineColor = unpack(custom.y)
    
//     var origin, size = imageSrc0Origin(), imageSrc0Size()
//     var uv = (texCoord - origin) / size
//     var halfSize = vec2(width, height) * 0.5
//     var pLocal = (uv - 0.5) * vec2(width, height)
    
//     var maxRadius = min(halfSize.x, halfSize.y)
//     var radius = abs(roundness) * maxRadius
    
//     var q = abs(pLocal) - halfSize + radius
//     var dShape = length(max(q, 0.0)) + min(max(q.x, q.y), 0.0) - radius
    
//     if roundness < 0.0 {
//         dShape = max(radius-length(max(q, 0.0)), max(abs(pLocal).x-halfSize.x, abs(pLocal).y-halfSize.y))
//     }
    
//     var dEdge = dShape
//     if outlineSize > 0.0 {
//         dEdge = dShape - outlineSize 
//     } else if outlineSize < 0.0 {
//         dEdge = dShape + abs(outlineSize) 
//     }
    
//     var af = 1.0 * Smoothing 
    
//     var sShape = 1.0 - smoothstep(-af, af, dShape)
//     var sEdge = 1.0 - smoothstep(-af, af, dEdge)

//     var finalFillColor = imageSrc0At(texCoord) * fillColor
    
//     if outlineSize > 0.0 {
//         var baseMask = outlineColor * sEdge
//         var topMask = finalFillColor * sShape
//         return topMask + baseMask * (1.0 - topMask.a)
//     } else if outlineSize < 0.0 {
//         var baseMask = finalFillColor * sShape
//         var innerEdgeMask = smoothstep(-af, af, dEdge) 
//         var topMask = outlineColor * innerEdgeMask * sShape
//         return topMask + baseMask * (1.0 - topMask.a)
//     }
    
//     return finalFillColor * sShape
// }