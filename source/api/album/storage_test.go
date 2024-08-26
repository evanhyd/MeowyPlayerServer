package album

import (
	"meowyplayerserver/api/account"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"
)

func makeStubAlbumStorage(t *testing.T) albumStorage {
	storage := makeStorage()
	storage.albumDir = filepath.Join(t.TempDir(), storage.albumDir)
	return storage
}

func makeGoodAlbum() Album {
	return Album{
		key:   newRandomAlbumKey(),
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

	userID := account.NewUserID()
	if err := s.allocateStorage(userID); err != nil {
		t.Fatal(err)
	}
	album1 := makeGoodAlbum()
	if err := s.uploadAlbum(userID, album1); err != nil {
		t.Fatal(err)
	}

	album2, err := s.getAlbum(userID, album1.Key())
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

	userID := account.NewUserID()
	if err := s.allocateStorage(userID); err != nil {
		t.Fatal(err)
	}

	albumKey := newRandomAlbumKey()
	_, err := s.getAlbum(userID, albumKey)
	if err == nil {
		t.Errorf("getAlbum expected error")
	}
}

func TestAlbumStorage_DownloadAllAlbumsEmpty(t *testing.T) {
	s := makeStubAlbumStorage(t)
	if err := s.initialize(); err != nil {
		t.Fatal(err)
	}

	userID := account.NewUserID()
	if err := s.allocateStorage(userID); err != nil {
		t.Fatal(err)
	}

	albums, err := s.getAllAlbums(userID)
	if err != nil {
		t.Fatal(err)
	}

	if l := len(albums); l != 0 {
		t.Errorf("len(albums) = %v, wanted 0", l)
	}
}

func TestAlbumStorage_DownloadAllAlbumsTwo(t *testing.T) {
	s := makeStubAlbumStorage(t)
	if err := s.initialize(); err != nil {
		t.Fatal(err)
	}

	userID := account.NewUserID()
	if err := s.allocateStorage(userID); err != nil {
		t.Fatal(err)
	}

	albums := []Album{makeGoodAlbum(), makeGoodAlbum()}
	slices.SortFunc(albums, func(lhs, rhs Album) int {
		return strings.Compare(lhs.Key().String(), rhs.Key().String())
	})

	if err := s.uploadAlbum(userID, albums[0]); err != nil {
		t.Fatal(err)
	}
	if err := s.uploadAlbum(userID, albums[1]); err != nil {
		t.Fatal(err)
	}

	albums1, err := s.getAllAlbums(userID)
	if err != nil {
		t.Errorf("getAllAlbum received error %v, wanted nil", err)
	}

	if len(albums) != len(albums1) {
		t.Errorf("len(albums1) = %v, wanted %v", len(albums1), len(albums))
	}
	expectAlbumEqual(t, &albums[0], &albums1[0])
	expectAlbumEqual(t, &albums[1], &albums1[1])
}
