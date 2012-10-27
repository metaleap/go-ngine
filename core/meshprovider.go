package core

// import (
// 	gl "github.com/chsc/gogl/gl42"
// )

type FMeshProvider func (args ... interface {}) (*TMesh, error)

type tMeshProviders struct {
	PrefabCube, PrefabPlane, PrefabPyramid, PrefabQuad, PrefabTri FMeshProvider
}

var (
	MeshProviders = &tMeshProviders { meshProviderPrefabCube, meshProviderPrefabPlane, meshProviderPrefabPyramid, meshProviderPrefabQuad, meshProviderPrefabTri }
)

func meshProviderPrefabCube (args ... interface {}) (mesh *TMesh, err error) {
	var raw = newMeshData()
	mesh = Core.Meshes.New()
	raw.addPositions(
		tVa3 { -1, -1, 1 }, tVa3 { 1, -1, 1 }, tVa3 { 1, 1, 1 },
		tVa3 { -1, 1, 1 }, tVa3 { -1, -1, -1 }, tVa3 { -1, 1, -1 },
		tVa3 { 1, 1, -1 }, tVa3 { 1, -1, -1 })
	raw.addTexCoords(tVa2 { 0, 0 }, tVa2 { 1, 0 }, tVa2 { 1, 1 }, tVa2 { 0, 1 })
	raw.addNormals(tVa3 { 0, 0, 1 }, tVa3 { 0, 0, -1 }, tVa3 { 0, 1, 0 }, tVa3 { 0, -1, 0 }, tVa3 { 1, 0, 0 }, tVa3 { -1, 0, 0 })
	raw.addFaces(
		tMeshFace3 { tVe { 0, 0, 0 }, tVe { 1, 1, 0 }, tVe { 2, 2, 0 } },	tMeshFace3 { tVe { 0, 0, 0 }, tVe { 2, 2, 0 }, tVe { 3, 3, 0 } },		//	front
		tMeshFace3 { tVe { 4, 0, 1 }, tVe { 5, 1, 1 }, tVe { 6, 2, 1 } },	tMeshFace3 { tVe { 4, 0, 1 }, tVe { 6, 2, 1 }, tVe { 7, 3, 1 } },		//	back
		tMeshFace3 { tVe { 5, 0, 2 }, tVe { 3, 1, 2 }, tVe { 2, 2, 2 } },	tMeshFace3 { tVe { 5, 0, 2 }, tVe { 2, 2, 2 }, tVe { 6, 3, 2 } },		//	top
		tMeshFace3 { tVe { 4, 0, 3 }, tVe { 7, 1, 3 }, tVe { 1, 2, 3 } },	tMeshFace3 { tVe { 4, 0, 3 }, tVe { 1, 2, 3 }, tVe { 0, 3, 3 } },		//	bottom
		tMeshFace3 { tVe { 7, 0, 4 }, tVe { 6, 1, 4 }, tVe { 2, 2, 4 } },	tMeshFace3 { tVe { 7, 0, 4 }, tVe { 2, 2, 4 }, tVe { 1, 3, 4 } },		//	right
		tMeshFace3 { tVe { 4, 0, 5 }, tVe { 0, 1, 5 }, tVe { 3, 2, 5 } },	tMeshFace3 { tVe { 4, 0, 5 }, tVe { 3, 2, 5 }, tVe { 5, 3, 5 } },		//	left
		)
	mesh.load(raw)
	return
}

func meshProviderPrefabPlane (args ... interface {}) (mesh *TMesh, err error) {
	var raw = newMeshData()
	mesh = Core.Meshes.New()
	raw.addPositions(tVa3 { -1, 0, 1 }, tVa3 { 1, 0, 1 }, tVa3 { -1, 0, -1 }, tVa3 { 1, 0, -1 })
	raw.addTexCoords(tVa2 { 0, 0 }, tVa2 { 10, 0 }, tVa2 { 0, 10 }, tVa2 { 10, 10 })
	raw.addNormals(tVa3 { 0, 1, 0 })
	raw.addFaces(tMeshFace3 { tVe { 0, 0, 0 }, tVe { 1, 1, 0 }, tVe { 2, 2, 0 } },
		tMeshFace3 { tVe { 3, 3, 0 }, tVe { 2, 2, 0 }, tVe { 1, 1, 0 } })
	mesh.load(raw)
	return
}

func meshProviderPrefabPyramid (args ... interface {}) (mesh *TMesh, err error) {
	var raw = newMeshData()

	mesh = Core.Meshes.New()
	raw.addPositions(tVa3 { 0, 1, 0 }, tVa3 { -1, -1, 1 }, tVa3 { 1, -1, 1 }, tVa3 { 1, -1, -1 }, tVa3 { -1, -1, -1 })
	raw.addTexCoords(tVa2 { 0, 0 }, tVa2 { 1, 0 }, tVa2 { 1, 1 }, tVa2 { 0, 1})
	raw.addNormals(tVa3 { 0, 0, 1 }, tVa3 { 1, 0, 0 }, tVa3 { 0, 0, -1 }, tVa3 { -1, 0, 0 })
	raw.addFaces(
		tMeshFace3 { tVe { 0, 0, 0 }, tVe { 1, 1, 0 }, tVe { 2, 2, 0 } },
		tMeshFace3 { tVe { 0, 1, 1 }, tVe { 2, 2, 1 }, tVe { 3, 3, 1 } },
		tMeshFace3 { tVe { 0, 1, 2 }, tVe { 3, 2, 2 }, tVe { 4, 3, 2 } },
		tMeshFace3 { tVe { 0, 0, 3 }, tVe { 4, 1, 3 }, tVe { 1, 2, 3 } },
		)
	mesh.load(raw)
	return
}

func meshProviderPrefabQuad (args ... interface {}) (mesh *TMesh, err error) {
	var raw = newMeshData()
	mesh = Core.Meshes.New()
	raw.addPositions(tVa3 { 1, 1, 0 }, tVa3 { -1, 1, 0 }, tVa3 { -1, -1, 0 }, tVa3 { 1, -1, 0 })
	raw.addTexCoords(tVa2 { -0.125, 0 }, tVa2 { -0.125, 3 }, tVa2 { 1.125, 3 }, tVa2 { 1.125, 0 })
	raw.addNormals(tVa3 { 0, 0, 1 })
	raw.addFaces(tMeshFace3 { tVe { 0, 0, 0 }, tVe { 1, 1, 0 }, tVe { 2, 2, 0 } },
		tMeshFace3 { tVe { 0, 0, 0 }, tVe { 2, 2, 0 }, tVe { 3, 3, 0 } })
	mesh.load(raw)
	return
}

func meshProviderPrefabTri (args ... interface {}) (mesh *TMesh, err error) {
	var raw = newMeshData()
	raw.addPositions(tVa3 { 0, 1, 0 }, tVa3 { -1, -1, 0 }, tVa3 { 1, -1, 0 })
	raw.addTexCoords(tVa2 { 0, 0 }, tVa2 { 3, 0 }, tVa2 { 3, 2 })
	raw.addNormals(tVa3 { 0, 0, 1 })
	raw.addFaces(tMeshFace3 { tVe { 0, 0, 0 }, tVe { 1, 1, 0 }, tVe { 2, 2, 0 } })
	mesh = Core.Meshes.New()
	mesh.load(raw)
	return
}
