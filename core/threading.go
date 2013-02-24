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
		curCam       *Camera
		curCanv      *RenderCanvas
		curScene     *Scene
		curNode      *Node
		curTechScene *RenderTechniqueScene
		nodePreBatch nodeBatchPreps
	}
	thrRend struct {
		curCanv              *RenderCanvas
		curCam               *Camera
		curEffect, tmpEffect *FxEffect
		curMatId             string
		curMat, tmpMat       *FxMaterial
		curMeshBuf           *MeshBuffer
		curNode              *Node
		curProg, tmpProg     *ugl.Program
		curScene             *Scene
		curSampler           map[gl.Uint]*ugl.Sampler
		curTech, tmpTech     RenderTechnique
		tmpSampler           *ugl.Sampler
		tmpQuadTex           *ugl.Texture2D
		tmpFaceTag           string
		tmpFace              *meshRawFace
	}
)

func init() {
	thrRend.curSampler = map[gl.Uint]*ugl.Sampler{}
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
		if thrPrep.curCam.Enabled {
			thrPrep.curCam.copyAppToPrep()
		}
	}
}

func (me *RenderCanvas) copyPrepToRend() {
	for _, thrPrep.curCam = range me.Cameras {
		if thrPrep.curCam.Enabled {
			thrPrep.curCam.copyPrepToRend()
		}
	}
}

func (me *Camera) copyAppToPrep() {
	me.thrPrep.matProj = me.thrApp.matProj
	me.Controller.thrPrep.mat = me.Controller.thrApp.mat
	me.thrPrep.matPos.Translation(&me.Controller.Pos)
	if thrPrep.curScene = me.scene; thrPrep.curScene != nil {
		thrPrep.curScene.RootNode.copyAppToPrep()
	}
}

func (me *Camera) copyPrepToRend() {
	me.thrRend.states = me.Rendering.States
	if thrPrep.curScene = me.scene; thrPrep.curScene != nil {
		thrPrep.curScene.RootNode.copyPrepToRend()
	}
}

func (me *Node) copyAppToPrep() {
	if !me.thrPrep.copyDone {
		me.thrPrep.copyDone = true
		me.thrPrep.model = me.model
		me.thrPrep.matModelView = me.Transform.matModelView
		for _, me.thrPrep.curSubNode = range me.ChildNodes.M {
			me.thrPrep.curSubNode.copyAppToPrep()
		}
	}
}

func (me *Node) copyPrepToRend() {
	if !me.thrRend.copyDone {
		me.thrRend.copyDone = true
		for me.thrPrep.tmpCam, me.thrPrep.tmpMat = range me.thrPrep.camProjMats {
			me.thrRend.camProjMats[me.thrPrep.tmpCam].Load(me.thrPrep.tmpMat)
		}
		for _, me.thrPrep.curSubNode = range me.ChildNodes.M {
			me.thrPrep.curSubNode.copyPrepToRend()
		}
		me.thrPrep.done = false
	}
}
