package account

import (
	"path/filepath"
	"testing"
)

func makeStubAccountStorage(t *testing.T) accountStorage {
	s := makeStorage()
	s.accountDir = filepath.Join(t.TempDir(), s.accountDir)
	return s
}

func TestAccountStorage_Create(t *testing.T) {
	s := makeStubAccountStorage(t)
	if err := s.initialize(); err != nil {
		t.Fatal(err)
	}

	if !s.create(Account{username: "UnboxTheCat", id: "id", salt: []byte("salt"), hash: []byte("hash")}) {
		t.Fatal("create() = false, wanted true")
	}
	if s.create(Account{username: "UnboxTheCat", id: "id", salt: []byte("salt"), hash: []byte("hash")}) {
		t.Fatal("create() = true, wanted false")
	}
}

func TestAccountStorage_Get(t *testing.T) {
	s := makeStubAccountStorage(t)
	if err := s.initialize(); err != nil {
		t.Fatal(err)
	}

	if !s.create(Account{username: "UnboxTheCat", id: "id", salt: []byte("salt"), hash: []byte("hash")}) {
		t.Fatal("create() = false, wanted true")
	}

	if _, exist := s.get("UnboxTheCat"); !exist {
		t.Fatal("get() = false, wanted true")
	}

	if _, exist := s.get("nonamer"); exist {
		t.Fatal("get() = true, wanted false")
	}
}

func TestAccountStorage_Save(t *testing.T) {
	s := makeStubAccountStorage(t)
	if err := s.initialize(); err != nil {
		t.Fatal(err)
	}

	if !s.create(Account{username: "UnboxTheCat", id: "id", salt: []byte("salt"), hash: []byte("hash")}) {
		t.Fatal("create() = false, wanted true")
	}
	if !s.create(Account{username: "Guest", id: "id1", salt: []byte("salt1"), hash: []byte("hash1")}) {
		t.Fatal("create() = false, wanted true")
	}
	if err := s.save(); err != nil {
		t.Fatal(err)
	}

	s1 := makeStubAccountStorage(t)
	s1.accountDir = s.accountDir
	if err := s1.initialize(); err != nil {
		t.Fatal(err)
	}
	if _, exist := s1.get("UnboxTheCat"); !exist {
		t.Fatal("get() = false, wanted true")
	}
	if _, exist := s1.get("Guest"); !exist {
		t.Fatal("get() = false, wanted true")
	}
}
