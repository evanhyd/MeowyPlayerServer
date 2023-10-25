package resource

import (
	"os"
	"path/filepath"

	"meowyplayerserver.com/utility/assert"
)

const (
	userFolder = "user"
)

func UserPath() string {
	return userFolder
}

func CollectionPath(userName string) string {
	return filepath.Join(userFolder, userName)
}

func MakeNecessaryPath() {
	assert.NoErr(os.MkdirAll(userFolder, 0777), "failed to create user folder")
}
