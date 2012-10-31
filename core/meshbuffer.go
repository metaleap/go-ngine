package core

import (
	"fmt"

	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/go3d/go-util/gl"
)

type tMeshBufferParams struct {
	HugeMeshSupport, MostlyStatic, CompressTexCoords, CompressTexCoordsNeg, CompressNormals, CompressPositions bool
	NumVerts, NumIndices int32
}

type tMeshBuffers struct {
	bufs map[string]*TMeshBuffer
}

	func newMeshBuffers () *tMeshBuffers {
		var meshBuffers = &tMeshBuffers {}
		meshBuffers.bufs = map[string]*TMeshBuffer {}
		return meshBuffers
	}

	func (me *tMeshBuffers) Add (name string, params *tMeshBufferParams) (buf *TMeshBuffer, err error) {
		buf = me.bufs[name]
		if buf == nil {
			if buf, err = newMeshBuffer(name, params); err == nil {
				me.bufs[name] = buf
			} else if buf != nil {
				buf.dispose()
			}
		} else {
			err = fmt.Errorf("Cannot add a new mesh buffer with name '%v': already exists", name)
		}
		return
	}

	func (me *tMeshBuffers) dispose () {
		for _, buf := range me.bufs { buf.dispose() }
		me.bufs = map[string]*TMeshBuffer {}
	}

	func (me *tMeshBuffers) FloatsPerVertex () int32 {
		const numVertPosFloats, numVertTexCoordFloats, numVertNormalFloats int32 = 3, 2, 3
		return numVertPosFloats + numVertNormalFloats + numVertTexCoordFloats
	}

	func (me *tMeshBuffers) MemSizePerIndex () int32 {
		return 4
	}

	func (me *tMeshBuffers) MemSizePerVertex () int32 {
		const sizePerFloat int32 = 4
		return sizePerFloat * me.FloatsPerVertex()
	}

	func (me *tMeshBuffers) NewParams (numVerts, numIndices int32) *tMeshBufferParams {
		var params = &tMeshBufferParams {}
		params.MostlyStatic, params.NumIndices, params.NumVerts = true, numIndices, numVerts
		return params
	}

	func (me *tMeshBuffers) Remove (name string) {
		var buf = me.bufs[name]
		if buf != nil { buf.dispose(); delete(me.bufs, name) }
	}

type TMeshBuffer struct {
	MemSizeIndices, MemSizeVertices int32
	Params *tMeshBufferParams

	offsetBaseIndex, offsetIndices, offsetVerts int32
	name string
	meshes tMeshes
	glIbo, glVbo gl.Uint
	glVaos map[string]gl.Uint
}

	func newMeshBuffer (name string, params *tMeshBufferParams) (buf *TMeshBuffer, err error) {
		var glVao gl.Uint
		buf = &TMeshBuffer {}
		buf.name = name
		buf.meshes = tMeshes {}
		buf.Params = params
		buf.glVaos = map[string]gl.Uint {}
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
		gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(buf.MemSizeVertices), gl.Pointer(nil), glutil.IfE(params.MostlyStatic, gl.STATIC_DRAW, gl.DYNAMIC_DRAW))
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(buf.MemSizeIndices), gl.Pointer(nil), glutil.IfE(params.MostlyStatic, gl.STATIC_DRAW, gl.DYNAMIC_DRAW))
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		if err = glutil.LastError("newMeshBuffer(%v numVerts=%v numIndices=%v)", name, params.NumVerts, params.NumIndices); err != nil {
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

	func (me *TMeshBuffer) Add (mesh *TMesh) (err error) {
		if mesh.meshBuffer != nil {
			err = fmt.Errorf("Cannot add mesh '%v' to mesh buffer '%v': already belongs to mesh buffer '%v'.", mesh.name, me.name, mesh.meshBuffer.name)
		} else if me.meshes.Add(mesh) != nil {
			mesh.gpuSynced = false
			mesh.meshBuffer = me
		} else {
			err = fmt.Errorf("Cannot add mesh '%v' to mesh buffer '%v': already has a mesh with that name.", mesh.name, me.name)
		}
		return
	}

	func (me *TMeshBuffer) use () {
		curMeshBuf = me
		gl.BindVertexArray(me.glVaos[curTechnique.name()])
	}

	func (me *TMeshBuffer) dispose () {
		for _, mesh := range me.meshes { mesh.meshBuffer, mesh.gpuSynced = nil, false }
		gl.DeleteBuffers(1, &me.glIbo)
		gl.DeleteBuffers(1, &me.glVbo)
		for _, glVao := range me.glVaos {
			gl.DeleteVertexArrays(1, &glVao)
		}
	}

	func (me *TMeshBuffer) Remove (mesh *TMesh) {
		if mesh.meshBuffer == me {
			mesh.GpuDelete()
			mesh.meshBuffer = nil
			delete(me.meshes, mesh.name)
		}
	}
