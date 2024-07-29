package user

import (
	"testing"
)

func newStubAccountManager(t *testing.T) *userManager {
	acc := NewUserManager()
	acc.storage = newStubAccountStorage(t)
	return acc
}

func TestManager_Register(t *testing.T) {
	manager := newStubAccountManager(t)
	if err := manager.initialize(); err != nil {
		t.Fatal(err)
	}

	if err := manager.register("UnboxTheCat", "test"); err != nil {
		t.Fatal(err)
	}

	if err := manager.register("UnboxTheCat", "test"); err == nil {
		t.Fatal(err)
	}

	if err := manager.register("", "test"); err == nil {
		t.Fatal(err)
	}
}

func TestManager_Authorize(t *testing.T) {
	manager := newStubAccountManager(t)
	if err := manager.initialize(); err != nil {
		t.Fatal(err)
	}

	if err := manager.register("UnboxTheCat", "test"); err != nil {
		t.Fatal(err)
	}

	if !manager.authorize("UnboxTheCat", "test") {
		t.Fatalf("authorize() = false, wanted true")
	}

	if manager.authorize("UnboxTheCat", "test1") {
		t.Fatalf("authorize() = true, wanted false")
	}

	if manager.authorize("not_exist", "abc") {
		t.Fatalf("authorize() = true, wanted false")
	}

	if manager.authorize("", "") {
		t.Fatalf("authorize() = true, wanted false")
	}
}

func TestManager_IsValidUsername(t *testing.T) {
	manager := newStubAccountManager(t)
	if err := manager.initialize(); err != nil {
		t.Fatal(err)
	}

	badUsername := []string{
		"",        //no name
		" ",       //space
		"(abc",    //bad prefix
		"abc)",    //bad suffix
		"(abc)",   //bad pre/suffix
		"abc(123", //bad character in the middle
	}

	goodUsername := []string{
		"abc",         //letter only
		"123",         //digit only
		"abc123",      //letter + digit
		"hello_world", //underscore
		"hello-world", //dash
		"UnboxTheCat", //my name
		"__MAIN__",    //underscore pre/suffix
		"-_-",         //-_-
	}

	for _, username := range badUsername {
		if manager.isValidUsername(username) {
			t.Errorf("isValidUser(%v) = true, wanted false", username)
		}
	}

	for _, username := range goodUsername {
		if !manager.isValidUsername(username) {
			t.Errorf("isValidUser(%v) = false, wanted true", username)
		}
	}
}
