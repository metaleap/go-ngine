package core

type LibElemIDChangedHandler func(oldNewIDs map[int]int)

type LibElemIDChangedHandlers []LibElemIDChangedHandler

func (me *LibElemIDChangedHandlers) Add(fn LibElemIDChangedHandler) {
	if fn != nil {
		*me = append(*me, fn)
	}
}

func (me LibElemIDChangedHandlers) callAll(oldNewIDs map[int]int) {
	for _, fn := range me {
		fn(oldNewIDs)
	}
}

//	Only used for Core.Libs.
type EngineLibs struct {
	Effects   FxEffectLib
	Materials FxMaterialLib
	Images    struct {
		SplashScreen FxImage2D
		TexCube      FxImageCubeLib
		Tex2D        FxImage2DLib
	}
	Meshes MeshLib
	Models ModelLib
	Scenes SceneLib
}

func (me *EngineLibs) dispose() {
	for _, disp := range []interface {
		dispose()
	}{
		&Core.Render.Canvases, &Core.MeshBuffers,
		&me.Models, &me.Materials, &me.Effects,
		&me.Images.Tex2D, &me.Images.TexCube,
		&me.Meshes, &me.Scenes,
	} {
		disp.dispose()
	}
}

func (me *EngineLibs) init() {
	for _, c := range []interface {
		init()
	}{
		&Core.Render.Canvases, &Core.MeshBuffers,
		&me.Models, &me.Materials, &me.Effects,
		&me.Images.Tex2D, &me.Images.TexCube,
		&me.Meshes, &me.Scenes,
	} {
		c.init()
	}
}

func (_ *EngineLibs) UpdateIDRef(oldNewIDs map[int]int, ptr *int) {
	if newID, ok := oldNewIDs[*ptr]; ok {
		*ptr = newID
	}
}

func (_ *EngineLibs) UpdateIDRefs(oldNewIDs map[int]int, ptrs ...*int) {
	var newID int
	var ok bool
	for _, ptr := range ptrs {
		if newID, ok = oldNewIDs[*ptr]; ok {
			*ptr = newID
		}
	}
}

func (_ *EngineLibs) UpdatedIDRef(oldNewIDs map[int]int, in int) (out int) {
	var ok bool
	if out, ok = oldNewIDs[in]; !ok {
		out = in
	}
	return
}

func (_ *EngineLibs) onFxImageIDsChanged(procID string, oldNewIDs map[int]int) {
	onOps := func(ops FxProcs) {
		if len(ops) > 0 {
			for i := 0; i < len(ops); i++ {
				if ops[i].procID == procID {
					Core.Libs.UpdateIDRef(oldNewIDs, &ops[i].Tex.ImageID)
				}
			}
		}
	}
	var id, camID int
	for id = 0; id < len(Core.Libs.Effects); id++ {
		if Core.Libs.Effects.Ok(id) {
			onOps(Core.Libs.Effects[id].FxProcs)
		}
	}
	for id = 0; id < len(Core.Render.Canvases); id++ {
		for camID = 0; camID < len(Core.Render.Canvases[id].Cams); camID++ {
			onOps(Core.Render.Canvases[id].Cams[camID].Render.FxProcs)
		}
	}
}

func (_ FxEffectLib) onFxEffectIDsChanged(oldNewIDs map[int]int) {
	for mid := 0; mid < len(Core.Libs.Materials); mid++ {
		if Core.Libs.Materials.Ok(mid) {
			Core.Libs.UpdateIDRef(oldNewIDs, &Core.Libs.Materials[mid].DefaultEffectID)
			if Core.Libs.Materials[mid].HasFaceEffects() {
				var m int
				var f string
				for f, m = range Core.Libs.Materials[mid].FaceEffects.ByID {
					Core.Libs.Materials[mid].FaceEffects.ByID[f] = Core.Libs.UpdatedIDRef(oldNewIDs, m)
				}
				for f, m = range Core.Libs.Materials[mid].FaceEffects.ByTag {
					Core.Libs.Materials[mid].FaceEffects.ByTag[f] = Core.Libs.UpdatedIDRef(oldNewIDs, m)
				}
			}
		}
	}
	Options.Libs.OnIDsChanged.Effects.callAll(oldNewIDs)
}

func (_ FxImage2DLib) onFxImage2DIDsChanged(oldNewIDs map[int]int) {
	Core.Libs.onFxImageIDsChanged("Tex2D", oldNewIDs)
	Options.Libs.OnIDsChanged.Images.Tex2D.callAll(oldNewIDs)
}

func (_ FxImageCubeLib) onFxImageCubeIDsChanged(oldNewIDs map[int]int) {
	Core.Libs.onFxImageIDsChanged("TexCube", oldNewIDs)
	Options.Libs.OnIDsChanged.Images.TexCube.callAll(oldNewIDs)
}

func (_ FxMaterialLib) onFxMaterialIDsChanged(oldNewIDs map[int]int) {
	var id, v int
	for id = 0; id < len(Core.Libs.Scenes); id++ {
		if Core.Libs.Scenes.Ok(id) {
			Core.Libs.Scenes[id].RootNode.Walk(func(node *Node) {
				Core.Libs.UpdateIDRef(oldNewIDs, &node.MatID)
			})
			for v = 0; v < len(Core.Libs.Scenes[id].allNodes); v++ {
				if Core.Libs.Scenes[id].allNodes.Ok(v) {
					Core.Libs.UpdateIDRef(oldNewIDs, &Core.Libs.Scenes[id].allNodes[v].Render.MatID)
				}
			}
		}
	}
	for id = 0; id < len(Core.Libs.Models); id++ {
		if Core.Libs.Models.Ok(id) {
			Core.Libs.UpdateIDRef(oldNewIDs, &Core.Libs.Models[id].MatID)
		}
	}
	Options.Libs.OnIDsChanged.Materials.callAll(oldNewIDs)
}

func (_ MeshLib) onMeshIDsChanged(oldNewIDs map[int]int) {
	var id, v int
	for _, meshBuf := range Core.MeshBuffers {
		for id, v = range meshBuf.meshIDs {
			meshBuf.meshIDs[id] = Core.Libs.UpdatedIDRef(oldNewIDs, v)
		}
	}
	for id = 0; id < len(Core.Libs.Scenes); id++ {
		if Core.Libs.Scenes.Ok(id) {
			Core.Libs.Scenes[id].RootNode.Walk(func(node *Node) {
				Core.Libs.UpdateIDRef(oldNewIDs, &node.MeshID)
			})
			for v = 0; v < len(Core.Libs.Scenes[id].allNodes); v++ {
				if Core.Libs.Scenes[id].allNodes.Ok(v) {
					Core.Libs.UpdateIDRef(oldNewIDs, &Core.Libs.Scenes[id].allNodes[v].Render.MeshID)
				}
			}
		}
	}
	Options.Libs.OnIDsChanged.Meshes.callAll(oldNewIDs)
}

func (_ ModelLib) onModelIDsChanged(oldNewIDs map[int]int) {
	var id, i int
	for id = 0; id < len(Core.Libs.Scenes); id++ {
		if Core.Libs.Scenes.Ok(id) {
			Core.Libs.Scenes[id].RootNode.Walk(func(node *Node) {
				Core.Libs.UpdateIDRef(oldNewIDs, &node.ModelID)
			})
			for i = 0; i < len(Core.Libs.Scenes[id].allNodes); i++ {
				if Core.Libs.Scenes[id].allNodes.Ok(i) {
					Core.Libs.UpdateIDRef(oldNewIDs, &Core.Libs.Scenes[id].allNodes[i].Render.ModelID)
				}
			}
		}
	}
	for id = 0; id < len(Core.Libs.Meshes); id++ {
		if Core.Libs.Meshes.Ok(id) {
			Core.Libs.UpdateIDRef(oldNewIDs, &Core.Libs.Meshes[id].DefaultModelID)
		}
	}
	Options.Libs.OnIDsChanged.Models.callAll(oldNewIDs)
}

func (_ SceneLib) onSceneIDsChanged(oldNewIDs map[int]int) {
	var id, camID int
	for id = 0; id < len(Core.Render.Canvases); id++ {
		for camID = 0; camID < len(Core.Render.Canvases[id].Cams); camID++ {
			Core.Libs.UpdateIDRef(oldNewIDs, &Core.Render.Canvases[id].Cams[camID].sceneID)
		}
	}
	Options.Libs.OnIDsChanged.Scenes.callAll(oldNewIDs)
}

func (me *SceneNodeLib) onSceneNodeIDsChanged(oldNewIDs map[int]int) {
	var i int
	for i = 0; i < len(*me); i++ {
		if me.Ok(i) {
			Core.Libs.UpdateIDRef(oldNewIDs, &(*me)[i].parentID)
		}
	}
	for id := 0; id < len(Core.Libs.Scenes); id++ {
		if Core.Libs.Scenes.Ok(id) && &Core.Libs.Scenes[id].allNodes == me {
			break
		}
	}
}
