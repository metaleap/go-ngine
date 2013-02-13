/*

GLSL functions used for fragment shader permutation.
This file is "somewhat parsed" and processed, so indentation and
naming patterns are significant and not subject to personal taste.

*/

void fx_Grayscale (inout vec3 vCol) {
	vCol = vec3((vCol.r * 0.3) + (vCol.g * 0.59) + (vCol.b * 0.11));
}

void fx_RedTest (inout vec3 vCol) {
	vCol = vec3(0.9, 0.1, 0.1);
}

void fx_Tex2D (inout vec3 vCol) {
	vCol = texture(uni_sampler2D_Tex2D, var_vec2_Tex2D).rgb;
}
