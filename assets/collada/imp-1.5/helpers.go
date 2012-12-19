package collimp

import (
	"strings"

	xmlx "github.com/jteeuwen/go-pkg-xmlx"

	nga "github.com/go3d/go-ngine/assets"
	ugo "github.com/metaleap/go-util"
	ugfx "github.com/metaleap/go-util/gfx"
	unum "github.com/metaleap/go-util/num"
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

func arr_Uints(xn *xmlx.Node, l int, s func(int, uint64)) {
	for i, n := range list_Uints(xn) {
		if i >= l {
			break
		}
		s(i, n)
	}
}

func get_CubeFace(xn *xmlx.Node) (cf int) {
	switch strings.ToUpper(xas(xn, "face")) {
	case "NEGATIVE_X":
		cf = nga.FX_CUBE_FACE_NEGATIVE_X
	case "POSITIVE_Y":
		cf = nga.FX_CUBE_FACE_POSITIVE_Y
	case "NEGATIVE_Y":
		cf = nga.FX_CUBE_FACE_NEGATIVE_Y
	case "POSITIVE_Z":
		cf = nga.FX_CUBE_FACE_POSITIVE_Z
	case "NEGATIVE_Z":
		cf = nga.FX_CUBE_FACE_NEGATIVE_Z
	default:
		cf = nga.FX_CUBE_FACE_POSITIVE_X
	}
	return
}

func get_LightColor(xn *xmlx.Node, oc *nga.Float3) {
	if c := obj_Float3(xn, "color"); c != nil {
		*oc = *c
	}
}

func get_ParamRef(xn *xmlx.Node, name string) (v string) {
	if len(name) > 0 {
		xn = xcn(xn, name)
	}
	if xn != nil {
		if v = xas(xn, "ref"); len(v) == 0 {
			v = xn.Value
		}
	}
	return
}

func get_Transforms(xn *xmlx.Node) (ts []*nga.Transform) {
	for _, cn := range xn.Children {
		if cn.Type == xmlx.NT_ELEMENT {
			if t := obj_Transform(cn, ""); (t != nil) && (t.Type > 0) {
				ts = append(ts, t)
			}
		}
	}
	return
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

func has_Id(xn *xmlx.Node, obj *nga.HasId) {
	obj.Id = xas(xn, "id")
}

func has_Inputs(xn *xmlx.Node, obj *nga.HasInputs) {
	obj.Inputs = objs_Input(xn, "input")
}

func has_Name(xn *xmlx.Node, obj *nga.HasName) {
	obj.Name = xas(xn, "name")
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
	obj.Sid = xas(xn, "sid")
}

func has_Sources(xn *xmlx.Node, obj *nga.HasSources) {
	for _, src := range objs_Source(xn, "source") {
		if src != nil {
			obj.Sources[src.Id] = src
		}
	}
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

func list_StringsN(xn *xmlx.Node, name string) (sl []string) {
	cns := xcns(xn, name)
	sl = make([]string, len(cns))
	for i, cn := range cns {
		sl[i] = cn.Value
	}
	return
}

func list_Uints(xn *xmlx.Node) (sl []uint64) {
	var v xsdt.UnsignedLong
	vals := xsdt.ListValues(xn.Value)
	sl = make([]uint64, len(vals))
	for i, s := range vals {
		v.Set(s)
		sl[i] = v.N()
	}
	return
}

func listcn_Bools(xn *xmlx.Node, name string) (sl []bool) {
	if cn := xcn(xn, name); cn != nil {
		sl = list_Bools(cn)
	}
	return
}

func listcn_Floats(xn *xmlx.Node, name string) (sl []float64) {
	if cn := xcn(xn, name); cn != nil {
		sl = list_Floats(cn)
	}
	return
}

func listcn_Ints(xn *xmlx.Node, name string) (sl []int64) {
	if cn := xcn(xn, name); cn != nil {
		sl = list_Ints(cn)
	}
	return
}

func listcn_Strings(xn *xmlx.Node, name string) (sl []string) {
	if cn := xcn(xn, name); cn != nil {
		sl = list_Strings(cn)
	}
	return
}

func listcn_Uints(xn *xmlx.Node, name string) (sl []uint64) {
	if cn := xcn(xn, name); cn != nil {
		sl = list_Uints(cn)
	}
	return
}

func node_TechCommon(xn *xmlx.Node) *xmlx.Node {
	return xcn(xn, "technique_common")
}

func setInstDefRef(xn *xmlx.Node, inst *nga.BaseInst) {
	inst.DefRef.Set(xas1(xn, "url", "body", "constraint", "target"))
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

func xaf64(xn *xmlx.Node, name string) (v float64) {
	if v = xn.F64(xmlns, name); v == 0 {
		if v = xn.F64("", name); v == 0 {
			v = xn.F64("*", name)
		}
	}
	return
}

func xaf64d(xn *xmlx.Node, name string, def float64) (v float64) {
	if v = xaf64(xn, name); v == 0 {
		v = def
	}
	return
}

func xas(xn *xmlx.Node, name string) (v string) {
	if v = xn.As(xmlns, name); len(v) == 0 {
		if v = xn.As("", name); len(v) == 0 {
			v = xn.As("*", name)
		}
	}
	return
}

func xas1(xn *xmlx.Node, names ...string) (v string) {
	for _, n := range names {
		if v = xas(xn, n); len(v) > 0 {
			return
		}
	}
	return
}

func xasd(xn *xmlx.Node, name string, def string) (v string) {
	if v = xas(xn, name); len(v) == 0 {
		v = def
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

func xbp(xn *xmlx.Node, name string) (v *bool) {
	for _, cn := range xn.Children {
		if cn.Name.Local == name {
			*v = xb(xn, name)
			return
		}
	}
	return
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

func xcns(xn *xmlx.Node, names ...string) (cns []*xmlx.Node) {
	if len(names) == 1 {
		if cns = xn.SelectNodes(xmlns, names[0]); len(cns) == 0 {
			if cns = xn.SelectNodes("", names[0]); len(cns) == 0 {
				cns = xn.SelectNodes("*", names[0])
			}
		}
	} else {
		var scns []*xmlx.Node
		for _, n := range names {
			if scns = xcns(xn, n); len(scns) > 0 {
				cns = append(cns, scns...)
			}
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

func xf64p(xn *xmlx.Node, name string) (v *float64) {
	for _, cn := range xn.Children {
		if cn.Name.Local == name {
			*v = xf64(xn, name)
			return
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

func xi64p(xn *xmlx.Node, name string) (v *int64) {
	for _, cn := range xn.Children {
		if cn.Name.Local == name {
			*v = xi64(xn, name)
			return
		}
	}
	return
}

func xm4(xn *xmlx.Node, name string) (mat *unum.Mat4) {
	mat = unum.NewMat4Identity()
	arr_Floats(xn, 16, func(i int, f float64) {
		mat[i] = f
	})
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

func xsd(xn *xmlx.Node, name string, def string) (v string) {
	if v = xs(xn, name); len(v) == 0 {
		v = def
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
			case "sampler_image":
				val = obj_FxSamplerImage(xn, "")
			case "sampler_states":
				val = obj_FxSamplerStates(xn, "")
			case "SIDREF", "sidref":
				val = nga.RefSid(xn.Value)
			default:
				val = xn.Value
			}
		}
	}
	return
}

func xv3(xn *xmlx.Node, name string) *unum.Vec3 {
	var f3 [3]float64
	arr_Floats(xn, 3, func(i int, f float64) {
		f3[i] = f
	})
	return &unum.Vec3{f3[0], f3[1], f3[2]}
}
