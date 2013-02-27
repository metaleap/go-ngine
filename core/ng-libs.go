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
		if len(ops) > 0 {
			var ok bool
			var texOp *fxOpTexBase
			for i := 0; i < len(ops); i++ {
				if texOp, ok = ops[i].(*fxOpTexBase); ok {
					libElemIDChangeRef(oldNewIDs, &texOp.ImageID)
				}
			}
		}
	}
	var id, cid int
	for id = 0; id < len(Core.Libs.Effects); id++ {
		if Core.Libs.Effects.Ok(id) {
			onOps(Core.Libs.Effects[id].Ops)
		}
	}
	for id = 0; id < len(Core.Rendering.Canvases); id++ {
		if Core.Rendering.Canvases.Ok(id) {
			for cid = 0; cid < len(Core.Rendering.Canvases[id].Cameras); cid++ {
				if Core.Rendering.Canvases[id].Cameras.Ok(cid) {
					onOps(Core.Rendering.Canvases[id].Cameras[cid].Rendering.FxOps)
				}
			}
		}
	}
}

func (_ FxEffectLib) onFxEffectIDsChanged(oldNewIDs map[int]int) {
	for mid := 0; mid < len(Core.Libs.Materials); mid++ {
		if Core.Libs.Materials.Ok(mid) {
			libElemIDChangeRef(oldNewIDs, &Core.Libs.Materials[mid].DefaultEffectID)
			if Core.Libs.Materials[mid].HasFaceEffects() {
				var m int
				var f string
				for f, m = range Core.Libs.Materials[mid].FaceEffects.ByID {
					Core.Libs.Materials[mid].FaceEffects.ByID[f] = libElemIDChangedRef(oldNewIDs, m)
				}
				for f, m = range Core.Libs.Materials[mid].FaceEffects.ByTag {
					Core.Libs.Materials[mid].FaceEffects.ByTag[f] = libElemIDChangedRef(oldNewIDs, m)
				}
			}
		}
	}
	Options.Libs.OnIDsChanged.Effects(oldNewIDs)
}

func (_ FxImage2DLib) onFxImage2DIDsChanged(oldNewIDs map[int]int) {
	libOnFxImageIDsChanged(oldNewIDs)
	Options.Libs.OnIDsChanged.Images.Tex2D(oldNewIDs)
}

func (_ FxImageCubeLib) onFxImageCubeIDsChanged(oldNewIDs map[int]int) {
	libOnFxImageIDsChanged(oldNewIDs)
	Options.Libs.OnIDsChanged.Images.TexCube(oldNewIDs)
}

func (_ FxMaterialLib) onFxMaterialIDsChanged(oldNewIDs map[int]int) {
	var id int
	for id = 0; id < len(Core.Libs.Scenes); id++ {
		if Core.Libs.Scenes.Ok(id) {
			Core.Libs.Scenes[id].RootNode.Walk(func(node *Node) {
				libElemIDChangeRef(oldNewIDs, &node.MatID)
			})
		}
	}
	for id = 0; id < len(Core.Libs.Models); id++ {
		if Core.Libs.Models.Ok(id) {
			libElemIDChangeRef(oldNewIDs, &Core.Libs.Models[id].MatID)
		}
	}
	Options.Libs.OnIDsChanged.Materials(oldNewIDs)
}

func (_ MeshLib) onMeshIDsChanged(oldNewIDs map[int]int) {
	for id := 0; id < len(Core.Libs.Scenes); id++ {
		if Core.Libs.Scenes.Ok(id) {
			Core.Libs.Scenes[id].RootNode.Walk(func(node *Node) {
				libElemIDChangeRef(oldNewIDs, &node.MeshID)
			})
		}
	}
	Options.Libs.OnIDsChanged.Meshes(oldNewIDs)
}

func (_ ModelLib) onModelIDsChanged(oldNewIDs map[int]int) {
	var id int
	for id = 0; id < len(Core.Libs.Scenes); id++ {
		if Core.Libs.Scenes.Ok(id) {
			Core.Libs.Scenes[id].RootNode.Walk(func(node *Node) {
				libElemIDChangeRef(oldNewIDs, &node.ModelID)
			})
		}
	}
	for id = 0; id < len(Core.Libs.Meshes); id++ {
		if Core.Libs.Meshes.Ok(id) {
			libElemIDChangeRef(oldNewIDs, &Core.Libs.Meshes[id].DefaultModelID)
		}
	}
	Options.Libs.OnIDsChanged.Models(oldNewIDs)
}

func (_ SceneLib) onSceneIDsChanged(oldNewIDs map[int]int) {
	var id, cid int
	for id = 0; id < len(Core.Rendering.Canvases); id++ {
		if Core.Rendering.Canvases.Ok(id) {
			for cid = 0; cid < len(Core.Rendering.Canvases[id].Cameras); cid++ {
				if Core.Rendering.Canvases[id].Cameras.Ok(cid) {
					libElemIDChangeRef(oldNewIDs, &Core.Rendering.Canvases[id].Cameras[cid].sceneID)
				}
			}
		}
	}
	Options.Libs.OnIDsChanged.Scenes(oldNewIDs)
}

func (_ RenderCanvasLib) onRenderCanvasIDsChanged(oldNewIDs map[int]int) {
	Options.Libs.OnIDsChanged.RenderCanvases(oldNewIDs)
}

func (_ CameraLib) onCameraIDsChanged(oldNewIDs map[int]int) {
	Options.Libs.OnIDsChanged.Cameras(oldNewIDs)
}
