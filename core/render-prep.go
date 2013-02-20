package core

func (me *EngineCore) onPrep() {
	thrPrep.nodePreBatch.reset()
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
			thrPrep.curScene.RootNode.Walk(me.thrPrep.onPrepNode)
			if thrPrep.curTechScene = me.RenderTechniqueScene(); thrPrep.curTechScene != nil && thrPrep.curTechScene.Batch.Enabled {
				thrPrep.curTechScene.Batch.onPrep()
			}
		}
	}
}

func (me *Camera) onPrepNode(node *Node) {
	me.thrPrep.tmpCamRender = node.Rendering.Enabled && (node.parentNode == nil || node.parentNode.thrPrep.camRender[me]) // && inFrustum etc.
	if node.thrPrep.camRender[me] = me.thrPrep.tmpCamRender; me.thrPrep.tmpCamRender {
		if me.Perspective.Use {
			if node.Rendering.skyMode {
				node.thrPrep.camProjMats[me].SetFromMult4(&me.thrPrep.matCamProj, &me.thrPrep.matPos)
			} else {
				node.thrPrep.camProjMats[me].SetFromMult4(&me.thrPrep.matCamProj, &node.thrPrep.matModelView)
			}
		} else {
			node.thrPrep.camProjMats[me].CopyFrom(&node.thrPrep.matModelView)
		}
		thrPrep.nodePreBatch.prepNode(node)
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
