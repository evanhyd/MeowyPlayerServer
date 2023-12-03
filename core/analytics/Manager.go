package analytics

import (
	"errors"
	"io/fs"
	"os"
	"sync"

	"meowyplayerserver.com/utility/ujson"
)

var instance manager

type manager struct {
	mux     sync.RWMutex   //guard records
	Records map[string]int `json:"records"`
}

type Record struct {
	Action string
	Count  int
}

func save() error {
	return ujson.WriteFile(AnalyticsFile(), instance.Records)
}

func load() error {
	return ujson.ReadFile(AnalyticsFile(), &instance.Records)
}

func Initialize() error {
	_, err := os.Stat(AnalyticsFile())
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if errors.Is(err, fs.ErrNotExist) {
		instance.Records = map[string]int{}
		if err := save(); err != nil {
			return err
		}
	}
	return load()
}

func Log(action string) error {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	instance.Records[action]++
	return save()
}

func Stat() []Record {
	instance.mux.RLock()
	defer instance.mux.RUnlock()

	records := make([]Record, 0, len(instance.Records))
	for command, count := range instance.Records {
		records = append(records, Record{command, count})
	}
	return records
}
