package core

import (
	"sync"

	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
)

var (
	thrApp struct {
		sync.Mutex
	}
	thrPrep struct {
		sync.Mutex
	}
	thrRend struct {
		curCam                *Camera
		curView               *RenderView
		curEffect, nextEffect *FxEffect
		curTech, nextTech     RenderTechnique
		curProg               *ugl.Program
		quadTex               gl.Uint
	}
)

func init() {
}

func (_ *NgCore) copyAppToPrep() {
	for cid := 0; cid < len(Core.Render.Canvases); cid++ {
		if Core.Render.Canvases[cid].renderThisFrame() {
			Core.Render.Canvases[cid].copyAppToPrep()
		}
	}
}

func (_ *NgCore) copyPrepToRend() {
	for cid := 0; cid < len(Core.Render.Canvases); cid++ {
		if Core.Render.Canvases[cid].renderThisFrame() {
			Core.Render.Canvases[cid].copyPrepToRend()
		}
	}
}

func (me *RenderCanvas) copyAppToPrep() {
	for view := 0; view < len(me.Views); view++ {
		me.Views[view].copyAppToPrep()
	}
}

func (me *RenderCanvas) copyPrepToRend() {
	for view := 0; view < len(me.Views); view++ {
		me.Views[view].copyPrepToRend()
	}
}

func (me *RenderView) copyAppToPrep() {
	me.Technique.copyAppToPrep()
}

func (me *RenderView) copyPrepToRend() {
	me.Technique.copyPrepToRend()
}

func (me *RenderTechniqueScene) copyAppToPrep() {
	me.Camera.copyAppToPrep()
}

func (me *RenderTechniqueScene) copyPrepToRend() {
	me.Camera.copyPrepToRend()
	if me.Batch.Enabled {
		n := me.thrPrep.batch.n
		me.thrRend.batch.n = n
		if len(me.thrRend.batch.all) < n {
			me.thrRend.batch.all = make([]renderBatchEntry, n)
		}
		copy(me.thrRend.batch.all, me.thrPrep.batch.all)
	}
}

func (me *Camera) copyAppToPrep() {
	me.thrPrep.matProj = me.thrApp.matProj
	me.Controller.thrPrep.mat = me.Controller.thrApp.mat
	me.thrPrep.matPos.Translation(&me.Controller.Pos)
	if scene := me.Scene(); scene != nil {
		scene.copyAppToPrep()
	}
}

func (me *Camera) copyPrepToRend() {
	copy(me.thrRend.nodeRender, me.thrPrep.nodeRender)
	for i := 0; i < len(me.thrPrep.nodeProjMats); i++ {
		me.thrRend.nodeProjMats[i].Load(&me.thrPrep.nodeProjMats[i])
	}
	if scene := me.Scene(); scene != nil {
		scene.copyPrepToRend()
	}
}

func (me *Scene) copyAppToPrep() {
	if !me.thrPrep.copyDone {
		me.thrPrep.copyDone = true
		for i := 0; i < len(me.allNodes); i++ {
			if me.allNodes.Ok(i) {
				me.allNodes[i].copyAppToPrep()
			}
		}
	}
}

func (me *Scene) copyPrepToRend() {
	if !me.thrRend.copyDone {
		me.thrRend.copyDone = true
		for i := 0; i < len(me.allNodes); i++ {
			if me.allNodes.Ok(i) {
				me.allNodes[i].copyPrepToRend()
			}
		}
		me.thrPrep.done = false
	}
}

func (me *SceneNode) copyAppToPrep() {
	me.Transform.thrPrep.matModelView = me.Transform.thrApp.matModelView
}

func (me *SceneNode) copyPrepToRend() {
}
