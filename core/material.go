package core

import (
	"log"
	"path/filepath"
)

type TMaterial struct {
	texKey string
}

func NewMaterialFromLocalTextureImageFile (filePath string) *TMaterial {
	var err error
	var tex = &TTexture {}
	filePath = filepath.Join(Core.Options.AssetRootDirPath, filePath)
	tex.Params = Core.Options.DefaultTextureParams
	if true {
		tex.LoadAsync(TextureProviders.LocalFile, filePath)
	} else if err = tex.Load(TextureProviders.LocalFile, filePath); err != nil {
		log.Printf("ERROR loading texture %v: %v\n", filePath, err)
	}
	Core.Textures[filePath] = tex
	return NewMaterial(filePath)
}

func NewMaterialFromRemoteTextureImageFile (fileUrl string) *TMaterial {
	var err error
	var tex = &TTexture {}
	tex.Params = Core.Options.DefaultTextureParams
	if true {
		tex.LoadAsync(TextureProviders.RemoteFile, fileUrl)
	} else if err = tex.Load(TextureProviders.RemoteFile, fileUrl); err != nil {
		log.Printf("ERROR loading texture %v: %v\n", fileUrl, err)
	}
	Core.Textures[fileUrl] = tex
	return NewMaterial(fileUrl)
}

func NewMaterial (texKey string) *TMaterial {
	var mat = &TMaterial {}
	mat.texKey = texKey
	return mat
}
