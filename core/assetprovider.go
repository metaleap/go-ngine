package core

import (
	"io"

	xmlx "github.com/jteeuwen/go-pkg-xmlx"
)

type AssetProvider struct {
	xml *xmlx.Document
}

	func newAssetProvider (reader io.Reader) (prov *AssetProvider, err error) {
		prov = &AssetProvider {}
		prov.xml = xmlx.New()
		if err = prov.xml.LoadStream(reader, nil); err != nil { prov = nil; return }
		return
	}
