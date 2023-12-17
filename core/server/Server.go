package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"meowyplayerserver.com/core/analytics"
	"meowyplayerserver.com/core/authentication"
	"meowyplayerserver.com/core/collection"
)

var instance = Server{}

func Instance() *Server {
	return &instance
}

type Server struct {
}

func (s *Server) ServerStats(resp http.ResponseWriter, req *http.Request) {
	analytics.Log(req.URL.Path)

	records := analytics.Stat()
	switch req.URL.Query().Get("sort") {
	case "title":
		slices.SortStableFunc(records, func(r1, r2 analytics.Record) int {
			return strings.Compare(strings.ToUpper(r1.Action), strings.ToUpper(r2.Action))
		})

	case "count":
		fallthrough
	default:
		slices.SortStableFunc(records, func(r1, r2 analytics.Record) int {
			return r2.Count - r1.Count
		})
	}

	for _, record := range records {
		fmt.Fprintf(resp, "%v: %v\n", record.Action, record.Count)
	}
}

func (s *Server) ServerList(resp http.ResponseWriter, req *http.Request) {
	analytics.Log(req.URL.Path)

	infos, err := collection.List()
	if err != nil {
		sendError(resp, http.StatusInternalServerError, err.Error())
		return
	}
	json.NewEncoder(resp).Encode(infos)
}

func (s *Server) ServerRegister(resp http.ResponseWriter, req *http.Request) {
	analytics.Log(req.URL.Path)

	username, password, _ := req.BasicAuth()
	if err := authentication.RegisterAccount(username, password); err != nil {
		sendError(resp, http.StatusNotFound, err.Error())
		return
	}
}

func (s *Server) ServerUpload(resp http.ResponseWriter, req *http.Request) {
	analytics.Log(req.URL.Path)

	username, password, _ := req.BasicAuth()
	if !authentication.IsGoodAuth(username, password) {
		sendError(resp, http.StatusNotFound, "invalid username or password")
		return
	}

	files, _, err := req.FormFile("collection")
	if err != nil {
		sendError(resp, http.StatusBadRequest, err.Error()+" - missing collection file")
		return
	}
	defer files.Close()

	if err := collection.Update(files, username); err != nil {
		sendError(resp, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *Server) ServerDownload(resp http.ResponseWriter, req *http.Request) {
	analytics.Log(req.URL.Path)

	username, password, _ := req.BasicAuth()
	if !authentication.IsGoodAuth(username, password) {
		sendError(resp, http.StatusNotFound, "invalid username or password")
		return
	}

	collectionID := req.URL.Query().Get("collection")
	if !authentication.IsUserExist(collectionID) {
		sendError(resp, http.StatusNotFound, "invalid collection id")
		return
	}

	if err := collection.Fetch(resp, collectionID); err != nil {
		sendError(resp, http.StatusInternalServerError, err.Error())
		return
	}
}
