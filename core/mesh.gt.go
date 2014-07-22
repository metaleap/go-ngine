package core

import (
	"github.com/go-utils/unum"
	u3d "github.com/go3d/go-3dutil"
	gl "github.com/go3d/go-opengl/core"
)

//	Represents semi-processed loaded mesh data "almost ready" to core.Mesh.GpuUpload().
type meshRaw struct {
	lastNumIndices gl.Sizei

	//	Raw vertices
	verts []float32

	//	Vertex indices
	indices []uint32

	//	Raw face definitions
	faces []meshRawFace

	bounding u3d.Bounds
}

//	Represents a triangle face inside a meshRaw.
type meshRawFace struct {
	//	Indices of the triangle corners
	entries [3]uint32
	center  unum.Vec3

	base u3d.MeshFaceBase
}

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

func (me *Mesh) Load(provider u3d.MeshProvider) (err error) {
	var meshData *u3d.MeshDescriptor
	if meshData, err = provider(); err == nil && meshData != nil {
		err = me.load(meshData)
	}
	return
}

func (me *Mesh) load(meshData *u3d.MeshDescriptor) (err error) {
	var (
		offsetFloat, offsetIndex, offsetVertex, vindex uint32
		f                                              float64
		vreuse, offsetFace, ei, numFinalVerts, tvi     int
		vexists                                        bool
		ventry                                         u3d.MeshDescF3V
		tvp                                            [3]unum.Vec3
	)
	numVerts := 3 * int32(len(meshData.Faces))
	vertsMap := make(map[u3d.MeshDescF3V]uint32, numVerts)
	me.gpuSynced = false
	me.raw.bounding.Reset()
	me.raw.verts = make([]float32, Core.Mesh.Buffers.FloatsPerVertex()*numVerts)
	me.raw.indices = make([]uint32, numVerts)
	me.raw.lastNumIndices = gl.Sizei(numVerts)
	me.raw.faces = make([]meshRawFace, len(meshData.Faces))
	for fi := 0; fi < len(meshData.Faces); fi++ {
		me.raw.faces[offsetFace].base = meshData.Faces[fi].MeshFaceBase
		meshData.Positions[meshData.Faces[fi].V[0].PosIndex].ToVec3(&tvp[0])
		meshData.Positions[meshData.Faces[fi].V[1].PosIndex].ToVec3(&tvp[1])
		meshData.Positions[meshData.Faces[fi].V[2].PosIndex].ToVec3(&tvp[2])
		for tvi = 0; tvi < len(tvp); tvi++ {
			if f = tvp[tvi].DistanceFromZero(); f > me.raw.bounding.Sphere {
				me.raw.bounding.Sphere = f
			}
			me.raw.bounding.AaBox.UpdateMinMax(&tvp[tvi])
		}
		me.raw.faces[offsetFace].center.X = (tvp[0].X + tvp[1].X + tvp[2].X) / 3
		me.raw.faces[offsetFace].center.Y = (tvp[0].Y + tvp[1].Y + tvp[2].Y) / 3
		me.raw.faces[offsetFace].center.Z = (tvp[0].Z + tvp[1].Z + tvp[2].Z) / 3
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
	me.raw.bounding.AaBox.SetCenterExtent()
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

func (me *MeshLib) AddNewAndLoad(name string, meshProvider u3d.MeshProvider) (meshID int, err error) {
	meshID = me.AddNew()
	mesh := &(*me)[meshID]
	mesh.Name = name
	if err = mesh.Load(meshProvider); err != nil {
		me.Remove(meshID, 1)
		meshID = -1
	}
	return
}

func (_ MeshLib) MeshCube() u3d.MeshProvider {
	return u3d.MeshDescriptorCube
}

func (_ MeshLib) MeshPlane() u3d.MeshProvider {
	return u3d.MeshDescriptorPlane
}

func (_ MeshLib) MeshPyramid() u3d.MeshProvider {
	return u3d.MeshDescriptorPyramid
}

func (_ MeshLib) MeshQuad() u3d.MeshProvider {
	return u3d.MeshDescriptorQuad
}

func (_ MeshLib) MeshTri() u3d.MeshProvider {
	return u3d.MeshDescriptorTri
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
