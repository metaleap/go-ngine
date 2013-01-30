package core

//	prepares all renderBatches for the next onRender() call...
func (me *EngineCore) onPrep() {
	for _, thrPrep.curCanv = range me.Rendering.Canvases {
		if thrPrep.curCanv.renderThisFrame() {
			thrPrep.curCanv.onPrep()
		}
	}
}

func (me *RenderCanvas) onPrep() {
	for _, thrPrep.curCam = range me.Cameras {
		thrPrep.curCam.onPrep()
	}
}

func (me *Camera) onPrep() {
	if me.Enabled {
		me.thrPrep.matCamProj.SetFromMult4(&me.thrPrep.matProj, &me.Controller.thrPrep.mat)
		thrPrep.curScene = Core.Libs.Scenes[me.Rendering.SceneID]
		thrPrep.curScene.RootNode.onPrep()
	}
}

func (me *Node) onPrep() {
	if me.Enabled {
		if thrPrep.curNode = me; me.thrPrep.model != nil {
			if thrPrep.curCam.Perspective.Use {
				me.thrPrep.matModelProjs[thrPrep.curCam].SetFromMult4(&thrPrep.curCam.thrPrep.matCamProj, &me.thrPrep.matModelView)
			} else {
				*me.thrPrep.matModelProjs[thrPrep.curCam] = me.thrPrep.matModelView
			}
			me.thrPrep.glMatModelProjs[thrPrep.curCam].Load(me.thrPrep.matModelProjs[thrPrep.curCam])
		}
		for me.thrPrep.curId, me.thrPrep.curSubNode = range me.ChildNodes.M {
			me.thrPrep.curSubNode.onPrep()
		}
	}
}
