package resource

import (
	"os"
	"path/filepath"

	"meowyplayerserver.com/utility/logger"
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
	if err := os.MkdirAll(CollectionPath(), 0777); err != nil {
		logger.Error(err, 0)
	}
}
