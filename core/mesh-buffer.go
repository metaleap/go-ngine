package core

import (
	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
)

type meshBufferParams struct {
	HugeMeshSupport, MostlyStatic, CompressTexCoords, CompressTexCoordsNeg, CompressNormals, CompressPositions bool
	NumVerts, NumIndices                                                                                       int32
}

type MeshBuffers struct {
	bufs map[string]*MeshBuffer
}

func newMeshBuffers() (me *MeshBuffers) {
	me = &MeshBuffers{}
	me.bufs = map[string]*MeshBuffer{}
	return
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
		err = fmtErr("Cannot add a new mesh buffer with ID '%v': already exists", id)
	}
	return
}

func (me *MeshBuffers) dispose() {
	for _, buf := range me.bufs {
		buf.dispose()
	}
	me.bufs = map[string]*MeshBuffer{}
}

func (me *MeshBuffers) FloatsPerVertex() int32 {
	const numVertPosFloats, numVertTexCoordFloats, numVertNormalFloats int32 = 3, 2, 3
	return numVertPosFloats + numVertNormalFloats + numVertTexCoordFloats
}

func (me *MeshBuffers) MemSizePerIndex() int32 {
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
	Params                          *meshBufferParams

	offsetBaseIndex, offsetIndices, offsetVerts int32
	id                                          string
	glIbo, glVbo                                ugl.Buffer
	glVaos                                      map[string]*ugl.VertexArray
	meshes                                      map[*Mesh]bool
}

func newMeshBuffer(id string, params *meshBufferParams) (me *MeshBuffer, err error) {
	me = &MeshBuffer{}
	me.id = id
	me.meshes = map[*Mesh]bool{}
	me.Params = params
	me.glVaos = map[string]*ugl.VertexArray{}
	me.MemSizeIndices = Core.MeshBuffers.MemSizePerIndex() * params.NumIndices
	me.MemSizeVertices = Core.MeshBuffers.MemSizePerVertex() * params.NumVerts
	me.glVbo.Recreate(gl.ARRAY_BUFFER, gl.Sizeiptr(me.MemSizeVertices), gl.Pointer(nil), ugl.Typed.Ife(params.MostlyStatic, gl.STATIC_DRAW, gl.DYNAMIC_DRAW))
	me.glIbo.Recreate(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(me.MemSizeIndices), gl.Pointer(nil), ugl.Typed.Ife(params.MostlyStatic, gl.STATIC_DRAW, gl.DYNAMIC_DRAW))
	for techName, _ := range techs {
		glVao := &ugl.VertexArray{}
		glVao.Create()
		me.glVaos[techName] = glVao
	}
	if err = ugl.LastError("newMeshBuffer(%v numVerts=%v numIndices=%v)", id, params.NumVerts, params.NumIndices); err != nil {
		me.dispose()
		me = nil
	} else {
		for techName, glVao := range me.glVaos {
			glVao.Setup(techs[techName].initMeshBuffer(me), &me.glVbo, &me.glIbo)
		}
	}
	return
}

func (me *MeshBuffer) Add(mesh *Mesh) (err error) {
	if mesh.meshBuffer != nil {
		err = fmtErr("Cannot add mesh '%v' to mesh buffer '%v': already belongs to mesh buffer '%v'.", mesh.id, me.id, mesh.meshBuffer.id)
	} else if !me.meshes[mesh] {
		me.meshes[mesh] = true
		mesh.gpuSynced = false
		mesh.meshBuffer = me
	} else {
		err = fmtErr("Cannot add mesh '%v' to mesh buffer '%v': already added.", mesh.id, me.id)
	}
	return
}

func (me *MeshBuffer) use() {
	thrRend.curMeshBuf = me
	me.glVaos[thrRend.curTechnique.name()].Bind()
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
	me.glVaos = map[string]*ugl.VertexArray{}
}

func (me *MeshBuffer) Remove(mesh *Mesh) {
	if mesh.meshBuffer == me {
		mesh.GpuDelete()
		mesh.meshBuffer = nil
		delete(me.meshes, mesh)
	}
}
