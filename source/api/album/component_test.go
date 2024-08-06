package album

import (
	"testing"

	"github.com/google/uuid"
)

func makeStubAlbumComponent(t *testing.T) Component {
	comp := MakeComponent()
	comp.storage = makeStubAlbumStorage(t)
	return comp
}

func TestUploadGoodAlbumKey(t *testing.T) {
	comp := makeStubAlbumComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	album1 := Album{key: AlbumKey(uuid.NewString())}
	if comp.Upload(album1) != nil {
		t.Fatal("Upload() err is not nil, expected nil")
	}
}

func TestUploadBadAlbumKey(t *testing.T) {
	comp := makeStubAlbumComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	album1 := Album{key: AlbumKey("bad-key-obviously")}
	if comp.Upload(album1) == nil {
		t.Fatal("Upload() err is nil, expected not nil")
	}
}

func TestUploadEmptyAlbumKey(t *testing.T) {
	comp := makeStubAlbumComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	album1 := Album{}
	if comp.Upload(album1) == nil {
		t.Fatal("Upload() err is nil, expected not nil")
	}
}

func TestDownloadGoodAlbumKey(t *testing.T) {
	comp := makeStubAlbumComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	album1 := makeGoodAlbum()
	if err := comp.Upload(album1); err != nil {
		t.Fatal(err)
	}
	album2, err := comp.Download(album1.Key())
	if err != nil {
		t.Fatal(err)
	}
	expectAlbumEqual(t, &album1, &album2)
}

func TestDownloadEmptyAlbumKey(t *testing.T) {
	comp := makeStubAlbumComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	if _, err := comp.Download(""); err == nil {
		t.Fatal("Download() err is nil, expected not nil")
	}
}
