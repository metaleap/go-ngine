package core

import (
	"fmt"

	gl "github.com/chsc/gogl/gl42"

	glutil "github.com/go3d/go-util/gl"
)

type tMeshes map[string]*TMesh

	func (me tMeshes) Add (mesh *TMesh) *TMesh {
		if me[mesh.name] == nil {
			me[mesh.name] = mesh
			return mesh
		}
		return nil
	}

	func (me tMeshes) AddRange (meshes ... *TMesh) {
		for _, m := range meshes { me.Add(m) }
	}

	func (me tMeshes) Load (name string, provider FMeshProvider, args ... interface {}) (mesh *TMesh, err error) {
		var meshData *tMeshData
		mesh = me.New(name)
		if meshData, err = provider(args ...); err == nil {
			mesh.load(meshData)
		} else {
			mesh = nil
		}
		return
	}

	func (me tMeshes) New (name string) *TMesh {
		var mesh = &TMesh {}
		mesh.name = name
		return mesh
	}

type tMeshRaw struct {
	verts []gl.Float
	indices []gl.Uint
	faces []*TMeshFace
}

type TMeshFace struct {
	entries [3]gl.Uint
}

	func newMeshFace () *TMeshFace {
		var face = &TMeshFace {}
		return face
	}

type TMesh struct {
	name string
	meshBuffer *TMeshBuffer
	meshBufOffsetBaseIndex, meshBufOffsetIndices, meshBufOffsetVerts int32
	raw *tMeshRaw
	gpuSynced bool
}

	func (me *TMesh) GpuDelete () {
		if me.gpuSynced {
			me.gpuSynced = false
		}
	}

	func (me *TMesh) GpuUpload () (err error) {
		var sizeVerts, sizeIndices = gl.Sizeiptr(4 * len(me.raw.verts)), gl.Sizeiptr(4 * len(me.raw.indices))
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
			gl.BufferSubData(gl.ARRAY_BUFFER, gl.Intptr(me.meshBufOffsetVerts), sizeVerts, gl.Pointer(&me.raw.verts[0]))
			me.meshBuffer.offsetVerts += int32(sizeVerts)
			gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, gl.Intptr(me.meshBufOffsetIndices), sizeIndices, gl.Pointer(&me.raw.indices[0]))
			me.meshBuffer.offsetIndices += int32(sizeIndices)
			me.meshBuffer.offsetBaseIndex += int32(len(me.raw.indices))
			gl.BindBuffer(gl.ARRAY_BUFFER, 0)
			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
			if err = glutil.LastError("mesh[%v].GpuUpload()", me.name); err == nil { me.gpuSynced = true }
		}
		return 
	}

	func (me *TMesh) GpuUploaded () bool {
		return me.gpuSynced
	}

	func (me *TMesh) load (meshData *tMeshData) {
		var numVerts = 3 * int32(len(meshData.faces))
		var vertsMap = map[tVe]gl.Uint {}
		var offsetFloat, offsetIndex, offsetVertex, vindex gl.Uint
		var offsetFace, ei = 0, 0
		var face tMeshFace3
		var ventry tVe
		var vexists bool
		me.gpuSynced = false
		me.raw = &tMeshRaw {}
		me.raw.verts = make([]gl.Float, Core.MeshBuffers.FloatsPerVertex() * numVerts)
		me.raw.indices = make([]gl.Uint, numVerts)
		me.raw.faces = make([]*TMeshFace, len(meshData.faces))
		for _, face = range meshData.faces {
			me.raw.faces[offsetFace] = newMeshFace()
			for ei, ventry = range face {
				if vindex, vexists = vertsMap[ventry]; !vexists {
					vindex, vertsMap[ventry] = offsetVertex, offsetVertex
					copy(me.raw.verts[offsetFloat : (offsetFloat + 3)], meshData.positions[ventry.posIndex][0 : 3])
					offsetFloat += 3
					copy(me.raw.verts[offsetFloat : (offsetFloat + 2)], meshData.texCoords[ventry.texCoordIndex][0 : 2])
					offsetFloat += 2
					copy(me.raw.verts[offsetFloat : (offsetFloat + 3)], meshData.normals[ventry.normalIndex][0 : 3])
					offsetFloat += 3
					offsetVertex++
				}
				me.raw.indices[offsetIndex] = vindex
				me.raw.faces[offsetFace].entries[ei] = offsetIndex
				offsetIndex++
			}
			offsetFace++
		}
		fmt.Printf("meshload() gave %v faces, %v att floats for %v verts, %v indices\n", len(me.raw.faces), len(me.raw.verts), numVerts, len(me.raw.indices))
	}

	func (me *TMesh) Loaded () bool {
		return me.raw != nil
	}

	func (me *TMesh) render () {
		curMesh = me
		if curMeshBuf != me.meshBuffer { me.meshBuffer.bindVao() }
		curTechnique.onRenderMesh()
		gl.DrawElementsBaseVertex(gl.TRIANGLES, gl.Sizei(len(me.raw.indices)), gl.UNSIGNED_INT, gl.Offset(nil, uintptr(me.meshBufOffsetIndices)), gl.Int(me.meshBufOffsetBaseIndex))
//		gl.DrawElements(gl.TRIANGLES, gl.Sizei(len(me.raw.indices)), gl.UNSIGNED_INT, gl.Pointer(nil))
	}

	func (me *TMesh) Unload () {
		me.raw = nil
	}
