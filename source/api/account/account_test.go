package account

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"
)

func TestAccountMapMarshaler(t *testing.T) {
	accs := []Account{
		{Username: "a0", UserID: NewUserID(), Salt: []byte("b0"), Hash: []byte("c0")},
		{Username: "a1", UserID: NewUserID(), Salt: []byte("b1"), Hash: []byte("c1")},
		{Username: "a2", UserID: NewUserID(), Salt: []byte("b2"), Hash: []byte("c2")},
		{Username: "a3", UserID: NewUserID(), Salt: []byte("b3"), Hash: []byte("c3")},
	}

	mp := accountMap{}
	for _, acc := range accs {
		mp.Store(acc.Username, acc)
	}

	//marshal and unmarshal
	data, err := json.Marshal(mp)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(data))
	mp2 := accountMap{}
	if err := json.Unmarshal(data, &mp2); err != nil {
		t.Fatal(err)
	}

	//check equality
	for _, acc := range accs {
		val, ok := mp.Load(acc.Username)
		if !ok {
			t.Errorf("Load(%v) = false, expected true", acc.Username)
		}

		acc1, ok := val.(Account)
		if !ok {
			t.Fatalf("Load(%v) returns non Account type", val)
		}

		if acc.Username != acc1.Username || acc.UserID != acc1.UserID || !slices.Equal(acc.Salt, acc1.Salt) || !slices.Equal(acc.Hash, acc1.Hash) {
			t.Errorf("acc != acc1, %v %v", acc, acc1)
		}
	}
}
