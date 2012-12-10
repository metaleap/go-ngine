package collimp

import (
	xmlx "github.com/jteeuwen/go-pkg-xmlx"

	nga "github.com/go3d/go-ngine/assets"
	c141 "github.com/go3d/go-ngine/assets/collada/conv-1.4.1-to-1.5"
	// ugfx "github.com/metaleap/go-util/gfx"
)

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

func xsn(xn *xmlx.Node, name string) (sn *xmlx.Node) {
	sn = xn.SelectNode("", name)
	return
}

func xsns(xn *xmlx.Node, name string) []*xmlx.Node {
	return xn.SelectNodes("", name)
}

//	Imports the specified Collada document, using the import options specified in importBag.
func ImportCollada(colladaDoc []byte, importBag *ImportBag) (err error) {
	state = &importState{curAssetUnitInMeters: 1}
	c141.Force, c141.Strict = false, false
	if state.doc, err = c141.ConvertDoc(colladaDoc); err == nil {
		xn := state.doc.Root.Children[0]
		sg := new(nga.SceneGraph)
		sg.Asset = obj_Asset(xn, "asset")
		sg.Extras = objs_Extra(xn, "extra")
		sg.Scene = obj_Scene(xn, "scene")
	}
	return
}
