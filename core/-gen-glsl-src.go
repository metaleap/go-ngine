package core

func init() {
	glc.progMan.Reset()
	rs := &glc.progMan.RawSources
	rs.Fragment["postfx"] = "in vec2 vTexCoords;\n\n#ifdef PostFx_TextureRect\nuniform sampler2DRect uTexRendering;\n#else\nuniform sampler2D uTexRendering;\n#endif\n\nout vec3 oColor;\n\nvoid Grayscale() {\n\toColor.rgb = vec3((oColor.r * 0.3) + (oColor.g * 0.59) + (oColor.b * 0.11));\n}\n\nvoid TestEffect() {\n\toColor.r = 1;\n}\n\nvoid main() {\n\n#ifdef PostFx_TextureRect\n\toColor = texture(uTexRendering, gl_FragCoord.xy).rgb;\n#else\n\toColor = texture(uTexRendering, vTexCoords).rgb; \n#endif\n\n#ifdef FX_Grayscale\n\tGrayscale();\n#endif\n\n}\n"
	rs.Fragment["rt_unlit3"] = "in vec2 vTexCoords;\n\nuniform sampler2D uDiffuse;\n\nout vec3 oColor;\n\nvoid main () {\n\toColor = texture(uDiffuse, vTexCoords).rgb;\n}\n"
	rs.Vertex["postfx"] = "const float extent = 3;\nconst vec2 vPos[3] = vec2[](vec2(-1, -1), vec2(extent, -1), vec2(-1, extent));\nconst vec2 vTex[3] = vec2[](vec2(0, 0), vec2(extent - 1, 0), vec2(0, extent - 1));\n\nout vec2 vTexCoords;\n\nvoid main () {\n\tgl_Position = vec4(vPos[gl_VertexID], 0, 1);\n\tvTexCoords = vTex[gl_VertexID];\n}\n"
	rs.Vertex["rt_unlit3"] = "in vec3 aPos;\nin vec2 aTexCoords;\n\nuniform mat4 uMatModelProj;\n\nout vec2 vTexCoords;\n\nvoid main () {\n\tgl_Position = uMatModelProj * vec4(aPos, 1.0);\n\tvTexCoords = aTexCoords;\n}\n"
	glc.progMan.Names = []string{"postfx", "rt_unlit3"}
}
