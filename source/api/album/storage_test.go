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

func makeGoodAlbum() Album {
	return Album{
		key:   AlbumKey(uuid.NewString()),
		date:  time.Now().Round(time.Second),
		title: "owo",
		music: []Music{{
			date:     time.Now().Round(time.Second),
			title:    "title 1",
			length:   10 * time.Second,
			platform: "YouTube",
			id:       "123i1r101",
		}, {
			date:     time.Now().Round(time.Second),
			title:    "title 2",
			length:   1 * time.Minute,
			platform: "BiliBili",
			id:       "abcdef",
		}},
		cover: []byte("image cover"),
	}
}

func expectAlbumEqual(t *testing.T, album1 *Album, album2 *Album) {
	if album1.Key() != album2.Key() {
		t.Errorf("actual key %+v, expected %+v", album2.Key(), album1.Key())
	}
	if album1.Date() != album2.Date() {
		t.Errorf("actual date %+v, expected %+v", album2.Date(), album1.Date())
	}
	if album1.Title() != album2.Title() {
		t.Errorf("actual title %+v, expected %+v", album2.Title(), album1.Title())
	}
	if !slices.EqualFunc(album1.Music(), album2.Music(), func(m1, m2 Music) bool {
		return m1.Date() == m2.Date() && m1.Title() == m2.Title() && m1.Length() == m2.Length() && m1.Key() == m2.Key() && m1.platform == m2.platform && m1.id == m2.id
	}) {
		t.Errorf("actual music %+v, expected %+v", album2.Music(), album1.Music())
	}
	if !slices.Equal(album1.Cover(), album2.Cover()) {
		t.Errorf("actual cover %+v, expected %+v", album2.Cover(), album1.Cover())
	}
}

func TestAlbumStorage_UploadAndDownload(t *testing.T) {
	s := makeStubAlbumStorage(t)
	if err := s.initialize(); err != nil {
		t.Fatal(err)
	}

	album1 := makeGoodAlbum()
	if err := s.upload(album1); err != nil {
		t.Fatal(err)
	}

	album2, err := s.download(album1.Key())
	if err != nil {
		t.Fatal(err)
	}
	expectAlbumEqual(t, &album1, &album2)
}

func TestAlbumStorage_DownloadNonExistedAlbum(t *testing.T) {
	s := makeStubAlbumStorage(t)
	if err := s.initialize(); err != nil {
		t.Fatal(err)
	}

	album, err := s.download("nonExistedKey")
	if err == nil {
		t.Errorf("actual syncDownloadImpl %+v, expected %+v", album, "error")
	}
}

func TestAlbumStorage_DownloadEmptyAlbum(t *testing.T) {
	s := makeStubAlbumStorage(t)
	if err := s.initialize(); err != nil {
		t.Fatal(err)
	}

	album, err := s.download("")
	if err == nil {
		t.Errorf("actual syncDownloadImpl %+v, expected %+v", album, "error")
	}
}
