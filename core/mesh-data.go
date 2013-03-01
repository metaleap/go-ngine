package core

import (
	"strings"

	gl "github.com/go3d/go-opengl/core"
)

//	Represents semi-processed loaded mesh data "almost ready" to core.Mesh.GpuUpload().
type meshRaw struct {
	lastNumIndices gl.Sizei

	//	Raw vertices
	verts []float32

	//	Vertex indices
	indices []uint32

	//	Raw face definitions
	faces []meshRawFace
}

//	Represents a triangle face inside a meshRaw.
type meshRawFace struct {
	//	Indices of the triangle corners
	entries [3]uint32

	base MeshFaceBase
}

type MeshFaceBase struct {
	//	Mesh-unique identifier for this face.
	ID string

	//	Arbitrary classification tags for this face.
	Tags []string
}

//	Represents an indexed triangle face.
type MeshDescF3 struct {
	//	The indexed vertices making up this triangle face.
	V [3]MeshDescF3V

	//	ID, Tags
	MeshFaceBase
}

//	Creates and initializes a new MeshDescF3V with the specified tags,
//	ID and verts, and returns it. tags may be empty or contain multiple
//	classification tags separated by spaces, which will be split into Tags.
func NewMeshDescF3(tags, id string, verts ...MeshDescF3V) (me *MeshDescF3) {
	me = &MeshDescF3{V: [3]MeshDescF3V{verts[0], verts[1], verts[2]}}
	if me.ID = id; len(tags) > 0 {
		me.Tags = strings.Split(tags, " ")
	}
	return
}

//	Represents an indexed vertex in a MeshDescF3.
type MeshDescF3V struct {
	//	Index of the vertex position
	PosIndex uint32

	//	Index of the texture-coordinate.
	TexCoordIndex uint32

	//	Index of the vertex normal.
	NormalIndex uint32
}

//	Represents a 2-component vertex attribute in a MeshDescriptor.
//	(such as for example texture-coordinates)
type MeshDescVA2 [2]float32

//	Represents a 3-component vertex attribute in a MeshDescriptor
//	(such as for example vertex-normals)
type MeshDescVA3 [3]float32

//	Represents yet-unprocessed, descriptive mesh source data.
type MeshDescriptor struct {
	//	Vertex positions
	Positions []MeshDescVA3

	//	Vertex texture coordinates
	TexCoords []MeshDescVA2

	//	Vertex normals
	Normals []MeshDescVA3

	//	Indexed triangle definitions
	Faces []MeshDescF3
}

//	Adds all specified Faces to this MeshDescriptor.
func (me *MeshDescriptor) AddFaces(faces ...*MeshDescF3) {
	if len(me.Faces) == 0 {
		me.Faces = make([]MeshDescF3, 0, len(faces))
	}
	for i := 0; i < len(faces); i++ {
		me.Faces = append(me.Faces, *faces[i])
	}
}

//	Adds all specified Positions to this MeshDescriptor.
func (me *MeshDescriptor) AddPositions(positions ...MeshDescVA3) {
	me.Positions = append(me.Positions, positions...)
}

//	Adds all the specified Normals to this MeshDescriptor.
func (me *MeshDescriptor) AddNormals(normals ...MeshDescVA3) {
	me.Normals = append(me.Normals, normals...)
}

//	Adds all the specified TexCoords to this MeshDescriptor.
func (me *MeshDescriptor) AddTexCoords(texCoords ...MeshDescVA2) {
	me.TexCoords = append(me.TexCoords, texCoords...)
}
