package collimp

import (
	coll14 "github.com/metaleap/go-xsd-pkg/khronos.org/files/collada_schema_1_4_go"
	nga "github.com/go3d/go-ngine/assets"
)

func Imp14Camera (cam *coll14.TxsdCamera) {
	var def = nga.CameraDefs.AddNew(cam.Id.String())
	if (cam.Optics != nil) && (cam.Optics.TechniqueCommon != nil) {
		if o := cam.Optics.TechniqueCommon.Orthographic; o != nil {
			def.Ortho, def.MagX, def.MagY, def.Zfar, def.Znear = true, f64(o.Xmag), f64(o.Ymag), f64(o.Zfar), f64(o.Znear)
		} else if p := cam.Optics.TechniqueCommon.Perspective; p != nil {
			def.Ortho, def.FovX, def.FovY, def.Zfar, def.Znear = false, f64(p.Xfov), f64(p.Yfov), f64(p.Zfar), f64(p.Znear)
		}
	}
}
