/*

GLSL functions used for vertex shader permutation.
This file is "parsed" so all functions accept one inout arg and return void.
*/

void vx_Pos_Matrix (inout vec4 vPos) {
	vPos = uni_VertexMatrix * vec4(att_Pos, 1.0);
}

void vx_TexCoord_Att (inout vec2 vTexCoord) {
	vTexCoord = att_Tex0;
}

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

void vx_Pos_Main (inout vec4 vPos) {
	vx_Pos_Matrix(vPos);
	vx_Pos_Quad(vPos);
}

void vx_TexCoord_Main (inout vec2 vTex) {
	vx_TexCoord_Att(vTex);
	vx_TexCoord_Quad(vTex);
}
