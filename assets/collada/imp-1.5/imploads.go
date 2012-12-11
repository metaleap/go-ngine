package collimp

import (
	"encoding/hex"

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
	if tc, tcn := &obj.TechniqueCommon, node_TechCommon(xn); tcn != nil {
		tc.Ambient = obj_LightAmbient(tcn, "ambient")
		tc.Directional = obj_LightDirectional(tcn, "directional")
		tc.Point = obj_LightPoint(tcn, "point")
		tc.Spot = obj_LightSpot(tcn, "spot")
	}
}

func load_Source(xn *xmlx.Node, obj *nga.Source) {
	if tc, tcn := &obj.TechniqueCommon, node_TechCommon(xn); tcn != nil {
		tc.Accessor = obj_SourceAccessor(tcn, "accessor")
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
	switch xa(xn, "face") {
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
	obj.Source = xa(xn, "source")
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
	obj.Platform = xa(xn, "platform")
	obj.Ref = xa(xn, "ref")
	obj.Profile = xa(xn, "profile")
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
		switch xs(xn, n) {
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
	obj.JointAxis = xa(xn, "axis")
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
	obj.ParentRef = xa(xn, "parent")
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

}

func load_GeometryBrepCurve(xn *xmlx.Node, obj *nga.GeometryBrepCurve) {

}

func load_GeometryBrepNurbs(xn *xmlx.Node, obj *nga.GeometryBrepNurbs) {

}

func load_FormulaDef(xn *xmlx.Node, obj *nga.FormulaDef) {

}

func load_LightInst(xn *xmlx.Node, obj *nga.LightInst) {

}

func load_Asset(xn *xmlx.Node, obj *nga.Asset) {

}

func load_FxPassProgramShader(xn *xmlx.Node, obj *nga.FxPassProgramShader) {

}

func load_GeometryBrepCapsule(xn *xmlx.Node, obj *nga.GeometryBrepCapsule) {

}

func load_FxProfileCommon(xn *xmlx.Node, obj *nga.FxProfileCommon) {

}

func load_ControllerMorph(xn *xmlx.Node, obj *nga.ControllerMorph) {

}

func load_FxMaterialInstBindVertexInput(xn *xmlx.Node, obj *nga.FxMaterialInstBindVertexInput) {
	obj.InputSemantic, obj.Semantic = xa(xn, "input_semantic"), xa(xn, "semantic")
	obj.InputSet = xau64p(xn, "input_set")
}

func load_Float3(xn *xmlx.Node, obj *nga.Float3) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_PxSceneDef(xn *xmlx.Node, obj *nga.PxSceneDef) {

}

func load_InputShared(xn *xmlx.Node, obj *nga.InputShared) {

}

func load_ScopedBool(xn *xmlx.Node, obj *nga.ScopedBool) {

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

}

func load_GeometryMesh(xn *xmlx.Node, obj *nga.GeometryMesh) {

}

func load_BindMaterial(xn *xmlx.Node, obj *nga.BindMaterial) {

}

func load_FxEffectInst(xn *xmlx.Node, obj *nga.FxEffectInst) {

}

func load_CameraCommon(xn *xmlx.Node, obj *nga.CameraCommon) {

}

func load_KxJointLimits(xn *xmlx.Node, obj *nga.KxJointLimits) {

}

func load_FxSamplerWrapping(xn *xmlx.Node, obj *nga.FxSamplerWrapping) {

}

func load_AnimationInst(xn *xmlx.Node, obj *nga.AnimationInst) {

}

func load_GeometryBrepWires(xn *xmlx.Node, obj *nga.GeometryBrepWires) {

}

func load_FxCreateFormat(xn *xmlx.Node, obj *nga.FxCreateFormat) {

}

func load_FormulaInst(xn *xmlx.Node, obj *nga.FormulaInst) {

}

func load_CameraOptics(xn *xmlx.Node, obj *nga.CameraOptics) {

}

func load_PxRigidBodyInst(xn *xmlx.Node, obj *nga.PxRigidBodyInst) {

}

func load_FxAnnotation(xn *xmlx.Node, obj *nga.FxAnnotation) {

}

func load_FxCreateFormatHint(xn *xmlx.Node, obj *nga.FxCreateFormatHint) {

}

func load_GeometryBrepCircle(xn *xmlx.Node, obj *nga.GeometryBrepCircle) {

}

func load_PxSceneInst(xn *xmlx.Node, obj *nga.PxSceneInst) {

}

func load_FxParamDef(xn *xmlx.Node, obj *nga.FxParamDef) {

}

func load_FxCreateMips(xn *xmlx.Node, obj *nga.FxCreateMips) {

}

func load_KxArticulatedSystemKinematicsFrame(xn *xmlx.Node, obj *nga.KxArticulatedSystemKinematicsFrame) {

}

func load_KxJointInst(xn *xmlx.Node, obj *nga.KxJointInst) {

}

func load_Input(xn *xmlx.Node, obj *nga.Input) {

}

func load_ControllerDef(xn *xmlx.Node, obj *nga.ControllerDef) {

}

func load_FxCreateCommon(xn *xmlx.Node, obj *nga.FxCreateCommon) {

}

func load_FxTechniqueCommonConstant(xn *xmlx.Node, obj *nga.FxTechniqueCommonConstant) {

}

func load_PxRigidConstraintLimit(xn *xmlx.Node, obj *nga.PxRigidConstraintLimit) {

}

func load_ScopedVec3(xn *xmlx.Node, obj *nga.ScopedVec3) {

}

func load_NodeInst(xn *xmlx.Node, obj *nga.NodeInst) {

}

func load_CameraOrthographic(xn *xmlx.Node, obj *nga.CameraOrthographic) {

}

func load_Float3x3(xn *xmlx.Node, obj *nga.Float3x3) {
	arr_Floats(xn, len(obj), func(i int, f float64) {
		obj[i] = f
	})
}

func load_FxPassState(xn *xmlx.Node, obj *nga.FxPassState) {

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
		obj.Raw.Format = xa(hn, "format")
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

func load_IndexedInputs(xn *xmlx.Node, obj *nga.IndexedInputs) {

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
