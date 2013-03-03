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
		for i := 0; i < len(scene.allNodes); i++ {
			if scene.allNodes.Ok(i) {
				me.Camera.onPrepNode(scene, &scene.allNodes[i])
			}
		}
		// if thrPrep.curTechScene = me.RenderTechniqueScene(); thrPrep.curTechScene != nil && thrPrep.curTechScene.Batch.Enabled {
		// 	thrPrep.curTechScene.Batch.onPrep()
		// }
	}
}

func (me *Camera) onPrepNode(scene *Scene, node *SceneNode) {
	camNodeRender := node.Render.Enabled && (node.parentID == -1 || me.thrPrep.nodeRender[&scene.allNodes[node.parentID]]) // && inFrustum etc.
	if me.thrPrep.nodeRender[node] = camNodeRender; camNodeRender {
		if me.Perspective.Enabled {
			if node.Render.skyMode {
				me.thrPrep.nodeProjMats[node].SetFromMult4(&me.thrPrep.matCamProj, &me.thrPrep.matPos)
			} else {
				me.thrPrep.nodeProjMats[node].SetFromMult4(&me.thrPrep.matCamProj, &node.Transform.thrPrep.matModelView)
			}
		} else {
			me.thrPrep.nodeProjMats[node].CopyFrom(&node.Transform.thrPrep.matModelView)
		}
		// thrPrep.nodePreBatch.prepNode(node)
	}
}
