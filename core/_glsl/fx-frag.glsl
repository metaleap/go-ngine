/*

GLSL functions used for fragment shader permutation.
This file is "somewhat parsed" and processed.

*/

vec3 fx_Grayscale (const in vec3 vCol) {
	return vec3((vCol.r * 0.3) + (vCol.g * 0.59) + (vCol.b * 0.11));
}

vec3 fx_RedTest (const in vec3 vCol) {
	return vec3(1, vCol.gb);
}

vec3 fx_Tex0 (const in vec3 vCol) {
	return texture(uni_Tex0, var_Tex0).rgb;
}

void fx_Main (inout vec3 vCol) {
	vCol = fx_Tex0(vCol);
	vCol = fx_Grayscale(vCol);
	vCol = fx_RedTest(vCol);
}
