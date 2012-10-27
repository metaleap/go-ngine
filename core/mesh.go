package core

import (
	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/go3d/go-util/gl"
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
	Indices []gl.Uint
	Normals []gl.Float
	Verts []gl.Float
	TexCoords []gl.Float

	raw *tMeshData
	glInit, gpuSynced bool
	glElemBuf, glVnBuf, glVpBuf, glTcBuf gl.Uint
	glMode gl.Enum
	glNumIndices, glNumVerts gl.Sizei
}

func (me *TMesh) GpuDelete () {
	if me.glInit {
		me.gpuSynced = false
		gl.DeleteBuffers(1, &me.glVpBuf)
		gl.DeleteBuffers(1, &me.glElemBuf)
		gl.DeleteBuffers(1, &me.glVnBuf)
		gl.DeleteBuffers(1, &me.glTcBuf)
		me.glVpBuf, me.glElemBuf, me.glVnBuf, me.glTcBuf = 0, 0, 0, 0
	}
}

func (me *TMesh) GpuSync () {
	me.GpuDelete()
	gl.GenBuffers(1, &me.glVpBuf)
	gl.GenBuffers(1, &me.glElemBuf)
	gl.GenBuffers(1, &me.glVnBuf)
	gl.GenBuffers(1, &me.glTcBuf)

	gl.BindBuffer(gl.ARRAY_BUFFER, me.glVpBuf)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.Verts)), gl.Pointer(&me.Verts[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, me.glElemBuf)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.Indices)), gl.Pointer(&me.Indices[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, me.glVnBuf)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.Normals)), gl.Pointer(&me.Normals[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, me.glTcBuf)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.TexCoords)), gl.Pointer(&me.TexCoords[0]), gl.STATIC_DRAW)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	me.gpuSynced = true
}

func (me *TMesh) GpuSynced () bool {
	return me.gpuSynced
}

func (me *TMesh) load (raw *tMeshData) {
	var numVerts = 3 * len(raw.faces)
	var outVerts = make([]gl.Float, Core.MeshBuffer.FloatsPerVertex() * numVerts)
	var outElems = make([]gl.Uint, numVerts)
	var outFaces = make([]*TMeshFace, len(raw.faces))
	var vertsMap = map[gl.Uint]gl.Uint {}
	var voff, ioff, vindex, vhash gl.Uint
	var foff, ei = 0, 0
	var face tMeshFace3
	var ventry tVe
	var vexists bool
	for _, face = range raw.faces {
		outFaces[foff] = newMeshFace()
		for ei, ventry = range face {
			vhash = glutil.Hash3(ventry.posIndex, ventry.texCoordIndex, ventry.normalIndex)
			if _, vexists = vertsMap[vhash]; !vexists {
				vindex, vertsMap[vhash] = voff, voff
				copy(outVerts[voff : (voff + 3)], raw.positions[ventry.posIndex][0 : 3])
				voff += 3
				copy(outVerts[voff : (voff + 2)], raw.texCoords[ventry.texCoordIndex][0 : 2])
				voff += 2
				copy(outVerts[voff : (voff + 3)], raw.normals[ventry.normalIndex][0 : 3])
				voff += 3
			}
			outElems[ioff] = vindex
			outFaces[foff].entries[ei] = ioff
			ioff++
		}
		foff++
	}
}

func (me *TMesh) Loaded () bool {
	return len(me.Verts) > 0
}

func (me *TMesh) render () {
	curMesh = me
	curTechnique.OnRenderMesh()
	gl.BindBuffer(gl.ARRAY_BUFFER, me.glVpBuf)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, me.glElemBuf)
	gl.EnableVertexAttribArray(curProg.AttrLocs["aPos"])
	gl.VertexAttribPointer(curProg.AttrLocs["aPos"], 3, gl.FLOAT, gl.FALSE, 0, gl.Pointer(nil))
	gl.DrawElements(me.glMode, me.glNumIndices, gl.UNSIGNED_INT, gl.Pointer(nil))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (me *TMesh) Unload () {
	me.Verts, me.Normals, me.Indices, me.TexCoords = nil, nil, nil, nil
}
