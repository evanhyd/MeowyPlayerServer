package account

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
)

type UserID struct {
	id uuid.UUID
}

var _ json.Marshaler = UserID{}
var _ json.Unmarshaler = &UserID{}

func NewUserID() UserID         { return UserID{uuid.New()} }
func (k UserID) String() string { return k.id.String() }

func (k UserID) MarshalJSON() ([]byte, error) {
	return json.Marshal(k.String())
}

func (k *UserID) UnmarshalJSON(d []byte) error {
	var err error
	k.id, err = uuid.ParseBytes(d)
	return err
}

type Account struct {
	Username string `json:"username"`
	UserID   UserID `json:"userID"`
	Salt     []byte `json:"salt"`
	Hash     []byte `json:"hash"`
}

// username: Account
type accountMap struct {
	sync.Map
}

var _ json.Marshaler = accountMap{}
var _ json.Unmarshaler = &accountMap{}

func (m accountMap) MarshalJSON() ([]byte, error) {
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
