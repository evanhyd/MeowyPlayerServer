package authentication

import (
	"os"
	"testing"
)

func TestRegister(t *testing.T) {
	Initialize()
	defer os.RemoveAll(AccountFile())

	if err := RegisterAccount("UnboxTheCat", "test"); err != nil {
		t.Fatal(err)
	}

	if err := RegisterAccount("UnboxTheCat", "test"); err == nil {
		t.Fatal(err)
	}

	if !IsUserExist("UnboxTheCat") {
		t.Fatal()
	}

	if IsUserExist("Guest") {
		t.Fatal()
	}

	if IsUserExist("") {
		t.Fatal()
	}
}

func TestPassword(t *testing.T) {
	Initialize()
	defer os.RemoveAll(AccountFile())

	if err := RegisterAccount("UnboxTheCat", "test"); err != nil {
		t.Fatal(err)
	}

	if IsGoodAuth("UnboxTheCat", "abc") {
		t.Fatal()
	}

	if !IsGoodAuth("UnboxTheCat", "test") {
		t.Fatal()
	}

	if IsGoodAuth("Guest", "test") {
		t.Fatal()
	}
}

func TestUsernameValidator(t *testing.T) {
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
		if isUserValid(username) {
			t.Errorf("isValidUser(%v) = true, wanted false", username)
		}
	}

	for _, username := range goodUsername {
		if !isUserValid(username) {
			t.Errorf("isValidUser(%v) = false, wanted true", username)
		}
	}
}
