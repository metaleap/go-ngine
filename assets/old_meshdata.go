package assets

type MeshProvider func (args ... interface {}) (*MeshData, error)

type MeshFace3 [3]Vert

type MeshRaw struct {
	Verts []float32
	Indices []uint32
	Faces []*MeshFace
}

type Vert struct {
	PosIndex, TexCoordIndex, NormalIndex uint32
}

type VertAtt2 [2]float32

type VertAtt3 [3]float32

type MeshFace struct {
	Entries [3]uint32
}

	func NewMeshFace () (me *MeshFace) {
		me = &MeshFace {}
		return
	}

type MeshData struct {
	Positions []VertAtt3
	TexCoords []VertAtt2
	Normals []VertAtt3
	Faces []MeshFace3
}

	func NewMeshData () (me *MeshData) {
		me = &MeshData {}
		me.Positions = []VertAtt3 {}
		me.TexCoords = []VertAtt2 {}
		me.Normals = []VertAtt3 {}
		me.Faces = []MeshFace3 {}
		return
	}

	func (me *MeshData) AddFaces (faces ... MeshFace3) {
		me.Faces = append(me.Faces, faces ...)
	}

	func (me *MeshData) AddPositions (positions ... VertAtt3) {
		me.Positions = append(me.Positions, positions ...)
	}

	func (me *MeshData) AddNormals (normals ... VertAtt3) {
		me.Normals = append(me.Normals, normals ...)
	}

	func (me *MeshData) AddTexCoords (texCoords ... VertAtt2) {
		me.TexCoords = append(me.TexCoords, texCoords ...)
	}
