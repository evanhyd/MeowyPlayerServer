package album

import (
	"fmt"
	"meowyplayerserver/api/account"
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
	return true
}

func (c *Component) Upload(userID account.UserID, album Album) error {
	if !c.isValidAlbumKey(album.Key()) {
		return fmt.Errorf("invalid album key %v", album.Key())
	}
	return c.storage.upload(userID, album)
}

func (c *Component) Download(userID account.UserID, key AlbumKey) (Album, error) {
	if !c.isValidAlbumKey(key) {
		return Album{}, fmt.Errorf("invalid album key %v", key)
	}
	return c.storage.download(userID, key)
}
