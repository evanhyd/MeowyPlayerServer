package album

import (
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

func (c *Component) Register(userID account.UserID) error {
	return c.storage.allocateStorage(userID)
}

func (c *Component) Upload(userID account.UserID, album Album) error {
	return c.storage.uploadAlbum(userID, album)
}

func (c *Component) Download(userID account.UserID, key AlbumKey) (Album, error) {
	return c.storage.getAlbum(userID, key)
}

func (c *Component) DownloadAll(userID account.UserID) ([]Album, error) {
	return c.storage.getAllAlbums(userID)
}
