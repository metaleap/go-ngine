package core

import (
	gl "github.com/go3d/go-opengl/core"

	ugl "github.com/go3d/go-opengl/util"
)

type meshBufferParams struct {
	HugeMeshSupport, MostlyStatic, CompressTexCoords, CompressTexCoordsNeg, CompressNormals, CompressPositions bool
	NumVerts, NumIndices                                                                                       int32
}

type MeshBuffers struct {
	bufs map[string]*MeshBuffer
}

func (me *MeshBuffers) init() {
	me.bufs = map[string]*MeshBuffer{}
}

func (me *MeshBuffers) Add(id string, params *meshBufferParams) (buf *MeshBuffer, err error) {
	buf = me.bufs[id]
	if buf == nil {
		if buf, err = newMeshBuffer(id, params); err == nil {
			me.bufs[id] = buf
		} else if buf != nil {
			buf.dispose()
		}
	} else {
		err = errf("Cannot add a new mesh buffer with ID '%v': already exists", id)
	}
	return
}

func (me *MeshBuffers) dispose() {
	for _, buf := range me.bufs {
		buf.dispose()
	}
	me.bufs = map[string]*MeshBuffer{}
}

func (_ *MeshBuffers) FloatsPerVertex() int32 {
	const numVertPosFloats, numVertTexCoordFloats, numVertNormalFloats int32 = 3, 2, 3
	return numVertPosFloats + numVertNormalFloats + numVertTexCoordFloats
}

func (_ *MeshBuffers) MemSizePerIndex() int32 {
	return 4
}

func (me *MeshBuffers) MemSizePerVertex() int32 {
	const sizePerFloat int32 = 4
	return sizePerFloat * me.FloatsPerVertex()
}

func (me *MeshBuffers) NewParams(numVerts, numIndices int32) (params *meshBufferParams) {
	params = &meshBufferParams{MostlyStatic: true, NumIndices: numIndices, NumVerts: numVerts}
	return
}

func (me *MeshBuffers) Remove(id string) {
	if buf := me.bufs[id]; buf != nil {
		buf.dispose()
		delete(me.bufs, id)
	}
}

type MeshBuffer struct {
	MemSizeIndices, MemSizeVertices int32
	// Params                         meshBufferParams
	Name string

	offsetBaseIndex, offsetIndices, offsetVerts int32
	id                                          string
	glIbo, glVbo                                ugl.Buffer
	glVaos                                      map[*ugl.Program]*ugl.VertexArray
	meshes                                      map[*Mesh]bool
}

func (me *MeshBuffer) init() {
}

func newMeshBuffer(id string, params *meshBufferParams) (me *MeshBuffer, err error) {
	me = &MeshBuffer{}
	me.id = id
	me.meshes = map[*Mesh]bool{}
	// me.Params = *params
	me.glVaos = map[*ugl.Program]*ugl.VertexArray{}
	me.MemSizeIndices = Core.MeshBuffers.MemSizePerIndex() * params.NumIndices
	me.MemSizeVertices = Core.MeshBuffers.MemSizePerVertex() * params.NumVerts
	if err = me.glVbo.Recreate(gl.ARRAY_BUFFER, gl.Sizeiptr(me.MemSizeVertices), ugl.PtrNil, ugl.Typed.Ife(params.MostlyStatic, gl.STATIC_DRAW, gl.DYNAMIC_DRAW)); err == nil {
		err = me.glIbo.Recreate(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(me.MemSizeIndices), ugl.PtrNil, ugl.Typed.Ife(params.MostlyStatic, gl.STATIC_DRAW, gl.DYNAMIC_DRAW))
	}
	if err == nil {
		var tech RenderTechnique
		var ok bool
		for i := 0; i < len(ogl.progs.All); i++ {
			if tech, ok = ogl.progs.All[i].Tag.(RenderTechnique); ok {
				if err = me.setupVao(&ogl.progs.All[i], tech); err != nil {
					break
				}
			}
		}
	}
	if err != nil {
		me.dispose()
		me = nil
	}
	return
}

func (me *MeshBuffer) Add(mesh *Mesh) (err error) {
	if mesh.meshBuffer != nil {
		err = errf("Cannot add mesh '%v' to mesh buffer '%v': already belongs to mesh buffer '%v'.", mesh.Name, me.id, mesh.meshBuffer.id)
	} else if !me.meshes[mesh] {
		me.meshes[mesh] = true
		mesh.gpuSynced = false
		mesh.meshBuffer = me
	} else {
		err = errf("Cannot add mesh '%v' to mesh buffer '%v': already added.", mesh.Name, me.id)
	}
	return
}

func (me *MeshBuffer) use() {
	me.glVaos[thrRend.curProg].Bind()
}

func (me *MeshBuffer) dispose() {
	for mesh, _ := range me.meshes {
		mesh.meshBuffer, mesh.gpuSynced = nil, false
	}
	me.glIbo.Dispose()
	me.glVbo.Dispose()
	for _, glVao := range me.glVaos {
		glVao.Dispose()
	}
	me.glVaos = nil
}

func (me *MeshBuffer) Remove(mesh *Mesh) {
	if mesh.meshBuffer == me {
		mesh.GpuDelete()
		mesh.meshBuffer = nil
		delete(me.meshes, mesh)
	}
}

func (me *MeshBuffer) setupVao(prog *ugl.Program, tech RenderTechnique) (err error) {
	vao := &ugl.VertexArray{}
	if err = vao.Create(); err == nil {
		if sceneTech, ok := tech.(*RenderTechniqueScene); ok {
			if err = vao.Setup(prog, sceneTech.vertexAttribPointers(me), &me.glVbo, &me.glIbo); err != nil {
				vao.Dispose()
				vao = nil
			}
		}
		if vao != nil {
			me.glVaos[prog] = vao
		}
	}
	return
}

//#begin-gt -gen-reflib.gt T:MeshBuffer L:Core.MeshBuffers

//	Only used for Core.MeshBuffers
type MeshBufferLib []*MeshBuffer

func (me *MeshBufferLib) AddNew() (ref *MeshBuffer) {
	ref = new(MeshBuffer)
	*me = append(*me, ref)
	ref.init()
	return
}

func (me *MeshBufferLib) init() {
	*me = make(MeshBufferLib, 0, 4)
}

func (me *MeshBufferLib) dispose() {
	me.Remove(0, 0)
}

func (me MeshBufferLib) Get(id int) (ref *MeshBuffer) {
	if me.IsOk(id) {
		ref = me[id]
	}
	return
}

func (me MeshBufferLib) IsOk(id int) bool {
	return id > -1 && id < len(me)
}

func (me *MeshBufferLib) Remove(fromID, num int) {
	if l := len(*me); fromID > -1 && fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		for i := fromID; i < fromID+num; i++ {
			(*me)[fromID].dispose()
		}
		before, after := (*me)[:fromID], (*me)[fromID+num:]
		*me = append(before, after...)
	}
}

func (me MeshBufferLib) Walk(on func(ref *MeshBuffer)) {
	for id := 0; id < len(me); id++ {
		on(me[id])
	}
}

//#end-gt
