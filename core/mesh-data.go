package core

import (
	"strings"
)

//	Provides MeshData for constructing a Mesh. An implementation
//	might load a certain file format or procedurally generate
//	the returned MeshData.
type MeshProvider func(args ...interface{}) (*MeshData, error)

type MeshFaceBase struct {
	//	Mesh-unique identifier for this face.
	ID string

	//	Arbitrary classification(s) for this face.
	Classes []string
}

//	Represents an indexed triangle.
type MeshFace3 struct {
	//	The indexed vertices making up this triangle.
	V [3]MeshVert

	//	ID, Classes
	MeshFaceBase
}

//	Creates and initializes a new MeshVert with the specified classes,
//	ID and verts, and returns it. class may be empty or contain multiple
//	classifications separated by spaces, which will be split into Classes.
func NewMeshFace3(class, id string, verts ...MeshVert) (me *MeshFace3) {
	me = &MeshFace3{}
	for i := 0; i < 3; i++ {
		me.V[i] = verts[0]
	}
	me.ID = id
	if len(class) > 0 {
		me.Classes = strings.Split(class, " ")
	}
	return
}

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

	base MeshFaceBase
}

//	Initializes and returns a new *meshRawFace* instance.
func newMeshRawFace(base *MeshFaceBase) (me *meshRawFace) {
	me = &meshRawFace{}
	me.base = *base
	return
}

//	Represents an indexed vertex in a MeshFace3.
type MeshVert struct {
	//	Index of the vertex position
	PosIndex uint32

	//	Index of the texture-coordinate.
	TexCoordIndex uint32

	//	Index of the vertex normal.
	NormalIndex uint32
}

//	Represents a 2-component vertex attribute in a MeshData.
//	(such as for example texture-coordinates)
type MeshVertAtt2 [2]float32

//	Represents a 3-component vertex attribute in a MeshData
//	(such as for example vertex-normals)
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
func (me *MeshData) AddFaces(faces ...*MeshFace3) {
	for _, f := range faces {
		me.Faces = append(me.Faces, *f)
	}
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
