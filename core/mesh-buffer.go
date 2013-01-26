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
	glIbo, glVbo                                ugl.Buffer
	glVaos                                      map[string]gl.Uint
	meshes                                      map[*Mesh]bool
}

func newMeshBuffer(id string, params *meshBufferParams) (me *MeshBuffer, err error) {
	var glVao gl.Uint
	me = &MeshBuffer{}
	me.id = id
	me.meshes = map[*Mesh]bool{}
	me.Params = params
	me.glVaos = map[string]gl.Uint{}
	me.MemSizeIndices = Core.MeshBuffers.MemSizePerIndex() * params.NumIndices
	me.MemSizeVertices = Core.MeshBuffers.MemSizePerVertex() * params.NumVerts
	me.glVbo.Recreate(gl.ARRAY_BUFFER, gl.Sizeiptr(me.MemSizeVertices), gl.Pointer(nil), ugl.Typed.Ife(params.MostlyStatic, gl.STATIC_DRAW, gl.DYNAMIC_DRAW))
	me.glIbo.Recreate(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(me.MemSizeIndices), gl.Pointer(nil), ugl.Typed.Ife(params.MostlyStatic, gl.STATIC_DRAW, gl.DYNAMIC_DRAW))
	// gl.GenBuffers(1, &me.glVbo)
	// gl.GenBuffers(1, &me.glIbo)
	for techName, _ := range techs {
		gl.GenVertexArrays(1, &glVao)
		me.glVaos[techName] = glVao
	}
	// gl.BindBuffer(gl.ARRAY_BUFFER, me.glVbo)
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, me.glIbo)
	// gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(me.MemSizeVertices), gl.Pointer(nil), ugl.Ife(params.MostlyStatic, gl.STATIC_DRAW, gl.DYNAMIC_DRAW))
	// gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(me.MemSizeIndices), gl.Pointer(nil), ugl.Ife(params.MostlyStatic, gl.STATIC_DRAW, gl.DYNAMIC_DRAW))
	// gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	if err = ugl.LastError("newMeshBuffer(%v numVerts=%v numIndices=%v)", id, params.NumVerts, params.NumIndices); err != nil {
		me.dispose()
		me = nil
	} else {
		for techName, glVao := range me.glVaos {
			gl.BindVertexArray(glVao)
			me.glVbo.Bind()
			me.glIbo.Bind()
			// gl.BindBuffer(gl.ARRAY_BUFFER, me.glVbo)
			// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, me.glIbo)
			techs[techName].initMeshBuffer(me)
			gl.BindVertexArray(0)
			me.glIbo.Unbind()
			me.glVbo.Unbind()
			// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
			// gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		}
	}
	return
}

func (me *MeshBuffer) Add(mesh *Mesh) (err error) {
	if mesh.meshBuffer != nil {
		err = fmt.Errorf("Cannot add mesh '%v' to mesh buffer '%v': already belongs to mesh buffer '%v'.", mesh.id, me.id, mesh.meshBuffer.id)
	} else if !me.meshes[mesh] {
		me.meshes[mesh] = true
		mesh.gpuSynced = false
		mesh.meshBuffer = me
	} else {
		err = fmt.Errorf("Cannot add mesh '%v' to mesh buffer '%v': already added.", mesh.id, me.id)
	}
	return
}

func (me *MeshBuffer) use() {
	curMeshBuf = me
	gl.BindVertexArray(me.glVaos[curTechnique.name()])
}

func (me *MeshBuffer) dispose() {
	for mesh, _ := range me.meshes {
		mesh.meshBuffer, mesh.gpuSynced = nil, false
	}
	me.glIbo.Dispose()
	me.glVbo.Dispose()
	// gl.DeleteBuffers(1, &me.glIbo)
	// gl.DeleteBuffers(1, &me.glVbo)
	for _, glVao := range me.glVaos {
		gl.DeleteVertexArrays(1, &glVao)
	}
}

func (me *MeshBuffer) Remove(mesh *Mesh) {
	if mesh.meshBuffer == me {
		mesh.GpuDelete()
		mesh.meshBuffer = nil
		delete(me.meshes, mesh)
	}
}
