package collimp

import (
	xmlx "github.com/jteeuwen/go-pkg-xmlx"

	nga "github.com/go3d/go-ngine/assets"
	c141 "github.com/go3d/go-ngine/assets/collada/conv-1.4.1-to-1.5"
	// ugfx "github.com/metaleap/go-util/gfx"
)

var (
	bag   *ImportBag
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
func ImportCollada(colladaDoc []byte, importBag *ImportBag) (doc *nga.Document, err error) {
	bag = importBag
	state = &importState{curAssetUnitInMeters: 1}
	c141.Force, c141.Strict = false, false
	if state.doc, err = c141.ConvertDoc(colladaDoc); err == nil {
		xn := state.doc.Root.Children[0]
		doc = obj_Document(xn, "")
		libs_All(xn)
	}
	return
}
