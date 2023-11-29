package authentication

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"sync"

	"meowyplayerserver.com/utility/ujson"
)

var instance manager

type manager struct {
	accounts map[string]account //id -> account
	mux      sync.RWMutex       //guard accounts
}

func save() error {
	return ujson.WriteFile(AccountFile(), instance.accounts)
}

func load() error {
	return ujson.ReadFile(AccountFile(), &instance.accounts)
}

func Initialize() error {
	_, err := os.Stat(AccountFile())
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return err
	}

	if errors.Is(err, fs.ErrNotExist) {
		instance.accounts = map[string]account{}
		if err := save(); err != nil {
			return err
		}
	}
	return load()
}

func RegisterAccount(id string, plainPassword []byte) error {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	if _, ok := instance.accounts[id]; ok {
		return fmt.Errorf("user ID %v is already registered", id)
	}

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	instance.accounts[id] = account{id, salt, computeHash(plainPassword, salt)}
	return save()
}

func IsValidID(id string) bool {
	return idValidator.MatchString(id)
}

func IsAccountExist(id string) bool {
	instance.mux.RLock()
	defer instance.mux.RUnlock()

	_, ok := instance.accounts[id]
	return ok
}

func IsPasswordMatch(id string, plainPassword []byte) bool {
	instance.mux.RLock()
	defer instance.mux.RUnlock()

	account, ok := instance.accounts[id]
	if !ok {
		return false
	}
	return slices.Equal(account.Hash, computeHash(plainPassword, account.Salt))
}

func computeHash(plainPassword []byte, salt []byte) []byte {
	hash := sha256.Sum256(append(plainPassword, salt...))
	return hash[:]
}
