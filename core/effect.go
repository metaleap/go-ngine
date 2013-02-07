package core

/*

grayscale:

	func fx_grayscale (inout vec3 col)
		col.rgb = graymaths(col)


color:
	uniform vec3 uColN;

	func fx_color (inout vec3 col)
		col.rgb = uColN

texture:
	in vec2 vTexCoordN
	uniform sampler2D uTexN

	func fx_texture( inout vec3 col)
		col.rgb = texture(uTexN, vTexCoordN).rgb;


*/

type Effect struct {
}
