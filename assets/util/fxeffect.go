package assetutil

import (
	nga "github.com/go3d/go-ngine/assets"
)

func EffectAddProfileCommon(def *nga.FxEffectDef) {
	var prof = nga.NewProfile()
	prof.Common = &nga.FxProfileCommon{}
	def.Profiles = append(def.Profiles, prof)
}
