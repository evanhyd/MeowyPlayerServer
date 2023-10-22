package server

import (
	"bytes"
	"fmt"
	"net/http"
	"os"

	"meowyplayerserver.com/core/resource"
	"meowyplayerserver.com/utility/assert"
)

var state = makeServerState()

func GetInstance() *ServerState {
	return &state
}

type ServerState struct {
	ServerAnalytics
}

func makeServerState() ServerState {
	return ServerState{makeServerStats()}
}

func (s *ServerState) ServerStats(resp http.ResponseWriter, req *http.Request) {
	s.record(req.URL.Path)

	buffer := bytes.Buffer{}
	requests, counts := s.report()
	for i := range requests {
		fmt.Fprintf(&buffer, "%v: %v\n", requests[i], counts[i])
	}

	_, err := resp.Write(buffer.Bytes())
	assert.NoErr(err, fmt.Sprintf("failed to respond to %v", req.URL.Path))
}

func (s *ServerState) ServerList(resp http.ResponseWriter, req *http.Request) {
	s.record(req.URL.Path)

	entries, err := os.ReadDir("core/resource")
	assert.NoErr(err, "failed to list all the collections")

	buffer := bytes.Buffer{}
	for _, entry := range entries {
		info, err := entry.Info()
		assert.NoErr(err, "failed to read collection info")
		fmt.Fprintln(&buffer, info.Name(), info.ModTime().Format("2006-01-02 15:04"), info.Size())
	}

	_, err = resp.Write(buffer.Bytes())
	assert.NoErr(err, fmt.Sprintf("failed to respond to %v", req.URL.Path))
}

func (s *ServerState) ServerCollection(resp http.ResponseWriter, req *http.Request) {
	s.record(req.URL.Path)

	collectionName := req.URL.Query().Get("user")
	data, err := os.ReadFile(resource.CollectionPath(collectionName)) //seems dangerous? can be abused with name="../fileName"
	assert.NoErr(err, fmt.Sprintf("failed to find collection: %v", collectionName))

	_, err = resp.Write(data)
	assert.NoErr(err, fmt.Sprintf("failed to respond to %v", req.URL.Path))
}
