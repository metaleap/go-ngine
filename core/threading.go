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
		// nodePreBatch nodeBatchPreps
	}
	thrRend struct {
		curCam                *Camera
		curEffect, nextEffect *FxEffect
		curTech, nextTech     RenderTechnique
		curProg               *ugl.Program
		quadTex               gl.Uint
	}
)

func init() {
}

func (_ *EngineCore) copyAppToPrep() {
	for cid := 0; cid < len(Core.Render.Canvases); cid++ {
		if Core.Render.Canvases[cid].renderThisFrame() {
			Core.Render.Canvases[cid].copyAppToPrep()
		}
	}
}

func (_ *EngineCore) copyPrepToRend() {
	for cid := 0; cid < len(Core.Render.Canvases); cid++ {
		if Core.Render.Canvases[cid].renderThisFrame() {
			Core.Render.Canvases[cid].copyPrepToRend()
		}
	}
}

func (me *RenderCanvas) copyAppToPrep() {
	for cam := 0; cam < len(me.Cams); cam++ {
		me.Cams[cam].copyAppToPrep()
	}
}

func (me *RenderCanvas) copyPrepToRend() {
	for cam := 0; cam < len(me.Cams); cam++ {
		me.Cams[cam].copyPrepToRend()
	}
}

func (me *Camera) copyAppToPrep() {
	me.thrPrep.matProj = me.thrApp.matProj
	me.Controller.thrPrep.mat = me.Controller.thrApp.mat
	me.thrPrep.matPos.Translation(&me.Controller.Pos)
	if scene := me.scene(); scene != nil {
		scene.RootNode.copyAppToPrep()
	}
}

func (me *Camera) copyPrepToRend() {
	me.thrRend.states = me.Rendering.States
	if scene := me.scene(); scene != nil {
		scene.RootNode.copyPrepToRend()
	}
}

func (me *Node) copyAppToPrep() {
	if !me.thrPrep.copyDone {
		me.thrPrep.copyDone = true
		me.thrPrep.matModelView = me.Transform.matModelView
		for _, subNode := range me.ChildNodes.M {
			subNode.copyAppToPrep()
		}
	}
}

func (me *Node) copyPrepToRend() {
	if !me.thrRend.copyDone {
		me.thrRend.copyDone = true
		for cam, mat := range me.thrPrep.camProjMats {
			me.thrRend.camProjMats[cam].Load(mat)
		}
		for _, subNode := range me.ChildNodes.M {
			subNode.copyPrepToRend()
		}
		me.thrPrep.done = false
	}
}
