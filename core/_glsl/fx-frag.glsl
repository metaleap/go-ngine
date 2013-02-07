#ifdef FX_Grayscale
void fx_Grayscale (inout vec3 vCol) {
	vCol.rgb = vec3((vCol.r * 0.3) + (vCol.g * 0.59) + (vCol.b * 0.11));
}
#endif

#ifdef FX_RedTest
void fx_RedTest (inout vec3 vCol) {
	vCol.r = 1;
}
#endif

void fx_Tex0 (inout vec3 vCol) {
	vCol = texture(uni_Tex0, var_Tex0).rgb;
}

void fx_Main (inout vec3 vCol) {
	fx_Tex0(vCol);

#ifdef FX_Grayscale
	fx_Grayscale(vCol);
#endif

#ifdef FX_RedTest
	fx_RedTest(vCol);
#endif
}
