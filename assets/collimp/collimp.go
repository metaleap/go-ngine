package collimp

import (
	coll14 "github.com/metaleap/go-xsd-pkg/khronos.org/files/collada_schema_1_4_go"
	coll15 "github.com/metaleap/go-xsd-pkg/khronos.org/files/collada_schema_1_5_go"

	nga "github.com/go3d/go-ngine/assets"
)

type c14AssetStack []*coll14.TxsdAsset

	func (me *c14AssetStack) Pop (popAsset bool) (p bool, curAsset *coll14.TxsdAsset) {
		if popAsset && (len(*me) > 0) { sl := *me; *me = sl[1 :] }; if len(*me) > 0 { curAsset = (*me)[0] }; return
	}

	func (me *c14AssetStack) Push (el *coll14.TxsdAsset) (p bool, curAsset *coll14.TxsdAsset) {
		if p = (el != nil); p { curAsset = el; nu := []*coll14.TxsdAsset { el }; *me = append(nu, *me...) } else if len(*me) > 0 { curAsset = (*me)[0] }; return
	}

func ImportCollada (c14 *coll14.TxsdCollada, c15 *coll15.TxsdCollada) {
	var f64 = func (f *coll14.TargetableFloat) float64 {
		return float64(f.XsdGoPkgValue.ToXsdtDouble())
	}
	if c14 != nil {
		coll14.WalkHandlers.TxsdAsset = func (x *coll14.TxsdAsset) {
		}
		coll14.WalkHandlers.TxsdCamera = func (cam *coll14.TxsdCamera) {
			var def = nga.CameraDefs.AddNew(cam.Id.String())
			if (cam.Optics != nil) && (cam.Optics.TechniqueCommon != nil) {
				if o := cam.Optics.TechniqueCommon.Orthographic; o != nil {
					def.IsOrtho, def.FovOrMagX, def.FovOrMagY, def.ZfarPlane, def.ZnearPlane = true, f64(o.Xmag), f64(o.Ymag), f64(o.Zfar), f64(o.Znear)
				} else if p := cam.Optics.TechniqueCommon.Perspective; p != nil {
					def.IsOrtho, def.FovOrMagX, def.FovOrMagY, def.ZfarPlane, def.ZnearPlane = false, f64(p.Xfov), f64(p.Yfov), f64(p.Zfar), f64(p.Znear)
				}
			}
		}
		c14.Walk()
	}
	if c15 != nil {

	}
}
