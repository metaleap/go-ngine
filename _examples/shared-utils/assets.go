package exampleutils

import (
	ng "github.com/go3d/go-ngine/core"
)

//	Creates a new core.Scene, adds it to the Scenes library under
//	the specified ID, and returns it.
func AddScene(id string, mainCamScene bool, skyMeshID string) (me *ng.Scene) {
	me = ng.NewScene()
	if ng.Core.Libs.Scenes[id] = me; mainCamScene {
		SceneCam.SetScene(id)
	}
	if len(skyMeshID) > 0 {
		cubeMap := ng.Core.Libs.Images.TexCubes.AddNew("img_sky")
		cubeMap.InitFrom[0].RefUrl = "tex/sky/east.png"  // positive X
		cubeMap.InitFrom[1].RefUrl = "tex/sky/west.png"  // negative X
		cubeMap.InitFrom[2].RefUrl = "tex/sky/up.png"    // positive Y
		cubeMap.InitFrom[3].RefUrl = "tex/sky/down.png"  // negative Y
		cubeMap.InitFrom[4].RefUrl = "tex/sky/north.png" // positive Z
		cubeMap.InitFrom[5].RefUrl = "tex/sky/south.png" // negative Z
		fx := ng.Core.Libs.Effects.AddNew("fx_sky")
		fx.Ops.EnableTexCube(0).SetImageID("img_sky")
		fx.Ops.EnableColored(0)
		fx.Ops.DisableColored(0)
		fx.UpdateRoutine()
		ng.Core.Libs.Materials.AddNew("mat_sky").DefaultEffectID = "fx_sky"

		ng.Core.Libs.Meshes[skyMeshID].Models.Default().Clone("meshmodel_skybox").SetMatID("mat_sky")
		me.RootNode.SetMeshModelID(skyMeshID, "meshmodel_skybox")
	}
	return
}

/*
//	Sets up plain-color effects/materials with the specified IDs.
//	For each ID (such as "foo" and "bar"):
//	-	creates an ng.FxEffect with ID "fx_{ID}" (ie. "fx_foo" and "fx_bar") and adds it
//	to ng.Core.Libs.Effects; its Diffuse field pointing to the color
//	-	creates an ng.FxMaterial with ID "mat_{ID}" (ie. "mat_foo" and "mat_bar") and
//	adds it to ng.Core.Libs.Materials; its DefaultEffectID pointing to the ng.FxEffect.
// func AddColorMaterials(idsColors map[string][]float64) {
// 	for id, col := range idsColors {
// 		ng.Core.Libs.Effects.AddNew("fx_" + id).OldDiffuse = ng.NewFxColor(col...)
// 		ng.Core.Libs.Materials.AddNew("mat_" + id).DefaultEffectID = "fx_" + id
// 	}
// }
*/

//	Sets up textures and associated effects/materials with the specified IDs and image URLs.
//	For each ID (such as "foo" and "bar"):
//	-	creates an ng.FxImage2D with ID "img_{ID}" (ie. "img_foo" and "img_bar") and
//	adds it to ng.Core.Libs.Images.I2D
//	-	creates an ng.FxEffect with ID "fx_{ID}" (ie. "fx_foo" and "fx_bar") and adds it
//	to ng.Core.Libs.Effects; its Diffuse field pointing to the ng.FxImage2D
//	-	creates an ng.FxMaterial with ID "mat_{ID}" (ie. "mat_foo" and "mat_bar") and
//	adds it to ng.Core.Libs.Materials; its DefaultEffectID pointing to the ng.FxEffect.
func AddTextureMaterials(idsUrls map[string]string) {
	for id, refUrl := range idsUrls {
		img := ng.Core.Libs.Images.Tex2D.AddNew("img_" + id)
		img.InitFrom.RefUrl = refUrl
		fx := ng.Core.Libs.Effects.AddNew("fx_" + id)
		fx.Ops.EnableTex2D(0).SetImageID("img_" + id)
		fx.Ops.EnableColored(0)
		fx.Ops.DisableColored(0)
		fx.UpdateRoutine()
		ng.Core.Libs.Materials.AddNew("mat_" + id).DefaultEffectID = "fx_" + id
	}
}
