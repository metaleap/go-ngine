package core

import (
	gl "github.com/go3d/go-opengl/core"
	ugl "github.com/go3d/go-opengl/util"
	unum "github.com/metaleap/go-util/num"
)

type CameraPerspective struct {
	//	Vertical field-of-view. Defaults to 37.8493.
	//	After changing this value, you must call the ApplyMatrices() method.
	FovY float64

	//	Distance of the far-plane from the camera. Defaults to 30000.
	//	After changing this value, you must call the ApplyMatrices() method.
	ZFar float64

	//	Distance of the near-plane from the camera. Defaults to 0.3.
	//	After changing this value, you must call the ApplyMatrices() method.
	ZNear float64
}

//	A camera embodies the eye point of the viewer looking at the visual scene.
type Camera struct {
	//	Optical and imager properties for this camera.
	Perspective struct {
		//	FovY, ZFar, ZNear
		CameraPerspective

		//	Whether this is a perspective-projection camera. Defaults to true.
		//	If false, no projection transformation is applied.
		Use bool
	}

	//	Encapsulates the position and direction of this camera.
	Controller Controller

	Enabled bool

	Rendering struct {
		//	The device-relative or absolute view-port for this Camera.
		Viewport CameraViewport

		States ugl.RenderStatesBag

		Technique RenderTechnique

		FxOps FxOps
	}

	sceneID int

	thrApp struct {
		matProj unum.Mat4
	}
	thrPrep struct {
		onPrepNode                  func(*Node)
		matCamProj, matProj, matPos unum.Mat4
	}
	thrRend struct {
		states ugl.RenderStatesBag
	}
}

func (me *Camera) setup2D(canv *RenderCanvas, depth bool) {
	me.setup(canv, false, depth, Options.Rendering.DefaultTechnique2D)
	return
}

func (me *Camera) setup3D(canv *RenderCanvas) {
	me.setup(canv, true, true, Options.Rendering.DefaultTechnique3D)
	return
}

func (me *Camera) setupQuad(canv *RenderCanvas) {
	me.setup(canv, false, false, Options.Rendering.DefaultTechniqueQuad)
	return
}

func (me *Camera) setup(canv *RenderCanvas, persp3d bool, depth bool, technique string) {
	me.Rendering.Viewport.canvWidth, me.Rendering.Viewport.canvHeight = float64(canv.absViewWidth), float64(canv.absViewHeight)
	me.Perspective.Use = persp3d
	me.Rendering.States.DepthTest = depth
	if depth {
		me.Rendering.States.Other.ClearBits = gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT
	} else {
		me.Rendering.States.Other.ClearBits = gl.COLOR_BUFFER_BIT
	}
	me.Controller.init()
	me.Rendering.Viewport.init()
	me.ApplyMatrices()
	me.Rendering.Technique = Core.Render.KnownTechniques[technique]()
}

func (me *Camera) init() {
	me.Enabled = true
	me.sceneID = -1
	me.thrPrep.onPrepNode = func(n *Node) { me.onPrepNode(n) }
	me.Rendering.States.FaceCulling, me.Rendering.States.StencilTest = false, false
	me.Rendering.States.ClearColor = Options.Rendering.DefaultClearColor
	me.Perspective.CameraPerspective = Options.Cameras.PerspectiveDefaults
	unum.Mat4Identities(&me.thrApp.matProj, &me.thrPrep.matProj)
}

//	Applies changes made to the FovY, ZNear and/or ZFar parameters in me.Perspective.
func (me *Camera) ApplyMatrices() {
	me.thrApp.matProj.Perspective(me.Perspective.FovY, me.Rendering.Viewport.aspect, me.Perspective.ZNear, me.Perspective.ZFar)
}

func (me *Camera) clearNodeMats() {
	if scene := me.scene(); scene != nil {
		scene.RootNode.Walk(func(node *Node) {
			delete(node.thrPrep.camProjMats, me)
			delete(node.thrRend.camProjMats, me)
			delete(node.thrPrep.camRender, me)
		})
	}
}

func (me *Camera) dispose() {
	me.clearNodeMats()
	me.Rendering.Technique.dispose()
}

func (me *Camera) scene() *Scene {
	return Core.Libs.Scenes.get(me.sceneID)
}

func (me *Camera) SetScene(sceneID int) {
	if sceneID != me.sceneID {
		me.clearNodeMats()
		me.sceneID = sceneID
		if scene := me.scene(); scene != nil {
			scene.RootNode.Walk(func(node *Node) {
				node.initCamData(me)
			})
		}
	}
}

func (me *Camera) RenderTechniqueQuad() (tech *RenderTechniqueQuad) {
	tech, _ = me.Rendering.Technique.(*RenderTechniqueQuad)
	return
}

func (me *Camera) RenderTechniqueScene() (tech *RenderTechniqueScene) {
	tech, _ = me.Rendering.Technique.(*RenderTechniqueScene)
	return
}

//	Encapsulates a device-relative or absolute camera view-port.
type CameraViewport struct {
	relative, shouldScissor bool
	relX, relY, relW, relH  float64
	absX, absY, absW, absH  int
	aspect                  float64
	glVpX, glVpY            gl.Int
	glVpW, glVpH            gl.Sizei
	canvWidth, canvHeight   float64
}

func (me *CameraViewport) init() {
	me.SetRel(0, 0, 1, 1)
}

//	Sets the absolute viewport origin and dimensions in pixels.
func (me *CameraViewport) SetAbs(x, y, width, height int) {
	me.relative, me.absX, me.absY, me.absW, me.absH = false, x, y, width, height
	me.update()
}

//	Sets the device-relative viewport origin and dimensions, with the value 1.0
//	representing the maximum extent of the viewport on that respective axis.
func (me *CameraViewport) SetRel(x, y, width, height float64) {
	me.relative, me.relX, me.relY, me.relW, me.relH = true, x, y, width, height
	me.update()
}

func (me *CameraViewport) update() {
	me.shouldScissor = !(me.relative && me.relX == 0 && me.relY == 0 && me.relW == 1 && me.relH == 1)
	if me.relative {
		me.absW, me.absH = int(me.relW*float64(me.canvWidth)), int(me.relH*float64(me.canvHeight))
		me.absX, me.absY = int(me.relX*float64(me.canvWidth)), int(me.relY*float64(me.canvHeight))
	}
	me.glVpX, me.glVpY, me.glVpW, me.glVpH = gl.Int(me.absX), gl.Int(me.absY), gl.Sizei(me.absW), gl.Sizei(me.absH)
	me.aspect = float64(me.absW) / float64(me.absH)
}

//#begin-gt -gen-reflib.gt T:Camera L:Core.Render.Canvases[id].Cams

//	Only used for Core.Render.Canvases[id].Cams
type CameraLib []*Camera

func (me *CameraLib) add(ref *Camera) {
	*me = append(*me, ref)
	ref.init()
	return
}

func (me *CameraLib) init() {
	*me = make(CameraLib, 0, 4)
}

func (me *CameraLib) dispose() {
	me.Remove(0, 0)
}

func (me CameraLib) IsOk(id int) bool {
	return id > -1 && id < len(me)
}

func (me *CameraLib) Remove(fromID, num int) {
	if l := len(*me); fromID > -1 && fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		for i := fromID; i < fromID+num; i++ {
			(*me)[fromID].dispose()
		}
		before, after := (*me)[:fromID], (*me)[fromID+num:]
		*me = append(before, after...)
	}
}

func (me CameraLib) Walk(on func(ref *Camera)) {
	for id := 0; id < len(me); id++ {
		on(me[id])
	}
}

//#end-gt
