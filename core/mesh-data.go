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

	//	Arbitrary classification tags for this face.
	Tags []string
}

//	Represents an indexed triangle.
type MeshF3 struct {
	//	The indexed vertices making up this triangle.
	V [3]MeshV

	//	ID, Tags
	MeshFaceBase
}

//	Creates and initializes a new MeshV with the specified tags,
//	ID and verts, and returns it. tags may be empty or contain multiple
//	classification tags separated by spaces, which will be split into Tags.
func NewMeshF3(tags, id string, verts ...MeshV) (me *MeshF3) {
	me = &MeshF3{V: [3]MeshV{verts[0], verts[1], verts[2]}}
	if me.ID = id; len(tags) > 0 {
		me.Tags = strings.Split(tags, " ")
	}
	return
}

//	Represents semi-processed loaded mesh data "almost ready" to core.Mesh.GpuUpload().
type meshRaw struct {
	//	Raw vertices
	meshVs []float32

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

//	Represents an indexed vertex in a MeshF3.
type MeshV struct {
	//	Index of the vertex position
	PosIndex uint32

	//	Index of the texture-coordinate.
	TexCoordIndex uint32

	//	Index of the vertex normal.
	NormalIndex uint32
}

//	Represents a 2-component vertex attribute in a MeshData.
//	(such as for example texture-coordinates)
type MeshVA2 [2]float32

//	Represents a 3-component vertex attribute in a MeshData
//	(such as for example vertex-normals)
type MeshVA3 [3]float32

//	Represents yet-unprocessed mesh source data.
type MeshData struct {
	//	Vertex positions
	Positions []MeshVA3

	//	Vertex texture coordinates
	TexCoords []MeshVA2

	//	Vertex normals
	Normals []MeshVA3

	//	Indexed triangle definitions
	Faces []MeshF3
}

//	Initializes and returns a new *MeshData* instance.
func NewMeshData() (me *MeshData) {
	me = &MeshData{}
	return
}

//	Adds all specified Faces to this MeshData.
func (me *MeshData) AddFaces(faces ...*MeshF3) {
	for _, f := range faces {
		me.Faces = append(me.Faces, *f)
	}
}

//	Adds all specified Positions to this MeshData.
func (me *MeshData) AddPositions(positions ...MeshVA3) {
	me.Positions = append(me.Positions, positions...)
}

//	Adds all the specified Normals to this MeshData.
func (me *MeshData) AddNormals(normals ...MeshVA3) {
	me.Normals = append(me.Normals, normals...)
}

//	Adds all the specified TexCoords to this MeshData.
func (me *MeshData) AddTexCoords(texCoords ...MeshVA2) {
	me.TexCoords = append(me.TexCoords, texCoords...)
}
