package authentication_test

import (
	"os"
	"testing"

	"meowyplayerserver.com/core/authentication"
)

func TestRegister(t *testing.T) {
	authentication.Initialize()
	defer os.RemoveAll(authentication.AccountFile())

	if err := authentication.RegisterAccount("UnboxTheCat", []byte("test")); err != nil {
		t.Fatal(err)
	}

	if err := authentication.RegisterAccount("UnboxTheCat", []byte("test")); err != nil {
		t.Fatal(err)
	}

	if !authentication.IsAccountExist("UnboxTheCat") {
		t.Fatal()
	}

	if authentication.IsAccountExist("") {
		t.Fatal()
	}
}

func TestPassword(t *testing.T) {
	authentication.Initialize()
	defer os.RemoveAll(authentication.AccountFile())

	if err := authentication.RegisterAccount("UnboxTheCat", []byte("test")); err != nil {
		t.Fatal(err)
	}

	if authentication.IsPasswordMatch("UnboxTheCat", []byte("abc")) {
		t.Fatal()
	}

	if !authentication.IsPasswordMatch("UnboxTheCat", []byte("test")) {
		t.Fatal()
	}

	if authentication.IsPasswordMatch("Guest", []byte("test")) {
		t.Fatal()
	}
}

func TestIDValidator(t *testing.T) {
	badID := []string{
		"",        //no name
		" ",       //space
		"(abc",    //bad prefix
		"abc)",    //bad suffix
		"(abc)",   //bad pre/suffix
		"abc(123", //bad character in the middle
	}

	goodID := []string{
		"abc",         //letter only
		"123",         //digit only
		"abc123",      //letter + digit
		"hello_world", //underscore
		"hello-world", //dash
		"UnboxTheCat", //my name
		"__MAIN__",    //underscore pre/suffix
		"-_-",         //-_-
	}

	for _, id := range badID {
		if authentication.IsValidID(id) {
			t.Errorf("failed, expected bad ID: %v", id)
		}
	}

	for _, id := range goodID {
		if !authentication.IsValidID(id) {
			t.Errorf("failed, expected good ID: %v", id)
		}
	}
}
