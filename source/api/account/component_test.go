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

	//normal registration
	if _, ok := comp.Register("UnboxTheCat", "test"); !ok {
		t.Error("Register() = false, expected true")
	}

	//username already exists
	if _, ok := comp.Register("UnboxTheCat", "test"); ok {
		t.Error("Register() = true, expected false")
	}

	//username is too short
	if _, ok := comp.Register("", "test"); ok {
		t.Error("Register() = true, expected false")
	}

	//username is too long
	if _, ok := comp.Register("0123456789abcdef0123456789abcdef", "test"); ok {
		t.Error("Register() = true, expected false")
	}
}

func TestManager_Authenticate(t *testing.T) {
	comp := makeStubAccountComponent(t)
	if err := comp.Initialize(); err != nil {
		t.Fatal(err)
	}

	//normal registration
	if _, ok := comp.Register("UnboxTheCat", "test"); !ok {
		t.Error("Register() = false, expected true")
	}

	//normal login
	if _, ok := comp.Authenticate("UnboxTheCat", "test"); !ok {
		t.Fatalf("Authenticate() = false, wanted true")
	}

	//incorrect password
	if _, ok := comp.Authenticate("UnboxTheCat", "abcdefghijkl"); ok {
		t.Fatalf("Authenticate() = true, wanted false")
	}

	//login to non-existed user
	if _, ok := comp.Authenticate("not_exist", "abc"); ok {
		t.Fatalf("Authenticate() = true, wanted false")
	}

	//missing username
	if _, ok := comp.Authenticate("", ""); ok {
		t.Fatalf("Authenticate() = true, wanted false")
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
