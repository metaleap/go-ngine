package core

import (
	"fmt"
	gl "github.com/chsc/gogl/gl42"
)

type tMeshes map[string]*TMesh

	func (me *tMeshes) New () *TMesh {
		var mesh = &TMesh {}
		return mesh
	}

type TMesh struct {
	Indices []gl.Uint
	Normals []gl.Float
	Verts []gl.Float

	glInit, glSynced bool
	glElemBuf, glNormalBuf, glVertBuf gl.Uint
	glMode gl.Enum
	glNumIndices, glNumVerts gl.Sizei
}

func (me *TMesh) Dispose () {
	if me.glInit {
		me.glSynced, me.glInit = false, false
		gl.DeleteBuffers(1, &me.glVertBuf)
		if me.glElemBuf > 0 {
			gl.DeleteBuffers(1, &me.glElemBuf)
		}
		if me.glNormalBuf > 0 {
			gl.DeleteBuffers(1, &me.glNormalBuf)
		}
	}
}

func (me *TMesh) initBuffer () {
	if !me.glInit {
		me.glSynced, me.glInit = false, true
		gl.GenBuffers(1, &me.glVertBuf)
		if len(me.Indices) > 0 {
			gl.GenBuffers(1, &me.glElemBuf)
		}
		if len(me.Normals) > 0 {
			gl.GenBuffers(1, &me.glNormalBuf)
		}
	}
}

func (me *TMesh) render () {
	curMesh = me
	curTechnique.OnRenderMesh()
	gl.BindBuffer(gl.ARRAY_BUFFER, me.glVertBuf)
	if me.glElemBuf > 0 {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, me.glElemBuf)
		gl.VertexAttribPointer(curProg.AttrLocs["aPos"], 3, gl.FLOAT, gl.FALSE, 0, gl.Pointer(nil))
		gl.DrawElements(me.glMode, me.glNumIndices, gl.UNSIGNED_INT, gl.Pointer(nil))
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	} else {
		gl.VertexAttribPointer(curProg.AttrLocs["aPos"], 3, gl.FLOAT, gl.FALSE, 0, gl.Pointer(nil))
		if camFirst { fmt.Printf("NV=%v --- ", me.glMode) }
		gl.DrawArrays(me.glMode, 0, me.glNumVerts)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
}

func (me *TMesh) updateBuffer () {
	gl.BindBuffer(gl.ARRAY_BUFFER, me.glVertBuf)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.Verts)), gl.Pointer(&me.Verts[0]), gl.STATIC_DRAW)
	if me.glElemBuf > 0 {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, me.glElemBuf)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.Indices)), gl.Pointer(&me.Indices[0]), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	if me.glNormalBuf > 0 {
		gl.BindBuffer(gl.ARRAY_BUFFER, me.glNormalBuf)
		gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(4 * len(me.Normals)), gl.Pointer(&me.Normals[0]), gl.STATIC_DRAW)
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}
	me.glSynced = true
}
