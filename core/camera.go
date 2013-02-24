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

	scene *Scene

	thrApp struct {
		matProj unum.Mat4
	}
	thrPrep struct {
		tmpCamRender                bool
		onPrepNode                  func(*Node)
		matCamProj, matProj, matPos unum.Mat4
	}
	thrRend struct {
		states ugl.RenderStatesBag
	}
}

func newCamera2D(canv *RenderCanvas, depth bool) (me *Camera) {
	me = &Camera{}
	me.init(canv, false, depth, Options.Rendering.DefaultTechnique2D)
	return
}

func newCamera3D(canv *RenderCanvas) (me *Camera) {
	me = &Camera{}
	me.init(canv, true, true, Options.Rendering.DefaultTechnique3D)
	return
}

func newCameraQuad(canv *RenderCanvas) (me *Camera) {
	me = &Camera{}
	me.init(canv, false, false, Options.Rendering.DefaultTechniqueQuad)
	return
}

func (me *Camera) init(canv *RenderCanvas, persp3d bool, depth bool, technique string) {
	me.Enabled = true
	me.thrPrep.onPrepNode = func(n *Node) { me.onPrepNode(n) }
	rend := &me.Rendering
	rend.Viewport.canvas = canv
	rend.States.DepthTest, rend.States.FaceCulling, rend.States.StencilTest = depth, false, false
	rend.States.ClearColor = Options.Rendering.DefaultClearColor
	if depth {
		rend.States.Other.ClearBits = gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT
	} else {
		rend.States.Other.ClearBits = gl.COLOR_BUFFER_BIT
	}
	me.Perspective.Use, me.Perspective.CameraPerspective = persp3d, Options.Cameras.PerspectiveDefaults
	unum.Mat4Identities(&me.thrApp.matProj, &me.thrPrep.matProj)
	me.Controller.init()
	rend.Viewport.init()
	me.ApplyMatrices()
	me.Rendering.Technique = Core.Rendering.KnownTechniques[technique](me)
}

//	Applies changes made to the FovY, ZNear and/or ZFar parameters in me.Perspective.
func (me *Camera) ApplyMatrices() {
	me.thrApp.matProj.Perspective(me.Perspective.FovY, me.Rendering.Viewport.aspect, me.Perspective.ZNear, me.Perspective.ZFar)
}

func (me *Camera) clearNodeMats() {
	if me.scene != nil {
		me.scene.RootNode.Walk(func(node *Node) {
			delete(node.thrPrep.camProjMats, me)
			delete(node.thrRend.camProjMats, me)
			delete(node.thrPrep.camRender, me)
		})
	}
}

func (me *Camera) dispose() {
	me.clearNodeMats()
	me.Rendering.Technique.dispose()
	me.Rendering.Technique = nil
}

func (me *Camera) Scene() *Scene {
	return me.scene
}

func (me *Camera) setScene(scene *Scene) {
	if scene != me.scene {
		me.clearNodeMats()
		if me.scene = scene; me.scene != nil {
			me.scene.RootNode.Walk(func(node *Node) {
				node.initCamData(me)
			})
		}
	}
}

func (me *Camera) SetScene(id string) {
	me.setScene(Core.Libs.Scenes[id])
}

func (me *Camera) RenderTechniqueQuad() (tech *RenderTechniqueQuad) {
	tech, _ = me.Rendering.Technique.(*RenderTechniqueQuad)
	return
}

func (me *Camera) RenderTechniqueScene() (tech *RenderTechniqueScene) {
	tech, _ = me.Rendering.Technique.(*RenderTechniqueScene)
	return
}

type Cameras []*Camera

func (me *Cameras) dispose() {
	for _, cam := range *me {
		cam.dispose()
	}
	*me = nil
}

func (me *Cameras) Remove(camera *Camera) {
	for i, cam := range *me {
		if cam == camera {
			pre, post := (*me)[:i], (*me)[i+1:]
			*me = append(pre, post...)
		}
	}
}

//	Encapsulates a device-relative or absolute camera view-port.
type CameraViewport struct {
	relative, shouldScissor bool
	relX, relY, relW, relH  float64
	absX, absY, absW, absH  int
	aspect                  float64
	glVpX, glVpY            gl.Int
	glVpW, glVpH            gl.Sizei
	canvas                  *RenderCanvas
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
		me.absW, me.absH = int(me.relW*float64(me.canvas.absViewWidth)), int(me.relH*float64(me.canvas.absViewHeight))
		me.absX, me.absY = int(me.relX*float64(me.canvas.absViewWidth)), int(me.relY*float64(me.canvas.absViewHeight))
	}
	me.glVpX, me.glVpY, me.glVpW, me.glVpH = gl.Int(me.absX), gl.Int(me.absY), gl.Sizei(me.absW), gl.Sizei(me.absH)
	me.aspect = float64(me.absW) / float64(me.absH)
}
