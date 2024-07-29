package user

import (
	"encoding/json"
	"sync"
)

const ()

type account struct {
	username string `json:"username"`
	id       string `json:"id"`
	salt     []byte `json:"salt"`
	hash     []byte `json:"hash"`
}

type accountMap struct {
	sync.Map
}

func (m *accountMap) MarshalJSON() ([]byte, error) {
	mp := map[string]account{}
	m.Range(func(key, value any) bool {
		mp[key.(string)] = value.(account)
		return true
	})
	return json.Marshal(mp)
}

func (m *accountMap) UnmarshalJSON(data []byte) error {
	mp := map[string]account{}
	if err := json.Unmarshal(data, &mp); err != nil {
		return err
	}
	for key, val := range mp {
		m.Store(key, val)
	}
	return nil
}
