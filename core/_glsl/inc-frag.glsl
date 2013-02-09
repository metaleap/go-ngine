#ifdef FX_Grayscale
vec3 fx_Grayscale (const in vec3 vCol) {
	return vec3((vCol.r * 0.3) + (vCol.g * 0.59) + (vCol.b * 0.11));
}
#endif

#ifdef FX_RedTest
vec3 fx_RedTest (const in vec3 vCol) {
	return vec3(1, vCol.gb);
}
#endif

vec3 fx_Tex0 (const in vec3 vCol) {
	return texture(uni_Tex0, var_Tex0).rgb;
}

void fx_Main (inout vec3 vCol) {
	vCol = fx_Tex0(vCol);

#ifdef FX_Grayscale
	vCol = fx_Grayscale(vCol);
#endif

#ifdef FX_RedTest
	vCol = fx_RedTest(vCol);
#endif
}
