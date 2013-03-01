package core

type MeshProvider func() (*MeshDescriptor, error)

//	A MeshProvider that creates MeshDescriptor for a cube with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshDescriptor contains 12 triangle faces with IDs "t0" through "t11".
//	These faces are classified in 6 distinct tags: "front","back","top","bottom","right","left".
func MeshDescriptorCube() (meshDescriptor *MeshDescriptor, err error) {
	var md MeshDescriptor
	var u, d, l, r, f, b float32
	md.AddNormals(
		//	top
		MeshDescVA3{0, 1, 0},
		//	right
		MeshDescVA3{1, 0, 0},
		//	back
		MeshDescVA3{0, 0, 1},
		//	left
		MeshDescVA3{-1, 0, 0},
		//	front
		MeshDescVA3{0, 0, -1},
		//	bottom
		MeshDescVA3{0, -1, 0},
	)
	md.AddTexCoords(
		MeshDescVA2{1, 0},
		MeshDescVA2{1, 1},
		MeshDescVA2{0, 0},
		MeshDescVA2{0, 1},
	)
	u, d, l, r, f, b = 1, -1, -1, 1, -1, 1
	md.AddPositions(
		MeshDescVA3{l, u, f},
		MeshDescVA3{l, u, b},
		MeshDescVA3{r, u, f},
		MeshDescVA3{r, u, b},
		MeshDescVA3{r, d, f},
		MeshDescVA3{r, d, b},
		MeshDescVA3{l, d, b},
		MeshDescVA3{l, d, f},
	)
	md.AddFaces(
		NewMeshDescF3("top", "t0", MeshDescF3V{0, 0, 0}, MeshDescF3V{1, 1, 0}, MeshDescF3V{2, 2, 0}), NewMeshDescF3("top", "t1", MeshDescF3V{1, 1, 0}, MeshDescF3V{3, 3, 0}, MeshDescF3V{2, 2, 0}),
		NewMeshDescF3("right", "t2", MeshDescF3V{2, 1, 1}, MeshDescF3V{3, 3, 1}, MeshDescF3V{4, 0, 1}), NewMeshDescF3("right", "t3", MeshDescF3V{4, 0, 1}, MeshDescF3V{3, 3, 1}, MeshDescF3V{5, 2, 1}),
		NewMeshDescF3("back", "t4", MeshDescF3V{5, 0, 2}, MeshDescF3V{3, 1, 2}, MeshDescF3V{1, 3, 2}), NewMeshDescF3("back", "t5", MeshDescF3V{5, 0, 2}, MeshDescF3V{1, 3, 2}, MeshDescF3V{6, 2, 2}),
		NewMeshDescF3("left", "t6", MeshDescF3V{6, 0, 3}, MeshDescF3V{1, 1, 3}, MeshDescF3V{0, 3, 3}), NewMeshDescF3("left", "t7", MeshDescF3V{6, 0, 3}, MeshDescF3V{0, 3, 3}, MeshDescF3V{7, 2, 3}),
		NewMeshDescF3("front", "t8", MeshDescF3V{0, 1, 4}, MeshDescF3V{2, 3, 4}, MeshDescF3V{7, 0, 4}), NewMeshDescF3("front", "t9", MeshDescF3V{2, 3, 4}, MeshDescF3V{4, 2, 4}, MeshDescF3V{7, 0, 4}),
		NewMeshDescF3("bottom", "t10", MeshDescF3V{7, 2, 5}, MeshDescF3V{4, 0, 5}, MeshDescF3V{6, 3, 5}), NewMeshDescF3("bottom", "t11", MeshDescF3V{6, 3, 5}, MeshDescF3V{4, 0, 5}, MeshDescF3V{5, 1, 5}),
	)
	meshDescriptor = &md
	return
}

//	A MeshProvider that creates MeshDescriptor for a flat ground plane with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshDescriptor contains 2 triangle faces with IDs "t0" through "t1".
//	These faces are all classified with tag: "plane".
func MeshDescriptorPlane() (meshDescriptor *MeshDescriptor, err error) {
	var md MeshDescriptor
	md.AddPositions(MeshDescVA3{-1, 0, 1}, MeshDescVA3{1, 0, 1}, MeshDescVA3{-1, 0, -1}, MeshDescVA3{1, 0, -1})
	md.AddTexCoords(MeshDescVA2{1000, 1000}, MeshDescVA2{0, 1000}, MeshDescVA2{1000, 0}, MeshDescVA2{0, 0})
	md.AddNormals(MeshDescVA3{0, 1, 0})
	md.AddFaces(
		NewMeshDescF3("plane", "t0", MeshDescF3V{0, 0, 0}, MeshDescF3V{1, 1, 0}, MeshDescF3V{2, 2, 0}),
		NewMeshDescF3("plane", "t1", MeshDescF3V{3, 3, 0}, MeshDescF3V{2, 2, 0}, MeshDescF3V{1, 1, 0}))
	meshDescriptor = &md
	return
}

//	A MeshProvider that creates MeshDescriptor for a pyramid with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshDescriptor contains 4 triangle faces with IDs "t0" through "t3".
//	These faces are all classified with tag: "pyr".
func MeshDescriptorPyramid() (meshDescriptor *MeshDescriptor, err error) {
	var md MeshDescriptor
	md.AddPositions(MeshDescVA3{0, 1, 0}, MeshDescVA3{-1, -1, 1}, MeshDescVA3{1, -1, 1}, MeshDescVA3{1, -1, -1}, MeshDescVA3{-1, -1, -1})
	md.AddTexCoords(MeshDescVA2{0.5, 1}, MeshDescVA2{0, 0}, MeshDescVA2{1, 0})
	md.AddNormals(MeshDescVA3{0, 0, 1}, MeshDescVA3{1, 0, 0}, MeshDescVA3{0, 0, -1}, MeshDescVA3{-1, 0, 0})
	md.AddFaces(
		NewMeshDescF3("pyr", "t0", MeshDescF3V{0, 0, 0}, MeshDescF3V{1, 1, 0}, MeshDescF3V{2, 2, 0}),
		NewMeshDescF3("pyr", "t1", MeshDescF3V{0, 0, 1}, MeshDescF3V{2, 1, 1}, MeshDescF3V{3, 2, 1}),
		NewMeshDescF3("pyr", "t2", MeshDescF3V{0, 0, 2}, MeshDescF3V{3, 1, 2}, MeshDescF3V{4, 2, 2}),
		NewMeshDescF3("pyr", "t3", MeshDescF3V{0, 0, 3}, MeshDescF3V{4, 1, 3}, MeshDescF3V{1, 2, 3}))
	meshDescriptor = &md
	return
}

//	A MeshProvider that creates MeshDescriptor for a quad with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshDescriptor contains 2 triangle faces with IDs "t0" through "t1".
//	These faces are all classified with tag: "quad".
func MeshDescriptorQuad() (meshDescriptor *MeshDescriptor, err error) {
	var md MeshDescriptor
	md.AddPositions(MeshDescVA3{-1, -1, 0}, MeshDescVA3{-1, 1, 0}, MeshDescVA3{1, -1, 0}, MeshDescVA3{1, 1, 0})
	md.AddTexCoords(MeshDescVA2{1, 0}, MeshDescVA2{1, 1}, MeshDescVA2{0, 0}, MeshDescVA2{0, 1})
	md.AddNormals(MeshDescVA3{0, 0, 1})
	md.AddFaces(
		NewMeshDescF3("quad", "t0", MeshDescF3V{0, 0, 0}, MeshDescF3V{1, 1, 0}, MeshDescF3V{2, 2, 0}),
		NewMeshDescF3("quad", "t1", MeshDescF3V{1, 1, 0}, MeshDescF3V{3, 3, 0}, MeshDescF3V{2, 2, 0}))
	meshDescriptor = &md
	return
}

//	A MeshProvider that creates MeshDescriptor for a triangle with extents -1 .. 1.
//	args is ignored and err is always nil.
//	The returned MeshDescriptor contains 1 triangle face with ID "t0" and tag "tri".
func MeshDescriptorTri() (meshDescriptor *MeshDescriptor, err error) {
	var md MeshDescriptor
	md.AddPositions(MeshDescVA3{-1, -1, 0}, MeshDescVA3{0, 1, 0}, MeshDescVA3{1, -1, 0})
	md.AddTexCoords(MeshDescVA2{0, 0}, MeshDescVA2{0.5, 1}, MeshDescVA2{1, 0})
	md.AddNormals(MeshDescVA3{0, 0, 1})
	md.AddFaces(NewMeshDescF3("tri", "t0", MeshDescF3V{0, 0, 0}, MeshDescF3V{1, 1, 0}, MeshDescF3V{2, 2, 0}))
	meshDescriptor = &md
	return
}
