package core

import (
	ugl "github.com/go3d/go-glutil"
)

var (
	thrPrep struct {
		curCam   *Camera
		curCanv  *RenderCanvas
		curScene *Scene
		curNode  *Node
	}
	thrRend struct {
		asyncResources   map[asyncResource]bool
		curCanv          *RenderCanvas
		curCam           *Camera
		curMatId         string
		curMat           *FxMaterial
		curMeshBuf       *MeshBuffer
		curNode          *Node
		curProg, tmpProg *ugl.Program
		curScene         *Scene
		curTechnique     renderTechnique
		tmpEffect        *FxEffect
		tmpMat           *FxMaterial
	}
)

func init() {
	thrRend.asyncResources = map[asyncResource]bool{}
}

func (me *EngineCore) copyAppToPrep() {
	for _, thrPrep.curCanv = range me.Rendering.Canvases {
		if thrPrep.curCanv.renderThisFrame() {
			thrPrep.curCanv.copyAppToPrep()
		}
	}
}

func (me *EngineCore) copyPrepToRend() {
	for _, thrPrep.curCanv = range me.Rendering.Canvases {
		if thrPrep.curCanv.renderThisFrame() {
			thrPrep.curCanv.copyPrepToRend()
		}
	}
}

func (me *RenderCanvas) copyAppToPrep() {
	for _, thrPrep.curCam = range me.Cameras {
		thrPrep.curCam.copyAppToPrep()
	}
}

func (me *RenderCanvas) copyPrepToRend() {
	for _, thrPrep.curCam = range me.Cameras {
		thrPrep.curCam.copyPrepToRend()
	}
}

func (me *Camera) copyAppToPrep() {
	me.thrPrep.matProj = me.thrApp.matProj
	me.Controller.thrPrep.mat = me.Controller.thrApp.mat
	if thrPrep.curScene = Core.Libs.Scenes[me.Rendering.SceneID]; thrPrep.curScene != nil {
		thrPrep.curScene.RootNode.copyAppToPrep()
	}
}

func (me *Camera) copyPrepToRend() {
	me.thrRend.states = me.Rendering.States
	if thrRend.curScene = Core.Libs.Scenes[me.Rendering.SceneID]; thrRend.curScene != nil {
		thrRend.curScene.RootNode.copyPrepToRend()
	}
}

func (me *Node) copyAppToPrep() {
	me.thrPrep.model = me.model
	me.thrPrep.matModelView = me.Transform.matModelView
	for _, me.thrPrep.curSubNode = range me.ChildNodes.M {
		me.thrPrep.curSubNode.copyAppToPrep()
	}
}

func (me *Node) copyPrepToRend() {
	me.thrRend.glMatModelProjs = me.thrPrep.glMatModelProjs
	for _, me.thrPrep.curSubNode = range me.ChildNodes.M {
		me.thrPrep.curSubNode.copyPrepToRend()
	}
}
