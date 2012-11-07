package core

import (
	gl "github.com/chsc/gogl/gl42"
)

type meshFace3 [3]ve

type va2 [2]gl.Float

type va3 [3]gl.Float

type ve struct {
	posIndex, texCoordIndex, normalIndex gl.Uint
}

type meshData struct {
	positions []va3
	texCoords []va2
	normals []va3
	faces []meshFace3
}

	func newMeshData () *meshData {
		var raw = &meshData {}
		raw.positions = []va3 {}
		raw.texCoords = []va2 {}
		raw.normals = []va3 {}
		raw.faces = []meshFace3 {}
		return raw
	}

	func (me *meshData) addFaces (faces ... meshFace3) {
		me.faces = append(me.faces, faces ...)
	}

	func (me *meshData) addPositions (positions ... va3) {
		me.positions = append(me.positions, positions ...)
	}

	func (me *meshData) addNormals (normals ... va3) {
		me.normals = append(me.normals, normals ...)
	}

	func (me *meshData) addTexCoords (texCoords ... va2) {
		me.texCoords = append(me.texCoords, texCoords ...)
	}
