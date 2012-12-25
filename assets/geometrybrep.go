package assets

import (
	unum "github.com/metaleap/go-util/num"
)

//	Describes a boundary representation (B-rep) structure.
type GeometryBrep struct {
	//	Extras
	HasExtras

	//	Sources
	HasSources

	//	Describes all vertices of the B-rep.
	//	Vertices are the base topological entity for all B-rep structures.
	Vertices GeometryVertices

	//	Contains all curves used in this B-rep. Required if Edges are present.
	Curves *GeometryBrepCurves

	//	Contains all 2D curves used in this B-rep.
	//	This includes surfaces that describe the kind of the face. Required if Faces are present.
	SurfaceCurves *GeometryBrepSurfaceCurves

	//	Contains all surfaces used in this B-rep.
	Surfaces *GeometryBrepSurfaces

	//	Describes all edges of the B-rep.
	Edges *GeometryBrepEdges

	//	Describes all wires of the B-rep.
	Wires *GeometryBrepWires

	//	Describes all faces of the B-rep.
	Faces *GeometryBrepFaces

	//	Describes all pcurves of the B-rep.
	Pcurves *GeometryBrepPcurves

	//	Describes all shells of the B-rep.
	Shells *GeometryBrepShells

	//	Describes all solids of the B-rep.
	Solids *GeometryBrepSolids
}

//	Constructor
func NewGeometryBrep() (me *GeometryBrep) {
	me = &GeometryBrep{}
	me.Sources = Sources{}
	return
}

//	Declares an axis-aligned box that is centered around its local origin.
type GeometryBrepBox struct {
	//	Extras
	HasExtras

	//	Represents the extents of the box. The dimensions of the box are double the half-extents.
	HalfExtents unum.Vec3
}

//	A capsule is a cylinder with rounded caps.
type GeometryBrepCapsule struct {
	//	Extras
	HasExtras

	//	The length of the line segment connecting the centers of the capping hemispheres (ellipsoids).
	Height float64

	//	The x, y, and z radii of the capsule (it may be elliptical).
	Radii unum.Vec3
}

//	Describes a circle in 3D space.
type GeometryBrepCircle struct {
	//	Extras
	HasExtras

	//	The radius of the circle.
	Radius float64
}

//	Describes a conical surface. A cone is defined by the half-angle at its apex.
type GeometryBrepCone struct {
	//	Extras
	HasExtras

	//	The conical surface semi-angle.
	Angle float64

	//	Radius of the cone.
	Radius float64
}

//	Describes a specific curve.
type GeometryBrepCurve struct {
	//	Sid
	HasSid

	//	Name
	HasName

	//	Optional positioning of this surface to its correct location.
	Location GeometryPositioning

	//	The curve element. At least and at most one of these fields must be set (non-nil).
	Element struct {
		//	If set, curve element is a line.
		Line *GeometryBrepLine

		//	If set, curve element is a circle.
		Circle *GeometryBrepCircle

		//	If set, curve element is an ellipse.
		Ellipse *GeometryBrepEllipse

		//	If set, curve element is a parabola.
		Parabola *GeometryBrepParabola

		//	If set, curve element is a hyperbola.
		Hyperbola *GeometryBrepHyperbola

		//	If set, curve element is a NURBS curve.
		Nurbs *GeometryBrepNurbs
	}
}

//	Contains all curves that are used in the parent B-rep structure.
type GeometryBrepCurves struct {
	//	Extras
	HasExtras

	//	A container for all 3D curves used by the edges of the parent B-rep structure.
	All []*GeometryBrepCurve
}

//	Describes an unlimited cylindrical surface.
//	An unlimited cylinder has a radius but is assumed to extend to an infinite length.
type GeometryBrepCylinder struct {
	//	Extras
	HasExtras

	//	The first value is the major radius, the second is the minor radius (cylinder may be elliptical).
	Radii Float2
}

//	Describes the edges of a B-rep structure.
type GeometryBrepEdges struct {
	//	Id
	HasId

	//	Name
	HasName

	//	Extras
	HasExtras

	//	Four inputs are needed to define an edge:
	//	One with Semantic "CURVE" to reference the corresponding geometric element for the edge.
	//	Two with Semantic "VERTEX" to reference the two vertices that limit each edge.
	//	One with Semantic "PARAM" to set the parametric values (start and end parameters) of the curve.
	IndexedInputs
}

//	Describes an ellipse in 3D space.
type GeometryBrepEllipse struct {
	//	Extras
	HasExtras

	//	The first value is the major radius, the second is the minor radius.
	Radii Float2
}

//	Describes the faces of a B-rep structure.
type GeometryBrepFaces struct {
	//	Id
	HasId

	//	Name
	HasName

	//	Extras
	HasExtras

	//	There must be at least three inputs:
	//	One with Semantic "SURFACE" to reference the corresponding geometric element for the face.
	//	One with Semantic "WIRE" to reference the wires for each face.
	//	One with Semantic "ORIENTATION" defining the orientation of the referenced wire within the face.
	IndexedInputs
}

//	Describes a hyperbola in 3D space.
type GeometryBrepHyperbola struct {
	//	Extras
	HasExtras

	//	The first value is the major radius, the second is the minor radius.
	Radii Float2
}

//	Describes a single line in 3D space.
type GeometryBrepLine struct {
	//	Extras
	HasExtras

	//	The origin of the line.
	Origin unum.Vec3

	//	The direction of the line as a unit vector.
	Direction unum.Vec3
}

//	Describes a NURBS curve in 3D space.
type GeometryBrepNurbs struct {
	//	Extras
	HasExtras

	//	Sources
	HasSources

	//	Specifies the degree of the NURBS curve.
	Degree uint64

	//	Specifies whether this NURBS curve is closed.
	Closed bool

	//	Control vertices for curve interpolation.
	ControlVertices GeometryControlVertices
}

//	Constructor
func NewGeometryBrepNurbs() (me *GeometryBrepNurbs) {
	me = &GeometryBrepNurbs{}
	me.Sources = Sources{}
	return
}

//	Describes a NURBS surface in 3D space.
type GeometryBrepNurbsSurface struct {
	//	Extras
	HasExtras

	//	Sources
	HasSources

	//	The u and v directions for the NURBS curve.
	U, V struct {
		//	Specifies the degree of the NURBS curve for this direction.
		Degree uint64

		//	Specifies whether a NURBS curve is closed for this direction.
		Closed bool
	}
	//	Control vertices for curve interpolation.
	ControlVertices GeometryControlVertices
}

//	Constructor
func NewGeometryBrepNurbsSurface() (me *GeometryBrepNurbsSurface) {
	me = &GeometryBrepNurbsSurface{}
	me.Sources = Sources{}
	return
}

//	Describes the orientation of an object frame.
type GeometryBrepOrientation struct {
	//	The axis of rotation.
	Axis unum.Vec3

	//	The rotational angle in degrees.
	Angle float64
}

//	Describes a parabola in 3D space.
type GeometryBrepParabola struct {
	//	Extras
	HasExtras

	//	The distance between the parabola's focus and its apex.
	FocalLength float64
}

//	Specifies how an edge is represented in a face's parametric space.
type GeometryBrepPcurves struct {
	//	Id
	HasId

	//	Name
	HasName

	//	Extras
	HasExtras

	//	There must be at least three inputs:
	//	One with Semantic "CURVE2D" referencing a pcurve.
	//	One with Semantic "FACE" and one with Semantic "EDGE"
	//	to specify the connection between the edge and the face.
	IndexedInputs
}

//	Defines an infinite plane.
type GeometryBrepPlane struct {
	//	Extras
	HasExtras

	//	The four coefficients for the plane's equation: Ax + By + Cz + D = 0
	Equation Float4
}

//	Describes the shells of a B-rep structure.
//	A shell is the union of one or more faces. A closed shell can limit a solid.
type GeometryBrepShells struct {
	//	Id
	HasId

	//	Name
	HasName

	//	Extras
	HasExtras

	//	There must be at least two inputs:
	//	One with Semantic "FACE" to reference the faces for each shell.
	//	One with Semantic "ORIENTATION" defining the orientation of the referenced face within the shell.
	IndexedInputs
}

//	Describes the solids of a B-rep structure.
type GeometryBrepSolids struct {
	//	Id
	HasId

	//	Name
	HasName

	//	Extras
	HasExtras

	//	There must be at least two inputs:
	//	One with Semantic "SHELL" to reference the shells for each solid.
	//	One with Semantic "ORIENTATION" defining the orientation of the referenced shell within the solid.
	IndexedInputs
}

//	Describes a perfectly round sphere that is centered around its local origin.
type GeometryBrepSphere struct {
	//	Extras
	HasExtras

	//	The radius of this sphere.
	Radius float64
}

//	Describes a specific surface.
type GeometryBrepSurface struct {
	//	Sid
	HasSid

	//	Name
	HasName

	//	Optional positioning of this surface to its correct location.
	Location GeometryPositioning

	//	The surface element. At least and at most one of these fields must be set (non-nil).
	Element struct {
		//	Surface is described by a cone.
		Cone *GeometryBrepCone

		//	Surface is described by a plane.
		Plane *GeometryBrepPlane

		//	Surface is described by a cylinder.
		Cylinder *GeometryBrepCylinder

		//	Surface is described by a NURBS surface.
		NurbsSurface *GeometryBrepNurbsSurface

		//	Surface is described by a sphere.
		Sphere *GeometryBrepSphere

		//	Surface is described by a torus.
		Torus *GeometryBrepTorus

		//	Surface is described by an extruded or revolved curve.
		SweptSurface *GeometryBrepSweptSurface
	}
}

//	Contains all parametric curves (pcurves) that are used in a B-rep structure.
type GeometryBrepSurfaceCurves struct {
	//	Extras
	HasExtras

	//	Pcurves are curves in the parametric space of the surface on which they lie.
	All []*GeometryBrepCurve
}

//	Contains all surfaces that are used in a B-rep structure.
type GeometryBrepSurfaces struct {
	//	Extras
	HasExtras

	//	A container for all surfaces used by the faces of the parent B-rep structure.
	All []*GeometryBrepSurface
}

//	Describes a surface by extruding or revolving a curve.
type GeometryBrepSweptSurface struct {
	//	Extras
	HasExtras

	//	Describes the base curve being extruded or revolved.
	Curve *GeometryBrepCurve

	//	If Direction is set (non-nil), Revolution is ignored and this surface extrudes Curve.
	Extrusion struct {
		//	The direction of this curve extrusion.
		Direction *unum.Vec3
	}
	//	Only used if Extrusion.Direction is nil; then this surface revolves Curve.
	Revolution struct {
		//	The origin of the axis for revolution.
		Origin *unum.Vec3

		//	The axis' direction for revolution.
		Direction *unum.Vec3
	}
}

//	Returns true if this surface is described by extruding a curve.
func (me *GeometryBrepSweptSurface) IsExtrusion() bool {
	return me.Extrusion.Direction != nil
}

//	Returns true if this surface is described by revolving a curve.
func (me *GeometryBrepSweptSurface) IsRevolution() bool {
	return !me.IsExtrusion()
}

//	Describes a torus in 3D space.
type GeometryBrepTorus struct {
	//	Extras
	HasExtras

	//	The first value is the major radius, the second is the minor radius.
	Radii Float2
}

//	Describes the wires of a B-rep structure.
type GeometryBrepWires struct {
	//	Id
	HasId

	//	Name
	HasName

	//	Extras
	HasExtras

	//	There must be at least inputs:
	//	One with Semantic "EDGE" to reference the edges for each wire.
	//	One with Semantic "ORIENTATION" defining the orientation of the referenced edge within the wire.
	IndexedInputs
}

//	Used to position a surface or curve to its correct location.
type GeometryPositioning struct {
	//	If set, describes the origin of the object frame.
	Origin *unum.Vec3

	//	If set, these describe the orientation of the object frame.
	Orientations []*GeometryBrepOrientation
}
