package collimp

import (
	"encoding/hex"

	nga "github.com/go3d/go-ngine/assets"
	c15 "github.com/metaleap/go-xsd-pkg/khronos.org/files/collada_schema_1_5_go"
)

func c15_TassetType(obj *c15.TassetType) {
}

func c15_TcameraType(obj *c15.TcameraType) {
	def := nga.CameraDefs.AddNew(obj.Id.String())
	if (obj.Optics != nil) && (obj.Optics.TechniqueCommon != nil) {
		if o := obj.Optics.TechniqueCommon.Orthographic; o != nil {
			def.Ortho, def.MagX, def.MagY, def.Zfar, def.Znear = true, f64(o.Xmag), f64(o.Ymag), f64(o.Zfar), f64(o.Znear)
		} else if p := obj.Optics.TechniqueCommon.Perspective; p != nil {
			def.Ortho, def.FovX, def.FovY, def.Zfar, def.Znear = false, f64(p.Xfov), f64(p.Yfov), f64(p.Zfar), f64(p.Znear)
		}
	}
}

func c15_TimageType(obj *c15.TimageType) {
	def := nga.ImageDefs.AddNew(obj.Id.String())
	if init := obj.InitFrom; init != nil {
		if !init.MipsGenerate {
			def.InitFrom.AutoMip = false
		}
		def.InitFrom.RefUrl = init.Ref.String()
		if len(init.Hex.XsdGoPkgValue) > 0 {
			def.InitFrom.RawData, _ = hex.DecodeString(init.Hex.ToXsdtString().String())
		}
	}
}
