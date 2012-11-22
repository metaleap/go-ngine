package collimp

import (
	coll14 "github.com/metaleap/go-xsd-pkg/khronos.org/files/collada_schema_1_4_go"
	coll15 "github.com/metaleap/go-xsd-pkg/khronos.org/files/collada_schema_1_5_go"
	xsdt "github.com/metaleap/go-xsd/types"
)

type c14AssetStack []*coll14.TxsdAsset

	func (me *c14AssetStack) Pop (popAsset bool) (p bool, curAsset *coll14.TxsdAsset) {
		if popAsset && (len(*me) > 0) { sl := *me; *me = sl[1 :] }; if len(*me) > 0 { curAsset = (*me)[0] }; return
	}

	func (me *c14AssetStack) Push (el *coll14.TxsdAsset) (p bool, curAsset *coll14.TxsdAsset) {
		if p = (el != nil); p { curAsset = el; nu := []*coll14.TxsdAsset { el }; *me = append(nu, *me...) } else if len(*me) > 0 { curAsset = (*me)[0] }; return
	}

func f64 (f *coll14.TargetableFloat) float64 {
	return float64(f.XsdGoPkgValue.ToXsdtDouble())
}

func ImportCollada (c14 *coll14.TxsdCollada, c15 *coll15.TxsdCollada) {
	if c14 != nil {
		coll14.WalkHandlers.TxsdAsset = func (x *coll14.TxsdAsset) {
		}
		coll14.WalkHandlers.TxsdCamera = Imp14Camera
		c14.Walk()
	}
	if c15 != nil {
		coll15.WalkHandlers.TcameraType = Imp15Camera
	}
}
