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
		nodePreBatch nodeBatchPreps
	}
	thrRend struct {
		curCam                *Camera
		curEffect, nextEffect *FxEffect
		curMat                *FxMaterial
		curMeshBuf            *MeshBuffer
		curNode               *Node
		curProg               *ugl.Program
		curSampler            map[gl.Uint]*ugl.Sampler
		curTech, nextTech     RenderTechnique
		quadTex               *ugl.Texture2D
	}
)

func init() {
	thrRend.curSampler = map[gl.Uint]*ugl.Sampler{}
}

func (me *EngineCore) copyAppToPrep() {
	for _, canv := range me.Rendering.Canvases {
		if canv.renderThisFrame() {
			canv.copyAppToPrep()
		}
	}
}

func (me *EngineCore) copyPrepToRend() {
	for _, canv := range me.Rendering.Canvases {
		if canv.renderThisFrame() {
			canv.copyPrepToRend()
		}
	}
}

func (me *RenderCanvas) copyAppToPrep() {
	for _, cam := range me.Cameras {
		if cam.Enabled {
			cam.copyAppToPrep()
		}
	}
}

func (me *RenderCanvas) copyPrepToRend() {
	for _, cam := range me.Cameras {
		if cam.Enabled {
			cam.copyPrepToRend()
		}
	}
}

func (me *Camera) copyAppToPrep() {
	me.thrPrep.matProj = me.thrApp.matProj
	me.Controller.thrPrep.mat = me.Controller.thrApp.mat
	me.thrPrep.matPos.Translation(&me.Controller.Pos)
	if me.scene != nil {
		me.scene.RootNode.copyAppToPrep()
	}
}

func (me *Camera) copyPrepToRend() {
	me.thrRend.states = me.Rendering.States
	if me.scene != nil {
		me.scene.RootNode.copyPrepToRend()
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
