package server

import (
	"archive/zip"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"meowyplayerserver.com/core/resource"
	"meowyplayerserver.com/utility/uzip"
)

type serverManager struct {
}

func makeServerManager() serverManager {
	return serverManager{}
}

func (s *ServerState) isValidUser(user string) bool {
	_, err := os.Stat(resource.CollectionPath(user))
	return fs.ValidPath(user) && err == nil
}

func (s *ServerState) ServerUsers(resp http.ResponseWriter, req *http.Request) {
	s.log(req.URL.Path)

	entries, err := os.ReadDir(resource.UserPath())
	if err != nil {
		sendError(resp, http.StatusInternalServerError, fmt.Sprintf("%v: %v", err, "failed to read user path"))
		return
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			sendError(resp, http.StatusInternalServerError, fmt.Sprintf("%v: %v", err, "failed to read user info"))
			return
		}
		fmt.Fprintln(resp, info.Name(), info.ModTime().Format("2006-01-02 15:04"), info.Size())
	}
}

func (s *ServerState) ServerDownload(resp http.ResponseWriter, req *http.Request) {
	s.log(req.URL.Path)

	user := req.URL.Query().Get("user")
	if !s.isValidUser(user) {
		sendError(resp, http.StatusNotFound, fmt.Sprintf("invalid user id: %v", user))
		return
	}

	if err := uzip.Compress(resp, resource.CollectionPath(user)); err != nil {
		sendError(resp, http.StatusInternalServerError, fmt.Sprintf("failed to compress: %v", user))
	}
}

func (s *ServerState) ServerUpload(resp http.ResponseWriter, req *http.Request) {
	s.log(req.URL.Path)

	user := req.PostFormValue("user")
	if !s.isValidUser(user) {
		sendError(resp, http.StatusNotFound, fmt.Sprintf("invalid user id: %v", user))
		return
	}
	collectionPath := resource.CollectionPath(user)

	//reset user's collection directory
	if err := os.RemoveAll(collectionPath); err != nil {
		sendError(resp, http.StatusInternalServerError, fmt.Sprintf("failed to upload: %v", user))
		return
	}

	if err := os.MkdirAll(collectionPath, 0777); err != nil {
		sendError(resp, http.StatusInternalServerError, fmt.Sprintf("failed to upload: %v", user))
		return
	}

	//extract the file
	files, fileHeaders, err := req.FormFile("collection")
	if err != nil {
		sendError(resp, http.StatusNotFound, fmt.Sprintf("failed to parse the collection file: %v", fileHeaders))
		return
	}
	defer files.Close()

	zipHandle, err := zip.NewReader(files, fileHeaders.Size)
	if err != nil {
		sendError(resp, http.StatusInternalServerError, fmt.Sprintf("failed to upload: %v", user))
		return
	}

	if err := uzip.Extract(collectionPath, zipHandle); err != nil {
		sendError(resp, http.StatusInternalServerError, fmt.Sprintf("failed to upload: %v", user))
		return
	}
}
