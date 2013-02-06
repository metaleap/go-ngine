package core

func init() {
	glc.progMan.Reset()
	rs := &glc.progMan.RawSources
	rs.Fragment["rt_quad"] = "in vec2 vTexCoords;\n\nuniform sampler2D uTexRendering;\n\nout vec3 oColor;\n\nvoid Grayscale() {\n\toColor.rgb = vec3((oColor.r * 0.3) + (oColor.g * 0.59) + (oColor.b * 0.11));\n}\n\nvoid TestEffect() {\n\toColor.r = 1;\n}\n\nvoid main() {\n\toColor = texture(uTexRendering, vTexCoords).rgb;\n\n#ifdef FX_Grayscale\n\tGrayscale();\n#endif\n}\n"
	rs.Fragment["rt_unlit"] = "in vec2 vTexCoords;\n\nuniform sampler2D uDiffuse;\n\nout vec3 oColor;\n\nvoid main () {\n\toColor = texture(uDiffuse, vTexCoords).rgb;\n}\n"
	rs.Vertex["rt_quad"] = "const float extent = 3;\n\nconst vec2 vPos[3] = vec2[](vec2(-1, -1), vec2(extent, -1), vec2(-1, extent));\nconst vec2 vTex[3] = vec2[](vec2(0, 0), vec2(extent - 1, 0), vec2(0, extent - 1));\n\nout vec2 vTexCoords;\n\nvoid main () {\n\tgl_Position = vec4(vPos[gl_VertexID], 0, 1);\n\tvTexCoords = vTex[gl_VertexID];\n}\n"
	rs.Vertex["rt_unlit"] = "in vec3 aPos;\nin vec2 aTexCoords;\n\nuniform mat4 uMatModelProj;\n\nout vec2 vTexCoords;\n\nvoid main () {\n\tgl_Position = uMatModelProj * vec4(aPos, 1.0);\n\tvTexCoords = aTexCoords;\n}\n"
	glc.progMan.Names = []string{"rt_quad", "rt_unlit"}
}
