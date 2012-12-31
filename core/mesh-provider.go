package core

type meshProviders struct {
	PrefabCube, PrefabPlane, PrefabPyramid, PrefabQuad, PrefabTri MeshProvider
}

var (
	//	A collection of all "mesh providers" known to go:ngine.
	MeshProviders = &meshProviders{meshProviderPrefabCube, meshProviderPrefabPlane, meshProviderPrefabPyramid, meshProviderPrefabQuad, meshProviderPrefabTri}
)

func meshProviderPrefabCube(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	meshData.addPositions(
		meshVertAtt3{-1, -1, 1}, meshVertAtt3{1, -1, 1}, meshVertAtt3{1, 1, 1},
		meshVertAtt3{-1, 1, 1}, meshVertAtt3{-1, -1, -1}, meshVertAtt3{-1, 1, -1},
		meshVertAtt3{1, 1, -1}, meshVertAtt3{1, -1, -1})
	meshData.addTexCoords(meshVertAtt2{0, 0}, meshVertAtt2{1, 0}, meshVertAtt2{1, 1}, meshVertAtt2{0, 1})
	meshData.addNormals(meshVertAtt3{0, 0, 1}, meshVertAtt3{0, 0, -1}, meshVertAtt3{0, 1, 0}, meshVertAtt3{0, -1, 0}, meshVertAtt3{1, 0, 0}, meshVertAtt3{-1, 0, 0})
	meshData.addFaces(
		meshFace3{meshVert{0, 0, 0}, meshVert{1, 1, 0}, meshVert{2, 2, 0}}, meshFace3{meshVert{0, 0, 0}, meshVert{2, 2, 0}, meshVert{3, 3, 0}}, //	front
		meshFace3{meshVert{4, 0, 1}, meshVert{5, 1, 1}, meshVert{6, 2, 1}}, meshFace3{meshVert{4, 0, 1}, meshVert{6, 2, 1}, meshVert{7, 3, 1}}, //	back
		meshFace3{meshVert{5, 0, 2}, meshVert{3, 1, 2}, meshVert{2, 2, 2}}, meshFace3{meshVert{5, 0, 2}, meshVert{2, 2, 2}, meshVert{6, 3, 2}}, //	top
		meshFace3{meshVert{4, 0, 3}, meshVert{7, 1, 3}, meshVert{1, 2, 3}}, meshFace3{meshVert{4, 0, 3}, meshVert{1, 2, 3}, meshVert{0, 3, 3}}, //	bottom
		meshFace3{meshVert{7, 0, 4}, meshVert{6, 1, 4}, meshVert{2, 2, 4}}, meshFace3{meshVert{7, 0, 4}, meshVert{2, 2, 4}, meshVert{1, 3, 4}}, //	right
		meshFace3{meshVert{4, 0, 5}, meshVert{0, 1, 5}, meshVert{3, 2, 5}}, meshFace3{meshVert{4, 0, 5}, meshVert{3, 2, 5}, meshVert{5, 3, 5}}) //	left
	return
}

func meshProviderPrefabPlane(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	meshData.addPositions(meshVertAtt3{-1, 0, 1}, meshVertAtt3{1, 0, 1}, meshVertAtt3{-1, 0, -1}, meshVertAtt3{1, 0, -1})
	meshData.addTexCoords(meshVertAtt2{0, 0}, meshVertAtt2{1000, 0}, meshVertAtt2{0, 1000}, meshVertAtt2{1000, 1000})
	meshData.addNormals(meshVertAtt3{0, 1, 0})
	meshData.addFaces(
		meshFace3{meshVert{0, 0, 0}, meshVert{1, 1, 0}, meshVert{2, 2, 0}},
		meshFace3{meshVert{3, 3, 0}, meshVert{2, 2, 0}, meshVert{1, 1, 0}})
	return
}

func meshProviderPrefabPyramid(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	meshData.addPositions(meshVertAtt3{0, 1, 0}, meshVertAtt3{-1, -1, 1}, meshVertAtt3{1, -1, 1}, meshVertAtt3{1, -1, -1}, meshVertAtt3{-1, -1, -1})
	meshData.addTexCoords(meshVertAtt2{0, 0}, meshVertAtt2{1, 0}, meshVertAtt2{1, 1}, meshVertAtt2{0, 1})
	meshData.addNormals(meshVertAtt3{0, 0, 1}, meshVertAtt3{1, 0, 0}, meshVertAtt3{0, 0, -1}, meshVertAtt3{-1, 0, 0})
	meshData.addFaces(
		meshFace3{meshVert{0, 0, 0}, meshVert{1, 1, 0}, meshVert{2, 2, 0}},
		meshFace3{meshVert{0, 1, 1}, meshVert{2, 2, 1}, meshVert{3, 3, 1}},
		meshFace3{meshVert{0, 1, 2}, meshVert{3, 2, 2}, meshVert{4, 3, 2}},
		meshFace3{meshVert{0, 0, 3}, meshVert{4, 1, 3}, meshVert{1, 2, 3}})
	return
}

func meshProviderPrefabQuad(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	meshData.addPositions(meshVertAtt3{1, 1, 0}, meshVertAtt3{-1, 1, 0}, meshVertAtt3{-1, -1, 0}, meshVertAtt3{1, -1, 0})
	meshData.addTexCoords(meshVertAtt2{-0.125, 0}, meshVertAtt2{-0.125, 3}, meshVertAtt2{1.125, 3}, meshVertAtt2{1.125, 0})
	meshData.addNormals(meshVertAtt3{0, 0, 1})
	meshData.addFaces(
		meshFace3{meshVert{0, 0, 0}, meshVert{1, 1, 0}, meshVert{2, 2, 0}},
		meshFace3{meshVert{0, 0, 0}, meshVert{2, 2, 0}, meshVert{3, 3, 0}})
	return
}

func meshProviderPrefabTri(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	meshData.addPositions(meshVertAtt3{0, 1, 0}, meshVertAtt3{-1, -1, 0}, meshVertAtt3{1, -1, 0})
	meshData.addTexCoords(meshVertAtt2{0, 0}, meshVertAtt2{3, 0}, meshVertAtt2{3, 2})
	meshData.addNormals(meshVertAtt3{0, 0, 1})
	meshData.addFaces(meshFace3{meshVert{0, 0, 0}, meshVert{1, 1, 0}, meshVert{2, 2, 0}})
	return
}
