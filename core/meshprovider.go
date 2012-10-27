package core

import (
	gl "github.com/chsc/gogl/gl42"
)

type FMeshProvider func (args ... interface {}) (*TMesh, error)

type tMeshProviders struct {
	PrefabCube, PrefabPlane, PrefabPyramid, PrefabQuad, PrefabTri FMeshProvider
}

var (
	MeshProviders = &tMeshProviders { meshProviderPrefabCube, meshProviderPrefabPlane, meshProviderPrefabPyramid, meshProviderPrefabQuad, meshProviderPrefabTri }
)

func meshProviderPrefabCube (args ... interface {}) (mesh *TMesh, err error) {
	mesh = Core.Meshes.New()
	mesh.Verts = []gl.Float {
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
	mesh.Indices = []gl.Uint {
		0, 1, 2, 0, 2, 3, // Front face
		4, 5, 6, 4, 6, 7, // Back face
		8, 9, 10, 8, 10, 11, // Top face
		12, 13, 14, 12, 14, 15, // Bottom face
		16, 17, 18, 16, 18, 19, // Right face
		20, 21, 22, 20, 22, 23, // Left face
	}
	mesh.Normals = []gl.Float {
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
	mesh.TexCoords = []gl.Float {
		// Front face
		0, 0,
		1, 0,
		1, 1,
		0, 1,
		// Back face
		1, 0,
		1, 1,
		0, 1,
		0, 0,
		// Top face
		0, 1,
		0, 0,
		1, 0,
		1, 1,
		// Bottom face
		1, 1,
		0, 1,
		0, 0,
		1, 0,
		// Right face
		1, 0,
		1, 1,
		0, 1,
		0, 0,
		// Left face
		0, 0,
		1, 0,
		1, 1,
		0, 1,
	}
	mesh.glMode, mesh.glNumIndices, mesh.glNumVerts = gl.TRIANGLES, 36, 24
	return
}

func meshProviderPrefabPlane (args ... interface {}) (mesh *TMesh, err error) {
	mesh = Core.Meshes.New()
	mesh.Verts = []gl.Float {
		-1, 0, 1,
		1, 0, 1,
		-1, 0, -1,
		1, 0, -1,
	}
	mesh.Indices = []gl.Uint { 0, 1, 2, 3, 2, 1 }
	mesh.Normals = []gl.Float {
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
		0, 1, 0,
	}
	mesh.TexCoords = []gl.Float {
		0, 0,
		10, 0,
		0, 10,
		10, 10,
	}
	mesh.glMode, mesh.glNumIndices, mesh.glNumVerts = gl.TRIANGLES, 6, 4
	return
}

func meshProviderPrefabPyramid (args ... interface {}) (mesh *TMesh, err error) {
	mesh = Core.Meshes.New()
	mesh.Verts = []gl.Float {
		// Front face
		0, 1, 0,
		-1, -1, 1,
		1, -1, 1,
		// Right face
		// 0, 1, 0,
		// 1, -1, 1,
		1, -1, -1,
		// Back face
		// 0, 1, 0,
		// 1, -1, -1,
		-1, -1, -1,
		// Left face
		// 0, 1, 0,
		// -1, -1, -1,
		// -1, -1, 1,
	}
	mesh.Indices = []gl.Uint {
		0, 1, 2, // Front face
		0, 2, 3, // Right face
		0, 3, 4, // Back face
		0, 4, 1, // Left face
	}
	mesh.Normals = []gl.Float {
		// Front face
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		// Right face
		// 1, 0, 0,
		// 1, 0, 0,
		1, 0, 0,
		// Back face
		// 0, 0, -1,
		// 0, 0, -1,
		0, 0, -1,
		// Left face
		// -1, 0, 0,
		// -1, 0, 0,
		// -1, 0, 0,
	}
	mesh.TexCoords = []gl.Float {
		// Front face
		0, 0,
		1, 0,
		1, 1,
		// Right face
		1, 0,
		1, 1,
		0, 1,
		// Back face
		1, 0,
		1, 1,
		0, 1,
		// Left face
		0, 0,
		1, 0,
		1, 1,
	}
	mesh.glMode, mesh.glNumIndices, mesh.glNumVerts = gl.TRIANGLES, 12, 5
	return
}

func meshProviderPrefabQuad (args ... interface {}) (mesh *TMesh, err error) {
	mesh = Core.Meshes.New()
	mesh.raw = newMeshData()
	mesh.raw.addPositions()
	mesh.raw.addTexCoords()
	mesh.raw.addNormals()
	mesh.raw.addFace(tMeshFace3 {})
	mesh.Verts = []gl.Float {
		1, 1, 0,
		-1, 1, 0,
		-1, -1, 0,
		1, -1, 0,
	}
	mesh.Indices = []gl.Uint { 0, 1, 2, 0, 2, 3 }
	mesh.Normals = []gl.Float {
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
	}
	mesh.TexCoords = []gl.Float {
		0, 0,
		0, 1,
		1, 1,
		1, 0,
	}
	mesh.glMode, mesh.glNumIndices, mesh.glNumVerts = gl.TRIANGLES, 6, 4
	return
}

func meshProviderPrefabTri (args ... interface {}) (mesh *TMesh, err error) {
	mesh = Core.Meshes.New()
	mesh.raw = newMeshData()
	mesh.raw.addPositions(tVa3 { 0, 1, 0 }, tVa3 { -1, -1, 0 }, tVa3 { 1, -1, 0 })
	mesh.raw.addTexCoords(tVa2 { 0, 0 }, tVa2 { 4, 0 }, tVa2 { 2, 2 })
	mesh.raw.addNormals(tVa3 { 0, 0, 1 })
	mesh.raw.addFace(tMeshFace3 { tVe { 0, 0, 0 }, tVe { 1, 1, 0 }, tVe { 2, 2, 0 } })
	mesh.Verts = []gl.Float {
		0, 1, 0,
		-1, -1, 0,
		1, -1, 0,
	}
	mesh.Normals = []gl.Float {
		0, 0, 1,
		0, 0, 1,
		0, 0, 1,
	}
	mesh.TexCoords = []gl.Float {
		0, 0,
		3, 0,
		3, 3,
	}
	mesh.Indices = []gl.Uint { 0, 1, 2 }
	mesh.glMode, mesh.glNumIndices, mesh.glNumVerts = gl.TRIANGLES, 3, 3
	return
}
