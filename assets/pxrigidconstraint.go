package assets

//	Defines an attachment frame (or attachment frame of reference),
//	to a rigid body or a node, within a rigid constraint.
type PxRigidConstraintAttachment struct {
	//	Extras
	HasExtras
	//	Refers to a RigidBodyDef or NodeDef.
	RigidBody RefSid
	//	Zero or more translation and/or rotation transformations:
	//	The position of a TRANSFORM_TYPE_TRANSLATE Transform indicates
	//	the attachment point on the corresponding RigidBodyDef.
	//	The orientation of a TRANSFORM_TYPE_ROTATE Transform indicates
	//	the alignment of the joint frame for that RigidBodyDef.
	Transforms []*Transform
}

//	Degrees of freedom and ranges.
type PxRigidConstraintLimit struct {
	//	Lower bounds for this limit across all 3 axes.
	Min ScopedVec3
	//	Upper bounds for this limit across all 3 axes.
	Max ScopedVec3
}

//	Spring is based based on either distance (Linear) or angle (Angular), or both.
type PxRigidConstraintSpring struct {
	//	Also called spring coefficient.
	//	Has units of force/distance (for Linear) or force/angle (for Angular).
	Stiffness ScopedFloat
	//	How this spring is damped.
	Damping ScopedFloat
	//	Target value for this spring.
	TargetValue ScopedFloat
}

//	Constructor
func NewPxRigidConstraintSpring() (me *PxRigidConstraintSpring) {
	me = &PxRigidConstraintSpring{}
	me.Stiffness.F = 1
	return
}

//	Constrains rigid bodies to each other or to the world.
type PxRigidConstraintDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	Sid
	HasSid
	//	Techniques
	HasTechniques
	//	Defines the attachment frame of reference (to a rigid_body or a node) within this rigid constraint.
	RefAttachment PxRigidConstraintAttachment
	//	Defines an attachment frame (to a rigid body or a node) within this rigid constraint.
	Attachment PxRigidConstraintAttachment
	//	Common-technique profile
	TC struct {
		//	If false, this rigid constraint doesn't exert any force or influence on the rigid bodies.
		Enabled ScopedBool
		//	If true, the attached rigid bodies may interpenetrate.
		Interpenetrate ScopedBool
		//	Degrees of freedom and ranges.
		Limits struct {
			//	Describes the angular limits along each rotation axis in degrees.
			Angular *PxRigidConstraintLimit
			//	Describes linear (translational) limits along each axis.
			Linear *PxRigidConstraintLimit
		}
		//	Spring is based based on either distance (Linear) or angle (Angular), or both.
		Spring struct {
			//	Angle-based spring.
			Angular *PxRigidConstraintSpring
			//	Distance-based spring.
			Linear *PxRigidConstraintSpring
		}
	}
}

//	Initialization
func (me *PxRigidConstraintDef) Init() {
	me.TC.Enabled.B = true
}

//	Instantiates a rigid constraint resource.
type PxRigidConstraintInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
}

//	Initialization
func (me *PxRigidConstraintInst) Init() {
}
