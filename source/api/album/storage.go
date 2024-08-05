package album

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type albumStorage struct {
	albumDir string
	albums   albumMap
	requests chan func()
}

func makeStorage() albumStorage {
	const (
		kAlbumDir         = "album"
		kRequestQueueSize = 128
	)
	return albumStorage{albumDir: kAlbumDir, albums: albumMap{}, requests: make(chan func(), kRequestQueueSize)}
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

func (s *albumStorage) syncUploadImpl(album Album) error {
	data, err := json.Marshal(&album)
	if err != nil {
		return err
	}
	return os.WriteFile(s.getPath(album.key), data, 0600)
}

func (s *albumStorage) syncDownloadImpl(key AlbumKey) (album Album, err error) {
	data, err := os.ReadFile(s.getPath(key))
	if err == nil {
		err = json.Unmarshal(data, &album)
	}
	return
}

func (s *albumStorage) upload(w http.ResponseWriter, album Album) error {
	s.requests <- func() {
		if err := s.syncUploadImpl(album); err != nil {
			http.Error(w, fmt.Sprintf("[requestUpload] %v", err), http.StatusInternalServerError)
		}
	}
	return nil
}

func (s *albumStorage) download(w http.ResponseWriter, key AlbumKey) error {
	s.requests <- func() {
		album, err := s.syncDownloadImpl(key)
		if err != nil {
			http.Error(w, fmt.Sprintf("[requestDownload] %v", err), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(&album); err != nil {
			http.Error(w, fmt.Sprintf("[requestDownload] %v", err), http.StatusInternalServerError)
		}
	}
	return nil
}
