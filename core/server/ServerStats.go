package server

import (
	"sync"

	"golang.org/x/exp/maps"
)

type ServerAnalytics struct {
	mux           sync.Mutex
	requestRecord map[string]int
}

func makeServerStats() ServerAnalytics {
	return ServerAnalytics{requestRecord: make(map[string]int)}
}

func (s *ServerAnalytics) record(request string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.requestRecord[request]++
}

func (s *ServerAnalytics) report() ([]string, []int) {
	s.mux.Lock()
	defer s.mux.Unlock()
	requests := maps.Keys(s.requestRecord)
	counts := maps.Values(s.requestRecord)
	return requests, counts
}
