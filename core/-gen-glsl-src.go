package core

func init() {
	glc.progMan.Reset()
	rs := &glc.progMan.RawSources
	rs.Vertex["rt_quad"] = "out vec2 var_Tex0;\n\nvoid vx_Pos_Quad (inout vec4 vPos) {\n\tconst float extent = 3;\n\tconst vec2 pos[3] = vec2[](vec2(-1, -1), vec2(extent, -1), vec2(-1, extent));\n\tvPos = vec4(pos[gl_VertexID], 0, 1);\n}\n\nvoid vx_TexCoord_Quad (inout vec2 vTexCoord) {\n\tconst float extent = 3;\n\tconst vec2 vTex[3] = vec2[](vec2(0, 0), vec2(extent - 1, 0), vec2(0, extent - 1));\n\tvTexCoord = vTex[gl_VertexID];\n}\n\nvoid main () {\n\tvx_Pos_Quad(gl_Position);\n\tvx_TexCoord_Quad(var_Tex0);\n}\n"
	rs.Vertex["rt_unlit"] = "in vec3 att_Pos;\nin vec2 att_Tex0;\n\nuniform mat4 uni_VertexMatrix;\n\nout vec2 var_Tex0;\n\nvoid vx_Pos_Matrix (inout vec4 vPos) {\n\tvPos = uni_VertexMatrix * vec4(att_Pos, 1.0);\n}\n\nvoid vx_TexCoord_Attr (inout vec2 vTexCoord) {\n\tvTexCoord = att_Tex0;\n}\n\nvoid main () {\n\tvx_Pos_Matrix(gl_Position);\n\tvx_TexCoord_Attr(var_Tex0);\n}\n"
	rs.Fragment["rt_quad"] = "in vec2 var_Tex0;\n\nuniform sampler2D uni_Tex0;\n\nout vec3 out_Color;\n\n#ifdef FX_Grayscale\nvoid fx_Grayscale (inout vec3 vCol) {\n\tvCol.rgb = vec3((vCol.r * 0.3) + (vCol.g * 0.59) + (vCol.b * 0.11));\n}\n#endif\n\n#ifdef FX_RedTest\nvoid fx_RedTest (inout vec3 vCol) {\n\tvCol.r = 1;\n}\n#endif\n\nvoid fx_Tex0 (inout vec3 vCol) {\n\tvCol = texture(uni_Tex0, var_Tex0).rgb;\n}\n\nvoid main() {\n\tfx_Tex0(out_Color);\n\n#ifdef FX_Grayscale\n\tfx_Grayscale(out_Color);\n#endif\n\n#ifdef FX_RedTest\n\tfx_RedTest(out_Color);\n#endif\n}\n"
	rs.Fragment["rt_unlit"] = "in vec2 var_Tex0;\n\nuniform sampler2D uni_Tex0;\n\nout vec3 out_Color;\n\nvoid fx_Tex0 (inout vec3 vCol) {\n\tvCol = texture(uni_Tex0, var_Tex0).rgb;\n}\n\nvoid main () {\n\tfx_Tex0(out_Color);\n}\n"
	glc.progMan.Names = []string{"rt_quad", "rt_unlit"}
}
