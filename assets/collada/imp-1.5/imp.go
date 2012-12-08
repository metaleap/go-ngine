package collimp

import (
	"encoding/xml"
	"fmt"
	"log"

	nga "github.com/go3d/go-ngine/assets"
	c141 "github.com/go3d/go-ngine/assets/collada/conv-1.4.1-to-1.5"
	ugfx "github.com/metaleap/go-util/gfx"
	c15 "github.com/metaleap/go-xsd-pkg/khronos.org/files/collada_schema_1_5_go"
	xsdt "github.com/metaleap/go-xsd/types"
)

var (
	state *importState
)

//	Provides options for importing Collada documents.
type ImportBag struct {
}

//	Initializes and returns a newly created ImportBag instance.
func NewImportBag() (me *ImportBag) {
	me = &ImportBag{}
	return
}

type importState struct {
	d15                  *c15.TxsdCollada
	curAssetUnitInMeters float64
	assetStack           []*c15.TassetType
}

//	Imports the specified Collada document, using the import options specified in importBag.
func ImportCollada(colladaDoc []byte, importBag *ImportBag) (err error) {
	state = &importState{curAssetUnitInMeters: 1, d15: &c15.TxsdCollada{}}
	c141.Force, c141.Strict = false, false
	if colladaDoc, err = c141.Convert(colladaDoc); err == nil {
		err = xml.Unmarshal(colladaDoc, state.d15)
		if colladaDoc = nil; err == nil {
			c15.WalkHandlers.TcameraType = c15_TcameraType
			c15.WalkHandlers.TeffectType = c15_TeffectType
			c15.WalkHandlers.TimageType = c15_TimageType
			c15.WalkHandlers.TmaterialType = c15_TmaterialType
			c15.WalkHandlers.TxsdCollada = c15_TxsdCollada
			c15.WalkOnError = func(err error) {
				log.Printf("%v\n", err)
			}
			c15.WalkContinueOnError = true
			state.d15.Walk()
		}
	}
	return
}

func imp_TfxNewparamType(obj *c15.TfxNewparamType) (np *nga.FxParamDef) {
	np = &nga.FxParamDef{Modifier: obj.Modifier.String(), Semantic: obj.Semantic.String()}
	np.Sid = obj.Sid.String()
	if obj.Bool != nil {
		np.Value = obj.Bool.B()
	} else if len(obj.Enum) > 0 {
		np.Value = obj.Enum.String()
	} else if obj.Float != nil {
		np.Value = f64(obj.Float)
	} else if obj.Int != nil {
		np.Value = i64(obj.Int)
	} else if sampler := pickSampler(&obj.Sampler1D.TfxSamplerType, &obj.Sampler2D.TfxSamplerType, &obj.Sampler3D.TfxSamplerType, &obj.SamplerCUBE.TfxSamplerType, &obj.SamplerDEPTH.TfxSamplerType, &obj.SamplerRECT.TfxSamplerType); sampler != nil {
		np.Value = imp_TfxSamplerType(sampler)
	} else if bools := pickBools(obj.Bool2, obj.Bool3, obj.Bool4); len(bools) > 0 {
		np.Value = xsdt.ListValuesBoolean(bools)
	} else if floats := pickFloats(obj.Float2, obj.Float3, obj.Float4, obj.Float2X1, obj.Float2X2, obj.Float2X3, obj.Float2X4, obj.Float3X1, obj.Float3X2, obj.Float3X3, obj.Float3X4, obj.Float4X1, obj.Float4X2, obj.Float4X3, obj.Float4X4); len(floats) > 0 {
		np.Value = xsdt.ListValuesDouble(floats)
	} else if ints := pickInts(obj.Int2, obj.Int3, obj.Int4); len(ints) > 0 {
		np.Value = xsdt.ListValuesLong(ints)
	}
	return
}

func imp_TfxSamplerType(obj *c15.TfxSamplerType) (s *nga.Sampler) {
	isFilterMatch := func(o fmt.Stringer, i int) bool {
		return ((i == nga.SAMPLER_FILTER_ANISOTROPIC) && (o.String() == "ANISOTROPIC")) || ((i == nga.SAMPLER_FILTER_LINEAR) && (o.String() == "LINEAR")) || ((i == nga.SAMPLER_FILTER_NEAREST) && (o.String() == "NEAREST")) || ((i == nga.SAMPLER_FILTER_NONE) && (o.String() == "NONE"))
	}
	isWrapMatch := func(o c15.TfxSamplerWrapEnum, i int) bool {
		return (o.IsBorder() && (i == nga.SAMPLER_WRAP_BORDER)) || (o.IsClamp() && (i == nga.SAMPLER_WRAP_CLAMP)) || (o.IsMirror() && (i == nga.SAMPLER_WRAP_MIRROR)) || (o.IsMirrorOnce() && (i == nga.SAMPLER_WRAP_MIRROR_ONCE)) || (o.IsWrap() && (i == nga.SAMPLER_WRAP_WRAP))
	}
	s = nga.NewSampler()
	nuWrap, nuFilter := !(isWrapMatch(obj.WrapS, s.Wrapping.WrapS) && isWrapMatch(obj.WrapT, s.Wrapping.WrapT) && isWrapMatch(obj.WrapP, s.Wrapping.WrapP)), !(isFilterMatch(obj.Magfilter, s.Filtering.Mag) && isFilterMatch(obj.Minfilter, s.Filtering.Min) && isFilterMatch(obj.Mipfilter, s.Filtering.Mip) && (obj.MaxAnisotropy.N() == s.Filtering.MaxAnisotropy))
	if len(obj.BorderColor) > 0 {

	}
	return
}
