package core

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
		if thrPrep.curScene = me.scene; thrPrep.curScene != nil {
			thrPrep.curScene.RootNode.onPrep()
		}
	}
}

func (me *Node) onPrep() {
	if me.Rendering.Enabled && !me.thrPrep.done {
		me.thrPrep.done, me.thrPrep.copyDone, me.thrRend.copyDone = true, false, false
		for me.thrPrep.tmpCam, me.thrPrep.tmpMat = range me.thrPrep.matProjs {
			if me.thrPrep.tmpCam.Perspective.Use {
				if me.Rendering.skyMode {
					me.thrPrep.tmpMat.SetFromMult4(&me.thrPrep.tmpCam.thrPrep.matCamProj, &me.thrPrep.tmpCam.thrPrep.matCtlPos)
				} else {
					me.thrPrep.tmpMat.SetFromMult4(&me.thrPrep.tmpCam.thrPrep.matCamProj, &me.thrPrep.matModelView)
				}
			} else {
				me.thrPrep.tmpMat.CopyFrom(&me.thrPrep.matModelView)
			}
		}
		for me.thrPrep.curId, me.thrPrep.curSubNode = range me.ChildNodes.M {
			me.thrPrep.curSubNode.onPrep()
		}
	}
}
