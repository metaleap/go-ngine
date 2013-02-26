package core

type LibElemIDChangedHandler func(oldNewIDs map[int]int)

func (_ *EngineCore) disposeLibs() {
	for _, disp := range []disposable{
		&Core.Rendering.Canvases,
		&Core.Libs.Materials, &Core.Libs.Effects,
		&Core.Libs.Images.Tex2D, &Core.Libs.Images.TexCubes,
		&Core.Libs.Meshes, Core.MeshBuffers, &Core.Libs.Scenes,
	} {
		disp.dispose()
	}
}

func (_ *EngineCore) initLibs() {
	libs := &Core.Libs
	for _, c := range []ctorable{
		&libs.Images.Tex2D, &libs.Images.TexCubes,
		&libs.Effects, &libs.Materials, &libs.Meshes, &libs.Scenes,
	} {
		c.ctor()
	}
}

func (_ FxMaterialLib) onFxMaterialIDsChanged(oldNewIDs map[int]int) {
	var (
		ok    bool
		newID int
	)
	for _, scene := range Core.Libs.Scenes {
		scene.RootNode.Walk(func(node *Node) {
			if newID, ok = oldNewIDs[node.MatID]; ok {
				node.MatID = newID
			}
		})
	}
	for _, mesh := range Core.Libs.Meshes {
		for _, model := range mesh.Models {
			if newID, ok = oldNewIDs[model.MatID]; ok {
				model.MatID = newID
			}
		}
	}
}

func (_ SceneLib) onSceneIDsChanged(oldNewIDs map[int]int) {
	var (
		ok    bool
		newID int
	)
	for _, canv := range Core.Rendering.Canvases {
		for _, cam := range canv.Cameras {
			if newID, ok = oldNewIDs[cam.sceneID]; ok {
				cam.sceneID = newID
			}
		}
	}
}
