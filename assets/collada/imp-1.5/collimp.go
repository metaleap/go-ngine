package collimp

import (
	"encoding/xml"
	c141 "github.com/go3d/go-ngine/assets/collada/conv-1.4.1-to-1.5"
	c15 "github.com/metaleap/go-xsd-pkg/khronos.org/files/collada_schema_1_5_go"
)

var (
	bag *ImportBag
)

//	Contains the Collada documents to import, and provides options for the import.
type ImportBag struct {
	d15                  *c15.TxsdCollada
	curAssetUnitInMeters float64
}

//	Initializes and returns a newly created ImportBag instance.
func NewImportBag() (bag *ImportBag) {
	bag = &ImportBag{}
	return
}

//	Imports the Collada documents in the specified importBag, using its current options.
func ImportCollada(data []byte, importBag *ImportBag) (err error) {
	bag = importBag
	bag.d15 = &c15.TxsdCollada{}
	c141.Force, c141.Strict = false, false
	if data, err = c141.Convert(data); err == nil {
		err = xml.Unmarshal(data, bag.d15)
		data = nil
		if err == nil {
			c15.WalkHandlers.TassetType = c15_TassetType
			c15.WalkHandlers.TcameraType = c15_TcameraType
			c15.WalkHandlers.TimageType = c15_TimageType
			bag.d15.Walk()
		}
	}
	return
}
