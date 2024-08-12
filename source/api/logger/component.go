package logger

import (
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"sync/atomic"
)

type Record struct {
	Pattern string
	Count   int32
}

type Component struct {
	tracker   map[string]*atomic.Int32
	logBuffer circularBuffer
}

func MakeComponent() Component {
	const kBufferSize = 1 << 17 //128 KB
	return Component{
		tracker:   map[string]*atomic.Int32{},
		logBuffer: makeCircularBuffer(kBufferSize),
	}
}

func (c *Component) Initialize() error {
	log.SetOutput(&c.logBuffer)
	return nil
}

func (c *Component) GetRecords() []Record {
	records := make([]Record, 0, len(c.tracker))
	for pattern, cnt := range c.tracker {
		records = append(records, Record{Pattern: pattern, Count: cnt.Load()})
	}
	slices.SortFunc(records, func(l Record, r Record) int {
		return strings.Compare(l.Pattern, r.Pattern)
	})
	return records
}

func (c *Component) DumpLog(w io.Writer) (int64, error) {
	return c.logBuffer.WriteTo(w)
}

func (c *Component) RegisterAPI(pattern string, handler http.HandlerFunc) {
	c.tracker[pattern] = &atomic.Int32{}
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		c.tracker[pattern].Add(1)
		handler(w, r)
	})
}
