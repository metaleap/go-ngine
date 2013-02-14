package core

import (
	"sync"

	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
	unum "github.com/metaleap/go-util/num"
)

var (
	thrApp struct {
		sync.Mutex
		ctlTmps struct {
			tmpCopy           Controller
			matTrans, matLook unum.Mat4
			posNeg, axis      unum.Vec3
		}
		numBag    unum.Bag
		nodeTrans struct {
			matParent, matTrans, matScale, matRotX, matRotY, matRotZ unum.Mat4
		}
	}
	thrPrep struct {
		sync.Mutex
		numBag   unum.Bag
		curCam   *Camera
		curCanv  *RenderCanvas
		curScene *Scene
		curNode  *Node
	}
	thrRend struct {
		asyncResources       map[asyncResource]bool
		curCanvIndex         int
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
		numBag               unum.Bag
		tmpSampler           *ugl.Sampler
		tmpQuadTex           *ugl.Texture2D
		tmpTexUnit           gl.Uint
	}
)

func init() {
	thrRend.asyncResources = map[asyncResource]bool{}
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
	me.thrPrep.matCtlPos.Translation(&me.Controller.Pos)
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
		for me.thrPrep.tmpCam, me.thrPrep.tmpMat = range me.thrPrep.matProjs {
			me.thrRend.matProjs[me.thrPrep.tmpCam].Load(me.thrPrep.tmpMat)
		}
		for _, me.thrPrep.curSubNode = range me.ChildNodes.M {
			me.thrPrep.curSubNode.copyPrepToRend()
		}
		me.thrPrep.done = false
	}
}
