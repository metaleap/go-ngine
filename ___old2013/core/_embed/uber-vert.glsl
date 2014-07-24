/*

GLSL functions used for vertex shader permutation.
This file is "somewhat parsed" and processed, so indentation and
naming patterns are significant and not subject to personal taste.

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
	vec4 pos = uni_mat4_VertexMatrix * vec4(att_vec3_Pos, 1.0);
	if (uni_int_Sky == 1) {
		pos.z = pos.w;
	}
	return pos;
}

vec2 vx_Scene_var_vec2_Tex2D () {
	return att_vec2_Tex2D;
}

vec3 vx_Scene_var_vec3_TexCube () {
	return att_vec3_Pos;
}
