package album

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type albumStorage struct {
	albumDir string
	requests chan func()
}

func makeStorage() albumStorage {
	const (
		kAlbumDir         = "album"
		kRequestQueueSize = 128
	)
	return albumStorage{albumDir: kAlbumDir, requests: make(chan func(), kRequestQueueSize)}
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

func (s *albumStorage) getPath(key AlbumKey) string {
	return filepath.Join(s.albumDir, fmt.Sprintf("%v.json", key))
}

func (s *albumStorage) store(album Album) error {
	data, err := json.Marshal(&album)
	if err != nil {
		return err
	}
	return os.WriteFile(s.getPath(album.key), data, 0600)
}

func (s *albumStorage) load(key AlbumKey) (Album, error) {
	data, err := os.ReadFile(s.getPath(key))
	if err != nil {
		return Album{}, err
	}

	var album Album
	err = json.Unmarshal(data, &album)
	return album, err
}

func (s *albumStorage) upload(album Album) error {
	respC := make(chan error)
	s.requests <- func() {
		respC <- s.store(album)
	}
	return <-respC
}

func (s *albumStorage) download(key AlbumKey) (Album, error) {
	var album Album
	var err error
	readyC := make(chan struct{})
	s.requests <- func() {
		album, err = s.load(key)
		readyC <- struct{}{}
	}
	<-readyC
	return album, err
}
