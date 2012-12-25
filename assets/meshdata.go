package assets

//	Temporary concoction.
type MeshProvider func(args ...interface{}) (*MeshData, error)

//	Represents an indexed triangle.
type MeshFace3 [3]MeshVert

//	Represents semi-processed loaded mesh data "almost ready" to core.Mesh.GpuUpload().
type MeshRaw struct {
	//	Raw vertices
	MeshVerts []float32

	//	Vertex indices
	Indices []uint32

	//	Raw face definitions
	Faces []*MeshRawFace
}

//	Represents a triangle face inside a MeshRaw.
type MeshRawFace struct {
	//	Indices of the triangle corners
	Entries [3]uint32
}

//	Initializes and returns a new *MeshRawFace* instance.
func NewMeshRawFace() (me *MeshRawFace) {
	me = &MeshRawFace{}
	return
}

//	Represents an indexed vertex.
type MeshVert struct {
	//	Index of the vertex position
	PosIndex uint32

	//	Index of the texture-coordinate.
	TexCoordIndex uint32

	//	Index of the vertex normal.
	NormalIndex uint32
}

//	Represents a 2-component vertex attribute (such as for example texture-coordinates)
type MeshVertAtt2 [2]float32

//	Represents a 3-component vertex attribute (such as for example vertex-normals)
type MeshVertAtt3 [3]float32

//	Represents yet-unprocessed mesh source data.
type MeshData struct {
	//	Vertex positions
	Positions []MeshVertAtt3

	//	Vertex texture coordinates
	TexCoords []MeshVertAtt2

	//	Vertex normals
	Normals []MeshVertAtt3

	//	Indexed triangle definitions
	Faces []MeshFace3
}

//	Initializes and returns a new *MeshData* instance.
func NewMeshData() (me *MeshData) {
	me = &MeshData{}
	return
}

//	Adds all specified Faces to this MeshData.
func (me *MeshData) AddFaces(faces ...MeshFace3) {
	me.Faces = append(me.Faces, faces...)
}

//	Adds all specified Positions to this MeshData.
func (me *MeshData) AddPositions(positions ...MeshVertAtt3) {
	me.Positions = append(me.Positions, positions...)
}

//	Adds all the specified Normals to this MeshData.
func (me *MeshData) AddNormals(normals ...MeshVertAtt3) {
	me.Normals = append(me.Normals, normals...)
}

//	Adds all the specified TexCoords to this MeshData.
func (me *MeshData) AddTexCoords(texCoords ...MeshVertAtt2) {
	me.TexCoords = append(me.TexCoords, texCoords...)
}
