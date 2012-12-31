package core

//	Temporary concoction.
type MeshProvider func(args ...interface{}) (*MeshData, error)

//	Represents an indexed triangle.
type meshFace3 [3]meshVert

//	Represents semi-processed loaded mesh data "almost ready" to core.Mesh.GpuUpload().
type meshRaw struct {
	//	Raw vertices
	meshVerts []float32

	//	Vertex indices
	indices []uint32

	//	Raw face definitions
	faces []*meshRawFace
}

//	Represents a triangle face inside a meshRaw.
type meshRawFace struct {
	//	Indices of the triangle corners
	entries [3]uint32
}

//	Initializes and returns a new *meshRawFace* instance.
func newMeshRawFace() (me *meshRawFace) {
	me = &meshRawFace{}
	return
}

//	Represents an indexed vertex.
type meshVert struct {
	//	Index of the vertex position
	posIndex uint32

	//	Index of the texture-coordinate.
	texCoordIndex uint32

	//	Index of the vertex normal.
	normalIndex uint32
}

//	Represents a 2-component vertex attribute (such as for example texture-coordinates)
type meshVertAtt2 [2]float32

//	Represents a 3-component vertex attribute (such as for example vertex-normals)
type meshVertAtt3 [3]float32

//	Represents yet-unprocessed mesh source data.
type MeshData struct {
	//	Vertex positions
	positions []meshVertAtt3

	//	Vertex texture coordinates
	texCoords []meshVertAtt2

	//	Vertex normals
	normals []meshVertAtt3

	//	Indexed triangle definitions
	faces []meshFace3
}

//	Initializes and returns a new *MeshData* instance.
func NewMeshData() (me *MeshData) {
	me = &MeshData{}
	return
}

//	Adds all specified Faces to this MeshData.
func (me *MeshData) addFaces(faces ...meshFace3) {
	me.faces = append(me.faces, faces...)
}

//	Adds all specified Positions to this MeshData.
func (me *MeshData) addPositions(positions ...meshVertAtt3) {
	me.positions = append(me.positions, positions...)
}

//	Adds all the specified Normals to this MeshData.
func (me *MeshData) addNormals(normals ...meshVertAtt3) {
	me.normals = append(me.normals, normals...)
}

//	Adds all the specified TexCoords to this MeshData.
func (me *MeshData) addTexCoords(texCoords ...meshVertAtt2) {
	me.texCoords = append(me.texCoords, texCoords...)
}
