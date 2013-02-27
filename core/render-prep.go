package core

func (_ *EngineCore) onPrep() {
	// thrPrep.nodePreBatch.reset()
	for cid := 0; cid < len(Core.Rendering.Canvases); cid++ {
		if Core.Rendering.Canvases.Ok(cid) && Core.Rendering.Canvases[cid].renderThisFrame() {
			Core.Rendering.Canvases[cid].onPrep()
		}
	}
}

func (me *RenderCanvas) onPrep() {
	for cam := 0; cam < len(me.Cameras); cam++ {
		if me.Cameras.Ok(cam) {
			me.Cameras[cam].onPrep()
		}
	}
}

func (me *Camera) onPrep() {
	if me.Enabled {
		me.thrPrep.matCamProj.SetFromMult4(&me.thrPrep.matProj, &me.Controller.thrPrep.mat)
		if scene := me.scene(); scene != nil {
			scene.RootNode.onPrep()
			scene.RootNode.Walk(me.thrPrep.onPrepNode)
			// if thrPrep.curTechScene = me.RenderTechniqueScene(); thrPrep.curTechScene != nil && thrPrep.curTechScene.Batch.Enabled {
			// 	thrPrep.curTechScene.Batch.onPrep()
			// }
		}
	}
}

func (me *Camera) onPrepNode(node *Node) {
	camNodeRender := node.Rendering.Enabled && (node.parentNode == nil || node.parentNode.thrPrep.camRender[me]) // && inFrustum etc.
	if node.thrPrep.camRender[me] = camNodeRender; camNodeRender {
		if me.Perspective.Use {
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
	if me.Rendering.Enabled && !me.thrPrep.done {
		me.thrPrep.done, me.thrPrep.copyDone, me.thrRend.copyDone = true, false, false
		for _, subNode := range me.ChildNodes.M {
			subNode.onPrep()
		}
	}
}
