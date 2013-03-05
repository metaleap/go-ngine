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
	raw                                                              meshRaw
	meshBuffer                                                       *MeshBuffer
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
	sizeVerts, sizeIndices := gl.Sizeiptr(4*len(me.raw.verts)), gl.Sizeiptr(4*len(me.raw.indices))
	me.GpuDelete()
	if sizeVerts > gl.Sizeiptr(me.meshBuffer.memSizeVertices) {
		err = errf("Cannot upload mesh '%v': vertex size (%vB) exceeds mesh buffer's available vertex memory (%vB)", me.Name, sizeVerts, me.meshBuffer.memSizeVertices)
	} else if sizeIndices > gl.Sizeiptr(me.meshBuffer.memSizeIndices) {
		err = errf("Cannot upload mesh '%v': index size (%vB) exceeds mesh buffer's available index memory (%vB)", me.Name, sizeIndices, me.meshBuffer.memSizeIndices)
	} else {
		me.meshBufOffsetBaseIndex, me.meshBufOffsetIndices, me.meshBufOffsetVerts = me.meshBuffer.offsetBaseIndex, me.meshBuffer.offsetIndices, me.meshBuffer.offsetVerts
		Diag.LogMeshes("Upload %v at voff=%v ioff=%v boff=%v", me.Name, me.meshBufOffsetVerts, me.meshBufOffsetIndices, me.meshBufOffsetBaseIndex)
		me.meshBuffer.glIbo.Bind()
		defer me.meshBuffer.glIbo.Unbind()
		me.meshBuffer.glVbo.Bind()
		defer me.meshBuffer.glVbo.Unbind()
		if err = me.meshBuffer.glVbo.SubData(gl.Intptr(me.meshBufOffsetVerts), sizeVerts, gl.Ptr(&me.raw.verts[0])); err == nil {
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

func (me *Mesh) Load(provider MeshProvider) (err error) {
	var meshData *MeshDescriptor
	if meshData, err = provider(); err == nil && meshData != nil {
		err = me.load(meshData)
	}
	return
}

func (me *Mesh) load(meshData *MeshDescriptor) (err error) {
	var (
		offsetFloat, offsetIndex, offsetVertex, vindex uint32
		vreuse, offsetFace, ei, numFinalVerts          int
		vexists                                        bool
		ventry                                         MeshDescF3V
		tvp1, tvp2, tvp3                               MeshDescVA3
	)
	numVerts := 3 * int32(len(meshData.Faces))
	vertsMap := make(map[MeshDescF3V]uint32, numVerts)
	me.gpuSynced = false
	me.raw.verts = make([]float32, Core.MeshBuffers.FloatsPerVertex()*numVerts)
	me.raw.indices = make([]uint32, numVerts)
	me.raw.lastNumIndices = gl.Sizei(numVerts)
	me.raw.faces = make([]meshRawFace, len(meshData.Faces))
	for fi := 0; fi < len(meshData.Faces); fi++ {
		me.raw.faces[offsetFace].base = meshData.Faces[fi].MeshFaceBase
		tvp1 = meshData.Positions[meshData.Faces[fi].V[0].PosIndex]
		tvp2 = meshData.Positions[meshData.Faces[fi].V[1].PosIndex]
		tvp3 = meshData.Positions[meshData.Faces[fi].V[2].PosIndex]
		me.raw.faces[offsetFace].center.X = float64(tvp1[0]+tvp2[0]+tvp3[0]) / 3
		me.raw.faces[offsetFace].center.Y = float64(tvp1[1]+tvp2[1]+tvp3[1]) / 3
		me.raw.faces[offsetFace].center.Z = float64(tvp1[2]+tvp2[2]+tvp3[2]) / 3
		for ei, ventry = range meshData.Faces[fi].V {
			if vindex, vexists = vertsMap[ventry]; !vexists {
				vindex, vertsMap[ventry] = offsetVertex, offsetVertex
				copy(me.raw.verts[offsetFloat:(offsetFloat+3)], meshData.Positions[ventry.PosIndex][0:3])
				offsetFloat += 3
				copy(me.raw.verts[offsetFloat:(offsetFloat+2)], meshData.TexCoords[ventry.TexCoordIndex][0:2])
				offsetFloat += 2
				copy(me.raw.verts[offsetFloat:(offsetFloat+3)], meshData.Normals[ventry.NormalIndex][0:3])
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
	Diag.LogMeshes("mesh{%v}.Load() gave %v faces, %v att floats for %v final verts (%v source verts), %v indices (%vx vertex reuse)", me.Name, len(me.raw.faces), len(me.raw.verts), numFinalVerts, numVerts, len(me.raw.indices), vreuse)
	return
}

func (me *Mesh) Loaded() bool {
	return len(me.raw.verts) > 0
}

func (me *Mesh) Unload() {
	me.raw.verts, me.raw.indices = nil, nil
}

func (_ MeshLib) GpuSync() (err error) {
	for id := 0; id < len(Core.Libs.Meshes); id++ {
		if Core.Libs.Meshes.Ok(id) && !Core.Libs.Meshes[id].gpuSynced {
			if err = Core.Libs.Meshes[id].GpuUpload(); err != nil {
				return
			}
		}
	}
	return
}

func (me *MeshLib) AddNewAndLoad(name string, meshProvider MeshProvider) (meshID int, err error) {
	meshID = me.AddNew()
	mesh := &(*me)[meshID]
	mesh.Name = name
	if err = mesh.Load(meshProvider); err != nil {
		me.Remove(meshID, 1)
		meshID = -1
	}
	return
}

//#begin-gt -gen-lib.gt T:Mesh L:Core.Libs.Meshes

//	Only used for Core.Libs.Meshes
type MeshLib []Mesh

func (me *MeshLib) AddNew() (id int) {
	id = -1
	for i := 0; i < len(*me); i++ {
		if (*me)[i].ID == -1 {
			id = i
			break
		}
	}
	if id == -1 {
		if id = len(*me); id == cap(*me) {
			nu := make(MeshLib, id, id+Options.Libs.GrowCapBy)
			copy(nu, *me)
			*me = nu
		}
		*me = append(*me, Mesh{})
	}
	ref := &(*me)[id]
	ref.ID = id
	ref.init()
	return
}

func (me *MeshLib) Compact() {
	var (
		before, after []Mesh
		ref           *Mesh
		oldID, i      int
		compact       bool
	)
	for i = 0; i < len(*me); i++ {
		if (*me)[i].ID == -1 {
			compact, before, after = true, (*me)[:i], (*me)[i+1:]
			*me = append(before, after...)
		}
	}
	if compact {
		changed := make(map[int]int, len(*me))
		for i = 0; i < len(*me); i++ {
			if ref = &(*me)[i]; ref.ID != i {
				oldID, ref.ID = ref.ID, i
				changed[oldID] = i
			}
		}
		if len(changed) > 0 {
			me.onMeshIDsChanged(changed)
		}
	}
}

func (me *MeshLib) init() {
	*me = make(MeshLib, 0, Options.Libs.InitialCap)
}

func (me *MeshLib) dispose() {
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (me MeshLib) get(id int) (ref *Mesh) {
	if me.IsOk(id) {
		ref = &me[id]
	}
	return
}

func (me MeshLib) IsOk(id int) (ok bool) {
	if id > -1 && id < len(me) {
		ok = me[id].ID == id
	}
	return
}

func (me MeshLib) Ok(id int) bool {
	return me[id].ID == id
}

func (me MeshLib) Remove(fromID, num int) {
	if l := len(me); fromID > -1 && fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			me[id].dispose()
			changed[id], me[id].ID = -1, -1
		}
		me.onMeshIDsChanged(changed)
	}
}

func (me MeshLib) Walk(on func(ref *Mesh)) {
	for id := 0; id < len(me); id++ {
		if me.Ok(id) {
			on(&me[id])
		}
	}
}

//#end-gt
