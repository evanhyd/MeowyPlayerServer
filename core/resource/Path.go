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

func CollectionPath(collectionName string) string {
	return filepath.Join(collectionFolderPath, collectionName, collectionFileName)
}

func MakeNecessaryPath() {
	assert.NoErr(os.MkdirAll(collectionFolderPath, 0777), "failed to create collection folder")
}
