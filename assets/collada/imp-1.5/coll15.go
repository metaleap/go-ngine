package collimp

import (
	"encoding/hex"

	nga "github.com/go3d/go-ngine/assets"
	c15 "github.com/metaleap/go-xsd-pkg/khronos.org/files/collada_schema_1_5_go"
	xsdt "github.com/metaleap/go-xsd/types"
)

func onAsset(obj *c15.TassetType, enter bool) {
	if obj != nil {
		if enter {
			nu := append([]*c15.TassetType{obj}, state.assetStack...)
			state.assetStack = nu
		} else if len(state.assetStack) > 0 {
			state.assetStack = state.assetStack[1:]
		}
		if (len(state.assetStack) > 0) && (state.assetStack[0].Unit != nil) && (state.assetStack[0].Unit.Meter != 0) {
			state.curAssetUnitInMeters = f64(state.assetStack[0].Unit.Meter)
		} else {
			state.curAssetUnitInMeters = 1
		}
	}
}

func c15_TxsdCollada(obj *c15.TxsdCollada, enter bool) (err error) {
	onAsset(obj.Asset, enter)
	return
}

func c15_TcameraType(obj *c15.TcameraType, enter bool) (err error) {
	if onAsset(obj.Asset, enter); enter {
		def := nga.CameraDefs.AddNew(obj.Id.String())
		if (obj.Optics != nil) && (obj.Optics.TechniqueCommon != nil) {
			if o := obj.Optics.TechniqueCommon.Orthographic; o != nil {
				def.Ortho, def.MagX, def.MagY, def.Zfar, def.Znear = true, f64(o.Xmag), f64(o.Ymag), f64c(o.Zfar), f64c(o.Znear)
			} else if p := obj.Optics.TechniqueCommon.Perspective; p != nil {
				def.Ortho, def.FovX, def.FovY, def.Zfar, def.Znear = false, f64(p.Xfov), f64(p.Yfov), f64c(p.Zfar), f64c(p.Znear)
			}
		}
	}
	return
}

func c15_TeffectType(obj *c15.TeffectType, enter bool) (err error) {
	if onAsset(obj.Asset, enter); enter {
		def := nga.EffectDefs.AddNew(obj.Id.String())
		for _, np := range obj.Newparams {
			def.NewParams[np.Sid.String()] = imp_TfxNewparamType(np)
		}
	}
	return
}

func c15_TimageType(obj *c15.TimageType, enter bool) (err error) {
	if onAsset(obj.Asset, enter); enter {
		def := nga.ImageDefs.AddNew(obj.Id.String())
		if init := obj.InitFrom; init != nil {
			if !init.MipsGenerate {
				def.InitFrom.AutoMip = false
			}
			def.InitFrom.RefUrl = init.Ref.String()
			if (init.Hex != nil) && (len(init.Hex.XsdGoPkgValue) > 0) {
				def.InitFrom.RawData, err = hex.DecodeString(init.Hex.ToXsdtString().String())
			}
		}
	}
	return
}

func c15_TmaterialType(obj *c15.TmaterialType, enter bool) (err error) {
	if onAsset(obj.Asset, enter); enter {

	}
	return
}

func imp_TfxNewparamType(obj *c15.TfxNewparamType) (np *nga.FxNewParam) {
	var (
		val    interface{}
		bools  []xsdt.Boolean
		floats []xsdt.Double
	)
	np = nga.NewFxNewParam(obj.Modifier.String(), obj.Semantic.String())
	for _, ann := range obj.Annotates {
		if ann.Bool != nil {
			val = bool(*ann.Bool)
		} else if ann.Float != nil {
			val = f64(ann.Float)
		} else if ann.Int != nil {
			val = i64(ann.Int)
		} else if len(ann.Bool2) > 0 {
			bools = ann.Bool2.ToTlistOfBoolsType().Values()
		} else if len(ann.Bool3) > 0 {
			bools = ann.Bool3.ToTlistOfBoolsType().Values()
		} else if len(ann.Bool4) > 0 {
			bools = ann.Bool4.ToTlistOfBoolsType().Values()
		} else if len(ann.Float2) > 0 {
			floats = ann.Float2.ToTlistOfFloatsType().Values()
		}
		if len(bools) > 0 {
			val = xsdt.ListValuesBoolean(bools)
		} else if len(floats) > 0 {
			val = xsdt.ListValuesDouble(floats)
		}
		np.Annotations[ann.Name.String()] = val
	}
	return
}
