package core

//	A MeshProvider that creates MeshData for a cube with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshData contains 12 triangle faces with IDs "t0" through "t11".
//	These faces are classified in 6 distinct tags: "front","back","top","bottom","right","left".
func MeshProviderPrefabCube(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	var u, d, l, r, f, b float32
	meshData.AddNormals(
		//	top
		MeshVA3{0, 1, 0},
		//	right
		MeshVA3{1, 0, 0},
		//	back
		MeshVA3{0, 0, 1},
		//	left
		MeshVA3{-1, 0, 0},
		//	front
		MeshVA3{0, 0, -1},
		//	bottom
		MeshVA3{0, -1, 0},
	)
	meshData.AddTexCoords(
		MeshVA2{1, 0},
		MeshVA2{1, 1},
		MeshVA2{0, 0},
		MeshVA2{0, 1},
	)
	u, d, l, r, f, b = 1, -1, -1, 1, -1, 1
	meshData.AddPositions(
		MeshVA3{l, u, f},
		MeshVA3{l, u, b},
		MeshVA3{r, u, f},
		MeshVA3{r, u, b},
		MeshVA3{r, d, f},
		MeshVA3{r, d, b},
		MeshVA3{l, d, b},
		MeshVA3{l, d, f},
	)
	meshData.AddFaces(
		NewMeshF3("top", "t0", MeshV{0, 0, 0}, MeshV{1, 1, 0}, MeshV{2, 2, 0}), NewMeshF3("top", "t1", MeshV{1, 1, 0}, MeshV{3, 3, 0}, MeshV{2, 2, 0}),
		NewMeshF3("right", "t2", MeshV{2, 1, 1}, MeshV{3, 3, 1}, MeshV{4, 0, 1}), NewMeshF3("right", "t3", MeshV{4, 0, 1}, MeshV{3, 3, 1}, MeshV{5, 2, 1}),
		NewMeshF3("back", "t4", MeshV{5, 0, 2}, MeshV{3, 1, 2}, MeshV{1, 3, 2}), NewMeshF3("back", "t5", MeshV{5, 0, 2}, MeshV{1, 3, 2}, MeshV{6, 2, 2}),
		NewMeshF3("left", "t6", MeshV{6, 0, 3}, MeshV{1, 1, 3}, MeshV{0, 3, 3}), NewMeshF3("left", "t7", MeshV{6, 0, 3}, MeshV{0, 3, 3}, MeshV{7, 2, 3}),
		NewMeshF3("front", "t8", MeshV{0, 1, 4}, MeshV{2, 3, 4}, MeshV{7, 0, 4}), NewMeshF3("front", "t9", MeshV{2, 3, 4}, MeshV{4, 2, 4}, MeshV{7, 0, 4}),
		NewMeshF3("bottom", "t10", MeshV{7, 2, 5}, MeshV{4, 0, 5}, MeshV{6, 3, 5}), NewMeshF3("bottom", "t11", MeshV{6, 3, 5}, MeshV{4, 0, 5}, MeshV{5, 1, 5}),
	)
	return
}

//	A MeshProvider that creates MeshData for a flat ground plane with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshData contains 2 triangle faces with IDs "t0" through "t1".
//	These faces are all classified with tag: "plane".
func MeshProviderPrefabPlane(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	meshData.AddPositions(MeshVA3{-1, 0, 1}, MeshVA3{1, 0, 1}, MeshVA3{-1, 0, -1}, MeshVA3{1, 0, -1})
	meshData.AddTexCoords(MeshVA2{1000, 1000}, MeshVA2{0, 1000}, MeshVA2{1000, 0}, MeshVA2{0, 0})
	meshData.AddNormals(MeshVA3{0, 1, 0})
	meshData.AddFaces(
		NewMeshF3("plane", "t0", MeshV{0, 0, 0}, MeshV{1, 1, 0}, MeshV{2, 2, 0}),
		NewMeshF3("plane", "t1", MeshV{3, 3, 0}, MeshV{2, 2, 0}, MeshV{1, 1, 0}))
	return
}

//	A MeshProvider that creates MeshData for a pyramid with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshData contains 4 triangle faces with IDs "t0" through "t3".
//	These faces are all classified with tag: "pyr".
func MeshProviderPrefabPyramid(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	meshData.AddPositions(MeshVA3{0, 1, 0}, MeshVA3{-1, -1, 1}, MeshVA3{1, -1, 1}, MeshVA3{1, -1, -1}, MeshVA3{-1, -1, -1})
	meshData.AddTexCoords(MeshVA2{0, 0}, MeshVA2{1, 0}, MeshVA2{1, 1}, MeshVA2{0, 1})
	meshData.AddNormals(MeshVA3{0, 0, 1}, MeshVA3{1, 0, 0}, MeshVA3{0, 0, -1}, MeshVA3{-1, 0, 0})
	meshData.AddFaces(
		NewMeshF3("pyr", "t0", MeshV{0, 0, 0}, MeshV{1, 1, 0}, MeshV{2, 2, 0}),
		NewMeshF3("pyr", "t1", MeshV{0, 1, 1}, MeshV{2, 2, 1}, MeshV{3, 3, 1}),
		NewMeshF3("pyr", "t2", MeshV{0, 1, 2}, MeshV{3, 2, 2}, MeshV{4, 3, 2}),
		NewMeshF3("pyr", "t3", MeshV{0, 0, 3}, MeshV{4, 1, 3}, MeshV{1, 2, 3}))
	return
}

//	A MeshProvider that creates MeshData for a quad with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshData contains 2 triangle faces with IDs "t0" through "t1".
//	These faces are all classified with tag: "quad".
func MeshProviderPrefabQuad(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	meshData.AddPositions(MeshVA3{-1, -1, 0}, MeshVA3{-1, 1, 0}, MeshVA3{1, -1, 0}, MeshVA3{1, 1, 0})
	meshData.AddTexCoords(MeshVA2{1, 0}, MeshVA2{1, 1}, MeshVA2{0, 0}, MeshVA2{0, 1})
	meshData.AddNormals(MeshVA3{0, 0, 1})
	meshData.AddFaces(
		NewMeshF3("quad", "t0", MeshV{0, 0, 0}, MeshV{1, 1, 0}, MeshV{2, 2, 0}),
		NewMeshF3("quad", "t1", MeshV{1, 1, 0}, MeshV{3, 3, 0}, MeshV{2, 2, 0}))
	return
}

//	A MeshProvider that creates MeshData for a triangle with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshData contains 1 triangle face with ID "t0" and tag "tri".
func MeshProviderPrefabTri(args ...interface{}) (meshData *MeshData, err error) {
	meshData = NewMeshData()
	meshData.AddPositions(MeshVA3{-1, -1, 0}, MeshVA3{0, 1, 0}, MeshVA3{1, -1, 0})
	meshData.AddTexCoords(MeshVA2{0, 0}, MeshVA2{0.5, 1}, MeshVA2{1, 0})
	meshData.AddNormals(MeshVA3{0, 0, 1})
	meshData.AddFaces(NewMeshF3("tri", "t0", MeshV{0, 0, 0}, MeshV{1, 1, 0}, MeshV{2, 2, 0}))
	return
}
