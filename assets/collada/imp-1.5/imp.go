package collimp

import (
	xmlx "github.com/jteeuwen/go-pkg-xmlx"

	// nga "github.com/go3d/go-ngine/assets"
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

//	Imports the specified Collada document, using the import options specified in importBag.
func ImportCollada(colladaDoc []byte, importBag *ImportBag) (err error) {
	state = &importState{curAssetUnitInMeters: 1}
	c141.Force, c141.Strict = false, false
	if state.doc, err = c141.ConvertDoc(colladaDoc); err == nil {
		// err = xml.Unmarshal(colladaDoc, state.d15)
		// if colladaDoc = nil; err == nil {
		// 	c15.WalkHandlers.TcameraType = c15_TcameraType
		// 	c15.WalkHandlers.TeffectType = c15_TeffectType
		// 	c15.WalkHandlers.TimageType = c15_TimageType
		// 	c15.WalkHandlers.TmaterialType = c15_TmaterialType
		// 	c15.WalkHandlers.TxsdCollada = c15_TxsdCollada
		// 	c15.WalkOnError = func(err error) {
		// 		log.Printf("%v\n", err)
		// 	}
		// 	c15.WalkContinueOnError = true
		// 	state.d15.Walk()
		// }
	}
	return
}
