package assets

func (me *GeometryBrepSweptSurface) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Curve != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Curve, nil, me.Curve.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *GeometryBrep) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Curves != nil {
		if val = me.Curves.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.SurfaceCurves != nil {
		if val = me.SurfaceCurves.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Surfaces != nil {
		if val = me.Surfaces.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxSceneInst) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *ControllerInst) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.BindMaterial != nil {
		if val = me.BindMaterial.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *LightSpot) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if val = me.Attenuation.sidResolve(path, bag); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.Falloff.Angle, nil, me.Falloff.Angle.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.Falloff.Exponent, nil, me.Falloff.Exponent.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	return
}

func (me *FxColorOrTexture) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Color != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Color, nil, me.Color.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *VisualSceneDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Nodes {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.Evaluations {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *VisualSceneDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibVisualSceneDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *PxMaterialDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.DynamicFriction, nil, me.TC.DynamicFriction.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.Restitution, nil, me.TC.Restitution.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.StaticFriction, nil, me.TC.StaticFriction.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	return
}

func (me *PxMaterialDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibPxMaterialDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *ControllerSkin) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *GeometryBrepCurve) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Element.Nurbs != nil {
		if val = me.Element.Nurbs.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *LightAttenuation) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Constant, nil, me.Constant.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.Linear, nil, me.Linear.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.Quadratic, nil, me.Quadratic.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	return
}

func (me *FxMaterialDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Effect, nil, me.Effect.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	return
}

func (me *FxMaterialDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibFxMaterialDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *PxModelDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.RigidBodies {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.RigidConstraints {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.Insts {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *PxModelDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibPxModelDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *AnimationDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	for _, subItem := range me.AnimationDefs {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *AnimationDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibAnimationDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *PxRigidConstraintLimit) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Min, nil, me.Min.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.Max, nil, me.Max.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	return
}

func (me *KxEffector) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *Source) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.TC.Accessor != nil {
		if val = me.TC.Accessor.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *Source) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *GeometryBrepCurves) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.All {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *FxProfileCommon) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Technique, &me.Technique, me.Technique.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	return
}

func (me *PxSceneDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.ForceFields {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.Models {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if me.TC.Gravity != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.Gravity, nil, me.TC.Gravity.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if me.TC.TimeStep != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.TimeStep, nil, me.TC.TimeStep.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *PxSceneDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibPxSceneDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *CameraDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if val = me.Optics.sidResolve(path, bag); val != nil {
		return
	}
	return
}

func (me *CameraDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibCameraDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *KxModelInst) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *KxFrameObject) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *ParamOrSidFloat) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.F, nil, me.F.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	return
}

func (me *GeometryDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Brep != nil {
		if val = me.Brep.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Mesh != nil {
		if val = me.Mesh.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Spline != nil {
		if val = me.Spline.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *GeometryDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibGeometryDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *CameraOptics) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.TC.AspectRatio != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.AspectRatio, nil, me.TC.AspectRatio.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.Zfar, nil, me.TC.Zfar.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.Znear, nil, me.TC.Znear.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	if me.TC.Orthographic != nil {
		if val = me.TC.Orthographic.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.TC.Perspective != nil {
		if val = me.TC.Perspective.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxRigidConstraintDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if val = me.RefAttachment.sidResolve(path, bag); val != nil {
		return
	}
	if val = me.Attachment.sidResolve(path, bag); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.Enabled, nil, me.TC.Enabled.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.Interpenetrate, nil, me.TC.Interpenetrate.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	if me.TC.Limits.Angular != nil {
		if val = me.TC.Limits.Angular.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.TC.Limits.Linear != nil {
		if val = me.TC.Limits.Linear.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.TC.Spring.Angular != nil {
		if val = me.TC.Spring.Angular.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.TC.Spring.Linear != nil {
		if val = me.TC.Spring.Linear.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxRigidConstraintDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *GeometryBrepNurbsSurface) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxFrameTip) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *FxSampler) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Image != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Image, nil, me.Image.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *KxJoint) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Axis, nil, me.Axis.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	if me.Limits != nil {
		if val = me.Limits.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *LightDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.TC.Point != nil {
		if val = me.TC.Point.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.TC.Spot != nil {
		if val = me.TC.Spot.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *LightDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibLightDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *GeometryInst) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.MaterialBinding != nil {
		if val = me.MaterialBinding.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxEffectDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, subItem := range me.Profiles {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxEffectDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibFxEffectDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *PxRigidBodyCommon) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Dynamic, nil, me.Dynamic.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	if me.Mass != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Mass, nil, me.Mass.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.MassFrame {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if me.Inertia != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Inertia, nil, me.Inertia.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if val = me.Material.sidResolve(path, bag); val != nil {
		return
	}
	for _, subItem := range me.Shapes {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxSceneDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Models {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.ArticulatedSystems {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *KxSceneDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibKxSceneDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *PxMaterial) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Inst != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Inst, nil, me.Inst.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *GeometryBrepNurbs) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxTechniqueGlsl) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Passes {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *FxTechniqueGlsl) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *KxFrameOrigin) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *FxProfileGlsl) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.CodesIncludes {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.Techniques {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *KxModelDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.TC.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.TC.Links {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, subItem := range me.TC.Formulas {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxModelDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibKxModelDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *NodeDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.Insts.Camera {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.Insts.Controller {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.Insts.Geometry {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.Insts.Light {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, subItem := range me.Nodes {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *NodeDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibNodeDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *Scene) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Visual != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Visual, nil, me.Visual.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if me.Kinematics != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Kinematics, me.Kinematics, me.Kinematics.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.Physics {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *LightPoint) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if val = me.Attenuation.sidResolve(path, bag); val != nil {
		return
	}
	return
}

func (me *FxPass) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Evaluate != nil {
		if val = me.Evaluate.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxKinematicsAxis) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, subItem := range me.Formulas {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *ChildNode) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Inst != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Inst, nil, me.Inst.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *KxArticulatedSystemInst) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *KxMotionAxis) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *VisualSceneEvaluation) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.RenderPasses {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *VisualSceneEvaluation) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *KxMotionSystem) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.ArticulatedSystem != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.ArticulatedSystem, me.ArticulatedSystem, me.ArticulatedSystem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.TC.AxisInfos {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if me.TC.EffectorInfo != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.EffectorInfo, me.TC.EffectorInfo, me.TC.EffectorInfo.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *Document) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Scene != nil {
		if val = me.Scene.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxFrame) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *FxPassEvaluationTarget) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Image != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Image, nil, me.Image.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *FormulaDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *FormulaDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibFormulaDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *KxJointLimits) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Min != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Min, nil, me.Min.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if me.Max != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Max, nil, me.Max.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *ControllerMorph) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *Formula) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Inst != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Inst, nil, me.Inst.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *ControllerDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Morph != nil {
		if val = me.Morph.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Skin != nil {
		if val = me.Skin.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *ControllerDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibControllerDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *GeometryBrepSurfaces) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.All {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *AnimationClipDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Animations {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.Formulas {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *AnimationClipDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibAnimationClipDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *GeometryBrepSurface) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Element.NurbsSurface != nil {
		if val = me.Element.NurbsSurface.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Element.SweptSurface != nil {
		if val = me.Element.SweptSurface.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxRigidBodyInst) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.Dynamic, nil, me.TC.Dynamic.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	if me.TC.Mass != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.Mass, nil, me.TC.Mass.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.TC.MassFrame {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if me.TC.Inertia != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.Inertia, nil, me.TC.Inertia.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if val = me.TC.Material.sidResolve(path, bag); val != nil {
		return
	}
	for _, subItem := range me.TC.Shapes {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxRigidConstraintAttachment) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *PxShape) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Hollow, nil, me.Hollow.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	if me.Mass != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Mass, nil, me.Mass.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if me.Density != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Density, nil, me.Density.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if val = me.Material.sidResolve(path, bag); val != nil {
		return
	}
	if me.Geometry.Inst != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Geometry.Inst, me.Geometry.Inst, me.Geometry.Inst.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *KxArticulatedSystemDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Kinematics != nil {
		if val = me.Kinematics.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Motion != nil {
		if val = me.Motion.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxArticulatedSystemDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibKxArticulatedSystemDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *CameraPerspective) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.FovX != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.FovX, nil, me.FovX.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if me.FovY != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.FovY, nil, me.FovY.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *SourceAccessor) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Params {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *GeometrySpline) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxTechniqueCommon) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Emission != nil {
		if val = me.Emission.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Reflective != nil {
		if val = me.Reflective.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Reflectivity != nil {
		if val = me.Reflectivity.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Transparent != nil {
		if val = me.Transparent.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Transparency != nil {
		if val = me.Transparency.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.IndexOfRefraction != nil {
		if val = me.IndexOfRefraction.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Ambient != nil {
		if val = me.Ambient.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Diffuse != nil {
		if val = me.Diffuse.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Specular != nil {
		if val = me.Specular.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Shininess != nil {
		if val = me.Shininess.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxTechniqueCommon) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *GeometryMesh) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, subItem := range me.Sources {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *KxFrameTcp) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *KxLink) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, subItem := range me.Attachments {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxModelInst) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.ForceFields {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.RigidBodies {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.RigidConstraints {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *KxKinematicsSystem) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Models {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.TC.AxisInfos {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if val = me.TC.Frame.Origin.sidResolve(path, bag); val != nil {
		return
	}
	if val = me.TC.Frame.Tip.sidResolve(path, bag); val != nil {
		return
	}
	if me.TC.Frame.Tcp != nil {
		if val = me.TC.Frame.Tcp.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.TC.Frame.Object != nil {
		if val = me.TC.Frame.Object.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *CameraOrthographic) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.MagX != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.MagX, nil, me.MagX.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if me.MagY != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.MagY, nil, me.MagY.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *GeometryBrepSurfaceCurves) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.All {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *PxRigidBodyDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.TC.Dynamic, nil, me.TC.Dynamic.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	if me.TC.Mass != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.Mass, nil, me.TC.Mass.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.TC.MassFrame {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if me.TC.Inertia != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.TC.Inertia, nil, me.TC.Inertia.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if val = me.TC.Material.sidResolve(path, bag); val != nil {
		return
	}
	for _, subItem := range me.TC.Shapes {
		if val = subItem.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *PxRigidBodyDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *KxJointDef) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.All {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, sidItem, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *KxJointDef) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

func (me *LibKxJointDefs) sidResolver(id string) refSidResolver {
	return me.M[id]
}

func (me *PxRigidConstraintSpring) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	bag.valRaw, bag.valAsRes, bag.sid = &me.Stiffness, nil, me.Stiffness.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.Damping, nil, me.Damping.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	bag.valRaw, bag.valAsRes, bag.sid = &me.TargetValue, nil, me.TargetValue.Sid
	if val = bag.sidResolve(path); val != nil {
		return
	}
	return
}

func (me *MaterialBinding) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Params {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	for _, sidItem := range me.TC.Materials {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *KxAttachment) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.Transforms {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if me.Link != nil {
		bag.valRaw, bag.valAsRes, bag.sid = me.Link, me.Link, me.Link.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	return
}

func (me *FxPassEvaluation) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	if me.Color.Target != nil {
		if val = me.Color.Target.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Depth.Target != nil {
		if val = me.Depth.Target.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Stencil.Target != nil {
		if val = me.Stencil.Target.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxProfile) sidResolve(path []string, bag *refSidBag) (val interface{}) {
	for _, sidItem := range me.NewParams {
		bag.valRaw, bag.valAsRes, bag.sid = sidItem, nil, sidItem.Sid
		if val = bag.sidResolve(path); val != nil {
			return
		}
	}
	if me.Common != nil {
		if val = me.Common.sidResolve(path, bag); val != nil {
			return
		}
	}
	if me.Glsl != nil {
		if val = me.Glsl.sidResolve(path, bag); val != nil {
			return
		}
	}
	return
}

func (me *FxProfile) sidResolver(id string) (rsr refSidResolver) {
	if (id == me.Id) || (id == ".") {
		rsr = me
	}
	return
}

//	RefSidFielder implementation.
//	Supported field names: "Semantic", "Type".
func (me *Param) AccessField(fn string) interface{} {
	switch fn {
	case "Semantic":
		return &me.Semantic
	case "Type":
		return &me.Type
	}
	return nil
}

//	RefSidFielder implementation.
//	Supported field names: "IsInclude".
func (me *FxProfileGlslCodeInclude) AccessField(fn string) interface{} {
	switch fn {
	case "IsInclude":
		return &me.IsInclude
	}
	return nil
}

//	RefSidFielder implementation.
//	Supported field names: "Dynamic".
func (me *PxRigidBodyCommon) AccessField(fn string) interface{} {
	switch fn {
	case "Dynamic":
		return &me.Dynamic.B
	}
	return nil
}

//	RefSidFielder implementation.
//	Supported field names: "Modifier", "Semantic".
func (me *FxParamDef) AccessField(fn string) interface{} {
	switch fn {
	case "Modifier":
		return &me.Modifier
	case "Semantic":
		return &me.Semantic
	}
	return nil
}

//	RefSidFielder implementation.
//	Supported field names: "CameraNode".
func (me *VisualSceneRendering) AccessField(fn string) interface{} {
	switch fn {
	case "CameraNode":
		return &me.CameraNode
	}
	return nil
}

//	RefSidFielder implementation.
//	Supported field names: "Proxy".
func (me *NodeInst) AccessField(fn string) interface{} {
	switch fn {
	case "Proxy":
		return &me.Proxy
	}
	return nil
}

//	RefSidFielder implementation.
//	Supported field names: "Hollow".
func (me *PxShape) AccessField(fn string) interface{} {
	switch fn {
	case "Hollow":
		return &me.Hollow.B
	}
	return nil
}

//	RefSidFielder implementation.
//	Supported field names: "Constant", "Linear", "Quadratic".
func (me *LightAttenuation) AccessField(fn string) interface{} {
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

//	RefSidFielder implementation.
//	Supported field names: "IsSkinJoint".
func (me *NodeDef) AccessField(fn string) interface{} {
	switch fn {
	case "IsSkinJoint":
		return &me.IsSkinJoint
	}
	return nil
}

//	RefSidFielder implementation.
//	Supported field names: "Value".
func (me *ParamDef) AccessField(fn string) interface{} {
	switch fn {
	case "Value":
		return &me.Value
	}
	return nil
}

//	RefSidFielder implementation.
//	Supported field names: "Stiffness", "Damping", "TargetValue".
func (me *PxRigidConstraintSpring) AccessField(fn string) interface{} {
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

//	RefSidFielder implementation.
//	Supported field names: "TargetNode".
func (me *PxRigidBodyInst) AccessField(fn string) interface{} {
	switch fn {
	case "TargetNode":
		return &me.TargetNode
	}
	return nil
}

//	RefSidFielder implementation.
//	Supported field names: "Disabled".
func (me *VisualSceneEvaluation) AccessField(fn string) interface{} {
	switch fn {
	case "Disabled":
		return &me.Disabled
	}
	return nil
}

//	RefSidFielder implementation.
//	Supported field names: "Symbol".
func (me *FxMaterialInst) AccessField(fn string) interface{} {
	switch fn {
	case "Symbol":
		return &me.Symbol
	}
	return nil
}

//	RefSidFielder implementation.
//	Supported field names: "Parent".
func (me *PxModelInst) AccessField(fn string) interface{} {
	switch fn {
	case "Parent":
		return &me.Parent
	}
	return nil
}

//	RefSidFielder implementation.
//	Supported field names: "Min", "Max".
func (me *PxRigidConstraintLimit) AccessField(fn string) interface{} {
	switch fn {
	case "Min":
		return &me.Min.Vec3
	case "Max":
		return &me.Max.Vec3
	}
	return nil
}
