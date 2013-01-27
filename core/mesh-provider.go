package core

//	A MeshProvider that creates MeshData for a cube with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshData contains 12 triangle faces with IDs "t0" through "t11".
//	These faces are classified in 6 distinct tags: "front","back","top","bottom","right","left".
func MeshProviderPrefabCube(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	meshData.AddPositions(
		MeshVertAtt3{-1, -1, 1}, MeshVertAtt3{1, -1, 1}, MeshVertAtt3{1, 1, 1},
		MeshVertAtt3{-1, 1, 1}, MeshVertAtt3{-1, -1, -1}, MeshVertAtt3{-1, 1, -1},
		MeshVertAtt3{1, 1, -1}, MeshVertAtt3{1, -1, -1})
	meshData.AddTexCoords(MeshVertAtt2{0, 0}, MeshVertAtt2{1, 0}, MeshVertAtt2{1, 1}, MeshVertAtt2{0, 1})
	meshData.AddNormals(MeshVertAtt3{0, 0, 1}, MeshVertAtt3{0, 0, -1}, MeshVertAtt3{0, 1, 0}, MeshVertAtt3{0, -1, 0}, MeshVertAtt3{1, 0, 0}, MeshVertAtt3{-1, 0, 0})
	meshData.AddFaces(
		NewMeshFace3("front", "t0", MeshVert{0, 0, 0}, MeshVert{1, 1, 0}, MeshVert{2, 2, 0}), NewMeshFace3("front", "t1", MeshVert{0, 0, 0}, MeshVert{2, 2, 0}, MeshVert{3, 3, 0}),
		NewMeshFace3("back", "t2", MeshVert{4, 0, 1}, MeshVert{5, 1, 1}, MeshVert{6, 2, 1}), NewMeshFace3("back", "t3", MeshVert{4, 0, 1}, MeshVert{6, 2, 1}, MeshVert{7, 3, 1}),
		NewMeshFace3("top", "t4", MeshVert{5, 0, 2}, MeshVert{3, 1, 2}, MeshVert{2, 2, 2}), NewMeshFace3("top", "t5", MeshVert{5, 0, 2}, MeshVert{2, 2, 2}, MeshVert{6, 3, 2}),
		NewMeshFace3("bottom", "t6", MeshVert{4, 0, 3}, MeshVert{7, 1, 3}, MeshVert{1, 2, 3}), NewMeshFace3("bottom", "t7", MeshVert{4, 0, 3}, MeshVert{1, 2, 3}, MeshVert{0, 3, 3}),
		NewMeshFace3("right", "t8", MeshVert{7, 0, 4}, MeshVert{6, 1, 4}, MeshVert{2, 2, 4}), NewMeshFace3("right", "t9", MeshVert{7, 0, 4}, MeshVert{2, 2, 4}, MeshVert{1, 3, 4}),
		NewMeshFace3("left", "t10", MeshVert{4, 0, 5}, MeshVert{0, 1, 5}, MeshVert{3, 2, 5}), NewMeshFace3("left", "t11", MeshVert{4, 0, 5}, MeshVert{3, 2, 5}, MeshVert{5, 3, 5}))
	return
}

//	A MeshProvider that creates MeshData for a flat ground plane with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshData contains 2 triangle faces with IDs "t0" through "t1".
//	These faces are all classified with tag: "plane".
func MeshProviderPrefabPlane(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	meshData.AddPositions(MeshVertAtt3{-1, 0, 1}, MeshVertAtt3{1, 0, 1}, MeshVertAtt3{-1, 0, -1}, MeshVertAtt3{1, 0, -1})
	meshData.AddTexCoords(MeshVertAtt2{0, 0}, MeshVertAtt2{1000, 0}, MeshVertAtt2{0, 1000}, MeshVertAtt2{1000, 1000})
	meshData.AddNormals(MeshVertAtt3{0, 1, 0})
	meshData.AddFaces(
		NewMeshFace3("plane", "t0", MeshVert{0, 0, 0}, MeshVert{1, 1, 0}, MeshVert{2, 2, 0}),
		NewMeshFace3("plane", "t1", MeshVert{3, 3, 0}, MeshVert{2, 2, 0}, MeshVert{1, 1, 0}))
	return
}

//	A MeshProvider that creates MeshData for a pyramid with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshData contains 4 triangle faces with IDs "t0" through "t3".
//	These faces are all classified with tag: "pyr".
func MeshProviderPrefabPyramid(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	meshData.AddPositions(MeshVertAtt3{0, 1, 0}, MeshVertAtt3{-1, -1, 1}, MeshVertAtt3{1, -1, 1}, MeshVertAtt3{1, -1, -1}, MeshVertAtt3{-1, -1, -1})
	meshData.AddTexCoords(MeshVertAtt2{0, 0}, MeshVertAtt2{1, 0}, MeshVertAtt2{1, 1}, MeshVertAtt2{0, 1})
	meshData.AddNormals(MeshVertAtt3{0, 0, 1}, MeshVertAtt3{1, 0, 0}, MeshVertAtt3{0, 0, -1}, MeshVertAtt3{-1, 0, 0})
	meshData.AddFaces(
		NewMeshFace3("pyr", "t0", MeshVert{0, 0, 0}, MeshVert{1, 1, 0}, MeshVert{2, 2, 0}),
		NewMeshFace3("pyr", "t1", MeshVert{0, 1, 1}, MeshVert{2, 2, 1}, MeshVert{3, 3, 1}),
		NewMeshFace3("pyr", "t2", MeshVert{0, 1, 2}, MeshVert{3, 2, 2}, MeshVert{4, 3, 2}),
		NewMeshFace3("pyr", "t3", MeshVert{0, 0, 3}, MeshVert{4, 1, 3}, MeshVert{1, 2, 3}))
	return
}

//	A MeshProvider that creates MeshData for a quad with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshData contains 2 triangle faces with IDs "t0" through "t1".
//	These faces are all classified with tag: "quad".
func MeshProviderPrefabQuad(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	meshData.AddPositions(MeshVertAtt3{1, 1, 0}, MeshVertAtt3{-1, 1, 0}, MeshVertAtt3{-1, -1, 0}, MeshVertAtt3{1, -1, 0})
	meshData.AddTexCoords(MeshVertAtt2{0, 0}, MeshVertAtt2{0, 1}, MeshVertAtt2{1, 1}, MeshVertAtt2{1, 0})
	meshData.AddNormals(MeshVertAtt3{0, 0, 1})
	meshData.AddFaces(
		NewMeshFace3("quad", "t0", MeshVert{0, 0, 0}, MeshVert{1, 1, 0}, MeshVert{2, 2, 0}),
		NewMeshFace3("quad", "t1", MeshVert{0, 0, 0}, MeshVert{2, 2, 0}, MeshVert{3, 3, 0}))
	return
}

//	A MeshProvider that creates MeshData for a triangle with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshData contains 1 triangle face with ID "t0" and tag "tri".
func MeshProviderPrefabTri(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	meshData.AddPositions(MeshVertAtt3{0, 1, 0}, MeshVertAtt3{-1, -1, 0}, MeshVertAtt3{1, -1, 0})
	meshData.AddTexCoords(MeshVertAtt2{0, 0}, MeshVertAtt2{3, 0}, MeshVertAtt2{3, 2})
	meshData.AddNormals(MeshVertAtt3{0, 0, 1})
	meshData.AddFaces(NewMeshFace3("tri", "t0", MeshVert{0, 0, 0}, MeshVert{1, 1, 0}, MeshVert{2, 2, 0}))
	return
}
