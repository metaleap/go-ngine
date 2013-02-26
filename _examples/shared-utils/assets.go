package exampleutils

import (
	ng "github.com/go3d/go-ngine/core"
)

var LibIDs = struct {
	Fx, Img, Mat map[string]int
}{map[string]int{}, map[string]int{}, map[string]int{}}

func AddMainScene() (me *ng.Scene) {
	me = ng.Core.Libs.Scenes.AddNew()
	SceneCam.SetScene(me.ID)
	return
}

func AddSkyMesh(scene *ng.Scene, meshID string) {
	cubeMap := ng.Core.Libs.Images.TexCube.AddNew()
	LibIDs.Img["sky"] = cubeMap.ID
	cubeMap.InitFrom[0].RefUrl = "tex/sky/east.png"  // positive X
	cubeMap.InitFrom[1].RefUrl = "tex/sky/west.png"  // negative X
	cubeMap.InitFrom[2].RefUrl = "tex/sky/up.png"    // positive Y
	cubeMap.InitFrom[3].RefUrl = "tex/sky/down.png"  // negative Y
	cubeMap.InitFrom[4].RefUrl = "tex/sky/north.png" // positive Z
	cubeMap.InitFrom[5].RefUrl = "tex/sky/south.png" // negative Z
	fx := ng.Core.Libs.Effects.AddNew()
	LibIDs.Fx["sky"] = fx.ID
	fx.Ops.EnableTexCube(0).ImageID = cubeMap.ID
	fx.Ops.EnableCoords(0)
	fx.Ops.DisableCoords(0)
	fx.UpdateRoutine()
	matSky := ng.Core.Libs.Materials.AddNew()
	matSky.DefaultEffectID = LibIDs.Fx["sky"]
	LibIDs.Mat["sky"] = matSky.ID

	ng.Core.Libs.Meshes[meshID].Models.Default().Clone("meshmodel_skybox").MatID = matSky.ID
	scene.RootNode.SetMeshModelID(meshID, "meshmodel_skybox")
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
		img *ng.FxImage2D
		fx  *ng.FxEffect
		mat *ng.FxMaterial
	)
	for id, refUrl := range idsUrls {
		img = ng.Core.Libs.Images.Tex2D.AddNew()
		LibIDs.Img[id] = img.ID
		img.InitFrom.RefUrl = refUrl
		fx = ng.Core.Libs.Effects.AddNew()
		LibIDs.Fx[id] = fx.ID
		fx.Ops.EnableTex2D(0).ImageID = img.ID
		fx.Ops.EnableCoords(0)
		fx.Ops.DisableCoords(0)
		fx.UpdateRoutine()
		mat = ng.Core.Libs.Materials.AddNew()
		mat.DefaultEffectID = fx.ID
		LibIDs.Mat[id] = mat.ID
	}
	return
}
