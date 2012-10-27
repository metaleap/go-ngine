package core

import (
	gl "github.com/chsc/gogl/gl42"
)

type TMeshFace struct {
	entries [3]gl.Uint
}

func newMeshFace () *TMeshFace {
	var face = &TMeshFace {}
	return face
}
