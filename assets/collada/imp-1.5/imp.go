package collimp

import (
	xmlx "github.com/jteeuwen/go-pkg-xmlx"

	nga "github.com/go3d/go-ngine/assets"
	c141 "github.com/go3d/go-ngine/assets/collada/conv-1.4.1-to-1.5"
	// ugfx "github.com/metaleap/go-util/gfx"
)

type Root struct {
	nga.HasAsset
	nga.HasExtras
	Scene *nga.Scene
}

var (
	state *importState
)

//	Provides options for importing Collada documents.
type ImportBag struct {
}

//	Initializes and returns a newly created ImportBag instance.
func NewImportBag() (me *ImportBag) {
	me = &ImportBag{}
	return
}

type importState struct {
	doc                  *xmlx.Document
	curAssetUnitInMeters float64
}

func has_Asset(xn *xmlx.Node, obj *nga.HasAsset) {
	obj.Asset = obj_Asset(xn, "asset")
}

func has_Extras(xn *xmlx.Node, obj *nga.HasExtras) {
	obj.Extras = objs_Extra(xn, "extra")
}

func has_ID(xn *xmlx.Node, obj *nga.HasID) {
	obj.ID = xna(xn, "id")
}

func has_Name(xn *xmlx.Node, obj *nga.HasName) {
	obj.Name = xna(xn, "name")
}

func has_Sid(xn *xmlx.Node, obj *nga.HasSid) {
	obj.Sid = xna(xn, "sid")
}

func has_Techniques(xn *xmlx.Node, obj *nga.HasTechniques) {
	obj.Techniques = objs_Technique(xn, "technique")
}

//	Imports the specified Collada document, using the import options specified in importBag.
func ImportCollada(colladaDoc []byte, importBag *ImportBag) (root *Root, err error) {
	state = &importState{curAssetUnitInMeters: 1}
	c141.Force, c141.Strict = false, false
	if state.doc, err = c141.ConvertDoc(colladaDoc); err == nil {
		xn := state.doc.Root.Children[0]
		root = &Root{Scene: obj_Scene(xn, "scene")}
		root.Asset, root.Extras = obj_Asset(xn, "asset"), objs_Extra(xn, "extra")
		libs_All(xn)
	}
	return
}

func xna(xn *xmlx.Node, name string) string {
	return xn.As("", name)
}

func xnau64(xn *xmlx.Node, name string) uint64 {
	return xn.Au64("", name)
}

func xnau64p(xn *xmlx.Node, name string) (p *uint64) {
	for _, att := range xn.Attributes {
		if att.Name.Local == name {
			*p = xnau64(xn, name)
			break
		}
	}
	return
}

func xnf64(xn *xmlx.Node, name string) float64 {
	return xn.F64("", name)
}

func xsn(xn *xmlx.Node, name string) (sn *xmlx.Node) {
	sn = xn.SelectNode("", name)
	return
}

func xsns(xn *xmlx.Node, name string) []*xmlx.Node {
	return xn.SelectNodes("", name)
}
