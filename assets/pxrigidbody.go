package assets

import (
	unum "github.com/metaleap/go-util/num"
)

type PxCylinder struct {
	Height float64
	Radius Float2
}

type PxShape struct {
	HasExtras
	Hollow   *ScopedBool
	Mass     *ScopedFloat
	Density  *ScopedFloat
	Material struct {
		Def  *PxMaterialDef
		Inst *PxMaterialInst
	}
	Geometry struct {
		Plane    *GeometryBrepPlane
		Box      *GeometryBrepBox
		Sphere   *GeometryBrepSphere
		Cylinder *PxCylinder
		Capsule  *GeometryBrepCapsule
		Inst     *GeometryInst
	}
	Transforms []*Transform
}

type PxRigidBodyCommon struct {
	Dynamic   *ScopedBool
	Mass      *ScopedFloat
	MassFrame []*Transform
	Inertia   *ScopedFloat3
	Material  struct {
		Def  *PxMaterialDef
		Inst *PxMaterialInst
	}
	Shapes []*PxShape
}

type PxRigidBodyDef struct {
	BaseDef
	HasSid
	HasTechniques
	TechniqueCommon PxRigidBodyCommon
}

func (me *PxRigidBodyDef) init() {
}

type PxRigidBodyInst struct {
	BaseInst
	Def             *PxRigidBodyDef
	Name            string
	Body            string
	Target          string
	HasTechniques
	TechniqueCommon struct {
		PxRigidBodyCommon
		AngularVelocity unum.Vec3
		Velocity        unum.Vec3
	}
}

func (me *PxRigidBodyInst) init() {
}
