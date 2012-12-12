package assets

import (
	unum "github.com/metaleap/go-util/num"
)

type GeometryBrep struct {
	HasExtras
	HasSources
	Vertices      GeometryVertices
	Curves        *GeometryBrepCurves
	SurfaceCurves *GeometryBrepSurfaceCurves
	Surfaces      *GeometryBrepSurfaces
	Edges         *GeometryBrepEdges
	Wires         *GeometryBrepWires
	Faces         *GeometryBrepFaces
	Pcurves       *GeometryBrepPcurves
	Shells        *GeometryBrepShells
	Solids        *GeometryBrepSolids
}

func NewGeometryBrep() (me *GeometryBrep) {
	me = &GeometryBrep{}
	me.Sources = Sources{}
	return
}

type GeometryBrepBox struct {
	HasExtras
	HalfExtents Float3
}

type GeometryBrepCapsule struct {
	HasExtras
	Height float64
	Radius Float3
}

type GeometryBrepCircle struct {
	HasExtras
	Radius float64
}

type GeometryBrepCone struct {
	HasExtras
	Angle  float64
	Radius float64
}

type GeometryBrepCurve struct {
	HasSid
	HasName
	Orientations []*GeometryBrepOrientation
	Origin       *unum.Vec3
	Element      struct {
		Line      *GeometryBrepLine
		Circle    *GeometryBrepCircle
		Ellipse   *GeometryBrepEllipse
		Parabola  *GeometryBrepParabola
		Hyperbola *GeometryBrepHyperbola
		Nurbs     *GeometryBrepNurbs
	}
}

type GeometryBrepCurves struct {
	HasExtras
	Curves []*GeometryBrepCurve
}

type GeometryBrepCylinder struct {
	HasExtras
	Radius Float2
}

type GeometryBrepEdges struct {
	HasID
	HasName
	HasExtras
	IndexedInputs
}

type GeometryBrepEllipse struct {
	HasExtras
	Radius Float2
}

type GeometryBrepFaces struct {
	HasID
	HasName
	HasExtras
	IndexedInputsV
}

type GeometryBrepHyperbola struct {
	HasExtras
	Radius Float2
}

type GeometryBrepLine struct {
	HasExtras
	Origin    unum.Vec3
	Direction unum.Vec3
}

type GeometryBrepNurbs struct {
	HasExtras
	HasSources
	Degree          uint64
	Closed          bool
	ControlVertices GeometryControlVertices
}

func NewGeometryBrepNurbs() (me *GeometryBrepNurbs) {
	me = &GeometryBrepNurbs{}
	me.Sources = Sources{}
	return
}

type GeometryBrepNurbsSurface struct {
	HasExtras
	HasSources
	U, V struct {
		Degree uint64
		Closed bool
	}
	ControlVertices GeometryControlVertices
}

func NewGeometryBrepNurbsSurface() (me *GeometryBrepNurbsSurface) {
	me = &GeometryBrepNurbsSurface{}
	me.Sources = Sources{}
	return
}

type GeometryBrepOrientation struct {
	Axis  unum.Vec3
	Angle float64
}

type GeometryBrepParabola struct {
	HasExtras
	FocalLength float64
}

type GeometryBrepPcurves struct {
	HasID
	HasName
	HasExtras
	IndexedInputsV
}

type GeometryBrepPlane struct {
	HasExtras
	Equation Float4
}

type GeometryBrepShells struct {
	HasID
	HasName
	HasExtras
	IndexedInputsV
}

type GeometryBrepSolids struct {
	HasID
	HasName
	HasExtras
	IndexedInputsV
}

type GeometryBrepSphere struct {
	HasExtras
	Radius float64
}

type GeometryBrepSurface struct {
	HasSid
	HasName
	Orientations []*GeometryBrepOrientation
	Origin       *unum.Vec3
	Element      struct {
		Cone         *GeometryBrepCone
		Plane        *GeometryBrepPlane
		Cylinder     *GeometryBrepCylinder
		NurbsSurface *GeometryBrepNurbsSurface
		Sphere       *GeometryBrepSphere
		Torus        *GeometryBrepTorus
		SweptSurface *GeometryBrepSweptSurface
	}
}

type GeometryBrepSurfaceCurves struct {
	HasExtras
	Curves []*GeometryBrepCurve
}

type GeometryBrepSurfaces struct {
	HasExtras
	Surfaces []*GeometryBrepSurface
}

type GeometryBrepSweptSurface struct {
	HasExtras
	Curve              *GeometryBrepCurve
	ExtrusionDirection *unum.Vec3
	Revolution         struct {
		Origin *unum.Vec3
		Axis   *unum.Vec3
	}
}

type GeometryBrepTorus struct {
	HasExtras
	Radius Float2
}

type GeometryBrepWires struct {
	HasID
	HasName
	HasExtras
	IndexedInputsV
}
