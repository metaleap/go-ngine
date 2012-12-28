package assets

func (me *LightSpot) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if val = me.Attenuation.resolveSidPath(path, bag); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.Falloff.Angle, nil, me.Falloff.Angle.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.Falloff.Exponent, nil, me.Falloff.Exponent.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	return
}

func (me *PxMaterialDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.DynamicFriction, nil, me.TC.DynamicFriction.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.Restitution, nil, me.TC.Restitution.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.StaticFriction, nil, me.TC.StaticFriction.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	return
}

func (me *PxMaterialDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *KxKinematicsAxis) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, subItem := range me.Formulas {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *GeometryBrepCurves) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.All {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxProfile) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if me.Common != nil {
		if val = me.Common.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Glsl != nil {
		if val = me.Glsl.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxProfile) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *KxModelInst) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *ControllerDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Morph != nil {
		if val = me.Morph.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Skin != nil {
		if val = me.Skin.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *ControllerDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *PxRigidConstraintDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if val = me.RefAttachment.resolveSidPath(path, bag); val != nil {
		return
	}
	if val = me.Attachment.resolveSidPath(path, bag); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.Enabled, nil, me.TC.Enabled.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.Interpenetrate, nil, me.TC.Interpenetrate.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	if me.TC.Limits.Angular != nil {
		if val = me.TC.Limits.Angular.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.TC.Limits.Linear != nil {
		if val = me.TC.Limits.Linear.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.TC.Spring.Angular != nil {
		if val = me.TC.Spring.Angular.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.TC.Spring.Linear != nil {
		if val = me.TC.Spring.Linear.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxRigidConstraintDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *GeometryBrep) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Curves != nil {
		if val = me.Curves.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.SurfaceCurves != nil {
		if val = me.SurfaceCurves.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Surfaces != nil {
		if val = me.Surfaces.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxSceneInst) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxRigidConstraintLimit) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Min, nil, me.Min.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.Max, nil, me.Max.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	return
}

func (me *GeometryBrepCurve) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Element.Nurbs != nil {
		if val = me.Element.Nurbs.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *MaterialBinding) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Params {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.TC.Materials {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxAttachment) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if me.Link != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Link, me.Link, me.Link.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *VisualSceneEvaluation) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.RenderPasses {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *VisualSceneEvaluation) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *SourceAccessor) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Params {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxEffectDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, subItem := range me.Profiles {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxEffectDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *ChildNode) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Inst != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Inst, nil, me.Inst.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxMaterialDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Effect, nil, me.Effect.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	return
}

func (me *FxMaterialDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *FxSampler) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Image != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Image, nil, me.Image.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxShape) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Hollow, nil, me.Hollow.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	if me.Mass != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Mass, nil, me.Mass.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if me.Density != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Density, nil, me.Density.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if val = me.Material.resolveSidPath(path, bag); val != nil {
		return
	}
	if me.Geometry.Inst != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Geometry.Inst, me.Geometry.Inst, me.Geometry.Inst.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *CameraOrthographic) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.MagX != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.MagX, nil, me.MagX.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if me.MagY != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.MagY, nil, me.MagY.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *Document) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Scene != nil {
		if val = me.Scene.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *ControllerMorph) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *ParamOrSidFloat) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.F, nil, me.F.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	return
}

func (me *PxRigidBodyCommon) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Dynamic, nil, me.Dynamic.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	if me.Mass != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Mass, nil, me.Mass.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.MassFrame {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if me.Inertia != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Inertia, nil, me.Inertia.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if val = me.Material.resolveSidPath(path, bag); val != nil {
		return
	}
	for _, subItem := range me.Shapes {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxJointLimits) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Min != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Min, nil, me.Min.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if me.Max != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Max, nil, me.Max.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *GeometrySpline) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxRigidConstraintSpring) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Stiffness, nil, me.Stiffness.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.Damping, nil, me.Damping.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.TargetValue, nil, me.TargetValue.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	return
}

func (me *KxEffector) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxFrameTip) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxRigidBodyInst) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.Dynamic, nil, me.TC.Dynamic.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	if me.TC.Mass != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.Mass, nil, me.TC.Mass.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.TC.MassFrame {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if me.TC.Inertia != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.Inertia, nil, me.TC.Inertia.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if val = me.TC.Material.resolveSidPath(path, bag); val != nil {
		return
	}
	for _, subItem := range me.TC.Shapes {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *LightDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.TC.Point != nil {
		if val = me.TC.Point.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.TC.Spot != nil {
		if val = me.TC.Spot.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *LightDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *FxPassEvaluationTarget) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Image != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Image, nil, me.Image.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *VisualSceneDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Nodes {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.Evaluations {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *VisualSceneDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *CameraOptics) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.TC.AspectRatio != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.AspectRatio, nil, me.TC.AspectRatio.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.Zfar, nil, me.TC.Zfar.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.Znear, nil, me.TC.Znear.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	if me.TC.Orthographic != nil {
		if val = me.TC.Orthographic.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.TC.Perspective != nil {
		if val = me.TC.Perspective.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxLink) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, subItem := range me.Attachments {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxArticulatedSystemDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Kinematics != nil {
		if val = me.Kinematics.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Motion != nil {
		if val = me.Motion.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxArticulatedSystemDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *Scene) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Visual != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Visual, nil, me.Visual.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if me.Kinematics != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Kinematics, me.Kinematics, me.Kinematics.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.Physics {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxMaterial) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Inst != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Inst, nil, me.Inst.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *GeometryBrepSurface) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Element.NurbsSurface != nil {
		if val = me.Element.NurbsSurface.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Element.SweptSurface != nil {
		if val = me.Element.SweptSurface.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *GeometryInst) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.MaterialBinding != nil {
		if val = me.MaterialBinding.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *GeometryBrepSweptSurface) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Curve != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Curve, me.Curve, me.Curve.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *LightAttenuation) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Constant, nil, me.Constant.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.Linear, nil, me.Linear.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.Quadratic, nil, me.Quadratic.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	return
}

func (me *GeometryMesh) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *AnimationClipDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Animations {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.Formulas {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *AnimationClipDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *KxMotionAxis) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *GeometryBrepNurbs) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxSceneDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.ForceFields {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.Models {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if me.TC.Gravity != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.Gravity, nil, me.TC.Gravity.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if me.TC.TimeStep != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.TimeStep, nil, me.TC.TimeStep.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxSceneDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *GeometryBrepNurbsSurface) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxProfileGlsl) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.CodesIncludes {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.Techniques {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *GeometryDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Brep != nil {
		if val = me.Brep.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Mesh != nil {
		if val = me.Mesh.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Spline != nil {
		if val = me.Spline.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *GeometryDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *FormulaDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FormulaDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *NodeDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.Insts.Camera {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.Insts.Controller {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.Insts.Geometry {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.Insts.Light {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, subItem := range me.Nodes {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *NodeDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *KxFrameOrigin) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *ControllerInst) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.BindMaterial != nil {
		if val = me.BindMaterial.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *GeometryBrepSurfaceCurves) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.All {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxKinematicsSystem) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Models {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.TC.AxisInfos {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if val = me.TC.Frame.Origin.resolveSidPath(path, bag); val != nil {
		return
	}
	if val = me.TC.Frame.Tip.resolveSidPath(path, bag); val != nil {
		return
	}
	if me.TC.Frame.Tcp != nil {
		if val = me.TC.Frame.Tcp.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.TC.Frame.Object != nil {
		if val = me.TC.Frame.Object.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxModelDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.RigidBodies {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.RigidConstraints {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.Insts {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxModelDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *FxColorOrTexture) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Color != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Color, nil, me.Color.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxRigidConstraintAttachment) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *GeometryBrepSurfaces) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.All {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxArticulatedSystemInst) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *ControllerSkin) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxFrame) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxRigidBodyDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.Dynamic, nil, me.TC.Dynamic.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	if me.TC.Mass != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.Mass, nil, me.TC.Mass.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.TC.MassFrame {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if me.TC.Inertia != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.Inertia, nil, me.TC.Inertia.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if val = me.TC.Material.resolveSidPath(path, bag); val != nil {
		return
	}
	for _, subItem := range me.TC.Shapes {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxRigidBodyDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *CameraDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if val = me.Optics.resolveSidPath(path, bag); val != nil {
		return
	}
	return
}

func (me *CameraDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *Source) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.TC.Accessor != nil {
		if val = me.TC.Accessor.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *Source) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *KxFrameObject) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxPassEvaluation) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Color.Target != nil {
		if val = me.Color.Target.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Depth.Target != nil {
		if val = me.Depth.Target.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Stencil.Target != nil {
		if val = me.Stencil.Target.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxModelInst) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.ForceFields {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.RigidBodies {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.RigidConstraints {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxTechniqueGlsl) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Passes {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxTechniqueGlsl) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *KxSceneDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Models {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.ArticulatedSystems {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxSceneDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *AnimationDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	for _, subItem := range me.AnimationDefs {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *AnimationDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *FxProfileCommon) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Technique, &me.Technique, me.Technique.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	return
}

func (me *KxModelDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.TC.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.TC.Links {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, subItem := range me.TC.Formulas {
		if val = subItem.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxModelDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *FxTechniqueCommon) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Emission != nil {
		if val = me.Emission.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Reflective != nil {
		if val = me.Reflective.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Reflectivity != nil {
		if val = me.Reflectivity.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Transparent != nil {
		if val = me.Transparent.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Transparency != nil {
		if val = me.Transparency.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.IndexOfRefraction != nil {
		if val = me.IndexOfRefraction.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Ambient != nil {
		if val = me.Ambient.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Diffuse != nil {
		if val = me.Diffuse.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Specular != nil {
		if val = me.Specular.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	if me.Shininess != nil {
		if val = me.Shininess.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxTechniqueCommon) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LightPoint) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if val = me.Attenuation.resolveSidPath(path, bag); val != nil {
		return
	}
	return
}

func (me *KxJointDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.All {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxJointDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *FxPass) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Evaluate != nil {
		if val = me.Evaluate.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxFrameTcp) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *Formula) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.Inst != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Inst, nil, me.Inst.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxMotionSystem) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.ArticulatedSystem != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.ArticulatedSystem, me.ArticulatedSystem, me.ArticulatedSystem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	for _, sidItem := range me.TC.AxisInfos {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if me.TC.EffectorInfo != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.EffectorInfo, me.TC.EffectorInfo, me.TC.EffectorInfo.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxJoint) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Axis, nil, me.Axis.Sid
	if val = sidResolveCore(path, bag); val != nil {
		return
	}
	if me.Limits != nil {
		if val = me.Limits.resolveSidPath(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *CameraPerspective) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	if me.FovX != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.FovX, nil, me.FovX.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	if me.FovY != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.FovY, nil, me.FovY.Sid
		if val = sidResolveCore(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxForceFieldDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	return
}

func (me *PxForceFieldDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *FxImageDef) resolveSidPath(path []string, bag *refSidBag) (val interface{}) {
	return
}

func (me *FxImageDef) resolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *PxRigidConstraintLimit) accessField(fn string) interface{} {
	switch fn {
	case "Min":
		return &me.Min.Vec3
	case "Max":
		return &me.Max.Vec3
	}
	return nil
}

func (me *NodeInst) accessField(fn string) interface{} {
	switch fn {
	case "Proxy":
		return &me.Proxy
	}
	return nil
}

func (me *PxModelInst) accessField(fn string) interface{} {
	switch fn {
	case "Parent":
		return &me.Parent
	}
	return nil
}

func (me *VisualSceneRendering) accessField(fn string) interface{} {
	switch fn {
	case "CameraNode":
		return &me.CameraNode
	}
	return nil
}

func (me *FxParamDef) accessField(fn string) interface{} {
	switch fn {
	case "Modifier":
		return &me.Modifier
	case "Semantic":
		return &me.Semantic
	}
	return nil
}

func (me *Param) accessField(fn string) interface{} {
	switch fn {
	case "Semantic":
		return &me.Semantic
	case "Type":
		return &me.Type
	}
	return nil
}

func (me *PxShape) accessField(fn string) interface{} {
	switch fn {
	case "Hollow":
		return &me.Hollow.B
	}
	return nil
}

func (me *FxProfileGlslCodeInclude) accessField(fn string) interface{} {
	switch fn {
	case "IsInclude":
		return &me.IsInclude
	}
	return nil
}

func (me *VisualSceneEvaluation) accessField(fn string) interface{} {
	switch fn {
	case "Disabled":
		return &me.Disabled
	}
	return nil
}

func (me *PxRigidBodyInst) accessField(fn string) interface{} {
	switch fn {
	case "TargetNode":
		return &me.TargetNode
	}
	return nil
}

func (me *PxRigidConstraintSpring) accessField(fn string) interface{} {
	switch fn {
	case "Stiffness":
		return &me.Stiffness.F
	case "Damping":
		return &me.Damping.F
	case "TargetValue":
		return &me.TargetValue.F
	}
	return nil
}

func (me *ParamDef) accessField(fn string) interface{} {
	switch fn {
	case "Value":
		return &me.Value
	}
	return nil
}

func (me *LightAttenuation) accessField(fn string) interface{} {
	switch fn {
	case "Constant":
		return &me.Constant.F
	case "Linear":
		return &me.Linear.F
	case "Quadratic":
		return &me.Quadratic.F
	}
	return nil
}

func (me *NodeDef) accessField(fn string) interface{} {
	switch fn {
	case "IsSkinJoint":
		return &me.IsSkinJoint
	}
	return nil
}

func (me *PxRigidBodyCommon) accessField(fn string) interface{} {
	switch fn {
	case "Dynamic":
		return &me.Dynamic.B
	}
	return nil
}

func (me *FxMaterialInst) accessField(fn string) interface{} {
	switch fn {
	case "Symbol":
		return &me.Symbol
	}
	return nil
}
