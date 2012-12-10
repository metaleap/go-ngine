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
