package account

import (
	"testing"
)

func makeStubAccountComponent(t *testing.T) Component {
	acc := MakeComponent()
	acc.storage = makeStubAccountStorage(t)
	return acc
}

func TestManager_Register(t *testing.T) {
	comp := makeStubAccountComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	if err := comp.Register("UnboxTheCat", "test"); err != nil {
		t.Fatal(err)
	}

	if err := comp.Register("UnboxTheCat", "test"); err == nil {
		t.Fatal(err)
	}

	if err := comp.Register("", "test"); err == nil {
		t.Fatal(err)
	}
}

func TestManager_Authorize(t *testing.T) {
	comp := makeStubAccountComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	if err := comp.Register("UnboxTheCat", "test"); err != nil {
		t.Fatal(err)
	}

	if !comp.Authorize("UnboxTheCat", "test") {
		t.Fatalf("authorize() = false, wanted true")
	}

	if comp.Authorize("UnboxTheCat", "test1") {
		t.Fatalf("authorize() = true, wanted false")
	}

	if comp.Authorize("not_exist", "abc") {
		t.Fatalf("authorize() = true, wanted false")
	}

	if comp.Authorize("", "") {
		t.Fatalf("authorize() = true, wanted false")
	}
}

func TestManager_IsValidUsername(t *testing.T) {
	comp := makeStubAccountComponent(t)
	if err := comp.Initialize(); err != nil {
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
		if comp.isValidUsername(username) {
			t.Errorf("isValidUser(%v) = true, wanted false", username)
		}
	}

	for _, username := range goodUsername {
		if !comp.isValidUsername(username) {
			t.Errorf("isValidUser(%v) = false, wanted true", username)
		}
	}
}
