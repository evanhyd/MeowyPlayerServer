package api

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"sync/atomic"
)

type APIManager struct {
	tracker    map[string]*atomic.Int32
	logsBuffer circularBuffer
}

func NewAPIManager() *APIManager {
	const kBufferSize = 1 << 17 //128 KB
	m := APIManager{
		tracker:    map[string]*atomic.Int32{},
		logsBuffer: makeCircularBuffer(kBufferSize),
	}
	log.SetOutput(&m.logsBuffer)
	m.RegisterAPI("/stats", m.statsHandler)
	m.RegisterAPI("/logs", m.logsHandler)
	return &m
}
func (m *APIManager) RegisterAPI(pattern string, handler http.HandlerFunc) {
	m.tracker[pattern] = &atomic.Int32{}
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		m.tracker[pattern].Add(1)
		handler(w, r)
	})
}

func (m *APIManager) statsHandler(w http.ResponseWriter, r *http.Request) {
	type Record struct {
		Pattern string
		Count   int32
	}

	records := make([]Record, 0, len(m.tracker))
	for pattern, cnt := range m.tracker {
		records = append(records, Record{Pattern: pattern, Count: cnt.Load()})
	}
	slices.SortFunc(records, func(l Record, r Record) int {
		return strings.Compare(l.Pattern, r.Pattern)
	})

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(records); err != nil {
		log.Println(err)
	}
}

func (m *APIManager) logsHandler(w http.ResponseWriter, _ *http.Request) {
	n, err := m.logsBuffer.WriteTo(w)
	if err != nil {
		log.Printf("[logsHandler] %v bytes sent, %v\n", n, err)
	}
}
