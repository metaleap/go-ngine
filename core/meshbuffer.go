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
	bufs map[string]*tMeshBuffer
}

	func newMeshBuffers () *tMeshBuffers {
		var meshBuffers = &tMeshBuffers {}
		meshBuffers.bufs = map[string]*tMeshBuffer {}
		return meshBuffers
	}

	func (me *tMeshBuffers) Add (name string, params *tMeshBufferParams) (buf *tMeshBuffer, err error) {
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
		me.bufs = map[string]*tMeshBuffer {}
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

type tMeshBuffer struct {
	MemSizeIndices, MemSizeVertices int32
	Params *tMeshBufferParams

	glIbo, glVbo gl.Uint
	glVaos map[string]gl.Uint
}

	func newMeshBuffer (name string, params *tMeshBufferParams) (mem *tMeshBuffer, err error) {
		var glVao gl.Uint
		mem = &tMeshBuffer {}
		mem.Params = params
		mem.glVaos = map[string]gl.Uint {}
		mem.MemSizeIndices = Core.MeshBuffers.MemSizePerIndex() * params.NumIndices
		mem.MemSizeVertices = Core.MeshBuffers.MemSizePerVertex() * params.NumVerts
		gl.GenBuffers(1, &mem.glVbo)
		gl.GenBuffers(1, &mem.glIbo)
		for techName, _ := range techs {
			gl.GenVertexArrays(1, &glVao)
			mem.glVaos[techName] = glVao
		}
		gl.BindBuffer(gl.ARRAY_BUFFER, mem.glVbo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mem.glIbo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(mem.MemSizeVertices), gl.Pointer(nil), glutil.IfE(params.MostlyStatic, gl.STATIC_DRAW, gl.DYNAMIC_DRAW))
		gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(mem.MemSizeIndices), gl.Pointer(nil), glutil.IfE(params.MostlyStatic, gl.STATIC_DRAW, gl.DYNAMIC_DRAW))
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		if err = glutil.LastError("newMeshBuffer(%v numVerts=%v numIndices=%v)", name, params.NumVerts, params.NumIndices); err != nil {
			mem.dispose()
			mem = nil
		} else {
			for techName, glVao := range mem.glVaos {
				gl.BindVertexArray(glVao)
				gl.BindBuffer(gl.ARRAY_BUFFER, mem.glVbo)
				gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mem.glIbo)
				techs[techName].initMeshBuffer(mem)
				gl.BindVertexArray(0)
				gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
				gl.BindBuffer(gl.ARRAY_BUFFER, 0)
			}
		}
		return
	}

	func (me *tMeshBuffer) dispose () {
		gl.DeleteBuffers(1, &me.glIbo)
		gl.DeleteBuffers(1, &me.glVbo)
		for _, glVao := range me.glVaos {
			gl.DeleteVertexArrays(1, &glVao)
		}
	}
