package core

import (
	"fmt"

	gl "github.com/chsc/gogl/gl42"
)

type tMeshes map[string]*TMesh

	func (me *tMeshes) Load (provider FMeshProvider, args ... interface {}) (*TMesh, error) {
		return provider(args ...)
	}

	func (me *tMeshes) New () *TMesh {
		var mesh = &TMesh {}
		return mesh
	}

type tMeshRaw struct {
	verts []gl.Float
	indices []gl.Uint
	faces []*TMeshFace
}

type TMeshFace struct {
	entries [3]gl.Uint
}

	func newMeshFace () *TMeshFace {
		var face = &TMeshFace {}
		return face
	}

type TMesh struct {
	meshBuffer *tMeshBuffer
	raw *tMeshRaw
	gpuSynced bool
	glNewVao, glNewVbo, glNewIbo gl.Uint
}

	func (me *TMesh) GpuDelete () {
		me.gpuSynced = false
		gl.DeleteBuffers(1, &me.glNewIbo)
		gl.DeleteBuffers(1, &me.glNewVbo)
		gl.DeleteVertexArrays(1, &me.glNewVao)
	}

	func (me *TMesh) GpuSync () {
		me.GpuDelete()
		gl.GenVertexArrays(1, &me.glNewVao)
		gl.GenBuffers(1, &me.glNewIbo)
		gl.GenBuffers(1, &me.glNewVbo)

		gl.BindVertexArray(me.glNewVao)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, me.glNewIbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, me.glNewVbo)
		gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.raw.verts)), gl.Pointer(&me.raw.verts[0]), gl.STATIC_DRAW)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.raw.indices)), gl.Pointer(&me.raw.indices[0]), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

		gl.BindVertexArray(0)
		me.gpuSynced = true
	}

	func (me *TMesh) load (meshData *tMeshData) {
		var numVerts = 3 * len(meshData.faces)
		var vertsMap = map[tVe]gl.Uint {}
		var offsetFloat, offsetIndex, offsetVertex, vindex gl.Uint
		var offsetFace, ei = 0, 0
		var face tMeshFace3
		var ventry tVe
		var vexists bool
		me.gpuSynced = false
		me.raw = &tMeshRaw {}
		me.raw.verts = make([]gl.Float, Core.MeshBuffers.FloatsPerVertex() * numVerts)
		me.raw.indices = make([]gl.Uint, numVerts)
		me.raw.faces = make([]*TMeshFace, len(meshData.faces))
		for _, face = range meshData.faces {
			me.raw.faces[offsetFace] = newMeshFace()
			for ei, ventry = range face {
				if vindex, vexists = vertsMap[ventry]; !vexists {
					vindex, vertsMap[ventry] = offsetVertex, offsetVertex
					copy(me.raw.verts[offsetFloat : (offsetFloat + 3)], meshData.positions[ventry.posIndex][0 : 3])
					offsetFloat += 3
					copy(me.raw.verts[offsetFloat : (offsetFloat + 2)], meshData.texCoords[ventry.texCoordIndex][0 : 2])
					offsetFloat += 2
					copy(me.raw.verts[offsetFloat : (offsetFloat + 3)], meshData.normals[ventry.normalIndex][0 : 3])
					offsetFloat += 3
					offsetVertex++
				}
				me.raw.indices[offsetIndex] = vindex
				me.raw.faces[offsetFace].entries[ei] = offsetIndex
				offsetIndex++
			}
			offsetFace++
		}
		fmt.Printf("meshload() gave %v faces, %v att floats for %v verts, %v indices\n", len(me.raw.faces), len(me.raw.verts), numVerts, len(me.raw.indices))
	}

	func (me *TMesh) Loaded () bool {
		return me.raw != nil
	}

	func (me *TMesh) render () {
		curMesh = me
		gl.BindVertexArray(me.glNewVao)

		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, me.glNewIbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, me.glNewVbo)
		gl.EnableVertexAttribArray(curProg.AttrLocs["aPos"])
		gl.VertexAttribPointer(curProg.AttrLocs["aPos"], 3, gl.FLOAT, gl.FALSE, 8 * 4, gl.Pointer(nil))
		curTechnique.OnRenderMesh()
		gl.DrawElements(gl.TRIANGLES, gl.Sizei(len(me.raw.indices)), gl.UNSIGNED_INT, gl.Pointer(nil))
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

		gl.BindVertexArray(0)
	}

	func (me *TMesh) Unload () {
		me.raw = nil
	}
