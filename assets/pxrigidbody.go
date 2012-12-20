package assets

import (
	unum "github.com/metaleap/go-util/num"
)

//	Declares a cylinder primitive that is centered around its local origin
//	and aligned along its local y axis.
type PxCylinder struct {
	//	Radii, Extras
	GeometryBrepCylinder
	//	Contains a floating-point value that represents the length of the cylinder along the y axis.
	Height float64
}

//	Describes the physical surface properties for a rigid body or one of its shapes.
//	Either Def or Inst (or none), but not both, should be specified.
type PxMaterial struct {
	//	An inline physics material definition.
	Def *PxMaterialDef
	//	Instantiation of a previously defined physics material.
	Inst *PxMaterialInst
}

//	A component part of a rigid body's collection of bounding shapes for collision detection.
type PxShape struct {
	//	Extras
	HasExtras
	//	If true, the mass is distributed along the surface of this shape.
	Hollow SidBool
	//	The mass of this shape.
	Mass *SidFloat
	//	The density of this shape.
	Density *SidFloat
	//	Describes the physical surface properties for this shape.
	Material PxMaterial
	//	Geometry of the shape. At least and at most one of its fields should ever be set.
	Geometry struct {
		//	A flat plane.
		Plane *GeometryBrepPlane
		//	Axis-aligned box.
		Box *GeometryBrepBox
		//	A perfectly round sphere.
		Sphere *GeometryBrepSphere
		//	A cylinder primitive that is centered on its local origin and aligned along its y axis.
		Cylinder *PxCylinder
		//	A capsule primitive that is centered on the local origin and aligned along the y axis.
		Capsule *GeometryBrepCapsule
		//	Refers to a previously defined mesh or spline geometric primitive.
		Inst *GeometryInst
	}
	//	Zero or more TRANSFORM_TYPE_ROTATE and/or TRANSFORM_TYPE_TRANSLATE transformations for the shape.
	Transforms []*Transform
}

//	Common-technique profile for rigid body definitions and instances.
type PxRigidBodyCommon struct {
	//	Specifies whether this rigid body is movable. Defaults to true.
	Dynamic SidBool
	//	If set, specifies the total mass of this rigid body.
	Mass *SidFloat
	//	Zero or more TRANSFORM_TYPE_TRANSLATE and/or TRANSFORM_TYPE_ROTATE transformations defining the
	//	center and orientation of mass of the rigid-body relative to the local origin of the "root" shape.
	//	This makes the off-diagonal elements of the inertia tensor (products of inertia) all 0
	//	and allows us to just store the diagonal elements (moments of inertia).
	MassFrame []*Transform
	//	The diagonal elements of the inertia tensor (moments of inertia),
	//	represented in the local frame of the center of mass.
	Inertia *SidFloat3
	//	Describes the physical surface properties for this rigid body.
	Material PxMaterial
	//	Zero or more shapes for collision detection.
	Shapes []*PxShape
}

type PxRigidBodyDef struct {
	//	Id, Name, Asset, Extras
	BaseDef
	//	Sid
	HasSid
	//	Techniques
	HasTechniques
	//	Common-technique profile
	TC struct {
		//	Dynamic, Mass, MassFrame, Inertia, Material, Shapes
		PxRigidBodyCommon
	}
}

//	Initialization
func (me *PxRigidBodyDef) Init() {
	me.TC.Dynamic.B = true
}

//	Instantiates a rigid-body resource.
type PxRigidBodyInst struct {
	//	Sid, Name, Extras, DefRef
	BaseInst
	//	A pointer to the resource definition referenced by this instance.
	//	Is nil by default and meant to be set ONLY by the EnsureDef() method (which uses BaseInst.DefRef to find it).
	Def *PxRigidBodyDef
	//	Techniques
	HasTechniques
	//	Refers to the NodeDef influenced by this rigid body instance.
	TargetNode RefId
	//	Common-technique profile
	TC struct {
		//	Dynamic, Mass, MassFrame, Inertia, Material, Shapes
		PxRigidBodyCommon
		//	Initial spin or angular velocity, also known as the axis of rotation,
		//	with a magnitude equal to the rate of rotation in radians per second.
		AngularVelocity unum.Vec3
		//	Initial linear velocity.
		LinearVelocity unum.Vec3
	}
}

//	Initialization
func (me *PxRigidBodyInst) Init() {
	me.TC.PxRigidBodyCommon.Dynamic.B = true
}
