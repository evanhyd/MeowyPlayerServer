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

	if !comp.Register("UnboxTheCat", "test") {
		t.Error("Register() = false, expected true")
	}

	//username already exists
	if comp.Register("UnboxTheCat", "test") {
		t.Error("Register() = true, expected false")
	}

	//username is too short
	if comp.Register("", "test") {
		t.Error("Register() = true, expected false")
	}

	//username is too long
	if comp.Register("0123456789abcdef0123456789abcdef", "test") {
		t.Error("Register() = true, expected false")
	}
}

func TestManager_Authorize(t *testing.T) {
	comp := makeStubAccountComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	if !comp.Register("UnboxTheCat", "test") {
		t.Error("Register() = false, expected true")
	}

	//normal login
	if !comp.Authorize("UnboxTheCat", "test") {
		t.Fatalf("Authorize() = false, wanted true")
	}

	//incorrect password
	if comp.Authorize("UnboxTheCat", "abcdefghijkl") {
		t.Fatalf("Authorize() = true, wanted false")
	}

	//login to non-existed user
	if comp.Authorize("not_exist", "abc") {
		t.Fatalf("Authorize() = true, wanted false")
	}

	//missing username
	if comp.Authorize("", "") {
		t.Fatalf("Authorize() = true, wanted false")
	}
}

func TestManager_IsValidUsername(t *testing.T) {
	comp := makeStubAccountComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	badUsername := []string{
		"",                               //too short
		"012345678901234567890123456789", //too long
	}

	goodUsername := []string{
		"abc123",      //average username
		"UnboxTheCat", //my name
		"__MAIN__",    //underscore pre/suffix
		"用户名",         //utf8
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
