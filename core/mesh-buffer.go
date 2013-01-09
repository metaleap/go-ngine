package core

import (
	"fmt"

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
		err = fmt.Errorf("Cannot add a new mesh buffer with ID '%v': already exists", id)
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
	meshes                                      Meshes
	glIbo, glVbo                                gl.Uint
	glVaos                                      map[string]gl.Uint
}

func newMeshBuffer(id string, params *meshBufferParams) (buf *MeshBuffer, err error) {
	var glVao gl.Uint
	buf = &MeshBuffer{}
	buf.id = id
	buf.meshes = Meshes{}
	buf.Params = params
	buf.glVaos = map[string]gl.Uint{}
	buf.MemSizeIndices = Core.MeshBuffers.MemSizePerIndex() * params.NumIndices
	buf.MemSizeVertices = Core.MeshBuffers.MemSizePerVertex() * params.NumVerts
	gl.GenBuffers(1, &buf.glVbo)
	gl.GenBuffers(1, &buf.glIbo)
	for techName, _ := range techs {
		gl.GenVertexArrays(1, &glVao)
		buf.glVaos[techName] = glVao
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, buf.glVbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf.glIbo)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(buf.MemSizeVertices), gl.Pointer(nil), ugl.Ife(params.MostlyStatic, gl.STATIC_DRAW, gl.DYNAMIC_DRAW))
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(buf.MemSizeIndices), gl.Pointer(nil), ugl.Ife(params.MostlyStatic, gl.STATIC_DRAW, gl.DYNAMIC_DRAW))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	if err = ugl.LastError("newMeshBuffer(%v numVerts=%v numIndices=%v)", id, params.NumVerts, params.NumIndices); err != nil {
		buf.dispose()
		buf = nil
	} else {
		for techName, glVao := range buf.glVaos {
			gl.BindVertexArray(glVao)
			gl.BindBuffer(gl.ARRAY_BUFFER, buf.glVbo)
			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, buf.glIbo)
			techs[techName].initMeshBuffer(buf)
			gl.BindVertexArray(0)
			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
			gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		}
	}
	return
}

func (me *MeshBuffer) Add(mesh *Mesh) (err error) {
	if mesh.meshBuffer != nil {
		err = fmt.Errorf("Cannot add mesh '%v' to mesh buffer '%v': already belongs to mesh buffer '%v'.", mesh.id, me.id, mesh.meshBuffer.id)
	} else if me.meshes.Add(mesh) != nil {
		mesh.gpuSynced = false
		mesh.meshBuffer = me
	} else {
		err = fmt.Errorf("Cannot add mesh '%v' to mesh buffer '%v': already has a mesh with that ID.", mesh.id, me.id)
	}
	return
}

func (me *MeshBuffer) use() {
	curMeshBuf = me
	gl.BindVertexArray(me.glVaos[curTechnique.name()])
}

func (me *MeshBuffer) dispose() {
	for _, mesh := range me.meshes {
		mesh.meshBuffer, mesh.gpuSynced = nil, false
	}
	gl.DeleteBuffers(1, &me.glIbo)
	gl.DeleteBuffers(1, &me.glVbo)
	for _, glVao := range me.glVaos {
		gl.DeleteVertexArrays(1, &glVao)
	}
}

func (me *MeshBuffer) Remove(mesh *Mesh) {
	if mesh.meshBuffer == me {
		mesh.GpuDelete()
		mesh.meshBuffer = nil
		delete(me.meshes, mesh.id)
	}
}
