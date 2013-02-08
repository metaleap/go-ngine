/*

GLSL functions used for dynamic fragment shader composition.
This file is "parsed" so all functions accept one inout arg and return void.

*/

void fx_Grayscale (inout vec3 vCol) {
	vCol.rgb = vec3((vCol.r * 0.3) + (vCol.g * 0.59) + (vCol.b * 0.11));
}

void fx_RedTest (inout vec3 vCol) {
	vCol.r = 1;
}

void fx_Tex0 (inout vec3 vCol) {
	vCol = texture(uni_Tex0, var_Tex0).rgb;
}

void fx_Main (inout vec3 vCol) {
	fx_Tex0(vCol);
	fx_Grayscale(vCol);
	fx_RedTest(vCol);
}
