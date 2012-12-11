package collimp

import (
	"encoding/hex"
	"strings"

	xmlx "github.com/jteeuwen/go-pkg-xmlx"

	nga "github.com/go3d/go-ngine/assets"
	ustr "github.com/metaleap/go-util/str"
)

func load_FxEffectDef(xn *xmlx.Node, obj *nga.FxEffectDef) {
	obj.Annotations = objs_FxAnnotation(xn, "annotate")
	obj.Profiles.GlSl = objs_FxProfileGlSl(xn, "profile_GLSL")
	obj.Profiles.Common = objs_FxProfileCommon(xn, "profile_COMMON")
}

func load_PxCylinder(xn *xmlx.Node, obj *nga.PxCylinder) {
	obj.Height = xf64(xn, "height")
	if r := obj_Float2(xn, "radius"); r != nil {
		obj.Radius = *r
	}
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

func load_Source(xn *xmlx.Node, obj *nga.Source) {
	if tcn := node_TechCommon(xn); tcn != nil {
		obj.TC.Accessor = obj_SourceAccessor(tcn, "accessor")
	}
	if od, dn := &obj.Data, xcn1(xn, "bool_array", "float_array", "IDREF_array", "int_array", "Name_array", "SIDREF_array", "token_array"); dn != nil {
		has_ID(dn, &od.HasID)
		has_Name(dn, &od.HasName)
		switch dn.Name.Local {
		case "bool_array":
			od.Bools = list_Bools(dn)
		case "float_array":
			od.Floats = list_Floats(dn)
		case "IDREF_array":
			od.IdRefs = list_Strings(dn)
		case "int_array":
			od.Ints = list_Ints(dn)
		case "Name_array":
			od.Names = list_Strings(dn)
		case "SIDREF_array":
			od.SidRefs = list_Strings(dn)
		case "token_array":
			od.Tokens = list_Strings(dn)
		}
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
		obj.Radius = *f
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
	switch strings.ToUpper(xas(xn, "face")) {
	case "POSITIVE_X":
		obj.Face = nga.FX_CUBE_FACE_POSITIVE_X
	case "NEGATIVE_X":
		obj.Face = nga.FX_CUBE_FACE_NEGATIVE_X
	case "POSITIVE_Y":
		obj.Face = nga.FX_CUBE_FACE_POSITIVE_Y
	case "NEGATIVE_Y":
		obj.Face = nga.FX_CUBE_FACE_NEGATIVE_Y
	case "POSITIVE_Z":
		obj.Face = nga.FX_CUBE_FACE_POSITIVE_Z
	case "NEGATIVE_Z":
		obj.Face = nga.FX_CUBE_FACE_NEGATIVE_Z
	}
}

func load_KxModelInst(xn *xmlx.Node, obj *nga.KxModelInst) {
	obj.Bindings = objs_KxBind(xn, "bind")
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
	obj.Source = xas(xn, "source")
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

func load_KxArticulatedSystemKinematicsAxis(xn *xmlx.Node, obj *nga.KxArticulatedSystemKinematicsAxis) {
	obj.JointAxis = xas(xn, "axis")
	obj.Formulas.Defs = objs_FormulaDef(xn, "formula")
	obj.Formulas.Insts = objs_FormulaInst(xn, "instance_formula")
	obj.Limits = obj_KxArticulatedSystemAxisLimits(xn, "limits")
	obj.Indices = objs_KxArticulatedSystemAxisIndex(xn, "index")
	if l := obj_ParamBool(xn, "locked"); l != nil {
		obj.Locked = *l
	}
	if pb := obj_ParamBool(xn, "active"); pb != nil {
		obj.Active = *pb
	} else {
		obj.Active.B = true
	}
}

func load_PxModelInst(xn *xmlx.Node, obj *nga.PxModelInst) {
	obj.ParentRef = xas(xn, "parent")
	obj.ForceFields = objs_PxForceFieldInst(xn, "instance_force_field")
	obj.RigidBodies = objs_PxRigidBodyInst(xn, "instance_rigid_body")
	obj.RigidConstraints = objs_PxRigidConstraintInst(xn, "instance_rigid_constraint")
}

func load_Transform(xn *xmlx.Node, obj *nga.Transform) {
	obj.F = list_Floats(xn)
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
}

func load_KxJoint(xn *xmlx.Node, obj *nga.KxJoint) {
	if xa := xcn(xn, "axis"); xa != nil {
		has_Sid(xa, &obj.Axis.HasSid)
		has_Name(xa, &obj.Axis.HasName)
		if f := obj_Float3(xa, ""); f != nil {
			obj.Axis.F = *f
		}
	}
	obj.Limits = obj_KxJointLimits(xn, "limits")
	switch xn.Name.Local {
	case "revolute":
		obj.Type = nga.KX_JOINT_TYPE_REVOLUTE
	case "prismatic":
		obj.Type = nga.KX_JOINT_TYPE_PRISMATIC
	}
}

func load_ParamInt(xn *xmlx.Node, obj *nga.ParamInt) {
	obj.Ref = xs(xn, "param")
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
	obj.BindMaterial = obj_BindMaterial(xn, "bind_material")
	obj.Skeletons = list_StringsN(xn, "skeleton")
}

func load_GeometryBrepCurve(xn *xmlx.Node, obj *nga.GeometryBrepCurve) {
	obj.Orientations = objs_GeometryBrepOrientation(xn, "orient")
	obj.Origin = xv3(xn, "origin")
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

func load_FormulaDef(xn *xmlx.Node, obj *nga.FormulaDef) {
	if pf := obj_ParamFloat(xn, "target"); pf != nil {
		obj.Target = *pf
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
	obj.Created = xs(xn, "created")
	obj.Keywords = xs(xn, "keywords")
	obj.Modified = xs(xn, "modified")
	obj.Revision = xs(xn, "revision")
	obj.Subject = xs(xn, "subject")
	obj.Title = xs(xn, "title")
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
	case "TESSELATION":
		obj.Stage = nga.FX_PASS_PROGRAM_SHADER_STAGE_TESSELATION
	case "VERTEX":
		obj.Stage = nga.FX_PASS_PROGRAM_SHADER_STAGE_VERTEX
	}
	if sn := xcn(xn, "sources"); sn != nil {
		pss := nga.FxPassProgramShaderSources{}
		arr := make([]nga.FxPassProgramShaderSources, 0, len(sn.Children))
		for _, scn := range sn.Children {
			if pss.IsImportRef = (scn.Name.Local == "import"); pss.IsImportRef {
				pss.S = xas(scn, "ref")
			} else {
				pss.S = scn.Value
			}
			arr = append(arr, pss)
		}
		obj.Sources = arr
	}
}

func load_GeometryBrepCapsule(xn *xmlx.Node, obj *nga.GeometryBrepCapsule) {
	obj.Height = xf64(xn, "height")
	if f3 := obj_Float3(xn, "radius"); f3 != nil {
		obj.Radius = *f3
	}
}

func load_FxProfileCommon(xn *xmlx.Node, obj *nga.FxProfileCommon) {
	if fp := obj_FxProfile(xn, ""); fp != nil {
		obj.FxProfile = *fp
	}
	if tc := obj_FxTechniqueCommon(xn, "technique"); tc != nil {
		obj.Technique = *tc
	}
}

func load_ControllerMorph(xn *xmlx.Node, obj *nga.ControllerMorph) {
	obj.Source = xas(xn, "source")
	obj.Relative = (strings.ToUpper(xas(xn, "method")) == "RELATIVE")
	if t := obj_ControllerInputs(xn, "targets"); t != nil {
		obj.Targets = *t
	}
}

func load_FxMaterialInstBindVertexInput(xn *xmlx.Node, obj *nga.FxMaterialInstBindVertexInput) {
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
		obj.TC.Gravity = obj_ScopedVec3(tcn, "gravity")
		obj.TC.TimeStep = obj_ScopedFloat(tcn, "time_step")
	}
}

func load_InputShared(xn *xmlx.Node, obj *nga.InputShared) {
	if in := obj_Input(xn, ""); in != nil {
		obj.Input = *in
	}
	obj.Offset = xau64(xn, "offset")
	obj.Set = xau64p(xn, "set")
}

func load_ScopedBool(xn *xmlx.Node, obj *nga.ScopedBool) {
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
	obj.ConvexHullOf = xas(xn, "convex_hull_of")
	for _, pn := range xcns(xn, "lines", "linestrips", "polygons", "polylist", "triangles", "trifans", "tristrips") {
		if p := obj_GeometryPrimitives(pn, ""); p != nil {
			obj.Primitives = append(obj.Primitives, p)
		}
	}
	obj.Vertices = obj_GeometryVertices(xn, "vertices")
}

func load_BindMaterial(xn *xmlx.Node, obj *nga.BindMaterial) {
	obj.Params = objs_Param(xn, "param")
	if tcn := node_TechCommon(xn); tcn != nil {
		obj.TC.Materials = objs_FxMaterialInst(tcn, "instance_material")
	}
}

func load_FxEffectInst(xn *xmlx.Node, obj *nga.FxEffectInst) {
	obj.TechniqueHints = objs_FxEffectInstTechniqueHint(xn, "technique_hint")
}

func load_KxJointLimits(xn *xmlx.Node, obj *nga.KxJointLimits) {
	obj.Max = obj_ScopedFloat(xn, "max")
	obj.Min = obj_ScopedFloat(xn, "min")
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
	obj.Count = xau64(xn, "count")
	if iiv := obj_IndexedInputsV(xn, ""); iiv != nil {
		obj.IndexedInputsV = *iiv
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
			sf := obj_ScopedFloat(cn, "aspect_ratio")
			obj.TC.AspectRatio = sf
			if sf = obj_ScopedFloat(cn, "zfar"); sf != nil {
				obj.TC.Zfar = *sf
			}
			if sf = obj_ScopedFloat(cn, "znear"); sf != nil {
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
	obj.TargetNode = xas(xn, "target")
	if tcn := node_TechCommon(xn); tcn != nil {
		if rbc := obj_PxRigidBodyCommon(tcn, ""); rbc != nil {
			obj.TC.PxRigidBodyCommon = *rbc
		}
		v3 := xv3(tcn, "angular_velocity")
		if v3 != nil {
			obj.TC.AngularVelocity = *v3
		}
		if v3 = xv3(tcn, "velocity"); v3 != nil {
			obj.TC.Velocity = *v3
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
		obj.Channels = nga.FX_CREATE_FORMAT_HINT_CHANNELS_LUM
	case "LA":
		obj.Channels = nga.FX_CREATE_FORMAT_HINT_CHANNELS_LUMA
	case "D":
		obj.Channels = nga.FX_CREATE_FORMAT_HINT_CHANNELS_DEPTH
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

func load_KxArticulatedSystemKinematicsFrame(xn *xmlx.Node, obj *nga.KxArticulatedSystemKinematicsFrame) {
	obj.Link = xas(xn, "link")
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
	for _, cn := range xn.Children {
		if t := obj_Transform(cn, ""); t != nil {
			obj.Transforms = append(obj.Transforms, t)
		}
	}
}

func load_KxJointInst(xn *xmlx.Node, obj *nga.KxJointInst) {
}

func load_Input(xn *xmlx.Node, obj *nga.Input) {
	obj.Semantic = xas(xn, "semantic")
	obj.Source = xas(xn, "source")
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
	obj.IndexOfRefraction = obj_ParamScopedFloat(xn, "index_of_refraction")
	obj.Reflective = obj_FxColorOrTexture(xn, "reflective")
	obj.Reflectivity = obj_ParamScopedFloat(xn, "reflectivity")
	obj.Transparency = obj_ParamScopedFloat(xn, "transparency")
	obj.Transparent = obj_FxColorOrTexture(xn, "transparent")
}

func load_PxRigidConstraintLimit(xn *xmlx.Node, obj *nga.PxRigidConstraintLimit) {
	sf3 := obj_ScopedFloat3(xn, "max")
	if sf3 != nil {
		obj.Max = *sf3
	}
	if sf3 = obj_ScopedFloat3(xn, "min"); sf3 != nil {
		obj.Min = *sf3
	}
}

func load_ScopedVec3(xn *xmlx.Node, obj *nga.ScopedVec3) {
	if v3 := xv3(xn, ""); v3 != nil {
		obj.V = *v3
	}
}

func load_NodeInst(xn *xmlx.Node, obj *nga.NodeInst) {
	obj.ProxyRef = xas(xn, "proxy")
}

func load_CameraOrthographic(xn *xmlx.Node, obj *nga.CameraOrthographic) {
	obj.MagX = obj_ScopedFloat(xn, "xmag")
	obj.MagY = obj_ScopedFloat(xn, "ymag")
}

func load_Float3x3(xn *xmlx.Node, obj *nga.Float3x3) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_FxPassState(xn *xmlx.Node, obj *nga.FxPassState) {
	obj.Value = xas(xn, "value")
	obj.ParamRef = xas(xn, "param")
	obj.Index = xf64(xn, "index")
}

func load_GeometryBrepSurface(xn *xmlx.Node, obj *nga.GeometryBrepSurface) {

}

func load_AnimationSampler(xn *xmlx.Node, obj *nga.AnimationSampler) {

}

func load_ParamUint(xn *xmlx.Node, obj *nga.ParamUint) {
	obj.Ref = xs(xn, "param")
	obj.U = xu64(xn, "uint")
}

func load_ParamDefs(xn *xmlx.Node, obj *nga.ParamDefs) {

}

func load_AnimationChannel(xn *xmlx.Node, obj *nga.AnimationChannel) {

}

func load_LightAmbient(xn *xmlx.Node, obj *nga.LightAmbient) {

}

func load_CameraPerspective(xn *xmlx.Node, obj *nga.CameraPerspective) {

}

func load_PxForceFieldDef(xn *xmlx.Node, obj *nga.PxForceFieldDef) {

}

func load_PxRigidConstraintSpring(xn *xmlx.Node, obj *nga.PxRigidConstraintSpring) {

}

func load_GeometryBrepNurbsSurface(xn *xmlx.Node, obj *nga.GeometryBrepNurbsSurface) {

}

func load_CameraDef(xn *xmlx.Node, obj *nga.CameraDef) {

}

func load_GeometryBrepSweptSurface(xn *xmlx.Node, obj *nga.GeometryBrepSweptSurface) {

}

func load_FxPassProgramBindUniform(xn *xmlx.Node, obj *nga.FxPassProgramBindUniform) {

}

func load_KxArticulatedSystemMotionAxis(xn *xmlx.Node, obj *nga.KxArticulatedSystemMotionAxis) {

}

func load_AnimationClipDef(xn *xmlx.Node, obj *nga.AnimationClipDef) {

}

func load_Float7(xn *xmlx.Node, obj *nga.Float7) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_AnimationDef(xn *xmlx.Node, obj *nga.AnimationDef) {

}

func load_PxShape(xn *xmlx.Node, obj *nga.PxShape) {

}

func load_LightPoint(xn *xmlx.Node, obj *nga.LightPoint) {

}

func load_GeometrySpline(xn *xmlx.Node, obj *nga.GeometrySpline) {

}

func load_GeometryInst(xn *xmlx.Node, obj *nga.GeometryInst) {

}

func load_PxRigidConstraintDef(xn *xmlx.Node, obj *nga.PxRigidConstraintDef) {

}

func load_GeometryBrepSurfaceCurves(xn *xmlx.Node, obj *nga.GeometryBrepSurfaceCurves) {

}

func load_Int3x3(xn *xmlx.Node, obj *nga.Int3x3) {
	arr_Ints(xn, len(obj), func(i int, n int64) {
		obj[i] = n
	})
}

func load_FxTechniqueCommonBlinn(xn *xmlx.Node, obj *nga.FxTechniqueCommonBlinn) {

}

func load_ScopedString(xn *xmlx.Node, obj *nga.ScopedString) {

}

func load_GeometryDef(xn *xmlx.Node, obj *nga.GeometryDef) {

}

func load_FxCreate2DSizeExact(xn *xmlx.Node, obj *nga.FxCreate2DSizeExact) {

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

func load_KxBind(xn *xmlx.Node, obj *nga.KxBind) {

}

func load_FxPassEvaluationTarget(xn *xmlx.Node, obj *nga.FxPassEvaluationTarget) {

}

func load_FxColorOrTexture(xn *xmlx.Node, obj *nga.FxColorOrTexture) {

}

func load_GeometryVertices(xn *xmlx.Node, obj *nga.GeometryVertices) {

}

func load_KxArticulatedSystemMotion(xn *xmlx.Node, obj *nga.KxArticulatedSystemMotion) {

}

func load_Int4x4(xn *xmlx.Node, obj *nga.Int4x4) {
	arr_Ints(xn, len(obj), func(i int, n int64) {
		obj[i] = n
	})
}

func load_CameraImager(xn *xmlx.Node, obj *nga.CameraImager) {

}

func load_FxPassEvaluationClearDepth(xn *xmlx.Node, obj *nga.FxPassEvaluationClearDepth) {

}

func load_GeometryBrepEllipse(xn *xmlx.Node, obj *nga.GeometryBrepEllipse) {

}

func load_FxMaterialInstBind(xn *xmlx.Node, obj *nga.FxMaterialInstBind) {

}

func load_GeometryBrepParabola(xn *xmlx.Node, obj *nga.GeometryBrepParabola) {

}

func load_ScopedFloat(xn *xmlx.Node, obj *nga.ScopedFloat) {

}

func load_Bool3(xn *xmlx.Node, obj *nga.Bool3) {
	arr_Bools(xn, len(obj), func(i int, b bool) {
		obj[i] = b
	})
}

func load_VisualSceneDef(xn *xmlx.Node, obj *nga.VisualSceneDef) {

}

func load_KxSceneInst(xn *xmlx.Node, obj *nga.KxSceneInst) {

}

func load_LightDirectional(xn *xmlx.Node, obj *nga.LightDirectional) {

}

func load_KxArticulatedSystemAxisIndex(xn *xmlx.Node, obj *nga.KxArticulatedSystemAxisIndex) {

}

func load_VisualSceneRendering(xn *xmlx.Node, obj *nga.VisualSceneRendering) {

}

func load_GeometryPolygonHole(xn *xmlx.Node, obj *nga.GeometryPolygonHole) {

}

func load_Param(xn *xmlx.Node, obj *nga.Param) {

}

func load_FxTechniqueCommonLambert(xn *xmlx.Node, obj *nga.FxTechniqueCommonLambert) {

}

func load_ParamScopedFloat(xn *xmlx.Node, obj *nga.ParamScopedFloat) {
	obj.Ref = xs(xn, "param")
	if f := obj_ScopedFloat(xn, "float"); f != nil {
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

}

func load_KxSceneInstBindModel(xn *xmlx.Node, obj *nga.KxSceneInstBindModel) {

}

func load_GeometryBrepSolids(xn *xmlx.Node, obj *nga.GeometryBrepSolids) {

}

func load_FxMaterialDef(xn *xmlx.Node, obj *nga.FxMaterialDef) {

}

func load_ParamFloat(xn *xmlx.Node, obj *nga.ParamFloat) {
	obj.Ref = xs(xn, "param")
	obj.F = xf64(xn, "float")
}

func load_ScopedFloat3(xn *xmlx.Node, obj *nga.ScopedFloat3) {

}

func load_FxInitFrom(xn *xmlx.Node, obj *nga.FxInitFrom) {
	obj.RefUrl = xs(xn, "ref")
	if hn := xcn(xn, "hex"); hn != nil {
		obj.Raw.Format = xas(hn, "format")
		obj.Raw.Data, _ = hex.DecodeString(xs(hn, ""))
	}
}

func load_FxTexture(xn *xmlx.Node, obj *nga.FxTexture) {

}

func load_Int4(xn *xmlx.Node, obj *nga.Int4) {
	arr_Ints(xn, len(obj), func(i int, n int64) {
		obj[i] = n
	})
}

func load_GeometryBrepOrientation(xn *xmlx.Node, obj *nga.GeometryBrepOrientation) {

}

func load_GeometryPrimitives(xn *xmlx.Node, obj *nga.GeometryPrimitives) {

}

func load_Sources(xn *xmlx.Node, obj *nga.Sources) {

}

func load_PxRigidConstraintAttachment(xn *xmlx.Node, obj *nga.PxRigidConstraintAttachment) {

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

}

func load_Float3x4(xn *xmlx.Node, obj *nga.Float3x4) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_KxSceneInstBindJointAxis(xn *xmlx.Node, obj *nga.KxSceneInstBindJointAxis) {

}

func load_KxJointDef(xn *xmlx.Node, obj *nga.KxJointDef) {

}

func load_VisualSceneEvaluation(xn *xmlx.Node, obj *nga.VisualSceneEvaluation) {

}

func load_FxProfileGlSlCodeInclude(xn *xmlx.Node, obj *nga.FxProfileGlSlCodeInclude) {

}

func load_AssetGeographicLocation(xn *xmlx.Node, obj *nga.AssetGeographicLocation) {

}

func load_KxArticulatedSystemEffector(xn *xmlx.Node, obj *nga.KxArticulatedSystemEffector) {

}

func load_GeometryBrepTorus(xn *xmlx.Node, obj *nga.GeometryBrepTorus) {

}

func load_Layers(xn *xmlx.Node, obj *nga.Layers) {

}

func load_FxProfile(xn *xmlx.Node, obj *nga.FxProfile) {

}

func load_KxArticulatedSystemAxisLimits(xn *xmlx.Node, obj *nga.KxArticulatedSystemAxisLimits) {

}

func load_FxMaterialInst(xn *xmlx.Node, obj *nga.FxMaterialInst) {

}

func load_GeometryBrepCone(xn *xmlx.Node, obj *nga.GeometryBrepCone) {

}

func load_Extra(xn *xmlx.Node, obj *nga.Extra) {

}

func load_ControllerSkin(xn *xmlx.Node, obj *nga.ControllerSkin) {

}

func load_Bool2(xn *xmlx.Node, obj *nga.Bool2) {
	arr_Bools(xn, len(obj), func(i int, b bool) {
		obj[i] = b
	})
}

func load_VisualSceneRenderingMaterialInst(xn *xmlx.Node, obj *nga.VisualSceneRenderingMaterialInst) {

}

func load_NodeDef(xn *xmlx.Node, obj *nga.NodeDef) {

}

func load_PxForceFieldInst(xn *xmlx.Node, obj *nga.PxForceFieldInst) {

}

func load_KxArticulatedSystemKinematics(xn *xmlx.Node, obj *nga.KxArticulatedSystemKinematics) {

}

func load_GeometryBrepLine(xn *xmlx.Node, obj *nga.GeometryBrepLine) {

}

func load_GeometryBrep(xn *xmlx.Node, obj *nga.GeometryBrep) {

}

func load_KxModelDef(xn *xmlx.Node, obj *nga.KxModelDef) {

}

func load_FxImageInst(xn *xmlx.Node, obj *nga.FxImageInst) {

}

func load_GeometryBrepBox(xn *xmlx.Node, obj *nga.GeometryBrepBox) {

}

func load_ParamFloat2(xn *xmlx.Node, obj *nga.ParamFloat2) {
	obj.Ref = xs(xn, "param")
	if fn := xcn(xn, "float2"); fn != nil {
		arr_Floats(fn, 2, func(i int, f float64) {
			obj.F[i] = f
		})
	}
}

func load_Technique(xn *xmlx.Node, obj *nga.Technique) {

}

func load_FxPassEvaluationClearStencil(xn *xmlx.Node, obj *nga.FxPassEvaluationClearStencil) {

}

func load_FxCreate2DSizeRatio(xn *xmlx.Node, obj *nga.FxCreate2DSizeRatio) {

}

func load_GeometryBrepPlane(xn *xmlx.Node, obj *nga.GeometryBrepPlane) {

}

func load_FxPassProgramBindAttribute(xn *xmlx.Node, obj *nga.FxPassProgramBindAttribute) {

}

func load_FxParamDefs(xn *xmlx.Node, obj *nga.FxParamDefs) {

}

func load_PxMaterialDef(xn *xmlx.Node, obj *nga.PxMaterialDef) {

}

func load_Float4x2(xn *xmlx.Node, obj *nga.Float4x2) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_PxRigidConstraintInst(xn *xmlx.Node, obj *nga.PxRigidConstraintInst) {

}

func load_GeometryBrepSurfaces(xn *xmlx.Node, obj *nga.GeometryBrepSurfaces) {

}

func load_Bool4(xn *xmlx.Node, obj *nga.Bool4) {
	arr_Bools(xn, len(obj), func(i int, b bool) {
		obj[i] = b
	})
}

func load_FxCreate3D(xn *xmlx.Node, obj *nga.FxCreate3D) {

}

func load_GeometryBrepCurves(xn *xmlx.Node, obj *nga.GeometryBrepCurves) {

}

func load_PxRigidBodyDef(xn *xmlx.Node, obj *nga.PxRigidBodyDef) {

}

func load_KxArticulatedSystemDef(xn *xmlx.Node, obj *nga.KxArticulatedSystemDef) {

}

func load_GeometryBrepPcurves(xn *xmlx.Node, obj *nga.GeometryBrepPcurves) {

}

func load_GeometryBrepHyperbola(xn *xmlx.Node, obj *nga.GeometryBrepHyperbola) {

}

func load_PxRigidBodyCommon(xn *xmlx.Node, obj *nga.PxRigidBodyCommon) {

}

func load_KxArticulatedSystemInst(xn *xmlx.Node, obj *nga.KxArticulatedSystemInst) {

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

func load_FxPassProgramShaderSources(xn *xmlx.Node, obj *nga.FxPassProgramShaderSources) {

}

func load_ParamInst(xn *xmlx.Node, obj *nga.ParamInst) {

}

func load_FxTechniqueGlsl(xn *xmlx.Node, obj *nga.FxTechniqueGlsl) {

}

func load_KxLink(xn *xmlx.Node, obj *nga.KxLink) {

}

func load_FxProfileGlSl(xn *xmlx.Node, obj *nga.FxProfileGlSl) {

}

func load_FxCreateCube(xn *xmlx.Node, obj *nga.FxCreateCube) {

}

func load_FxCreate3DInitFrom(xn *xmlx.Node, obj *nga.FxCreate3DInitFrom) {

}

func load_KxSceneDef(xn *xmlx.Node, obj *nga.KxSceneDef) {

}

func load_ControllerInputs(xn *xmlx.Node, obj *nga.ControllerInputs) {
}

func load_IndexedInputs(xn *xmlx.Node, obj *nga.IndexedInputs) {
	obj.Count = xau64(xn, "count")
	obj.Inputs = objs_InputShared(xn, "input")
	obj.Indices = listcn_Ints(xn, "p")
}

func load_IndexedInputsV(xn *xmlx.Node, obj *nga.IndexedInputsV) {
	if ii := obj_IndexedInputs(xn, ""); ii != nil {
		obj.IndexedInputs = *ii
	}
	obj.Vcount = listcn_Ints(xn, "vcount")
}

func load_FxTechnique(xn *xmlx.Node, obj *nga.FxTechnique) {

}

func load_GeometryBrepShells(xn *xmlx.Node, obj *nga.GeometryBrepShells) {

}

func load_LightSpot(xn *xmlx.Node, obj *nga.LightSpot) {

}

func load_FxCreate2D(xn *xmlx.Node, obj *nga.FxCreate2D) {

}

func load_ParamBool(xn *xmlx.Node, obj *nga.ParamBool) {
	obj.Ref = xs(xn, "param")
	obj.B = xb(xn, "bool")
}

func load_GeometryBrepFaces(xn *xmlx.Node, obj *nga.GeometryBrepFaces) {

}

func load_KxAttachment(xn *xmlx.Node, obj *nga.KxAttachment) {

}

func load_AssetContributor(xn *xmlx.Node, obj *nga.AssetContributor) {

}

func load_ParamSidRef(xn *xmlx.Node, obj *nga.ParamSidRef) {
	obj.Ref = xs(xn, "param")
	obj.S = xss(xn, "sidref", "SIDREF")
}

func load_FxSampler(xn *xmlx.Node, obj *nga.FxSampler) {

}

func load_Scene(xn *xmlx.Node, obj *nga.Scene) {

}
