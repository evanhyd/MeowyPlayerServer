package collection_test

import (
	"bytes"
	"os"
	"testing"

	"meowyplayerserver.com/core/collection"
)

func TestList(t *testing.T) {
	collection.Initialize()
	defer os.RemoveAll(collection.CollectionPath())

	if err := os.WriteFile(collection.CollectionFile("UnboxTheCat"), []byte("test"), 0777); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(collection.CollectionFile("Guest"), []byte("test again"), 0777); err != nil {
		t.Fatal(err)
	}

	list, err := collection.List()
	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 2 {
		t.Fatal()
	}

	if !(list[0].Title == "UnboxTheCat" && list[1].Title == "Guest" || list[1].Title == "UnboxTheCat" && list[0].Title == "Guest") {
		t.Fatal()
	}
}

func TestUpdate(t *testing.T) {
	collection.Initialize()
	defer os.RemoveAll(collection.CollectionPath())

	if err := collection.Update(bytes.NewReader([]byte("test")), "UnboxTheCat"); err != nil {
		t.Fatal(err)
	}

	list, err := collection.List()
	if err != nil {
		t.Fatal(err)
	}

	if len(list) != 1 {
		t.Fatal()
	}
}

func TestFetch(t *testing.T) {
	collection.Initialize()
	defer os.RemoveAll(collection.CollectionPath())

	if err := os.WriteFile(collection.CollectionFile("UnboxTheCat"), []byte("test"), 0777); err != nil {
		t.Fatal(err)
	}

	buffer := bytes.Buffer{}
	if err := collection.Fetch(&buffer, "UnboxTheCat"); err != nil {
		t.Fatal(err)
	}

	if buffer.String() != "test" {
		t.Fatal()
	}
}

func TestFileNameValidator(t *testing.T) {
	validName := []string{
		"a",
		"123",
		"valid_name",
		"valid-name",
		"CON", // Windows reserved filename but we are on Linux
	}

	invalidName := []string{
		"",                       // Empty filename
		"   ",                    // Space filename
		".",                      // Not filename
		"..",                     // Not filename
		"",                       //Empty filename
		"Hello World",            //Filename contains space
		".hiddenfile",            // Hidden files
		"../log",                 // Contains path elements
		"/path/to/data",          // Contains path elements
		"/test",                  // Contains absolute path
		"report.pdf",             // Wrong file extension
		"backup.zip",             // Zip file extension
		"folder\\image",          // Contains path elements (Windows path)
		"notes/",                 // Trailing slash indicates a directory
		"../folder/",             // Contains path elements and trailing slash
		"../../../../etc/passwd", // Directory traversal attack
		"my\000file",             // Null byte injection
		"file; rm -rf /",         // Command injection
		"file|",                  // Pipe character could be used in command injection
		"file && rm -rf /",       // Command chaining
		"evil.jpg\x00",           // Null byte to hide true extension
		"file .",                 // Space and trailing dot
		"file..",                 // Trailing dots
		"symlink->/etc/passwd",   // Creating a symlink
		"\u202Efilezip",          // Unicode right-to-left override character
		"file /",                 // Trailing space and slash
		"file#",                  // Hash character could be used in scripts or URLs
		"image.svg\x00",          // Null byte to fake extension
		"file/../malicious",      // Directory traversal after valid extension
	}

	for _, name := range validName {
		if !collection.IsValidFileName(name) {
			t.Fatalf("IsValidFileName(%v)=false, wanted true", name)
		}
	}

	for _, name := range invalidName {
		if collection.IsValidFileName(name) {
			t.Fatalf("IsValidFileName(%v)=true, wanted false", name)
		}
	}
}
