package collimp

import (
	"strings"

	xmlx "github.com/jteeuwen/go-pkg-xmlx"

	nga "github.com/go3d/go-ngine/assets"
	ugo "github.com/metaleap/go-util"
	ugfx "github.com/metaleap/go-util/gfx"
	xsdt "github.com/metaleap/go-xsd/types"
)

const (
	xmlns = "http://www.collada.org/2008/03/COLLADASchema"
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

func has_FxParamDefs(xn *xmlx.Node, obj *nga.HasFxParamDefs) {
	for _, np := range objs_FxParamDef(xn, "newparam") {
		if np != nil {
			obj.NewParams[np.Sid] = np
		}
	}
}

func has_ID(xn *xmlx.Node, obj *nga.HasID) {
	obj.ID = xa(xn, "id")
}

func has_Name(xn *xmlx.Node, obj *nga.HasName) {
	obj.Name = xa(xn, "name")
}

func has_ParamDefs(xn *xmlx.Node, obj *nga.HasParamDefs) {
	for _, np := range objs_ParamDef(xn, "newparam") {
		if np != nil {
			obj.NewParams[np.Sid] = np
		}
	}
}

func has_ParamInsts(xn *xmlx.Node, obj *nga.HasParamInsts) {
	obj.SetParams = objs_ParamInst(xn, "setparam")
}

func has_Sid(xn *xmlx.Node, obj *nga.HasSid) {
	obj.Sid = xa(xn, "sid")
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
	return xcn(xn, "technique_common")
}

func setInstDefRef(xn *xmlx.Node, inst *nga.BaseInst) {
	inst.DefRef = xa1(xn, "url", "target", "body", "constraint")
}

func xa(xn *xmlx.Node, name string) (v string) {
	if v = xn.As(xmlns, name); len(v) == 0 {
		if v = xn.As("", name); len(v) == 0 {
			v = xn.As("*", name)
		}
	}
	return
}

func xa1(xn *xmlx.Node, names ...string) (v string) {
	for _, n := range names {
		if v = xa(xn, n); len(v) > 0 {
			return
		}
	}
	return
}

func xab(xn *xmlx.Node, name string) bool {
	return xn.Ab(xmlns, name) || xn.Ab("", name) || xn.Ab("*", name)
}

func xabp(xn *xmlx.Node, name string) (b *bool) {
	for _, att := range xn.Attributes {
		if att.Name.Local == name {
			*b = xab(xn, name)
			break
		}
	}
	return
}

func xau64(xn *xmlx.Node, name string) (v uint64) {
	if v = xn.Au64(xmlns, name); v == 0 {
		if v = xn.Au64("", name); v == 0 {
			v = xn.Au64("*", name)
		}
	}
	return
}

func xau64p(xn *xmlx.Node, name string) (p *uint64) {
	for _, att := range xn.Attributes {
		if att.Name.Local == name {
			*p = xau64(xn, name)
			break
		}
	}
	return
}

func xb(xn *xmlx.Node, name string) bool {
	return xn.B(xmlns, name) || xn.B("", name) || xn.B("*", name)
}

func xcn(xn *xmlx.Node, name string) (cn *xmlx.Node) {
	if cn = xn.SelectNode(xmlns, name); cn == nil {
		if cn = xn.SelectNode("", name); cn == nil {
			cn = xn.SelectNode("*", name)
		}
	}
	return
}

func xcn1(xn *xmlx.Node, names ...string) (cn *xmlx.Node) {
	for _, n := range names {
		if cn = xcn(xn, n); cn != nil {
			return
		}
	}
	return
}

func xcns(xn *xmlx.Node, name string) (cns []*xmlx.Node) {
	if cns = xn.SelectNodes(xmlns, name); len(cns) == 0 {
		if cns = xn.SelectNodes("", name); len(cns) == 0 {
			cns = xn.SelectNodes("*", name)
		}
	}
	return
}

func xf64(xn *xmlx.Node, name string) (v float64) {
	if v = xn.F64(xmlns, name); v == 0 {
		if v = xn.F64("", name); v == 0 {
			v = xn.F64("*", name)
		}
	}
	return
}

func xi64(xn *xmlx.Node, name string) (v int64) {
	if v = xn.I64(xmlns, name); v == 0 {
		if v = xn.I64("", name); v == 0 {
			v = xn.I64("*", name)
		}
	}
	return
}

func xs(xn *xmlx.Node, name string) (v string) {
	if v = xn.S(xmlns, name); len(v) == 0 {
		if v = xn.S("", name); len(v) == 0 {
			v = xn.S("*", name)
		}
	}
	return
}

func xss(xn *xmlx.Node, names ...string) (s string) {
	for _, n := range names {
		if s = xs(xn, n); len(s) > 0 {
			return
		}
	}
	return
}

func xu64(xn *xmlx.Node, name string) (v uint64) {
	if v = xn.U64(xmlns, name); v == 0 {
		if v = xn.U64("", name); v == 0 {
			v = xn.U64("*", name)
		}
	}
	return
}

func xu64p(xn *xmlx.Node, name string) (u *uint64) {
	if cn := xcn(xn, name); cn != nil {
		*u = xu64(cn, "")
	}
	return
}

func xv(xn *xmlx.Node) (val interface{}) {
	if xn != nil {
		if strings.HasPrefix(xn.Name.Local, "sampler") {
			val = obj_FxSampler(xn, "")
		} else {
			switch xn.Name.Local {
			case "array":
				if l := xau64(xn, "length"); (l > 0) && (len(xn.Children) > 0) {
					sl := make([]interface{}, l)
					for i := 0; i < int(l); i++ {
						sl[i] = xv(xn.Children[ugo.Ifi(i >= len(xn.Children), 0, i)])
					}
					val = sl
				}
			case "bool":
				val = xb(xn, "")
			case "bool2":
				if v := obj_Bool2(xn, ""); v != nil {
					val = *v
				}
			case "bool3":
				if v := obj_Bool3(xn, ""); v != nil {
					val = *v
				}
			case "bool4":
				if v := obj_Bool4(xn, ""); v != nil {
					val = *v
				}
			case "float":
				val = xf64(xn, "")
			case "float2", "float2x1":
				if v := obj_Float2(xn, ""); v != nil {
					val = *v
				}
			case "float3", "float3x1":
				if v := obj_Float3(xn, ""); v != nil {
					val = *v
				}
			case "float4", "float4x1":
				if v := obj_Float4(xn, ""); v != nil {
					val = *v
				}
			case "float2x2":
				if v := obj_Float2x2(xn, ""); v != nil {
					val = *v
				}
			case "float2x3":
				if v := obj_Float2x3(xn, ""); v != nil {
					val = *v
				}
			case "float2x4":
				if v := obj_Float2x4(xn, ""); v != nil {
					val = *v
				}
			case "float3x2":
				if v := obj_Float3x2(xn, ""); v != nil {
					val = *v
				}
			case "float3x3":
				if v := obj_Float3x3(xn, ""); v != nil {
					val = *v
				}
			case "float3x4":
				if v := obj_Float3x4(xn, ""); v != nil {
					val = *v
				}
			case "float4x2":
				if v := obj_Float4x2(xn, ""); v != nil {
					val = *v
				}
			case "float4x3":
				if v := obj_Float4x3(xn, ""); v != nil {
					val = *v
				}
			case "float4x4":
				if v := obj_Float4x4(xn, ""); v != nil {
					val = *v
				}
			case "int":
				val = xi64(xn, "")
			case "int2":
				if v := obj_Int2(xn, ""); v != nil {
					val = *v
				}
			case "int3":
				if v := obj_Int3(xn, ""); v != nil {
					val = *v
				}
			case "int4":
				if v := obj_Int4(xn, ""); v != nil {
					val = *v
				}
			default:
				val = xn.Value
			}
		}
	}
	return
}
