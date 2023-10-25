package server

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"meowyplayerserver.com/core/resource"
	"meowyplayerserver.com/utility/assert"
)

type serverManager struct {
}

func makeServerManager() serverManager {
	return serverManager{}
}

func (s *ServerState) ServerList(resp http.ResponseWriter, req *http.Request) {
	s.log(req.URL.Path)

	entries, err := os.ReadDir(resource.UserPath())
	assert.NoErr(err, "failed to read base users path")

	for _, entry := range entries {
		info, err := entry.Info()
		assert.NoErr(err, "failed to read user info")

		_, err = fmt.Fprintln(resp, info.Name(), info.ModTime().Format("2006-01-02 15:04"), info.Size())
		assert.NoErr(err, "failed to print user info")
	}
}

func (s *ServerState) ServerDownload(resp http.ResponseWriter, req *http.Request) {
	s.log(req.URL.Path)

	//sanatize user name
	//should be replaced by a more secure method
	user := req.URL.Query().Get("user")
	if !fs.ValidPath(user) {
		sendError(resp, http.StatusNotFound, fmt.Errorf("invalid user name: %v", user), "failed to fetch user profile")
		return
	}

	//check if user exists
	collectionPath := resource.CollectionPath(user)
	if _, err := os.Stat(collectionPath); err != nil {
		sendError(resp, http.StatusNotFound, fmt.Errorf("invalid user name: %v", user), "failed to fetch user profile")
		return
	}

	//prepare zip
	zipWriter := zip.NewWriter(resp)
	defer zipWriter.Close()

	addToZip := func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		//directory
		if info.IsDir() {
			_, err = zipWriter.Create(path + "/")
			return err
		}

		//file
		fileWriter, err := zipWriter.Create(path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(fileWriter, file)
		return err
	}

	//compress files into the zip buffer
	if err := filepath.WalkDir(collectionPath, addToZip); err != nil {
		sendError(resp, http.StatusNotFound, err, "failed to download "+user+" collection")
		return
	}
}
