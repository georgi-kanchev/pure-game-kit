#version 330

in vec2 fragTexCoord;
in vec4 fragColor;

out vec4 finalColor;

#define TEXTURE_W 0
#define TEXTURE_H 1
#define BLUR_X 2
#define BLUR_Y 3
#define GAMMA 4
#define SATURATION 5
#define CONTRAST 6
#define BRIGHTNESS 7
#define GRAYSCALE 8
#define INVERSION 9
#define PIXEL_SIZE 10
#define DEPTH_Z 11
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
#define SKIP_COLOR_ADJUST 26

uniform sampler2D texture0;
uniform sampler2D tileData;
uniform float u[32];

float map(float value, float min1, float max1, float min2, float max2) {
    return min2 + (value - min1) * (max2 - min2) / (max1 - min1);
}

vec2 compute_pixelated_uv(vec2 uv) {
    float pixelSize = u[PIXEL_SIZE];
    vec2 texSize = vec2(u[TEXTURE_W], u[TEXTURE_H]);
    vec2 numBlocks = texSize / max(pixelSize, 0.001); // Avoid div by zero
    vec2 pixelated = (floor(uv * numBlocks) + 0.5) / numBlocks;
    return pixelSize <= 1.0 ? uv : pixelated;
}
vec4 compute_blur(vec2 uv) {
    vec2 blur = vec2(u[BLUR_X], u[BLUR_Y]);
    if (blur.x == 0.0 && blur.y == 0.0)
        return texture(texture0, uv);
    
	vec2 texSize = vec2(u[TEXTURE_W], u[TEXTURE_H]);
    vec2 res = 1.0 / texSize;
    vec2 offset = (blur + 0.5) * res;
    vec4 sum = texture(texture0, uv + vec2(-offset.x, -offset.y));
    sum += texture(texture0, uv + vec2(offset.x, -offset.y));
    sum += texture(texture0, uv + vec2(-offset.x, offset.y));
    sum += texture(texture0, uv + vec2(offset.x, offset.y));
    return sum * 0.25;
}
vec4 compute_outline(vec4 color, vec2 uv) {
    float outline = u[OUTLINE_SIZE];
    if (color.a > 0 || outline == 0.0)
		return color;
    
	vec2 texSize = vec2(u[TEXTURE_W], u[TEXTURE_H]);
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
vec4 compute_color_adjust(vec4 color) {
    if (u[SKIP_COLOR_ADJUST] > 0.0)
        return color;

    float gam = u[GAMMA];
    float sat = u[SATURATION];
    float con = u[CONTRAST];
    float bri = u[BRIGHTNESS];
    float luminance = dot(color.rgb, vec3(0.2126, 0.7152, 0.0722));
    float gamma = gam < 0.5 ? map(gam, 0.0, 0.5, 6.0, 1.0) : map(gam, 0.5, 1.0, 1.0, 0.0);
    float saturation = sat < 0.5 ? map(sat, 0.0, 0.5, 0.0, 1.0) : map(sat, 0.5, 1.0, 1.0, 10.0);
    float contrast = con < 0.5 ? map(con, 0.0, 0.5, 0.0, 1.0) : map(con, 0.5, 1.0, 1.0, 3.0);
    float brightness = bri < 0.5 ? map(bri, 0.0, 0.5, 0.0, 1.0) : map(bri, 0.5, 1.0, 1.0, 4.0);
    color.rgb = pow(max(color.rgb, vec3(0.0)), vec3(gamma));
    color.rgb = mix(vec3(luminance), color.rgb, saturation);
    color.rgb = mix(vec3(0.5), color.rgb, contrast);
    color.rgb = mix(color.rgb, vec3(luminance), u[GRAYSCALE]);
    color.rgb = mix(color.rgb, 1.0 - color.rgb, u[INVERSION]);
    color.rgb *= brightness;
    return color;
}
vec4 compute_silhouette(vec4 color) {
    vec4 silColor = vec4(u[SILHOUETTE_R], u[SILHOUETTE_G], u[SILHOUETTE_B], u[SILHOUETTE_A]);
    color.rgb = mix(color.rgb, silColor.rgb, silColor.a);
    return color;
}
vec2 compute_tile(vec2 uv) {
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
    
    bool flip = (gid & 0x80000000u) != 0u;       // bits 31..31
    uint rot = (gid & 0x60000000u) >> 29;        // bits 30..29
    uint animCount = (gid & 0x1E000000u) >> 25;  // bits 28..25
    uint animOffset = (gid & 0x01E00000u) >> 21; // bits 24..21
    uint speedRaw = (gid & 0x001F0000u) >> 16;   // bits 20..16
    uint atlasBase = gid & 0xFFFFu;              // bits 15..00
    
    float s = float(speedRaw); // multiplier logic: 0..10 maps to 0.00..1.00; 11..31 maps to 1.33..10.00
    float multiplier = (s <= 10.0) ? (s * 0.1) : (1.0 + (s - 10.0) * 0.45);
    uint frameRange = animCount + 1u;
    uint currentFrame = uint(mod(floor(u[TIME] * multiplier) + float(animOffset), float(frameRange)));
    uint atlasIndex = atlasBase + currentFrame;
    
    float atlasCols = floor(u[TEXTURE_W] / u[TILE_W]);
    vec2 coord = vec2(mod(float(atlasIndex), atlasCols), floor(float(atlasIndex) / atlasCols));
    vec2 localUV = fract(uv * vec2(float(mapSize.x), float(mapSize.y)));
    localUV -= 0.5;
    localUV.x = flip ? -localUV.x : localUV.x;
    localUV = rot == 1u ? vec2(localUV.y, -localUV.x) : localUV;  // 90 degrees
    localUV = rot == 2u ? vec2(-localUV.x, -localUV.y) : localUV; // 180 degrees
    localUV = rot == 3u ? vec2(-localUV.y, localUV.x) : localUV;  // 270 degrees
    localUV += 0.5;
    
    vec2 atlasSizeInTiles = vec2(u[TEXTURE_W] / u[TILE_W], u[TEXTURE_H] / u[TILE_H]);
    return (coord + localUV) / atlasSizeInTiles;
}

void main() {
    vec2 uv = fragTexCoord;
    uv = compute_tile(uv);
    uv = compute_pixelated_uv(uv);
    vec4 color = compute_blur(uv);
    color = compute_outline(color, uv);

    if (color.a * fragColor.a < 0.004)
        discard;
     
    color = compute_color_adjust(color);
    color = compute_silhouette(color);

    finalColor = color * fragColor;
    gl_FragDepth = u[DEPTH_Z];
}