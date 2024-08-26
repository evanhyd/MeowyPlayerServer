package album

import (
	"meowyplayerserver/api/account"
	"slices"
	"strings"
	"testing"
)

func makeStubAlbumComponent(t *testing.T) Component {
	comp := MakeComponent()
	comp.storage = makeStubAlbumStorage(t)
	return comp
}

func TestUploadAlbum(t *testing.T) {
	comp := makeStubAlbumComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	userID := account.NewUserID()
	comp.Register(userID)
	album1 := makeGoodAlbum()
	if err := comp.Upload(userID, album1); err != nil {
		t.Fatal("Upload() err is not nil, expected nil", err)
	}
}

func TestDownloadAlbum(t *testing.T) {
	comp := makeStubAlbumComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	userID := account.NewUserID()
	comp.Register(userID)
	album1 := makeGoodAlbum()
	if err := comp.Upload(userID, album1); err != nil {
		t.Fatal(err)
	}
	album2, err := comp.Download(userID, album1.Key())
	if err != nil {
		t.Fatal(err)
	}
	expectAlbumEqual(t, &album1, &album2)
}

func TestDownloadAllAlbums(t *testing.T) {
	comp := makeStubAlbumComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	userID := account.NewUserID()
	comp.Register(userID)
	albums := []Album{makeGoodAlbum(), makeGoodAlbum()}
	slices.SortFunc(albums, func(lhs, rhs Album) int {
		return strings.Compare(lhs.Key().String(), rhs.Key().String())
	})

	if err := comp.Upload(userID, albums[0]); err != nil {
		t.Fatal(err)
	}
	if err := comp.Upload(userID, albums[1]); err != nil {
		t.Fatal(err)
	}

	albums1, err := comp.DownloadAll(userID)
	if err != nil {
		t.Fatal(err)
	}
	expectAlbumEqual(t, &albums[0], &albums1[0])
	expectAlbumEqual(t, &albums[1], &albums1[1])
}
