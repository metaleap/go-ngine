package core

import (
	"sort"

	"github.com/go-utils/uhash"
	"github.com/go-utils/unum"
	"github.com/go-utils/uslice"
)

const numPrios = 3

type RenderBatchCriteria int

const (
	BatchByProgram RenderBatchCriteria = 0
	BatchByTexture RenderBatchCriteria = 1
	BatchByBuffer  RenderBatchCriteria = 2
)

type renderBatchEntry struct {
	node, mesh, prog, fx int
	dist                 float64
	texes                []int
	face                 int32
}

type renderBatchList struct {
	all       []renderBatchEntry
	distStats renderBatchDistStats
	ki, kj    renderBatchDistStatsKey
	n         int
	prios     [numPrios]RenderBatchCriteria
}

func (me *renderBatchList) calcStats() {
	me.distStats = make(renderBatchDistStats, numPrios*len(me.all))
	var key renderBatchDistStatsKey
	for i := 0; i < len(me.all); i++ {
		for max := 0; max < numPrios; max++ {
			me.statsKey(i, max, &key)
			me.distStats.add(me.all[i].dist, key)
		}
	}
}

func (me *renderBatchList) statsKey(index, max int, key *renderBatchDistStatsKey) {
	for i := 0; i < len(me.prios); i++ {
		if i > max {
			key[i] = -1
		} else {
			switch me.prios[i] {
			case BatchByProgram:
				key[i] = me.all[index].prog
			case BatchByBuffer:
				key[i] = int(Core.Libs.Meshes[me.all[index].mesh].meshBuffer.glIbo.GlHandle)
			case BatchByTexture:
				key[i] = uhash.Fnv1a(me.all[index].texes)
			}
		}
	}
}

func (me *renderBatchList) compare(i, j int, crit RenderBatchCriteria) (less, equal bool) {
	switch crit {
	case BatchByProgram:
		less = me.all[i].prog < me.all[j].prog
		equal = me.all[i].prog == me.all[j].prog
	case BatchByBuffer:
		less = Core.Libs.Meshes[me.all[i].mesh].meshBuffer.glIbo.GlHandle < Core.Libs.Meshes[me.all[j].mesh].meshBuffer.glIbo.GlHandle
		equal = Core.Libs.Meshes[me.all[i].mesh].meshBuffer.glIbo.GlHandle == Core.Libs.Meshes[me.all[j].mesh].meshBuffer.glIbo.GlHandle
	case BatchByTexture:
		less, equal = false, true
		for t := 0; (less || equal) && t < len(me.all[i].texes); t++ {
			less = less || me.all[i].texes[t] < me.all[j].texes[t]
			equal = equal && me.all[i].texes[t] == me.all[j].texes[t]
		}
	}
	return
}

func (me *renderBatchList) Len() int {
	return me.n
}

func (me *renderBatchList) Less(i, j int) (less bool) {
	var mi, mj float64
	var lt, eq bool
	for crit := 0; crit < numPrios; crit++ {
		if lt, eq = me.compare(i, j, me.prios[crit]); !eq {
			me.statsKey(i, crit, &me.ki)
			me.statsKey(j, crit, &me.kj)
			mi, mj = me.distStats[me.ki].mean, me.distStats[me.kj].mean
			if mi == mj {
				return lt
			} else {
				return mi < mj
			}
		}
	}
	return me.all[i].dist < me.all[j].dist
}

// func (me *renderBatchList) oldLess(i, j int) (less bool) {
// 	var eq bool
// 	if less, eq = me.compare(i, j, me.prios[0]); eq {
// 		if less, eq = me.compare(i, j, me.prios[1]); eq {
// 			if less, eq = me.compare(i, j, me.prios[2]); eq {
// 				less = me.all[i].dist < me.all[j].dist
// 			}
// 		}
// 	}
// 	return
// }

func (me *renderBatchList) Swap(i, j int) {
	me.all[i], me.all[j] = me.all[j], me.all[i]
}

type renderBatchDistStatsKey [numPrios]int

type renderBatchDistStats map[renderBatchDistStatsKey]struct {
	count, mean, sum float64
}

func (me renderBatchDistStats) add(dist float64, key renderBatchDistStatsKey) {
	v := me[key]
	v.count++
	v.sum += dist
	v.mean = v.sum / v.count
	me[key] = v
}

type RenderBatcher struct {
	Enabled  bool
	Priority [numPrios]RenderBatchCriteria
}

func (me *RenderTechniqueScene) prepEntry(n, nid, fx int, fi int32) {
	entry := &me.thrPrep.batch.all[n]
	entry.node, entry.fx, entry.face, entry.mesh = nid, fx, fi, Core.Libs.Scenes[me.Camera.sceneID].allNodes[nid].Render.meshID
	var distPos *unum.Vec3
	if fi == -1 {
		distPos = &Core.Libs.Scenes[me.Camera.sceneID].allNodes[nid].Transform.Pos
	} else {
		distPos = Core.Libs.Scenes[me.Camera.sceneID].allNodes[nid].Transform.Pos.AddTo(&Core.Libs.Meshes[entry.mesh].raw.faces[fi].center)
	}
	entry.dist = me.Camera.Controller.Pos.ManhattanDistanceFrom(distPos)
	entry.prog = ogl.progs.Index(Core.Libs.Effects[fx].uberPnames[me.name()])
	uslice.IntEnsureLen(&entry.texes, Stats.Programs.maxTexUnits)
	var ti int
	for ti = 0; ti < len(entry.texes); ti++ {
		entry.texes[ti] = -1
	}
	ti = 0
	for oi := 0; oi < len(Core.Libs.Effects[fx].FxProcs); oi++ {
		if Core.Libs.Effects[fx].FxProcs[oi].IsTex() {
			entry.texes[ti] = Core.Libs.Effects[fx].FxProcs[oi].Tex.ImageID
			ti++
		}
	}
	me.thrPrep.Done()
}

func (me *RenderTechniqueScene) prepBatch(scene *Scene, size int) {
	var (
		mesh   *Mesh
		mat    *FxMaterial
		effect *FxEffect
		fi, fl int32
		nid    int
	)
	b := &me.thrPrep.batch
	b.n = 0
	if len(b.all) < size {
		b.all = make([]renderBatchEntry, size)
	}
	for nid = 1; nid < len(scene.allNodes); nid++ {
		if scene.allNodes.Ok(nid) && me.Camera.thrPrep.nodeRender[nid] {
			if mesh, mat = scene.allNodes[nid].meshMat(); mat.HasFaceEffects() {
				for fi, fl = 0, int32(len(mesh.raw.faces)); fi < fl; fi++ {
					if effect = mat.faceEffect(&mesh.raw.faces[fi]); effect != nil {
						me.thrPrep.Add(1)
						go me.prepEntry(b.n, nid, effect.ID, fi)
						b.n++
					}
				}
			} else if effect = Core.Libs.Effects.get(mat.DefaultEffectID); effect != nil {
				me.thrPrep.Add(1)
				go me.prepEntry(b.n, nid, effect.ID, -1)
				b.n++
			}
		}
	}
	b.prios = me.Batch.Priority
	me.thrPrep.Wait()
	b.calcStats()
	sort.Sort(b)
}
