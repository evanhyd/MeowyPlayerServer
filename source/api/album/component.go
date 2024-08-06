package album

import (
	"fmt"

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

func (c *Component) Upload(album Album) error {
	if !c.isValidAlbumKey(album.Key()) {
		return fmt.Errorf("invalid album key %v", album.Key())
	}
	return c.storage.upload(album)
}

func (c *Component) Download(key AlbumKey) (Album, error) {
	if !c.isValidAlbumKey(key) {
		return Album{}, fmt.Errorf("invalid album key %v", key)
	}
	return c.storage.download(key)
}
