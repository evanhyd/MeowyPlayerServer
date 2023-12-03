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

	if !(list[0].Title == "UnboxTheCat.zip" && list[1].Title == "Guest.zip" || list[1].Title == "UnboxTheCat.zip" && list[0].Title == "Guest.zip") {
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
		"a.zip",
		"123.zip",
		"valid_name.zip",
		"valid-name.zip",
		"CON.zip", // Windows reserved filename but we are on Linux
	}

	invalidName := []string{
		"",                           // Empty filename
		"   ",                        // Space filename
		".",                          // Not filename
		"..",                         // Not filename
		".zip",                       //Empty filename
		"Hello World.zip",            //Filename contains space
		".hiddenfile.zip",            // Hidden files
		"../log.zip",                 // Contains path elements
		"/path/to/data.zip",          // Contains path elements
		"/test.zip",                  // Contains absolute path
		"report.pdf",                 // Wrong file extension
		"backup",                     // No file extension
		"folder\\image.zip",          // Contains path elements (Windows path)
		"notes.zip/",                 // Trailing slash indicates a directory
		"../folder.zip/",             // Contains path elements and trailing slash
		"../../../../etc/passwd.zip", // Directory traversal attack
		"my\000file.zip",             // Null byte injection
		"file.zip; rm -rf /",         // Command injection
		"file|.zip",                  // Pipe character could be used in command injection
		"file.zip && rm -rf /",       // Command chaining
		"evil.jpg\x00.zip",           // Null byte to hide true extension
		"file.zip .",                 // Space and trailing dot
		"file.zip..",                 // Trailing dots
		"file.zip--",                 // Double dash
		"symlink.zip->/etc/passwd",   // Creating a symlink
		"\u202Efilezip.zip",          // Unicode right-to-left override character
		"file.zip /",                 // Trailing space and slash
		"file.zip#",                  // Hash character could be used in scripts or URLs
		"image.svg\x00.zip",          // Null byte to fake extension
		"file.zip/../malicious.zip",  // Directory traversal after valid extension
	}

	for _, name := range validName {
		if !collection.IsValidFileName(name) {
			t.Fatalf("failed, expected valid filename: %v", name)
		}
	}

	for _, name := range invalidName {
		if collection.IsValidFileName(name) {
			t.Fatalf("failed, expected invalid filename: %v", name)
		}
	}
}
