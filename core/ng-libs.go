package core

type LibElemIDChangedHandler func(oldNewIDs map[int]int)

func (_ *EngineCore) disposeLibs() {
	for _, disp := range []disposable{
		&Core.Rendering.Canvases,
		&Core.Libs.Materials, &Core.Libs.Effects,
		&Core.Libs.Images.Tex2D, &Core.Libs.Images.TexCube,
		&Core.Libs.Meshes, Core.MeshBuffers, &Core.Libs.Scenes,
	} {
		disp.dispose()
	}
}

func (_ *EngineCore) initLibs() {
	libs := &Core.Libs
	for _, c := range []ctorable{
		&libs.Images.Tex2D, &libs.Images.TexCube,
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

func libOnFxImageIDsChanged(oldNewIDs map[int]int) {
	onOp := func(op FxOp) {
		if texOp, ok := op.(*fxOpTexBase); ok {
			libElemIDRewrite(oldNewIDs, &texOp.ImageID)
		}
	}
	var (
		fx *FxEffect
		op FxOp
	)
	for id, _ := range Core.Libs.Effects {
		if fx = Core.Libs.Effects.Get(id); fx != nil {
			for _, op = range fx.Ops {
				onOp(op)
			}
		}
	}
	for _, canv := range Core.Rendering.Canvases {
		for _, cam := range canv.Cams {
			for _, op = range cam.Rendering.FxOps {
				onOp(op)
			}
		}
	}
}

func (_ FxEffectLib) onFxEffectIDsChanged(oldNewIDs map[int]int) {
	var mat *FxMaterial
	for id, _ := range Core.Libs.Materials {
		if mat = Core.Libs.Materials.Get(id); mat != nil {
			libElemIDRewrite(oldNewIDs, &mat.DefaultEffectID)
			for id, _ := range mat.FaceEffects.ByID {
				mat.FaceEffects.ByID[id] = libElemIDRewrite2(oldNewIDs, mat.FaceEffects.ByID[id])
			}
			for tag, _ := range mat.FaceEffects.ByTag {
				mat.FaceEffects.ByTag[tag] = libElemIDRewrite2(oldNewIDs, mat.FaceEffects.ByTag[tag])
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
		mesh  *Mesh
	)
	for sceneID, _ := range Core.Libs.Scenes {
		if scene = Core.Libs.Scenes.Get(sceneID); scene != nil {
			scene.RootNode.Walk(func(node *Node) {
				libElemIDRewrite(oldNewIDs, &node.MatID)
			})
		}
	}
	for meshID, _ := range Core.Libs.Meshes {
		if mesh = Core.Libs.Meshes.Get(meshID); mesh != nil {
			for _, model := range Core.Libs.Meshes[meshID].Models {
				libElemIDRewrite(oldNewIDs, &model.MatID)
			}
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
				libElemIDRewrite(oldNewIDs, &node.MeshID)
			})
		}
	}
}

func (_ SceneLib) onSceneIDsChanged(oldNewIDs map[int]int) {
	for _, canv := range Core.Rendering.Canvases {
		for _, cam := range canv.Cams {
			libElemIDRewrite(oldNewIDs, &cam.sceneID)
		}
	}
}
