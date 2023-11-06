package resource

import (
	"os"
	"path/filepath"

	"meowyplayerserver.com/utility/assert"
)

const (
	collectionPath = "collection"
)

func CollectionPath() string {
	return collectionPath
}

func CollectionFile(collection string) string {
	return filepath.Join(collectionPath, collection)
}

func MakeNecessaryPath() {
	assert.NoErr(os.MkdirAll(CollectionPath(), 0777), "failed to create collection path")
}
