package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"meowyplayerserver.com/core/resource"
)

type serverManager struct {
	access    sync.Map //collection, *RWMutex
	validator *regexp.Regexp
}

func makeServerManager() serverManager {
	const pattern = `^(\w+\.(zip|txt|json))$`
	return serverManager{validator: regexp.MustCompile(pattern)}
}

func (s *Server) isValidCollection(collection string) bool {
	return s.validator.MatchString(collection)
}

func (s *Server) getMux(collection string) *sync.RWMutex {
	mux, _ := s.access.LoadOrStore(collection, &sync.RWMutex{})
	return mux.(*sync.RWMutex)
}

func (s *Server) ServerRequestList(resp http.ResponseWriter, req *http.Request) {
	s.log(req.URL.Path)

	entries, err := os.ReadDir(resource.CollectionPath())
	if err != nil {
		sendError(resp, http.StatusInternalServerError, err.Error()+" - collection path doesn't exist")
		return
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			sendError(resp, http.StatusInternalServerError, err.Error()+" - failed to read entry info: "+entry.Name())
			return
		}
		fmt.Fprintln(resp, info.Name(), info.ModTime().Format("2006-01-02 15:04"), info.Size())
	}
}

func (s *Server) ServerRequestUpload(resp http.ResponseWriter, req *http.Request) {
	s.log(req.URL.Path)

	files, fileHeaders, err := req.FormFile("collection")
	if err != nil {
		sendError(resp, http.StatusBadRequest, err.Error()+" - failed to upload the collection")
		return
	}
	defer files.Close()

	mux := s.getMux(fileHeaders.Filename)
	mux.Lock()
	defer mux.Unlock()

	file, err := os.OpenFile(filepath.Join(resource.CollectionPath(), fileHeaders.Filename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
	if err != nil {
		sendError(resp, http.StatusInternalServerError, err.Error()+" - failed to create the collection: "+fileHeaders.Filename)
		return
	}
	defer file.Close()
	io.Copy(file, files)
}

func (s *Server) ServerRequestDownload(resp http.ResponseWriter, req *http.Request) {
	s.log(req.URL.Path)

	collection := req.URL.Query().Get("collection")
	if !s.isValidCollection(collection) {
		sendError(resp, http.StatusNotFound, "invalid collection: "+collection)
		return
	}

	mux := s.getMux(collection)
	mux.RLock()
	defer mux.RUnlock()

	file, err := os.Open(filepath.Join(resource.CollectionPath(), collection))
	if err != nil {
		sendError(resp, http.StatusInternalServerError, err.Error()+" - failed to download the collection: "+collection)
		return
	}
	defer file.Close()
	io.Copy(resp, file)
}

func (s *Server) ServerRequestRemove(resp http.ResponseWriter, req *http.Request) {
	s.log(req.URL.Path)

	collection := req.PostFormValue("collection")
	if !s.isValidCollection(collection) {
		sendError(resp, http.StatusNotFound, "invalid collection: "+collection)
		return
	}

	mux := s.getMux(collection)
	mux.Lock()
	defer mux.Unlock()

	if err := os.RemoveAll(filepath.Join(resource.CollectionPath(), collection)); err != nil {
		sendError(resp, http.StatusInternalServerError, err.Error()+" - failed to remove the collection: "+collection)
	}
}
