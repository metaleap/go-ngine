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
	if scene := me.Camera.Scene(); scene != nil {
		if !scene.thrPrep.done {
			scene.thrPrep.done = true
			scene.thrPrep.copyDone, scene.thrRend.copyDone = false, false
		}
		me.Camera.onPrepNode(scene.allNodes, 0)
		// if thrPrep.curTechScene = me.RenderTechniqueScene(); thrPrep.curTechScene != nil && thrPrep.curTechScene.Batch.Enabled {
		// 	thrPrep.curTechScene.Batch.onPrep()
		// }
	}
}

func (me *Camera) onPrepNode(all SceneNodeLib, nodeID int) {
	camNodeRender := all[nodeID].Render.Enabled && (nodeID == 0 || me.thrPrep.nodeRender[all[nodeID].parentID]) // && inFrustum etc.
	if me.thrPrep.nodeRender[nodeID] = camNodeRender; camNodeRender {
		if me.Perspective.Enabled {
			if all[nodeID].Render.skyMode {
				me.thrPrep.nodeProjMats[nodeID].SetFromMult4(&me.thrPrep.matCamProj, &me.thrPrep.matPos)
			} else {
				me.thrPrep.nodeProjMats[nodeID].SetFromMult4(&me.thrPrep.matCamProj, &all[nodeID].Transform.thrPrep.matModelView)
			}
		} else {
			me.thrPrep.nodeProjMats[nodeID].CopyFrom(&all[nodeID].Transform.thrPrep.matModelView)
		}
		// thrPrep.nodePreBatch.prepNode(node)
	}
	for i := 0; i < len(all[nodeID].childNodeIDs); i++ {
		if all.IsOk(all[nodeID].childNodeIDs[i]) {
			me.onPrepNode(all, all[nodeID].childNodeIDs[i])
		}
	}
}
