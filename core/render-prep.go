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
		// thrPrep.curScene = Core.Libs.Scenes[me.Rendering.SceneID]
		// thrPrep.curScene.RootNode.onPrep()
	}
}

// func (me *Node) onPrep() {
// 	if me.Enabled {
// 		for me.thrPrep.curId, me.thrPrep.curSubNode = range me.ChildNodes.M {
// 			me.thrPrep.curSubNode.onPrep()
// 		}
// 	}
// }
