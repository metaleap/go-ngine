package core

import (
	gl "github.com/go3d/go-opengl/gogl"

	ugl "github.com/go3d/go-opengl/util"
)

type Mesh struct {
	Models Models

	meshBufOffsetBaseIndex, meshBufOffsetIndices, meshBufOffsetVerts int32
	gpuSynced                                                        bool
	id                                                               string
	meshBuffer                                                       *MeshBuffer
	raw                                                              *meshRaw
}

func (me *Mesh) dispose() {
	me.GpuDelete()
	me.Unload()
}

func (me *Mesh) init() {
	me.Models = Models{"": newModel("", me)}
}

func (me *Mesh) GpuDelete() {
	if me.gpuSynced {
		me.gpuSynced = false
	}
}

func (me *Mesh) GpuUpload() (err error) {
	sizeVerts, sizeIndices := gl.Sizeiptr(4*len(me.raw.meshVerts)), gl.Sizeiptr(4*len(me.raw.indices))
	me.GpuDelete()

	if sizeVerts > gl.Sizeiptr(me.meshBuffer.MemSizeVertices) {
		err = fmtErr("Cannot upload mesh '%v': vertex size (%vB) exceeds mesh buffer's available vertex memory (%vB)", me.id, sizeVerts, me.meshBuffer.MemSizeVertices)
	} else if sizeIndices > gl.Sizeiptr(me.meshBuffer.MemSizeIndices) {
		err = fmtErr("Cannot upload mesh '%v': index size (%vB) exceeds mesh buffer's available index memory (%vB)", me.id, sizeIndices, me.meshBuffer.MemSizeIndices)
	} else {
		me.meshBufOffsetBaseIndex, me.meshBufOffsetIndices, me.meshBufOffsetVerts = me.meshBuffer.offsetBaseIndex, me.meshBuffer.offsetIndices, me.meshBuffer.offsetVerts
		Diag.LogMeshes("Upload %v at voff=%v ioff=%v boff=%v", me.id, me.meshBufOffsetVerts, me.meshBufOffsetIndices, me.meshBufOffsetBaseIndex)
		me.meshBuffer.glIbo.Bind()
		me.meshBuffer.glVbo.Bind()
		// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, me.meshBuffer.glIbo)
		// gl.BindBuffer(gl.ARRAY_BUFFER, me.meshBuffer.glVbo)
		// gl.BufferSubData(gl.ARRAY_BUFFER, gl.Intptr(me.meshBufOffsetVerts), sizeVerts, gl.Pointer(&me.raw.meshVerts[0]))
		me.meshBuffer.glVbo.Update(gl.Intptr(me.meshBufOffsetVerts), sizeVerts, gl.Pointer(&me.raw.meshVerts[0]))
		me.meshBuffer.offsetVerts += int32(sizeVerts)
		// gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, gl.Intptr(me.meshBufOffsetIndices), sizeIndices, gl.Pointer(&me.raw.indices[0]))
		me.meshBuffer.glIbo.Update(gl.Intptr(me.meshBufOffsetIndices), sizeIndices, gl.Pointer(&me.raw.indices[0]))
		me.meshBuffer.offsetIndices += int32(sizeIndices)
		me.meshBuffer.offsetBaseIndex += int32(len(me.raw.indices))
		me.meshBuffer.glIbo.Unbind()
		me.meshBuffer.glVbo.Unbind()
		// gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		if err = ugl.LastError("mesh[%v].GpuUpload()", me.id); err == nil {
			me.gpuSynced = true
		}
	}
	return
}

func (me *Mesh) GpuUploaded() bool {
	return me.gpuSynced
}

func (me *Mesh) Load(provider MeshProvider, args ...interface{}) (err error) {
	var meshData *MeshData
	if meshData, err = provider(args...); err == nil {
		err = me.load(meshData)
	}
	return
}

func (me *Mesh) load(meshData *MeshData) (err error) {
	var (
		numVerts                                       = 3 * int32(len(meshData.Faces))
		numFinalVerts                                  = 0
		offsetFace, ei                                 = 0, 0
		vertsMap                                       = map[MeshVert]uint32{}
		offsetFloat, offsetIndex, offsetVertex, vindex uint32
		vexists                                        bool
		vreuse                                         int
		ventry                                         MeshVert
	)
	me.gpuSynced = false
	me.raw = &meshRaw{}
	me.raw.meshVerts = make([]float32, Core.MeshBuffers.FloatsPerVertex()*numVerts)
	me.raw.indices = make([]uint32, numVerts)
	me.raw.faces = make([]*meshRawFace, len(meshData.Faces))
	for _, face := range meshData.Faces {
		me.raw.faces[offsetFace] = newMeshRawFace(&face.MeshFaceBase)
		for ei, ventry = range face.V {
			if vindex, vexists = vertsMap[ventry]; !vexists {
				vindex, vertsMap[ventry] = offsetVertex, offsetVertex
				copy(me.raw.meshVerts[offsetFloat:(offsetFloat+3)], meshData.Positions[ventry.PosIndex][0:3])
				offsetFloat += 3
				copy(me.raw.meshVerts[offsetFloat:(offsetFloat+2)], meshData.TexCoords[ventry.TexCoordIndex][0:2])
				offsetFloat += 2
				copy(me.raw.meshVerts[offsetFloat:(offsetFloat+3)], meshData.Normals[ventry.NormalIndex][0:3])
				offsetFloat += 3
				offsetVertex++
				numFinalVerts++
			} else {
				vreuse++
			}
			me.raw.indices[offsetIndex] = vindex
			me.raw.faces[offsetFace].entries[ei] = offsetIndex
			offsetIndex++
		}
		offsetFace++
	}
	Diag.LogMeshes("mesh{%v}.Load() gave %v faces, %v att floats for %v final verts (%v source verts), %v indices (%vx vertex reuse)", me.id, len(me.raw.faces), len(me.raw.meshVerts), numFinalVerts, numVerts, len(me.raw.indices), vreuse)
	return
}

func (me *Mesh) Loaded() bool {
	return me.raw != nil
}

func (me *Mesh) Unload() {
	me.raw = nil
}

//	Initializes and returns a new Mesh with default parameters.
func NewMesh(id string) (me *Mesh) {
	me = &Mesh{id: id}
	me.init()
	return
}

//	A hash-table of Meshs associated by IDs. Only for use in Core.Libs.
type LibMeshes map[string]*Mesh

//	Creates and initializes a new Mesh with default parameters,
//	adds it to me under the specified ID, and returns it.
func (me LibMeshes) AddNew(id string) (obj *Mesh) {
	obj = NewMesh(id)
	me[id] = obj
	return
}

func (me LibMeshes) AddLoad(id string, meshProvider MeshProvider, args ...interface{}) (mesh *Mesh, err error) {
	if me[id] == nil {
		mesh = me.AddNew(id)
		if err = mesh.Load(meshProvider, args...); err != nil {
			delete(me, id)
			mesh.dispose()
			mesh = nil
		}
	}
	return
}

func (me *LibMeshes) ctor() {
	*me = LibMeshes{}
}

func (me *LibMeshes) dispose() {
	for _, o := range *me {
		o.dispose()
	}
	me.ctor()
}
