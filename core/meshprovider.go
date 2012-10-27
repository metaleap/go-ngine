package core

type FMeshProvider func (args ... interface {}) (*tMeshData, error)

type tMeshProviders struct {
	PrefabCube, PrefabPlane, PrefabPyramid, PrefabQuad, PrefabTri FMeshProvider
}

var (
	MeshProviders = &tMeshProviders { meshProviderPrefabCube, meshProviderPrefabPlane, meshProviderPrefabPyramid, meshProviderPrefabQuad, meshProviderPrefabTri }
)

func meshProviderPrefabCube (args ... interface {}) (meshData *tMeshData, err error) {
	meshData = newMeshData()
	meshData.addPositions(
		tVa3 { -1, -1, 1 }, tVa3 { 1, -1, 1 }, tVa3 { 1, 1, 1 },
		tVa3 { -1, 1, 1 }, tVa3 { -1, -1, -1 }, tVa3 { -1, 1, -1 },
		tVa3 { 1, 1, -1 }, tVa3 { 1, -1, -1 })
	meshData.addTexCoords(tVa2 { 0, 0 }, tVa2 { 1, 0 }, tVa2 { 1, 1 }, tVa2 { 0, 1 })
	meshData.addNormals(tVa3 { 0, 0, 1 }, tVa3 { 0, 0, -1 }, tVa3 { 0, 1, 0 }, tVa3 { 0, -1, 0 }, tVa3 { 1, 0, 0 }, tVa3 { -1, 0, 0 })
	meshData.addFaces(
		tMeshFace3 { tVe { 0, 0, 0 }, tVe { 1, 1, 0 }, tVe { 2, 2, 0 } },	tMeshFace3 { tVe { 0, 0, 0 }, tVe { 2, 2, 0 }, tVe { 3, 3, 0 } },		//	front
		tMeshFace3 { tVe { 4, 0, 1 }, tVe { 5, 1, 1 }, tVe { 6, 2, 1 } },	tMeshFace3 { tVe { 4, 0, 1 }, tVe { 6, 2, 1 }, tVe { 7, 3, 1 } },		//	back
		tMeshFace3 { tVe { 5, 0, 2 }, tVe { 3, 1, 2 }, tVe { 2, 2, 2 } },	tMeshFace3 { tVe { 5, 0, 2 }, tVe { 2, 2, 2 }, tVe { 6, 3, 2 } },		//	top
		tMeshFace3 { tVe { 4, 0, 3 }, tVe { 7, 1, 3 }, tVe { 1, 2, 3 } },	tMeshFace3 { tVe { 4, 0, 3 }, tVe { 1, 2, 3 }, tVe { 0, 3, 3 } },		//	bottom
		tMeshFace3 { tVe { 7, 0, 4 }, tVe { 6, 1, 4 }, tVe { 2, 2, 4 } },	tMeshFace3 { tVe { 7, 0, 4 }, tVe { 2, 2, 4 }, tVe { 1, 3, 4 } },		//	right
		tMeshFace3 { tVe { 4, 0, 5 }, tVe { 0, 1, 5 }, tVe { 3, 2, 5 } },	tMeshFace3 { tVe { 4, 0, 5 }, tVe { 3, 2, 5 }, tVe { 5, 3, 5 } })		//	left
	return
}

func meshProviderPrefabPlane (args ... interface {}) (meshData *tMeshData, err error) {
	meshData = newMeshData()
	meshData.addPositions(tVa3 { -1, 0, 1 }, tVa3 { 1, 0, 1 }, tVa3 { -1, 0, -1 }, tVa3 { 1, 0, -1 })
	meshData.addTexCoords(tVa2 { 0, 0 }, tVa2 { 30, 0 }, tVa2 { 0, 30 }, tVa2 { 30, 30 })
	meshData.addNormals(tVa3 { 0, 1, 0 })
	meshData.addFaces(
		tMeshFace3 { tVe { 0, 0, 0 }, tVe { 1, 1, 0 }, tVe { 2, 2, 0 } },
		tMeshFace3 { tVe { 3, 3, 0 }, tVe { 2, 2, 0 }, tVe { 1, 1, 0 } })
	return
}

func meshProviderPrefabPyramid (args ... interface {}) (meshData *tMeshData, err error) {
	meshData = newMeshData()
	meshData.addPositions(tVa3 { 0, 1, 0 }, tVa3 { -1, -1, 1 }, tVa3 { 1, -1, 1 }, tVa3 { 1, -1, -1 }, tVa3 { -1, -1, -1 })
	meshData.addTexCoords(tVa2 { 0, 0 }, tVa2 { 1, 0 }, tVa2 { 1, 1 }, tVa2 { 0, 1})
	meshData.addNormals(tVa3 { 0, 0, 1 }, tVa3 { 1, 0, 0 }, tVa3 { 0, 0, -1 }, tVa3 { -1, 0, 0 })
	meshData.addFaces(
		tMeshFace3 { tVe { 0, 0, 0 }, tVe { 1, 1, 0 }, tVe { 2, 2, 0 } },
		tMeshFace3 { tVe { 0, 1, 1 }, tVe { 2, 2, 1 }, tVe { 3, 3, 1 } },
		tMeshFace3 { tVe { 0, 1, 2 }, tVe { 3, 2, 2 }, tVe { 4, 3, 2 } },
		tMeshFace3 { tVe { 0, 0, 3 }, tVe { 4, 1, 3 }, tVe { 1, 2, 3 } })
	return
}

func meshProviderPrefabQuad (args ... interface {}) (meshData *tMeshData, err error) {
	meshData = newMeshData()
	meshData.addPositions(tVa3 { 1, 1, 0 }, tVa3 { -1, 1, 0 }, tVa3 { -1, -1, 0 }, tVa3 { 1, -1, 0 })
	meshData.addTexCoords(tVa2 { -0.125, 0 }, tVa2 { -0.125, 3 }, tVa2 { 1.125, 3 }, tVa2 { 1.125, 0 })
	meshData.addNormals(tVa3 { 0, 0, 1 })
	meshData.addFaces(
		tMeshFace3 { tVe { 0, 0, 0 }, tVe { 1, 1, 0 }, tVe { 2, 2, 0 } },
		tMeshFace3 { tVe { 0, 0, 0 }, tVe { 2, 2, 0 }, tVe { 3, 3, 0 } })
	return
}

func meshProviderPrefabTri (args ... interface {}) (meshData *tMeshData, err error) {
	meshData = newMeshData()
	meshData.addPositions(tVa3 { 0, 1, 0 }, tVa3 { -1, -1, 0 }, tVa3 { 1, -1, 0 })
	meshData.addTexCoords(tVa2 { 0, 0 }, tVa2 { 3, 0 }, tVa2 { 3, 2 })
	meshData.addNormals(tVa3 { 0, 0, 1 })
	meshData.addFaces(tMeshFace3 { tVe { 0, 0, 0 }, tVe { 1, 1, 0 }, tVe { 2, 2, 0 } })
	return
}
