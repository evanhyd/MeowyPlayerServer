package album

import (
	"encoding/json"
	"fmt"
	"meowyplayerserver/api/account"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type albumStorage struct {
	albumDir string
	requests chan func()
	albums   map[uuid.UUID][]Album //id: a set of album
}

func makeStorage() albumStorage {
	const (
		kAlbumDir         = "album"
		kRequestQueueSize = 32
	)
	return albumStorage{albumDir: kAlbumDir, requests: make(chan func(), kRequestQueueSize), albums: map[uuid.UUID][]Album{}}
}

func (s *albumStorage) initialize() error {
	if err := os.MkdirAll(s.albumDir, 0700); err != nil {
		return err
	}
	go func() {
		for req := range s.requests {
			req()
		}
	}()
	return nil
}

func (s *albumStorage) userIDPath(userID account.UserID) string {
	return filepath.Join(s.albumDir, userID.String())
}

func (s *albumStorage) albumPath(userID account.UserID, key AlbumKey) string {
	return filepath.Join(s.userIDPath(userID), fmt.Sprintf("%v.json", key))
}

func (s *albumStorage) store(userID account.UserID, album Album) error {
	data, err := json.Marshal(album)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(s.userIDPath(userID), 0700); err != nil {
		return err
	}
	return os.WriteFile(s.albumPath(userID, album.Key()), data, 0600)
}

func (s *albumStorage) load(userID account.UserID, key AlbumKey) (Album, error) {
	data, err := os.ReadFile(s.albumPath(userID, key))
	if err != nil {
		return Album{}, err
	}

	var album Album
	err = json.Unmarshal(data, &album)
	return album, err
}

func (s *albumStorage) upload(userID account.UserID, album Album) error {
	respC := make(chan error)
	s.requests <- func() {
		respC <- s.store(userID, album)
	}
	return <-respC
}

func (s *albumStorage) download(userID account.UserID, key AlbumKey) (Album, error) {
	var album Album
	var err error
	readyC := make(chan struct{})
	s.requests <- func() {
		album, err = s.load(userID, key)
		readyC <- struct{}{}
	}
	<-readyC
	return album, err
}
