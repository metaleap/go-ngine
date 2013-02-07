#if defined(VX_Uni_Matrix) && defined(VX_Att_Pos)
void vx_Pos_Matrix (inout vec4 vPos) {
	vPos = uni_VertexMatrix * vec4(att_Pos, 1.0);
}
#endif

#ifdef VX_Att_TexCoord
void vx_TexCoord_Attr (inout vec2 vTexCoord) {
	vTexCoord = att_Tex0;
}
#endif

#ifdef VX_Quad
void vx_Pos_Quad (inout vec4 vPos) {
	const float extent = 3;
	const vec2 pos[3] = vec2[](vec2(-1, -1), vec2(extent, -1), vec2(-1, extent));
	vPos = vec4(pos[gl_VertexID], 0, 1);
}

void vx_TexCoord_Quad (inout vec2 vTexCoord) {
	const float extent = 3;
	const vec2 vTex[3] = vec2[](vec2(0, 0), vec2(extent - 1, 0), vec2(0, extent - 1));
	vTexCoord = vTex[gl_VertexID];
}
#endif

void vx_Pos_Main (inout vec4 vPos) {
#if defined(VX_Uni_Matrix) && defined(VX_Att_Pos)
	vx_Pos_Matrix(vPos);
#endif

#ifdef VX_Quad
	vx_Pos_Quad(vPos);
#endif
}

void vx_TexCoord_Main (inout vec2 vTex) {
#ifdef VX_Att_TexCoord
	vx_TexCoord_Attr(vTex);
#endif

#ifdef VX_Quad
	vx_TexCoord_Quad(vTex);
#endif
}
