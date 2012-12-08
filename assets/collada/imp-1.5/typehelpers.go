package collimp

import (
	nga "github.com/go3d/go-ngine/assets"
	c15 "github.com/metaleap/go-xsd-pkg/khronos.org/files/collada_schema_1_5_go"
	xsdt "github.com/metaleap/go-xsd/types"
)

func pickBools(to ...c15.ToTlistOfBoolsType) []xsdt.Boolean {
	for _, t := range to {
		if b := t.ToTlistOfBoolsType(); len(b) > 0 {
			return b.Values()
		}
	}
	return nil
}

func pickFloats(to ...c15.ToTlistOfFloatsType) []xsdt.Double {
	for _, t := range to {
		if f := t.ToTlistOfFloatsType(); len(f) > 0 {
			return f.ValuesXsdtDouble()
		}
	}
	return nil
}

func pickInts(to ...c15.ToTlistOfIntsType) []xsdt.Long {
	for _, t := range to {
		if i := t.ToTlistOfIntsType(); len(i) > 0 {
			return i.ValuesXsdtLong()
		}
	}
	return nil
}

func pickSampler(to ...*c15.TfxSamplerType) (t *c15.TfxSamplerType) {
	for _, t = range to {
		if t != nil {
			return
		}
	}
	return
}

func f64(d xsdt.ToXsdtDouble) float64 {
	if d == nil {
		return 0
	}
	return d.ToXsdtDouble().N()
}

func f64c(d xsdt.ToXsdtDouble) float64 {
	return (f64(d) * state.curAssetUnitInMeters) / nga.UnitInMeters
}

func i64(d xsdt.ToXsdtLong) int64 {
	if d == nil {
		return 0
	}
	return d.ToXsdtLong().N()
}
