package resource

import (
	"os"

	"meowyplayerserver.com/utility/assert"
)

const (
	collectionPath = "collection"
)

func CollectionPath() string {
	return collectionPath
}

func MakeNecessaryPath() {
	assert.NoErr(os.MkdirAll(CollectionPath(), 0777), "failed to create collection path")
}
