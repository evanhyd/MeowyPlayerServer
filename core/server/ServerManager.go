package server

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
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

	infos := make([]resource.CollectionInfo, 0, len(entries))
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			sendError(resp, http.StatusInternalServerError, err.Error()+" - failed to read entry info: "+entry.Name())
			return
		}
		infos = append(infos, resource.CollectionInfo{Title: info.Name(), Date: info.ModTime(), Size: info.Size()})
	}
	json.NewEncoder(resp).Encode(infos)
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

	file, err := os.OpenFile(resource.CollectionFile(fileHeaders.Filename), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
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

	file, err := os.Open(resource.CollectionFile(collection))
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

	if err := os.RemoveAll(resource.CollectionFile(collection)); err != nil {
		sendError(resp, http.StatusInternalServerError, err.Error()+" - failed to remove the collection: "+collection)
	}
}
