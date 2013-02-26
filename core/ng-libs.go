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

func libElemIDRewrite(oldNewIDs map[int]int, ptr *int) {
	if newID, ok := oldNewIDs[*ptr]; ok {
		*ptr = newID
	}
}

func libElemIDRewrite2(oldNewIDs map[int]int, in int) (out int) {
	var ok bool
	if out, ok = oldNewIDs[in]; !ok {
		out = in
	}
	return
}

func (_ FxEffectLib) onFxEffectIDsChanged(oldNewIDs map[int]int) {
	Core.Libs.Materials.Walk(func(mat *FxMaterial) {
		libElemIDRewrite(oldNewIDs, &mat.DefaultEffectID)
		for id, _ := range mat.FaceEffects.ByID {
			mat.FaceEffects.ByID[id] = libElemIDRewrite2(oldNewIDs, mat.FaceEffects.ByID[id])
		}
		for tag, _ := range mat.FaceEffects.ByTag {
			mat.FaceEffects.ByTag[tag] = libElemIDRewrite2(oldNewIDs, mat.FaceEffects.ByTag[tag])
		}
	})
}

func (_ FxMaterialLib) onFxMaterialIDsChanged(oldNewIDs map[int]int) {
	Core.Libs.Scenes.Walk(func(scene *Scene) {
		scene.RootNode.Walk(func(node *Node) {
			libElemIDRewrite(oldNewIDs, &node.MatID)
		})
	})
	// for meshID, _ := range Core.Libs.Meshes {
	// 	for modelID, _ := range Core.Libs.Meshes[meshID].Models {
	// 		if newID, ok = oldNewIDs[model.MatID]; ok {
	// 			model.MatID = newID
	// 		}
	// 	}
	// }
}

func (_ SceneLib) onSceneIDsChanged(oldNewIDs map[int]int) {
	for _, canv := range Core.Rendering.Canvases {
		for _, cam := range canv.Cameras {
			libElemIDRewrite(oldNewIDs, &cam.sceneID)
		}
	}
}
