package core

import (
	"sort"
)

type RenderBatchCriteria int

const (
	BatchByProgram RenderBatchCriteria = iota + 1
	BatchByTexture
	BatchByBuffer
)

type renderBatchEntry struct {
	node, fx, mesh int
	face           int32
}

type renderBatchList struct {
	all []renderBatchEntry
	n   int
}

func (me *renderBatchList) Len() int {
	return me.n
}

func (me *renderBatchList) Less(i, j int) bool {
	// lessProg := me.all[i].fx < me.all[j].fx
	// less
	return j < i
}

func (me *renderBatchList) Swap(i, j int) {
	me.all[i], me.all[j] = me.all[j], me.all[i]
}

type RenderBatcher struct {
	Enabled  bool
	Priority [3]RenderBatchCriteria
}

func (me *RenderTechniqueScene) prepBatch(scene *Scene, size int) {
	var (
		entry  *renderBatchEntry
		mesh   *Mesh
		mat    *FxMaterial
		fi, fl int32
	)
	b := &me.thrPrep.batch
	b.n = 0
	if len(b.all) < size {
		b.all = make([]renderBatchEntry, size)
	}
	for i := 0; i < len(scene.allNodes); i++ {
		if scene.allNodes.Ok(i) && me.Camera.thrPrep.nodeRender[i] {
			if mesh, mat = scene.allNodes[i].meshMat(); mat.HasFaceEffects() {
				for fi, fl = 0, int32(len(mesh.raw.faces)); fi < fl; fi++ {
					entry = &b.all[b.n]
					entry.face = fi
					entry.mesh = mesh.ID
					entry.node = i
					entry.fx = mat.faceEffectID(&mesh.raw.faces[fi])
					b.n++
				}
			} else {
				entry = &b.all[b.n]
				entry.face = -1
				entry.mesh = mesh.ID
				entry.node = i
				entry.fx = mat.DefaultEffectID
				b.n++
			}
		}
	}
	sort.Sort(b)
}
