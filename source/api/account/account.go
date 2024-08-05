package account

import (
	"encoding/json"
	"sync"
)

type Account struct {
	username string `json:"username"`
	id       string `json:"id"`
	salt     []byte `json:"salt"`
	hash     []byte `json:"hash"`
}

// username: account
type accountMap struct {
	sync.Map
}

func (m *accountMap) MarshalJSON() ([]byte, error) {
	mp := map[string]Account{}
	m.Range(func(key, value any) bool {
		mp[key.(string)] = value.(Account)
		return true
	})
	return json.Marshal(mp)
}

func (m *accountMap) UnmarshalJSON(data []byte) error {
	mp := map[string]Account{}
	if err := json.Unmarshal(data, &mp); err != nil {
		return err
	}
	for key, val := range mp {
		m.Store(key, val)
	}
	return nil
}
