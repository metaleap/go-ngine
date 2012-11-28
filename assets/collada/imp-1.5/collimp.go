package collimp

import (
	c14 "github.com/metaleap/go-xsd-pkg/khronos.org/files/collada_schema_1_4_go"
	c15 "github.com/metaleap/go-xsd-pkg/khronos.org/files/collada_schema_1_5_go"
)

var (
	bag *ImportBag
)

//	Contains the Collada documents to import, and provides options for the import.
type ImportBag struct {
	//	The Collada 1.4.1 document to import, if any.
	C14 *c14.TxsdCollada
	//	The Collada 1.5 document to import, if any.
	C15 *c15.TxsdCollada

	curAssetUnitInMeters float64
}

func NewImportBag() (bag *ImportBag) {
	bag = &ImportBag{}
	return
}

//	Imports the Collada documents in the specified importBag, using its current options.
//	importBag must not be nil, but its document fields (C14 and C15) may be.
func ImportCollada(importBag *ImportBag) {
	var ()
	bag = importBag
	if bag.C14 != nil {
		c14.WalkHandlers.TxsdAsset = c14_TxsdAsset
		c14.WalkHandlers.TxsdCamera = c14_TxsdCamera
		c14.WalkHandlers.TxsdImage = c14_TxsdImage
		bag.C14.Walk()
	}
	if bag.C15 != nil {
		c15.WalkHandlers.TassetType = c15_TassetType
		c15.WalkHandlers.TcameraType = c15_TcameraType
		c15.WalkHandlers.TimageType = c15_TimageType
		bag.C15.Walk()
	}
}
