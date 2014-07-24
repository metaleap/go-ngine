package exampleutils

import (
	ng "github.com/go3d/go-ngine/___old2013/core"
)

var LibIDs struct {
	Fx, Img2D, ImgCube, Mat map[string]int
}

func init() {
	const cap = 16
	LibIDs.Fx = make(map[string]int, cap)
	LibIDs.Mat = make(map[string]int, cap)
	LibIDs.Img2D = make(map[string]int, cap)
	LibIDs.ImgCube = make(map[string]int, 1)
}

func AddMainScene() (me *ng.Scene) {
	id := ng.Core.Libs.Scenes.AddNew()
	me = &ng.Core.Libs.Scenes[id]
	SceneCam.SetScene(id)
	return
}

func AddNode(scene *ng.Scene, parentNodeID, meshID, matID, modelID int) (node *ng.SceneNode) {
	nodeID := scene.AddNewChildNode(parentNodeID, meshID)
	node = scene.Node(nodeID)
	node.Render.MatID, node.Render.ModelID = matID, modelID
	return
}

func AddSkyMesh(scene *ng.Scene, meshID int) {
	cubeMapID := ng.Core.Libs.Images.TexCube.AddNew()
	cubeMap := &ng.Core.Libs.Images.TexCube[cubeMapID]
	LibIDs.ImgCube["sky"] = cubeMapID
	cubeMap.InitFrom[0].RefUrl = "tex/sky/east.png"  // positive X
	cubeMap.InitFrom[1].RefUrl = "tex/sky/west.png"  // negative X
	cubeMap.InitFrom[2].RefUrl = "tex/sky/up.png"    // positive Y
	cubeMap.InitFrom[3].RefUrl = "tex/sky/down.png"  // negative Y
	cubeMap.InitFrom[4].RefUrl = "tex/sky/north.png" // positive Z
	cubeMap.InitFrom[5].RefUrl = "tex/sky/south.png" // negative Z
	fxID := ng.Core.Libs.Effects.AddNew()
	LibIDs.Fx["sky"] = fxID
	fx := &ng.Core.Libs.Effects[fxID]
	fx.EnableTexCube(0).Tex.ImageID = cubeMap.ID
	fx.EnableCoords(0)
	fx.DisableCoords(0)
	fx.UpdateRoutine()
	matID := ng.Core.Libs.Materials.AddNew()
	ng.Core.Libs.Materials[matID].DefaultEffectID = LibIDs.Fx["sky"]
	LibIDs.Mat["sky"] = matID

	scene.SetNodeMeshID(0, meshID)
	scene.Root().Render.MatID = matID
}

//	Sets up textures and associated effects/materials with the specified IDs and image URLs.
//	For each ID (such as "foo" and "bar"):
//	-	creates an ng.FxImage2D with ID "img_{ID}" (ie. "img_foo" and "img_bar") and
//	adds it to ng.Core.Libs.Images.I2D
//	-	creates an ng.FxEffect with ID "fx_{ID}" (ie. "fx_foo" and "fx_bar") and adds it
//	to ng.Core.Libs.Effects; its Diffuse field pointing to the ng.FxImage2D
//	-	creates an ng.FxMaterial with ID "mat_{ID}" (ie. "mat_foo" and "mat_bar") and
//	adds it to ng.Core.Libs.Materials; its DefaultEffectID pointing to the ng.FxEffect.
func AddTextureMaterials(idsUrls map[string]string) {
	var (
		imgID, fxID, matID int
		fx                 *ng.FxEffect
	)
	for id, refUrl := range idsUrls {
		imgID = ng.Core.Libs.Images.Tex2D.AddNew()
		LibIDs.Img2D[id] = imgID
		ng.Core.Libs.Images.Tex2D[imgID].InitFrom.RefUrl = refUrl
		fxID = ng.Core.Libs.Effects.AddNew()
		fx = &ng.Core.Libs.Effects[fxID]
		LibIDs.Fx[id] = fxID
		fx.EnableTex2D(0).Tex.ImageID = imgID
		fx.EnableCoords(0)
		fx.DisableCoords(0)
		fx.UpdateRoutine()
		matID = ng.Core.Libs.Materials.AddNew()
		ng.Core.Libs.Materials[matID].DefaultEffectID = fxID
		LibIDs.Mat[id] = matID
	}
	return
}

func setupIDHandlers() {
	setup := func(handlers ng.LibElemIDChangedHandlers, namedIDs map[string]int) {
		handlers.Add(func(oldNewIDs map[int]int) {
			for k, v := range namedIDs {
				namedIDs[k] = ng.Core.Libs.UpdatedIDRef(oldNewIDs, v)
			}
		})
	}
	setup(ng.Options.Libs.OnIDsChanged.Effects, LibIDs.Fx)
	setup(ng.Options.Libs.OnIDsChanged.Images.Tex2D, LibIDs.Img2D)
	setup(ng.Options.Libs.OnIDsChanged.Images.TexCube, LibIDs.ImgCube)
	setup(ng.Options.Libs.OnIDsChanged.Materials, LibIDs.Mat)
}
