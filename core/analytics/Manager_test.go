package analytics_test

import (
	"os"
	"slices"
	"testing"

	"meowyplayerserver.com/core/analytics"
)

func TestLog(t *testing.T) {
	analytics.Initialize()
	defer os.RemoveAll(analytics.AnalyticsFile())

	//
	records := analytics.Stat()
	if len(records) > 0 {
		t.Fatal()
	}

	// hello: 1
	analytics.Log("hello")
	records = analytics.Stat()
	if len(records) != 1 {
		t.Fatal()
	}
	if slices.Index(records, analytics.Record{"hello", 1}) == -1 {
		t.Fatal()
	}

	// hello: 3
	analytics.Log("hello")
	analytics.Log("hello")
	records = analytics.Stat()
	if len(records) != 1 {
		t.Fatal()
	}
	if slices.Index(records, analytics.Record{"hello", 3}) == -1 {
		t.Fatal()
	}

	// hello: 3
	// world: 2
	analytics.Log("world")
	analytics.Log("world")
	records = analytics.Stat()
	if len(records) != 2 {
		t.Fatal()
	}
	if slices.Index(records, analytics.Record{"hello", 3}) == -1 {
		t.Fatal()
	}
	if slices.Index(records, analytics.Record{"world", 2}) == -1 {
		t.Fatal()
	}
}
