/*

GLSL functions used for vertex shader permutation.
This file is "somewhat parsed" and processed.

*/

vec4 vx_Pos_Matrix () {
	return uni_VertexMatrix * vec4(att_Pos, 1.0);
}

vec2 vx_TexCoord_Att () {
	return att_Tex0;
}

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

void vx_Main_Pos (inout vec4 vPos) {
	vPos = vx_Pos_Matrix();
	vPos = vx_Pos_Quad();
}

void vx_Main_TexCoord (inout vec2 vTex) {
	vTex = vx_TexCoord_Att();
	vTex = vx_TexCoord_Quad();
}
