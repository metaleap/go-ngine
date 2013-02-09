#if defined(VX_Uni_Matrix) && defined(VX_Att_Pos)
vec4 vx_Pos_Matrix () {
	return uni_VertexMatrix * vec4(att_Pos, 1.0);
}
#endif

#ifdef VX_Att_TexCoord
vec2 vx_TexCoord_Att () {
	return att_Tex0;
}
#endif

#ifdef VX_Quad
vec4 vx_Pos_Quad () {
	const float extent = 3;
	const vec2 pos[3] = vec2[](vec2(-1, -1), vec2(extent, -1), vec2(-1, extent));
	return vec4(pos[gl_VertexID], 0, 1);
}

vec2 vx_TexCoord_Quad () {
	const float extent = 3;
	const vec2 vTex[3] = vec2[](vec2(0, 0), vec2(extent - 1, 0), vec2(0, extent - 1));
	return vTex[gl_VertexID];
}
#endif

void vx_Pos_Main (inout vec4 vPos) {
#if defined(VX_Uni_Matrix) && defined(VX_Att_Pos)
	vPos = vx_Pos_Matrix();
#endif

#ifdef VX_Quad
	vPos = vx_Pos_Quad();
#endif
}

void vx_TexCoord_Main (inout vec2 vTex) {
#ifdef VX_Att_TexCoord
	vTex = vx_TexCoord_Att();
#endif

#ifdef VX_Quad
	vTex = vx_TexCoord_Quad();
#endif
}
