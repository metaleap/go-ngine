package collimp

import (
	xmlx "github.com/jteeuwen/go-pkg-xmlx"

	nga "github.com/go3d/go-ngine/assets"
	ugfx "github.com/metaleap/go-util/gfx"
	xsdt "github.com/metaleap/go-xsd/types"
)

func arr_Bools(xn *xmlx.Node, l int, s func(int, bool)) {
	for i, b := range list_Bools(xn) {
		if i >= l {
			break
		}
		s(i, b)
	}
}

func arr_Floats(xn *xmlx.Node, l int, s func(int, float64)) {
	for i, f := range list_Floats(xn) {
		if i >= l {
			break
		}
		s(i, f)
	}
}

func arr_Ints(xn *xmlx.Node, l int, s func(int, int64)) {
	for i, n := range list_Ints(xn) {
		if i >= l {
			break
		}
		s(i, n)
	}
}

func has_Asset(xn *xmlx.Node, obj *nga.HasAsset) {
	obj.Asset = obj_Asset(xn, "asset")
}

func has_Extras(xn *xmlx.Node, obj *nga.HasExtras) {
	obj.Extras = objs_Extra(xn, "extra")
}

func has_ID(xn *xmlx.Node, obj *nga.HasID) {
	obj.ID = xna(xn, "id")
}

func has_Name(xn *xmlx.Node, obj *nga.HasName) {
	obj.Name = xna(xn, "name")
}

func has_Sid(xn *xmlx.Node, obj *nga.HasSid) {
	obj.Sid = xna(xn, "sid")
}

func has_Techniques(xn *xmlx.Node, obj *nga.HasTechniques) {
	obj.Techniques = objs_Technique(xn, "technique")
}

func list_Bools(xn *xmlx.Node) (sl []bool) {
	var v xsdt.Boolean
	vals := xsdt.ListValues(xn.Value)
	sl = make([]bool, len(vals))
	for i, s := range vals {
		v.Set(s)
		sl[i] = v.B()
	}
	return
}

func list_Floats(xn *xmlx.Node) (sl []float64) {
	var v xsdt.Double
	vals := xsdt.ListValues(xn.Value)
	sl = make([]float64, len(vals))
	for i, s := range vals {
		v.Set(s)
		sl[i] = v.N()
	}
	return
}

func list_Ints(xn *xmlx.Node) (sl []int64) {
	var v xsdt.Long
	vals := xsdt.ListValues(xn.Value)
	sl = make([]int64, len(vals))
	for i, s := range vals {
		v.Set(s)
		sl[i] = v.N()
	}
	return
}

func list_Rgba32(xn *xmlx.Node, obj *ugfx.Rgba32) {
	if f := list_Floats(xn); len(f) > 0 {
		if obj.R = float32(f[0]); len(f) > 1 {
			if obj.G = float32(f[1]); len(f) > 2 {
				if obj.B = float32(f[2]); len(f) > 3 {
					obj.A = float32(f[3])
				}
			}
		}
	}
}

func list_Strings(xn *xmlx.Node) []string {
	return xsdt.ListValues(xn.Value)
}

func node_TechCommon(xn *xmlx.Node) *xmlx.Node {
	return xsn(xn, "technique_common")
}

func setInstDefRef(xn *xmlx.Node, inst *nga.BaseInst) {
	inst.DefRef = xna1(xn, "url", "target", "body", "constraint")
}

func xna(xn *xmlx.Node, name string) string {
	return xn.As("", name)
}

func xna1(xn *xmlx.Node, names ...string) (v string) {
	for _, n := range names {
		if v = xna(xn, n); len(v) > 0 {
			return
		}
	}
	return
}

func xnau64(xn *xmlx.Node, name string) uint64 {
	return xn.Au64("", name)
}

func xnau64p(xn *xmlx.Node, name string) (p *uint64) {
	for _, att := range xn.Attributes {
		if att.Name.Local == name {
			*p = xnau64(xn, name)
			break
		}
	}
	return
}

func xnf64(xn *xmlx.Node, name string) float64 {
	return xn.F64("", name)
}

func xsn(xn *xmlx.Node, name string) (sn *xmlx.Node) {
	sn = xn.SelectNode("", name)
	return
}

func xsn1(xn *xmlx.Node, names ...string) (sn *xmlx.Node) {
	for _, n := range names {
		if sn = xsn(xn, n); sn != nil {
			return
		}
	}
	return
}

func xsns(xn *xmlx.Node, name string) []*xmlx.Node {
	return xn.SelectNodes("", name)
}
