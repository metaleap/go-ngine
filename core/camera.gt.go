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
	ID int

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
	me.Rendering.Technique = Core.Rendering.KnownTechniques[technique](me)
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
	me.Rendering.Technique = nil
}

func (me *Camera) scene() *Scene {
	return Core.Libs.Scenes.Get(me.sceneID)
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

//#begin-gt -gen-lib.gt T:Camera L:Core.Rendering.Canvases[id].Cameras

//	Only used for Core.Rendering.Canvases[id].Cameras
type CameraLib []Camera

func (me *CameraLib) AddNew() (ref *Camera) {
	id := -1
	for i := 0; i < len(*me); i++ {
		if (*me)[i].ID < 0 {
			id = i
			break
		}
	}
	if id < 0 {
		if id = len(*me); id == cap(*me) {
			nu := make(CameraLib, id, id+Options.Libs.GrowCapBy)
			copy(nu, *me)
			*me = nu
		}
		*me = append(*me, Camera{})
	}
	ref = &(*me)[id]
	ref.ID = id
	ref.init()
	return
}

func (me *CameraLib) Compact() {
	var (
		before, after []Camera
		ref           *Camera
		oldID, i      int
	)
	for i = 0; i < len(*me); i++ {
		if (*me)[i].ID < 0 {
			before, after = (*me)[:i], (*me)[i+1:]
			*me = append(before, after...)
		}
	}
	changed := make(map[int]int, len(*me))
	for i = 0; i < len(*me); i++ {
		if (*me)[i].ID != i {
			ref = &(*me)[i]
			oldID, ref.ID = ref.ID, i
			changed[oldID] = i
		}
	}
	if len(changed) > 0 {
		me.onCameraIDsChanged(changed)
	}
}

func (me *CameraLib) ctor() {
	*me = make(CameraLib, 0, Options.Libs.InitialCap)
}

func (me *CameraLib) dispose() {
	me.Remove(0, 0)
	*me = (*me)[:0]
}

func (me CameraLib) Get(id int) (ref *Camera) {
	if id > -1 && id < len(me) {
		if ref = &me[id]; ref.ID != id {
			ref = nil
		}
	}
	return
}

func (me CameraLib) IsOk(id int) (ok bool) {
	if id > -1 && id < len(me) {
		ok = me[id].ID == id
	}
	return
}

func (me CameraLib) Ok(id int) bool {
	return me[id].ID == id
}

func (me CameraLib) Remove(fromID, num int) {
	if l := len(me); fromID < l {
		if num < 1 || num > (l-fromID) {
			num = l - fromID
		}
		changed := make(map[int]int, num)
		for id := fromID; id < fromID+num; id++ {
			me[id].dispose()
			changed[id], me[id].ID = -1, -1
		}
		me.onCameraIDsChanged(changed)
	}
}

func (me CameraLib) Walk(on func(ref *Camera)) {
	for id := 0; id < len(me); id++ {
		if me[id].ID > -1 {
			on(&me[id])
		}
	}
}

//#end-gt
