package core

func (_ *NgCore) onPrep() {
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
			scene.onPrep()
		}
		me.Camera.onPrep(scene.allNodes, 0)
	}
}

func (me *Camera) onPrep(all SceneNodeLib, nodeID int) {
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
	}
	for i := 0; i < len(all[nodeID].childNodeIDs); i++ {
		if all.IsOk(all[nodeID].childNodeIDs[i]) {
			me.onPrep(all, all[nodeID].childNodeIDs[i])
		}
	}
}

func (me *Scene) onPrep() {
	me.thrPrep.copyDone, me.thrRend.copyDone = false, false
}
