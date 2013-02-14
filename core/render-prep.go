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
			thrPrep.curScene.RootNode.Walk(me.thrPrep.onPrepNode)
			thrPrep.curScene.RootNode.onPrep()
		}
	}
}

func (me *Camera) onPrepNode(node *Node) {
	if node.Rendering.Enabled {
		if me.Perspective.Use {
			if node.Rendering.skyMode {
				node.thrPrep.matProjs[me].SetFromMult4(&me.thrPrep.matCamProj, &me.thrPrep.matPos)
			} else {
				node.thrPrep.matProjs[me].SetFromMult4(&me.thrPrep.matCamProj, &node.thrPrep.matModelView)
			}
		} else {
			node.thrPrep.matProjs[me].CopyFrom(&node.thrPrep.matModelView)
		}
	}

}

func (me *Node) onPrep() {
	if me.Rendering.Enabled && !me.thrPrep.done {
		me.thrPrep.done, me.thrPrep.copyDone, me.thrRend.copyDone = true, false, false
		for me.thrPrep.curId, me.thrPrep.curSubNode = range me.ChildNodes.M {
			me.thrPrep.curSubNode.onPrep()
		}
	}
}
