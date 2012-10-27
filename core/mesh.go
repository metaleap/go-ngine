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

type TMesh struct {
	OldIndices []gl.Uint
	OldNormals []gl.Float
	OldVerts []gl.Float
	OldTexCoords []gl.Float

	raw *tMeshData
	newVerts []gl.Float
	newElems []gl.Uint
	newFaces []*TMeshFace
	glVao gl.Uint
	glInit, gpuSynced bool
	glNewVbo, glNewIbo gl.Uint
	glOldElemBuf, glOldVnBuf, glOldVpBuf, glOldTcBuf gl.Uint
	glOldMode gl.Enum
	glOldNumIndices, glOldNumVerts gl.Sizei
}

func (me *TMesh) GpuDelete () {
	if me.glInit {
		me.gpuSynced = false
		gl.DeleteBuffers(1, &me.glOldVpBuf)
		gl.DeleteBuffers(1, &me.glOldElemBuf)
		gl.DeleteBuffers(1, &me.glOldVnBuf)
		gl.DeleteBuffers(1, &me.glOldTcBuf)
		gl.DeleteBuffers(1, &me.glNewIbo)
		gl.DeleteBuffers(1, &me.glNewVbo)
		gl.DeleteVertexArrays(1, &me.glVao)
		me.glOldVpBuf, me.glOldElemBuf, me.glOldVnBuf, me.glOldTcBuf = 0, 0, 0, 0
	}
}

func (me *TMesh) GpuSync () {
	me.GpuDelete()
	gl.GenVertexArrays(1, &me.glVao)
	gl.GenBuffers(1, &me.glOldVpBuf)
	gl.GenBuffers(1, &me.glOldElemBuf)
	gl.GenBuffers(1, &me.glOldVnBuf)
	gl.GenBuffers(1, &me.glOldTcBuf)
	gl.GenBuffers(1, &me.glNewIbo)
	gl.GenBuffers(1, &me.glNewVbo)

	gl.BindVertexArray(me.glVao)
	if me.raw == nil {
		gl.BindBuffer(gl.ARRAY_BUFFER, me.glOldVpBuf)
		gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.OldVerts)), gl.Pointer(&me.OldVerts[0]), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, me.glOldElemBuf)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.OldIndices)), gl.Pointer(&me.OldIndices[0]), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		gl.BindBuffer(gl.ARRAY_BUFFER, me.glOldVnBuf)
		gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.OldNormals)), gl.Pointer(&me.OldNormals[0]), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, me.glOldTcBuf)
		gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.OldTexCoords)), gl.Pointer(&me.OldTexCoords[0]), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, me.glNewIbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, me.glNewVbo)
	if me.raw != nil {
		gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.newVerts)), gl.Pointer(&me.newVerts[0]), gl.STATIC_DRAW)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.newElems)), gl.Pointer(&me.newElems[0]), gl.STATIC_DRAW)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	gl.BindVertexArray(0)
	me.gpuSynced = true
}

func (me *TMesh) GpuSynced () bool {
	return me.gpuSynced
}

func (me *TMesh) load (raw *tMeshData) {
	var numVerts = 3 * len(raw.faces)
	var vertsMap = map[tVe]gl.Uint {}
	var offsetFloat, offsetIndex, offsetVertex, vindex gl.Uint
	var offsetFace, ei = 0, 0
	var face tMeshFace3
	var ventry tVe
	var vexists bool
	me.raw, me.gpuSynced = raw, false
	me.newVerts = make([]gl.Float, Core.MeshBuffer.FloatsPerVertex() * numVerts)
	me.newElems = make([]gl.Uint, numVerts)
	me.newFaces = make([]*TMeshFace, len(raw.faces))
	for _, face = range raw.faces {
		me.newFaces[offsetFace] = newMeshFace()
		for ei, ventry = range face {
			if vindex, vexists = vertsMap[ventry]; !vexists {
				vindex, vertsMap[ventry] = offsetVertex, offsetVertex
				copy(me.newVerts[offsetFloat : (offsetFloat + 3)], raw.positions[ventry.posIndex][0 : 3])
				offsetFloat += 3
				copy(me.newVerts[offsetFloat : (offsetFloat + 2)], raw.texCoords[ventry.texCoordIndex][0 : 2])
				offsetFloat += 2
				copy(me.newVerts[offsetFloat : (offsetFloat + 3)], raw.normals[ventry.normalIndex][0 : 3])
				offsetFloat += 3
				offsetVertex++
			}
			me.newElems[offsetIndex] = vindex
			me.newFaces[offsetFace].entries[ei] = offsetIndex
			offsetIndex++
		}
		offsetFace++
	}
	fmt.Printf("meshload() gave %v faces, %v att floats for %v verts, %v indices\n", len(me.newFaces), len(me.newVerts), numVerts, len(me.newElems))
}

func (me *TMesh) Loaded () bool {
	return (len(me.OldVerts) > 0) || (me.raw != nil)
}

func (me *TMesh) render () {
	curMesh = me
	gl.BindVertexArray(me.glVao)
	if me.raw != nil {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, me.glNewIbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, me.glNewVbo)
		gl.EnableVertexAttribArray(curProg.AttrLocs["aPos"])
		gl.VertexAttribPointer(curProg.AttrLocs["aPos"], 3, gl.FLOAT, gl.FALSE, 8 * 4, gl.Pointer(nil))
		curTechnique.OnRenderMesh()
		gl.DrawElements(gl.TRIANGLES, gl.Sizei(len(me.newElems)), gl.UNSIGNED_INT, gl.Pointer(nil))
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	} else {
		curTechnique.OnRenderMesh()
		gl.BindBuffer(gl.ARRAY_BUFFER, me.glOldVpBuf)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, me.glOldElemBuf)
		gl.EnableVertexAttribArray(curProg.AttrLocs["aPos"])
		gl.VertexAttribPointer(curProg.AttrLocs["aPos"], 3, gl.FLOAT, gl.FALSE, 0, gl.Pointer(nil))
		gl.DrawElements(me.glOldMode, me.glOldNumIndices, gl.UNSIGNED_INT, gl.Pointer(nil))
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}
	gl.BindVertexArray(0)
}

func (me *TMesh) Unload () {
	me.OldVerts, me.OldNormals, me.OldIndices, me.OldTexCoords = nil, nil, nil, nil
}
