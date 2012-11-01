package core

import (
	"io"

	xmlx "github.com/jteeuwen/go-pkg-xmlx"
)

type TAssetProvider struct {
	xml *xmlx.Document
}

	func newAssetProvider (reader io.Reader) (prov *TAssetProvider, err error) {
		prov = &TAssetProvider {}
		prov.xml = xmlx.New()
		if err = prov.xml.LoadStream(reader, nil); err != nil { prov = nil; return }
		return
	}
