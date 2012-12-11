package assets

type PxRigidConstraintAttachment struct {
	HasExtras
	RigidBody  string
	Transforms []*Transform
}

type PxRigidConstraintLimit struct {
	Min ScopedFloat3
	Max ScopedFloat3
}

type PxRigidConstraintSpring struct {
	Stiffness   ScopedFloat
	Damping     ScopedFloat
	TargetValue ScopedFloat
}

type PxRigidConstraintDef struct {
	BaseDef
	HasSid
	HasTechniques
	RefAttachment   PxRigidConstraintAttachment
	Attachment      PxRigidConstraintAttachment
	TechniqueCommon struct {
		Disabled       ScopedBool
		Interpenetrate ScopedBool
		Limits         struct {
			SwingConeAndTwist *PxRigidConstraintLimit
			Linear            *PxRigidConstraintLimit
		}
		Spring struct {
			Angular *PxRigidConstraintSpring
			Linear  *PxRigidConstraintSpring
		}
	}
}

func (me *PxRigidConstraintDef) Init() {
}

type PxRigidConstraintInst struct {
	BaseInst
}

func (me *PxRigidConstraintInst) Init() {
}
