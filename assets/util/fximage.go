package assetutil

import (
	nga "github.com/go3d/go-ngine/assets"
)

//	Adds a new assets.FxImageDef with the specified Id to the assets.FxImageDefs
//	library, to be initialized from the specified refUrl.
func FxAddImage(id, refUrl string) *nga.FxImageDef {
	return nga.FxImageDefs.AddNewUrl(id, refUrl)
}
