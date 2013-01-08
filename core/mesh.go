package core

import (
	"fmt"

	gl "github.com/chsc/gogl/gl42"

	ugl "github.com/go3d/go-glutil"
)

type Meshes map[string]*Mesh

func (me Meshes) Add(mesh *Mesh) *Mesh {
	if me[mesh.name] == nil {
		me[mesh.name] = mesh
		return mesh
	}
	return nil
}

func (me Meshes) AddRange(meshes ...*Mesh) {
	for _, m := range meshes {
		me.Add(m)
	}
}

func (me Meshes) Load(name string, provider MeshProvider, args ...interface{}) (mesh *Mesh, err error) {
	var meshData *MeshData
	mesh = me.New(name)
	if meshData, err = provider(args...); err == nil {
		mesh.load(meshData)
	} else {
		mesh = nil
	}
	return
}

func (me Meshes) New(name string) (mesh *Mesh) {
	mesh = &Mesh{name: name, Models: models{}}
	return mesh
}

type Mesh struct {
	Models models

	name                                                             string
	meshBuffer                                                       *MeshBuffer
	meshBufOffsetBaseIndex, meshBufOffsetIndices, meshBufOffsetVerts int32
	raw                                                              *meshRaw
	gpuSynced                                                        bool
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
		err = fmt.Errorf("Cannot upload mesh '%v': vertex size (%vB) exceeds mesh buffer's available vertex memory (%vB)", me.name, sizeVerts, me.meshBuffer.MemSizeVertices)
	} else if sizeIndices > gl.Sizeiptr(me.meshBuffer.MemSizeIndices) {
		err = fmt.Errorf("Cannot upload mesh '%v': index size (%vB) exceeds mesh buffer's available index memory (%vB)", me.name, sizeIndices, me.meshBuffer.MemSizeIndices)
	} else {
		me.meshBufOffsetBaseIndex, me.meshBufOffsetIndices, me.meshBufOffsetVerts = me.meshBuffer.offsetBaseIndex, me.meshBuffer.offsetIndices, me.meshBuffer.offsetVerts
		fmt.Printf("Upload %v at voff=%v ioff=%v boff=%v\n", me.name, me.meshBufOffsetVerts, me.meshBufOffsetIndices, me.meshBufOffsetBaseIndex)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, me.meshBuffer.glIbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, me.meshBuffer.glVbo)
		gl.BufferSubData(gl.ARRAY_BUFFER, gl.Intptr(me.meshBufOffsetVerts), sizeVerts, gl.Pointer(&me.raw.meshVerts[0]))
		me.meshBuffer.offsetVerts += int32(sizeVerts)
		gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, gl.Intptr(me.meshBufOffsetIndices), sizeIndices, gl.Pointer(&me.raw.indices[0]))
		me.meshBuffer.offsetIndices += int32(sizeIndices)
		me.meshBuffer.offsetBaseIndex += int32(len(me.raw.indices))
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
		if err = ugl.LastError("mesh[%v].GpuUpload()", me.name); err == nil {
			me.gpuSynced = true
		}
	}
	return
}

func (me *Mesh) GpuUploaded() bool {
	return me.gpuSynced
}

func (me *Mesh) load(meshData *MeshData) {
	var (
		numVerts                                       = 3 * int32(len(meshData.faces))
		numFinalVerts                                  = 0
		offsetFace, ei                                 = 0, 0
		vertsMap                                       = map[meshVert]uint32{}
		offsetFloat, offsetIndex, offsetVertex, vindex uint32
		vexists                                        bool
		vreuse                                         int
		ventry                                         meshVert
	)
	me.Models = models{}
	me.gpuSynced = false
	me.raw = &meshRaw{}
	me.raw.meshVerts = make([]float32, Core.MeshBuffers.FloatsPerVertex()*numVerts)
	me.raw.indices = make([]uint32, numVerts)
	me.raw.faces = make([]*meshRawFace, len(meshData.faces))
	for _, face := range meshData.faces {
		me.raw.faces[offsetFace] = newMeshRawFace()
		for ei, ventry = range face {
			if vindex, vexists = vertsMap[ventry]; !vexists {
				vindex, vertsMap[ventry] = offsetVertex, offsetVertex
				copy(me.raw.meshVerts[offsetFloat:(offsetFloat+3)], meshData.positions[ventry.posIndex][0:3])
				offsetFloat += 3
				copy(me.raw.meshVerts[offsetFloat:(offsetFloat+2)], meshData.texCoords[ventry.texCoordIndex][0:2])
				offsetFloat += 2
				copy(me.raw.meshVerts[offsetFloat:(offsetFloat+3)], meshData.normals[ventry.normalIndex][0:3])
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
	me.Models[""] = newModel("", me)
	fmt.Printf("meshload(%v) gave %v faces, %v att floats for %v final verts (%v source verts), %v indices (%vx vertex reuse)\n", me.name, len(me.raw.faces), len(me.raw.meshVerts), numFinalVerts, numVerts, len(me.raw.indices), vreuse)
}

func (me *Mesh) Loaded() bool {
	return me.raw != nil
}

func (me *Mesh) render() {
	if curMeshBuf != me.meshBuffer {
		me.meshBuffer.use()
	}
	curTechnique.onRenderMesh()
	gl.DrawElementsBaseVertex(gl.TRIANGLES, gl.Sizei(len(me.raw.indices)), gl.UNSIGNED_INT, gl.Offset(nil, uintptr(me.meshBufOffsetIndices)), gl.Int(me.meshBufOffsetBaseIndex))
	// gl.DrawElements(gl.TRIANGLES, gl.Sizei(len(me.raw.indices)), gl.UNSIGNED_INT, gl.Pointer(nil))
}

func (me *Mesh) Unload() {
	me.raw = nil
}
