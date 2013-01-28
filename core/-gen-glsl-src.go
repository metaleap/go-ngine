package core

func init() {
	glc.progMan.Reset()
	rs := &glc.progMan.RawSources
	rs.Fragment["postfx"] = "uniform sampler2DRect uTexRendering;\n\nout vec3 oColor;\n\nvoid Grayscale() {\n\toColor.rgb = vec3((oColor.r * 0.3) + (oColor.g * 0.59) + (oColor.b * 0.11));\n}\n\nvoid TestEffect() {\n\toColor.r = 1;\n}\n\nvoid main() {\n\toColor = texture(uTexRendering, gl_FragCoord.xy).rgb;\n\n#ifdef FX_Grayscale\n\tGrayscale();\n#endif\n\n}\n"
	rs.Fragment["rt_unlit3"] = "in vec2 vTexCoords;\n\nuniform sampler2D uDiffuse;\n\nout vec3 oColor;\n\nvoid main () {\n\toColor = texture(uDiffuse, vTexCoords).rgb;\n}\n"
	rs.Vertex["postfx"] = "const vec2 vPos[3] = vec2[](vec2(-1, -1), vec2(3, -1), vec2(-1, 3));\n\nvoid main () {\n\tgl_Position = vec4(vPos[gl_VertexID], 0, 1);\n}\n"
	rs.Vertex["rt_unlit3"] = "in vec3 aPos;\nin vec2 aTexCoords;\n\nuniform mat4 uMatModelProj;\n\nout vec2 vTexCoords;\n\nvoid main () {\n\tgl_Position = uMatModelProj * vec4(aPos, 1.0);\n\tvTexCoords = aTexCoords;\n}\n"
	glc.progMan.Names = []string{"postfx", "rt_unlit3"}
}
