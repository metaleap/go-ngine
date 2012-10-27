package core

import (
	gl "github.com/chsc/gogl/gl42"
)

type tGpuMeshMemory struct {
	NumVerts, NumIndices gl.Sizeiptr

	glElemBuf, glVertBuf gl.Uint
	glElemOffset, glVertOffset int
}

func newGpuMeshMemory (numVerts, numIndices gl.Sizeiptr) *tGpuMeshMemory {
	var mem = &tGpuMeshMemory {}
	mem.NumIndices, mem.NumVerts = numIndices, numVerts
	gl.GenBuffers(1, &mem.glVertBuf)
	gl.GenBuffers(1, &mem.glElemBuf)
	gl.BindBuffer(gl.ARRAY_BUFFER, mem.glVertBuf)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mem.glElemBuf)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, mem.MemSizePerVertex() * numVerts, gl.Pointer(nil), gl.DYNAMIC_DRAW)
	gl.BufferData(gl.ARRAY_BUFFER, mem.MemSizePerIndex() * numIndices, gl.Pointer(nil), gl.DYNAMIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	return mem
}

func (me *tGpuMeshMemory) dispose () {
	gl.DeleteBuffers(1, &me.glElemBuf)
	gl.DeleteBuffers(1, &me.glVertBuf)
}

func (me *tGpuMeshMemory) MemSizePerIndex () gl.Sizeiptr {
	return 4 * 3
}

func (me *tGpuMeshMemory) MemSizePerVertex () gl.Sizeiptr {
	const sizePerFloat gl.Sizeiptr = 4
	const numVertPosFloats gl.Sizeiptr = 3
	const numVertNormalFloats gl.Sizeiptr = 3
	const numVertTexCoordFloats = 2
	return (sizePerFloat * numVertNormalFloats) + (sizePerFloat * numVertNormalFloats) + (sizePerFloat * numVertTexCoordFloats)
}
