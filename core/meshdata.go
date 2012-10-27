package core

import (
	gl "github.com/chsc/gogl/gl42"
)

type tMeshFace3 [3]tVe

type tVa2 [2]gl.Float

type tVa3 [3]gl.Float

type tVe struct {
	posIndex, texCoordIndex, normalIndex gl.Uint
}

type tMeshData struct {
	positions []tVa3
	texCoords []tVa2
	normals []tVa3
	faces []tMeshFace3
}

	func newMeshData () *tMeshData {
		var raw = &tMeshData {}
		raw.positions = []tVa3 {}
		raw.texCoords = []tVa2 {}
		raw.normals = []tVa3 {}
		raw.faces = []tMeshFace3 {}
		return raw
	}

	func (me *tMeshData) addFaces (faces ... tMeshFace3) {
		me.faces = append(me.faces, faces ...)
	}

	func (me *tMeshData) addPositions (positions ... tVa3) {
		me.positions = append(me.positions, positions ...)
	}

	func (me *tMeshData) addNormals (normals ... tVa3) {
		me.normals = append(me.normals, normals ...)
	}

	func (me *tMeshData) addTexCoords (texCoords ... tVa2) {
		me.texCoords = append(me.texCoords, texCoords ...)
	}
