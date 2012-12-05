package collimp

import (
	"encoding/xml"
	"log"

	nga "github.com/go3d/go-ngine/assets"
	c141 "github.com/go3d/go-ngine/assets/collada/conv-1.4.1-to-1.5"
	c15 "github.com/metaleap/go-xsd-pkg/khronos.org/files/collada_schema_1_5_go"
	xsdt "github.com/metaleap/go-xsd/types"
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
	d15                  *c15.TxsdCollada
	curAssetUnitInMeters float64
	assetStack           []*c15.TassetType
}

func f64(d xsdt.ToXsdtDouble) float64 {
	if d == nil {
		return 0
	}
	return float64(d.ToXsdtDouble())
}

func f64c(d xsdt.ToXsdtDouble) float64 {
	return (f64(d) * state.curAssetUnitInMeters) / nga.UnitInMeters
}

//	Imports the specified Collada document, using the import options specified in importBag.
func ImportCollada(colladaDoc []byte, importBag *ImportBag) (err error) {
	state = &importState{curAssetUnitInMeters: 1, d15: &c15.TxsdCollada{}}
	c141.Force, c141.Strict = false, false
	if colladaDoc, err = c141.Convert(colladaDoc); err == nil {
		err = xml.Unmarshal(colladaDoc, state.d15)
		if colladaDoc = nil; err == nil {
			c15.WalkHandlers.TcameraType = c15_TcameraType
			c15.WalkHandlers.TeffectType = c15_TeffectType
			c15.WalkHandlers.TimageType = c15_TimageType
			c15.WalkHandlers.TmaterialType = c15_TmaterialType
			c15.WalkHandlers.TxsdCollada = c15_TxsdCollada
			c15.WalkOnError = func(err error) {
				log.Printf("%v\n", err)
			}
			c15.WalkContinueOnError = true
			state.d15.Walk()
		}
	}
	return
}
