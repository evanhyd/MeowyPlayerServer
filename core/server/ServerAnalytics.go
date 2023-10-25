package server

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"sync"

	"meowyplayerserver.com/utility/assert"
)

type queryRecord struct {
	query string
	count int
}

type serverAnalytics struct {
	mux     sync.Mutex
	records map[string]int
}

func makeServerAnalytics() serverAnalytics {
	return serverAnalytics{records: make(map[string]int)}
}

func (s *serverAnalytics) log(query string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.records[query]++
}

func (s *serverAnalytics) report() []queryRecord {
	s.mux.Lock()
	defer s.mux.Unlock()

	records := make([]queryRecord, 0, len(s.records))
	for query, count := range s.records {
		records = append(records, queryRecord{query, count})
	}
	return records
}

func (s *serverAnalytics) ServerStats(resp http.ResponseWriter, req *http.Request) {
	s.log(req.URL.Path)

	records := s.report()
	switch req.URL.Query().Get("sort") {
	case "title":
		slices.SortStableFunc(records, func(q1, q2 queryRecord) int {
			return strings.Compare(strings.ToUpper(q1.query), strings.ToUpper(q2.query))
		})

	case "count":
		fallthrough
	default:
		slices.SortStableFunc(records, func(q1, q2 queryRecord) int {
			return q1.count - q2.count
		})
	}

	for _, record := range records {
		_, err := fmt.Fprintf(resp, "%v: %v\n", record.query, record.count)
		assert.NoErr(err, "failed to print analytics stats")
	}
}
