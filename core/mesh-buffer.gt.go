package core

import (
	gl "github.com/go3d/go-opengl/core"

	"github.com/go-utils/uslice"
	ugl "github.com/go3d/go-opengl/util"
)

type MeshBuffer struct {
	Name string

	memSizeIndices, memSizeVertices             int32
	offsetBaseIndex, offsetIndices, offsetVerts int32

	glIbo, glVbo ugl.Buffer
	glVaos       []ugl.VertexArray
	meshIDs      []int
}

func newMeshBuffer(name string, capacity int32) (me *MeshBuffer, err error) {
	me = &MeshBuffer{}
	me.Name = name
	me.meshIDs = make([]int, 0, 256)
	me.glVaos = make([]ugl.VertexArray, 16)
	numVerts, numIndices := capacity, capacity
	me.memSizeIndices = Core.Mesh.Buffers.MemSizePerIndex() * numIndices
	me.memSizeVertices = Core.Mesh.Buffers.MemSizePerVertex() * numVerts
	if err = me.glVbo.Recreate(gl.ARRAY_BUFFER, gl.Sizeiptr(me.memSizeVertices), ugl.PtrNil, gl.STATIC_DRAW); err == nil {
		err = me.glIbo.Recreate(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(me.memSizeIndices), ugl.PtrNil, gl.STATIC_DRAW)
	}
	// if err == nil {
	// 	var ok bool
	// 	for i := 0; i < len(ogl.progs.All); i++ {
	// 		if err = me.setupVao(i); err != nil {
	// 			break
	// 		}
	// 	}
	// }
	if err != nil {
		me.dispose()
		me = nil
	}
	return
}

func (me *MeshBuffer) init() {
}

func (me *MeshBuffer) dispose() {
	for _, meshID := range me.meshIDs {
		if Core.Libs.Meshes.IsOk(meshID) {
			Core.Libs.Meshes[meshID].meshBuffer, Core.Libs.Meshes[meshID].gpuSynced = nil, false
		}
	}
	me.glIbo.Dispose()
	me.glVbo.Dispose()
	for i := 0; i < len(me.glVaos); i++ {
		me.glVaos[i].Dispose()
	}
	me.glVaos = nil
}

func (me *MeshBuffer) Add(meshID int) (err error) {
	if mesh := Core.Libs.Meshes.get(meshID); mesh != nil && mesh.meshBuffer != me {
		if mesh.meshBuffer != nil {
			err = errf("Cannot add mesh '%v' to mesh buffer '%v': already belongs to mesh buffer '%v'.", mesh.Name, me.Name, mesh.meshBuffer.Name)
		} else {
			uslice.IntAppendUnique(&me.meshIDs, meshID)
			mesh.gpuSynced, mesh.meshBuffer = false, me
		}
	}
	return
}

func (me *MeshBuffer) use() {
	if thrRend.curProg.Index >= len(me.glVaos) || me.glVaos[thrRend.curProg.Index].GlHandle == 0 {
		me.setupVao(thrRend.curProg.Index)
	}
	me.glVaos[thrRend.curProg.Index].Bind()
}

func (me *MeshBuffer) Remove(meshID int) {
	if mesh := Core.Libs.Meshes.get(meshID); mesh != nil && mesh.meshBuffer == me {
		mesh.GpuDelete()
		mesh.meshBuffer = nil
	}
	for i, mid := range me.meshIDs {
		if mid == meshID {
			before, after := me.meshIDs[:i], me.meshIDs[i+1:]
			me.meshIDs = append(before, after...)
			break
		}
	}
}

func (me *MeshBuffer) setupVao(progIndex int) (err error) {
	if progIndex >= len(me.glVaos) {
		nuVaos := make([]ugl.VertexArray, len(me.glVaos)+16)
		copy(nuVaos, me.glVaos)
		me.glVaos = nuVaos
	}
	if err = me.glVaos[progIndex].Create(); err == nil {
		if sceneTech, ok := ogl.progs.All[progIndex].Tag.(*RenderTechniqueScene); ok {
			if err = me.glVaos[progIndex].Setup(&ogl.progs.All[progIndex], sceneTech.vertexAttribPointers(), &me.glVbo, &me.glIbo); err != nil {
				return
			}
		}
	}
	return
}

func (me *MeshBufferLib) AddNew(name string, capacity int32) (buf *MeshBuffer, err error) {
	if buf, err = newMeshBuffer(name, capacity); err == nil {
		me.add(buf)
	} else if buf != nil {
		buf.dispose()
		buf = nil
	}
	return
}

func (_ MeshBufferLib) FloatsPerVertex() int32 {
	const numVertPosFloats, numVertTexCoordFloats, numVertNormalFloats int32 = 3, 2, 3
	return numVertPosFloats + numVertNormalFloats + numVertTexCoordFloats
}

func (_ MeshBufferLib) MemSizePerIndex() int32 {
	return 4
}

func (_ MeshBufferLib) MemSizePerVertex() int32 {
	const sizePerFloat int32 = 4
	return sizePerFloat * Core.Mesh.Buffers.FloatsPerVertex()
}











//#begin-gt -gen-reflib.gt T:MeshBuffer L:Core.Mesh.Buffers

//	Only used for Core.Mesh.Buffers
type MeshBufferLib []*MeshBuffer

func (me *MeshBufferLib) add(ref *MeshBuffer) {
	*me = append(*me, ref)
	ref.init()
	return
}

func (me *MeshBufferLib) init() {
	*me = make(MeshBufferLib, 0, 4)
}

func (me *MeshBufferLib) dispose() {
	me.Remove(0, 0)
}

func (me MeshBufferLib) IsOk(id int) bool {
	return id > -1 && id < len(me)
}

func (me *MeshBufferLib) Remove(fromID, num int) {
	if l := len(*me); fromID > -1 && fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		for i := fromID; i < fromID+num; i++ {
			(*me)[fromID].dispose()
		}
		before, after := (*me)[:fromID], (*me)[fromID+num:]
		*me = append(before, after...)
	}
}

func (me MeshBufferLib) Walk(on func(ref *MeshBuffer)) {
	for id := 0; id < len(me); id++ {
		on(me[id])
	}
}

//#end-gt
