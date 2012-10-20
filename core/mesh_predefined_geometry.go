package core

import (
	gl "github.com/chsc/gogl/gl42"
)

func (me *TMesh) SetVertsCube () {
	me.Verts = []gl.Float {
		// Front face
		-1, -1, 1,
		1, -1, 1,
		1, 1, 1,
		-1, 1, 1,
		// Back face
		-1, -1, -1,
		-1, 1, -1,
		1, 1, -1,
		1, -1, -1,
		// Top face
		-1, 1, -1,
		-1, 1, 1,
		1, 1, 1,
		1, 1, -1,
		// Bottom face
		-1, -1, -1,
		1, -1, -1,
		1, -1, 1,
		-1, -1, 1,
		// Right face
		1, -1, -1,
		1, 1, -1,
		1, 1, 1,
		1, -1, 1,
		// Left face
		-1, -1, -1,
		-1, -1, 1,
		-1, 1, 1,
		-1, 1, -1,
	}
	me.Normals = []gl.Float {
		// Front face
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		// Back face
		0, 0, -1,
		0, 0, -1,
		0, 0, -1,
		0, 0, -1,
		// Top face
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		// Bottom face
		0, -1, 0,
		0, -1, 0,
		0, -1, 0,
		0, -1, 0,
		// Right face
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		// Left face
		-1, 0, 0,
		-1, 0, 0,
		-1, 0, 0,
		-1, 0, 0,
	}
	me.Indices = []gl.Uint {
		0, 1, 2, 0, 2, 3, // Front face
		4, 5, 6, 4, 6, 7, // Back face
		8, 9, 10, 8, 10, 11, // Top face
		12, 13, 14, 12, 14, 15, // Bottom face
		16, 17, 18, 16, 18, 19, // Right face
		20, 21, 22, 20, 22, 23, // Left face
	}
	me.glMode, me.glNumIndices, me.glNumVerts = gl.TRIANGLES, 36, 24
}

func (me *TMesh) SetVertsFace3 () {
	me.Indices = nil
	me.Verts = []gl.Float {
		0, 1, 0,
		-1, -1, 0,
		1, -1, 0,
	}
	me.glMode, me.glNumVerts = gl.TRIANGLES, 3
}

func (me *TMesh) SetVertsFace4 () {
	me.Indices = nil
	me.Verts = []gl.Float {
		1, 1, 0,
		-1, 1, 0,
		1, -1, 0,
		-1, -1, 0,
	}
	me.glMode, me.glNumVerts = gl.TRIANGLE_STRIP, 4
}

func (me *TMesh) SetVertsPlane () {
	me.Indices = nil
	me.Verts = []gl.Float {
		1, 0, -1,
		-1, 0, -1,
		1, -0, 1,
		-1, -0, 1,
	}
	me.Normals = []gl.Float {
		0, 10, 0,
		0, 10, 0,
		0, 10, 0,
		0, 10, 0,
	}
	me.glMode, me.glNumVerts = gl.TRIANGLE_STRIP, 4
}

func (me *TMesh) SetVertsPyramid () {
	me.Indices = nil
	me.Verts = []gl.Float {
		// Front face
		0, 1, 0,
		-1, -1, 1,
		1, -1, 1,
		// Right face
		0, 1, 0,
		1, -1, 1,
		1, -1, -1,
		// Back face
		0, 1, 0,
		1, -1, -1,
		-1, -1, -1,
		// Left face
		0, 1, 0,
		-1, -1, -1,
		-1, -1, 1,
	}
	me.Normals = []gl.Float {
		// Front face
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		// Back face
		0, 0, -1,
		0, 0, -1,
		0, 0, -1,
		// Right face
		1, 0, 0,
		1, 0, 0,
		1, 0, 0,
		// Left face
		-1, 0, 0,
		-1, 0, 0,
		-1, 0, 0,
	}
	me.glMode, me.glNumVerts = gl.TRIANGLES, 12
}
