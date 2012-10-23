package core

import (
	"path/filepath"
)

type tAssetManager struct {
}

func newAssetManager () *tAssetManager {
	var assMan = &tAssetManager {}
	return assMan
}

func (me *tAssetManager) LocalFilePath (assetRootRelativeFilePath string) string {
	return filepath.Join(Core.Options.AssetRootDirPath, assetRootRelativeFilePath)
}
