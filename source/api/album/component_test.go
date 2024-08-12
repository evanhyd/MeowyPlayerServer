package album

import (
	"meowyplayerserver/api/account"
	"testing"
)

func makeStubAlbumComponent(t *testing.T) Component {
	comp := MakeComponent()
	comp.storage = makeStubAlbumStorage(t)
	return comp
}

func TestUploadAlbumKey(t *testing.T) {
	comp := makeStubAlbumComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	userID := account.NewUserID()
	album1 := makeGoodAlbum()
	if err := comp.Upload(userID, album1); err != nil {
		t.Fatal("Upload() err is not nil, expected nil", err)
	}
}

func TestDownloadAlbumKey(t *testing.T) {
	comp := makeStubAlbumComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	userID := account.NewUserID()
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
