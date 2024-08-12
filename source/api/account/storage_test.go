package account

import (
	"path/filepath"
	"slices"
	"testing"
)

func makeStubAccountStorage(t *testing.T) accountStorage {
	s := makeStorage()
	s.accountDir = filepath.Join(t.TempDir(), s.accountDir)
	return s
}

func TestStorageReadConfig(t *testing.T) {
	s := makeStubAccountStorage(t)
	if err := s.initialize(); err != nil {
		t.Fatal(err)
	}

	accs := []Account{
		{Username: "a0", UserID: NewUserID(), Salt: []byte("b0"), Hash: []byte("c0")},
		{Username: "a1", UserID: NewUserID(), Salt: []byte("b1"), Hash: []byte("c1")},
		{Username: "a2", UserID: NewUserID(), Salt: []byte("b2"), Hash: []byte("c2")},
		{Username: "a3", UserID: NewUserID(), Salt: []byte("b3"), Hash: []byte("c3")},
	}
	for _, acc := range accs {
		s.store(acc)
	}

	//re-read the config
	s1 := makeStubAccountStorage(t)
	s1.accountDir = s.accountDir
	if err := s1.initialize(); err != nil {
		t.Fatal(err)
	}

	for _, acc := range accs {
		acc1, ok := s1.load(acc.Username)
		if !ok {
			t.Errorf("Load(%v) = false, expected true", acc.Username)
		}

		if acc.Username != acc1.Username || acc.UserID != acc1.UserID || !slices.Equal(acc.Salt, acc1.Salt) || !slices.Equal(acc.Hash, acc1.Hash) {
			t.Errorf("acc != acc1, %v %v", acc, acc1)
		}
	}
}
