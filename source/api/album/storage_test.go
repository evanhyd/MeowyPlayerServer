package album

import (
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"
)

func makeStubAlbumStorage(t *testing.T) albumStorage {
	storage := makeStorage()
	storage.albumDir = filepath.Join(t.TempDir(), storage.albumDir)
	return storage
}

func TestAlbumStorage_UploadAndDownload(t *testing.T) {
	s := makeStubAlbumStorage(t)
	if err := s.initialize(); err != nil {
		t.Fatal(err)
	}
	album1 := Album{
		key:   AlbumKey(uuid.NewString()),
		date:  time.Now().Round(time.Second),
		title: "owo",
		music: []Music{{title: "m1"}, {title: "m2"}},
		cover: []byte("image cover"),
	}

	if err := s.syncUploadImpl(album1); err != nil {
		t.Fatal(err)
	}

	album2, err := s.syncDownloadImpl(album1.Key())
	if err != nil {
		t.Fatal(err)
	}

	if album1.key != album2.key {
		t.Errorf("actual key %+v, expected %+v", album2.key, album1.key)
	}
	if album1.Date() != album2.Date() {
		t.Errorf("actual date %+v, expected %+v", album2.Date(), album1.Date())
	}
	if album1.title != album2.title {
		t.Errorf("actual title %+v, expected %+v", album2.title, album1.title)
	}
	if !slices.EqualFunc(album1.music, album2.music, func(m1, m2 Music) bool {
		//comparing date is not reliable
		return m1.title == m2.title && m1.Length() == m2.Length() && m1.platform == m2.platform && m1.id == m2.id
	}) {
		t.Errorf("actual music %+v, expected %+v", album2.music, album1.music)
	}
	if !slices.Equal(album1.cover, album2.cover) {
		t.Errorf("actual cover %+v, expected %+v", album2.cover, album1.cover)
	}
}

func TestAlbumStorage_DownloadNonExistedAlbum(t *testing.T) {
	s := makeStubAlbumStorage(t)
	if err := s.initialize(); err != nil {
		t.Fatal(err)
	}

	album, err := s.syncDownloadImpl("nonExistedKey")
	if err == nil {
		t.Errorf("actual syncDownloadImpl %+v, expected %+v", album, "error")
	}
}

func TestAlbumStorage_DownloadEmptyAlbum(t *testing.T) {
	s := makeStubAlbumStorage(t)
	if err := s.initialize(); err != nil {
		t.Fatal(err)
	}

	album, err := s.syncDownloadImpl("")
	if err == nil {
		t.Errorf("actual syncDownloadImpl %+v, expected %+v", album, "error")
	}
}
