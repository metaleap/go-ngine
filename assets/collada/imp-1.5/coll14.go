package collimp

import (
	"encoding/hex"

	nga "github.com/go3d/go-ngine/assets"
	c14 "github.com/metaleap/go-xsd-pkg/khronos.org/files/collada_schema_1_4_go"
)

func c14_TxsdAsset(obj *c14.TxsdAsset) {
}

func c14_TxsdCamera(obj *c14.TxsdCamera) {
	var def = nga.CameraDefs.AddNew(obj.Id.String())
	if (obj.Optics != nil) && (obj.Optics.TechniqueCommon != nil) {
		if o := obj.Optics.TechniqueCommon.Orthographic; o != nil {
			def.Ortho, def.MagX, def.MagY, def.Zfar, def.Znear = true, f64(o.Xmag), f64(o.Ymag), f64(o.Zfar), f64(o.Znear)
		} else if p := obj.Optics.TechniqueCommon.Perspective; p != nil {
			def.Ortho, def.FovX, def.FovY, def.Zfar, def.Znear = false, f64(p.Xfov), f64(p.Yfov), f64(p.Zfar), f64(p.Znear)
		}
	}
}

func c14_TxsdImage(obj *c14.TxsdImage) {
	var def = nga.ImageDefs.AddNew(obj.Id.String())
	if len(obj.Data) > 0 {
		def.InitFrom.RawData, _ = hex.DecodeString(obj.Data.ToXsdtString().String())
	}
	if len(obj.InitFrom) > 0 {
		def.InitFrom.RefUrl = obj.InitFrom.String()
	}
}
