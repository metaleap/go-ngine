package assets

//	References a resource by its unique identifier (Id).
type RefId string

//	Searches (all LibAnimationDefs contained in AllAnimationDefLibs) for the AnimationSampler
//	whose Id is referenced by me, returning the first match found.
func (me RefId) AnimationSampler() (as *AnimationSampler) {
	var (
		def *AnimationDef
		id  = me.S()
	)
	for _, lib := range AllAnimationDefLibs {
		for _, def = range lib.M {
			for _, as = range def.Samplers {
				if as.Id == id {
					return
				}
			}
			as = nil
		}
	}
	return
}

//	Searches (all LibAnimationDefs contained in AllAnimationDefLibs) for the SourceArray
//	whose Id is referenced by me, returning the first match found.
func (me RefId) ArrayInAnimationDef() *SourceArray {
	var (
		s   *Source
		def *AnimationDef
		id  = me.S()
	)
	for _, lib := range AllAnimationDefLibs {
		for _, def = range lib.M {
			for _, s = range def.Sources {
				if s.Array.Id == id {
					return &s.Array
				}
			}
		}
	}
	return nil
}

//	Calls the ArrayInAnimationDef(), ArrayInControllerDef() and ArrayInGeometryDef() methods in that order to find srcArr.
func (me RefId) ArrayInAnyDef() (srcArr *SourceArray) {
	if srcArr = me.ArrayInAnimationDef(); srcArr == nil {
		if srcArr = me.ArrayInControllerDef(); srcArr == nil {
			srcArr = me.ArrayInGeometryDef()
		}
	}
	return
}

//	Searches (all LibControllerDefs contained in AllControllerDefLibs) for the SourceArray
//	whose Id is referenced by me, returning the first match found.
func (me RefId) ArrayInControllerDef() *SourceArray {
	var (
		s   *Source
		cs  Sources
		def *ControllerDef
		id  = me.S()
	)
	for _, lib := range AllControllerDefLibs {
		for _, def = range lib.M {
			if cs = nil; def.Morph != nil {
				cs = def.Morph.Sources
			} else if def.Skin != nil {
				cs = def.Skin.Sources
			}
			if cs != nil {
				for _, s = range cs {
					if s.Array.Id == id {
						return &s.Array
					}
				}
			}
		}
	}
	return nil
}

//	Searches (all LibGeometryDefs contained in AllGeometryDefLibs) for the SourceArray
//	whose Id is referenced by me, returning the first match found.
func (me RefId) ArrayInGeometryDef() (sa *SourceArray) {
	var (
		gbc *GeometryBrepCurve
		gbs *GeometryBrepSurface
		def *GeometryDef
		s   *Source
		id  = me.S()
		ret = func(cs Sources) bool {
			for _, s = range cs {
				if s.Array.Id == id {
					sa = &s.Array
					return true
				}
			}
			return false
		}
	)
	for _, lib := range AllGeometryDefLibs {
		for _, def = range lib.M {
			if (def.Mesh != nil) && ret(def.Mesh.Sources) {
				return
			}
			if (def.Spline != nil) && ret(def.Spline.Sources) {
				return
			}
			if def.Brep != nil {
				if ret(def.Brep.Sources) {
					return
				}
				if def.Brep.Surfaces != nil {
					for _, gbs = range def.Brep.Surfaces.All {
						if (gbs.Element.NurbsSurface != nil) && ret(gbs.Element.NurbsSurface.Sources) {
							return
						} else if (gbs.Element.SweptSurface != nil) && (gbs.Element.SweptSurface.Curve != nil) && (gbs.Element.SweptSurface.Curve.Element.Nurbs != nil) && ret(gbs.Element.SweptSurface.Curve.Element.Nurbs.Sources) {
							return
						}
					}
				}
				if def.Brep.Curves != nil {
					for _, gbc = range def.Brep.Curves.All {
						if (gbc.Element.Nurbs != nil) && ret(gbc.Element.Nurbs.Sources) {
							return
						}
					}
				}
				if def.Brep.SurfaceCurves != nil {
					for _, gbc = range def.Brep.SurfaceCurves.All {
						if (gbc.Element.Nurbs != nil) && ret(gbc.Element.Nurbs.Sources) {
							return
						}
					}
				}
			}
		}
	}
	return
}

//	Searches (all LibFxEffectDefs contained in AllFxEffectDefLibs) for the FxProfile
//	whose Id is referenced by me, returning the first match found.
func (me RefId) FxProfile() (fp *FxProfile) {
	var (
		def *FxEffectDef
		id  = me.S()
	)
	for _, lib := range AllFxEffectDefLibs {
		for _, def = range lib.M {
			for _, fp = range def.Profiles {
				if fp.Id == id {
					return
				}
			}
			fp = nil
		}
	}
	return
}

//	Searches (all LibFxEffectDefs contained in AllFxEffectDefLibs) for the FxTechniqueCommon
//	whose Id is referenced by me, returning the first match found.
func (me RefId) FxTechniqueCommon() *FxTechniqueCommon {
	var (
		def *FxEffectDef
		fp  *FxProfile
		id  = me.S()
	)
	for _, lib := range AllFxEffectDefLibs {
		for _, def = range lib.M {
			for _, fp = range def.Profiles {
				if (fp.Common != nil) && (fp.Common.Technique.Id == id) {
					return &fp.Common.Technique
				}
			}
		}
	}
	return nil
}

//	Searches (all LibFxEffectDefs contained in AllFxEffectDefLibs) for the FxTechniqueGlsl
//	whose Id is referenced by me, returning the first match found.
func (me RefId) FxTechniqueGlsl() (t *FxTechniqueGlsl) {
	var (
		def *FxEffectDef
		fp  *FxProfile
		id  = me.S()
	)
	for _, lib := range AllFxEffectDefLibs {
		for _, def = range lib.M {
			for _, fp = range def.Profiles {
				if fp.Glsl != nil {
					for _, t = range fp.Glsl.Techniques {
						if t.Id == id {
							return
						}
					}
					t = nil
				}
			}
		}
	}
	return
}

//	Searches (all LibGeometryDefs contained in AllGeometryDefLibs) for the GeometryBrepEdges
//	whose Id is referenced by me, returning the first match found.
func (me RefId) GeometryBrepEdges() *GeometryBrepEdges {
	var (
		def *GeometryDef
		id  = me.S()
	)
	for _, lib := range AllGeometryDefLibs {
		for _, def = range lib.M {
			if (def.Brep != nil) && (def.Brep.Edges != nil) && (def.Brep.Edges.Id == id) {
				return def.Brep.Edges
			}
		}
	}
	return nil
}

//	Searches (all LibGeometryDefs contained in AllGeometryDefLibs) for the GeometryBrepFaces
//	whose Id is referenced by me, returning the first match found.
func (me RefId) GeometryBrepFaces() *GeometryBrepFaces {
	var (
		def *GeometryDef
		id  = me.S()
	)
	for _, lib := range AllGeometryDefLibs {
		for _, def = range lib.M {
			if (def.Brep != nil) && (def.Brep.Faces != nil) && (def.Brep.Faces.Id == id) {
				return def.Brep.Faces
			}
		}
	}
	return nil
}

//	Searches (all LibGeometryDefs contained in AllGeometryDefLibs) for the GeometryBrepPcurves
//	whose Id is referenced by me, returning the first match found.
func (me RefId) GeometryBrepPcurves() *GeometryBrepPcurves {
	var (
		def *GeometryDef
		id  = me.S()
	)
	for _, lib := range AllGeometryDefLibs {
		for _, def = range lib.M {
			if (def.Brep != nil) && (def.Brep.Pcurves != nil) && (def.Brep.Pcurves.Id == id) {
				return def.Brep.Pcurves
			}
		}
	}
	return nil
}

//	Searches (all LibGeometryDefs contained in AllGeometryDefLibs) for the GeometryBrepShells
//	whose Id is referenced by me, returning the first match found.
func (me RefId) GeometryBrepShells() *GeometryBrepShells {
	var (
		def *GeometryDef
		id  = me.S()
	)
	for _, lib := range AllGeometryDefLibs {
		for _, def = range lib.M {
			if (def.Brep != nil) && (def.Brep.Shells != nil) && (def.Brep.Shells.Id == id) {
				return def.Brep.Shells
			}
		}
	}
	return nil
}

//	Searches (all LibGeometryDefs contained in AllGeometryDefLibs) for the GeometryBrepSolids
//	whose Id is referenced by me, returning the first match found.
func (me RefId) GeometryBrepSolids() *GeometryBrepSolids {
	var (
		def *GeometryDef
		id  = me.S()
	)
	for _, lib := range AllGeometryDefLibs {
		for _, def = range lib.M {
			if (def.Brep != nil) && (def.Brep.Solids != nil) && (def.Brep.Solids.Id == id) {
				return def.Brep.Solids
			}
		}
	}
	return nil
}

//	Searches (all LibGeometryDefs contained in AllGeometryDefLibs) for the GeometryDef
//	whose Id is referenced by me, returning the Mesh of the first match found.
func (me RefId) GeometryMesh() (gm *GeometryMesh) {
	if def := me.GeometryDef(); def != nil {
		gm = def.Mesh
	}
	return
}

//	Searches (all LibGeometryDefs contained in AllGeometryDefLibs) for the GeometryVertices
//	whose Id is referenced by me, returning the first match found.
func (me RefId) GeometryVertices() *GeometryVertices {
	var (
		def *GeometryDef
		id  = me.S()
	)
	for _, lib := range AllGeometryDefLibs {
		for _, def = range lib.M {
			if (def.Mesh != nil) && (def.Mesh.Vertices != nil) && (def.Mesh.Vertices.Id == id) {
				return def.Mesh.Vertices
			} else if (def.Brep != nil) && (def.Brep.Vertices.Id == id) {
				return &def.Brep.Vertices
			}
		}
	}
	return nil
}

//	Searches (all LibGeometryDefs contained in AllGeometryDefLibs) for the GeometryBrepWires
//	whose Id is referenced by me, returning the first match found.
func (me RefId) GeometryBrepWires() *GeometryBrepWires {
	var (
		def *GeometryDef
		id  = me.S()
	)
	for _, lib := range AllGeometryDefLibs {
		for _, def = range lib.M {
			if (def.Brep != nil) && (def.Brep.Wires != nil) && (def.Brep.Wires.Id == id) {
				return def.Brep.Wires
			}
		}
	}
	return nil
}

//	Returns the Id currently referenced by me.
func (me RefId) S() string {
	return string(me)
}

//	Modifies the Id currently referenced by me.
func (me *RefId) SetIdRef(v string) {
	*me = RefId(v)
}

//	Searches (all LibAnimationDefs contained in AllAnimationDefLibs) for the Source
//	whose Id is referenced by me, returning the first match found.
func (me RefId) SourceInAnimationDef() (s *Source) {
	var (
		def *AnimationDef
		id  = me.S()
	)
	for _, lib := range AllAnimationDefLibs {
		for _, def = range lib.M {
			if s = def.Sources[id]; s != nil {
				return
			}
		}
	}
	return
}

//	Calls the SourceInAnimationDef(), SourceInControllerDef() and SourceInGeometryDef() methods in that order to find src.
func (me RefId) SourceInAnyDef() (src *Source) {
	if src = me.SourceInAnimationDef(); src == nil {
		if src = me.SourceInControllerDef(); src == nil {
			src = me.SourceInGeometryDef()
		}
	}
	return
}

//	Searches (all LibControllerDefs contained in AllControllerDefLibs) for the Source
//	whose Id is referenced by me, returning the first match found.
func (me RefId) SourceInControllerDef() (s *Source) {
	var (
		def *ControllerDef
		cs  Sources
		id  = me.S()
	)
	for _, lib := range AllControllerDefLibs {
		for _, def = range lib.M {
			if cs = nil; def.Morph != nil {
				cs = def.Morph.Sources
			} else if def.Skin != nil {
				cs = def.Skin.Sources
			}
			if cs != nil {
				if s = cs[id]; s != nil {
					return
				}
			}
		}
	}
	return
}

//	Searches (all LibGeometryDefs contained in AllGeometryDefLibs) for the Source
//	whose Id is referenced by me, returning the first match found.
func (me RefId) SourceInGeometryDef() (s *Source) {
	var (
		gbc *GeometryBrepCurve
		gbs *GeometryBrepSurface
		def *GeometryDef
		id  = me.S()
	)
	for _, lib := range AllGeometryDefLibs {
		for _, def = range lib.M {
			if def.Mesh != nil {
				if s = def.Mesh.Sources[id]; s != nil {
					return
				}
			}
			if def.Spline != nil {
				if s = def.Spline.Sources[id]; s != nil {
					return
				}
			}
			if def.Brep != nil {
				if s = def.Brep.Sources[id]; s != nil {
					return
				}
				if def.Brep.Surfaces != nil {
					for _, gbs = range def.Brep.Surfaces.All {
						if gbs.Element.NurbsSurface != nil {
							if s = gbs.Element.NurbsSurface.Sources[id]; s != nil {
								return
							}
						} else if (gbs.Element.SweptSurface != nil) && (gbs.Element.SweptSurface.Curve != nil) && (gbs.Element.SweptSurface.Curve.Element.Nurbs != nil) {
							if s = gbs.Element.SweptSurface.Curve.Element.Nurbs.Sources[id]; s != nil {
								return
							}
						}
					}
				}
				if def.Brep.Curves != nil {
					for _, gbc = range def.Brep.Curves.All {
						if gbc.Element.Nurbs != nil {
							if s = gbc.Element.Nurbs.Sources[id]; s != nil {
								return
							}
						}
					}
				}
				if def.Brep.SurfaceCurves != nil {
					for _, gbc = range def.Brep.SurfaceCurves.All {
						if gbc.Element.Nurbs != nil {
							if s = gbc.Element.Nurbs.Sources[id]; s != nil {
								return
							}
						}
					}
				}
			}
		}
	}
	return
}
