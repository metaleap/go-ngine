package core

import (
	gl "github.com/go3d/go-opengl/core"

	ugl "github.com/go3d/go-opengl/util"
)

type MeshBufferParams struct {
	NumVerts, NumIndices int32
}

type MeshBuffer struct {
	MemSizeIndices, MemSizeVertices int32
	Params                          MeshBufferParams
	Name                            string

	offsetBaseIndex, offsetIndices, offsetVerts int32
	glIbo, glVbo                                ugl.Buffer
	glVaos                                      map[*ugl.Program]*ugl.VertexArray
	meshes                                      map[*Mesh]bool
}

func newMeshBuffer(name string, params MeshBufferParams) (me *MeshBuffer, err error) {
	me = &MeshBuffer{}
	me.Name = name
	me.meshes = map[*Mesh]bool{}
	me.Params = params
	me.glVaos = map[*ugl.Program]*ugl.VertexArray{}
	me.MemSizeIndices = Core.MeshBuffers.MemSizePerIndex() * params.NumIndices
	me.MemSizeVertices = Core.MeshBuffers.MemSizePerVertex() * params.NumVerts
	if err = me.glVbo.Recreate(gl.ARRAY_BUFFER, gl.Sizeiptr(me.MemSizeVertices), ugl.PtrNil, gl.STATIC_DRAW); err == nil {
		err = me.glIbo.Recreate(gl.ELEMENT_ARRAY_BUFFER, gl.Sizeiptr(me.MemSizeIndices), ugl.PtrNil, gl.STATIC_DRAW)
	}
	if err == nil {
		var tech RenderTechnique
		var ok bool
		for i := 0; i < len(ogl.progs.All); i++ {
			if tech, ok = ogl.progs.All[i].Tag.(RenderTechnique); ok {
				if err = me.setupVao(&ogl.progs.All[i], tech); err != nil {
					break
				}
			}
		}
	}
	if err != nil {
		me.dispose()
		me = nil
	}
	return
}

func (me *MeshBuffer) init() {
}

func (me *MeshBuffer) dispose() {
	for mesh, _ := range me.meshes {
		mesh.meshBuffer, mesh.gpuSynced = nil, false
	}
	me.glIbo.Dispose()
	me.glVbo.Dispose()
	for _, glVao := range me.glVaos {
		glVao.Dispose()
	}
	me.glVaos = nil
}

func (me *MeshBuffer) Add(mesh *Mesh) (err error) {
	if mesh.meshBuffer != nil {
		err = errf("Cannot add mesh '%v' to mesh buffer '%v': already belongs to mesh buffer '%v'.", mesh.Name, me.Name, mesh.meshBuffer.Name)
	} else if !me.meshes[mesh] {
		me.meshes[mesh] = true
		mesh.gpuSynced = false
		mesh.meshBuffer = me
	} else {
		err = errf("Cannot add mesh '%v' to mesh buffer '%v': already added.", mesh.Name, me.Name)
	}
	return
}

func (me *MeshBuffer) use() {
	me.glVaos[thrRend.curProg].Bind()
}

func (me *MeshBuffer) Remove(mesh *Mesh) {
	if mesh.meshBuffer == me {
		mesh.GpuDelete()
		mesh.meshBuffer = nil
		delete(me.meshes, mesh)
	}
}

func (me *MeshBuffer) setupVao(prog *ugl.Program, tech RenderTechnique) (err error) {
	vao := &ugl.VertexArray{}
	if err = vao.Create(); err == nil {
		if sceneTech, ok := tech.(*RenderTechniqueScene); ok {
			if err = vao.Setup(prog, sceneTech.vertexAttribPointers(me), &me.glVbo, &me.glIbo); err != nil {
				vao.Dispose()
				vao = nil
			}
		}
		if vao != nil {
			me.glVaos[prog] = vao
		}
	}
	return
}

func (me *MeshBufferLib) AddNew(name string, params MeshBufferParams) (buf *MeshBuffer, err error) {
	if buf, err = newMeshBuffer(name, params); err == nil {
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
	return sizePerFloat * Core.MeshBuffers.FloatsPerVertex()
}

//#begin-gt -gen-reflib.gt T:MeshBuffer L:Core.MeshBuffers

//	Only used for Core.MeshBuffers
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

func (me MeshBufferLib) Get(id int) (ref *MeshBuffer) {
	if me.IsOk(id) {
		ref = me[id]
	}
	return
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
