package core

type LibElemIDChangedHandler func(oldNewIDs map[int]int)

func (_ *EngineCore) disposeLibs() {
	for _, disp := range []disposable{
		&Core.Rendering.Canvases, Core.MeshBuffers,
		&Core.Libs.Materials, &Core.Libs.Effects,
		&Core.Libs.Images.Tex2D, &Core.Libs.Images.TexCube,
		&Core.Libs.Meshes, &Core.Libs.Models, &Core.Libs.Scenes,
	} {
		disp.dispose()
	}
}

func (_ *EngineCore) initLibs() {
	libs := &Core.Libs
	for _, c := range []ctorable{
		&libs.Images.Tex2D, &libs.Images.TexCube,
		&libs.Effects, &libs.Materials, &libs.Meshes, &libs.Models, &libs.Scenes,
	} {
		c.ctor()
	}
}

func libCreateChangedIDsList(num int) (changedIDs []int) {
	changedIDs = make([]int, num)
	for i := 0; i < len(changedIDs); i++ {
		changedIDs[i] = -2
	}
	return
}

func libElemIDChangeRef(oldNewIDs map[int]int, ptr *int) {
	if newID, ok := oldNewIDs[*ptr]; ok {
		*ptr = newID
	}
}

func libElemIDChangedRef(oldNewIDs map[int]int, in int) (out int) {
	var ok bool
	if out, ok = oldNewIDs[in]; !ok {
		out = in
	}
	return
}

func libOnFxImageIDsChanged(oldNewIDs map[int]int) {
	onOps := func(ops FxOps) {
		var ok bool
		var texOp *fxOpTexBase
		for _, op := range ops {
			if texOp, ok = op.(*fxOpTexBase); ok {
				libElemIDChangeRef(oldNewIDs, &texOp.ImageID)
			}
		}
	}
	var fx *FxEffect
	for id, _ := range Core.Libs.Effects {
		if fx = Core.Libs.Effects.Get(id); fx != nil {
			onOps(fx.Ops)
		}
	}
	for _, canv := range Core.Rendering.Canvases {
		for _, cam := range canv.Cams {
			onOps(cam.Rendering.FxOps)
		}
	}
}

func (_ FxEffectLib) onFxEffectIDsChanged(oldNewIDs map[int]int) {
	var mat *FxMaterial
	for id, _ := range Core.Libs.Materials {
		if mat = Core.Libs.Materials.Get(id); mat != nil {
			libElemIDChangeRef(oldNewIDs, &mat.DefaultEffectID)
			for id, _ := range mat.FaceEffects.ByID {
				mat.FaceEffects.ByID[id] = libElemIDChangedRef(oldNewIDs, mat.FaceEffects.ByID[id])
			}
			for tag, _ := range mat.FaceEffects.ByTag {
				mat.FaceEffects.ByTag[tag] = libElemIDChangedRef(oldNewIDs, mat.FaceEffects.ByTag[tag])
			}
		}
	}
}

func (_ FxImage2DLib) onFxImage2DIDsChanged(oldNewIDs map[int]int) {
	libOnFxImageIDsChanged(oldNewIDs)
}

func (_ FxImageCubeLib) onFxImageCubeIDsChanged(oldNewIDs map[int]int) {
	libOnFxImageIDsChanged(oldNewIDs)
}

func (_ FxMaterialLib) onFxMaterialIDsChanged(oldNewIDs map[int]int) {
	var (
		scene *Scene
		model *Model
	)
	for sceneID, _ := range Core.Libs.Scenes {
		if scene = Core.Libs.Scenes.Get(sceneID); scene != nil {
			scene.RootNode.Walk(func(node *Node) {
				libElemIDChangeRef(oldNewIDs, &node.MatID)
			})
		}
	}
	for modelID, _ := range Core.Libs.Models {
		if model = Core.Libs.Models.Get(modelID); model != nil {
			libElemIDChangeRef(oldNewIDs, &model.MatID)
		}
	}
}

func (_ MeshLib) onMeshIDsChanged(oldNewIDs map[int]int) {
	var (
		scene *Scene
	)
	for sceneID, _ := range Core.Libs.Scenes {
		if scene = Core.Libs.Scenes.Get(sceneID); scene != nil {
			scene.RootNode.Walk(func(node *Node) {
				libElemIDChangeRef(oldNewIDs, &node.MeshID)
			})
		}
	}
}

func (_ ModelLib) onModelIDsChanged(oldNewIDs map[int]int) {
	var scene *Scene
	for sceneID, _ := range Core.Libs.Scenes {
		if scene = Core.Libs.Scenes.Get(sceneID); scene != nil {
			scene.RootNode.Walk(func(node *Node) {
				libElemIDChangeRef(oldNewIDs, &node.ModelID)
			})
		}
	}
	var mesh *Mesh
	for meshID, _ := range Core.Libs.Meshes {
		if mesh = Core.Libs.Meshes.Get(meshID); mesh != nil {
			libElemIDChangeRef(oldNewIDs, &mesh.DefaultModelID)
		}
	}
}

func (_ SceneLib) onSceneIDsChanged(oldNewIDs map[int]int) {
	for _, canv := range Core.Rendering.Canvases {
		for _, cam := range canv.Cams {
			libElemIDChangeRef(oldNewIDs, &cam.sceneID)
		}
	}
}
