package resource

import (
	"os"
	"path/filepath"

	"meowyplayerserver.com/utility/assert"
)

const (
	collectionFolderPath = "collection"
	collectionFileName   = "collection.json"
)

func CollectionPath() string {
	return collectionFolderPath
}

func CollectionFilePath(userName string) string {
	return filepath.Join(collectionFolderPath, userName, collectionFileName)
}

func MakeNecessaryPath() {
	assert.NoErr(os.MkdirAll(collectionFolderPath, 0777), "failed to create collection folder")
}
