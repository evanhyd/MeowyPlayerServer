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

var state = Server{}

func GetInstance() *Server {
	return &state
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

	id := req.PostFormValue("id")
	password := req.PostFormValue("password")
	if id == "" || authentication.IsAccountExist(id) {
		sendError(resp, http.StatusBadRequest, "user ID is invalid or has been registered")
		return
	}

	if err := authentication.RegisterAccount(id, []byte(password)); err != nil {
		sendError(resp, http.StatusBadRequest, err.Error())
		return
	}
}

func (s *Server) ServerUpload(resp http.ResponseWriter, req *http.Request) {
	analytics.Log(req.URL.Path)

	id := req.PostFormValue("id")
	if !authentication.IsAccountExist(id) {
		sendError(resp, http.StatusNotFound, "user id does not exist")
		return
	}

	files, _, err := req.FormFile("collection")
	if err != nil {
		sendError(resp, http.StatusBadRequest, err.Error()+" - missing collection file")
		return
	}
	defer files.Close()

	if err := collection.Update(files, id); err != nil {
		sendError(resp, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *Server) ServerDownload(resp http.ResponseWriter, req *http.Request) {
	analytics.Log(req.URL.Path)

	id := req.URL.Query().Get("id")
	if !authentication.IsAccountExist(id) {
		sendError(resp, http.StatusNotFound, "user id does not exist")
		return
	}

	if err := collection.Fetch(resp, id); err != nil {
		sendError(resp, http.StatusInternalServerError, err.Error())
		return
	}
}
