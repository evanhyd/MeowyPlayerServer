package album

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

type Component struct {
	storage albumStorage
}

func MakeComponent() Component {
	return Component{storage: makeStorage()}
}

func (c *Component) Initialize() error {
	return c.storage.initialize()
}

func (c *Component) isValidAlbumKey(key AlbumKey) bool {
	return uuid.Validate(string(key)) == nil
}

func (c *Component) Upload(w http.ResponseWriter, album Album) error {
	if !c.isValidAlbumKey(album.Key()) {
		return fmt.Errorf("nice try")
	}
	return c.storage.upload(w, album)
}

func (c *Component) Download(w http.ResponseWriter, key AlbumKey) error {
	if !c.isValidAlbumKey(key) {
		return fmt.Errorf("you thought")
	}
	return c.storage.download(w, key)
}
