/*

GLSL functions used for fragment shader permutation.
This file is "somewhat parsed" and processed, so indentation and
naming patterns are significant and not subject to personal taste.

*/

vec3 fx_Color (const in vec3 vCol) {
	return uni_vec3_Color_Rgb;
}

vec3 fx_Gamma (const in vec3 vCol) {
	return pow(vCol, vec3(1 / 2.2));
}

vec3 fx_Grayscale (const in vec3 vCol) {
	return vec3((vCol.r * 0.3) + (vCol.g * 0.59) + (vCol.b * 0.11));
}

vec3 fx_Orangify (const in vec3 vCol) {
	return vec3(vCol.r + 0.65, vCol.g + 0.25, vCol.b);
}

vec3 fx_Coords (const in vec3 vCol) {
	vec2 vCoords = clamp(var_vec2_Tex2D, 0, 1) * 0.33;
	return vec3(0, vCoords);
}

vec3 fx_Tex2D (const in vec3 vCol) {
	return texture(uni_sampler2D_Tex2D_Img, var_vec2_Tex2D).rgb;
}

vec3 fx_TexCube (const in vec3 vCol) {
	return texture(uni_samplerCube_TexCube_Img, var_vec3_TexCube).rgb;
}
