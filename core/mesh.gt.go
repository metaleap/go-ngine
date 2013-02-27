package core

import (
	gl "github.com/go3d/go-opengl/core"
)

type Mesh struct {
	ID             int
	DefaultModelID int
	Name           string

	meshBufOffsetBaseIndex, meshBufOffsetIndices, meshBufOffsetVerts int32
	gpuSynced                                                        bool
	meshBuffer                                                       *MeshBuffer
	raw                                                              *meshRaw
}

func (me *Mesh) dispose() {
	me.GpuDelete()
}

func (me *Mesh) init() {
	me.DefaultModelID = -1
}

func (me *Mesh) GpuDelete() {
	if me.gpuSynced {
		me.gpuSynced = false
	}
}

func (me *Mesh) GpuUpload() (err error) {
	sizeVerts, sizeIndices := gl.Sizeiptr(4*len(me.raw.meshVs)), gl.Sizeiptr(4*len(me.raw.indices))
	me.GpuDelete()

	if sizeVerts > gl.Sizeiptr(me.meshBuffer.MemSizeVertices) {
		err = errf("Cannot upload mesh '%v': vertex size (%vB) exceeds mesh buffer's available vertex memory (%vB)", me.Name, sizeVerts, me.meshBuffer.MemSizeVertices)
	} else if sizeIndices > gl.Sizeiptr(me.meshBuffer.MemSizeIndices) {
		err = errf("Cannot upload mesh '%v': index size (%vB) exceeds mesh buffer's available index memory (%vB)", me.Name, sizeIndices, me.meshBuffer.MemSizeIndices)
	} else {
		me.meshBufOffsetBaseIndex, me.meshBufOffsetIndices, me.meshBufOffsetVerts = me.meshBuffer.offsetBaseIndex, me.meshBuffer.offsetIndices, me.meshBuffer.offsetVerts
		Diag.LogMeshes("Upload %v at voff=%v ioff=%v boff=%v", me.Name, me.meshBufOffsetVerts, me.meshBufOffsetIndices, me.meshBufOffsetBaseIndex)
		me.meshBuffer.glIbo.Bind()
		defer me.meshBuffer.glIbo.Unbind()
		me.meshBuffer.glVbo.Bind()
		defer me.meshBuffer.glVbo.Unbind()
		if err = me.meshBuffer.glVbo.SubData(gl.Intptr(me.meshBufOffsetVerts), sizeVerts, gl.Ptr(&me.raw.meshVs[0])); err == nil {
			me.meshBuffer.offsetVerts += int32(sizeVerts)
			if err = me.meshBuffer.glIbo.SubData(gl.Intptr(me.meshBufOffsetIndices), sizeIndices, gl.Ptr(&me.raw.indices[0])); err == nil {
				me.meshBuffer.offsetIndices += int32(sizeIndices)
				me.meshBuffer.offsetBaseIndex += int32(len(me.raw.indices))
				if err == nil {
					me.gpuSynced = true
				}
			}
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
		offsetFloat, offsetIndex, offsetVertex, vindex uint32
		vreuse, offsetFace, ei, numFinalVerts          int
		vexists                                        bool
		ventry                                         MeshV
	)
	vertsMap, numVerts := map[MeshV]uint32{}, 3*int32(len(meshData.Faces))
	me.gpuSynced = false
	me.raw = &meshRaw{}
	me.raw.meshVs = make([]float32, Core.MeshBuffers.FloatsPerVertex()*numVerts)
	me.raw.indices = make([]uint32, numVerts)
	me.raw.faces = make([]*meshRawFace, len(meshData.Faces))
	for _, face := range meshData.Faces {
		me.raw.faces[offsetFace] = newMeshRawFace(face.MeshFaceBase)
		for ei, ventry = range face.V {
			if vindex, vexists = vertsMap[ventry]; !vexists {
				vindex, vertsMap[ventry] = offsetVertex, offsetVertex
				copy(me.raw.meshVs[offsetFloat:(offsetFloat+3)], meshData.Positions[ventry.PosIndex][0:3])
				offsetFloat += 3
				copy(me.raw.meshVs[offsetFloat:(offsetFloat+2)], meshData.TexCoords[ventry.TexCoordIndex][0:2])
				offsetFloat += 2
				copy(me.raw.meshVs[offsetFloat:(offsetFloat+3)], meshData.Normals[ventry.NormalIndex][0:3])
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
	Diag.LogMeshes("mesh{%v}.Load() gave %v faces, %v att floats for %v final verts (%v source verts), %v indices (%vx vertex reuse)", me.Name, len(me.raw.faces), len(me.raw.meshVs), numFinalVerts, numVerts, len(me.raw.indices), vreuse)
	return
}

func (me *Mesh) Loaded() bool {
	return me.raw != nil
}

func (me MeshLib) GpuSync() (err error) {
	var mesh *Mesh
	for id, _ := range Core.Libs.Meshes {
		if mesh = Core.Libs.Meshes.Get(id); mesh != nil && !mesh.gpuSynced {
			if err = mesh.GpuUpload(); err != nil {
				return
			}
		}
	}
	return
}

func (me *MeshLib) AddNewAndLoad(name string, meshProvider MeshProvider, args ...interface{}) (mesh *Mesh, err error) {
	mesh = me.AddNew()
	mesh.Name = name
	if err = mesh.Load(meshProvider, args...); err != nil {
		me.Remove(mesh.ID, 1)
		mesh = nil
	}
	return
}

//#begin-gt -gen-lib.gt T:Mesh L:Meshes

//	Only used for Core.Libs.Meshes.
type MeshLib []Mesh

func (me *MeshLib) AddNew() (ref *Mesh) {
	id := -1
	for i := 0; i < len(*me); i++ {
		if (*me)[i].ID < 0 {
			id = i
			break
		}
	}
	if id < 0 {
		if id = len(*me); id == cap(*me) {
			nu := make(MeshLib, id, id+Options.Libs.GrowCapBy)
			copy(nu, *me)
			*me = nu
		}
		*me = append(*me, Mesh{})
	}
	ref = &(*me)[id]
	ref.ID = id
	ref.init()
	return
}

func (me *MeshLib) Compact() {
	var (
		before, after []Mesh
		ref           *Mesh
		oldID, i      int
	)
	for i = 0; i < len(*me); i++ {
		if (*me)[i].ID < 0 {
			before, after = (*me)[:i], (*me)[i+1:]
			*me = append(before, after...)
		}
	}
	changed := make(map[int]int, len(*me))
	for i = 0; i < len(*me); i++ {
		if (*me)[i].ID != i {
			ref = &(*me)[i]
			oldID, ref.ID = ref.ID, i
			changed[oldID] = i
		}
	}
	if len(changed) > 0 {
		me.onMeshIDsChanged(changed)
		Options.Libs.OnIDsChanged.Meshes(changed)
	}
}

func (me *MeshLib) ctor() {
	*me = make(MeshLib, 0, Options.Libs.InitialCap)
}

func (me *MeshLib) dispose() {
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (me MeshLib) Get(id int) (ref *Mesh) {
	if id > -1 && id < len(me) {
		if ref = &me[id]; ref.ID != id {
			ref = nil
		}
	}
	return
}

func (me MeshLib) Has(id int) (has bool) {
	if id > -1 && id < len(me) {
		has = me[id].ID == id
	}
	return
}

func (me MeshLib) Ok(id int) bool {
	return me[id].ID > -1
}

func (me MeshLib) Remove(fromID, num int) {
	if l := len(me); fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			me[id].dispose()
			changed[id], me[id].ID = -1, -1
		}
		me.onMeshIDsChanged(changed)
		Options.Libs.OnIDsChanged.Meshes(changed)
	}
}

func (me MeshLib) Walk(on func(ref *Mesh)) {
	for id := 0; id < len(me); id++ {
		if me[id].ID > -1 {
			on(&me[id])
		}
	}
}

//#end-gt
