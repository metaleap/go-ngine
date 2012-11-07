package core

type MeshProvider func (args ... interface {}) (*meshData, error)

type meshProviders struct {
	PrefabCube, PrefabPlane, PrefabPyramid, PrefabQuad, PrefabTri MeshProvider
}

var (
	MeshProviders = &meshProviders { meshProviderPrefabCube, meshProviderPrefabPlane, meshProviderPrefabPyramid, meshProviderPrefabQuad, meshProviderPrefabTri }
)

func meshProviderPrefabCube (args ... interface {}) (meshData *meshData, err error) {
	meshData = newMeshData()
	meshData.addPositions(
		va3 { -1, -1, 1 }, va3 { 1, -1, 1 }, va3 { 1, 1, 1 },
		va3 { -1, 1, 1 }, va3 { -1, -1, -1 }, va3 { -1, 1, -1 },
		va3 { 1, 1, -1 }, va3 { 1, -1, -1 })
	meshData.addTexCoords(va2 { 0, 0 }, va2 { 1, 0 }, va2 { 1, 1 }, va2 { 0, 1 })
	meshData.addNormals(va3 { 0, 0, 1 }, va3 { 0, 0, -1 }, va3 { 0, 1, 0 }, va3 { 0, -1, 0 }, va3 { 1, 0, 0 }, va3 { -1, 0, 0 })
	meshData.addFaces(
		meshFace3 { ve { 0, 0, 0 }, ve { 1, 1, 0 }, ve { 2, 2, 0 } },	meshFace3 { ve { 0, 0, 0 }, ve { 2, 2, 0 }, ve { 3, 3, 0 } },		//	front
		meshFace3 { ve { 4, 0, 1 }, ve { 5, 1, 1 }, ve { 6, 2, 1 } },	meshFace3 { ve { 4, 0, 1 }, ve { 6, 2, 1 }, ve { 7, 3, 1 } },		//	back
		meshFace3 { ve { 5, 0, 2 }, ve { 3, 1, 2 }, ve { 2, 2, 2 } },	meshFace3 { ve { 5, 0, 2 }, ve { 2, 2, 2 }, ve { 6, 3, 2 } },		//	top
		meshFace3 { ve { 4, 0, 3 }, ve { 7, 1, 3 }, ve { 1, 2, 3 } },	meshFace3 { ve { 4, 0, 3 }, ve { 1, 2, 3 }, ve { 0, 3, 3 } },		//	bottom
		meshFace3 { ve { 7, 0, 4 }, ve { 6, 1, 4 }, ve { 2, 2, 4 } },	meshFace3 { ve { 7, 0, 4 }, ve { 2, 2, 4 }, ve { 1, 3, 4 } },		//	right
		meshFace3 { ve { 4, 0, 5 }, ve { 0, 1, 5 }, ve { 3, 2, 5 } },	meshFace3 { ve { 4, 0, 5 }, ve { 3, 2, 5 }, ve { 5, 3, 5 } })		//	left
	return
}

func meshProviderPrefabPlane (args ... interface {}) (meshData *meshData, err error) {
	meshData = newMeshData()
	meshData.addPositions(va3 { -1, 0, 1 }, va3 { 1, 0, 1 }, va3 { -1, 0, -1 }, va3 { 1, 0, -1 })
	meshData.addTexCoords(va2 { 0, 0 }, va2 { 30, 0 }, va2 { 0, 30 }, va2 { 30, 30 })
	meshData.addNormals(va3 { 0, 1, 0 })
	meshData.addFaces(
		meshFace3 { ve { 0, 0, 0 }, ve { 1, 1, 0 }, ve { 2, 2, 0 } },
		meshFace3 { ve { 3, 3, 0 }, ve { 2, 2, 0 }, ve { 1, 1, 0 } })
	return
}

func meshProviderPrefabPyramid (args ... interface {}) (meshData *meshData, err error) {
	meshData = newMeshData()
	meshData.addPositions(va3 { 0, 1, 0 }, va3 { -1, -1, 1 }, va3 { 1, -1, 1 }, va3 { 1, -1, -1 }, va3 { -1, -1, -1 })
	meshData.addTexCoords(va2 { 0, 0 }, va2 { 1, 0 }, va2 { 1, 1 }, va2 { 0, 1})
	meshData.addNormals(va3 { 0, 0, 1 }, va3 { 1, 0, 0 }, va3 { 0, 0, -1 }, va3 { -1, 0, 0 })
	meshData.addFaces(
		meshFace3 { ve { 0, 0, 0 }, ve { 1, 1, 0 }, ve { 2, 2, 0 } },
		meshFace3 { ve { 0, 1, 1 }, ve { 2, 2, 1 }, ve { 3, 3, 1 } },
		meshFace3 { ve { 0, 1, 2 }, ve { 3, 2, 2 }, ve { 4, 3, 2 } },
		meshFace3 { ve { 0, 0, 3 }, ve { 4, 1, 3 }, ve { 1, 2, 3 } })
	return
}

func meshProviderPrefabQuad (args ... interface {}) (meshData *meshData, err error) {
	meshData = newMeshData()
	meshData.addPositions(va3 { 1, 1, 0 }, va3 { -1, 1, 0 }, va3 { -1, -1, 0 }, va3 { 1, -1, 0 })
	meshData.addTexCoords(va2 { -0.125, 0 }, va2 { -0.125, 3 }, va2 { 1.125, 3 }, va2 { 1.125, 0 })
	meshData.addNormals(va3 { 0, 0, 1 })
	meshData.addFaces(
		meshFace3 { ve { 0, 0, 0 }, ve { 1, 1, 0 }, ve { 2, 2, 0 } },
		meshFace3 { ve { 0, 0, 0 }, ve { 2, 2, 0 }, ve { 3, 3, 0 } })
	return
}

func meshProviderPrefabTri (args ... interface {}) (meshData *meshData, err error) {
	meshData = newMeshData()
	meshData.addPositions(va3 { 0, 1, 0 }, va3 { -1, -1, 0 }, va3 { 1, -1, 0 })
	meshData.addTexCoords(va2 { 0, 0 }, va2 { 3, 0 }, va2 { 3, 2 })
	meshData.addNormals(va3 { 0, 0, 1 })
	meshData.addFaces(meshFace3 { ve { 0, 0, 0 }, ve { 1, 1, 0 }, ve { 2, 2, 0 } })
	return
}
