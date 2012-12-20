package collimp

import (
	"encoding/hex"
	"strings"

	xmlx "github.com/jteeuwen/go-pkg-xmlx"

	nga "github.com/go3d/go-ngine/assets"
	ustr "github.com/metaleap/go-util/str"
	xsdt "github.com/metaleap/go-xsd/types"
)

func load_Document(xn *xmlx.Node, obj *nga.Document) {
	obj.Scene = obj_Scene(xn, "scene")
}

func load_FxEffectDef(xn *xmlx.Node, obj *nga.FxEffectDef) {
	obj.Annotations = objs_FxAnnotation(xn, "annotate")
	for _, cn := range xn.Children {
		if (cn.Type == xmlx.NT_ELEMENT) && strings.HasPrefix(cn.Name.Local, "profile_") {
			obj.Profiles = append(obj.Profiles, obj_FxProfile(cn, ""))
		}
	}
}

func load_PxCylinder(xn *xmlx.Node, obj *nga.PxCylinder) {
	if bc := obj_GeometryBrepCylinder(xn, ""); bc != nil {
		obj.GeometryBrepCylinder = *bc
	}
	obj.Height = xf64(xn, "height")
}

func load_Float4x4(xn *xmlx.Node, obj *nga.Float4x4) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_FxTechniqueCommonPhong(xn *xmlx.Node, obj *nga.FxTechniqueCommonPhong) {
	if b := obj_FxTechniqueCommonBlinn(xn, ""); b != nil {
		obj.FxTechniqueCommonBlinn = *b
	}
}

func load_LightDef(xn *xmlx.Node, obj *nga.LightDef) {
	if tcn := node_TechCommon(xn); tcn != nil {
		obj.TC.Ambient = obj_LightAmbient(tcn, "ambient")
		obj.TC.Directional = obj_LightDirectional(tcn, "directional")
		obj.TC.Point = obj_LightPoint(tcn, "point")
		obj.TC.Spot = obj_LightSpot(tcn, "spot")
	}
}

func load_SourceArray(xn *xmlx.Node, obj *nga.SourceArray) {
	if dn := xcn1(xn, "bool_array", "float_array", "IDREF_array", "idref_array", "int_array", "Name_array", "name_array", "SIDREF_array", "sidref_array", "token_array"); dn != nil {
		switch strings.ToLower(dn.Name.Local) {
		case "bool_array":
			obj.Bools = list_Bools(dn)
		case "float_array":
			obj.Floats = list_Floats(dn)
		case "idref_array":
			obj.IdRefs = list_Strings(dn)
		case "int_array":
			obj.Ints = list_Ints(dn)
		case "name_array":
			obj.Names = list_Strings(dn)
		case "sidref_array":
			obj.SidRefs = list_Strings(dn)
		case "token_array":
			obj.Tokens = list_Strings(dn)
		}
	}
}

func load_Source(xn *xmlx.Node, obj *nga.Source) {
	if tcn := node_TechCommon(xn); tcn != nil {
		obj.TC.Accessor = obj_SourceAccessor(tcn, "accessor")
	}
	if arr := obj_SourceArray(xn, ""); arr != nil {
		obj.Array = *arr
	}
}

func load_FxPass(xn *xmlx.Node, obj *nga.FxPass) {
	obj.Annotations = objs_FxAnnotation(xn, "annotate")
	obj.Evaluate = obj_FxPassEvaluation(xn, "evaluate")
	obj.Program = obj_FxPassProgram(xn, "program")
	if sn := xcn(xn, "states"); sn != nil {
		for _, scn := range sn.Children {
			if scn.Type == xmlx.NT_ELEMENT {
				obj.States[scn.Name.Local] = obj_FxPassState(scn, "")
			}
		}
	}
}

func load_FxPassEvaluationClearColor(xn *xmlx.Node, obj *nga.FxPassEvaluationClearColor) {
	obj.Index = xau64(xn, "index")
	list_Rgba32(xn, &obj.Rgba32)
}

func load_GeometryBrepCylinder(xn *xmlx.Node, obj *nga.GeometryBrepCylinder) {
	if f := obj_Float2(xn, "radius"); f != nil {
		obj.Radii = *f
	}
}

func load_PxModelDef(xn *xmlx.Node, obj *nga.PxModelDef) {
	obj.Insts = objs_PxModelInst(xn, "instance_physics_model")
	obj.RigidBodies = objs_PxRigidBodyDef(xn, "rigid_body")
	obj.RigidConstraints = objs_PxRigidConstraintDef(xn, "rigid_constraint")
}

func load_FxCreateCubeInitFrom(xn *xmlx.Node, obj *nga.FxCreateCubeInitFrom) {
	if i := obj_FxCreateInitFrom(xn, ""); i != nil {
		obj.FxCreateInitFrom = *i
	}
	obj.Face = get_CubeFace(xn)
}

func load_KxModelInst(xn *xmlx.Node, obj *nga.KxModelInst) {
	obj.Bindings = objs_KxBinding(xn, "bind")
}

func load_GeometryBrepSphere(xn *xmlx.Node, obj *nga.GeometryBrepSphere) {
	obj.Radius = xf64(xn, "radius")
}

func load_Float4(xn *xmlx.Node, obj *nga.Float4) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_SourceAccessor(xn *xmlx.Node, obj *nga.SourceAccessor) {
	obj.Count = xau64(xn, "count")
	obj.Offset = xau64(xn, "offset")
	obj.Source.Set(xas(xn, "source"))
	if u := xau64p(xn, "stride"); u != nil {
		obj.Stride = *u
	} else {
		obj.Stride = 1
	}
	obj.Params = objs_Param(xn, "param")
}

func load_FxPassEvaluation(xn *xmlx.Node, obj *nga.FxPassEvaluation) {
	obj.Draw = xs(xn, "draw")
	obj.Color.Clear = obj_FxPassEvaluationClearColor(xn, "color_clear")
	obj.Color.Target = obj_FxPassEvaluationTarget(xn, "color_target")
	obj.Depth.Clear = obj_FxPassEvaluationClearDepth(xn, "depth_clear")
	obj.Depth.Target = obj_FxPassEvaluationTarget(xn, "depth_target")
	obj.Stencil.Clear = obj_FxPassEvaluationClearStencil(xn, "stencil_clear")
	obj.Stencil.Target = obj_FxPassEvaluationTarget(xn, "stencil_target")
}

func load_FxEffectInstTechniqueHint(xn *xmlx.Node, obj *nga.FxEffectInstTechniqueHint) {
	obj.Platform = xas(xn, "platform")
	obj.Ref = xas(xn, "ref")
	obj.Profile = xas(xn, "profile")
}

func load_FxSamplerFiltering(xn *xmlx.Node, obj *nga.FxSamplerFiltering) {
	obj.MipBias = xf64(xn, "mip_bias")
	obj.MipMinLevel = uint8(xu64(xn, "mip_min_level"))
	obj.MipMaxLevel = uint8(xu64(xn, "mip_max_level"))
	if u := xu64p(xn, "max_anisotropy"); u != nil {
		obj.MaxAnisotropy = uint32(*u)
	} else {
		obj.MaxAnisotropy = 1
	}
	for n, i := range map[string]*int{"minfilter": &obj.FilterMin, "magfilter": &obj.FilterMag, "mipfilter": &obj.FilterMip} {
		switch strings.ToUpper(xs(xn, n)) {
		case "ANISOTROPIC":
			*i = nga.FX_SAMPLER_FILTER_ANISOTROPIC
		case "NEAREST":
			*i = nga.FX_SAMPLER_FILTER_NEAREST
		case "NONE":
			*i = nga.FX_SAMPLER_FILTER_NONE
		default:
			*i = nga.FX_SAMPLER_FILTER_LINEAR
		}
	}
}

func load_FxPassProgram(xn *xmlx.Node, obj *nga.FxPassProgram) {
	obj.Shaders = objs_FxPassProgramShader(xn, "shader")
	obj.BindAttributes = objs_FxPassProgramBindAttribute(xn, "bind_attribute")
	obj.BindUniforms = objs_FxPassProgramBindUniform(xn, "bind_uniform")
}

func load_CameraInst(xn *xmlx.Node, obj *nga.CameraInst) {
}

func load_PxMaterialInst(xn *xmlx.Node, obj *nga.PxMaterialInst) {
}

func load_FxImageDef(xn *xmlx.Node, obj *nga.FxImageDef) {
	obj.Create2D = obj_FxCreate2D(xn, "create_2d")
	obj.Create3D = obj_FxCreate3D(xn, "create_3d")
	obj.CreateCube = obj_FxCreateCube(xn, "create_cube")
	obj.InitFrom = obj_FxImageInitFrom(xn, "init_from")
	if rn := xcn(xn, "renderable"); rn != nil {
		obj.Renderable.Is = true
		obj.Renderable.Shared = xab(rn, "share")
	}
}

func load_Float4x3(xn *xmlx.Node, obj *nga.Float4x3) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_KxKinematicsAxis(xn *xmlx.Node, obj *nga.KxKinematicsAxis) {
	var f *nga.Formula
	obj.Axis.Set(xas(xn, "axis"))
	for _, cn := range xcns(xn, "formula", "instance_formula") {
		if f = obj_Formula(cn, ""); f != nil {
			obj.Formulas = append(obj.Formulas, *f)
		}
	}
	obj.Limits = obj_KxAxisLimits(xn, "limits")
	obj.Indices = objs_KxAxisIndex(xn, "index")
	if l := obj_ParamOrBool(xn, "locked"); l != nil {
		obj.Locked = *l
	}
	if pb := obj_ParamOrBool(xn, "active"); pb != nil {
		obj.Active = *pb
	} else {
		obj.Active.B = true
	}
}

func load_PxModelInst(xn *xmlx.Node, obj *nga.PxModelInst) {
	obj.Parent.Set(xas(xn, "parent"))
	obj.ForceFields = objs_PxForceFieldInst(xn, "instance_force_field")
	obj.RigidBodies = objs_PxRigidBodyInst(xn, "instance_rigid_body")
	obj.RigidConstraints = objs_PxRigidConstraintInst(xn, "instance_rigid_constraint")
}

func load_Transform(xn *xmlx.Node, obj *nga.Transform) {
	switch xn.Name.Local {
	case "lookat":
		obj.Type = nga.TRANSFORM_TYPE_LOOKAT
	case "matrix":
		obj.Type = nga.TRANSFORM_TYPE_MATRIX
	case "rotate":
		obj.Type = nga.TRANSFORM_TYPE_ROTATE
	case "scale":
		obj.Type = nga.TRANSFORM_TYPE_SCALE
	case "skew":
		obj.Type = nga.TRANSFORM_TYPE_SKEW
	case "translate":
		obj.Type = nga.TRANSFORM_TYPE_TRANSLATE
	}
	if obj.Type > 0 {
		obj.F = list_Floats(xn)
	}
}

func load_KxJoint(xn *xmlx.Node, obj *nga.KxJoint) {
	switch xn.Name.Local {
	case "revolute":
		obj.Type = nga.KX_JOINT_TYPE_REVOLUTE
	case "prismatic":
		obj.Type = nga.KX_JOINT_TYPE_PRISMATIC
	}
	if obj.Type > 0 {
		if xa := xcn(xn, "axis"); xa != nil {
			if sv3 := obj_SidVec3(xa, ""); sv3 != nil {
				obj.Axis.SidVec3 = *sv3
			}
			has_Name(xa, &obj.Axis.HasName)
		}
		obj.Limits = obj_KxJointLimits(xn, "limits")
	}
}

func load_ParamOrInt(xn *xmlx.Node, obj *nga.ParamOrInt) {
	obj.Param.Set(xs(xn, "param"))
	obj.I = xi64(xn, "int")
}

func load_ParamDef(xn *xmlx.Node, obj *nga.ParamDef) {
	for _, cn := range xn.Children {
		if (cn.Type == xmlx.NT_ELEMENT) && !ustr.IsOneOf(cn.Name.Local, "annotate", "semantic", "modifier") {
			obj.Value = xv(cn)
			break
		}
	}
}

func load_Int2x2(xn *xmlx.Node, obj *nga.Int2x2) {
	arr_Ints(xn, len(obj), func(i int, n int64) {
		obj[i] = n
	})
}

func load_ControllerInst(xn *xmlx.Node, obj *nga.ControllerInst) {
	obj.BindMaterial = obj_MaterialBinding(xn, "bind_material")
	obj.SkinSkeletons = list_StringsN(xn, "skeleton")
}

func load_GeometryPositioning(xn *xmlx.Node, obj *nga.GeometryPositioning) {
	obj.Orientations = objs_GeometryBrepOrientation(xn, "orient")
	obj.Origin = xv3(xn, "origin")
}

func load_GeometryBrepCurve(xn *xmlx.Node, obj *nga.GeometryBrepCurve) {
	if pos := obj_GeometryPositioning(xn, ""); pos != nil {
		obj.Location = *pos
	}
	obj.Element.Line = obj_GeometryBrepLine(xn, "line")
	obj.Element.Circle = obj_GeometryBrepCircle(xn, "circle")
	obj.Element.Ellipse = obj_GeometryBrepEllipse(xn, "ellipse")
	obj.Element.Parabola = obj_GeometryBrepParabola(xn, "parabola")
	obj.Element.Hyperbola = obj_GeometryBrepHyperbola(xn, "hyperbola")
	obj.Element.Nurbs = obj_GeometryBrepNurbs(xn, "nurbs")
}

func load_GeometryBrepNurbs(xn *xmlx.Node, obj *nga.GeometryBrepNurbs) {
	obj.Degree = xau64(xn, "degree")
	obj.Closed = xab(xn, "closed")
	if cv := obj_GeometryControlVertices(xn, "control_vertices"); cv != nil {
		obj.ControlVertices = *cv
	}
}

func load_Formula(xn *xmlx.Node, obj *nga.Formula) {
	switch xn.Name.Local {
	case "instance_formula":
		obj.Inst = obj_FormulaInst(xn, "")
	case "formula":
		obj.Def = obj_FormulaDef(xn, "")
	}
}

func load_FormulaDef(xn *xmlx.Node, obj *nga.FormulaDef) {
	if pf := obj_ParamOrFloat(xn, "target"); pf != nil {
		obj.Target = *pf
	}
	if tcn := node_TechCommon(xn); tcn != nil {
		obj.TC.MathML = tcn.Children
	}
}

func load_LightInst(xn *xmlx.Node, obj *nga.LightInst) {
}

func load_Asset(xn *xmlx.Node, obj *nga.Asset) {
	obj.Contributors = objs_AssetContributor(xn, "contributor")
	if cn := xcn(xn, "coverage"); cn != nil {
		obj.Coverage = obj_AssetGeographicLocation(cn, "geographic_location")
	}
	if cn := xcn(xn, "unit"); cn != nil {
		obj.Unit.Meter, obj.Unit.Name = xaf64d(cn, "meter", obj.Unit.Meter), xasd(cn, "name", obj.Unit.Name)
	}
	// obj.Created = xs(xn, "created")
	// obj.Keywords = xs(xn, "keywords")
	// obj.Modified = xs(xn, "modified")
	// obj.Revision = xs(xn, "revision")
	// obj.Subject = xs(xn, "subject")
	// obj.Title = xs(xn, "title")
	if obj.UpAxis = xs(xn, "up_axis"); len(obj.UpAxis) == 0 {
		obj.UpAxis = "Y"
	} else if obj.UpAxis = strings.ToUpper(obj.UpAxis[:1]); (obj.UpAxis != "X") && (obj.UpAxis != "Z") {
		obj.UpAxis = "Y"
	}
}

func load_FxPassProgramShader(xn *xmlx.Node, obj *nga.FxPassProgramShader) {
	switch strings.ToUpper(xas(xn, "stage")) {
	case "COMPUTE":
		obj.Stage = nga.FX_PASS_PROGRAM_SHADER_STAGE_COMPUTE
	case "FRAGMENT":
		obj.Stage = nga.FX_PASS_PROGRAM_SHADER_STAGE_FRAGMENT
	case "GEOMETRY":
		obj.Stage = nga.FX_PASS_PROGRAM_SHADER_STAGE_GEOMETRY
	case "TESSELATION", "TESSELLATION":
		obj.Stage = nga.FX_PASS_PROGRAM_SHADER_STAGE_TESSELLATION
	case "VERTEX":
		obj.Stage = nga.FX_PASS_PROGRAM_SHADER_STAGE_VERTEX
	}
	if sn := xcn(xn, "sources"); sn != nil {
		pss := nga.FxPassProgramShaderSources{}
		arr := make([]nga.FxPassProgramShaderSources, 0, len(sn.Children))
		for _, scn := range sn.Children {
			if pss.S, pss.IsImportRef = scn.Value, (scn.Name.Local == "import"); pss.IsImportRef {
				pss.S = xas(scn, "ref")
			}
			arr = append(arr, pss)
		}
		obj.Sources = arr
	}
}

func load_GeometryBrepCapsule(xn *xmlx.Node, obj *nga.GeometryBrepCapsule) {
	obj.Height = xf64(xn, "height")
	if v3 := xv3(xn, "radius"); v3 != nil {
		obj.Radii = *v3
	}
}

func load_FxProfileCommon(xn *xmlx.Node, obj *nga.FxProfileCommon) {
	if tc := obj_FxTechniqueCommon(xn, "technique"); tc != nil {
		obj.Technique = *tc
	}
}

func load_ControllerMorph(xn *xmlx.Node, obj *nga.ControllerMorph) {
	obj.Source.Set(xas(xn, "source"))
	obj.Relative = (strings.ToUpper(xas(xn, "method")) == "RELATIVE")
	if t := obj_ControllerInputs(xn, "targets"); t != nil {
		obj.Targets = *t
	}
}

func load_FxVertexInputBinding(xn *xmlx.Node, obj *nga.FxVertexInputBinding) {
	obj.InputSemantic, obj.Semantic = xas(xn, "input_semantic"), xas(xn, "semantic")
	obj.InputSet = xau64p(xn, "input_set")
}

func load_Float3(xn *xmlx.Node, obj *nga.Float3) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_PxSceneDef(xn *xmlx.Node, obj *nga.PxSceneDef) {
	obj.ForceFields = objs_PxForceFieldInst(xn, "instance_force_field")
	obj.Models = objs_PxModelInst(xn, "instance_physics_model")
	if tcn := node_TechCommon(xn); tcn != nil {
		obj.TC.Gravity = obj_SidVec3(tcn, "gravity")
		obj.TC.TimeStep = obj_SidFloat(tcn, "time_step")
	}
}

func load_InputShared(xn *xmlx.Node, obj *nga.InputShared) {
	if in := obj_Input(xn, ""); in != nil {
		obj.Input = *in
	}
	obj.Offset = xau64(xn, "offset")
	obj.Set = xau64p(xn, "set")
}

func load_SidBool(xn *xmlx.Node, obj *nga.SidBool) {
	obj.B = xb(xn, "")
}

func load_Int3(xn *xmlx.Node, obj *nga.Int3) {
	arr_Ints(xn, len(obj), func(i int, n int64) {
		obj[i] = n
	})
}

func load_Float2x2(xn *xmlx.Node, obj *nga.Float2x2) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_GeometryBrepEdges(xn *xmlx.Node, obj *nga.GeometryBrepEdges) {
	if ii := obj_IndexedInputs(xn, ""); ii != nil {
		obj.IndexedInputs = *ii
	}
}

func load_GeometryMesh(xn *xmlx.Node, obj *nga.GeometryMesh) {
	obj.ConvexHullOf.Set(xas(xn, "convex_hull_of"))
	for _, pn := range xcns(xn, "lines", "linestrips", "polygons", "polylist", "triangles", "trifans", "tristrips") {
		if p := obj_GeometryPrimitives(pn, ""); p != nil {
			obj.Primitives = append(obj.Primitives, p)
		}
	}
	obj.Vertices = obj_GeometryVertices(xn, "vertices")
}

func load_MaterialBinding(xn *xmlx.Node, obj *nga.MaterialBinding) {
	obj.Params = objs_Param(xn, "param")
	if tcn := node_TechCommon(xn); tcn != nil {
		obj.TC.Materials = objs_FxMaterialInst(tcn, "instance_material")
	}
}

func load_FxEffectInst(xn *xmlx.Node, obj *nga.FxEffectInst) {
	obj.TechniqueHints = objs_FxEffectInstTechniqueHint(xn, "technique_hint")
}

func load_KxJointLimits(xn *xmlx.Node, obj *nga.KxJointLimits) {
	obj.Max = obj_SidFloat(xn, "max")
	obj.Min = obj_SidFloat(xn, "min")
}

func load_FxSamplerWrapping(xn *xmlx.Node, obj *nga.FxSamplerWrapping) {
	if cn := xcn(xn, "border_color"); cn != nil {
		list_Rgba32(cn, &obj.BorderColor)
	}
	for n, i := range map[string]*int{"wrap_s": &obj.WrapS, "wrap_t": &obj.WrapT, "wrap_p": &obj.WrapP} {
		switch strings.ToUpper(xs(xn, n)) {
		case "BORDER":
			*i = nga.FX_SAMPLER_WRAP_BORDER
		case "CLAMP":
			*i = nga.FX_SAMPLER_WRAP_CLAMP
		case "MIRROR":
			*i = nga.FX_SAMPLER_WRAP_MIRROR
		case "MIRROR_ONCE":
			*i = nga.FX_SAMPLER_WRAP_MIRROR_ONCE
		default:
			*i = nga.FX_SAMPLER_WRAP_WRAP
		}
	}
}

func load_AnimationInst(xn *xmlx.Node, obj *nga.AnimationInst) {
}

func load_GeometryBrepWires(xn *xmlx.Node, obj *nga.GeometryBrepWires) {
	if iiv := obj_IndexedInputs(xn, ""); iiv != nil {
		obj.IndexedInputs = *iiv
	}
}

func load_FxCreateFormat(xn *xmlx.Node, obj *nga.FxCreateFormat) {
	obj.Exact = xs(xn, "exact")
	obj.Hint = obj_FxCreateFormatHint(xn, "hint")
}

func load_FormulaInst(xn *xmlx.Node, obj *nga.FormulaInst) {
}

func load_CameraOptics(xn *xmlx.Node, obj *nga.CameraOptics) {
	if tcn := node_TechCommon(xn); tcn != nil {
		cn, on, pn := tcn, xcn(tcn, "orthographic"), xcn(tcn, "perspective")
		if cn = pn; cn == nil {
			cn = on
		}
		if cn != nil {
			sf := obj_SidFloat(cn, "aspect_ratio")
			obj.TC.AspectRatio = sf
			if sf = obj_SidFloat(cn, "zfar"); sf != nil {
				obj.TC.Zfar = *sf
			}
			if sf = obj_SidFloat(cn, "znear"); sf != nil {
				obj.TC.Znear = *sf
			}
		}
		if pn != nil {
			obj.TC.Perspective = obj_CameraPerspective(pn, "")
		}
		if on != nil {
			obj.TC.Orthographic = obj_CameraOrthographic(on, "")
		}
	}
}

func load_PxRigidBodyInst(xn *xmlx.Node, obj *nga.PxRigidBodyInst) {
	obj.TargetNode.Set(xas(xn, "target"))
	if tcn := node_TechCommon(xn); tcn != nil {
		if rbc := obj_PxRigidBodyCommon(tcn, ""); rbc != nil {
			obj.TC.PxRigidBodyCommon = *rbc
		}
		v3 := xv3(tcn, "angular_velocity")
		if v3 != nil {
			obj.TC.AngularVelocity = *v3
		}
		if v3 = xv3(tcn, "velocity"); v3 != nil {
			obj.TC.LinearVelocity = *v3
		}
	}
}

func load_FxAnnotation(xn *xmlx.Node, obj *nga.FxAnnotation) {
	for _, cn := range xn.Children {
		if cn.Type == xmlx.NT_ELEMENT {
			if obj.Value = xv(cn); obj.Value != nil {
				break
			}
		}
	}
}

func load_FxCreateFormatHint(xn *xmlx.Node, obj *nga.FxCreateFormatHint) {
	obj.Space = xas(xn, "space")
	switch strings.ToUpper(xas(xn, "channels")) {
	case "RGB":
		obj.Channels = nga.FX_CREATE_FORMAT_HINT_CHANNELS_RGB
	case "RGBA":
		obj.Channels = nga.FX_CREATE_FORMAT_HINT_CHANNELS_RGBA
	case "RGBE":
		obj.Channels = nga.FX_CREATE_FORMAT_HINT_CHANNELS_RGBE
	case "L":
		obj.Channels = nga.FX_CREATE_FORMAT_HINT_CHANNELS_L
	case "LA":
		obj.Channels = nga.FX_CREATE_FORMAT_HINT_CHANNELS_LA
	case "D":
		obj.Channels = nga.FX_CREATE_FORMAT_HINT_CHANNELS_D
	}
	switch strings.ToUpper(xas(xn, "range")) {
	case "FLOAT":
		obj.Range = nga.FX_CREATE_FORMAT_HINT_RANGE_FLOAT
	case "SINT":
		obj.Range = nga.FX_CREATE_FORMAT_HINT_RANGE_SINT
	case "SNORM":
		obj.Range = nga.FX_CREATE_FORMAT_HINT_RANGE_SNORM
	case "UINT":
		obj.Range = nga.FX_CREATE_FORMAT_HINT_RANGE_UINT
	case "UNORM":
		obj.Range = nga.FX_CREATE_FORMAT_HINT_RANGE_UNORM
	}
	switch strings.ToUpper(xas(xn, "precision")) {
	case "LOW":
		obj.Precision = nga.FX_CREATE_FORMAT_HINT_PRECISION_LOW
	case "MID":
		obj.Precision = nga.FX_CREATE_FORMAT_HINT_PRECISION_MID
	case "HIGH":
		obj.Precision = nga.FX_CREATE_FORMAT_HINT_PRECISION_HIGH
	case "MAX":
		obj.Precision = nga.FX_CREATE_FORMAT_HINT_PRECISION_MAX
	default:
		obj.Precision = nga.FX_CREATE_FORMAT_HINT_PRECISION_DEFAULT
	}
}

func load_GeometryBrepCircle(xn *xmlx.Node, obj *nga.GeometryBrepCircle) {
	obj.Radius = xf64(xn, "radius")
}

func load_PxSceneInst(xn *xmlx.Node, obj *nga.PxSceneInst) {
}

func load_FxParamDef(xn *xmlx.Node, obj *nga.FxParamDef) {
	if pd := obj_ParamDef(xn, ""); pd != nil {
		obj.ParamDef = *pd
	}
	obj.Annotations = objs_FxAnnotation(xn, "annotate")
	obj.Modifier = xs(xn, "modifier")
	obj.Semantic = xs(xn, "semantic")
}

func load_FxCreateMips(xn *xmlx.Node, obj *nga.FxCreateMips) {
	obj.Levels = xau64(xn, "levels")
	if b := xabp(xn, "auto_generate"); b != nil {
		obj.NoAutoGen = !*b
	}
}

func load_KxFrame(xn *xmlx.Node, obj *nga.KxFrame) {
	switch xn.Name.Local {
	case "frame_object":
		obj.Type = nga.KX_FRAME_TYPE_OBJECT
	case "frame_origin":
		obj.Type = nga.KX_FRAME_TYPE_ORIGIN
	case "frame_tcp":
		obj.Type = nga.KX_FRAME_TYPE_TCP
	case "frame_tip":
		obj.Type = nga.KX_FRAME_TYPE_TIP
	}
	if obj.Type > 0 {
		obj.Link.Set(xas(xn, "link"))
		obj.Transforms = get_Transforms(xn)
	}
}

func load_KxJointInst(xn *xmlx.Node, obj *nga.KxJointInst) {
}

func load_Input(xn *xmlx.Node, obj *nga.Input) {
	obj.Semantic = xas(xn, "semantic")
	obj.Source.Set(xas(xn, "source"))
}

func load_ControllerDef(xn *xmlx.Node, obj *nga.ControllerDef) {
	obj.Morph = obj_ControllerMorph(xn, "morph")
	obj.Skin = obj_ControllerSkin(xn, "skin")
}

func load_FxCreate(xn *xmlx.Node, obj *nga.FxCreate) {
	obj.Format = obj_FxCreateFormat(xn, "format")
	if an := xcn(xn, "array"); an != nil {
		obj.ArrayLength = xau64(an, "length")
	}
}

func load_FxTechniqueCommonConstant(xn *xmlx.Node, obj *nga.FxTechniqueCommonConstant) {
	obj.Emission = obj_FxColorOrTexture(xn, "emission")
	obj.IndexOfRefraction = obj_ParamOrSidFloat(xn, "index_of_refraction")
	obj.Reflective = obj_FxColorOrTexture(xn, "reflective")
	obj.Reflectivity = obj_ParamOrSidFloat(xn, "reflectivity")
	obj.Transparency = obj_ParamOrSidFloat(xn, "transparency")
	obj.Transparent = obj_FxColorOrTexture(xn, "transparent")
}

func load_PxRigidConstraintLimit(xn *xmlx.Node, obj *nga.PxRigidConstraintLimit) {
	sv3 := obj_SidVec3(xn, "max")
	if sv3 != nil {
		obj.Max = *sv3
	}
	if sv3 = obj_SidVec3(xn, "min"); sv3 != nil {
		obj.Min = *sv3
	}
}

func load_SidVec3(xn *xmlx.Node, obj *nga.SidVec3) {
	if v3 := xv3(xn, ""); v3 != nil {
		obj.Vec3 = *v3
	}
}

func load_NodeInst(xn *xmlx.Node, obj *nga.NodeInst) {
	obj.Proxy.Set(xas(xn, "proxy"))
}

func load_CameraOrthographic(xn *xmlx.Node, obj *nga.CameraOrthographic) {
	obj.MagX = obj_SidFloat(xn, "xmag")
	obj.MagY = obj_SidFloat(xn, "ymag")
}

func load_Float3x3(xn *xmlx.Node, obj *nga.Float3x3) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_FxPassState(xn *xmlx.Node, obj *nga.FxPassState) {
	obj.Value = xas(xn, "value")
	obj.Param.Set(xas(xn, "param"))
	obj.Index = xf64(xn, "index")
}

func load_GeometryBrepSurface(xn *xmlx.Node, obj *nga.GeometryBrepSurface) {
	if pos := obj_GeometryPositioning(xn, ""); pos != nil {
		obj.Location = *pos
	}
	obj.Element.Cone = obj_GeometryBrepCone(xn, "cone")
	obj.Element.Plane = obj_GeometryBrepPlane(xn, "plane")
	obj.Element.Cylinder = obj_GeometryBrepCylinder(xn, "cylinder")
	obj.Element.NurbsSurface = obj_GeometryBrepNurbsSurface(xn, "nurbs_surface")
	obj.Element.Sphere = obj_GeometryBrepSphere(xn, "sphere")
	obj.Element.Torus = obj_GeometryBrepTorus(xn, "torus")
	obj.Element.SweptSurface = obj_GeometryBrepSweptSurface(xn, "swept_surface")
}

func load_AnimationSampler(xn *xmlx.Node, obj *nga.AnimationSampler) {
	for n, pb := range map[string]*int{"pre_behavior": &obj.PreBehavior, "post_behavior": &obj.PostBehavior} {
		switch strings.ToUpper(xas(xn, n)) {
		case "CONSTANT":
			*pb = nga.ANIM_SAMPLER_BEHAVIOR_CONSTANT
		case "CYCLE":
			*pb = nga.ANIM_SAMPLER_BEHAVIOR_CYCLE
		case "CYCLE_RELATIVE":
			*pb = nga.ANIM_SAMPLER_BEHAVIOR_CYCLE_RELATIVE
		case "GRADIENT":
			*pb = nga.ANIM_SAMPLER_BEHAVIOR_GRADIENT
		case "OSCILLATE":
			*pb = nga.ANIM_SAMPLER_BEHAVIOR_OSCILLATE
		default:
			*pb = nga.ANIM_SAMPLER_BEHAVIOR_UNDEFINED
		}
	}
}

func load_ParamOrUint(xn *xmlx.Node, obj *nga.ParamOrUint) {
	obj.Param.Set(xs(xn, "param"))
	obj.U = xu64(xn, "uint")
}

func load_ParamDefs(xn *xmlx.Node, obj *nga.ParamDefs) {
}

func load_ParamInsts(xn *xmlx.Node, obj *nga.ParamInsts) {
}

func load_AnimationChannel(xn *xmlx.Node, obj *nga.AnimationChannel) {
	obj.Source.Set(xas(xn, "source"))
	obj.Target.Set(xas(xn, "target"))
}

func load_LightAmbient(xn *xmlx.Node, obj *nga.LightAmbient) {
	get_LightColor(xn, &obj.Color)
}

func load_CameraPerspective(xn *xmlx.Node, obj *nga.CameraPerspective) {
	obj.FovX = obj_SidFloat(xn, "xfov")
	obj.FovY = obj_SidFloat(xn, "yfov")
}

func load_PxForceFieldDef(xn *xmlx.Node, obj *nga.PxForceFieldDef) {
}

func load_PxRigidConstraintSpring(xn *xmlx.Node, obj *nga.PxRigidConstraintSpring) {
	sf := obj_SidFloat(xn, "damping")
	if sf != nil {
		obj.Damping = *sf
	}
	if sf = obj_SidFloat(xn, "target_value"); sf != nil {
		obj.TargetValue = *sf
	}
	if sf = obj_SidFloat(xn, "stiffness"); sf != nil {
		obj.Stiffness = *sf
	} else {
		obj.Stiffness.F = 1
	}
}

func load_GeometryBrepNurbsSurface(xn *xmlx.Node, obj *nga.GeometryBrepNurbsSurface) {
	obj.U.Degree = xau64(xn, "degree_u")
	obj.U.Closed = xab(xn, "closed_u")
	obj.V.Degree = xau64(xn, "degree_v")
	obj.V.Closed = xab(xn, "closed_v")
	if cv := obj_GeometryControlVertices(xn, "control_vertices"); cv != nil {
		obj.ControlVertices = *cv
	}
}

func load_CameraDef(xn *xmlx.Node, obj *nga.CameraDef) {
	obj.Imager = obj_CameraImager(xn, "imager")
	if op := obj_CameraOptics(xn, "optics"); op != nil {
		obj.Optics = *op
	}
}

func load_GeometryBrepSweptSurface(xn *xmlx.Node, obj *nga.GeometryBrepSweptSurface) {
	obj.Curve = obj_GeometryBrepCurve(xn, "curve")
	obj.Extrusion.Direction = xv3(xn, "direction")
	obj.Revolution.Origin = xv3(xn, "origin")
	obj.Revolution.Direction = xv3(xn, "axis")
}

func load_FxPassProgramBindUniform(xn *xmlx.Node, obj *nga.FxPassProgramBindUniform) {
	obj.Symbol = xas(xn, "symbol")
	obj.ParamRef.Set(get_ParamRef(xn, "param"))
	for _, cn := range xn.Children {
		if cn.Type == xmlx.NT_ELEMENT {
			if obj.Value = xv(cn); obj.Value != nil {
				break
			}
		}
	}
}

func load_KxMotionAxis(xn *xmlx.Node, obj *nga.KxMotionAxis) {
	obj.Axis.Set(xas(xn, "axis"))
	obj.Bindings = objs_KxBinding(xn, "bind")
	obj.Speed = obj_ParamOrFloat(xn, "speed")
	obj.Acceleration = obj_ParamOrFloat(xn, "acceleration")
	obj.Deceleration = obj_ParamOrFloat(xn, "deceleration")
	obj.Jerk = obj_ParamOrFloat(xn, "jerk")
}

func load_AnimationClipDef(xn *xmlx.Node, obj *nga.AnimationClipDef) {
	obj.Start = xf64(xn, "start")
	obj.End = xf64(xn, "end")
	obj.Animations = objs_AnimationInst(xn, "instance_animation")
	obj.Formulas = objs_FormulaInst(xn, "instance_formula")
}

func load_Float7(xn *xmlx.Node, obj *nga.Float7) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_AnimationDef(xn *xmlx.Node, obj *nga.AnimationDef) {
	obj.AnimationDefs = objs_AnimationDef(xn, "animation")
	obj.Samplers = objs_AnimationSampler(xn, "sampler")
	obj.Channels = objs_AnimationChannel(xn, "channel")
}

func load_PxMaterial(xn *xmlx.Node, obj *nga.PxMaterial) {
	obj.Def = obj_PxMaterialDef(xn, "physics_material")
	obj.Inst = obj_PxMaterialInst(xn, "instance_physics_material")
}

func load_PxShape(xn *xmlx.Node, obj *nga.PxShape) {
	if sb := obj_SidBool(xn, "hollow"); sb != nil {
		obj.Hollow = *sb
	}
	obj.Density = obj_SidFloat(xn, "density")
	obj.Mass = obj_SidFloat(xn, "mass")
	obj.Transforms = get_Transforms(xn)
	if pm := obj_PxMaterial(xn, ""); pm != nil {
		obj.Material = *pm
	}
	obj.Geometry.Plane = obj_GeometryBrepPlane(xn, "plane")
	obj.Geometry.Box = obj_GeometryBrepBox(xn, "box")
	obj.Geometry.Sphere = obj_GeometryBrepSphere(xn, "sphere")
	obj.Geometry.Cylinder = obj_PxCylinder(xn, "cylinder")
	obj.Geometry.Capsule = obj_GeometryBrepCapsule(xn, "capsule")
	obj.Geometry.Inst = obj_GeometryInst(xn, "instance_geometry")
}

func load_LightPoint(xn *xmlx.Node, obj *nga.LightPoint) {
	get_LightColor(xn, &obj.Color)
	if a := obj_LightAttenuation(xn, ""); a != nil {
		obj.Attenuation = *a
	}
}

func load_GeometrySpline(xn *xmlx.Node, obj *nga.GeometrySpline) {
	obj.Closed = xab(xn, "closed")
	if cv := obj_GeometryControlVertices(xn, "control_vertices"); cv != nil {
		obj.ControlVertices = *cv
	}
}

func load_GeometryInst(xn *xmlx.Node, obj *nga.GeometryInst) {
	obj.MaterialBinding = obj_MaterialBinding(xn, "bind_material")
}

func load_PxRigidConstraintDef(xn *xmlx.Node, obj *nga.PxRigidConstraintDef) {
	a := obj_PxRigidConstraintAttachment(xn, "attachment")
	if a != nil {
		obj.Attachment = *a
	}
	if a = obj_PxRigidConstraintAttachment(xn, "ref_attachment"); a != nil {
		obj.RefAttachment = *a
	}
	if tcn := node_TechCommon(xn); tcn != nil {
		sb := obj_SidBool(tcn, "enabled")
		if sb != nil {
			obj.TC.Enabled = *sb
		}
		if sb = obj_SidBool(tcn, "interpenetrate"); sb != nil {
			obj.TC.Interpenetrate = *sb
		}
		cn := xcn(tcn, "limits")
		if cn != nil {
			obj.TC.Limits.Angular = obj_PxRigidConstraintLimit(cn, "swing_cone_and_twist")
			obj.TC.Limits.Linear = obj_PxRigidConstraintLimit(cn, "linear")
		}
		if cn = xcn(tcn, "spring"); cn != nil {
			obj.TC.Spring.Angular = obj_PxRigidConstraintSpring(cn, "angular")
			obj.TC.Spring.Linear = obj_PxRigidConstraintSpring(cn, "linear")
		}
	}
}

func load_GeometryBrepSurfaceCurves(xn *xmlx.Node, obj *nga.GeometryBrepSurfaceCurves) {
	obj.All = objs_GeometryBrepCurve(xn, "curve")
}

func load_Int3x3(xn *xmlx.Node, obj *nga.Int3x3) {
	arr_Ints(xn, len(obj), func(i int, n int64) {
		obj[i] = n
	})
}

func load_FxTechniqueCommonBlinn(xn *xmlx.Node, obj *nga.FxTechniqueCommonBlinn) {
	if l := obj_FxTechniqueCommonLambert(xn, ""); l != nil {
		obj.FxTechniqueCommonLambert = *l
	}
	obj.Specular = obj_FxColorOrTexture(xn, "specular")
	obj.Shininess = obj_ParamOrSidFloat(xn, "shininess")
}

func load_SidString(xn *xmlx.Node, obj *nga.SidString) {
	obj.S = xn.Value
}

func load_GeometryDef(xn *xmlx.Node, obj *nga.GeometryDef) {
	obj.Brep = obj_GeometryBrep(xn, "brep")
	obj.Spline = obj_GeometrySpline(xn, "spline")
	if obj.Mesh = obj_GeometryMesh(xn, "mesh"); obj.Mesh == nil {
		obj.Mesh = obj_GeometryMesh(xn, "convex_mesh")
	}
}

func load_FxCreate2DSizeExact(xn *xmlx.Node, obj *nga.FxCreate2DSizeExact) {
	obj.Width = xau64(xn, "width")
	obj.Height = xau64(xn, "height")
}

func load_VisualSceneInst(xn *xmlx.Node, obj *nga.VisualSceneInst) {
}

func load_Int2(xn *xmlx.Node, obj *nga.Int2) {
	arr_Ints(xn, len(obj), func(i int, n int64) {
		obj[i] = n
	})
}

func load_Float2(xn *xmlx.Node, obj *nga.Float2) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_KxBinding(xn *xmlx.Node, obj *nga.KxBinding) {
	obj.Symbol = xas(xn, "symbol")
	obj.Param.Set(get_ParamRef(xn, "param"))
	for _, cn := range xn.Children {
		if cn.Type == xmlx.NT_ELEMENT {
			if obj.Value = xv(cn); obj.Value != nil {
				break
			}
		}
	}
}

func load_FxPassEvaluationTarget(xn *xmlx.Node, obj *nga.FxPassEvaluationTarget) {
	if u := xau64p(xn, "index"); u != nil {
		obj.Index = *u
	} else {
		obj.Index = 1
	}
	obj.Slice = xau64(xn, "slice")
	obj.Mip = xau64(xn, "mip")
	obj.CubeFace = get_CubeFace(xn)
	obj.Sampler.Set(get_ParamRef(xn, "param"))
	obj.Image = obj_FxImageInst(xn, "instance_image")
}

func load_FxColorOrTexture(xn *xmlx.Node, obj *nga.FxColorOrTexture) {
	obj.ParamRef.Set(get_ParamRef(xn, "param"))
	obj.Texture = obj_FxTexture(xn, "texture")
	if cn := xcn(xn, "color"); cn != nil {
		list_Rgba32(cn, obj.Color)
	}
	switch strings.ToUpper(xas(xn, "opaque")) {
	case "RGB_ZERO":
		obj.Opaque = nga.FX_COLOR_TEXTURE_OPAQUE_RGB_ZERO
	case "RGB_ONE":
		obj.Opaque = nga.FX_COLOR_TEXTURE_OPAQUE_RGB_ONE
	case "A_ZERO":
		obj.Opaque = nga.FX_COLOR_TEXTURE_OPAQUE_A_ZERO
	default:
		obj.Opaque = nga.FX_COLOR_TEXTURE_OPAQUE_A_ONE
	}
}

func load_GeometryVertices(xn *xmlx.Node, obj *nga.GeometryVertices) {
}

func load_KxMotionSystem(xn *xmlx.Node, obj *nga.KxMotionSystem) {
	obj.ArticulatedSystem = obj_KxArticulatedSystemInst(xn, "instance_articulated_system")
	if tcn := node_TechCommon(xn); tcn != nil {
		obj.TC.AxisInfos = objs_KxMotionAxis(tcn, "axis_info")
		obj.TC.EffectorInfo = obj_KxEffector(tcn, "effector_info")
	}
}

func load_Int4x4(xn *xmlx.Node, obj *nga.Int4x4) {
	arr_Ints(xn, len(obj), func(i int, n int64) {
		obj[i] = n
	})
}

func load_CameraImager(xn *xmlx.Node, obj *nga.CameraImager) {
}

func load_FxPassEvaluationClearDepth(xn *xmlx.Node, obj *nga.FxPassEvaluationClearDepth) {
	obj.Index = xau64(xn, "index")
	obj.F = xf64(xn, "")
}

func load_GeometryBrepEllipse(xn *xmlx.Node, obj *nga.GeometryBrepEllipse) {
	if f2 := obj_Float2(xn, "radius"); f2 != nil {
		obj.Radii = *f2
	}
}

func load_FxBinding(xn *xmlx.Node, obj *nga.FxBinding) {
	obj.Semantic = xas(xn, "semantic")
	obj.Target.Set(xas(xn, "target"))
}

func load_GeometryBrepParabola(xn *xmlx.Node, obj *nga.GeometryBrepParabola) {
	obj.FocalLength = xf64(xn, "focal")
}

func load_SidFloat(xn *xmlx.Node, obj *nga.SidFloat) {
	obj.F = xf64(xn, "")
}

func load_Bool3(xn *xmlx.Node, obj *nga.Bool3) {
	arr_Bools(xn, len(obj), func(i int, b bool) {
		obj[i] = b
	})
}

func load_VisualSceneDef(xn *xmlx.Node, obj *nga.VisualSceneDef) {
	obj.Evaluations = objs_VisualSceneEvaluation(xn, "evaluate_scene")
	obj.Nodes = objs_NodeDef(xn, "node")
}

func load_KxSceneInst(xn *xmlx.Node, obj *nga.KxSceneInst) {
	obj.ModelBindings = objs_KxModelBinding(xn, "bind_kinematics_model")
	obj.JointAxisBindings = objs_KxJointAxisBinding(xn, "bind_joint_axis")
}

func load_LightDirectional(xn *xmlx.Node, obj *nga.LightDirectional) {
	get_LightColor(xn, &obj.Color)
}

func load_KxAxisIndex(xn *xmlx.Node, obj *nga.KxAxisIndex) {
	obj.Semantic = xas(xn, "semantic")
	if pi := obj_ParamOrInt(xn, ""); pi != nil {
		obj.I = *pi
	}
}

func load_VisualSceneRendering(xn *xmlx.Node, obj *nga.VisualSceneRendering) {
	obj.CameraNode.Set(xas(xn, "camera_node"))
	for _, ln := range xcns(xn, "layer") {
		obj.Layers[ln.Value] = true
	}
	obj.MaterialInst = obj_VisualSceneRenderingMaterialInst(xn, "instance_material")
}

func load_GeometryPolygonHole(xn *xmlx.Node, obj *nga.GeometryPolygonHole) {
	obj.Indices = listcn_Uints(xn, "p")
	for _, cn := range xcns(xn, "h") {
		obj.Holes = append(obj.Holes, list_Uints(cn))
	}
}

func load_Param(xn *xmlx.Node, obj *nga.Param) {
	obj.Type = xas(xn, "type")
	obj.Semantic = xas(xn, "semantic")
}

func load_FxTechniqueCommonLambert(xn *xmlx.Node, obj *nga.FxTechniqueCommonLambert) {
	if c := obj_FxTechniqueCommonConstant(xn, ""); c != nil {
		obj.FxTechniqueCommonConstant = *c
	}
	obj.Ambient = obj_FxColorOrTexture(xn, "ambient")
	obj.Diffuse = obj_FxColorOrTexture(xn, "diffuse")
}

func load_ParamOrSidFloat(xn *xmlx.Node, obj *nga.ParamOrSidFloat) {
	obj.Param.Set(xs(xn, "param"))
	if f := obj_SidFloat(xn, "float"); f != nil {
		obj.F = *f
	}
}

func load_Float2x3(xn *xmlx.Node, obj *nga.Float2x3) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_FxCreateInitFrom(xn *xmlx.Node, obj *nga.FxCreateInitFrom) {
	if i := obj_FxInitFrom(xn, ""); i != nil {
		obj.FxInitFrom = *i
	}
	obj.ArrayIndex = xau64(xn, "array_index")
	obj.MipIndex = xau64(xn, "mip_index")
}

func load_LightAttenuation(xn *xmlx.Node, obj *nga.LightAttenuation) {
	var sf *nga.SidFloat
	for n, f := range map[string]*nga.SidFloat{"constant_attenuation": &obj.Constant, "linear_attenuation": &obj.Linear, "quadratic_attenuation": &obj.Quadratic} {
		if sf = obj_SidFloat(xn, n); sf != nil {
			*f = *sf
		}
	}
}

func load_KxModelBinding(xn *xmlx.Node, obj *nga.KxModelBinding) {
	obj.Node.Set(xas(xn, "node"))
	obj.Model.SidRef.Set(xss(xn, "SIDREF", "sidref"))
	obj.Model.ParamRef.Set(get_ParamRef(xn, "param"))
}

func load_GeometryBrepSolids(xn *xmlx.Node, obj *nga.GeometryBrepSolids) {
	if iiv := obj_IndexedInputs(xn, ""); iiv != nil {
		obj.IndexedInputs = *iiv
	}
}

func load_FxMaterialDef(xn *xmlx.Node, obj *nga.FxMaterialDef) {
	if e := obj_FxEffectInst(xn, "instance_effect"); e != nil {
		obj.Effect = *e
	}
}

func load_ParamOrFloat(xn *xmlx.Node, obj *nga.ParamOrFloat) {
	obj.Param.Set(xs(xn, "param"))
	obj.F = xf64(xn, "float")
}

func load_SidFloat3(xn *xmlx.Node, obj *nga.SidFloat3) {
	if f3 := obj_Float3(xn, ""); f3 != nil {
		obj.F = *f3
	}
}

func load_FxInitFrom(xn *xmlx.Node, obj *nga.FxInitFrom) {
	obj.RefUrl = xs(xn, "ref")
	if hn := xcn(xn, "hex"); hn != nil {
		obj.Raw.Format = xas(hn, "format")
		obj.Raw.Data, _ = hex.DecodeString(xs(hn, ""))
	}
}

func load_FxTexture(xn *xmlx.Node, obj *nga.FxTexture) {
	obj.TexCoord = xas(xn, "texcoord")
	obj.Sampler2D.Set(xas(xn, "texture"))
}

func load_Int4(xn *xmlx.Node, obj *nga.Int4) {
	arr_Ints(xn, len(obj), func(i int, n int64) {
		obj[i] = n
	})
}

func load_GeometryBrepOrientation(xn *xmlx.Node, obj *nga.GeometryBrepOrientation) {
	fs := list_Floats(xn)
	if len(fs) > 0 {
		if obj.Axis.X = fs[0]; len(fs) > 1 {
			if obj.Axis.Y = fs[1]; len(fs) > 2 {
				if obj.Axis.Z = fs[2]; len(fs) > 3 {
					obj.Angle = fs[3]
				}
			}
		}
	}
}

func load_GeometryPrimitives(xn *xmlx.Node, obj *nga.GeometryPrimitives) {
	switch xn.Name.Local {
	case "lines":
		obj.Type = nga.GEOMETRY_PRIMITIVE_TYPE_LINES
	case "linestrips":
		obj.Type = nga.GEOMETRY_PRIMITIVE_TYPE_LINE_STRIPS
	case "polygons":
		obj.Type = nga.GEOMETRY_PRIMITIVE_TYPE_POLYGONS
	case "polylist":
		obj.Type = nga.GEOMETRY_PRIMITIVE_TYPE_POLYLIST
	case "triangles":
		obj.Type = nga.GEOMETRY_PRIMITIVE_TYPE_TRIANGLES
	case "trifans":
		obj.Type = nga.GEOMETRY_PRIMITIVE_TYPE_TRIFANS
	case "tristrips":
		obj.Type = nga.GEOMETRY_PRIMITIVE_TYPE_TRISTRIPS
	}
	if obj.Type > 0 {
		if ii := obj_IndexedInputs(xn, ""); ii != nil {
			obj.IndexedInputs = *ii
		}
		obj.Material = xas(xn, "material")
		obj.PolyHoles = objs_GeometryPolygonHole(xn, "ph")
	}
}

func load_Sources(xn *xmlx.Node, obj *nga.Sources) {
}

func load_PxRigidConstraintAttachment(xn *xmlx.Node, obj *nga.PxRigidConstraintAttachment) {
	obj.RigidBody.Set(xas(xn, "rigid_body"))
	obj.Transforms = get_Transforms(xn)
}

func load_FxImageInitFrom(xn *xmlx.Node, obj *nga.FxImageInitFrom) {
	if i := obj_FxInitFrom(xn, ""); i != nil {
		obj.FxInitFrom = *i
	}
	if b := xabp(xn, "mips_generate"); b != nil {
		obj.NoAutoMip = !*b
	}
}

func load_FxTechniqueCommon(xn *xmlx.Node, obj *nga.FxTechniqueCommon) {
	if t := obj_FxTechnique(xn, ""); t != nil {
		obj.FxTechnique = *t
	}
	obj.Blinn = obj_FxTechniqueCommonBlinn(xn, "blinn")
	obj.Constant = obj_FxTechniqueCommonConstant(xn, "constant")
	obj.Lambert = obj_FxTechniqueCommonLambert(xn, "lambert")
	obj.Phong = obj_FxTechniqueCommonPhong(xn, "phong")
}

func load_Float3x4(xn *xmlx.Node, obj *nga.Float3x4) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_KxJointAxisBinding(xn *xmlx.Node, obj *nga.KxJointAxisBinding) {
	obj.Target.Set(xas(xn, "target"))
	if pf := obj_ParamOrFloat(xn, "value"); pf != nil {
		obj.Value = *pf
	}
	if psr := obj_ParamOrRefSid(xn, "axis"); psr != nil {
		obj.Axis = *psr
	}
}

func load_KxJointDef(xn *xmlx.Node, obj *nga.KxJointDef) {
	var j *nga.KxJoint
	for _, cn := range xcns(xn, "prismatic", "revolute") {
		if j = obj_KxJoint(cn, ""); (j != nil) && (j.Type > 0) {
			obj.All = append(obj.All, j)
		}
	}
}

func load_VisualSceneEvaluation(xn *xmlx.Node, obj *nga.VisualSceneEvaluation) {
	if b := xabp(xn, "enable"); b != nil {
		obj.Disabled = !*b
	}
	obj.RenderPasses = objs_VisualSceneRendering(xn, "render")
}

func load_FxProfileGlslCodeInclude(xn *xmlx.Node, obj *nga.FxProfileGlslCodeInclude) {
	if obj.S, obj.IsInclude = xn.Value, (xn.Name.Local == "include"); obj.IsInclude {
		obj.S = xas(xn, "url")
	}
}

func load_AssetGeographicLocation(xn *xmlx.Node, obj *nga.AssetGeographicLocation) {
	obj.Longitude = xf64(xn, "longitude")
	obj.Latitude = xf64(xn, "latitude")
	if an := xcn(xn, "altitude"); an != nil {
		obj.Altitude = xf64(an, "")
		obj.AltitudeAbsolute = xab(an, "absolute")
	}
}

func load_KxEffector(xn *xmlx.Node, obj *nga.KxEffector) {
	obj.Bindings = objs_KxBinding(xn, "bind")
	obj.Speed = obj_ParamOrFloat2(xn, "speed")
	obj.Acceleration = obj_ParamOrFloat2(xn, "acceleration")
	obj.Deceleration = obj_ParamOrFloat2(xn, "deceleration")
	obj.Jerk = obj_ParamOrFloat2(xn, "jerk")
}

func load_GeometryBrepTorus(xn *xmlx.Node, obj *nga.GeometryBrepTorus) {
	if f2 := obj_Float2(xn, "radius"); f2 != nil {
		obj.Radii = *f2
	}
}

func load_Layers(xn *xmlx.Node, obj *nga.Layers) {
}

func load_FxProfile(xn *xmlx.Node, obj *nga.FxProfile) {
	switch xn.Name.Local {
	case "profile_GLSL":
		obj.Glsl = obj_FxProfileGlsl(xn, "")
	case "profile_COMMON":
		obj.Common = obj_FxProfileCommon(xn, "")
	}
}

func load_KxAxisLimits(xn *xmlx.Node, obj *nga.KxAxisLimits) {
	pf := obj_ParamOrFloat(xn, "max")
	if pf != nil {
		obj.Max = *pf
	}
	if pf = obj_ParamOrFloat(xn, "min"); pf != nil {
		obj.Min = *pf
	}
}

func load_FxMaterialInst(xn *xmlx.Node, obj *nga.FxMaterialInst) {
	obj.Symbol = xas(xn, "symbol")
	obj.VertexInputBindings = objs_FxVertexInputBinding(xn, "bind_vertex_input")
	obj.Bindings = objs_FxBinding(xn, "bind")
}

func load_GeometryBrepCone(xn *xmlx.Node, obj *nga.GeometryBrepCone) {
	obj.Angle = xf64(xn, "angle")
	obj.Radius = xf64(xn, "radius")
}

func load_Extra(xn *xmlx.Node, obj *nga.Extra) {
	obj.Type = xas(xn, "type")
}

func load_ControllerSkin(xn *xmlx.Node, obj *nga.ControllerSkin) {
	obj.Source.Set(xas(xn, "source"))
	obj.BindShapeMatrix = *xm4(xn, "bind_shape_matrix")
	if ci := obj_ControllerInputs(xn, "joints"); ci != nil {
		obj.Joints = *ci
	}
	if iiv := obj_IndexedInputs(xn, "vertex_weights"); iiv != nil {
		obj.VertexWeights = *iiv
	}
}

func load_Bool2(xn *xmlx.Node, obj *nga.Bool2) {
	arr_Bools(xn, len(obj), func(i int, b bool) {
		obj[i] = b
	})
}

func load_VisualSceneRenderingMaterialInst(xn *xmlx.Node, obj *nga.VisualSceneRenderingMaterialInst) {
	obj.Bindings = objs_FxBinding(xn, "bind")
	if tn := xcn(xn, "technique_override"); tn != nil {
		obj.OverrideTechnique.Ref.Set(xas(tn, "ref"))
		obj.OverrideTechnique.Pass.Set(xas(tn, "pass"))
	}
}

func load_ChildNode(xn *xmlx.Node, obj *nga.ChildNode) {
	switch xn.Name.Local {
	case "instance_node":
		obj.Inst = obj_NodeInst(xn, "")
	case "node":
		obj.Def = obj_NodeDef(xn, "")
	}
}

func load_NodeDef(xn *xmlx.Node, obj *nga.NodeDef) {
	var nc *nga.ChildNode
	obj.IsSkinJoint = (xas(xn, "type") == "JOINT")
	if l := xsdt.ListValues(xas(xn, "layer")); len(l) > 0 {
		for _, n := range l {
			obj.Layers[n] = true
		}
	}
	for _, cn := range xcns(xn, "node", "instance_node") {
		if nc = obj_ChildNode(cn, ""); nc != nil {
			obj.Nodes = append(obj.Nodes, *nc)
		}
	}
	obj.Transforms = get_Transforms(xn)
	obj.Insts.Camera = objs_CameraInst(xn, "instance_camera")
	obj.Insts.Controller = objs_ControllerInst(xn, "instance_controller")
	obj.Insts.Geometry = objs_GeometryInst(xn, "instance_geometry")
	obj.Insts.Light = objs_LightInst(xn, "instance_light")
}

func load_PxForceFieldInst(xn *xmlx.Node, obj *nga.PxForceFieldInst) {
}

func load_KxKinematicsSystem(xn *xmlx.Node, obj *nga.KxKinematicsSystem) {
	obj.Models = objs_KxModelInst(xn, "instance_kinematics_model")
	if tcn := node_TechCommon(xn); tcn != nil {
		obj.TC.AxisInfos = objs_KxKinematicsAxis(tcn, "axis_info")
		obj.TC.Frame.Object = obj_KxFrame(tcn, "frame_object")
		obj.TC.Frame.Tcp = obj_KxFrame(tcn, "frame_tcp")
		f := obj_KxFrame(tcn, "frame_tip")
		if f != nil {
			obj.TC.Frame.Tip = *f
		}
		if f = obj_KxFrame(tcn, "frame_origin"); f != nil {
			obj.TC.Frame.Origin = *f
		}
	}
}

func load_GeometryBrepLine(xn *xmlx.Node, obj *nga.GeometryBrepLine) {
	v3 := xv3(xn, "origin")
	if v3 != nil {
		obj.Origin = *v3
	}
	if v3 = xv3(xn, "direction"); v3 != nil {
		obj.Direction = *v3
	}
}

func load_GeometryBrep(xn *xmlx.Node, obj *nga.GeometryBrep) {
	if v := obj_GeometryVertices(xn, "vertices"); v != nil {
		obj.Vertices = *v
	}
	obj.Curves = obj_GeometryBrepCurves(xn, "curves")
	obj.Edges = obj_GeometryBrepEdges(xn, "edges")
	obj.Faces = obj_GeometryBrepFaces(xn, "faces")
	obj.Pcurves = obj_GeometryBrepPcurves(xn, "pcurves")
	obj.Shells = obj_GeometryBrepShells(xn, "shells")
	obj.Solids = obj_GeometryBrepSolids(xn, "solids")
	obj.SurfaceCurves = obj_GeometryBrepSurfaceCurves(xn, "surface_curves")
	obj.Wires = obj_GeometryBrepWires(xn, "wires")
}

func load_KxModelDef(xn *xmlx.Node, obj *nga.KxModelDef) {
	var f *nga.Formula
	if tcn := node_TechCommon(xn); tcn != nil {
		has_ParamDefs(tcn, &obj.TC.HasParamDefs)
		for _, cn := range xcns(tcn, "formula", "instance_formula") {
			if f = obj_Formula(cn, ""); f != nil {
				obj.TC.Formulas = append(obj.TC.Formulas, *f)
			}
		}
		obj.TC.Links = objs_KxLink(tcn, "link")
	}
}

func load_FxImageInst(xn *xmlx.Node, obj *nga.FxImageInst) {
}

func load_GeometryBrepBox(xn *xmlx.Node, obj *nga.GeometryBrepBox) {
	if v3 := xv3(xn, "half_extents"); v3 != nil {
		obj.HalfExtents = *v3
	}
}

func load_ParamOrFloat2(xn *xmlx.Node, obj *nga.ParamOrFloat2) {
	obj.Param.Set(xs(xn, "param"))
	if fn := xcn(xn, "float2"); fn != nil {
		arr_Floats(fn, 2, func(i int, f float64) {
			obj.F[i] = f
		})
	}
}

func load_Technique(xn *xmlx.Node, obj *nga.Technique) {
	obj.Profile = xas(xn, "profile")
	obj.Data = xn.Children
}

func load_FxPassEvaluationClearStencil(xn *xmlx.Node, obj *nga.FxPassEvaluationClearStencil) {
	obj.Index = xau64(xn, "index")
	obj.B = byte(xu64(xn, ""))
}

func load_FxCreate2DSizeRatio(xn *xmlx.Node, obj *nga.FxCreate2DSizeRatio) {
	obj.Width = xaf64(xn, "width")
	obj.Height = xaf64(xn, "height")
}

func load_GeometryBrepPlane(xn *xmlx.Node, obj *nga.GeometryBrepPlane) {
	if f4 := obj_Float4(xn, "equation"); f4 != nil {
		obj.Equation = *f4
	}
}

func load_FxPassProgramBindAttribute(xn *xmlx.Node, obj *nga.FxPassProgramBindAttribute) {
	obj.Symbol = xas(xn, "symbol")
	obj.Semantic = xs(xn, "semantic")
}

func load_FxParamDefs(xn *xmlx.Node, obj *nga.FxParamDefs) {
}

func load_PxMaterialDef(xn *xmlx.Node, obj *nga.PxMaterialDef) {
	var sf *nga.SidFloat
	if tcn := node_TechCommon(xn); tcn != nil {
		for n, f := range map[string]*nga.SidFloat{"dynamic_friction": &obj.TC.DynamicFriction, "restitution": &obj.TC.Restitution, "static_friction": &obj.TC.StaticFriction} {
			if sf = obj_SidFloat(tcn, n); sf != nil {
				*f = *sf
			}
		}
	}
}

func load_Float4x2(xn *xmlx.Node, obj *nga.Float4x2) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_PxRigidConstraintInst(xn *xmlx.Node, obj *nga.PxRigidConstraintInst) {
}

func load_GeometryBrepSurfaces(xn *xmlx.Node, obj *nga.GeometryBrepSurfaces) {
	obj.All = objs_GeometryBrepSurface(xn, "surface")
}

func load_Bool4(xn *xmlx.Node, obj *nga.Bool4) {
	arr_Bools(xn, len(obj), func(i int, b bool) {
		obj[i] = b
	})
}

func load_FxCreate3D(xn *xmlx.Node, obj *nga.FxCreate3D) {
	if c := obj_FxCreate(xn, ""); c != nil {
		obj.FxCreate = *c
	}
	if m := obj_FxCreateMips(xn, "mips"); m != nil {
		obj.Mips = *m
	}
	obj.InitFrom = objs_FxCreate3DInitFrom(xn, "init_from")
	if sn := xcn(xn, "size"); sn != nil {
		obj.Size.Width = xau64(sn, "width")
		obj.Size.Height = xau64(sn, "height")
		obj.Size.Depth = xau64(sn, "depth")
	}
}

func load_GeometryBrepCurves(xn *xmlx.Node, obj *nga.GeometryBrepCurves) {
	obj.All = objs_GeometryBrepCurve(xn, "curve")
}

func load_PxRigidBodyDef(xn *xmlx.Node, obj *nga.PxRigidBodyDef) {
	if tcn := node_TechCommon(xn); tcn != nil {
		if rbc := obj_PxRigidBodyCommon(tcn, ""); rbc != nil {
			obj.TC.PxRigidBodyCommon = *rbc
		}
	}
}

func load_KxArticulatedSystemDef(xn *xmlx.Node, obj *nga.KxArticulatedSystemDef) {
	obj.Kinematics = obj_KxKinematicsSystem(xn, "kinematics")
	obj.Motion = obj_KxMotionSystem(xn, "motion")
}

func load_GeometryBrepPcurves(xn *xmlx.Node, obj *nga.GeometryBrepPcurves) {
	if iiv := obj_IndexedInputs(xn, ""); iiv != nil {
		obj.IndexedInputs = *iiv
	}
}

func load_GeometryBrepHyperbola(xn *xmlx.Node, obj *nga.GeometryBrepHyperbola) {
	if f2 := obj_Float2(xn, "radius"); f2 != nil {
		obj.Radii = *f2
	}
}

func load_PxRigidBodyCommon(xn *xmlx.Node, obj *nga.PxRigidBodyCommon) {
	if sb := obj_SidBool(xn, "dynamic"); sb != nil {
		obj.Dynamic = *sb
	}
	obj.Mass = obj_SidFloat(xn, "mass")
	if tn := xcn(xn, "mass_frame"); tn != nil {
		obj.MassFrame = get_Transforms(tn)
	}
	obj.Inertia = obj_SidFloat3(xn, "inertia")
	obj.Shapes = objs_PxShape(xn, "shape")
	if pm := obj_PxMaterial(xn, ""); pm != nil {
		obj.Material = *pm
	}
}

func load_KxArticulatedSystemInst(xn *xmlx.Node, obj *nga.KxArticulatedSystemInst) {
	obj.Bindings = objs_KxBinding(xn, "bind")
}

func load_Float3x2(xn *xmlx.Node, obj *nga.Float3x2) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_GeometryControlVertices(xn *xmlx.Node, obj *nga.GeometryControlVertices) {
}

func load_Float2x4(xn *xmlx.Node, obj *nga.Float2x4) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_ParamInst(xn *xmlx.Node, obj *nga.ParamInst) {
	obj.Ref.Set(xas(xn, "ref"))
	for _, cn := range xn.Children {
		if cn.Type == xmlx.NT_ELEMENT {
			if cn.Name.Local == "connect_param" {
				obj.IsConnectParamRef = true
				obj.Value = xas(cn, "ref")
				break
			} else if obj.Value = xv(cn); obj.Value != nil {
				break
			}
		}
	}
}

func load_FxTechniqueGlsl(xn *xmlx.Node, obj *nga.FxTechniqueGlsl) {
	if ft := obj_FxTechnique(xn, ""); ft != nil {
		obj.FxTechnique = *ft
	}
	obj.Annotations = objs_FxAnnotation(xn, "annotate")
	obj.Passes = objs_FxPass(xn, "pass")
}

func load_KxLink(xn *xmlx.Node, obj *nga.KxLink) {
	var a *nga.KxAttachment
	obj.Transforms = get_Transforms(xn)
	for _, cn := range xcns(xn, "attachment_full", "attachment_start", "attachment_end") {
		if a = obj_KxAttachment(cn, ""); a != nil {
			obj.Attachments = append(obj.Attachments, a)
		}
	}
}

func load_FxProfileGlsl(xn *xmlx.Node, obj *nga.FxProfileGlsl) {
	var ci *nga.FxProfileGlslCodeInclude
	obj.Platform = xasd(xn, "platform", obj.Platform)
	obj.Techniques = objs_FxTechniqueGlsl(xn, "technique")
	for _, cn := range xcns(xn, "code", "include") {
		if ci = obj_FxProfileGlslCodeInclude(cn, ""); ci != nil {
			obj.CodesIncludes = append(obj.CodesIncludes, *ci)
		}
	}
}

func load_FxCreateCube(xn *xmlx.Node, obj *nga.FxCreateCube) {
	if c := obj_FxCreate(xn, ""); c != nil {
		obj.FxCreate = *c
	}
	if m := obj_FxCreateMips(xn, "mips"); m != nil {
		obj.Mips = *m
	}
	obj.InitFrom = objs_FxCreateCubeInitFrom(xn, "init_from")
	if sn := xcn(xn, "size"); sn != nil {
		obj.Size.Width = xau64(sn, "width")
	}
}

func load_FxCreate3DInitFrom(xn *xmlx.Node, obj *nga.FxCreate3DInitFrom) {
	if ci := obj_FxCreateInitFrom(xn, ""); ci != nil {
		obj.FxCreateInitFrom = *ci
	}
	obj.Depth = xau64(xn, "depth")
}

func load_KxSceneDef(xn *xmlx.Node, obj *nga.KxSceneDef) {
	obj.ArticulatedSystems = objs_KxArticulatedSystemInst(xn, "instance_articulated_systems")
	obj.Models = objs_KxModelInst(xn, "instance_kinematics_model")
}

func load_ControllerInputs(xn *xmlx.Node, obj *nga.ControllerInputs) {
}

func load_IndexedInputs(xn *xmlx.Node, obj *nga.IndexedInputs) {
	obj.Count = xau64(xn, "count")
	obj.Inputs = objs_InputShared(xn, "input")
	obj.Indices = listcn_Uints(xn, "p")
	obj.Vcount = listcn_Ints(xn, "vcount")
}

func load_FxTechnique(xn *xmlx.Node, obj *nga.FxTechnique) {
}

func load_GeometryBrepShells(xn *xmlx.Node, obj *nga.GeometryBrepShells) {
}

func load_LightSpot(xn *xmlx.Node, obj *nga.LightSpot) {
	get_LightColor(xn, &obj.Color)
	if a := obj_LightAttenuation(xn, ""); a != nil {
		obj.Attenuation = *a
	}
	sf := obj_SidFloat(xn, "falloff_exponent")
	if sf != nil {
		obj.Falloff.Exponent = *sf
	}
	if sf = obj_SidFloat(xn, "falloff_angle"); sf != nil {
		obj.Falloff.Angle = *sf
	}
}

func load_FxCreate2D(xn *xmlx.Node, obj *nga.FxCreate2D) {
	if c := obj_FxCreate(xn, ""); c != nil {
		obj.FxCreate = *c
	}
	obj.Mips = obj_FxCreateMips(xn, "mips")
	obj.InitFrom = objs_FxCreateInitFrom(xn, "init_from")
	obj.Size.Exact = obj_FxCreate2DSizeExact(xn, "size_exact")
	obj.Size.Ratio = obj_FxCreate2DSizeRatio(xn, "size_ratio")
	obj.Unnormalized = (xcns(xn, "unnormalized") != nil)
}

func load_ParamOrBool(xn *xmlx.Node, obj *nga.ParamOrBool) {
	obj.Param.Set(xs(xn, "param"))
	obj.B = xb(xn, "bool")
}

func load_GeometryBrepFaces(xn *xmlx.Node, obj *nga.GeometryBrepFaces) {
	if iiv := obj_IndexedInputs(xn, ""); iiv != nil {
		obj.IndexedInputs = *iiv
	}
}

func load_KxAttachment(xn *xmlx.Node, obj *nga.KxAttachment) {
	switch xn.Name.Local {
	case "attachment_full":
		obj.Type = nga.KX_ATTACHMENT_TYPE_FULL
	case "attachment_start":
		obj.Type = nga.KX_ATTACHMENT_TYPE_START
	case "attachment_end":
		obj.Type = nga.KX_ATTACHMENT_TYPE_END
	}
	if obj.Type > 0 {
		obj.Joint.Set(xas(xn, "joint"))
		obj.Transforms = get_Transforms(xn)
		obj.Link = obj_KxLink(xn, "link")
	}
}

func load_AssetContributor(xn *xmlx.Node, obj *nga.AssetContributor) {
	// obj.Author = xs(xn, "author")
	// obj.AuthorEmail = xs(xn, "author_email")
	// obj.AuthorWebsite = xs(xn, "author_website")
	// obj.AuthoringTool = xs(xn, "authoring_tool")
	// obj.Comments = xs(xn, "comments")
	// obj.Copyright = xs(xn, "copyright")
	// obj.SourceData = xs(xn, "source_data")
}

func load_ParamOrRefSid(xn *xmlx.Node, obj *nga.ParamOrRefSid) {
	obj.Param.Set(xs(xn, "param"))
	obj.Sr.Set(xss(xn, "sidref", "SIDREF"))
}

func load_FxSampler(xn *xmlx.Node, obj *nga.FxSampler) {
	switch xn.Name.Local {
	case "sampler1D":
		obj.Type = nga.FX_SAMPLER_TYPE_1D
	case "sampler2D":
		obj.Type = nga.FX_SAMPLER_TYPE_2D
	case "sampler3D":
		obj.Type = nga.FX_SAMPLER_TYPE_3D
	case "samplerCUBE":
		obj.Type = nga.FX_SAMPLER_TYPE_CUBE
	case "samplerDEPTH":
		obj.Type = nga.FX_SAMPLER_TYPE_DEPTH
	case "samplerRECT":
		obj.Type = nga.FX_SAMPLER_TYPE_RECT
	}
	if obj.Type > 0 {
		obj.Image = obj_FxImageInst(xn, "instance_image")
		if ss := obj_FxSamplerStates(xn, ""); ss != nil {
			obj.FxSamplerStates = *ss
		}
	}
}

func load_Scene(xn *xmlx.Node, obj *nga.Scene) {
	obj.Kinematics = obj_KxSceneInst(xn, "instance_kinematics_scene")
	obj.Visual = obj_VisualSceneInst(xn, "instance_visual_scene")
	obj.Physics = objs_PxSceneInst(xn, "instance_physics_scene")
}

func load_FxPassProgramShaderSources(xn *xmlx.Node, obj *nga.FxPassProgramShaderSources) {
}

func load_FxSamplerImage(xn *xmlx.Node, obj *nga.FxSamplerImage) {
	if ii := obj_FxImageInst(xn, ""); ii != nil {
		obj.FxImageInst = *ii
	}
}

func load_FxSamplerStates(xn *xmlx.Node, obj *nga.FxSamplerStates) {
	obj.Filtering = obj_FxSamplerFiltering(xn, "")
	obj.Wrapping = obj_FxSamplerWrapping(xn, "")
}

func load_RefId(xn *xmlx.Node, obj *nga.RefId) {
	panic("load_RefId() should never be called in theory?!")
}

func load_RefParam(xn *xmlx.Node, obj *nga.RefParam) {
	panic("load_RefParam() should never be called in theory?!")
}

func load_RefSid(xn *xmlx.Node, obj *nga.RefSid) {
	panic("load_RefSid() should never be called in theory?!")
}
