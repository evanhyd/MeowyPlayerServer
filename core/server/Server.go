package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"meowyplayerserver.com/core/analytics"
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

	records := analytics.Read()
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

func (s *Server) ServerUpload(resp http.ResponseWriter, req *http.Request) {
	analytics.Log(req.URL.Path)

	//to do: verify id

	files, _, err := req.FormFile("collection")
	if err != nil {
		sendError(resp, http.StatusBadRequest, err.Error()+" - collection file is missing from the POST request")
		return
	}
	defer files.Close()

	if err := collection.Update(files, "Guest"); err != nil {
		sendError(resp, http.StatusInternalServerError, err.Error())
		return
	}
}

func (s *Server) ServerDownload(resp http.ResponseWriter, req *http.Request) {
	analytics.Log(req.URL.Path)

	//to do: verify id

	if err := collection.Fetch(resp, "Guest"); err != nil {
		sendError(resp, http.StatusInternalServerError, err.Error())
		return
	}
}
