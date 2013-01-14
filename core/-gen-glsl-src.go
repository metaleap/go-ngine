package core

func init() {
	var rss glShaderSources
	rss.init()
	rss.vertex["rt_unlit"] = "in vec3 aPos;\nin vec2 aTexCoords;\n\nuniform mat4 uMatModelProj;\n\nout vec2 vTexCoords;\n\nvoid main () {\n\tgl_Position = uMatModelProj * vec4(aPos, 1.0);\n\tvTexCoords = aTexCoords;\n}\n"
	rss.fragment["rt_unlit"] = "in vec2 vTexCoords;\n\nuniform sampler2D uDiffuse;\n\nout vec4 oColor;\n\nvoid main () {\n\toColor = texture(uDiffuse, vTexCoords);\n}\n"
	glShaderMan.sources = rss
	glShaderMan.names = []string{"rt_unlit"}
}
