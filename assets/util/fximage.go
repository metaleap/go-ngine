package assetutil

import (
	nga "github.com/go3d/go-ngine/assets"
)

//	Creates and returns a new assets.FxImageDef with the specified Id,
//	to be initialized from the specified refUrl.
func NewFxImageDef(id, refUrl string) (me *nga.FxImageDef) {
	me = nga.FxImageDefs.New(id)
	me.InitFrom = nga.NewFxImageInitFrom(refUrl)
	return
}
