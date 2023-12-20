package collection

import (
	"path/filepath"
	"regexp"
)

const (
	collectionPath  = "collection"
	fileNamePattern = `^[0-9A-Za-z\-_]+$`
)

var fileNameValidator = regexp.MustCompile(fileNamePattern)

func CollectionPath() string {
	return collectionPath
}

func CollectionFile(id string) string {
	return filepath.Join(collectionPath, id)
}
