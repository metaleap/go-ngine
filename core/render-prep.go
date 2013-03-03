package core

func (_ *NgCore) onPrep() {
	// thrPrep.nodePreBatch.reset()
	for cid := 0; cid < len(Core.Render.Canvases); cid++ {
		if Core.Render.Canvases[cid].renderThisFrame() {
			Core.Render.Canvases[cid].onPrep()
		}
	}
}

func (me *RenderCanvas) onPrep() {
	for view := 0; view < len(me.Views); view++ {
		me.Views[view].onPrep()
	}
}

func (me *RenderView) onPrep() {
	if me.Enabled {
		me.Technique.onPrep()
	}
}

func (me *RenderTechniqueScene) onPrep() {
	me.Camera.thrPrep.matCamProj.SetFromMult4(&me.Camera.thrPrep.matProj, &me.Camera.Controller.thrPrep.mat)
	if scene := me.Camera.scene(); scene != nil {
		if !scene.thrPrep.done {
			scene.thrPrep.done = true
			scene.thrPrep.copyDone, scene.thrRend.copyDone = false, false
			scene.RootNode.onPrep()
		}
		scene.RootNode.Walk(me.Camera.thrPrep.onPrepNode)
		// if thrPrep.curTechScene = me.RenderTechniqueScene(); thrPrep.curTechScene != nil && thrPrep.curTechScene.Batch.Enabled {
		// 	thrPrep.curTechScene.Batch.onPrep()
		// }
	}
}

func (me *Camera) onPrepNode(node *Node) {
	camNodeRender := node.Rendering.Enabled && (node.parentNode == nil || node.parentNode.thrPrep.camRender[me]) // && inFrustum etc.
	if node.thrPrep.camRender[me] = camNodeRender; camNodeRender {
		if me.Perspective.Enabled {
			if node.Rendering.skyMode {
				node.thrPrep.camProjMats[me].SetFromMult4(&me.thrPrep.matCamProj, &me.thrPrep.matPos)
			} else {
				node.thrPrep.camProjMats[me].SetFromMult4(&me.thrPrep.matCamProj, &node.thrPrep.matModelView)
			}
		} else {
			node.thrPrep.camProjMats[me].CopyFrom(&node.thrPrep.matModelView)
		}
		// thrPrep.nodePreBatch.prepNode(node)
	}
}

func (me *Node) onPrep() {
	if me.Rendering.Enabled {
		for _, subNode := range me.ChildNodes.M {
			subNode.onPrep()
		}
	}
}
