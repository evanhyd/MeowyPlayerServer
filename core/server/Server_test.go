package server

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"meowyplayerserver.com/core/resource"
)

const (
	kServerUrl   = "http://localhost"
	kTestFile    = "test.txt"
	kTestContent = "abcdefg"
)

func TestCollectionValidator(t *testing.T) {
	s := MakeServer()

	//valid
	if !s.isValidCollection("a.zip") {
		t.Fail()
	}
	if !s.isValidCollection("a.txt") {
		t.Fail()
	}
	if !s.isValidCollection("a.json") {
		t.Fail()
	}

	//invalid
	if s.isValidCollection("") {
		t.Fail()
	}
	if s.isValidCollection(" ") {
		t.Fail()
	}
	if s.isValidCollection(".") {
		t.Fail()
	}
	if s.isValidCollection("..") {
		t.Fail()
	}
	if s.isValidCollection("/a") {
		t.Fail()
	}
	if s.isValidCollection("/a/b") {
		t.Fail()
	}
	if s.isValidCollection("a/") {
		t.Fail()
	}
	if s.isValidCollection("a/b") {
		t.Fail()
	}
	if s.isValidCollection(".zip") {
		t.Fail()
	}
	if s.isValidCollection("a.sh") {
		t.Fail()
	}
	if s.isValidCollection("a.zip.txt") {
		t.Fail()
	}
	if s.isValidCollection("/") {
		t.Fail()
	}
	if s.isValidCollection(`\`) {
		t.Fail()
	}
	if s.isValidCollection("\a") {
		t.Fail()
	}
	if s.isValidCollection("\n") {
		t.Fail()
	}
}

func TestUpload(t *testing.T) {
	resource.MakeNecessaryPath()
	defer os.RemoveAll(resource.CollectionPath())

	go startTestServer()
	if err := uploadTestFile(); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(resource.CollectionPath(), kTestFile))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != kTestContent {
		t.Fatal("content mismatch")
	}
}

func TestDownload(t *testing.T) {
	resource.MakeNecessaryPath()
	defer os.RemoveAll(resource.CollectionPath())

	go startTestServer()
	if err := uploadTestFile(); err != nil {
		t.Fatal(err)
	}
	data, err := downloadTestFile()
	if err != nil {
		t.Fatal(err)
	}
	if string(data) == kTestContent {
		t.Fatal("content mismatch")
	}
}

func startTestServer() {
	s := MakeServer()
	http.HandleFunc("/stats", s.ServerStats)
	http.HandleFunc("/list", s.ServerRequestList)
	http.HandleFunc("/upload", s.ServerRequestUpload)
	http.HandleFunc("/download", s.ServerRequestDownload)
	http.ListenAndServe("localhost:80", nil)
}

func uploadTestFile() error {
	//prepare the fields
	fieldBody := bytes.Buffer{}
	fieldWriter := multipart.NewWriter(&fieldBody)

	//set file
	fieldPart, err := fieldWriter.CreateFormFile("collection", kTestFile)
	if err != nil {
		return err
	}
	if _, err = fieldPart.Write([]byte(kTestContent)); err != nil {
		return err
	}
	fieldWriter.Close()

	//send post
	resp, err := http.Post(kServerUrl+"/upload", fieldWriter.FormDataContentType(), &fieldBody)
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, resp.Body)
	return nil
}

func downloadTestFile() ([]byte, error) {
	resp, err := http.Get(kServerUrl + "/download")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
