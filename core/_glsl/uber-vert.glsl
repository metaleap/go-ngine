/*
GLSL functions used for vertex shader permutation.
This file is "somewhat parsed" and processed.
*/

vec4 vx_Quad_gl_Position () {
	const float extent = 3;
	const vec2 pos[3] = vec2[](vec2(-1, -1), vec2(extent, -1), vec2(-1, extent));
	return vec4(pos[gl_VertexID], 0, 1);
}

vec2 vx_Quad_var_vec2_Tex2D () {
	const float extent = 3;
	const vec2 vTex[3] = vec2[](vec2(0, 0), vec2(extent - 1, 0), vec2(0, extent - 1));
	return vTex[gl_VertexID];
}

vec4 vx_Scene_gl_Position () {
	return uni_mat4_VertexMatrix * vec4(att_vec3_Pos, 1.0);
}

vec2 vx_Scene_var_vec2_Tex2D () {
	return att_vec2_Tex2D;
}

/*
	void vx_Main_Pos (inout vec4 vPos) {
		vPos = vx_Pos_Matrix();
		vPos = vx_Pos_Quad();
	}

	void vx_Main_TexCoord (inout vec2 vTex) {
		vTex = vx_TexCoord_Att();
		vTex = vx_TexCoord_Quad();
	}
*/
