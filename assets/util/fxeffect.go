package assetutil

import (
	nga "github.com/go3d/go-ngine/assets"
)

//	Adds a new assets.FxEffectDef with the specified Id to the assets.FxEffectDefs
//	library, and optionally adds a common profile, a GLSL profile, or both to it.
func FxAddEffect(id string, ensureCommonProfile, ensureGlslProfile bool) (me *nga.FxEffectDef) {
	me = nga.FxEffectDefs.AddNew(id)
	if ensureCommonProfile {
		FxEnsureProfileCommon(me)
	}
	if ensureGlslProfile {
		FxEnsureProfileGlsl(me)
	}
	return
}

//	Adds a new common profile to the specified effect definition (even if it already contains one).
func FxAddProfileCommon(def *nga.FxEffectDef) (prof *nga.FxProfile) {
	prof = nga.NewProfile()
	prof.Common = &nga.FxProfileCommon{}
	def.Profiles = append(def.Profiles, prof)
	return
}

//	Adds a new GLSL profile to the specified effect definition (even if it already contains one).
func FxAddProfileGlsl(def *nga.FxEffectDef) (prof *nga.FxProfile) {
	prof = nga.NewProfile()
	prof.Glsl = nga.NewFxProfileGlsl()
	def.Profiles = append(def.Profiles, prof)
	return
}

//	Ensures the specified effect definition contains a common profile and returns it.
func FxEnsureProfileCommon(def *nga.FxEffectDef) (prof *nga.FxProfile) {
	for _, prof = range def.Profiles {
		if prof.Common != nil {
			break
		}
		prof = nil
	}
	if prof == nil {
		prof = FxAddProfileCommon(def)
	}
	return
}

//	Ensures the specified effect definition contains a GLSL profile and returns it.
func FxEnsureProfileGlsl(def *nga.FxEffectDef) (prof *nga.FxProfile) {
	for _, prof = range def.Profiles {
		if prof.Glsl != nil {
			break
		}
		prof = nil
	}
	if prof == nil {
		prof = FxAddProfileGlsl(def)
	}
	return
}

//	Creates and returns a new assets.FxColor initialized with the specified
//	Sid and red, green, blue and alpha channel values.
func NewFxColor(sid string, r, g, b, a float32) (me *nga.FxColor) {
	me = &nga.FxColor{}
	me.Sid, me.R, me.G, me.B, me.A = sid, r, g, b, a
	return me
}

//	Creates and returns a new assets.FxColorOrTexture referring to either the
//	specified texture or the specified color. Only at least and at most one of
//	the arguments (texture, color, paramRef) must be non-nil/non-empty.
func NewFxColorOrTexture(texture *nga.FxTexture, color *nga.FxColor, paramRef string) (ct *nga.FxColorOrTexture) {
	ct = &nga.FxColorOrTexture{}
	ct.Texture = texture
	ct.Color = color
	ct.ParamRef.SetParamRef(paramRef)
	return
}

//	Creates and returns a new assets.FxTexture sampling from the specified 2D sampler.
func NewFxTexture(sampler2DParamRef, texCoord string) (me *nga.FxTexture) {
	me = &nga.FxTexture{}
	me.Sampler2D.SetParamRef(sampler2DParamRef)
	me.TexCoord = texCoord
	return me
}
