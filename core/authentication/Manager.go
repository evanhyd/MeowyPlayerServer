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

func RegisterAccount(username string, password string) error {
	if !isUserValid(username) || IsUserExist(username) {
		return fmt.Errorf("username %v is invalid or registered", username)
	}

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	instance.mux.Lock()
	instance.accounts[username] = account{username, salt, computeHash([]byte(password), salt)}
	instance.mux.Unlock()

	return save()
}

func isUserValid(username string) bool {
	return usernameValidator.MatchString(username)
}

func IsUserExist(username string) bool {
	instance.mux.RLock()
	_, ok := instance.accounts[username]
	instance.mux.RUnlock()
	return ok
}

func IsGoodAuth(username string, password string) bool {
	instance.mux.RLock()
	account, ok := instance.accounts[username]
	instance.mux.RUnlock()

	if !ok {
		return false
	}
	return slices.Equal(account.Hash, computeHash([]byte(password), account.Salt))
}

func computeHash(password []byte, salt []byte) []byte {
	hash := sha256.Sum256(append(password, salt...))
	return hash[:]
}
