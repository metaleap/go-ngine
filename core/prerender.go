package core

//	prepares all renderBatches for the next onRender() call...
func (me *EngineCore) onPrep() {
	for _, curCanvas = range me.Rendering.Canvases {
		if curCanvas.renderThisFrame() {
			curCanvas.onPrep()
		}
	}
}

func (me *RenderCanvas) onPrep() {
	for _, curCam = range me.Cameras {
		curCam.onPrep()
	}
}

func (me *Camera) onPrep() {
	if me.Rendering.Enabled {
		me.matCamProj.SetFromMult4(&me.matProj, &me.Controller.mat)
		curScene = Core.Libs.Scenes[me.Rendering.SceneID]
		curScene.RootNode.onPrep()
	}
}

func (me *Node) onPrep() {
	if !me.Disabled {
		if curNode, curMesh, curModel = me, me.mesh, me.model; curModel != nil {
			if curCam.Perspective.Use {
				me.matModelProjs[curCam].SetFromMult4(&curCam.matCamProj, &me.Transform.matModelView)
			} else {
				*me.matModelProjs[curCam] = me.Transform.matModelView
			}
			me.glMatModelProjs[curCam].Load(me.matModelProjs[curCam])
		}
		for me.curKey, me.curSubNode = range me.ChildNodes.M {
			me.curSubNode.onPrep()
		}
	}
}

func (me *Model) preRender() {
	me.mesh.preRender()
}

func (me *Mesh) preRender() {
}
