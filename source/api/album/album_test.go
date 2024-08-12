package album

import (
	"testing"

	"github.com/google/uuid"
)

func TestParseAlbumKey_Good(t *testing.T) {
	id := uuid.NewString()
	if _, err := ParseAlbumKey(id); err != nil {
		t.Errorf("ParseAlbumKey(%v) emits error %v", id, err)
	}
}

func TestParseAlbumKey_Empty(t *testing.T) {
	id := ""
	if _, err := ParseAlbumKey(id); err == nil {
		t.Errorf("ParseAlbumKey(%v) passed, expected an error", id)
	}
}
